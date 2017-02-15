package oneandone

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// /logs tests

func TestListLogs(t *testing.T) {
	fmt.Println("Listing all logs...")

	res, err := api.ListLogs("LAST_24H", nil, nil)
	if err != nil {
		t.Errorf("ListLogs failed. Error: " + err.Error())
	}
	if len(res) == 0 {
		t.Errorf("No log found.")
	}

	res, err = api.ListLogs("LAST_24H", nil, nil, 1, 7, "action", "", "action,start_date,end_date")

	if err != nil {
		t.Errorf("ListLogs with parameter options failed. Error: " + err.Error())
		return
	}
	if len(res) == 0 {
		t.Errorf("No log found.")
	}
	if len(res) != 7 {
		t.Errorf("Wrong number of objects per page.")
	}
	for index := 0; index < len(res); index += 1 {
		if res[index].Status != nil || res[index].Resource != nil || res[index].User != nil {
			t.Errorf("Filtering a list of logs failed.")
		}
		if index < len(res)-1 {
			if res[index].Action > res[index+1].Action {
				t.Errorf("Sorting a list of logs failed.")
			}
		}
	}

	sd, _ := time.Parse(time.RFC3339, res[len(res)/2].StartDate)
	ed := sd.Add(time.Hour)
	res, err = api.ListLogs("CUSTOM", &sd, &ed, 0, 0, "start_date", "", "action,start_date,end_date")

	if err != nil {
		t.Errorf("Getting logs in custom date range failed. Error: " + err.Error())
		return
	}
	if len(res) > 0 {
		sd1, _ := time.Parse(time.RFC3339, res[0].StartDate)
		ed1, _ := time.Parse(time.RFC3339, res[len(res)-1].EndDate)
		if sd1.Before(sd) || ed1.After(ed) {
			t.Errorf("Getting logs in custom date range failed.")
		}
	}

	res, err = api.ListLogs("LAST_7D", nil, nil, 0, 0, "", "CREATE", "action")

	if err != nil {
		t.Errorf("ListLogs with parameter options failed. Error: " + err.Error())
		return
	}
	for _, log := range res {
		if !strings.Contains(strings.ToUpper(log.Action), "CREATE") {
			t.Errorf("Search parameter failed.")
		}
	}
	// Test for error response
	res, err = api.ListLogs("LAST_24H", nil, nil, 2, 5, 5)
	if res != nil || err == nil {
		t.Errorf("ListLogs failed to handle incorrect argument type.")
	}
}

func TestGetLog(t *testing.T) {
	logs, _ := api.ListLogs("LAST_24H", nil, nil, 0, 0, "", "", "id")
	fmt.Printf("Getting log '%s'...\n", logs[0].Id)
	log, err := api.GetLog(logs[0].Id)

	if err != nil {
		t.Errorf("GetLog failed. Error: " + err.Error())
		return
	}
	if log.Id != logs[0].Id {
		t.Errorf("Wrong server ID.")
	}
}
