package oneandone

import (
	"fmt"
	"testing"
	"time"
)

// /monitoring_center tests

func TestListMonitoringServersUsages(t *testing.T) {
	sync_server.Do(func() { deploy_test_server(false) })
	fmt.Println("Listing all server usages...")

	res, err := api.ListMonitoringServersUsages()
	if err != nil {
		t.Errorf("ListMonitoringServersUsages failed. Error: " + err.Error())
	}
	if len(res) == 0 {
		t.Errorf("No server usage found.")
	}

	res, err = api.ListMonitoringServersUsages(1, 2, "name", "", "id,name")

	if err != nil {
		t.Errorf("ListMonitoringServersUsages with parameter options failed. Error: " + err.Error())
		return
	}
	if len(res) == 0 {
		t.Errorf("No server usage found.")
	}
	if len(res) > 2 {
		t.Errorf("Wrong number of objects per page.")
	}
	if len(res) == 2 {
		if res[0].Name > res[1].Name {
			t.Errorf("Sorting a list of server usages failed.")
		}
	}
	for index := 0; index < len(res); index += 1 {
		if res[index].Agent != nil {
			t.Errorf("Filtering a list of server usages failed.")
		}
		if res[index].Alerts != nil {
			t.Errorf("Filtering a list of server usages failed.")
		}
		if res[index].Status != nil {
			t.Errorf("Filtering a list of server usages failed.")
		}
	}

	res, err = api.ListMonitoringServersUsages(0, 0, "", test_server.Name, "")

	if err != nil {
		t.Errorf("ListMonitoringServersUsages with parameter options failed. Error: " + err.Error())
	}
	if len(res) != 1 {
		t.Errorf("Wrong number of objects found.")
	}
}

func TestGetMonitoringServerUsage(t *testing.T) {
	sync_server.Do(func() { deploy_test_server(false) })
	fmt.Printf("Getting server usage '%s'...\n", test_server.Name)
	msu, err := api.GetMonitoringServerUsage(test_server.Id, "LAST_HOUR")

	if err != nil {
		t.Errorf("GetMonitoringServerUsage failed. Error: " + err.Error())
		return
	}
	if msu.Id != test_server.Id {
		t.Errorf("Wrong server ID.")
	}
	if msu.Name != test_server.Name {
		t.Errorf("Wrong server name.")
	}

	if msu.TransferStatus != nil && len(msu.TransferStatus.Data) > 0 {
		date, _ := time.Parse(time.RFC3339, msu.TransferStatus.Data[0].Date)
		msu, err = api.GetMonitoringServerUsage(test_server.Id, "CUSTOM", date, date)

		if err != nil {
			t.Errorf("GetMonitoringServerUsage failed. Error: " + err.Error())
			return
		}
		if len(msu.TransferStatus.Data) != 1 {
			t.Errorf("Getting server usage in custom date range failed.")
		}
	}
	// Test error response
	msu, err = api.GetMonitoringServerUsage(test_server.Id, "")
	if msu != nil || err == nil {
		t.Errorf("GetMonitoringServerUsage failed to handle empty date period.")
	}
	msu, err = api.GetMonitoringServerUsage(test_server.Id, "CUSTOM", time.Now(), time.Now().Add(-time.Hour*25))
	if msu != nil || err == nil {
		t.Errorf("GetMonitoringServerUsage failed to handle 'start_date > end_date' case.")
	}
	msu, err = api.GetMonitoringServerUsage(test_server.Id, "CUSTOM", time.Now())
	if msu != nil || err == nil {
		t.Errorf("GetMonitoringServerUsage failed to handle invalid number of date period parameters.")
	}
}
