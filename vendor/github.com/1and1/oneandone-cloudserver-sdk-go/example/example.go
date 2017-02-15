package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/1and1/oneandone-cloudserver-sdk-go"
)

func main() {
	//Set an authentication token
	token := oneandone.SetToken("82ee732b8d47e451be5c6ad5b7b56c81")
	//Create an API client
	api := oneandone.New(token, oneandone.BaseUrl)

	// List server appliances
	saps, err := api.ListServerAppliances()

	//	printObject(saps)
	//	time.Sleep(time.Second * 10)

	var sa oneandone.ServerAppliance
	for _, a := range saps {
		if a.Type == "IMAGE" {
			sa = a
		}
	}

	// Create a server
	req := oneandone.ServerRequest{
		Name:        "Example Server",
		Description: "Example server description.",
		ApplianceId: sa.Id,
		PowerOn:     true,
		Hardware: oneandone.Hardware{
			Vcores:            1,
			CoresPerProcessor: 1,
			Ram:               2,
			Hdds: []oneandone.Hdd{
				{
					Size:   sa.MinHddSize,
					IsMain: true,
				},
			},
		},
	}
	server_id, server, err := api.CreateServer(&req)
	if err == nil {
		// Wait until server is created and powered on for at most 60 x 10 seconds
		err = api.WaitForState(server, "POWERED_ON", 10, 60)
	}

	// Get a server
	server, err = api.GetServer(server_id)

	//	printObject(server)
	//	time.Sleep(time.Second * 10)

	// Create a load balancer
	lbr := oneandone.LoadBalancerRequest{
		Name:                "Load Balancer Example",
		Description:         "API created load balancer.",
		Method:              "ROUND_ROBIN",
		Persistence:         oneandone.Bool2Pointer(true),
		PersistenceTime:     oneandone.Int2Pointer(1200),
		HealthCheckTest:     "TCP",
		HealthCheckInterval: oneandone.Int2Pointer(40),
		Rules: []oneandone.LoadBalancerRule{
			{
				Protocol:     "TCP",
				PortBalancer: 80,
				PortServer:   80,
				Source:       "0.0.0.0",
			},
		},
	}
	var lb *oneandone.LoadBalancer
	var lb_id string
	lb_id, lb, err = api.CreateLoadBalancer(&lbr)
	if err != nil {
		api.WaitForState(lb, "ACTIVE", 10, 30)
	}

	// Get a load balancer
	lb, err = api.GetLoadBalancer(lb.Id)

	//	printObject(lb)
	//	time.Sleep(time.Second * 10)

	// Assign a load balancer to server's IP
	server, err = api.AssignServerIpLoadBalancer(server.Id, server.Ips[0].Id, lb_id)

	//	printObject(server)
	//	time.Sleep(time.Second * 10)

	// Create a firewall policy
	fpr := oneandone.FirewallPolicyRequest{
		Name:        "Firewall Policy Example",
		Description: "API created firewall policy.",
		Rules: []oneandone.FirewallPolicyRule{
			{
				Protocol: "TCP",
				PortFrom: oneandone.Int2Pointer(80),
				PortTo:   oneandone.Int2Pointer(80),
			},
		},
	}
	var fp *oneandone.FirewallPolicy
	var fp_id string
	fp_id, fp, err = api.CreateFirewallPolicy(&fpr)
	if err == nil {
		api.WaitForState(fp, "ACTIVE", 10, 30)
	}

	// Get a firewall policy
	fp, err = api.GetFirewallPolicy(fp_id)

	// Add servers IPs to a firewall policy.
	ips := []string{server.Ips[0].Id}
	fp, err = api.AddFirewallPolicyServerIps(fp.Id, ips)
	if err == nil {
		api.WaitForState(fp, "ACTIVE", 10, 60)
	}

	//	printObject(fp)
	//	time.Sleep(time.Second * 10)

	//Shutdown a server using 'SOFTWARE' method
	server, err = api.ShutdownServer(server.Id, false)
	if err != nil {
		err = api.WaitForState(server, "POWERED_OFF", 5, 20)
	}

	//	printObject(server)
	//	time.Sleep(time.Second * 10)

	// Delete a load balancer
	lb, err = api.DeleteLoadBalancer(lb.Id)
	if err != nil {
		err = api.WaitUntilDeleted(lb)
	}

	//	printObject(lb)
	//	time.Sleep(time.Second * 10)

	// Delete a firewall policy
	fp, err = api.DeleteFirewallPolicy(fp.Id)
	if err != nil {
		err = api.WaitUntilDeleted(fp)
	}

	//	printObject(fp)
	//	time.Sleep(time.Second * 10)

	// List usages in last 24h
	var usages *oneandone.Usages
	usages, err = api.ListUsages("LAST_24H", nil, nil)

	printObject(usages)
	//	time.Sleep(time.Second * 10)

	// List usages in last 5 hours
	n := time.Now()
	ed := time.Date(n.Year(), n.Month(), n.Day(), n.Hour(), n.Minute(), n.Second(), 0, time.UTC)
	sd := ed.Add(-(time.Hour * 5))
	usages, err = api.ListUsages("CUSTOM", &sd, &ed)

	//	printObject(usages)
	//	time.Sleep(time.Second * 10)

	//Create a shared storage
	ssr := oneandone.SharedStorageRequest{
		Name:        "Shared Storage Example",
		Description: "API alocated 100 GB disk.",
		Size:        oneandone.Int2Pointer(100),
	}
	var ss *oneandone.SharedStorage
	var ss_id string
	ss_id, ss, err = api.CreateSharedStorage(&ssr)
	if err != nil {
		api.WaitForState(ss, "ACTIVE", 10, 30)
	}

	//	printObject(ss)
	//	time.Sleep(time.Second * 10)

	// List shared storages on page 1, 5 results per page and sort by 'name' field.
	// Include only 'name', 'size' and 'minimum_size_allowed' fields in the result.
	var shs []oneandone.SharedStorage
	shs, err = api.ListSharedStorages(1, 5, "name", "", "name,size,minimum_size_allowed")

	printObject(shs)
	//	time.Sleep(time.Second * 10)

	// List all shared storages that contain 'example' string
	shs, err = api.ListSharedStorages(0, 0, "", "example", "")

	//	printObject(shs)
	//	time.Sleep(time.Second * 10)

	// Delete a shared storage
	ss, err = api.DeleteSharedStorage(ss_id)
	if err == nil {
		err = api.WaitUntilDeleted(ss)
	}

	//	printObject(ss)
	//	time.Sleep(time.Second * 10)

	// Delete a server
	server, err = api.DeleteServer(server.Id, false)
	if err == nil {
		err = api.WaitUntilDeleted(server)
	}

	//	printObject(server)

	//	The next example illustrates how to create a `TYPO3` application server
	//  of a fixed size with an initial password and a firewall policy that has just been created.

	// Create a new firewall policy
	fpr = oneandone.FirewallPolicyRequest{
		Name: "HTTPS Traffic Policy",
		Rules: []oneandone.FirewallPolicyRule{
			{
				Protocol: "TCP",
				PortFrom: oneandone.Int2Pointer(443),
				PortTo:   oneandone.Int2Pointer(443),
			},
		},
	}

	_, fp, err = api.CreateFirewallPolicy(&fpr)
	if fp != nil && err == nil {
		api.WaitForState(fp, "ACTIVE", 5, 60)

		// Look for the TYPO3 application appliance
		saps, _ := api.ListServerAppliances(0, 0, "", "typo3", "")

		var sa oneandone.ServerAppliance
		for _, a := range saps {
			if a.Type == "APPLICATION" {
				sa = a
				break
			}
		}

		var fixed_flavours []oneandone.FixedInstanceInfo
		var fixed_size_id string

		fixed_flavours, err = api.ListFixedInstanceSizes()
		for _, fl := range fixed_flavours {
			//look for 'M' size
			if fl.Name == "M" {
				fixed_size_id = fl.Id
				break
			}
		}

		req := oneandone.ServerRequest{
			Name:        "TYPO3 Server",
			ApplianceId: sa.Id,
			PowerOn:     true,
			Password:    "ucr_kXW8,.2SdMU",
			Hardware: oneandone.Hardware{
				FixedInsSizeId: fixed_size_id,
			},
			FirewallPolicyId: fp.Id,
		}
		_, server, _ := api.CreateServer(&req)
		if server != nil {
			api.WaitForState(server, "POWERED_ON", 10, 90)
		}
	}
}

func printObject(in interface{}) {
	bytes, _ := json.MarshalIndent(in, "", "    ")
	fmt.Printf("%v\n", string(bytes))
}
