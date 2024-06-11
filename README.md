# go-localenvironment

[![Version](https://img.shields.io/github/tag/coreybutler/go-localenvironment.svg)](https://github.com/coreybutler/go-localenvironment)
[![GoDoc](https://godoc.org/github.com/coreybutler/go-localenvironment?status.svg)](https://godoc.org/github.com/coreybutler/go-localenvironment)

Apply environment variables sourced from `env.json`, `.env`, or custom sources in the current working directory. This is a port of the [coreybutler/localenvironment](https://github.com/coreybutler/localenvironment) Node.js module.

**Basic Usage**

```go
import "github.com/coreybutler/go-localenvironment"

func main() {
  localenvironment.Apply()
}
```

## Overview

This package provides a lightweight approach to environment variable management within an application. It will look for a file called `env.json` or `.env`. If one or both files exist, each key is added as an environment variable, accessible via the [os.Getenv](https://golang.org/pkg/os/#Getenv) method. If the file does not exist, it is silently ignored.

### Example Usage

Consider the following directory contents:

```sh
> dir
  - env.json
  - .env
  - main.go
```

**env.json**

```json
{
  "MY_API_KEY": "12345"
}
```

**.env**

```sh
MY_OTHER_API_KEY=ABC67890
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
  err := localenvironment.Apply() // Apply the env.json/.env attributes to the environment variables.
  if err != nil {
    log.Printf("Error: %s", err)
  }

  apiKey := os.Getenv("MY_API_KEY")
  otherApiKey := os.Getenv("MY_OTHER_API_KEY")

  log.Printf("My API key is %s. The other is %s.\n", apiKey, otherApiKey)
}
```

Running `main.go` will log `My API key is 12345. The other is ABC67890.`.

Variables applied by localenvironment are _ephemerally added_ to the environment. They are not persisted to the user/system environment variable store. If `MY_API_URL` is defined as a user/system variable, it will still be available whether localenvironment applies `env.json`/`.env` variables or not.

In the case of a conflicting variable, localenvironment will override values at runtime only. For example, if `MY_API_KEY=abcde` is defined as a user environment variable, localenvironment will override the value of `MY_API_KEY` with the value from the `env.json`/`.env` file (i.e. it will be `12345`).

### Variable Flattening/Expansion (JSON only)

Nested JSON properties are automatically flattened.

For example:

```javascript
{
  "a": {
    "b": {
      "c": "something"
    }
  }
}
```

The data structure above would be flattened into an environment variable called `A_B_C`, with a value of `something`.

### Custom Sources (other than `env.json`/`.env`)

It is possible to specify an alternative JSON/KV files using the `ApplyFile` method:

```go
package main

import (
    "os"
    "log"
    "github.com/coreybutler/go-localenvironment"
)

func main() {
    err := localenvironment.ApplyFile("/path/to/config.json") // Apply your own attributes to the environment variables.
    if err != nil {
      log.Printf("Error: %s", err)
    }

    apiKey := os.Getenv("...")
}
```

### Multiple Sources

While less common, there are circumstances where it is useful to apply more than one configuration to the environment variables. The `ApplyFiles` method supports this.


```go
package main

import (
    "os"
    "log"
    "github.com/coreybutler/go-localenvironment"
)

func main() {
    err := localenvironment.ApplyFiles(
      "/path/to/config.json",
      "/path/to/other.json"
    ) // Apply your own attributes to the environment variables.

    if err != nil {
      log.Printf("Error: %s", err)
    }

    apiKey := os.Getenv("...")
}
```

---

## Why?

This is a simple approach for adding environment variable management to code. Alternatively (or in addition), it may be desirable to define environment variables in the build process. For this, [QuikGo](https://github.com/quikdev/go) provides a robust option leveraging similar techniques.
