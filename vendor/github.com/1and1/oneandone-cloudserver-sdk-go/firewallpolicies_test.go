package oneandone

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

var (
	set_fp       sync.Once
	test_fp_name string
	test_fp_desc string
	test_fp      *FirewallPolicy
)

// Helper functions

func create_firewall_policy() *FirewallPolicy {
	rand.Seed(time.Now().UnixNano())
	rint := rand.Intn(999)
	test_fp_name = fmt.Sprintf("FirewallPol_%d", rint)
	test_fp_desc = fmt.Sprintf("FirewallPol_%d description", rint)
	req := FirewallPolicyRequest{
		Name:        test_fp_name,
		Description: test_fp_desc,
		Rules: []FirewallPolicyRule{
			{
				Protocol: "UDP",
				PortFrom: Int2Pointer(161),
				PortTo:   Int2Pointer(162),
			},
		},
	}
	fmt.Printf("Creating new firewall policy '%s'...\n", test_fp_name)
	fp_id, fp, err := api.CreateFirewallPolicy(&req)
	if err != nil {
		fmt.Printf("Unable to create firewall policy '%s'. Error: %s", test_fp_name, err.Error())
		return nil
	}
	if fp_id == "" || fp.Id == "" {
		fmt.Printf("Unable to create firewall policy '%s'.", test_fp_name)
		return nil
	}

	api.WaitForState(fp, "ACTIVE", 10, 30)

	return fp
}

func set_firewall_policy() {
	test_fp = create_firewall_policy()
}

// /firewall_policies tests

func TestCreateFirewallPolicy(t *testing.T) {
	set_fp.Do(set_firewall_policy)

	if test_fp == nil {
		t.Errorf("CreateFirewallPolicy failed.")
	}
	if test_fp.Name != test_fp_name {
		t.Errorf("Wrong name of the firewall policy.")
	}
	if test_fp.Description != test_fp_desc {
		t.Errorf("Wrong description of the firewall policy.")
	}
}

func TestGetFirewallPolicy(t *testing.T) {
	set_fp.Do(set_firewall_policy)

	fmt.Printf("Getting firewall policy '%s'...\n", test_fp.Name)
	fp, err := api.GetFirewallPolicy(test_fp.Id)

	if err != nil {
		t.Errorf("GetFirewallPolicy failed. Error: " + err.Error())
	}
	if fp.Id != test_fp.Id {
		t.Errorf("Wrong ID of the firewall policy.")
	}
	if fp.Name != test_fp.Name {
		t.Errorf("Wrong name of the firewall policy.")
	}
	if fp.Description != test_fp.Description {
		t.Errorf("Wrong description of the firewall policy.")
	}
}

func TestListFirewallPolicies(t *testing.T) {
	set_fp.Do(set_firewall_policy)
	fmt.Println("Listing all firewall policies...")

	res, err := api.ListFirewallPolicies()
	if err != nil {
		t.Errorf("ListFirewallPolicies failed. Error: " + err.Error())
	}
	if len(res) == 0 {
		t.Errorf("No firewall policy found.")
	}

	res, err = api.ListFirewallPolicies(1, 4, "name", "", "id,name")

	if err != nil {
		t.Errorf("ListFirewallPolicies with parameter options failed. Error: " + err.Error())
	}
	if len(res) == 0 {
		t.Errorf("No firewall policy found.")
	}
	// Here we consider two default policies as well, Linux and Windows.
	if len(res) < 3 || len(res) > 4 {
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
	if res[0].Name > res[1].Name {
		t.Errorf("Sorting parameters failed.")
	}
	// Test for error response
	res, err = api.ListFirewallPolicies("name", 0)
	if res != nil || err == nil {
		t.Errorf("ListFirewallPolicies failed to handle incorrect argument type.")
	}

	res, err = api.ListFirewallPolicies(0, 0, "", test_fp.Name, "")

	if err != nil {
		t.Errorf("ListFirewallPolicies with parameter options failed. Error: " + err.Error())
	}
	if len(res) != 1 {
		t.Errorf("Search parameter failed.")
	}
	if res[0].Name != test_fp.Name {
		t.Errorf("Search parameter failed.")
	}
}

func TestAddFirewallPolicyServerIps(t *testing.T) {
	set_fp.Do(set_firewall_policy)

	fmt.Printf("Assigning server IPs to firewall policy '%s'...\n", test_fp.Name)
	sync_server.Do(func() { deploy_test_server(false) })

	test_server, _ = api.GetServer(test_server.Id)

	ips := []string{test_server.Ips[0].Id}
	fp, err := api.AddFirewallPolicyServerIps(test_fp.Id, ips)

	if err != nil {
		t.Errorf("AddFirewallPolicyServerIps failed. Error: " + err.Error())
	} else {
		api.WaitForState(fp, "ACTIVE", 10, 60)
		fp, err = api.GetFirewallPolicy(fp.Id)
		if err != nil {
			t.Errorf("FirewallPolicyAddServerIps failed. Error: " + err.Error())
		}
		if len(fp.ServerIps) != 1 {
			t.Errorf("Found no server IP attached to the firewall policy.")
		}
		if fp.ServerIps[0].Id != test_server.Ips[0].Id {
			t.Errorf("Wrong server IP attached to the firewall policy.")
		}
		test_fp = fp
	}
}

func TestGetFirewallPolicyServerIp(t *testing.T) {
	set_fp.Do(set_firewall_policy)

	fmt.Printf("Getting server IPs of firewall policy '%s'...\n", test_fp.Name)
	fp_ser_ip, err := api.GetFirewallPolicyServerIp(test_fp.Id, test_fp.ServerIps[0].Id)

	if err != nil {
		t.Errorf("GetFirewallPolicyServerIp failed. Error: " + err.Error())
	}
	if fp_ser_ip.Id != test_fp.ServerIps[0].Id {
		t.Errorf("Wrong ID of the server IP attached to firewall policy '%s'.", test_fp.Name)
	}
	if fp_ser_ip.Ip != test_fp.ServerIps[0].Ip {
		t.Errorf("Wrong server IP address attached to firewall policy '%s'.", test_fp.Name)
	}
}

func TestListFirewallPolicyServerIps(t *testing.T) {
	set_fp.Do(set_firewall_policy)
	sync_server.Do(func() { deploy_test_server(false) })

	fmt.Printf("Listing server IPs of firewall policy '%s'...\n", test_fp.Name)
	fp_ips, err := api.ListFirewallPolicyServerIps(test_fp.Id)

	if err != nil {
		t.Errorf("ListFirewallPolicyServerIps failed. Error: " + err.Error())
	}
	if len(fp_ips) != 1 {
		t.Errorf("Wrong number of server IPs added to firewall policy '%s'.", test_fp.Name)
	}
	if fp_ips[0].Id != test_server.Ips[0].Id {
		t.Errorf("Wrong server IP added to firewall policy '%s'.", test_fp.Name)
	}
}

func TestDeleteFirewallPolicyServerIp(t *testing.T) {
	set_fp.Do(set_firewall_policy)
	sync_server.Do(func() { deploy_test_server(false) })

	fmt.Printf("Deleting server IP from firewall policy '%s'...\n", test_fp.Name)
	fp, err := api.DeleteFirewallPolicyServerIp(test_fp.Id, test_server.Ips[0].Id)

	if err != nil {
		t.Errorf("DeleteFirewallPolicyServerIp failed. Error: " + err.Error())
	}

	api.WaitForState(fp, "ACTIVE", 10, 60)
	fp, err = api.GetFirewallPolicy(fp.Id)

	if err != nil {
		t.Errorf("Deleting server IP from the firewall policy failed.")
	}
	if len(fp.ServerIps) > 0 {
		t.Errorf("IP not deleted from the firewall policy.")
	}
}

func TestGetFirewallPolicyRule(t *testing.T) {
	set_fp.Do(set_firewall_policy)

	fmt.Printf("Getting the rule of firewall policy '%s'...\n", test_fp.Name)
	fp_rule, err := api.GetFirewallPolicyRule(test_fp.Id, test_fp.Rules[0].Id)

	if err != nil {
		t.Errorf("GetFirewallPolicyRule failed. Error: " + err.Error())
	}
	if fp_rule.Id != test_fp.Rules[0].Id {
		t.Errorf("Wrong rule ID.")
	}
	if *fp_rule.PortFrom != *test_fp.Rules[0].PortFrom {
		t.Errorf("Wrong rule port_from field.")
	}
	if *fp_rule.PortTo != *test_fp.Rules[0].PortTo {
		t.Errorf("Wrong rule port_to field.")
	}
	if fp_rule.Protocol != test_fp.Rules[0].Protocol {
		t.Errorf("Wrong rule protocol.")
	}
	if fp_rule.SourceIp != test_fp.Rules[0].SourceIp {
		t.Errorf("Wrong rule source IP.")
	}
}

func TestUpdateFirewallPolicy(t *testing.T) {
	set_fp.Do(set_firewall_policy)

	fmt.Printf("Updating firewall policy '%s'...\n", test_fp.Name)
	new_name := test_fp.Name + "_updated"
	new_desc := test_fp.Description + "_updated"
	fp, err := api.UpdateFirewallPolicy(test_fp.Id, new_name, new_desc)

	if err != nil {
		t.Errorf("UpdateFirewallPolicy failed. Error: " + err.Error())
	}
	if fp.Id != test_fp.Id {
		t.Errorf("Wrong firewall policy ID.")
	}
	if fp.Name != new_name {
		t.Errorf("Wrong firewall policy name.")
	}
	if fp.Description != new_desc {
		t.Errorf("Wrong firewall policy description.")
	}
	test_fp = fp
}

func TestAddFirewallPolicyRules(t *testing.T) {
	set_fp.Do(set_firewall_policy)

	fmt.Printf("Adding rules to firewall policy '%s'...\n", test_fp.Name)
	rules := []FirewallPolicyRule{
		{
			Protocol: "TCP",
			PortFrom: Int2Pointer(4567),
			PortTo:   Int2Pointer(4567),
			SourceIp: "0.0.0.0",
		},
		{
			Protocol: "TCP/UDP",
			PortFrom: Int2Pointer(143),
			PortTo:   Int2Pointer(143),
		},
		{
			Protocol: "GRE", // PortFrom & PortTo are optional for GRE, ICMP and IPSEC protocols.
		},
		{
			Protocol: "ICMP",
			PortFrom: nil,
			PortTo:   nil,
		},
		{
			Protocol: "IPSEC",
		},
	}
	fp, err := api.AddFirewallPolicyRules(test_fp.Id, rules)

	if err != nil {
		t.Errorf("AddFirewallPolicyRules failed. Error: " + err.Error())
	} else {
		api.WaitForState(fp, "ACTIVE", 10, 60)
	}
	fp, _ = api.GetFirewallPolicy(fp.Id)
	if len(fp.Rules) != 6 {
		t.Errorf("Unable to add rules to firewall policy '%s'.\n", test_fp.Name)
	}
}

func TestListFirewallPolicyRules(t *testing.T) {
	set_fp.Do(set_firewall_policy)

	fmt.Printf("Listing firewall policy '%s' rules...\n", test_fp.Name)
	fp_rules, err := api.ListFirewallPolicyRules(test_fp.Id)

	if err != nil {
		t.Errorf("ListFirewallPolicyRules failed. Error: " + err.Error())
	}
	if len(fp_rules) != 6 {
		t.Errorf("Wrong number of rules found at firewall policy '%s'.", test_fp.Name)
	}
}

func TestDeleteFirewallPolicyRule(t *testing.T) {
	set_fp.Do(set_firewall_policy)

	fmt.Printf("Deleting rule '%s' from firewall policy '%s'...\n", test_fp.Rules[0].Id, test_fp.Name)
	fp, err := api.DeleteFirewallPolicyRule(test_fp.Id, test_fp.Rules[0].Id)

	if err != nil {
		t.Errorf("DeleteFirewallPolicyRule failed. Error: " + err.Error())
	}

	api.WaitForState(fp, "ACTIVE", 10, 60)
	fp, err = api.GetFirewallPolicy(fp.Id)

	if err != nil {
		t.Errorf("Deleting rule from the firewall policy failed.")
	}
	if len(fp.Rules) != 5 {
		t.Errorf("Rule not deleted from the firewall policy.")
	}
	for _, rule := range fp.Rules {
		if rule.Id == test_fp.Rules[0].Id {
			t.Errorf("Rule not deleted from the firewall policy.")
		}
	}
}

func TestDeleteFirewallPolicy(t *testing.T) {
	set_fp.Do(set_firewall_policy)

	fmt.Printf("Deleting firewall policy '%s'...\n", test_fp.Name)
	fp, err := api.DeleteFirewallPolicy(test_fp.Id)

	if err != nil {
		t.Errorf("DeleteFirewallPolicy failed. Error: " + err.Error())
	} else {
		api.WaitUntilDeleted(fp)
	}

	fp, _ = api.GetFirewallPolicy(fp.Id)

	if fp != nil {
		t.Errorf("Unable to delete the firewall policy.")
	} else {
		test_fp = nil
	}
}
