package main

import (
	"fmt"
	"os"

	"github.com/cocooma/consul-cleaner/awsdiscovery"

	flag "github.com/docker/docker/pkg/mflag"
	"github.com/hashicorp/consul/api"
)

var (
	str, url, port, srvState                                                                          string
	showMember, showMmemberStatus, showChecks, showSrv, deregisterSrv, listSrvInState, showNodeStatus bool
	nsc                                                                                               int
)

// var wg sync.WaitGroup

func listmembers(consulClient *api.Client) {
	members, _ := consulClient.Agent().Members(false)
	for _, server := range members {
		fmt.Println(server.Name)
	}
}

func shownodestaus(consulClient *api.Client) {
	members, _ := consulClient.Agent().Members(false)
	for _, server := range members {
		fmt.Printf("%s %v ", server.Name, server.Status)
	}
}

func forceLeaveBadNode(consulClient *api.Client, nodeStatusCode int) {
	members, _ := consulClient.Agent().Members(false)
	for _, server := range members {
		if server.Status == nodeStatusCode {
			err := consulClient.Agent().ForceLeave(server.Name)
			if err != nil {
				panic(err)
			}
		}
	}
}

// func showmembers(members []*api.AgentMember) {
// 	for _, server := range members {
// 		fmt.Println(server.Status)
// 	}
// }

func listChecks(consulClient *api.Client) {
	checks, _ := consulClient.Agent().Checks()
	for _, check := range checks {
		fmt.Println(check.Node, check.Name, check.Status)
	}
}

func listServices(consulClient *api.Client) {
	services, _ := consulClient.Agent().Services()
	for _, service := range services {
		fmt.Println(service.ID, service.Service, service.Tags)
	}
}

func serviceNameServiceID(connection *api.Client, serviceCheckStatus string) map[string]string {
	services := map[string]string{}
	serv, _, _ := connection.Health().State(serviceCheckStatus, nil)
	for _, key := range serv {
		services[key.ServiceName] = key.ServiceID
	}
	return services
}

func connection(uurl, pport string) *api.Client {
	connection, err := api.NewClient(&api.Config{Address: uurl + ":" + pport})
	if err != nil {
		panic(err)
	}
	return connection
}

func listServicesInState(consulConnection *api.Client, serviceCheckStatus string) {
	service := serviceNameServiceID(consulConnection, serviceCheckStatus)
	for serviceName, serviceID := range service {
		fmt.Println(serviceName + " " + serviceID)
	}
}

func deregisterService(consulConnection *api.Client, serviceCheckStatus string) {
	service := serviceNameServiceID(consulConnection, serviceCheckStatus)
	for serviceName, serviceID := range service {
		fmt.Println(serviceName + " " + serviceID + " has been deregistered!!!")
		err := consulConnection.Agent().ServiceDeregister(serviceID)
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	flag.StringVar(&url, []string{"u", "-url"}, "localhost", "Consul members endpoint. Default: localhost")
	flag.StringVar(&port, []string{"p", "-port"}, "8500", "Consul members endpoint port. Default: 8500")
	flag.IntVar(&nsc, []string{"nsc", "-nodestatuscode"}, 4, "Node status code. Default: 4")
	flag.BoolVar(&showMember, []string{"sm", "-showMembers"}, false, "Show a list of members")
	flag.BoolVar(&showNodeStatus, []string{"sns", "-showNodeStatus"}, false, "Show node status")
	flag.BoolVar(&showChecks, []string{"schk", "-showChecks"}, false, "Show a list of checks")
	flag.BoolVar(&showSrv, []string{"sasrv", "-showAllServices"}, false, "Show a list of services")
	flag.StringVar(&srvState, []string{"ss", "-serviceState"}, "critical", "Deregister Service State. Default: critical")
	flag.BoolVar(&deregisterSrv, []string{"dsrv", "-deregisterService"}, false, "Deregister service")
	flag.BoolVar(&listSrvInState, []string{"lsrvis", "-listServiceInState"}, false, "Show a list of services")
	flag.Parse()

	// var status = flag.String("status", "", "Consul Member Status for eviction")
	// var Debug = flag.String("debug", "", "Debug information output")
	// varDryrun = flag.String("cmd", "", "Do not remove the nodes, just preten")

	consulClient := connection(url, port)

	// fmt.Printf("%T", consulClient)

	if deregisterSrv == true {
		deregisterService(consulClient, srvState)
		os.Exit(0)
	}

	if listSrvInState == true {
		listServicesInState(consulClient, srvState)
		os.Exit(0)
	}

	if showNodeStatus == true {
		shownodestaus(consulClient)
		os.Exit(0)
	}

	if showMember == true {
		listmembers(consulClient)
		os.Exit(0)
	}

	if showChecks == true {
		listChecks(consulClient)
		os.Exit(0)
	}

	if showSrv == true {
		listServices(consulClient)
		os.Exit(0)
	}

	session := awsdiscovery.AwsSessIon("eu-west-1")
	filter := awsdiscovery.AwsFilter("Location", "qa")
	ips := awsdiscovery.AwsInstancePrivateIP(session, filter)

	for _, ip := range ips {
		fmt.Println(ip)
	}
}
