# go-localenvironment

[![Version](https://img.shields.io/github/tag/coreybutler/go-localenvironment.svg)](https://github.com/coreybutler/go-localenvironment)
[![GoDoc](https://godoc.org/github.com/coreybutler/go-localenvironment?status.svg)](https://godoc.org/github.com/coreybutler/go-localenvironment)
[![Build Status](https://travis-ci.org/coreybutler/go-localenvironment.svg?branch=master)](https://travis-ci.org/coreybutler/go-localenvironment)

Apply environment variables sourced from a `env.json` in the current working directory. This is a port of the [coreybutler/localenvironment](https://github.com/coreybutler/localenvironment) Node.js module.

**Install** `go get github.com/coreybutler/go-localenvironment`

#### Overview

This package enables a lightweight approach to environment management within an application. It will look for a file called `env.json`. If the file exists, each key is added as as an environment variable, accessible via the [os.Getenv](https://golang.org/pkg/os/#Getenv) method. If the file does not exist, it is silently ignored.

### Example Usage

Consider the following directory contents:

```sh
> dir
  - env.json
  - main.exe (or whatever executable you generate with `go build`)
```

**env.json**

```json
{
  "MY_API_KEY": "12345"
}
```

**main.go**:

```go
package main

import (
  "os"
  "log"
  "github.com/coreybutler/go-localenvironment"
)

func main() {
  err := localenvironment.Apply() // Apply the env.json attributes to the environment variables.
  if err != nil {
    panic(err)
  }

  apiKey := os.Getenv("MY_API_KEY")

  log.Printf("My API key is %s.", apiKey)
}
```

Running `main.exe` (or equivalent binary) will log `My API key is 12345.`. The same Go app can be run in any directory, each with a different `env.json` file, potentially yielding different results. It's also possible to change a value in the `env.json` file, yielding a different result the next time the app is executed.

Variables applied by localenvironment are _added_ to the environment. If `MY_API_URL` is defined somewhere else, it will still be available whether localenvironment applies `env.json` variables or not.

In the case of a conflicting variable, localenvironment will override existing values at runtime only. In the scenario where `MY_API_KEY=abcde` is defined beforehand, the localenvironment variable will override the value of `MY_API_KEY` (i.e. it will be `12345`).

---

#### Why?

Building configurable systems can be more time consuming than necessary. Sometimes it's just a
nice smooth workflow to define environment variables in a JSON file. This allows every developer
to tweak environment variables while working on their own local instance of the code, without having to write complicated build tooling to configure a bunch of environment variables.

We use this pretty regularly when building applications that will ultimately run in Docker, as
background daemons/services, and even command line tools. This also helps prevent publishing
sensitive configuration information (like API keys and secrets).

There are also programs written that function differently _if an environment variable exists_. This
module won't fire an error if it doesn't find the `env.json`, so it's easy to manipulate the environment simply by commenting out some code or changing a file name.

Your mileage will vary, but we've found this to be much simpler than writing shell/batch files and other wrappers to manage environment variables, especially when there are a large number of variables.
