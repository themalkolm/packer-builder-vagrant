package oneandone

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

var (
	set_ss       sync.Once
	test_ss_name string
	test_ss_desc string
	test_ss      *SharedStorage
)

// Helper functions

func create_shared_storage() *SharedStorage {
	rand.Seed(time.Now().UnixNano())
	rint := rand.Intn(999)
	test_ss_name = fmt.Sprintf("SharedStorage_%d", rint)
	test_ss_desc = fmt.Sprintf("SharedStorage_%d description", rint)
	req := SharedStorageRequest{
		Name:        test_ss_name,
		Description: test_ss_desc,
		Size:        Int2Pointer(50),
	}
	fmt.Printf("Creating new shared storage '%s'...\n", test_ss_name)
	ss_id, ss, err := api.CreateSharedStorage(&req)
	if err != nil {
		fmt.Printf("Unable to create a shared storage. Error: %s", err.Error())
		return nil
	}
	if ss_id == "" || ss.Id == "" {
		fmt.Printf("Unable to create shared storage '%s'.", test_ss_name)
		return nil
	}

	api.WaitForState(ss, "ACTIVE", 10, 30)
	return ss
}

func set_shared_storage() {
	test_ss = create_shared_storage()
}

// /shared_storages tests

func TestCreateSharedStorage(t *testing.T) {
	set_ss.Do(set_shared_storage)

	if test_ss == nil {
		t.Errorf("CreateSharedStorage failed.")
		return
	}
	if test_ss.Id == "" {
		t.Errorf("Missing shared storage ID.")
	}
	if test_ss.Name != test_ss_name {
		t.Errorf("Wrong name of the shared storage.")
	}
	if test_ss.Description != test_ss_desc {
		t.Errorf("Wrong shared storage description.")
	}
	if test_ss.Size != 50 {
		t.Errorf("Wrong size of shared storage '%s'.", test_ss.Name)
	}
}

func TestGetSharedStorage(t *testing.T) {
	set_ss.Do(set_shared_storage)

	fmt.Printf("Getting shared storage '%s'...\n", test_ss.Name)
	ss, err := api.GetSharedStorage(test_ss.Id)

	if err != nil {
		t.Errorf("GetSharedStorage failed. Error: " + err.Error())
	} else {
		if ss.Id != test_ss.Id {
			t.Errorf("Wrong shared storage ID.")
		}
	}
}

func TestListSharedStorages(t *testing.T) {
	set_ss.Do(set_shared_storage)
	fmt.Println("Listing all shared storages...")

	res, err := api.ListSharedStorages()
	if err != nil {
		t.Errorf("ListSharedStorages failed. Error: " + err.Error())
	}
	if len(res) == 0 {
		t.Errorf("No shared storage found.")
	}

	res, err = api.ListSharedStorages(1, 1, "", "", "id,name")

	if err != nil {
		t.Errorf("ListSharedStorages with parameter options failed. Error: " + err.Error())
		return
	}
	if len(res) == 0 {
		t.Errorf("No shared storage found.")
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
	if res[0].CifsPath != "" {
		t.Errorf("Filtering parameters failed.")
	}
	if res[0].NfsPath != "" {
		t.Errorf("Filtering parameters failed.")
	}
	// Test for error response
	res, err = api.ListSharedStorages(true, false, "id", "name", "")
	if res != nil || err == nil {
		t.Errorf("ListSharedStorages failed to handle incorrect argument type.")
	}

	res, err = api.ListSharedStorages(0, 0, "", test_ss.Name, "")

	if err != nil {
		t.Errorf("ListSharedStorages with parameter options failed. Error: " + err.Error())
		return
	}
	if len(res) != 1 {
		t.Errorf("Search parameter failed.")
	}
	if res[0].Name != test_ss.Name {
		t.Errorf("Search parameter failed.")
	}
}

func TestGetSharedStorageCredentials(t *testing.T) {
	fmt.Printf("Getting access credentials for shared storage '%s'...\n", test_ss.Name)
	ssc, err := api.GetSharedStorageCredentials()

	if err != nil {
		t.Errorf("GetSharedStorageCredentials failed. Error: " + err.Error())
	} else {
		if ssc[0].KerberosContentFile == "" {
			t.Errorf("Missing Kerberos key.")
		}
		if ssc[0].UserDomain == "" {
			t.Errorf("Missing user domain.")
		}
	}
}

func TestUpdateSharedStorageCredentials(t *testing.T) {
	fmt.Printf("Updating credentials of shared storage '%s'...\n", test_ss.Name)
	rand.Seed(time.Now().UnixNano())
	rint := rand.Intn(999999)
	new_pass := fmt.Sprintf("N.u&A@V^d_g!b;e%d", rint)
	ssc, err := api.UpdateSharedStorageCredentials(new_pass)

	if err != nil {
		t.Errorf("UpdateSharedStorageCredentials failed. Error: " + err.Error())
		return
	}

	cr, _ := api.GetSharedStorageCredentials()

	if ssc[0].KerberosContentFile != cr[0].KerberosContentFile {
		t.Errorf("Credentials for shared storage not updated.")
	}
}

func TestAddSharedStorageServers(t *testing.T) {
	set_ss.Do(set_shared_storage)
	sync_server.Do(func() { deploy_test_server(false) })

	fmt.Printf("Adding server to shared storage '%s'...\n", test_ss.Name)

	servers := []SharedStorageServer{
		{
			Id:     test_server.Id,
			Rights: "RW",
		},
	}
	ss, err := api.AddSharedStorageServers(test_ss.Id, servers)

	if err != nil {
		t.Errorf("AddSharedStorageServers failed. Error: " + err.Error())
		return
	}

	api.WaitForState(ss, "ACTIVE", 10, 30)
	ss, _ = api.GetSharedStorage(ss.Id)

	if len(ss.Servers) != 1 {
		t.Errorf("Found no server added to the shared storage.")
	}
	if ss.Servers[0].Id != test_server.Id {
		t.Errorf("Wrong server added to the shared storage.")
	}
	test_ss = ss
}

func TestGetSharedStorageServer(t *testing.T) {
	set_ss.Do(set_shared_storage)

	fmt.Printf("Getting server added to shared storage '%s'...\n", test_ss.Name)
	ss_ser, err := api.GetSharedStorageServer(test_ss.Id, test_ss.Servers[0].Id)

	if err != nil {
		t.Errorf("GetSharedStorageServer failed. Error: " + err.Error())
		return
	}
	if ss_ser.Id != test_ss.Servers[0].Id {
		t.Errorf("Wrong ID of the server added to shared storage '%s'.", test_ss.Name)
	}
	if ss_ser.Name != test_ss.Servers[0].Name {
		t.Errorf("Wrong name of the server added to shared storage '%s'.", test_ss.Name)
	}
	if ss_ser.Rights != test_ss.Servers[0].Rights {
		t.Errorf("Wrong access rights of the server added to shared storage '%s'.", test_ss.Name)
	}
}

func TestListSharedStorageServers(t *testing.T) {
	set_ss.Do(set_shared_storage)
	sync_server.Do(func() { deploy_test_server(false) })

	fmt.Printf("Listing servers added to shared storage '%s'...\n", test_ss.Name)
	ss_srvs, err := api.ListSharedStorageServers(test_ss.Id)

	if err != nil {
		t.Errorf("ListSharedStorageServers failed. Error: " + err.Error())
		return
	}
	if len(ss_srvs) != 1 {
		t.Errorf("Wrong number of servers added to shared storage '%s'.", test_ss.Name)
	}
	if ss_srvs[0].Id != test_server.Id {
		t.Errorf("Wrong server added to shared storage '%s'.", test_ss.Name)
	}
}

func TestDeleteSharedStorageServer(t *testing.T) {
	set_ss.Do(set_shared_storage)
	sync_server.Do(func() { deploy_test_server(false) })

	fmt.Printf("Deleting server attached to shared storage '%s'...\n", test_ss.Name)
	ss, err := api.DeleteSharedStorageServer(test_ss.Id, test_server.Id)

	if err != nil {
		t.Errorf("DeleteSharedStorageServer failed. Error: " + err.Error())
		return
	}

	api.WaitForState(ss, "ACTIVE", 10, 30)
	ss, err = api.GetSharedStorage(ss.Id)

	if err != nil {
		t.Errorf("Deleting server attached to the shared storage failed.")
	} else {
		if len(ss.Servers) > 0 {
			t.Errorf("Server not deleted from the shared storage.")
		}
	}
}

func TestUpdateSharedStorage(t *testing.T) {
	set_ss.Do(set_shared_storage)

	fmt.Printf("Updating shared storage '%s'...\n", test_ss.Name)
	new_name := test_ss.Name + "_updated"
	new_desc := test_ss.Description + "_updated"
	new_size := 100
	ssu := SharedStorageRequest{
		Name:        new_name,
		Description: new_desc,
		Size:        &new_size,
	}
	ss, err := api.UpdateSharedStorage(test_ss.Id, &ssu)

	if err != nil {
		t.Errorf("UpdateSharedStorage failed. Error: " + err.Error())
	} else {
		api.WaitForState(ss, "ACTIVE", 10, 30)
	}
	ss, _ = api.GetSharedStorage(ss.Id)
	if ss.Name != new_name {
		t.Errorf("Failed to update shared storage name.")
	}
	if ss.Description != new_desc {
		t.Errorf("Failed to update shared storage description.")
	}
	if ss.Size != new_size {
		t.Errorf("Failed to update shared storage size.")
	}
}

func TestDeleteSharedStorage(t *testing.T) {
	set_ss.Do(set_shared_storage)

	fmt.Printf("Deleting shared storage '%s'...\n", test_ss.Name)
	ss, err := api.DeleteSharedStorage(test_ss.Id)

	if err != nil {
		t.Errorf("DeleteSharedStorage failed. Error: " + err.Error())
		return
	} else {
		api.WaitUntilDeleted(ss)
	}

	ss, err = api.GetSharedStorage(ss.Id)

	if ss != nil {
		t.Errorf("Unable to delete the shared storage.")
	} else {
		test_ss = nil
	}
}
