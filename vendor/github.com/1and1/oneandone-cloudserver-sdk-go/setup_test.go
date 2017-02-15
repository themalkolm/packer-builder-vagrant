package oneandone

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"testing"
)

var (
	api         *API
	sync_server sync.Once
	test_server *Server
)

func Init() {
	fmt.Println("Initializing tests...")
	token, err := getEnvironmentVar("ONEANDONE_TOKEN")
	if err != nil {
		fmt.Printf("The 1&1 cloud server api token must be set in the environment variable 'ONEANDONE_TOKEN'")
		os.Exit(1)
	}

	SetToken(token)

	apiEndpoint, err := getEnvironmentVar("ONEANDONE_ENDPOINT")
	if err != nil {
		// Use default endpoint
		apiEndpoint = BaseUrl
	}

	api = New(token, apiEndpoint)
}

func getEnvironmentVar(name string) (string, error) {
	osVar := os.Getenv(name)
	if osVar == "" {
		return "", fmt.Errorf("The environment variable '%s' is not set.", name)
	}
	return osVar, nil
}

func setEnvironmentVar(name string, value string) {
	err := os.Setenv(name, value)
	if err != nil {
		fmt.Printf("The environment variable '%s' is not set.", name)
	}
}

func printObject(in interface{}) {
	bytes, _ := json.MarshalIndent(in, "", "    ")
	fmt.Printf("%v\n", string(bytes))
}

func deploy_test_server(power_on bool) {
	_, test_server, _ = create_test_server(power_on)
	if power_on {
		api.WaitForState(test_server, "POWERED_ON", 10, 90)
	} else {
		api.WaitForState(test_server, "POWERED_OFF", 10, 90)
	}
}

func Cleanup() {
	if server != nil {
		api.DeleteServer(server.Id, false)
	}
	if test_server != nil {
		api.DeleteServer(test_server.Id, false)
	}
	if test_image != nil {
		api.DeleteImage(test_image.Id)
	}
	if test_fp != nil {
		api.DeleteFirewallPolicy(test_fp.Id)
	}
	if test_lb != nil {
		api.DeleteLoadBalancer(test_lb.Id)
	}
	if test_mp != nil {
		api.DeleteMonitoringPolicy(test_mp.Id)
	}
	if test_pn != nil {
		api.DeletePrivateNetwork(test_pn.Id)
	}
	if test_ip != nil {
		api.DeletePublicIp(test_ip.Id)
	}
	if test_ss != nil {
		api.DeleteSharedStorage(test_ss.Id)
	}
	if ser_lb != nil {
		api.DeleteLoadBalancer(ser_lb.Id)
	}
	if image_serv != nil {
		api.DeleteServer(image_serv.Id, false)
	}
	if test_vpn != nil {
		api.DeleteVPN(test_vpn.Id)
	}
	if test_role != nil {
		api.DeleteRole(test_role.Id)
	}
}

func TestMain(m *testing.M) {
	Init()
	rc := m.Run()
	Cleanup()
	os.Exit(rc)
}
