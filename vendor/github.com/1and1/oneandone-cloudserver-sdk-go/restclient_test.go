package oneandone

import (
	"testing"
)

func TestCreateUrl_1(t *testing.T) {
	api := New("token", "http://test.de/v1")

	result := createUrl(api)
	if result != "http://test.de/v1" {
		t.Errorf("Failed to create url.")
	}
}

func TestCreateUrl_2(t *testing.T) {
	api := New("token", "http://test.de/v1")

	result := createUrl(api, "servers")
	if result != "http://test.de/v1/servers" {
		t.Errorf("Failed to create url.")
	}
}

func TestCreateUrl_3(t *testing.T) {
	api := New("token", "http://test.de/v1")

	result := createUrl(api, "servers", 1)
	if result != "http://test.de/v1/servers/1" {
		t.Errorf("Failed to create url.")
	}
}

func TestAppendQueryParams_1(t *testing.T) {
	params := map[string]interface{}{
		"foo": "bar",
	}
	result := appendQueryParams("http://test/", params)
	if result != "http://test/?foo=bar" {
		t.Errorf("Failed to create url with query parameters.")
	}
}

func TestAppendQueryParams_2(t *testing.T) {
	params := map[string]interface{}{
		"foo":  "bar",
		"size": 5,
	}
	result := appendQueryParams("http://test/", params)
	if result != "http://test/?foo=bar&size=5" {
		t.Errorf("Failed to create url with query parameters.")
	}
}

func TestAppendQueryParams_3(t *testing.T) {
	params := map[string]interface{}{}
	result := appendQueryParams("http://test/", params)
	if result != "http://test/" {
		t.Errorf("Failed to create url with query parameters.")
	}
}

func TestAppendQueryParams_UrlEncode_1(t *testing.T) {
	params := map[string]interface{}{
		"test": "1&2=3",
	}
	result := appendQueryParams("http://test/", params)
	if result != "http://test/?test=1%262%3D3" {
		t.Errorf("Failed to create url with query parameters.")
	}
}

func TestAppendQueryParams_UrlEncode_2(t *testing.T) {
	params := map[string]interface{}{
		"test!": "1&2=3",
	}
	result := appendQueryParams("http://test/", params)
	if result != "http://test/?test%21=1%262%3D3" {
		t.Errorf("Failed to create url with query parameters.")
	}
}
