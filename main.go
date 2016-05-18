package main

import (
	"fmt"
	"os"
	"sync"

	"consul-cleaner/awsdiscovery"

	flag "github.com/docker/docker/pkg/mflag"
	"github.com/hashicorp/consul/api"
)

var (
	str, url, port, srvState, awsregion, tag, tagvalue, hostdiscovery                                                         string
	showTargetHosts, showMmemberStatus, showChecks, showAllSrv, deregisterSrv, listSrvInState, showNodeStatus, forceLeaveNode bool
	nsc                                                                                                                       int
	hosts                                                                                                                     []string
	wg                                                                                                                        sync.WaitGroup
)

func connection(uurl, pport string) *api.Client {
	connection, err := api.NewClient(&api.Config{Address: uurl + ":" + pport})
	if err != nil {
		panic(err)
	}
	return connection
}

func consulmembers(consulClient *api.Client) []string {
	var ips []string
	members, _ := consulClient.Agent().Members(false)
	for _, server := range members {
		ips = append(ips, server.Name)
	}
	return ips
}

func awshosts(awsregion, tag, tagvalue string) []string {
	session := awsdiscovery.AwsSessIon(awsregion)
	filter := awsdiscovery.AwsFilter(tag, tagvalue)
	ips := awsdiscovery.AwsInstancePrivateIP(session, filter)
	return ips
}

func showtargethost(ips []string) {
	for _, server := range ips {
		fmt.Println(server)
	}
}

func showNodeStaus(consulClient *api.Client) {
	members, _ := consulClient.Agent().Members(false)
	for _, server := range members {
		fmt.Printf("%s %v \n", server.Name, server.Status)
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
	flag.IntVar(&nsc, []string{"nsc", "-nodeStatusCode"}, 4, "Node status code. Default: 4")
	flag.BoolVar(&showTargetHosts, []string{"sth", "-showTargetHosts"}, false, "Show target hosts")
	flag.BoolVar(&showNodeStatus, []string{"sns", "-showNodeStatus"}, false, "Show node status")
	flag.BoolVar(&showChecks, []string{"schk", "-showChecks"}, false, "Show a list of checks")
	flag.BoolVar(&showAllSrv, []string{"sasrv", "-showAllServices"}, false, "Show a list of services")
	flag.StringVar(&srvState, []string{"ss", "-serviceState"}, "critical", "Deregister Service State. Default: critical")
	flag.BoolVar(&deregisterSrv, []string{"drsrv", "-deregisterService"}, false, "Deregister service")
	flag.BoolVar(&listSrvInState, []string{"lsrvis", "-listServiceInState"}, false, "Show a list of services")
	flag.StringVar(&awsregion, []string{"ar", "-awsRegion"}, "eu-west-1", "AWS Region. Default: eu-west-1")
	flag.StringVar(&tag, []string{"t", "-tag"}, "", "AWS tag")
	flag.StringVar(&tagvalue, []string{"tv", "-tagValue"}, "", "AWS tag value")
	flag.StringVar(&hostdiscovery, []string{"hd", "-hostDiscovery"}, "aws", "Host discovery. 'consul' or 'aws'")
	flag.BoolVar(&forceLeaveNode, []string{"fl", "-forceLeaveNode"}, false, "Force leave consul node")
	flag.Parse()

	// var status = flag.String("status", "", "Consul Member Status for eviction")
	// var Debug = flag.String("debug", "", "Debug information output")
	// varDryrun = flag.String("cmd", "", "Do not remove the nodes, just preten")

	consulClient := connection(url, port)

	switch hostdiscovery {
	case "consul":
		hosts = consulmembers(consulClient)
	case "aws":
		hosts = awshosts(awsregion, tag, tagvalue)
	}

	if showTargetHosts == true {
		showtargethost(hosts)
		os.Exit(0)
	}

	if deregisterSrv == true {
		wg.Add(len(hosts))
		for _, ip := range hosts {
			go func(ip, srvState string) {
				deregisterService(connection(ip, "8500"), srvState)
				wg.Done()
			}(ip, srvState)
		}
	}

	if listSrvInState == true {
		wg.Add(len(hosts))
		for _, ip := range hosts {
			go func(ip, srvState string) {
				listServicesInState(connection(ip, "8500"), srvState)
				wg.Done()
			}(ip, srvState)
		}
	}

	if showNodeStatus == true {
		wg.Add(len(hosts))
		for _, ip := range hosts {
			go func(ip string) {
				showNodeStaus(connection(ip, "8500"))
				wg.Done()
			}(ip)
		}
	}

	if showChecks == true {
		wg.Add(len(hosts))
		for _, ip := range hosts {
			go func(ip string) {
				listChecks(connection(ip, "8500"))
				wg.Done()
			}(ip)
		}
	}

	if showAllSrv == true {
		wg.Add(len(hosts))
		for _, ip := range hosts {
			go func(ip string) {
				listServices(connection(ip, "8500"))
				wg.Done()
			}(ip)
		}
	}

	if showAllSrv == true {
		wg.Add(len(hosts))
		for _, ip := range hosts {
			go func(ip string) {
				listServices(connection(ip, "8500"))
				wg.Done()
			}(ip)
		}
	}

	if forceLeaveNode == true {
		wg.Add(len(hosts))
		for _, ip := range hosts {
			go func(ip string) {
				forceLeaveBadNode(connection(ip, "8500"), nsc)
				wg.Done()
			}(ip)
		}
	}
	wg.Wait()
}
