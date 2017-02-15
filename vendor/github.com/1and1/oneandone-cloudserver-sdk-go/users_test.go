package oneandone

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"
)

var (
	set_user  sync.Once
	getc_user sync.Once
	tuname    string
	tudesc    string
	tumail    string
	tu_id     string
	cur_user  *User
	test_user *User
)

// Helper functions

func create_user() *User {
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(999999)
	tuname = fmt.Sprintf("testuser_%d", r)
	tudesc = fmt.Sprintf("testuser_%d description", r)
	tumail = "testuser@1and1.com"
	req := UserRequest{
		Name:        tuname,
		Description: tudesc,
		Password:    fmt.Sprintf("Ss&)Hg3@&9!hJ&5%d", r),
		Email:       tumail,
	}
	fmt.Printf("Creating new user '%s'...\n", tuname)
	u_id, u, err := api.CreateUser(&req)
	if err != nil {
		fmt.Printf("Unable to create a user. Error: %s", err.Error())
		return nil
	}
	if u_id == "" || u.Id == "" {
		fmt.Printf("Unable to create user '%s'.", tuname)
		return nil
	}

	return u
}

func get_my_pub_ip() string {
	res, err := http.Get("http://echoip.com/")
	if err != nil {
		return ""
	}
	defer res.Body.Close()
	ip, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ""
	}
	return string(ip)
}

func set_test_users() {
	test_user = create_user()
}

func get_current_user() {
	res, _ := api.ListUsers()

	if res != nil {
		for _, u := range res {
			if Token == u.Api.Key {
				cur_user = &u
				break
			}
		}
	}
}

// /users tests

func TestCreateUser(t *testing.T) {
	t.Skip("TestCreateUser is skipped at the moment.")
	set_user.Do(set_test_users)

	if test_user == nil {
		t.Errorf("CreateUser failed.")
		return
	}
	if !strings.Contains(test_user.Name, tuname) {
		t.Errorf("Wrong user name.")
	}
	if test_user.Description != tudesc {
		t.Errorf("Wrong user description.")
	}
	if test_user.Email != tumail {
		t.Errorf("Wrong user email.")
	}
}

func TestListUsers(t *testing.T) {
	//	set_user.Do(set_test_users)
	fmt.Println("Listing all users...")

	res, err := api.ListUsers()
	if err != nil {
		t.Errorf("ListUsers failed. Error: " + err.Error())
	}
	if len(res) == 0 {
		t.Errorf("No user found.")
	}

	for _, u := range res {
		if Token == u.Api.Key {
			cur_user = &u
			break
		}
	}

	res, err = api.ListUsers(1, 3, "name", "", "id,name")

	if err != nil {
		t.Errorf("ListUsers with parameter options failed. Error: " + err.Error())
		return
	}
	if len(res) == 0 {
		t.Errorf("No user found.")
	}
	if len(res) != 3 {
		t.Errorf("Wrong number of objects per page.")
	}

	for i, _ := range res {
		if res[i].Id == "" || res[i].Name == "" || res[i].State != "" ||
			res[i].Api != nil || res[i].Role != nil {
			t.Errorf("Filtering parameters failed.")
		}
		if i < len(res)-1 {
			if res[i].Name > res[i+1].Name {
				t.Errorf("Sorting list of users failed.")
			}
		}
	}
	// Test for error response
	res, err = api.ListUsers("", nil, 10, "15", "")
	if res != nil || err == nil {
		t.Errorf("ListUsers failed to handle incorrect argument type.")
	}

	res, err = api.ListUsers(0, 0, "", cur_user.Name, "")

	if err != nil {
		t.Errorf("ListUsers with parameter options failed. Error: " + err.Error())
		return
	}
	if len(res) != 1 {
		t.Errorf("Search parameter failed.")
	}
	if res[0].Name != cur_user.Name {
		t.Errorf("Search parameter failed.")
	}
}

func TestGetUser(t *testing.T) {
	getc_user.Do(get_current_user)

	fmt.Printf("Getting user '%s'...\n", cur_user.Name)
	u, err := api.GetUser(cur_user.Id)

	if err != nil {
		t.Errorf("GetUser failed. Error: " + err.Error())
	} else {
		if u.Id != cur_user.Id {
			t.Errorf("Wrong user ID.")
		}
	}
}

func TestGetCurrentUserPermissions(t *testing.T) {
	getc_user.Do(get_current_user)

	fmt.Printf("Getting current user permissions ...\n")
	p, err := api.GetCurrentUserPermissions()

	if err != nil {
		t.Errorf("GetCurrentUserPermissions failed. Error: " + err.Error())
	} else {
		if p == nil || p.Backups == nil || p.Firewalls == nil || p.Images == nil || p.Invoice == nil ||
			p.IPs == nil || p.LoadBalancers == nil || p.Logs == nil || p.MonitorCenter == nil ||
			p.MonitorPolicies == nil || p.PrivateNetworks == nil || p.Roles == nil || p.Servers == nil ||
			p.SharedStorage == nil || p.Usages == nil || p.Users == nil || p.VPNs == nil {
			t.Errorf("Some permissions objects are missing.")
		}
	}
}

func TestGetUserApi(t *testing.T) {
	getc_user.Do(get_current_user)

	fmt.Printf("Getting API data of user '%s'...\n", cur_user.Name)
	ua, err := api.GetUserApi(cur_user.Id)

	if err != nil {
		t.Errorf("GetUserApi failed. Error: " + err.Error())
	} else {
		if ua.Key != cur_user.Api.Key {
			t.Errorf("Wrong user key.")
		}
		if !ua.Active {
			t.Errorf("Wrong user active state.")
		}
	}
}

func TestModifyUserApi(t *testing.T) {
	getc_user.Do(get_current_user)

	fmt.Printf("Modify API state of user '%s'...\n", cur_user.Name)
	// Just making sure that the request pass.
	// TODO: test with active=false once the REST functionality is completed.
	u, err := api.ModifyUserApi(cur_user.Id, true)

	if err != nil {
		t.Errorf("ModifyUserApi failed. Error: " + err.Error())
	} else {
		if u.Api.Key != cur_user.Api.Key {
			t.Errorf("Wrong user key.")
		}
		if !u.Api.Active {
			t.Errorf("Wrong user active state.")
		}
	}
}

func TestGetUserApiKey(t *testing.T) {
	getc_user.Do(get_current_user)

	fmt.Printf("Getting API data of user '%s'...\n", cur_user.Name)
	key, err := api.GetUserApiKey(cur_user.Id)

	if err != nil {
		t.Errorf("GetUserApiKey failed. Error: " + err.Error())
	} else {
		if key.Key != cur_user.Api.Key {
			t.Errorf("Wrong user key.")
		}
	}
}

func TestAddUserApiAlowedIps(t *testing.T) {
	getc_user.Do(get_current_user)

	fmt.Printf("Adding API allowed IPs to user '%s'...\n", cur_user.Name)
	my_ip := get_my_pub_ip()
	if my_ip == "" {
		fmt.Println("Not able to obtain its own public IP. Skipping the test.")
		return
	}
	ips := []string{my_ip, "192.168.7.77", "10.81.12.101"}
	u, err := api.AddUserApiAlowedIps(cur_user.Id, ips)

	if err != nil {
		t.Errorf("AddUserApiAlowedIps failed. Error: " + err.Error())
	} else {
		if len(u.Api.AllowedIps) != 3 {
			t.Errorf("Unable to add API allowed IPs to the user.")
		}
		for _, a := range u.Api.AllowedIps {
			if a != my_ip && a != "192.168.7.77" && a != "10.81.12.101" {
				t.Errorf("Wrong IP added to user's API allowed list.")
			}
		}
	}
}

func TestListUserApiAllowedIps(t *testing.T) {
	getc_user.Do(get_current_user)

	fmt.Printf("Listing API allowed IPs to user '%s'...\n", cur_user.Name)
	ips, err := api.ListUserApiAllowedIps(cur_user.Id)

	if err != nil {
		t.Errorf("ListUserApiAllowedIps failed. Error: " + err.Error())
	} else {
		if len(ips) != 3 {
			t.Errorf("Wrong number of API allowed IPs found.")
		}
	}
}

func TestRemoveUserApiAllowedIp(t *testing.T) {
	getc_user.Do(get_current_user)

	fmt.Printf("Removing API allowed IPs to user '%s'...\n", cur_user.Name)
	my_ip := get_my_pub_ip()
	if my_ip == "" {
		fmt.Println("Not able to obtain its own public IP. Skipping the test.")
		return
	}
	ips := []string{"192.168.7.77", "10.81.12.101", my_ip}

	for _, ip := range ips {
		_, err := api.RemoveUserApiAllowedIp(cur_user.Id, ip)

		if err != nil {
			t.Errorf("RemoveUserApiAllowedIp failed. Error: " + err.Error())
		}
	}
	u, _ := api.GetUser(cur_user.Id)
	if len(u.Api.AllowedIps) != 0 {
		t.Errorf("RemoveUserApiAllowedIp failed.")
	}
}

func TestRenewUserApiKey(t *testing.T) {
	t.Skip("TestRenewUserApiKey is skipped at the moment.")
	getc_user.Do(get_current_user)

	fmt.Printf("Renewing API key of user '%s'...\n", cur_user.Name)
	u, err := api.RenewUserApiKey(cur_user.Id)

	if err != nil {
		t.Errorf("RenewUserApiKey failed. Error: " + err.Error())
	} else {
		if u.Api.Key == cur_user.Api.Key {
			t.Errorf("Unable to renew user key.")
		} else {
			cur_user = u
			api = New(u.Api.Key, BaseUrl)
			SetToken(u.Api.Key)
			setEnvironmentVar("ONEANDONE_TOKEN", u.Api.Key)
		}
	}
}

func TestModifyUser(t *testing.T) {
	getc_user.Do(get_current_user)

	fmt.Printf("Modifying user '%s'...\n", cur_user.Name)
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(999999)
	new_pass := fmt.Sprintf("%d^&*bbYhgf%djv;HF", r, r)
	new_desc := tudesc + "_updated"
	new_mail := "test@oneandone.com"
	req := UserRequest{
		Description: new_desc,
		Email:       new_mail,
		Password:    new_pass,
		//		State: "DISABLED",
	}
	u, err := api.ModifyUser(cur_user.Id, &req)

	if err != nil {
		t.Errorf("ModifyUser failed. Error: " + err.Error())
		return
	}

	if u.Description != new_desc {
		t.Errorf("User description not updated.")
	}
	if u.Email != new_mail {
		t.Errorf("User email not updated.")
	}
}

func TestDeleteUser(t *testing.T) {
	t.Skip("TestDeleteUser is skipped at the moment.")
	getc_user.Do(get_current_user)

	fmt.Printf("Deleting user '%s'...\n", cur_user.Name)
	u, err := api.DeleteUser(cur_user.Id)

	if err != nil {
		t.Errorf("DeleteUser failed. Error: " + err.Error())
		return
	}
	if u != nil {
		u, err = api.GetUser(u.Id)

		if u != nil {
			t.Errorf("Unable to delete the user.")
		} else {
			cur_user = nil
		}
	}
}
