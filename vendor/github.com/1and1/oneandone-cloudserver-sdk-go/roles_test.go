package oneandone

import (
	"fmt"
	"math/rand"
	"reflect"
	"sync"
	"testing"
	"time"
)

var (
	set_role  sync.Once
	role_name string
	test_role *Role
)

// Helper functions

func create_role() *Role {
	rand.Seed(time.Now().UnixNano())
	rint := rand.Intn(1000)
	role_name = fmt.Sprintf("TestRole_%d", rint)

	fmt.Printf("Creating role '%s'...\n", role_name)
	role_id, role, err := api.CreateRole(role_name)
	if err != nil {
		fmt.Printf("Unable to create a role. Error: %s", err.Error())
		return nil
	}
	if role_id == "" || role.Id == "" {
		fmt.Printf("Unable to create role '%s'.", role_name)
		return nil
	}

	api.WaitForState(role, "ACTIVE", 2, 30)

	return role
}

func set_role_once() {
	test_role = create_role()
}

// /roles tests

func TestCreateRole(t *testing.T) {
	set_role.Do(set_role_once)

	if test_role == nil {
		t.Errorf("CreateRole failed.")
	} else {
		if test_role.Name != role_name {
			t.Errorf("Wrong name of the role.")
		}
		if test_role.Permissions == nil {
			t.Errorf("Missing role permissions.")
		}
	}
}

func TestGetRole(t *testing.T) {
	set_role.Do(set_role_once)

	fmt.Printf("Getting role '%s'...\n", role_name)
	role, err := api.GetRole(test_role.Id)

	if err != nil {
		t.Errorf("GetRole failed. Error: " + err.Error())
		return
	}
	if role.Id != test_role.Id {
		t.Errorf("Wrong role ID.")
	}
	if test_role.Permissions == nil {
		t.Errorf("Missing role permissions.")
	}
}

func TestListRoles(t *testing.T) {
	set_role.Do(set_role_once)
	fmt.Println("Listing all roles...")

	res, err := api.ListRoles()
	if err != nil {
		t.Errorf("ListRoles failed. Error: " + err.Error())
	}
	if len(res) == 0 {
		t.Errorf("No role found.")
	}

	res, err = api.ListRoles(1, 2, "name", "", "id,name")

	if err != nil {
		t.Errorf("ListRoles with parameter options failed. Error: " + err.Error())
		return
	}
	if len(res) == 0 {
		t.Errorf("No role found.")
	}
	if len(res) != 2 {
		t.Errorf("Wrong number of objects per page.")
		return
	}
	if res[0].Id == "" || res[0].Name == "" {
		t.Errorf("Filtering parameters failed.")
	}
	if res[0].State != "" || res[0].Permissions != nil {
		t.Errorf("Filtering parameters failed.")
	}
	if res[0].Name > res[1].Name {
		t.Errorf("Sorting list of roles failed.")
	}

	res, err = api.ListRoles(0, 0, "", test_role.Name, "")

	if err != nil {
		t.Errorf("ListRoles with parameter options failed. Error: " + err.Error())
		return
	}
	if len(res) != 1 || res[0].Name != test_role.Name {
		t.Errorf("Search parameter failed.")
	}
}

func TestModifyRole(t *testing.T) {
	set_role.Do(set_role_once)

	fmt.Printf("Modifying role '%s'...\n", role_name)
	new_name := test_role.Name + "_updated"
	new_desc := test_role.Description + "_updated"

	role, err := api.ModifyRole(test_role.Id, new_name, new_desc, "")

	if err != nil {
		t.Errorf("ModifyRole failed. Error: " + err.Error())
		return
	}
	if role.Id != test_role.Id {
		t.Errorf("Wrong role ID.")
	}
	if role.Name != new_name {
		t.Errorf("Wrong role name.")
	}
	if role.Description != new_desc {
		t.Errorf("Wrong role description.")
	}

	test_role = role
}

func TestGetRolePermissions(t *testing.T) {
	set_role.Do(set_role_once)

	fmt.Printf("Getting role's permissions...\n")
	p, err := api.GetRolePermissions(test_role.Id)

	if err != nil {
		t.Errorf("GetRolePermissions failed. Error: " + err.Error())
		return
	}
	if p == nil || p.Backups == nil || p.Firewalls == nil || p.Images == nil || p.Invoice == nil ||
		p.IPs == nil || p.LoadBalancers == nil || p.Logs == nil || p.MonitorCenter == nil ||
		p.MonitorPolicies == nil || p.PrivateNetworks == nil || p.Roles == nil || p.Servers == nil ||
		p.SharedStorage == nil || p.Usages == nil || p.Users == nil || p.VPNs == nil {
		t.Errorf("The role '%s' is missimg some permissions objects.", test_role.Name)
	}
}

func TestListRoleUsers(t *testing.T) {
	roles, err := api.ListRoles()
	if len(roles) > 0 {
		for _, r := range roles {
			if len(r.Users) > 0 {
				fmt.Printf("Getting role's users...\n")
				users, err := api.ListRoleUsers(r.Id)

				if err != nil {
					t.Errorf("ListRoleUsers failed. Error: " + err.Error())
					return
				}
				if !reflect.DeepEqual(r.Users, users) {
					t.Errorf("ListRoleUsers failed. Users do not match.")
				}
			}
			break
		}
	} else {
		t.Errorf("ListRoles failed. Error: " + err.Error())
	}
}

func TestAssignRemoveRoleUser(t *testing.T) {
	set_role.Do(set_role_once)

	users, _ := api.ListUsers(0, 0, "", "go_test_user", "")
	if len(users) > 0 {
		fmt.Printf("Assigning user '%s' to role '%s'...\n", users[0].Name, role_name)
		usl := []string{users[0].Id}
		role, err := api.AssignRoleUsers(test_role.Id, usl)

		if err != nil {
			t.Errorf("AssignRoleUsers failed. Error: " + err.Error())
			return
		}
		if len(role.Users) != 1 {
			t.Errorf("AssignRoleUsers failed.")
			return
		}
		if role.Users[0].Id != users[0].Id {
			t.Errorf("Wrong user assigned to the role.")
		}

		// Removing the user
		role, err = api.RemoveRoleUser(test_role.Id, role.Users[0].Id)

		if err != nil {
			t.Errorf("RemoveRoleUser failed. Error: " + err.Error())
			return
		}
		if len(role.Users) != 0 {
			t.Errorf("RemoveRoleUser failed.")
		}
	} else {
		t.Errorf("No user found for the role assignment.")
	}
}

func TestGetRoleUser(t *testing.T) {
	roles, err := api.ListRoles()
	if len(roles) > 0 {
		for _, r := range roles {
			if len(r.Users) > 0 {
				fmt.Printf("Getting role's user '%s'...\n", r.Name)
				user, err := api.GetRoleUser(r.Id, r.Users[0].Id)

				if err != nil {
					t.Errorf("GetRoleUser failed. Error: " + err.Error())
					return
				}
				if r.Users[0].Id != user.Id || r.Users[0].Name != user.Name {
					t.Errorf("User's ID or Name does not match.")
				}
			}
			break
		}
	} else {
		t.Errorf("ListRoles failed. Error: " + err.Error())
	}
}

func TestCloneRole(t *testing.T) {
	set_role.Do(set_role_once)

	fmt.Printf("Cloning role '%s'...\n", role_name)
	cn := role_name + " Clone"
	role, err := api.CloneRole(test_role.Id, cn)

	if err != nil {
		t.Errorf("CloneRole failed. Error: " + err.Error())
		return
	}
	if role.Name != cn {
		t.Errorf("Wrong name of the clone role.")
	}
	// cleaning
	api.DeleteRole(role.Id)
}

func TestDeleteRole(t *testing.T) {
	set_role.Do(set_role_once)

	fmt.Printf("Deleting role '%s'...\n", test_role.Name)
	role, err := api.DeleteRole(test_role.Id)

	if err != nil {
		t.Errorf("DeleteRole failed. Error: " + err.Error())
		return
	}

	api.WaitUntilDeleted(role)
	role, err = api.GetRole(role.Id)

	if role != nil {
		t.Errorf("Unable to delete the role.")
	} else {
		test_role = nil
	}
}
