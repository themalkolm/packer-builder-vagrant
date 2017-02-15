package oneandone

import (
	"fmt"
	"testing"
	"time"
)

// /usages tests

func TestListUsages(t *testing.T) {
	fmt.Println("Listing all usages...")

	usages, err := api.ListUsages("LAST_7D", nil, nil)
	if err != nil {
		t.Errorf("ListUsages failed. Error: " + err.Error())
		return
	}
	if len(usages.Servers) == 0 && len(usages.Images) == 0 && len(usages.LoadBalancers) == 0 &&
		len(usages.PublicIPs) == 0 && len(usages.SharedStorages) == 0 {
		t.Errorf("No usage found.")
	}

	usages, err = api.ListUsages("LAST_24H", nil, nil, 0, 0, "", "", "SERVERS.id,SERVERS.name")

	if err != nil {
		t.Errorf("ListUsages with parameter options failed. Error: " + err.Error())
		return
	}
	if len(usages.Servers) == 0 {
		t.Errorf("No usage found.")
	}
	if len(usages.Images) != 0 {
		t.Errorf("Filtering a list of usages failed.")
	}
	if len(usages.LoadBalancers) != 0 {
		t.Errorf("Filtering a list of usages failed.")
	}
	if len(usages.SharedStorages) != 0 {
		t.Errorf("Filtering a list of usages failed.")
	}
	if len(usages.PublicIPs) != 0 {
		t.Errorf("Filtering a list of usages failed.")
	}

	n := time.Now()
	ed := time.Date(n.Year(), n.Month(), n.Day(), n.Hour(), n.Minute(), n.Second(), 0, time.UTC)
	sd := ed.Add(-(time.Hour * 5))

	usages, err = api.ListUsages("CUSTOM", &sd, &ed)

	if err != nil {
		t.Errorf("Getting usages in custom date range failed. Error: " + err.Error())
		return
	}
	if len(usages.Servers) > 0 {
		sd1, _ := time.Parse(time.RFC3339, usages.Servers[0].Services[0].Details[0].StartDate)
		ed1, _ := time.Parse(time.RFC3339, usages.Servers[0].Services[0].Details[0].EndDate)
		if sd1.Unix() > sd.Unix() || ed1.Unix() > ed.Unix() {
			t.Errorf("Getting usages in custom date range failed.")
		}
	}
	// Tests for error response
	usages, err = api.ListUsages("LAST_24H", nil, nil, true)
	if usages != nil || err == nil {
		t.Errorf("ListUsages failed to handle incorrect argument type.")
	}
	usages, err = api.ListUsages("LAST_24H", &ed, &sd)
	if usages != nil || err == nil {
		t.Errorf("ListUsages failed to handle 'start_date > end_date' case.")
	}
}
