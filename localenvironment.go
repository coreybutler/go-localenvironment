// Package localenvironment creates environment variables if they exist in
// `env.json`.
// For example, the directory structure might look like:
//
// 	> dir
// 		- main.exe (or whatever executable you generate with `go build`)
// 		- env.json
//
// Consider the following `env.json` file (in the current working directory):
//
// 	{
//   	"MY_API_KEY": "12345"
// 	}
//
// ... and the main.go file:
//
// 	package main
//
// 	import (
//   		"os"
//   		"log"
//   		"github.com/coreybutler/go-localenvironment"
// 	)
//
// 	func main() {
//   		err := localenvironment.Apply() // Apply the env.json attributes to the environment variables.
//   		if err != nil {
//     		log.Print(err)
//   		}
//
//   		apiKey := os.Getenv("MY_API_KEY")
//
//   		log.Printf("My API key is %s.", apiKey)
// 	}
//
// Running this will output `My API key is 12345.`. The same Go app can be run in any
// directory, each with a different `env.json` file, potentially yielding different results.
package localenvironment

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

var knownEnvVars map[string]string

// Apply key/value pairs from a local `env.json` file (if it exists).
// Each key will be available as an environment variable.
func Apply() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	raw, err := ioutil.ReadFile(filepath.Join(cwd, "env.json"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	err = json.Unmarshal(raw, &knownEnvVars)
	if err != nil {
		return err
	}

	for key, value := range knownEnvVars {
		os.Setenv(key, value)
	}

	return nil
}

// Clear removes environment variables applied with localenvironment.
func Clear() {
	for key := range knownEnvVars {
		os.Unsetenv(key)
	}
}
