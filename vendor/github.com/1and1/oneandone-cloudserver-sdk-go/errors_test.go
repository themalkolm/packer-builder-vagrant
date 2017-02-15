package oneandone

import (
	"testing"
)

func TestCreateError(t *testing.T) {
	err := apiError{httpStatusCode: 404, message: "Not found"}

	if err.HttpStatusCode() != 404 {
		t.Errorf("Wrong HTTP status code.")
	}
	if err.Message() != "Not found" {
		t.Errorf("Wrong HTTP error message.")
	}
}
