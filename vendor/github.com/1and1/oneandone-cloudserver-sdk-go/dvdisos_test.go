package oneandone

import (
	"fmt"
	"strings"
	"testing"
)

// /dvd_isos tests

func TestListDvdIsos(t *testing.T) {
	fmt.Println("Listing all dvd isos...")

	res, err := api.ListDvdIsos()
	if err != nil {
		t.Errorf("ListDvdIsos failed. Error: " + err.Error())
	}
	if len(res) == 0 {
		t.Errorf("No dvd found.")
	}

	res, err = api.ListDvdIsos(1, 6, "name", "", "id,name")

	if err != nil {
		t.Errorf("ListDvdIsos with parameter options failed. Error: " + err.Error())
	}
	if len(res) == 0 {
		t.Errorf("No dvd found.")
	}
	if len(res) != 6 {
		t.Errorf("Wrong number of objects per page.")
	}
	for index := 0; index < len(res); index += 1 {
		if res[index].Id == "" {
			t.Errorf("Filtering a list of dvd isos failed.")
		}
		if res[index].Name == "" {
			t.Errorf("Filtering a list of dvd isos failed.")
		}
		if res[index].Type != "" {
			t.Errorf("Filtering parameters failed.")
		}
		if index < len(res)-1 {
			if res[index].Name > res[index+1].Name {
				t.Errorf("Sorting a list of dvd isos failed.")
			}
		}
	}
	// Test for error response
	res, err = api.ListDvdIsos(0, 0, "name", "dvd", "id", "type")
	if res != nil || err == nil {
		t.Errorf("ListDvdIsos failed to handle incorrect number of passed arguments.")
	}

	res, err = api.ListDvdIsos(0, 0, "", "freebsd", "")

	if err != nil {
		t.Errorf("ListDvdIsos with parameter options failed. Error: " + err.Error())
	}

	for _, dvd := range res {
		if !strings.Contains(strings.ToLower(dvd.Name), "freebsd") {
			t.Errorf("Search parameter failed.")
		}
	}
}

func TestGetDvdIso(t *testing.T) {
	dvds, _ := api.ListDvdIsos(1, 1, "", "", "")
	fmt.Printf("Getting dvd iso '%s'...\n", dvds[0].Name)
	dvd, err := api.GetDvdIso(dvds[0].Id)

	if err != nil {
		t.Errorf("GetDvdIso failed. Error: " + err.Error())
	}
	if dvd.Id != dvds[0].Id {
		t.Errorf("Wrong ID of the dvd iso.")
	}
}
