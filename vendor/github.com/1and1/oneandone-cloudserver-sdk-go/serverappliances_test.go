package oneandone

import (
	"fmt"
	"strings"
	"testing"
)

// /server_appliances tests

func TestListServerAppliances(t *testing.T) {
	fmt.Println("Listing all server appliances...")

	res, err := api.ListServerAppliances()
	if err != nil {
		t.Errorf("ListServerAppliances failed. Error: " + err.Error())
	}
	if len(res) == 0 {
		t.Errorf("No server appliance found.")
	}

	res, err = api.ListServerAppliances(1, 7, "name", "", "id,name,os_family")

	if err != nil {
		t.Errorf("ListServerAppliances with parameter options failed. Error: " + err.Error())
	}
	if len(res) == 0 {
		t.Errorf("No server appliance found.")
	}
	if len(res) != 7 {
		t.Errorf("Wrong number of objects per page.")
	}
	for index := 0; index < len(res); index += 1 {
		if res[index].Id == "" {
			t.Errorf("Filtering a list of server appliances failed.")
		}
		if res[index].Name == "" {
			t.Errorf("Filtering a list of server appliances failed.")
		}
		if res[index].Type != "" {
			t.Errorf("Filtering parameters failed.")
		}
		if index < len(res)-1 {
			if res[index].Name > res[index+1].Name {
				t.Errorf("Sorting a list of server appliances failed.")
			}
		}
	}
	// Test for error response
	res, err = api.ListServerAppliances(nil, nil, nil, nil, nil)
	if res != nil || err == nil {
		t.Errorf("ListServerAppliances failed to handle incorrect argument type.")
	}

	res, err = api.ListServerAppliances(0, 0, "", "linux", "")

	if err != nil {
		t.Errorf("ListServerAppliances with parameter options failed. Error: " + err.Error())
	}

	for _, sa := range res {
		if !strings.Contains(strings.ToLower(sa.OsFamily), "linux") {
			t.Errorf("Search parameter failed.")
		}
	}
}

func TestGetServerAppliance(t *testing.T) {
	saps, _ := api.ListServerAppliances(1, 1, "", "", "")
	fmt.Printf("Getting server appliance '%s'...\n", saps[0].Name)
	sa, err := api.GetServerAppliance(saps[0].Id)

	if sa == nil || err != nil {
		t.Errorf("GetServerAppliance failed. Error: " + err.Error())
	}
	if sa.Id != saps[0].Id {
		t.Errorf("Wrong ID of the server appliance.")
	}
}
