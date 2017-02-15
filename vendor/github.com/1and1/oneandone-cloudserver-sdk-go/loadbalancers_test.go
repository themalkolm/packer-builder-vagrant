package oneandone

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

var (
	set_lb       sync.Once
	test_lb_name string
	test_lb_desc string
	test_lb      *LoadBalancer
)

const (
	lb_round_robin = "ROUND_ROBIN"
	lb_least_con   = "LEAST_CONNECTIONS"
)

// Helper functions

func create_load_balancer() *LoadBalancer {
	rand.Seed(time.Now().UnixNano())
	rint := rand.Intn(999)
	test_lb_name = fmt.Sprintf("LoadBalancer_%d", rint)
	test_lb_desc = fmt.Sprintf("LoadBalancer_%d description", rint)
	req := LoadBalancerRequest{
		Name:                test_lb_name,
		Description:         test_lb_desc,
		Method:              lb_round_robin,
		Persistence:         Bool2Pointer(true),
		PersistenceTime:     Int2Pointer(60),
		HealthCheckTest:     "TCP",
		HealthCheckInterval: Int2Pointer(300),
		Rules: []LoadBalancerRule{
			{
				Protocol:     "TCP",
				PortBalancer: 8080,
				PortServer:   8089,
				Source:       "0.0.0.0",
			},
		},
	}
	fmt.Printf("Creating new load balancer '%s'...\n", test_lb_name)
	lb_id, lb, err := api.CreateLoadBalancer(&req)
	if err != nil {
		fmt.Printf("Unable to create a load balancer. Error: %s", err.Error())
		return nil
	}
	if lb_id == "" || lb.Id == "" {
		fmt.Printf("Unable to create load balancer '%s'.", test_lb_name)
		return nil
	}

	api.WaitForState(lb, "ACTIVE", 10, 30)
	return lb
}

func set_load_balancer() {
	test_lb = create_load_balancer()
}

// /load_balancers tests

func TestCreateLoadBalancer(t *testing.T) {
	set_lb.Do(set_load_balancer)

	if test_lb == nil {
		t.Errorf("CreateLoadBalancer failed.")
	}
	if test_lb.Id == "" {
		t.Errorf("Missing load balancer ID.")
	}
	if test_lb.Name != test_lb_name {
		t.Errorf("Wrong name of the load balancer.")
	}
	if test_lb.Description != test_lb_desc {
		t.Errorf("Wrong load balancer description.")
	}
	if !test_lb.Persistence {
		t.Errorf("Wrong load balancer persistence.")
	}
	if test_lb.PersistenceTime != 60 {
		t.Errorf("Wrong persistence time for load balancer '%s'.", test_lb.Name)
	}
	if test_lb.HealthCheckInterval != 300 {
		t.Errorf("Wrong health check interval for load balancer '%s'.", test_lb.Name)
	}
	if test_lb.HealthCheckTest != "TCP" {
		t.Errorf("Wrong health check test for load balancer '%s'.", test_lb.Name)
	}
	if test_lb.Method != lb_round_robin {
		t.Errorf("Wrong method for load balancer '%s'.", test_lb.Name)
	}
	if len(test_lb.Rules) != 1 {
		t.Errorf("Wrong number of load balancer '%s' rules.", test_lb.Name)
	} else {
		if test_lb.Rules[0].Protocol != "TCP" {
			t.Errorf("Wrong protocol of load balancer '%s' rule.", test_lb.Name)
		}
		if test_lb.Rules[0].PortBalancer != 8080 {
			t.Errorf("Wrong rule balancer port of load balancer '%s'.", test_lb.Name)
		}
		if test_lb.Rules[0].PortServer != 8089 {
			t.Errorf("Wrong rule server port of load balancer '%s'.", test_lb.Name)
		}
		if test_lb.Rules[0].Source != "0.0.0.0" {
			t.Errorf("Wrong source in load balancer '%s' rule.", test_lb.Name)
		}
	}
}

func TestGetLoadBalancer(t *testing.T) {
	set_lb.Do(set_load_balancer)

	fmt.Printf("Getting load balancer '%s'...\n", test_lb.Name)
	lb, err := api.GetLoadBalancer(test_lb.Id)

	if err != nil {
		t.Errorf("GetLoadBalancer failed. Error: " + err.Error())
	}
	if lb.Id != test_lb.Id {
		t.Errorf("Wrong load balancer ID.")
	}
}

func TestListLoadBalancers(t *testing.T) {
	set_lb.Do(set_load_balancer)
	fmt.Println("Listing all load balancers...")

	res, err := api.ListLoadBalancers()
	if err != nil {
		t.Errorf("ListLoadBalancers failed. Error: " + err.Error())
	}
	if len(res) == 0 {
		t.Errorf("No load balancer found.")
	}

	res, err = api.ListLoadBalancers(1, 1, "", "", "id,name")

	if err != nil {
		t.Errorf("ListLoadBalancers with parameter options failed. Error: " + err.Error())
	}
	if len(res) == 0 {
		t.Errorf("No load balancer found.")
	}
	if len(res) > 1 {
		t.Errorf("Wrong number of objects per page.")
	}
	if res[0].Id == "" {
		t.Errorf("Filtering parameters failed.")
	}
	if res[0].Name == "" {
		t.Errorf("Filtering parameters failed.")
	}
	if res[0].State != "" {
		t.Errorf("Filtering parameters failed.")
	}
	// Test for error response
	res, err = api.ListLoadBalancers(0, "name")
	if res != nil || err == nil {
		t.Errorf("ListLoadBalancers failed to handle incorrect argument type.")
	}

	res, err = api.ListLoadBalancers(0, 0, "", test_lb.Name, "")

	if err != nil {
		t.Errorf("ListLoadBalancers with parameter options failed. Error: " + err.Error())
	}
	if len(res) != 1 {
		t.Errorf("Search parameter failed.")
	}
	if res[0].Name != test_lb.Name {
		t.Errorf("Search parameter failed.")
	}
}

func TestAddLoadBalancerServerIps(t *testing.T) {
	set_lb.Do(set_load_balancer)
	sync_server.Do(func() { deploy_test_server(false) })

	fmt.Printf("Assigning server IPs to load balancer '%s'...\n", test_lb.Name)

	test_server, _ = api.GetServer(test_server.Id)

	ips := []string{test_server.Ips[0].Id}
	lb, err := api.AddLoadBalancerServerIps(test_lb.Id, ips)

	if err != nil {
		t.Errorf("AddLoadBalancerServerIps failed. Error: " + err.Error())
	}

	api.WaitForState(lb, "ACTIVE", 10, 60)
	lb, _ = api.GetLoadBalancer(lb.Id)

	if len(lb.ServerIps) != 1 {
		t.Errorf("Found no server IP attached to the load balancer.")
	}
	if lb.ServerIps[0].Id != test_server.Ips[0].Id {
		t.Errorf("Wrong server IP attached to the load balancer.")
	}
	test_lb = lb
}

func TestGetLoadBalancerServerIp(t *testing.T) {
	set_lb.Do(set_load_balancer)

	fmt.Printf("Getting server IPs of load balancer '%s'...\n", test_lb.Name)
	lb_ser_ip, err := api.GetLoadBalancerServerIp(test_lb.Id, test_lb.ServerIps[0].Id)

	if err != nil {
		t.Errorf("GetLoadBalancerServerIp failed. Error: " + err.Error())
	}
	if lb_ser_ip.Id != test_lb.ServerIps[0].Id {
		t.Errorf("Wrong ID of the server IP attached to load balancer '%s'.", test_lb.Name)
	}
	if lb_ser_ip.Ip != test_lb.ServerIps[0].Ip {
		t.Errorf("Wrong server IP address attached to load balancer '%s'.", test_lb.Name)
	}
}

func TestListLoadBalancerServerIps(t *testing.T) {
	set_lb.Do(set_load_balancer)
	sync_server.Do(func() { deploy_test_server(false) })

	fmt.Printf("Listing server IPs of load balancer '%s'...\n", test_lb.Name)
	lb_ips, err := api.ListLoadBalancerServerIps(test_lb.Id)

	if err != nil {
		t.Errorf("ListLoadBalancerServerIps failed. Error: " + err.Error())
	}
	if len(lb_ips) != 1 {
		t.Errorf("Wrong number of server IPs added to load balancer '%s'.", test_lb.Name)
	}
	if lb_ips[0].Id != test_server.Ips[0].Id {
		t.Errorf("Wrong server IP added to load balancer '%s'.", test_lb.Name)
	}
}

func TestDeleteLoadBalancerServerIp(t *testing.T) {
	set_lb.Do(set_load_balancer)
	sync_server.Do(func() { deploy_test_server(false) })

	fmt.Printf("Deleting server IP from load balancer '%s'...\n", test_lb.Name)
	lb, err := api.DeleteLoadBalancerServerIp(test_lb.Id, test_server.Ips[0].Id)

	if err != nil {
		t.Errorf("DeleteLoadBalancerServerIp failed. Error: " + err.Error())
	}

	api.WaitForState(lb, "ACTIVE", 10, 60)
	lb, err = api.GetLoadBalancer(lb.Id)

	if err != nil {
		t.Errorf("Deleting server IP from the load balancer failed.")
	}
	if len(lb.ServerIps) > 0 {
		t.Errorf("IP not deleted from the load balancer.")
	}
}

func TestGetLoadBalancerRule(t *testing.T) {
	set_lb.Do(set_load_balancer)

	fmt.Printf("Getting the rule of load balancer '%s'...\n", test_lb.Name)
	lb_rule, err := api.GetLoadBalancerRule(test_lb.Id, test_lb.Rules[0].Id)

	if err != nil {
		t.Errorf("GetLoadBalancerRule failed. Error: " + err.Error())
	}
	if lb_rule.Id != test_lb.Rules[0].Id {
		t.Errorf("Wrong rule ID.")
	}
	if lb_rule.PortBalancer != test_lb.Rules[0].PortBalancer {
		t.Errorf("Wrong rule port_balancer field.")
	}
	if lb_rule.PortServer != test_lb.Rules[0].PortServer {
		t.Errorf("Wrong rule port_server field.")
	}
	if lb_rule.Protocol != test_lb.Rules[0].Protocol {
		t.Errorf("Wrong rule protocol.")
	}
	if lb_rule.Source != test_lb.Rules[0].Source {
		t.Errorf("Wrong rule source IP.")
	}
}

func TestAddLoadBalancerRules(t *testing.T) {
	set_lb.Do(set_load_balancer)

	fmt.Printf("Adding rules to load balancer '%s'...\n", test_lb.Name)
	rules := []LoadBalancerRule{
		{
			Protocol:     "TCP",
			PortBalancer: 35367,
			PortServer:   35367,
			Source:       "0.0.0.0",
		},
		{
			Protocol:     "TCP",
			PortBalancer: 815,
			PortServer:   815,
		},
	}
	lb, err := api.AddLoadBalancerRules(test_lb.Id, rules)

	if err != nil {
		t.Errorf("AddLoadBalancerRules failed. Error: " + err.Error())
	} else {
		api.WaitForState(lb, "ACTIVE", 10, 60)
	}
	lb, _ = api.GetLoadBalancer(lb.Id)
	if len(lb.Rules) != 3 {
		t.Errorf("Unable to add rules to load balancer '%s'.\n", test_lb.Name)
	}
}

func TestListLoadBalancerRules(t *testing.T) {
	set_lb.Do(set_load_balancer)

	fmt.Printf("Listing load balancer '%s' rules...\n", test_lb.Name)
	lb_rules, err := api.ListLoadBalancerRules(test_lb.Id)

	if err != nil {
		t.Errorf("ListLoadBalancerRules failed. Error: " + err.Error())
	}
	if len(lb_rules) != 3 {
		t.Errorf("Wrong number of rules found at load balancer '%s'.", test_lb.Name)
	}
}

func TestDeleteLoadBalancerRule(t *testing.T) {
	set_lb.Do(set_load_balancer)

	lbr, _ := api.ListLoadBalancerRules(test_lb.Id)
	fmt.Printf("Deleting rule '%s' from load balancer '%s'...\n", lbr[0].Id, test_lb.Name)
	lb, err := api.DeleteLoadBalancerRule(test_lb.Id, lbr[0].Id)

	if err != nil {
		t.Errorf("DeleteLoadBalancerRule failed. Error: " + err.Error())
	}

	api.WaitForState(lb, "ACTIVE", 10, 60)
	lb, err = api.GetLoadBalancer(lb.Id)

	if err != nil {
		t.Errorf("Deleting rule from the load balancer failed.")
	}
	if len(lb.Rules) != 2 {
		t.Errorf("Rule not deleted from the load balancer.")
	}
	for _, rule := range lb.Rules {
		if rule.Id == lbr[0].Id {
			t.Errorf("Rule not deleted from the load balancer.")
		}
	}
}

func TestUpdateLoadBalancer(t *testing.T) {
	set_lb.Do(set_load_balancer)

	fmt.Printf("Updating load balancer '%s'...\n", test_lb.Name)
	new_name := test_lb.Name + "_updated"
	new_desc := test_lb.Description + "_updated"
	lbu := LoadBalancerRequest{
		Name:        new_name,
		Description: new_desc,
		Method:      lb_least_con,
		Persistence: Bool2Pointer(false),
	}
	lb, err := api.UpdateLoadBalancer(test_lb.Id, &lbu)

	if err != nil {
		t.Errorf("UpdateLoadBalancer failed. Error: " + err.Error())
	} else {
		api.WaitForState(lb, "ACTIVE", 10, 30)
	}
	lb, _ = api.GetLoadBalancer(lb.Id)
	if lb.Name != new_name {
		t.Errorf("Failed to update load balancer name.")
	}
	if lb.Description != new_desc {
		t.Errorf("Failed to update load balancer description.")
	}
	if lb.Method != lb_least_con {
		t.Errorf("Failed to update load balancer method.")
	}
	if lb.Persistence {
		t.Errorf("Failed to update load balancer persistence.")
	}
}

func TestDeleteLoadBalancer(t *testing.T) {
	set_lb.Do(set_load_balancer)

	fmt.Printf("Deleting load balancer '%s'...\n", test_lb.Name)
	lb, err := api.DeleteLoadBalancer(test_lb.Id)

	if err != nil {
		t.Errorf("DeleteLoadBalancer failed. Error: " + err.Error())
	} else {
		api.WaitUntilDeleted(lb)
	}
	lb, err = api.GetLoadBalancer(lb.Id)

	if lb != nil {
		t.Errorf("Unable to delete the load balancer.")
	} else {
		test_lb = nil
	}
}
