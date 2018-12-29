// Apply environment variables if they exist in `env.json`.
// For example, the directory structure might look like:
//
// ```sh
// > dir
//   - main.exe (or whatever executable you generate with `go build`)
//   - env.json
// ```
//
// Consider the following `env.json` file (in the current working directory):
//
// ```json
// {
//   "MY_API_KEY": "12345"
// }
// ```
//
// ... and the main.go file:
//
// ```go
// package main
//
// import (
//   "os"
//   "log"
//   "github.com/coreybutler/go-localenvironment"
// )
//
// func main() {
//   localenvironment.Apply() // Apply the env.json attributes to the environment variables.
//
//   apiKey := os.Getenv("MY_API_KEY")
//
//   log.Printf("My API key is %s.", apiKey)
// }
// ```
//
// Running this will output `My API key is 12345.`. The same Go app can be run in any
// directory, each with a different `env.json` file, potentially yielding different results.
package localenvironment

import (
  "os"
  "io/ioutil"
  "path/filepath"
  "encoding/json"
  "log"
)

var knownEnvVars map[string]interface{}

// Apply key/value pairs from a local `env.json` file (if it exists).
// Each key will be available as an environment variable.
func Apply () {
  cwd, err := os.Getwd()

  if err != nil {
    return
  }

  raw, fileError := ioutil.ReadFile(filepath.Join(cwd, "env.json"))

  if fileError != nil {
    return
  }

  json.Unmarshal(raw, &knownEnvVars)

  for key, value := range knownEnvVars {
    os.Setenv(key, value.(string))
  }
}

// Removes environment variables applied with localenvironment.
func Clear () {
  for key, value := range knownEnvVars {
    os.Unsetenv(key)
  }
}
