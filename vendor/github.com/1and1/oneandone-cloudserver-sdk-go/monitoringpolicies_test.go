package oneandone

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

var (
	set_mp       sync.Once
	test_mp_name string
	test_mp_desc string
	test_mp_mail string
	test_mp      *MonitoringPolicy
	mp_request   MonitoringPolicy
)

const (
	port_responding     = "RESPONDING"
	port_not_responding = "NOT_RESPONDING"
	process_running     = "RUNNING"
	process_not_running = "NOT_RUNNING"
	mp_email            = "test-go-sdk@oneandone.com"
)

// Helper functions

func create_monitoring_policy() *MonitoringPolicy {
	rand.Seed(time.Now().UnixNano())
	rint := rand.Intn(999)
	test_mp_name = fmt.Sprintf("MonitoringPolicy_%d", rint)
	test_mp_desc = fmt.Sprintf("MonitoringPolicy_%d description", rint)
	test_mp_mail = "test-go-sdk@oneandone.com"
	mp_request = MonitoringPolicy{
		Name:        test_mp_name,
		Description: test_mp_desc,
		Email:       test_mp_mail,
		Agent:       true,
		Thresholds: &MonitoringThreshold{
			Cpu: &MonitoringLevel{
				Warning: &MonitoringValue{
					Value: 90,
					Alert: false,
				},
				Critical: &MonitoringValue{
					Value: 95,
					Alert: false,
				},
			},
			Ram: &MonitoringLevel{
				Warning: &MonitoringValue{
					Value: 90,
					Alert: false,
				},
				Critical: &MonitoringValue{
					Value: 95,
					Alert: false,
				},
			},
			Disk: &MonitoringLevel{
				Warning: &MonitoringValue{
					Value: 80,
					Alert: false,
				},
				Critical: &MonitoringValue{
					Value: 90,
					Alert: false,
				},
			},
			Transfer: &MonitoringLevel{
				Warning: &MonitoringValue{
					Value: 1000,
					Alert: false,
				},
				Critical: &MonitoringValue{
					Value: 2000,
					Alert: false,
				},
			},
			InternalPing: &MonitoringLevel{
				Warning: &MonitoringValue{
					Value: 50,
					Alert: false,
				},
				Critical: &MonitoringValue{
					Value: 100,
					Alert: true,
				},
			},
		},
		Ports: []MonitoringPort{
			{
				Protocol:          "TCP",
				Port:              443,
				AlertIf:           port_not_responding,
				EmailNotification: true,
			},
		},
		Processes: []MonitoringProcess{
			{
				Process:           "httpdeamon",
				AlertIf:           process_not_running,
				EmailNotification: false,
			},
		},
	}
	fmt.Printf("Creating new monitoring policy '%s'...\n", test_mp_name)
	mp_id, mp, err := api.CreateMonitoringPolicy(&mp_request)
	if err != nil {
		fmt.Printf("Unable to create a monitoring policy. Error: %s", err.Error())
		return nil
	}
	if mp_id == "" || mp.Id == "" {
		fmt.Printf("Unable to create monitoring policy '%s'.", test_mp_name)
		return nil
	}
	api.WaitForState(mp, "ACTIVE", 30, 30)
	return mp
}

func set_monitoring_policy() {
	test_mp = create_monitoring_policy()
}

// /monitoring_policies tests

func TestCreateMonitoringPolicy(t *testing.T) {
	set_mp.Do(set_monitoring_policy)

	if test_mp == nil {
		t.Errorf("CreateMonitoringPolicy failed.")
	}
	if test_mp.Id == "" {
		t.Errorf("Missing monitoring policy ID.")
	}
	if test_mp.Name != test_mp_name {
		t.Errorf("Wrong name of the monitoring policy.")
	}
	if test_mp.Description != test_mp_desc {
		t.Errorf("Wrong monitoring policy description.")
	}
	if !test_mp.Agent {
		t.Errorf("Missing monitoring policy agent.")
	}
	if test_mp.CloudPanelId == "" {
		t.Errorf("Missing cloud panel ID in monitoring policy '%s' data.", test_mp.Name)
	}
	if test_mp.Email != test_mp_mail {
		t.Errorf("Wrong email of monitoring policy '%s'.", test_mp.Name)
	}
	if test_mp.Thresholds == nil {
		t.Errorf("Not thresholds found for monitoring policy '%s'.", test_mp.Name)
	} else {
		if test_mp.Thresholds.Cpu.Critical.Alert {
			t.Errorf("Wrong alerting state for CPU critical alert of monitoring policy '%s'.", test_mp.Name)
		}
		if test_mp.Thresholds.Cpu.Critical.Value != 95 {
			t.Errorf("Wrong value for CPU critical alert of monitoring policy '%s'.", test_mp.Name)
		}
		if test_mp.Thresholds.Cpu.Warning.Alert {
			t.Errorf("Wrong alerting state for CPU warning alert of monitoring policy '%s'.", test_mp.Name)
		}
		if test_mp.Thresholds.Cpu.Warning.Value != 90 {
			t.Errorf("Wrong value for CPU warning alert of monitoring policy '%s'.", test_mp.Name)
		}
		if test_mp.Thresholds.Ram.Critical.Alert {
			t.Errorf("Wrong alerting state for RAM critical alert of monitoring policy '%s'.", test_mp.Name)
		}
		if test_mp.Thresholds.Ram.Critical.Value != 95 {
			t.Errorf("Wrong value for RAM critical alert of monitoring policy '%s'.", test_mp.Name)
		}
		if test_mp.Thresholds.Ram.Warning.Alert {
			t.Errorf("Wrong alerting state for RAM warning alert of monitoring policy '%s'.", test_mp.Name)
		}
		if test_mp.Thresholds.Ram.Warning.Value != 90 {
			t.Errorf("Wrong value for RAM warning alert of monitoring policy '%s'.", test_mp.Name)
		}
		if test_mp.Thresholds.Disk.Critical.Alert {
			t.Errorf("Wrong alerting state for disk critical alert of monitoring policy '%s'.", test_mp.Name)
		}
		if test_mp.Thresholds.Disk.Critical.Value != 90 {
			t.Errorf("Wrong value for disk critical alert of monitoring policy '%s'.", test_mp.Name)
		}
		if test_mp.Thresholds.Disk.Warning.Alert {
			t.Errorf("Wrong alerting state for disk warning alert of monitoring policy '%s'.", test_mp.Name)
		}
		if test_mp.Thresholds.Disk.Warning.Value != 80 {
			t.Errorf("Wrong value for disk warning alert of monitoring policy '%s'.", test_mp.Name)
		}
		if test_mp.Thresholds.Transfer.Critical.Alert {
			t.Errorf("Wrong alerting state for transfer critical alert of monitoring policy '%s'.", test_mp.Name)
		}
		if test_mp.Thresholds.Transfer.Critical.Value != 2000 {
			t.Errorf("Wrong value for transfer critical alert of monitoring policy '%s'.", test_mp.Name)
		}
		if test_mp.Thresholds.Transfer.Warning.Alert {
			t.Errorf("Wrong alerting state for transfer warning alert of monitoring policy '%s'.", test_mp.Name)
		}
		if test_mp.Thresholds.Transfer.Warning.Value != 1000 {
			t.Errorf("Wrong value for transfer warning alert of monitoring policy '%s'.", test_mp.Name)
		}
		if !test_mp.Thresholds.InternalPing.Critical.Alert {
			t.Errorf("Wrong alerting state for ping critical alert of monitoring policy '%s'.", test_mp.Name)
		}
		if test_mp.Thresholds.InternalPing.Critical.Value != 100 {
			t.Errorf("Wrong value for ping critical alert of monitoring policy '%s'.", test_mp.Name)
		}
		if test_mp.Thresholds.InternalPing.Warning.Alert {
			t.Errorf("Wrong alerting state for ping warning alert of monitoring policy '%s'.", test_mp.Name)
		}
		if test_mp.Thresholds.InternalPing.Warning.Value != 50 {
			t.Errorf("Wrong value for ping warning alert of monitoring policy '%s'.", test_mp.Name)
		}
	}
	if len(test_mp.Ports) != 1 {
		t.Errorf("Wrong number of monitoring policy '%s' ports.", test_mp.Name)
	} else {
		if test_mp.Ports[0].Protocol != "TCP" {
			t.Errorf("Wrong monitoring protocol of policy '%s'.", test_mp.Name)
		}
		if test_mp.Ports[0].AlertIf != port_not_responding {
			t.Errorf("Wrong alerting state of monitoring policy '%s' port.", test_mp.Name)
		}
		if test_mp.Ports[0].Port != 443 {
			t.Errorf("Wrong monitoring policy '%s' port number.", test_mp.Name)
		}
		if !test_mp.Ports[0].EmailNotification {
			t.Errorf("Wrong email notification state in monitoring policy '%s' port.", test_mp.Name)
		}
	}
	if len(test_mp.Processes) != 1 {
		t.Errorf("Wrong number of monitoring policy '%s' processes.", test_mp.Name)
	} else {
		if test_mp.Processes[0].Process != "httpdeamon" {
			t.Errorf("Wrong monitoring policy '%s' process.", test_mp.Name)
		}
		if test_mp.Processes[0].AlertIf != process_not_running {
			t.Errorf("Wrong alerting state of monitoring policy '%s' process.", test_mp.Name)
		}
		if test_mp.Processes[0].EmailNotification {
			t.Errorf("Wrong email notification state in monitoring policy '%s' process.", test_mp.Name)
		}
	}
}

func TestGetMonitoringPolicy(t *testing.T) {
	set_mp.Do(set_monitoring_policy)

	fmt.Printf("Getting monitoring policy '%s'...\n", test_mp.Name)
	mp, err := api.GetMonitoringPolicy(test_mp.Id)

	if err != nil {
		t.Errorf("GetMonitoringPolicy failed. Error: " + err.Error())
	}
	if mp.Id != test_mp.Id {
		t.Errorf("Wrong monitoring policy ID.")
	}
}

func TestListMonitoringPolicies(t *testing.T) {
	set_mp.Do(set_monitoring_policy)
	fmt.Println("Listing all monitoring policies...")

	res, err := api.ListMonitoringPolicies()
	if err != nil {
		t.Errorf("ListMonitoringPolicies failed. Error: " + err.Error())
	}
	if len(res) == 0 {
		t.Errorf("No monitoring policy found.")
	}

	res, err = api.ListMonitoringPolicies(1, 1, "", "", "id,name,agent")

	if err != nil {
		t.Errorf("ListMonitoringPolicies with parameter options failed. Error: " + err.Error())
	}
	if len(res) == 0 {
		t.Errorf("No monitoring policy found.")
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

	res, err = api.ListMonitoringPolicies(0, 0, "", test_mp.Name, "")
	if err != nil {
		t.Errorf("ListMonitoringPolicies with parameter options failed. Error: " + err.Error())
	}
	if len(res) != 1 {
		t.Errorf("Search parameter failed.")
	}
	if res[0].Name != test_mp.Name {
		t.Errorf("Search parameter failed.")
	}
}

func TestAttachMonitoringPolicyServers(t *testing.T) {
	set_mp.Do(set_monitoring_policy)

	fmt.Printf("Attaching servers to monitoring policy '%s'...\n", test_mp.Name)
	sync_server.Do(func() { deploy_test_server(false) })

	servers := []string{test_server.Id}
	mp, err := api.AttachMonitoringPolicyServers(test_mp.Id, servers)

	if err != nil {
		t.Errorf("AttachMonitoringPolicyServers failed. Error: " + err.Error())
	}

	api.WaitForState(mp, "ACTIVE", 30, 60)
	mp, _ = api.GetMonitoringPolicy(mp.Id)

	if len(mp.Servers) != 1 {
		t.Errorf("Found no server attached to the monitoring policy.")
	}
	if mp.Servers[0].Id != test_server.Id {
		t.Errorf("Wrong server IP attached to the monitoring policy.")
	}
	test_mp = mp
}

func TestGetMonitoringPolicyServer(t *testing.T) {
	set_mp.Do(set_monitoring_policy)

	fmt.Printf("Getting server identity attached to monitoring policy '%s'...\n", test_mp.Name)
	mp_ser, err := api.GetMonitoringPolicyServer(test_mp.Id, test_mp.Servers[0].Id)

	if err != nil {
		t.Errorf("GetMonitoringPolicyServer failed. Error: " + err.Error())
	}
	if mp_ser.Id != test_mp.Servers[0].Id {
		t.Errorf("Wrong ID of the server attached to monitoring policy '%s'.", test_mp.Name)
	}
	if mp_ser.Name != test_mp.Servers[0].Name {
		t.Errorf("Wrong server name attached to monitoring policy '%s'.", test_mp.Name)
	}
}

func TestListMonitoringPolicyServers(t *testing.T) {
	set_mp.Do(set_monitoring_policy)
	sync_server.Do(func() { deploy_test_server(false) })

	fmt.Printf("Listing servers attached to monitoring policy '%s'...\n", test_mp.Name)
	mp_srvs, err := api.ListMonitoringPolicyServers(test_mp.Id)

	if err != nil {
		t.Errorf("ListMonitoringPolicyServers failed. Error: " + err.Error())
	}
	if len(mp_srvs) != 1 {
		t.Errorf("Wrong number of servers attached to monitoring policy '%s'.", test_mp.Name)
	}
	if mp_srvs[0].Id != test_server.Id {
		t.Errorf("Wrong server attached to monitoring policy '%s'.", test_mp.Name)
	}
}

func TestRemoveMonitoringPolicyServer(t *testing.T) {
	set_mp.Do(set_monitoring_policy)
	sync_server.Do(func() { deploy_test_server(false) })

	fmt.Printf("Removing server attached to monitoring policy '%s'...\n", test_mp.Name)
	mp, err := api.RemoveMonitoringPolicyServer(test_mp.Id, test_server.Id)

	if err != nil {
		t.Errorf("RemoveMonitoringPolicyServer failed. Error: " + err.Error())
	}

	api.WaitForState(mp, "ACTIVE", 30, 60)
	mp, err = api.GetMonitoringPolicy(mp.Id)

	if err != nil {
		t.Errorf("Removing server from the monitoring policy failed.")
	}
	if len(mp.Servers) > 0 {
		t.Errorf("Server not removed from the monitoring policy.")
	}
}

func TestGetMonitoringPolicyPort(t *testing.T) {
	set_mp.Do(set_monitoring_policy)

	fmt.Printf("Getting monitoring policy '%s' port...\n", test_mp.Name)
	mp_port, err := api.GetMonitoringPolicyPort(test_mp.Id, test_mp.Ports[0].Id)

	if err != nil {
		t.Errorf("GetMonitoringPolicyPort failed. Error: " + err.Error())
	}
	if mp_port.Id != test_mp.Ports[0].Id {
		t.Errorf("Wrong port ID.")
	}
	if mp_port.Port != test_mp.Ports[0].Port {
		t.Errorf("Wrong port number in the monitoring policy.")
	}
	if mp_port.AlertIf != test_mp.Ports[0].AlertIf {
		t.Errorf("Wrong alert_if field in the monitoring policy port.")
	}
	if mp_port.Protocol != test_mp.Ports[0].Protocol {
		t.Errorf("Wrong port protocol in the monitoring policy.")
	}
	if mp_port.EmailNotification != test_mp.Ports[0].EmailNotification {
		t.Errorf("Wrong email notification state in the monitoring policy port.")
	}
}

func TestModifyMonitoringPolicyPort(t *testing.T) {
	set_mp.Do(set_monitoring_policy)

	fmt.Printf("Modifying monitoring policy '%s' port...\n", test_mp.Name)
	mp_port, err := api.GetMonitoringPolicyPort(test_mp.Id, test_mp.Ports[0].Id)

	if err != nil {
		t.Errorf("GetMonitoringPolicyPort failed. Error: " + err.Error())
	}

	mp_port.AlertIf = port_responding
	mp_port.EmailNotification = false

	test_mp, err = api.ModifyMonitoringPolicyPort(test_mp.Id, mp_port.Id, mp_port)
	if err != nil {
		t.Errorf("ModifyMonitoringPolicyPort failed. Error: " + err.Error())
	}
	if test_mp.Ports[0].AlertIf != port_responding {
		t.Errorf("Unable to modify alert_if field in the monitoring policy port.")
	}
	if test_mp.Ports[0].EmailNotification {
		t.Errorf("Unable to modify email notification state in the monitoring policy port.")
	}
}

func TestAddMonitoringPolicyPorts(t *testing.T) {
	set_mp.Do(set_monitoring_policy)

	fmt.Printf("Adding ports to monitoring policy '%s'...\n", test_mp.Name)
	ports := []MonitoringPort{
		{
			Protocol:          "TCP",
			Port:              43215,
			AlertIf:           port_not_responding,
			EmailNotification: false,
		},
		{
			Protocol:          "UDP",
			Port:              161,
			AlertIf:           port_responding,
			EmailNotification: true,
		},
	}
	mp, err := api.AddMonitoringPolicyPorts(test_mp.Id, ports)

	if err != nil {
		t.Errorf("AddMonitoringPolicyPorts failed. Error: " + err.Error())
	} else {
		api.WaitForState(mp, "ACTIVE", 30, 60)
	}
	mp, _ = api.GetMonitoringPolicy(mp.Id)
	if len(mp.Ports) != 3 {
		t.Errorf("Unable to add ports to monitoring policy '%s'.\n", test_mp.Name)
	}
}

func TestListMonitoringPolicyPorts(t *testing.T) {
	set_mp.Do(set_monitoring_policy)

	fmt.Printf("Listing monitoring policy '%s' ports...\n", test_mp.Name)
	mp_ports, err := api.ListMonitoringPolicyPorts(test_mp.Id)

	if err != nil {
		t.Errorf("ListMonitoringPolicyPorts failed. Error: " + err.Error())
	}
	if len(mp_ports) != 3 {
		t.Errorf("Wrong number of ports found in monitoring policy '%s'.", test_mp.Name)
	}
}

func TestDeleteMonitoringPort(t *testing.T) {
	set_mp.Do(set_monitoring_policy)

	mp_ports, _ := api.ListMonitoringPolicyPorts(test_mp.Id)
	fmt.Printf("Deleting port '%s' from monitoring policy '%s'...\n", mp_ports[0].Id, test_mp.Name)
	mp, err := api.DeleteMonitoringPolicyPort(test_mp.Id, mp_ports[0].Id)

	if err != nil {
		t.Errorf("DeleteMonitoringPolicyPort failed. Error: " + err.Error())
	}

	api.WaitForState(mp, "ACTIVE", 30, 60)
	mp, err = api.GetMonitoringPolicy(mp.Id)

	if err != nil {
		t.Errorf("Deleting port from the monitoring policy failed.")
	}
	if len(mp.Ports) != 2 {
		t.Errorf("Port not deleted from the monitoring policy.")
	}
	for _, port := range mp.Ports {
		if port.Id == mp_ports[0].Id {
			t.Errorf("Port not deleted from the monitoring policy.")
		}
	}
}

func TestGetMonitoringPolicyProcess(t *testing.T) {
	set_mp.Do(set_monitoring_policy)

	fmt.Printf("Getting monitoring policy '%s' process...\n", test_mp.Name)
	mp_process, err := api.GetMonitoringPolicyProcess(test_mp.Id, test_mp.Processes[0].Id)

	if err != nil {
		t.Errorf("GetMonitoringPolicyProcess failed. Error: " + err.Error())
	}
	if mp_process.Id != test_mp.Processes[0].Id {
		t.Errorf("Wrong process ID.")
	}
	if mp_process.Process != test_mp.Processes[0].Process {
		t.Errorf("Wrong process name in the monitoring policy.")
	}
	if mp_process.AlertIf != test_mp.Processes[0].AlertIf {
		t.Errorf("Wrong alert_if field in the monitoring policy process.")
	}
	if mp_process.EmailNotification != test_mp.Processes[0].EmailNotification {
		t.Errorf("Wrong email notification state in the monitoring policy process.")
	}
}

func TestModifyMonitoringPolicyProcess(t *testing.T) {
	set_mp.Do(set_monitoring_policy)

	fmt.Printf("Modifying monitoring policy '%s' process...\n", test_mp.Name)
	mp_process, err := api.GetMonitoringPolicyProcess(test_mp.Id, test_mp.Processes[0].Id)

	if err != nil {
		t.Errorf("GetMonitoringPolicyProcess failed. Error: " + err.Error())
	}

	mp_process.AlertIf = process_running
	mp_process.EmailNotification = true

	test_mp, err = api.ModifyMonitoringPolicyProcess(test_mp.Id, mp_process.Id, mp_process)
	if err != nil {
		t.Errorf("ModifyMonitoringPolicyProcess failed. Error: " + err.Error())
	}
	if test_mp.Processes[0].AlertIf != process_running {
		t.Errorf("Unable to modify alert_if field in the monitoring policy process.")
	}
	if !test_mp.Processes[0].EmailNotification {
		t.Errorf("Unable to modify email notification state in the monitoring policy process.")
	}
}

func TestAddMonitoringPolicyProcesses(t *testing.T) {
	set_mp.Do(set_monitoring_policy)

	fmt.Printf("Adding processes to monitoring policy '%s'...\n", test_mp.Name)
	processes := []MonitoringProcess{
		{
			Process:           "iexplorer",
			AlertIf:           process_not_running,
			EmailNotification: false,
		},
		{
			Process:           "taskmgr",
			AlertIf:           process_running,
			EmailNotification: true,
		},
	}
	mp, err := api.AddMonitoringPolicyProcesses(test_mp.Id, processes)

	if err != nil {
		t.Errorf("AddMonitoringPolicyProcesses failed. Error: " + err.Error())
	} else {
		api.WaitForState(mp, "ACTIVE", 30, 60)
	}
	mp, _ = api.GetMonitoringPolicy(mp.Id)
	if len(mp.Processes) != 3 {
		t.Errorf("Unable to add processes to monitoring policy '%s'.\n", test_mp.Name)
	}
}

func TestListMonitoringPolicyProcesses(t *testing.T) {
	set_mp.Do(set_monitoring_policy)

	fmt.Printf("Listing monitoring policy '%s' processes...\n", test_mp.Name)
	mp_processes, err := api.ListMonitoringPolicyProcesses(test_mp.Id)

	if err != nil {
		t.Errorf("ListMonitoringPolicyProcesses failed. Error: " + err.Error())
	}
	if len(mp_processes) != 3 {
		t.Errorf("Wrong number of processes found in monitoring policy '%s'.", test_mp.Name)
	}
}

func TestDeleteMonitoringProcess(t *testing.T) {
	set_mp.Do(set_monitoring_policy)

	mp_processes, _ := api.ListMonitoringPolicyProcesses(test_mp.Id)
	fmt.Printf("Deleting process '%s' from monitoring policy '%s'...\n", mp_processes[0].Id, test_mp.Name)
	mp, err := api.DeleteMonitoringPolicyProcess(test_mp.Id, mp_processes[0].Id)

	if err != nil {
		t.Errorf("DeleteMonitoringPolicyProcess failed. Error: " + err.Error())
	}

	api.WaitForState(mp, "ACTIVE", 60, 60)
	mp, err = api.GetMonitoringPolicy(mp.Id)

	if err != nil {
		t.Errorf("Deleting process from the monitoring policy failed.")
	}
	if len(mp.Processes) != 2 {
		t.Errorf("Process not deleted from the monitoring policy.")
	}
	for _, process := range mp.Processes {
		if process.Id == mp_processes[0].Id {
			t.Errorf("Process not deleted from the monitoring policy.")
		}
	}
}

func TestUpdateMonitoringPolicy(t *testing.T) {
	set_mp.Do(set_monitoring_policy)

	fmt.Printf("Updating monitoring policy '%s'...\n", test_mp.Name)
	new_name := test_mp.Name + "_updated"
	new_desc := test_mp.Description + "_updated"
	new_mail := "test-go-sdk@oneandone.de"
	mp_request.Name = new_name
	mp_request.Description = new_desc
	mp_request.Email = new_mail
	mp_request.Thresholds.Cpu.Critical.Value = 99
	mp_request.Thresholds.Transfer.Critical.Alert = true

	mp, err := api.UpdateMonitoringPolicy(test_mp.Id, &mp_request)

	if err != nil {
		t.Errorf("UpdateMonitoringPolicy failed. Error: " + err.Error())
	} else {
		api.WaitForState(mp, "ACTIVE", 30, 30)
		mp, _ = api.GetMonitoringPolicy(mp.Id)
		if mp.Name != new_name {
			t.Errorf("Failed to update monitoring policy name.")
		}
		if mp.Description != new_desc {
			t.Errorf("Failed to update monitoring policy description.")
		}
		if mp.Email != new_mail {
			t.Errorf("Failed to update monitoring policy email.")
		}
		if mp.Thresholds.Cpu.Critical.Value != 99 {
			t.Errorf("Failed to update critical CPU threshold of the monitoring policy.")
		}
		if !mp.Thresholds.Transfer.Critical.Alert {
			t.Errorf("Failed to update alerting state for critical transfer threshold of the monitoring policy.")
		}
	}
}

func TestDeleteMonitoringPolicy(t *testing.T) {
	set_mp.Do(set_monitoring_policy)

	fmt.Printf("Deleting monitoring policy '%s'...\n", test_mp.Name)
	time.Sleep(time.Second)
	mp, err := api.DeleteMonitoringPolicy(test_mp.Id)

	if err != nil {
		t.Errorf("DeleteMonitoringPolicy failed. Error: " + err.Error())
		return
	}
	api.WaitUntilDeleted(mp)
	time.Sleep(time.Second)
	mp, err = api.GetMonitoringPolicy(mp.Id)

	if mp != nil {
		t.Errorf("Unable to delete the monitoring policy.")
	} else {
		test_mp = nil
	}
}
