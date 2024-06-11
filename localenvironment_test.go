package localenvironment

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
	// "log"
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

func CreateKVEnvFile(t *testing.T, vars map[string]string) {
	t.Helper()
	var b []byte
	for key, value := range vars {
		b = append(b, []byte(key+"="+value+"\n")...)
	}

	err := ioutil.WriteFile(".env", b, 0644)
	if err != nil {
		t.Fatalf("Unexpected error writing file: '%v'", err)
	}
}

func DeleteKVEnvFile(t *testing.T) {
	t.Helper()
	err := os.Remove(".env")
	if err != nil {
		t.Fatalf("Unexpected error deleting file: '%v'", err)
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
		t.Errorf("Unexpected value received for TEST. Expected 'Success', received '%s'", expectedvalue)
	}

	Clear()

	clearedValue := os.Getenv("TEST")
	if clearedValue != "" {
		t.Errorf("Unexpected value received for TEST after Clear: '%s' (expected an empty value)", clearedValue)
	}
}

func TestLocalEnvironmentKV(t *testing.T) {
	CreateKVEnvFile(t, map[string]string{
		"TEST": "Success",
	})

	defer DeleteKVEnvFile(t)

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
		t.Errorf("Unexpected value received for TEST. Expected 'Success', received '%s'", expectedvalue)
	}

	Clear()

	clearedValue := os.Getenv("TEST")
	if clearedValue != "" {
		t.Errorf("Unexpected value received for TEST after Clear: '%s' (expected an empty value)", clearedValue)
	}
}

func TestNestedParsing(t *testing.T) {
	content := `{
		"a": {
			"b": {
				"c": "ok"
			}
		},
		"simple": "test",
		"int": 1,
		"dec": 1.5,
		"boolean": true
	}`

	data, _ := parse([]byte(content))

	v, ok := data["a_b_c"]
	if !ok {
		t.Error("Expected a variable called a_b_c. This does not exist.")
	}

	if v != "ok" {
		t.Errorf("Expected the flattened value to be 'ok', received '%v'", v)
	}
}

func TestHandlesMissingFile(t *testing.T) {
	err := Apply()
	if err != nil {
		t.Fatalf("Unexpected error received from Apply: '%v'", err)
	}
}
