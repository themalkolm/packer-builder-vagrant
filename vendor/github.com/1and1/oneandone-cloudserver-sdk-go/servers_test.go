package oneandone

import (
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

var (
	set_server  sync.Once
	set_dvd     sync.Once
	dvd_iso_id  string
	server_id   string
	server_name string
	ser_app_id  string
	server_ip   *ServerIp
	server_hdd  *Hdd
	server      *Server
	ser_pn      *PrivateNetwork
	ser_lb      *LoadBalancer
)

const (
	v_cores  = 1
	c_per_pr = 1
	ram      = 1
	hdd_size = 40
)

// Helper functions

func wait_for_action_done(srv *Server, sec time.Duration, count int) *Server {
	for i := 0; i < count; i++ {
		status, err := api.GetServerStatus(srv.Id)
		if err == nil {
			if status.Percent == 0 {
				srv, _ = api.GetServer(srv.Id)
				return srv
			}
		}
		time.Sleep(sec * time.Second)
	}
	return srv
}

func setup_server() {
	fmt.Println("Deploying a test server...")
	srv_id, srv, err := create_test_server(false)

	if err != nil {
		fmt.Printf("Unable to create the server. Error: %s", err.Error())
		return
	}
	if srv_id == "" || srv.Id == "" {
		fmt.Printf("Unable to create the server.")
		return
	} else {
		server_id = srv.Id
	}

	err = api.WaitForState(srv, "POWERED_OFF", 10, 90)

	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	server = srv
}

func get_random_datacenterID() string {
	rand.Seed(time.Now().UnixNano())
	dcs, _ := api.ListDatacenters()
	if len(dcs) > 0 {
		i := rand.Intn(len(dcs))
		return dcs[i].Id
	}
	return ""
}

func get_random_appliance(max_disk_size int) ServerAppliance {
	rand.Seed(time.Now().UnixNano())
	saps, _ := api.ListServerAppliances()
	for {
		i := rand.Intn(len(saps))
		if saps[i].MinHddSize <= max_disk_size && saps[i].Type == "IMAGE" && !strings.Contains(saps[i].Name, "2008") {
			return saps[i]
		}
	}
}

func get_default_mon_policy() MonitoringPolicy {
	mps, _ := api.ListMonitoringPolicies(0, 0, "default", "", "id,default")
	return mps[len(mps)-1]
}

func create_test_server(power_on bool) (string, *Server, error) {
	rand.Seed(time.Now().UnixNano())
	server_name = fmt.Sprintf("TestServer_%d", rand.Intn(1000000))
	fmt.Printf("Creating test server '%s'...\n", server_name)

	sap := get_random_appliance(hdd_size)
	ser_app_id = sap.Id
	mp := get_default_mon_policy()

	req := ServerRequest{
		Name:               server_name,
		Description:        server_name + " description",
		ApplianceId:        ser_app_id,
		MonitoringPolicyId: mp.Id,
		PowerOn:            power_on,
		Hardware: Hardware{
			Vcores:            v_cores,
			CoresPerProcessor: c_per_pr,
			Ram:               ram,
			Hdds: []Hdd{
				Hdd{
					Size:   hdd_size,
					IsMain: true,
				},
			},
		},
	}
	ser_id, server, err := api.CreateServer(&req)
	return ser_id, server, err
}

func load_server_dvd(ser_id string) {
	dvds, _ := api.ListDvdIsos()
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(len(dvds))

	fmt.Printf("Loading dvd '%s' in server '%s' virtula drive...\n", dvds[i].Name, server.Name)

	srv, err := api.LoadServerDvd(ser_id, dvds[i].Id)

	if err != nil {
		fmt.Printf("Loading server's dvd failed. Error: " + err.Error())
		return
	}

	for srv.Dvd == nil || srv.Status.Percent != 0 {
		srv = wait_for_action_done(srv, 10, 90)
	}
	dvd_iso_id = dvds[i].Id
	server = srv
}

// /servers tests

func TestCreateServer(t *testing.T) {
	set_server.Do(setup_server)

	if server == nil {
		t.Errorf("CreateServer failed.")
		return
	}
	if server.Name != server_name {
		t.Errorf("Wrong server name.")
	}
	if server.Image.Id != ser_app_id {
		t.Errorf("Wrong server image on server '%s'.", server.Name)
	}
	if server.Hardware.Vcores != v_cores {
		t.Errorf("Wrong number of processor cores on server '%s'.", server.Name)
	}
	if server.Hardware.CoresPerProcessor != c_per_pr {
		t.Errorf("Wrong number of cores per processor on server '%s'.", server.Name)
	}
	if server.Hardware.Ram != ram {
		t.Errorf("Wrong RAM size on server '%s'.", server.Name)
	}
}

func TestCreateServerEx(t *testing.T) {
	fmt.Println("Creating a fixed-size server...")
	var size_s FixedInstanceInfo
	fixed_flavours, _ := api.ListFixedInstanceSizes()
	for _, fl := range fixed_flavours {
		if fl.Name == "S" {
			size_s = fl
			break
		}
	}
	sap := get_random_appliance(size_s.Hardware.Hdds[0].Size)

	req := ServerRequest{
		DatacenterId: get_random_datacenterID(),
		Name:         "Random S Server",
		ApplianceId:  sap.Id,
		Hardware: Hardware{
			FixedInsSizeId: size_s.Id,
		},
	}
	ip, password, err := api.CreateServerEx(&req, 1800)

	if ip == "" {
		t.Errorf("CreateServerEx failed. Server IP address cannot be blank.")
	}
	if password == "" {
		t.Errorf("CreateServerEx failed. Password cannot be blank.")
	}
	if err != nil {
		t.Errorf("CreateServerEx failed. Error: " + err.Error())
		return
	}
	srvs, _ := api.ListServers(0, 0, "", ip, "")
	if len(srvs) > 0 {
		api.DeleteServer(srvs[0].Id, false)
	}
}

func TestListServers(t *testing.T) {
	set_server.Do(setup_server)
	fmt.Println("Listing all servers...")

	res, err := api.ListServers()
	if err != nil {
		t.Errorf("ListServers failed. Error: " + err.Error())
	}
	if len(res) == 0 {
		t.Errorf("No server found.")
	}

	res, err = api.ListServers(1, 2, "name", "", "id,name")

	if err != nil {
		t.Errorf("ListServers with parameter options failed. Error: " + err.Error())
		return
	}
	if len(res) == 0 {
		t.Errorf("No server found.")
	}
	if len(res) > 2 {
		t.Errorf("Wrong number of objects per page.")
	}
	if res[0].Hardware != nil {
		t.Errorf("Filtering parameters failed.")
	}
	if res[0].Name == "" {
		t.Errorf("Filtering parameters failed.")
	}
	if len(res) == 2 && res[0].Name >= res[1].Name {
		t.Errorf("Sorting parameters failed.")
	}
	// Test for error response
	res, err = api.ListServers(0, 0, true, "name", "")
	if res != nil || err == nil {
		t.Errorf("ListServers failed to handle incorrect argument type.")
	}

	res, err = api.ListServers(0, 0, "", server_name, "")

	if err != nil {
		t.Errorf("ListServers with parameter options failed. Error: " + err.Error())
		return
	}
	if len(res) != 1 {
		t.Errorf("Search parameter failed.")
	}
	if res[0].Name != server_name {
		t.Errorf("Search parameter failed.")
	}
}

func TestListFixedInstanceSizes(t *testing.T) {
	fmt.Println("Listing fixed instance sizes...")

	res, err := api.ListFixedInstanceSizes()

	if err != nil {
		t.Errorf("ListFixedInstanceSizes failed. Error: " + err.Error())
		return
	}

	exp_type := reflect.Slice
	rec_type := reflect.TypeOf(res).Kind()

	if rec_type != exp_type {
		t.Errorf("Expected type: " + exp_type.String() + ", received: " + rec_type.String())
	}
}

func TestGetFixedInstanceSize(t *testing.T) {
	fmt.Println("Getting a fixed instance size...")

	res, err := api.ListFixedInstanceSizes()

	if err != nil {
		t.Errorf("GetFixedInstanceSize failed. Error: " + err.Error())
	}

	if len(res) > 0 {
		fix_in_id := res[0].Id
		res, err := api.GetFixedInstanceSize(fix_in_id)

		if err != nil {
			fmt.Println("GetFixedInstanceSize failed. Error: " + err.Error())
			t.Errorf("GetFixedInstanceSize failed. Error: " + err.Error())
		}
		if res.Id != fix_in_id {
			t.Errorf("Wrong fixed instance ID.")
		}
		if res.Hardware == nil {
			t.Errorf("Missing fixed instance hardware info.")
		}
	} else {
		t.Errorf("Empty fixed instance sizes list.")
	}
}

func TestGetServer(t *testing.T) {
	set_server.Do(setup_server)

	fmt.Println("Getting the server...")
	srv, err := api.GetServer(server_id)

	if err != nil {
		t.Errorf("GetServer failed. Error: " + err.Error())
		return
	}
	if srv.Id != server_id {
		t.Errorf("Wrong server ID.")
	}
}

func TestGetServerStatus(t *testing.T) {
	set_server.Do(setup_server)

	fmt.Println("Getting the server's status...")
	status, err := api.GetServerStatus(server_id)

	if err != nil {
		t.Errorf("GetServerStatus failed. Error: " + err.Error())
		return
	}
	if status.State != "POWERED_OFF" {
		t.Errorf("Wrong server status.")
	}
}

func TestStartServer(t *testing.T) {
	set_server.Do(setup_server)

	fmt.Println("Starting the server...")
	srv, err := api.StartServer(server_id)

	if err != nil {
		t.Errorf("StartServer failed. Error: " + err.Error())
		return
	}

	err = api.WaitForState(srv, "POWERED_ON", 10, 60)

	if err != nil {
		t.Errorf("Starting the server failed. Error: " + err.Error())
	}

	server, err = api.GetServer(srv.Id)

	if err != nil {
		t.Errorf("GetServer failed. Error: " + err.Error())
	} else if server.Status.State != "POWERED_ON" {
		t.Errorf("Wrong server state. Expected: POWERED_ON. Found: %s.", server.Status.State)
	}
}

func TestRebootServer(t *testing.T) {
	set_server.Do(setup_server)

	for i := 1; i < 3; i++ {
		is_hardware := i%2 == 0
		var method string
		if is_hardware {
			method = "HARDWARE"
		} else {
			method = "SOFTWARE"
		}
		fmt.Printf("Rebooting the server using '%s' method...\n", method)
		srv, err := api.RebootServer(server_id, is_hardware)

		if err != nil {
			t.Errorf("RebootServer using '%s' method failed. Error: %s", method, err.Error())
			return
		}

		err = api.WaitForState(srv, "REBOOTING", 10, 60)

		if err != nil {
			t.Errorf("Rebooting the server using '%s' method failed. Error:  %s", method, err.Error())
		}

		err = api.WaitForState(srv, "POWERED_ON", 10, 60)

		if err != nil {
			t.Errorf("Rebooting the server using '%s' method failed. Error:  %s", method, err.Error())
		}
	}
}

func TestRenameServer(t *testing.T) {
	set_server.Do(setup_server)

	fmt.Println("Renaming the server...")

	new_name := server.Name + "_renamed"
	new_desc := server.Description + "_renamed"

	srv, err := api.RenameServer(server_id, new_name, new_desc)

	if err != nil {
		t.Errorf("Renaming server failed. Error: " + err.Error())
		return
	}
	if srv.Name != new_name {
		t.Errorf("Wrong server name.")
	}
	if srv.Description != new_desc {
		t.Errorf("Wrong server description.")
	}
}

func TestLoadServerDvd(t *testing.T) {
	set_server.Do(setup_server)
	set_dvd.Do(func() { load_server_dvd(server_id) })

	if server.Dvd == nil {
		t.Errorf("No dvd loaded.")
	} else if server.Dvd.Id != dvd_iso_id {
		t.Errorf("Wrong dvd loaded.")
	}
}

func TestGetServerDvd(t *testing.T) {
	set_server.Do(setup_server)
	set_dvd.Do(func() { load_server_dvd(server_id) })

	fmt.Printf("Getting server '%s' virtual dvd image...\n", server.Name)
	dvd, err := api.GetServerDvd(server_id)

	if err != nil {
		t.Errorf("GetServerDvd failed. Error: " + err.Error())
		return
	}
	time.Sleep(time.Second)
	server, _ = api.GetServer(server_id)
	if dvd == nil || server.Dvd == nil {
		t.Errorf("No dvd loaded.")
	} else if dvd.Id != server.Dvd.Id {
		t.Errorf("Wrong dvd loaded.")
	}
}

func TestEjectServerDvd(t *testing.T) {
	set_server.Do(setup_server)
	set_dvd.Do(func() { load_server_dvd(server_id) })

	fmt.Printf("Ejecting server '%s' virtual dvd drive...\n", server.Name)
	srv, err := api.EjectServerDvd(server_id)

	if err != nil {
		t.Errorf("EjectServerDvd failed. Error: " + err.Error())
		return
	}

	for srv.Dvd != nil || srv.Status.Percent != 0 {
		srv = wait_for_action_done(srv, 10, 60)
	}

	if srv.Dvd != nil {
		t.Errorf("Unable to eject the server's dvd.")
	}
}

func TestGetServerHardware(t *testing.T) {
	set_server.Do(setup_server)

	fmt.Println("Getting the server's hardware...")
	hardware, err := api.GetServerHardware(server_id)

	if err != nil {
		t.Errorf("GetServerHardware failed. Error: " + err.Error())
		return
	}
	if hardware == nil {
		t.Errorf("Unable to get the server's hardware.")
	} else {
		if hardware.Vcores != server.Hardware.Vcores {
			t.Errorf("Wrong number of processor cores on server '%s'.", server.Name)
		}
		if hardware.CoresPerProcessor != server.Hardware.CoresPerProcessor {
			t.Errorf("Wrong number of cores per processor on server '%s'.", server.Name)
		}
		if hardware.Ram != server.Hardware.Ram {
			t.Errorf("Wrong RAM size on server '%s'.", server.Name)
		}
		if len(hardware.Hdds) == 0 {
			t.Errorf("Missing HDD on server '%s'.", server.Name)
		}
	}
}

func TestListServerHdds(t *testing.T) {
	set_server.Do(setup_server)

	fmt.Println("Listing all the server's HDDs...")
	hdds, err := api.ListServerHdds(server_id)

	if err != nil {
		t.Errorf("ListServerHdds failed. Error: " + err.Error())
		return
	}
	if len(hdds) != 1 {
		t.Errorf("Wrong number of the server's hard disks.")
	}
	if hdds[0].Id == "" {
		t.Errorf("Wrong HDD id.")
	}
	if hdds[0].Size != 40 {
		t.Errorf("Wrong HDD size.")
	}
	if hdds[0].IsMain != true {
		t.Errorf("Wrong main HDD.")
	}
	server_hdd = &hdds[0]
}

func TestGetServerHdd(t *testing.T) {
	set_server.Do(setup_server)
	hdds, _ := api.ListServerHdds(server_id)

	fmt.Println("Getting server HDD...")
	hdd, err := api.GetServerHdd(server_id, hdds[0].Id)

	if err != nil {
		t.Errorf("GetServerHdd failed. Error: " + err.Error())
		return
	}
	if hdd.Id != hdds[0].Id {
		t.Errorf("Wrong HDD id.")
	}
	if hdd.Size != hdds[0].Size {
		t.Errorf("Wrong HDD size.")
	}
	if hdd.IsMain != hdds[0].IsMain {
		t.Errorf("Wrong main HDD.")
	}
}

func TestResizeServerHdd(t *testing.T) {
	set_server.Do(setup_server)
	hdds, _ := api.ListServerHdds(server_id)

	fmt.Println("Resizing the server's HDD...")
	srv, err := api.ResizeServerHdd(server_id, hdds[0].Id, 50)

	if err != nil {
		t.Errorf("ResizeServerHdd failed. Error: " + err.Error())
		return
	}

	srv = wait_for_action_done(srv, 10, 30)

	if err != nil {
		t.Errorf("GetServer failed. Error: " + err.Error())
	} else {
		if srv.Hardware.Hdds[0].Id != hdds[0].Id {
			t.Errorf("Wrong HDD id.")
		}
		if srv.Hardware.Hdds[0].Size != 50 {
			t.Errorf("HDD not resized.")
		}
	}
}

func TestAddServerHdds(t *testing.T) {
	set_server.Do(setup_server)
	hdds := ServerHdds{
		Hdds: []Hdd{
			{
				Size:   20,
				IsMain: false,
			},
		},
	}
	fmt.Println("Adding a HDD to the server...")
	srv, err := api.AddServerHdds(server_id, &hdds)

	if err != nil {
		t.Errorf("AddServerHdds failed. Error: " + err.Error())
		return
	}

	srv = wait_for_action_done(srv, 10, 120)

	if len(srv.Hardware.Hdds) != 2 {
		t.Errorf("Wrong number of hard disks.")
	}

	var new_hdd *Hdd
	if srv.Hardware.Hdds[0].Id != server_hdd.Id {
		new_hdd = &srv.Hardware.Hdds[0]
	} else {
		new_hdd = &srv.Hardware.Hdds[1]
	}
	if new_hdd.Size != hdds.Hdds[0].Size {
		t.Errorf("Wrong HDD size.")
	}
	if new_hdd.IsMain != hdds.Hdds[0].IsMain {
		t.Errorf("Wrong main HDD.")
	}

	server_hdd = new_hdd
}

func TestDeleteServerHdd(t *testing.T) {
	set_server.Do(setup_server)

	fmt.Println("Deleting the server's HDD...")
	srv, err := api.DeleteServerHdd(server_id, server_hdd.Id)

	if err != nil {
		t.Errorf("DeleteServerHdd failed. Error: " + err.Error())
		return
	}
	srv = wait_for_action_done(srv, 10, 90)
	if len(srv.Hardware.Hdds) != 1 {
		t.Errorf("Wrong number of the server's hard disks. The HDD was not deleted.")
	}

	server_hdd = &srv.Hardware.Hdds[0]
}

func TestGetServerImage(t *testing.T) {
	set_server.Do(setup_server)

	fmt.Println("Getting the server's image...")
	img, err := api.GetServerImage(server_id)

	if err != nil {
		t.Errorf("GetServerImage failed. Error: " + err.Error())
		return
	}
	if img.Id != server.Image.Id {
		t.Errorf("Wrong image ID.")
	}
	if img.Name != server.Image.Name {
		t.Errorf("Wrong image name.")
	}
}

func TestReinstallServerImage(t *testing.T) {
	set_server.Do(setup_server)
	sap := get_random_appliance(hdd_size)
	fps, _ := api.ListFirewallPolicies(0, 0, "creation_date", sap.OsFamily, "id,name,default")
	fp_id := ""
	for _, fp := range fps {
		if fp.DefaultPolicy == 1 {
			fp_id = fp.Id
			break
		}
	}
	fmt.Printf("Reinstalling the server to '%s'...\n", sap.Name)
	srv, err := api.ReinstallServerImage(server_id, sap.Id, "", fp_id)
	if err != nil {
		t.Errorf("ReinstallServerImage failed. Error: " + err.Error())
	} else {
		err = api.WaitForState(srv, "POWERED_ON", 30, 120)
		if err != nil {
			t.Errorf("ReinstallServerImage failed. Error: " + err.Error())
		}
		time.Sleep(time.Second)
		srv, _ = api.GetServer(server_id)

		if srv.Image.Name != sap.Name {
			t.Errorf("Wrong image installed.")
		}
		if srv.Ips[0].Firewall == nil || srv.Ips[0].Firewall.Id != fp_id {
			t.Errorf("Failed to assign firewall policy '%s' to the server when reinstalling the server.", fp_id)
		}
	}
}

func TestStopServer(t *testing.T) {
	set_server.Do(setup_server)

	fmt.Println("Stopping the server...")
	srv, err := api.ShutdownServer(server_id, false)

	if err != nil {
		t.Errorf("ShutdownServer failed. Error: " + err.Error())
	} else {
		err = api.WaitForState(srv, "POWERED_OFF", 10, 60)
		if err != nil {
			t.Errorf("Stopping the server failed. Error: " + err.Error())
		}

		server, err = api.GetServer(server_id)
		if err != nil {
			t.Errorf("GetServer failed. Error: " + err.Error())
		}
		if server.Status.State != "POWERED_OFF" {
			t.Errorf("Wrong server state. Expected: POWERED_OFF. Found: %s.", server.Status.State)
		}
	}
}

func TestUpdateServerHardware(t *testing.T) {
	// The server should be POWERED_OFF for this test to pass.
	// Hot increase of CoresPerProcessor is not allowed.
	set_server.Do(setup_server)

	fmt.Println("Updating the server's hardware...")

	hw := Hardware{
		Vcores:            2,
		CoresPerProcessor: 2,
		Ram:               2,
	}

	srv, err := api.UpdateServerHardware(server_id, &hw)

	if err != nil {
		t.Errorf("UpdateServersHardware failed. Error: " + err.Error())
	} else {
		srv = wait_for_action_done(srv, 10, 90)
		if srv.Hardware.Vcores != hw.Vcores {
			t.Errorf("Wrong number of processor cores. Expected: %d ; Found: %d", hw.Vcores, srv.Hardware.Vcores)
		}
		if srv.Hardware.CoresPerProcessor != hw.CoresPerProcessor {
			t.Errorf("Wrong number of cores per processor. Expected: %d ; Found: %d", hw.CoresPerProcessor, srv.Hardware.CoresPerProcessor)
		}
		if srv.Hardware.Ram != hw.Ram {
			t.Errorf("Wrong RAM size. Expected: %f ; Found: %f", hw.Ram, srv.Hardware.Ram)
		}
		server = srv
	}
}

func TestCreateServerSnapshot(t *testing.T) {
	set_server.Do(setup_server)

	fmt.Println("Creating a new snapshot for the server...")
	srv, err := api.CreateServerSnapshot(server_id)

	if err != nil {
		t.Errorf("CreateServerSnapshot failed. Error: " + err.Error())
	} else {
		if srv.Snapshot == nil {
			t.Errorf("No snapshot created.")
		} else {
			time.Sleep(180 * time.Second)
		}
	}
}

func TestGetServerSnapshot(t *testing.T) {
	set_server.Do(setup_server)

	fmt.Println("Getting the server's snapshots...")
	ss, err := api.GetServerSnapshot(server_id)

	if err != nil {
		t.Errorf("GetServerSnapshot failed. Error: " + err.Error())
	}
	if ss == nil {
		t.Errorf("No server snapshot found.")
	}
}

func TestRestoreServerSnapshot(t *testing.T) {
	set_server.Do(setup_server)
	ss, _ := api.GetServerSnapshot(server_id)

	fmt.Println("Restoring the server's snapshot...")
	_, err := api.RestoreServerSnapshot(server_id, ss.Id)

	if err != nil {
		t.Errorf("RestoreServerSnapshot failed. Error: " + err.Error())
	} else {
		time.Sleep(180 * time.Second)
	}
}

func TestDeleteServerSnapshot(t *testing.T) {
	set_server.Do(setup_server)
	ss, _ := api.GetServerSnapshot(server_id)

	fmt.Println("Deleting the server's snapshot...")
	srv, err := api.DeleteServerSnapshot(server_id, ss.Id)

	if err != nil {
		t.Errorf("DeleteServerSnapshot failed. Error: " + err.Error())
	} else {
		time.Sleep(180 * time.Second)
		srv, _ = api.GetServer(server_id)
		if srv.Snapshot != nil {
			t.Errorf("The snapshot was not deleted.")
		}
	}
}

func TestListServerIps(t *testing.T) {
	set_server.Do(setup_server)

	srv, e := api.GetServer(server_id)
	if e == nil {
		server = srv
	}

	fmt.Println("Listing the server's IPs...")
	ips, err := api.ListServerIps(server_id)

	if err != nil {
		t.Errorf("ListServerIps failed. Error: " + err.Error())
		return
	}
	if len(ips) != len(server.Ips) {
		t.Errorf("Not all IPs were obtained.")
	}
	if ips[0].Id != server.Ips[0].Id {
		t.Errorf("Wrong IP ID.")
	}
	if ips[0].Ip != server.Ips[0].Ip {
		t.Errorf("Wrong IP address.")
	}
}

func TestAssignServerIp(t *testing.T) {
	set_server.Do(setup_server)

	fmt.Println("Assigning new IP addresses to the server...")
	for i := 2; i < 4; i++ {
		time.Sleep(time.Second)
		srv, err := api.AssignServerIp(server_id, "IPV4")
		if err != nil {
			t.Errorf("AssignServerIp failed. Error: " + err.Error())
			return
		}
		srv = wait_for_action_done(srv, 10, 30)
		if len(srv.Ips) != i {
			t.Errorf("IP address not assigned to the server.")
		}
		server = srv
	}
}

func TestGetServerIp(t *testing.T) {
	set_server.Do(setup_server)

	fmt.Println("Getting the server's IP...")
	server, _ = api.GetServer(server_id)
	if server == nil {
		t.Errorf("GetServer failed.")
		return
	}
	time.Sleep(time.Second)
	ip, err := api.GetServerIp(server_id, server.Ips[0].Id)

	if err != nil {
		t.Errorf("GetServerIps failed. Error: " + err.Error())
		return
	}
	if ip.Id != server.Ips[0].Id {
		t.Errorf("Wrong IP ID.")
	}
	if ip.Ip != server.Ips[0].Ip {
		t.Errorf("Wrong IP address.")
	}
}

func TestDeleteServerIp(t *testing.T) {
	set_server.Do(setup_server)

	if len(server.Ips) <= 1 {
		for i := 0; i < 2; i++ {
			time.Sleep(time.Second)
			s, e := api.AssignServerIp(server_id, "IPV4")
			if s != nil && e == nil {
				s = wait_for_action_done(s, 10, 30)
				server = s
			}
		}
	}
	ip_no := len(server.Ips)
	for i := 1; i < ip_no; i++ {
		keep_ip := i%2 == 0
		fmt.Printf("Deleting the server's IP '%s' (keep_ip = %s)...\n", server.Ips[i].Ip, strconv.FormatBool(keep_ip))
		srv, err := api.DeleteServerIp(server_id, server.Ips[i].Id, keep_ip)

		if err != nil {
			t.Errorf("DeleteServerIp failed. Error: " + err.Error())
			return
		}
		if len(srv.Ips) != ip_no-i {
			t.Errorf("IP address '%s' is not removed from the server.", server.Ips[i].Ip)
		}
		ip, _ := api.GetPublicIp(server.Ips[i].Id)
		if keep_ip {
			if ip == nil {
				t.Errorf("Failed to keep public IP '%s' when removed from server.", server.Ips[i].Ip)
			} else {
				fmt.Printf("Deleting IP address '%s' after removing from the server...\n", server.Ips[i].Ip)
				ip, err = api.DeletePublicIp(ip.Id)
			}
		} else if ip != nil {
			t.Errorf("Failed to delete public IP '%s' when removed from server.", server.Ips[i].Ip)
			fmt.Printf("Cleaning up. Deleting IP address '%s' directly...\n", server.Ips[i].Ip)
			ip, err = api.DeletePublicIp(ip.Id)
		}
	}
}

func TestAssignServerPrivateNetwork(t *testing.T) {
	set_server.Do(setup_server)
	ser_pn = create_private_netwok()

	fmt.Println("Assigning the private network to the server...")
	srv, err := api.AssignServerPrivateNetwork(server_id, ser_pn.Id)

	if err != nil {
		t.Errorf("AssignServerPrivateNetwork failed. Error: " + err.Error())
		return
	}
	prn, _ := api.GetServerPrivateNetwork(server_id, ser_pn.Id)

	if len(srv.PrivateNets) == 0 {
		t.Errorf("The private network was not assigned to the server.")
	} else if srv.PrivateNets[0].Id != ser_pn.Id {
		t.Errorf("The private network was not assigned to the server.")
	}
	//	prn = wait_for_state(prn, 20, 30, "ACTIVE")
	api.WaitForState(prn, "ACTIVE", 10, 60)
	ser_pn, _ = api.GetPrivateNetwork(prn.Id)
}

func TestListServerPrivateNetworks(t *testing.T) {
	set_server.Do(setup_server)

	fmt.Println("Listing the server's private networks...")
	pns, err := api.ListServerPrivateNetworks(server_id)

	if err != nil {
		t.Errorf("ListServerPrivateNetworks failed. Error: " + err.Error())
		return
	}
	if len(pns) > 0 {
		if pns[0].Id == "" {
			t.Errorf("The private network ID is missing.")
		}
		if pns[0].Name == "" {
			t.Errorf("The private network name is missing.")
		}
	} else {
		t.Errorf("No private networks were assigned to the server.")
	}
}

func TestGetServerPrivateNetwork(t *testing.T) {
	set_server.Do(setup_server)

	fmt.Println("Getting the server's private network...")
	privn, err := api.GetServerPrivateNetwork(server_id, ser_pn.Id)

	if err != nil {
		t.Errorf("GetServerPrivateNetwork failed. Error: " + err.Error())
		return
	}
	if len(privn.Servers) == 0 {
		t.Errorf("The private network server is missing.")
	}

	found_srv := false
	for _, s := range privn.Servers {
		if s.Id == server_id {
			found_srv = true
			break
		}
	}

	if !found_srv {
		t.Errorf("Private network server is missing.")
	}
}

func TestRemoveServerPrivateNetwork(t *testing.T) {
	set_server.Do(setup_server)

	fmt.Println("Unassigning the private network from the server...")
	srv, err := api.RemoveServerPrivateNetwork(server_id, ser_pn.Id)

	if err != nil {
		t.Errorf("RemoveServerPrivateNetwork failed. Error: " + err.Error())
		return
	}

	prn, _ := api.GetServerPrivateNetwork(server_id, ser_pn.Id)
	//	prn = wait_for_state(prn, 20, 30, "ACTIVE")
	api.WaitForState(prn, "ACTIVE", 10, 60)
	srv, err = api.GetServer(server_id)

	if srv == nil || err != nil {
		t.Errorf("GetServer failed. Error: " + err.Error())
	} else if len(srv.PrivateNets) > 0 {
		t.Errorf("Private network not unassigned from the server.")
	}
	// cleanup
	api.DeletePrivateNetwork(ser_pn.Id)
}

func TestAssignServerIpLoadBalancer(t *testing.T) {
	set_server.Do(setup_server)
	ips, _ := api.ListServerIps(server_id)
	lb := create_load_balancer()

	fmt.Printf("Assigning a load balancer to the server's IP '%s'...\n", ips[0].Ip)
	srv, err := api.AssignServerIpLoadBalancer(server_id, ips[0].Id, lb.Id)

	if err != nil {
		t.Errorf("AssignServerIpLoadBalancer failed. Error: " + err.Error())
		return
	}
	if len(srv.Ips[0].LoadBalancers) == 0 {
		t.Errorf("Load balancer not assigned.")
	}
	if srv.Ips[0].LoadBalancers[0].Id != lb.Id {
		t.Errorf("Wrong load balancer assigned.")
	}
	ser_lb = lb
}

func TestListServerIpLoadBalancers(t *testing.T) {
	set_server.Do(setup_server)
	ips, _ := api.ListServerIps(server_id)

	fmt.Println("Listing load balancers assigned to the server's IP...")
	lbs, err := api.ListServerIpLoadBalancers(server_id, ips[0].Id)

	if err != nil {
		t.Errorf("ListServerIpLoadBalancers failed. Error: " + err.Error())
		return
	}
	if len(lbs) == 0 {
		t.Errorf("No load balancer was assigned to the server's IP.")
	}
	if lbs[0].Id != ser_lb.Id {
		t.Errorf("Wrong load balancer assigned.")
	}
}

func TestUnassignServerIpLoadBalancer(t *testing.T) {
	set_server.Do(setup_server)
	ips, _ := api.ListServerIps(server_id)

	fmt.Println("Unassigning the load balancer from the server's IP...")
	srv, err := api.UnassignServerIpLoadBalancer(server_id, ips[0].Id, ser_lb.Id)

	if err != nil {
		t.Errorf("UnassignServerIpLoadBalancer failed. Error: " + err.Error())
		return
	}
	if len(srv.Ips[0].LoadBalancers) > 0 {
		t.Errorf("Unassigning the load balancer failed.")
	}
	ser_lb, err = api.DeleteLoadBalancer(ser_lb.Id)
	if err == nil {
		api.WaitUntilDeleted(ser_lb)
	}
	ser_lb, _ = api.GetLoadBalancer(ser_lb.Id)
}

func TestAssignServerIpFirewallPolicy(t *testing.T) {
	set_server.Do(setup_server)
	ips, _ := api.ListServerIps(server_id)

	fmt.Println("Assigning a firewall policy to the server's IP...")
	fps, err := api.ListFirewallPolicies(0, 1, "creation_date", "linux", "id,name")
	if err != nil {
		t.Errorf("ListFirewallPolicies failed. Error: " + err.Error())
		return
	}
	srv, err := api.AssignServerIpFirewallPolicy(server_id, ips[0].Id, fps[0].Id)

	if err != nil {
		t.Errorf("AssignServerIpFirewallPolicy failed. Error: " + err.Error())
		return
	}
	if srv.Ips[0].Firewall == nil {
		t.Errorf("Firewall policy not assigned.")
	}
	if srv.Ips[0].Firewall.Id != fps[0].Id {
		t.Errorf("Wrong firewall policy assigned.")
	}
}

func TestGetServerIpFirewallPolicy(t *testing.T) {
	set_server.Do(setup_server)
	ips, _ := api.ListServerIps(server_id)

	fmt.Println("Getting the firewall policy assigned to the server's IP...")
	fps, err := api.ListFirewallPolicies(0, 1, "creation_date", "linux", "id,name")
	if err != nil {
		t.Errorf("ListFirewallPolicies failed. Error: " + err.Error())
	}
	fp, err := api.GetServerIpFirewallPolicy(server_id, ips[0].Id)

	if err != nil {
		t.Errorf("GetServerIpFirewallPolicy failed. Error: " + err.Error())
	}
	if fp == nil {
		t.Errorf("No firewall policy assigned to the server's IP.")
	}
	if fp.Id != fps[0].Id {
		t.Errorf("Wrong firewall policy assigned to the server's IP.")
	}
}

func TestUnassignServerIpFirewallPolicy(t *testing.T) {
	set_server.Do(setup_server)
	ips, _ := api.ListServerIps(server_id)

	fmt.Println("Unassigning the firewall policy from the server's IP...")
	srv, err := api.UnassignServerIpFirewallPolicy(server_id, ips[0].Id)

	if err != nil {
		t.Errorf("UnassignServerIpFirewallPolicy failed. Error: " + err.Error())
		return
	}
	if srv.Ips[0].Firewall != nil {
		t.Errorf("Unassigning the firewall policy failed.")
	}
}

func TestCloneServer(t *testing.T) {
	set_server.Do(setup_server)

	fmt.Println("Cloning the server...")
	new_name := server_name + "_Copy"
	var dc_id string
	if server != nil && server.Datacenter != nil {
		dc_id = server.Datacenter.Id
	}
	srv, err := api.CloneServer(server_id, new_name, dc_id)

	if err != nil {
		t.Errorf("CloneServer failed. Error: " + err.Error())
	} else {
		if srv.Name != new_name {
			t.Errorf("Cloning the server failed. Wrong server name.")
		}
		if srv.Hardware.Vcores != server.Hardware.Vcores {
			t.Errorf("Cloning the server failed. Vcores values differ.")
		}
		if srv.Hardware.CoresPerProcessor != server.Hardware.CoresPerProcessor {
			t.Errorf("Cloning the server failed. CoresPerProcessor values differ.")
		}
		if srv.Hardware.Ram != server.Hardware.Ram {
			t.Errorf("Cloning the server failed. Ram values differ.")
		}
		if srv.Image.Id != server.Image.Id {
			t.Errorf("Cloning the server failed. Wrng server image.")
		}

		time.Sleep(500 * time.Second)
		fmt.Println("Deleting the clone...")
		srv, err = api.DeleteServer(srv.Id, false)

		if err != nil {
			t.Errorf("Unable to delete cloned server. Error: " + err.Error())
		}
	}
}

func TestDeleteServer(t *testing.T) {
	set_server.Do(setup_server)

	srv, err := api.DeleteServer(server_id, true)
	fmt.Printf("Deleting server '%s', keeping server's IP '%s'...\n", srv.Name, srv.Ips[0].Ip)
	ip_id := srv.Ips[0].Id

	if err != nil {
		t.Errorf("DeleteServer server failed. Error: " + err.Error())
		return
	}

	err = api.WaitUntilDeleted(srv)

	if err != nil {
		t.Errorf("Deleting the server failed. Error: " + err.Error())
	}

	srv, err = api.GetServer(server_id)

	if srv != nil {
		t.Errorf("Unable to delete the server.")
	} else {
		server = nil
	}

	ip, _ := api.GetPublicIp(ip_id)
	if ip == nil {
		t.Errorf("Failed to keep IP after deleting the server.")
	} else {
		fmt.Printf("Deleting server's IP '%s' after deleting the server...\n", ip.IpAddress)
		ip, err = api.DeletePublicIp(ip_id)
		if err != nil {
			t.Errorf("Unable to delete server's IP after deleting the server.")
		} else {
			api.WaitUntilDeleted(ip)
		}
	}
}
