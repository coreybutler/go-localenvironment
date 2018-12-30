package localenvironment

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

func CreateEnvFile(t *testing.T, vars map[string]string) {
	t.Helper()
	b, err := json.Marshal(vars)
	if err != nil {
		t.Fatalf("Unexpected error marshalling vars: '%v'", err)
	}
	err = ioutil.WriteFile("env.json", b, 0644)
	if err != nil {
		t.Fatalf("Unexpected error writing file: '%v'", err)
	}
}

func DeleteEnvFile(t *testing.T) {
	t.Helper()
	err := os.Remove("env.json")
	if err != nil {
		t.Fatalf("Unexpected error deleting file: '%v'", err)
	}
}

func TestLocalEnvironment(t *testing.T) {
	CreateEnvFile(t, map[string]string{
		"TEST": "Success",
	})

	defer DeleteEnvFile(t)

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

func TestHandlesMissingFile(t *testing.T) {
	err := Apply()
	if err != nil {
		t.Fatalf("Unexpected error received from Apply: '%v'", err)
	}
}
