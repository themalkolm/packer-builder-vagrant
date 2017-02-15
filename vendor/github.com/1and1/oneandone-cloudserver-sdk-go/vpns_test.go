package oneandone

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

var (
	set_vpn  sync.Once
	vpn_name string
	vpn_desc string
	test_vpn *VPN
)

const vpn_dc_code = "US"

// Helper functions

func create_vpn() *VPN {
	rand.Seed(time.Now().UnixNano())
	rint := rand.Intn(9999)
	vpn_name = fmt.Sprintf("Test_VPN_%d", rint)
	vpn_desc = fmt.Sprintf("Test_VPN_%d description", rint)

	fmt.Printf("Creating VPN '%s'...\n", vpn_name)
	vpn_id, vpn, err := api.CreateVPN(vpn_name, vpn_desc, "")
	if err != nil {
		fmt.Printf("Unable to create a VPN. Error: %s", err.Error())
		return nil
	}
	if vpn_id == "" || vpn.Id == "" {
		fmt.Printf("Unable to create VPN '%s'.", vpn_name)
		return nil
	}

	api.WaitForState(vpn, "ACTIVE", 10, 60)

	return vpn
}

func set_vpn_once() {
	test_vpn = create_vpn()
}

// /vpns tests

func TestCreateVPN(t *testing.T) {
	set_vpn.Do(set_vpn_once)

	if test_vpn == nil {
		t.Errorf("CreateVPN failed.")
	} else {
		if test_vpn.Name != vpn_name {
			t.Errorf("Wrong name of the VPN.")
		}
		if test_vpn.Description != vpn_desc {
			t.Errorf("Wrong VPN description.")
		}
		if test_vpn.Datacenter == nil {
			t.Errorf("Missing VPN Data Center.")
		} else if test_vpn.Datacenter.CountryCode != vpn_dc_code {
			t.Errorf("Wrong VPN Data Center.")
		}
	}
}

func TestGetVPN(t *testing.T) {
	set_vpn.Do(set_vpn_once)

	fmt.Printf("Getting VPN '%s'...\n", vpn_name)
	vpn, err := api.GetVPN(test_vpn.Id)

	if err != nil {
		t.Errorf("GetVPN failed. Error: " + err.Error())
		return
	}
	if vpn.Id != test_vpn.Id {
		t.Errorf("Wrong VPN ID.")
	}
	if len(test_vpn.IPs) == 0 {
		t.Errorf("Missing VPN IPs.")
	}
}

func TestListVPNs(t *testing.T) {
	set_vpn.Do(set_vpn_once)
	fmt.Println("Listing all VPNs...")

	res, err := api.ListVPNs()
	if err != nil {
		t.Errorf("ListVPNs failed. Error: " + err.Error())
	}
	if len(res) == 0 {
		t.Errorf("No VPN found.")
	}

	res, err = api.ListVPNs(1, 1, "", "", "id,name")

	if err != nil {
		t.Errorf("ListVPNs with parameter options failed. Error: " + err.Error())
		return
	}
	if len(res) == 0 {
		t.Errorf("No VPN found.")
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
	res, err = api.ListVPNs(0, 0, "", 1, "")
	if res != nil || err == nil {
		t.Errorf("ListVPNs failed to handle incorrect argument type.")
	}

	res, err = api.ListVPNs(0, 0, "", test_vpn.Name, "")

	if err != nil {
		t.Errorf("ListVPNs with parameter options failed. Error: " + err.Error())
		return
	}
	if len(res) != 1 {
		t.Errorf("Search parameter failed.")
	}
	if res[0].Name != test_vpn.Name {
		t.Errorf("Search parameter failed.")
	}
}

func TestGetVPNConfigFile(t *testing.T) {
	set_vpn.Do(set_vpn_once)

	fmt.Printf("Getting VPN's configuration...\n")
	base64_str, err := api.GetVPNConfigFile(test_vpn.Id)

	if err != nil {
		t.Errorf("GetVPNConfigFile failed. Error: " + err.Error())
		return
	}

	_, err = base64.StdEncoding.DecodeString(base64_str)

	if err != nil {
		t.Errorf("Unable to decode config file string. Error: " + err.Error())
	}
}

func TestModifyVPN(t *testing.T) {
	set_vpn.Do(set_vpn_once)

	fmt.Printf("Modifying VPN '%s'...\n", test_vpn.Id)
	new_name := test_vpn.Name + "_updated"
	new_desc := test_vpn.Description + " updated"

	vpn, err := api.ModifyVPN(test_vpn.Id, new_name, new_desc)

	if err != nil {
		t.Errorf("ModifyVPN failed. Error: " + err.Error())
		return
	}
	if vpn.Id != test_vpn.Id {
		t.Errorf("Wrong VPN ID.")
	}
	if vpn.Name != new_name {
		t.Errorf("Wrong VPN name.")
	}
	if vpn.Description != new_desc {
		t.Errorf("Wrong VPN description.")
	}

	test_vpn = vpn
}

func TestDeleteVPN(t *testing.T) {
	set_vpn.Do(set_vpn_once)

	fmt.Printf("Deleting VPN '%s'...\n", test_vpn.Name)
	vpn, err := api.DeleteVPN(test_vpn.Id)

	if err != nil {
		t.Errorf("DeleteVPN failed. Error: " + err.Error())
		return
	}

	api.WaitUntilDeleted(vpn)
	vpn, err = api.GetVPN(vpn.Id)

	if vpn != nil {
		t.Errorf("Unable to delete the VPN.")
	} else {
		test_vpn = nil
	}
}
