package localenvironment

import "testing"
import "os"

func TestLocalEnvironment(t *testing.T) {
	err := Apply()
	if err != nil {
		t.Fatalf("Unexpected error received from Apply: '%v'", err)
	}

	nonexistant := os.Getenv("I_DO_NOT_EXIST")
	if nonexistant != "" {
		t.Errorf("An environment variable (I_DO_NOT_EXIST) is recognized when it shouldn't be.")
	}

	expectedvalue := os.Getenv("TEST")
	if expectedvalue != "Success" {
		t.Errorf("Unexpected value received for TEST attribute: '%s'", expectedvalue)
	}

	Clear()

	clearedValue := os.Getenv("TEST")
	if clearedValue != "" {
		t.Errorf("Unexpected value received for TEST after Clear: '%s'", clearedValue)
	}
}
