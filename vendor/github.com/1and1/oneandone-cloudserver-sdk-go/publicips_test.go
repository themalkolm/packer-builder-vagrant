package oneandone

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

var (
	set_ip  sync.Once
	test_ip *PublicIp
)

const (
	ip_dns = "test.oneandone.com"
)

// Helper functions

func wait_for_ip_ready(ip_address *PublicIp, sec time.Duration, count int, state string) *PublicIp {
	if ip_address.State == state {
		return ip_address
	}
	for i := 0; i < count; i++ {
		ip, err := api.GetPublicIp(ip_address.Id)
		if err == nil {
			if ip.State == state {
				return ip
			}
		}
		time.Sleep(sec * time.Second)
	}
	return ip_address
}

func create_public_ip() *PublicIp {
	fmt.Printf("Creating an IPV4 public ip...\n")
	pip_id, pip, err := api.CreatePublicIp(IpTypeV4, ip_dns, "")
	if err != nil {
		fmt.Printf("Unable to create a public ip address. Error: %s", err.Error())
		return nil
	}
	if pip_id == "" || pip.Id == "" {
		fmt.Printf("Unable to create a public ip address.")
		return nil
	}
	return wait_for_ip_ready(pip, 5, 30, "ACTIVE")
}

func set_public_ip() {
	test_ip = create_public_ip()
}

// /public_ips tests

func TestCreatePublicIp(t *testing.T) {
	set_ip.Do(set_public_ip)

	if test_ip == nil {
		t.Errorf("CreatePublicIp failed.")
	}
	if test_ip.IpAddress == "" {
		t.Errorf("Missing IP address.")
	}
	if test_ip.Type != IpTypeV4 {
		t.Errorf("Wrong IP type.")
	}
	if test_ip.ReverseDns != ip_dns {
		t.Errorf("Wrong reverse dns of ip address '%s'.", test_ip.IpAddress)
	}
}

func TestGetPublicIp(t *testing.T) {
	set_ip.Do(set_public_ip)

	fmt.Printf("Getting public ip '%s'...\n", test_ip.IpAddress)
	ip, err := api.GetPublicIp(test_ip.Id)

	if err != nil {
		t.Errorf("GetPublicIp failed. Error: " + err.Error())
	}
	if ip.IpAddress != test_ip.IpAddress {
		t.Errorf("Wrong IP address.")
	}
	if ip.Type != test_ip.Type {
		t.Errorf("Wrong IP type.")
	}
	if ip.ReverseDns != test_ip.ReverseDns {
		t.Errorf("Wrong reverse dns of ip address '%s'.", ip.IpAddress)
	}
}

func TestListPublicIps(t *testing.T) {
	set_ip.Do(set_public_ip)
	fmt.Println("Listing all public ip addresses...")

	ips, err := api.ListPublicIps()
	if err != nil {
		t.Errorf("ListPublicIps failed. Error: " + err.Error())
	}
	if len(ips) == 0 {
		t.Errorf("No public ip found.")
	}

	ips, err = api.ListPublicIps(1, 3, "id", "", "id,ip")
	if err != nil {
		t.Errorf("ListPublicIps with parameter options failed. Error: " + err.Error())
	}
	if len(ips) == 0 {
		t.Errorf("No public ip found.")
	}
	if len(ips) > 3 {
		t.Errorf("Wrong number of objects per page.")
	}
	if ips[0].Id == "" {
		t.Errorf("Filtering parameters failed.")
	}
	if ips[0].IpAddress == "" {
		t.Errorf("Filtering parameters failed.")
	}
	if ips[0].State != "" {
		t.Errorf("Filtering parameters failed.")
	}
	if len(ips) >= 2 && ips[0].Id >= ips[1].Id {
		t.Errorf("Sorting parameters failed.")
	}
	// Test for error response
	ips, err = api.ListPublicIps(0, 0, "", "", false)
	if ips != nil || err == nil {
		t.Errorf("ListPublicIps failed to handle incorrect argument type.")
	}

	ips, err = api.ListPublicIps(0, 0, "", test_ip.IpAddress, "")
	if err != nil {
		t.Errorf("ListPublicIps with parameter options failed. Error: " + err.Error())
	}
	if len(ips) != 1 {
		t.Errorf("Search parameter failed.")
	}
	if ips[0].IpAddress != test_ip.IpAddress {
		t.Errorf("Search parameter failed.")
	}
}

func TestUpdatePublicIp(t *testing.T) {
	set_ip.Do(set_public_ip)

	fmt.Printf("Updating public ip '%s'...\n", test_ip.IpAddress)
	new_dns := "test.oneandone.de"

	ip, err := api.UpdatePublicIp(test_ip.Id, new_dns)

	if err != nil {
		t.Errorf("UpdatePublicIp failed. Error: " + err.Error())
	}
	ip = wait_for_ip_ready(ip, 5, 10, "ACTIVE")
	if ip.Id != test_ip.Id {
		t.Errorf("Wrong IP address ID.")
	}
	if ip.ReverseDns != new_dns {
		t.Errorf("Wrong reverse dns of ip address '%s'.", ip.IpAddress)
	}
}

func TestDeletePublicIp(t *testing.T) {
	set_ip.Do(set_public_ip)

	fmt.Printf("Deleting public ip '%s'...\n", test_ip.IpAddress)
	ip, err := api.DeletePublicIp(test_ip.Id)

	if err != nil {
		t.Errorf("DeletePublicIp failed. Error: " + err.Error())
	}

	ip = wait_for_ip_ready(ip, 5, 30, "REMOVING")
	time.Sleep(30 * time.Second)
	ip, _ = api.GetPublicIp(test_ip.Id)

	if ip != nil {
		t.Errorf("Unable to delete the public ip.")
	} else {
		test_ip = nil
	}
}
