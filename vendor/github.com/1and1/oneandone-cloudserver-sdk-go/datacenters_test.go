package oneandone

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// /datacenters tests

func TestListDatacenters(t *testing.T) {
	fmt.Println("Listing datacenters...")

	res, err := api.ListDatacenters()
	if err != nil {
		t.Errorf("ListDatacenters failed. Error: " + err.Error())
	}
	if len(res) == 0 {
		t.Errorf("No datacenter found.")
	}

	res, err = api.ListDatacenters(1, 2, "location", "", "id,location")

	if err != nil {
		t.Errorf("ListDatacenters with parameter options failed. Error: " + err.Error())
	}
	if len(res) == 0 {
		t.Errorf("No datacenter found.")
	}
	if len(res) != 2 {
		t.Errorf("Wrong number of objects per page.")
	}
	for index := 0; index < len(res); index += 1 {
		if res[index].Id == "" {
			t.Errorf("Filtering a list of datacenters failed.")
		}
		if res[index].Location == "" {
			t.Errorf("Filtering a list of datacenters failed.")
		}
		if index < len(res)-1 {
			if res[index].Location > res[index+1].Location {
				t.Errorf("Sorting a list of datacenters failed.")
			}
		}
	}
	// Test for error response
	res, err = api.ListDatacenters(0, 0, "location", "Spain", "id", "country_code")
	if res != nil || err == nil {
		t.Errorf("ListDatacenters failed to handle incorrect number of passed arguments.")
	}

	res, err = api.ListDatacenters(0, 0, "", "Germany", "")

	if err != nil {
		t.Errorf("ListDatacenters with parameter options failed. Error: " + err.Error())
	}

	for _, dc := range res {
		if !strings.Contains(dc.Location, "Germany") {
			t.Errorf("Search parameter failed.")
		}
	}
}

func TestGetDatacenter(t *testing.T) {
	dcs, err := api.ListDatacenters()

	if len(dcs) == 0 {
		t.Errorf("No datacenter found. " + err.Error())
		return
	}

	for i, _ := range dcs {
		time.Sleep(time.Second)
		fmt.Printf("Getting datacenter '%s'...\n", dcs[i].CountryCode)
		dc, err := api.GetDatacenter(dcs[i].Id)

		if err != nil {
			t.Errorf("GetDatacenter failed. Error: " + err.Error())
			return
		}
		if dc.Id != dcs[i].Id {
			t.Errorf("Wrong datacenter ID.")
		}
		if dc.CountryCode != dcs[i].CountryCode {
			t.Errorf("Wrong country code of the datacenter.")
		}
		if dc.Location != dcs[i].Location {
			t.Errorf("Wrong datacenter location.")
		}
	}
}
