// localenvironment populates environment variables if they exist in
// `env.json` or an alternative JSON file.
//
// For example, the directory structure might look like:
//
//	> dir
//		- main.exe (or whatever executable you generate with `go build`)
//		- env.json
//
// Consider the following `env.json` file (in the current working directory):
//
//		{
//	  	"MY_API_KEY": "12345"
//		}
//
// ... and the main.go file:
//
//		package main
//
//		import (
//	  		"os"
//	  		"log"
//	  		"github.com/coreybutler/go-localenvironment"
//		)
//
//		func main() {
//	  		err := localenvironment.Apply() // Apply the env.json attributes to the environment variables.
//	  		if err != nil {
//	    		log.Print(err)
//	  		}
//
//	  		apiKey := os.Getenv("MY_API_KEY")
//
//	  		log.Printf("My API key is %s.", apiKey)
//		}
//
// Running this will output `My API key is 12345.`. The same Go app can be run in any
// directory, each with a different `env.json` file, potentially yielding different results.
//
// ## Nested JSON:
// This module automatically expands nested JSON properties (flattens attribute names).
// For example:
//
//	{
//	  "a": {
//	    "b": {
//			    "c": "something"
//			 }
//		 }
//	}
//
// The data structure above would be flattened into an environment
// variable called `A_B_C` whose value is `something`.
package localenvironment

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	// "log"
	"reflect"
	"strconv"
)

var knownEnvVars map[string]string

// Apply key/value pairs from a local `env.json` or `.env` file (if it exists).
// Each key will be available as an environment variable.
func Apply() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	return ApplyFiles(filepath.Join(cwd, "env.json"), filepath.Join(cwd, ".env"))
}

// ApplyFile will process any properly formatted JSON file,
// allowing for configuration files with a different name from "env.json".
func ApplyFile(path string) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	knownEnvVars, err = parse(raw)
	if err != nil {
		return err
	}

	for key, value := range knownEnvVars {
		os.Setenv(key, value)
	}

	return nil
}

// ApplyFiles will loop through a list of file paths
// and apply each file to the environment.
func ApplyFiles(paths ...string) error {
	for _, path := range paths {
		e := ApplyFile(path)
		if e != nil {
			return e
		}
	}

	return nil
}

func parsekv(content []byte) map[string]interface{} {
	data := make(map[string]interface{})
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		keyval := strings.Split(line, "=")
		if len(keyval) != 2 {
			continue
		}

		data[strings.TrimSpace(keyval[0])] = strings.TrimSpace(keyval[1])
	}

	return data
}

func parse(content []byte) (map[string]string, error) {
	data := make(map[string]string)
	keypairs := make(map[string]interface{})

	err := json.Unmarshal(content, &keypairs)
	if err != nil {
		return mapKeyPairs(data, parsekv(content)), nil
	}

	return mapKeyPairs(data, keypairs), nil
}

func mapKeyPairs(data map[string]string, keypairs map[string]interface{}) map[string]string {
	for key, value := range keypairs {
		switch reflect.TypeOf(value).String() {
		case "bool":
			data[key] = strconv.FormatBool(value.(bool))
		case "float64":
			data[key] = strconv.FormatFloat(value.(float64), 'f', -1, 64)
		case "int64":
			data[key] = strconv.FormatInt(value.(int64), 10)
		case "string":
			data[key] = value.(string)
		case "interface":
			// log.Print(key, " is an interface.")
		default:
			extra := mapKeyPairs(make(map[string]string), value.(map[string]interface{}))
			for subkey, subval := range extra {
				data[key+"_"+subkey] = subval
			}
		}
	}

	return data
}

// Clear removes environment variables applied with localenvironment.
func Clear() {
	for key := range knownEnvVars {
		os.Unsetenv(key)
	}
}
