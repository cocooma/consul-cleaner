package main

import (
	"fmt"
	"os"

	flag "github.com/docker/docker/pkg/mflag"
	"github.com/hashicorp/consul/api"
)

var (
	str, url, port, srvState                                       string
	showMember, showChecks, showSrv, deregisterSrv, listSrvInState bool
)

// var wg sync.WaitGroup

func listmembers(consul_client *api.Client) {
	members, _ := consul_client.Agent().Members(false)
	for _, server := range members {
		fmt.Println(server.Name)
	}
}

func showmembers(members []*api.AgentMember) {
	for _, server := range members {
		fmt.Println(server.Status)
	}
}

func listChecks(consul_client *api.Client) {
	checks, _ := consul_client.Agent().Checks()
	for _, check := range checks {
		fmt.Println(check.Node, check.Name, check.Status)
	}
}

func listServices(consul_client *api.Client) {
	services, _ := consul_client.Agent().Services()
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

func listServicesInState(consul_connection *api.Client, serviceCheckStatus string) {
	service := serviceNameServiceID(consul_connection, serviceCheckStatus)
	for serviceName, serviceID := range service {
		fmt.Println(serviceName + " " + serviceID)
	}
}

func deregisterService(consul_connection *api.Client, serviceCheckStatus string) {
	service := serviceNameServiceID(consul_connection, serviceCheckStatus)
	for serviceName, serviceID := range service {
		fmt.Println(serviceName + " " + serviceID + " has been deregistered!!!")
		err := consul_connection.Agent().ServiceDeregister(serviceID)
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	flag.StringVar(&url, []string{"u", "-url"}, "localhost", "Consul members endpoint. Default: localhost")
	flag.StringVar(&port, []string{"p", "-port"}, "8500", "Consul members endpoint port. Default: 8500")
	flag.BoolVar(&showMember, []string{"sm", "-showMembers"}, false, "Show a list of members")
	flag.BoolVar(&showChecks, []string{"schk", "-showChecks"}, false, "Show a list of checks")
	flag.BoolVar(&showSrv, []string{"sasrv", "-showAllServices"}, false, "Show a list of services")
	flag.StringVar(&srvState, []string{"ss", "-serviceState"}, "critical", "Deregister Service State. Default: critical")
	flag.BoolVar(&deregisterSrv, []string{"dsrv", "-deregisterService"}, false, "Deregister service")
	flag.BoolVar(&listSrvInState, []string{"lsrvis", "-listServiceInState"}, false, "Show a list of services")
	flag.Parse()

	// var status = flag.String("status", "", "Consul Member Status for eviction")
	// var Debug = flag.String("debug", "", "Debug information output")
	// varDryrun = flag.String("cmd", "", "Do not remove the nodes, just preten")

	consul_client := connection(url, port)

	// fmt.Printf("%T", consul_client)

	if deregisterSrv == true {
		deregisterService(consul_client, srvState)
		os.Exit(0)
	}

	if listSrvInState == true {
		listServicesInState(consul_client, srvState)
		os.Exit(0)
	}

	if showMember == true {
		listmembers(consul_client)
		os.Exit(0)
	}

	if showChecks == true {
		listChecks(consul_client)
		os.Exit(0)
	}

	if showSrv == true {
		listServices(consul_client)
		os.Exit(0)
	}

}
