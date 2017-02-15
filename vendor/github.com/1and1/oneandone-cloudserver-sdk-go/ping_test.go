package oneandone

import (
	"fmt"
	"testing"
)

// ping tests

func TestPing(t *testing.T) {
	fmt.Println("PING...")
	// API client with no token
	client := New("", BaseUrl)
	pong, err := client.Ping()
	if err != nil {
		t.Errorf("Ping failed. Error: " + err.Error())
	}
	if len(pong) == 0 {
		t.Errorf("Empty PING response.")
		return
	}
	if pong[0] != "PONG" {
		t.Errorf("Invalid PING response.")
	}
}

func TestPingAuth(t *testing.T) {
	fmt.Println("PING with authorization check...")
	pong, err := api.PingAuth()
	if err != nil {
		t.Errorf("Ping with authorization check failed. Error: " + err.Error())
	}
	if len(pong) == 0 {
		t.Errorf("Empty PING authorization response.")
		return
	}
	if pong[0] != "PONG" {
		t.Errorf("Invalid PING authorization response.")
	}
}
