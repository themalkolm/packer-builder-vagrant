package oneandone

import (
	"fmt"
	"testing"
)

// /pricing tests

func TestGetPricing(t *testing.T) {
	fmt.Println("Getting pricing")
	pricing, err := api.GetPricing()

	if err != nil {
		t.Errorf("GetPricing failed. Error: " + err.Error())
		return
	}
	if pricing.Currency == "" {
		t.Errorf("Missing pricing currency.")
	}
	if pricing.Plan == nil {
		t.Errorf("Missing pricing plan.")
		return
	}
	if pricing.Plan.Image == nil {
		t.Errorf("Missing image pricing.")
	}
	if pricing.Plan.SharedStorage == nil {
		t.Errorf("Missing shared storage pricing.")
	}
	if len(pricing.Plan.PublicIPs) == 0 {
		t.Errorf("Missing public IP pricing.")
	}
	if len(pricing.Plan.SoftwareLicenses) == 0 {
		t.Errorf("Missing software license pricing.")
	}
	if pricing.Plan.Servers == nil {
		t.Errorf("Missing server pricing.")
		return
	}
	if len(pricing.Plan.Servers.FixedServers) == 0 {
		t.Errorf("Missing fixed server pricing.")
	}
	if len(pricing.Plan.Servers.FlexServers) == 0 {
		t.Errorf("Missing flex server pricing.")
	}
}
