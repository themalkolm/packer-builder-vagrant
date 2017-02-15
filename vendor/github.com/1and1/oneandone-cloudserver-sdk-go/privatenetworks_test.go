package oneandone

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

var (
	set_pn  sync.Once
	pn_name string
	pn_desc string
	test_pn *PrivateNetwork
)

const (
	network_add = "192.168.46.0"
	sub_mask    = "255.255.255.0"
)

// Helper functions

func create_private_netwok() *PrivateNetwork {
	rand.Seed(time.Now().UnixNano())
	rint := rand.Intn(999)
	pn_name = fmt.Sprintf("PrivateNet_%d", rint)
	pn_desc = fmt.Sprintf("PrivateNet_%d description", rint)
	req := PrivateNetworkRequest{
		Name:           pn_name,
		Description:    pn_desc,
		NetworkAddress: network_add,
		SubnetMask:     sub_mask,
	}
	fmt.Printf("Creating new private network '%s'...\n", pn_name)
	prn_id, prn, err := api.CreatePrivateNetwork(&req)
	if err != nil {
		fmt.Printf("Unable to create a private network. Error: %s", err.Error())
		return nil
	}
	if prn_id == "" || prn.Id == "" {
		fmt.Printf("Unable to create private network '%s'.", pn_name)
		return nil
	}

	//prn = wait_for_state(prn, 5, 30, "ACTIVE")
	api.WaitForState(prn, "ACTIVE", 5, 60)

	return prn
}

func set_private_network() {
	test_pn = create_private_netwok()
}

// /private_networks tests

func TestCreatePrivateNetwork(t *testing.T) {
	set_pn.Do(set_private_network)

	if test_pn == nil {
		t.Errorf("CreatePrivateNetwork failed.")
	} else {
		if test_pn.Name != pn_name {
			t.Errorf("Wrong name of the private network.")
		}
		if test_pn.Description != pn_desc {
			t.Errorf("Wrong private network description.")
		}
		if test_pn.NetworkAddress != network_add {
			t.Errorf("Wrong private network address.")
		}
		if test_pn.SubnetMask != sub_mask {
			t.Errorf("Wrong private network subnet mask.")
		}
	}
}

func TestGetPrivateNetwork(t *testing.T) {
	set_pn.Do(set_private_network)

	fmt.Printf("Getting private network '%s'...\n", test_pn.Name)
	prn, err := api.GetPrivateNetwork(test_pn.Id)

	if err != nil {
		t.Errorf("GetPrivateNetwork failed. Error: " + err.Error())
	} else if prn.Id != test_pn.Id {
		t.Errorf("Wrong private network ID.")
	}
}

func TestListPrivateNetworks(t *testing.T) {
	set_pn.Do(set_private_network)
	fmt.Println("Listing all private networks...")

	res, err := api.ListPrivateNetworks()
	if err != nil {
		t.Errorf("ListPrivateNetworks failed. Error: " + err.Error())
	}
	if len(res) == 0 {
		t.Errorf("No private network found.")
	}

	res, err = api.ListPrivateNetworks(1, 1, "", "", "id,name")

	if err != nil {
		t.Errorf("ListPrivateNetworks with parameter options failed. Error: " + err.Error())
		return
	}
	if len(res) == 0 {
		t.Errorf("No private network found.")
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
	res, err = api.ListPrivateNetworks(0, 0, "", 7, "")
	if res != nil || err == nil {
		t.Errorf("ListPrivateNetworks failed to handle incorrect argument type.")
	}

	res, err = api.ListPrivateNetworks(0, 0, "", test_pn.Name, "")

	if err != nil {
		t.Errorf("ListPrivateNetworks with parameter options failed. Error: " + err.Error())
		return
	}
	if len(res) != 1 {
		t.Errorf("Search parameter failed.")
	}
	if res[0].Name != test_pn.Name {
		t.Errorf("Search parameter failed.")
	}
}

func TestAttachPrivateNetworkServers(t *testing.T) {
	set_pn.Do(set_private_network)
	sync_server.Do(func() { deploy_test_server(false) })

	fmt.Printf("Attaching servers to private network '%s'...\n", test_pn.Name)

	servers := make([]string, 1)
	servers[0] = test_server.Id

	time.Sleep(time.Second)
	prn, err := api.AttachPrivateNetworkServers(test_pn.Id, servers)

	if err != nil {
		t.Errorf("AttachPrivateNetworkServers failed. Error: " + err.Error())
		return
	}

	//	prn = wait_for_state(prn, 10, 60, "ACTIVE")
	api.WaitForState(prn, "ACTIVE", 10, 60)
	prn, err = api.GetPrivateNetwork(prn.Id)
	if prn == nil || err != nil {
		t.Errorf("Attaching servers to the private network failed.")
		return
	}
	if len(prn.Servers) != 1 {
		t.Errorf("No servers attached to the private network were found.")
	}
	if prn.Servers[0].Id != test_server.Id {
		t.Errorf("Wrong servers attached to the private network.")
	}
	test_pn = prn
}

func TestGetPrivateNetworkServer(t *testing.T) {
	set_pn.Do(set_private_network)
	sync_server.Do(func() { deploy_test_server(false) })

	fmt.Printf("Getting '%s' server attached to private network '%s'...\n", test_server.Name, test_pn.Name)
	pns, err := api.GetPrivateNetworkServer(test_pn.Id, test_server.Id)

	if err != nil {
		t.Errorf("GetPrivateNetworkServer failed. Error: " + err.Error())
		return
	}
	if pns.Id != test_server.Id {
		t.Errorf("Wrong server attached to the private network was found.")
	}
}

func TestListPrivateNetworkServers(t *testing.T) {
	set_pn.Do(set_private_network)

	fmt.Printf("Listing servers attached to private network '%s'...\n", test_pn.Name)
	pnss, err := api.ListPrivateNetworkServers(test_pn.Id)

	if err != nil {
		t.Errorf("ListPrivateNetworkServers failed. Error: " + err.Error())
	}
	if len(pnss) != 1 {
		t.Errorf("Wrong number of servers attached to the private network.")
	}
}

func TestDetachPrivateNetworkServer(t *testing.T) {
	set_pn.Do(set_private_network)
	sync_server.Do(func() { deploy_test_server(false) })

	fmt.Printf("Detaching servers from private network '%s'...\n", test_pn.Name)
	prn, err := api.DetachPrivateNetworkServer(test_pn.Id, test_server.Id)

	if err != nil {
		t.Errorf("DetachPrivateNetworkServer failed. Error: " + err.Error())
		return
	}

	//	prn = wait_for_state(prn, 10, 60, "ACTIVE")
	api.WaitForState(prn, "ACTIVE", 10, 60)
	prn, err = api.GetPrivateNetwork(prn.Id)
	if prn == nil || err != nil {
		t.Errorf("Detaching servers from the private network failed.")
		return
	}
	if len(prn.Servers) > 0 {
		t.Errorf("No all servers detached from the private network.")
	}
}

func TestUpdatePrivateNetwork(t *testing.T) {
	set_pn.Do(set_private_network)

	fmt.Printf("Updating private network '%s'...\n", test_pn.Id)
	new_name := test_pn.Name + "_updated"
	new_desc := test_pn.Description + "_updated"
	pnset := PrivateNetworkRequest{
		Name:           new_name,
		Description:    new_desc,
		NetworkAddress: "192.168.7.0",
		SubnetMask:     "255.255.255.0",
	}
	prn, err := api.UpdatePrivateNetwork(test_pn.Id, &pnset)

	if err != nil {
		t.Errorf("UpdatePrivateNetwork failed. Error: " + err.Error())
		return
	}
	//	prn = wait_for_state(prn, 5, 60, "ACTIVE")
	api.WaitForState(prn, "ACTIVE", 5, 60)
	prn, err = api.GetPrivateNetwork(prn.Id)
	if prn == nil || err != nil {
		t.Errorf("GetPrivateNetwork failed. Error: " + err.Error())
		return
	}
	if prn.Id != test_pn.Id {
		t.Errorf("Wrong private network ID.")
	}
	if prn.Name != new_name {
		t.Errorf("Wrong private network name.")
	}
	if prn.Description != new_desc {
		t.Errorf("Wrong private network description.")
	}
	if prn.NetworkAddress != "192.168.7.0" {
		t.Errorf("Wrong private network address.")
	}
	if prn.SubnetMask != "255.255.255.0" {
		t.Errorf("Wrong private network subnet mask.")
	}
	test_pn = prn
}

func TestDeletePrivateNetwork(t *testing.T) {
	set_pn.Do(set_private_network)

	fmt.Printf("Deleting private network '%s'...\n", test_pn.Name)
	prn, err := api.DeletePrivateNetwork(test_pn.Id)

	if err != nil {
		t.Errorf("DeletePrivateNetwork failed. Error: " + err.Error())
		return
	}

	api.WaitUntilDeleted(prn)
	prn, err = api.GetPrivateNetwork(prn.Id)

	if prn != nil {
		t.Errorf("Unable to delete the private network.")
	} else {
		test_pn = nil
	}
}
