package localenvironment

import "testing"
import "os"

func TestLocalEnvironment(t *testing.T) {
	Apply()

	nonexistant := os.Getenv("I_DO_NOT_EXIST")
	if nonexistant != "" {
		t.Errorf("An environment variable (I_DO_NOT_EXIST) is recognized when it shouldn't be.")
	}

	expectedvalue := os.Getenv("TEST")
	if expectedvalue != "Success" {
		t.Errorf("Unexpected value received for TEST attribute: '%s'", expectedvalue)
	}
}
