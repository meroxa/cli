# cased-go

A Cased client for Go applications in your organization to control and monitor the access of information within your organization.

![Build Status](https://github.com/cased/cased-go/workflows/cased/badge.svg)
[![go.dev](https://img.shields.io/badge/go.dev-pkg-007d9c.svg?style=flat)](https://pkg.go.dev/github.com/cased/cased-go)

## Overview

- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
  - [Publishing events to Cased](#publishing-events-to-cased)
  - [Masking & filtering sensitive information](#masking--filtering-sensitive-information)
- [Contributing](#contributing)

## Installation

To make cased-go available for use in your package, run:

```sh
$ go get github.com/cased/cased-go
```

After installing the cased-go package and configuring cased-go, the first thing you'll want to do is start publishing audit events:

```go
package main

import "github.com/cased/cased-go"

func main() {
	p := cased.NewPublisher(
		cased.WithPublishKey("publish_test_1mY8qb355NWIa3uY00H2fk7elpT"),
	)
	cased.SetPublisher(p)

	ae := cased.AuditEvent{
		"action":          "user.login",
		"actor":           cased.NewSensitiveValue("John Doe", "username"),
		"actor_id":        "User;1",
		"organization":    "cased",
		"organization_id": "Organization;1",
	}
	cased.Publish(ae)

	fmt.Printf("%+v\n", ae)
}
```

cased-go publishes audit events asynchronously by default. To ensure no audit events are dropped at the end of your process, you'll want to make sure `cased.Flush` is called at some point in your shutdown process.

```go
package main

import "github.com/cased/cased-go"

func main() {
	p := cased.NewPublisher(
		cased.WithPublishKey("publish_test_1mY8qb355NWIa3uY00H2fk7elpT"),
	)
	cased.SetPublisher(p)

	// The process will wait 30 seconds to publish all events to Cased before
	// exiting the process.
	defer cased.Flush(30 * time.Second)

	ae := cased.AuditEvent{
		"action":          "user.login",
		"actor":           cased.NewSensitiveValue("John Doe", "username"),
		"actor_id":        "User;1",
		"organization":    "cased",
		"organization_id": "Organization;1",
	}
	cased.Publish(ae)

	fmt.Printf("%+v\n", ae)
}
```

You've now installed cased-go properly and have published your first event. For more details on [publishing audit events](#publishing-events-to-cased) and [protecting sensitive values](#masking--filtering-sensitive-information), keep reading on.

## Configuration

All configuration options available in cased-go are available to be configured by an environment variable or manually.

```go
package main

import (
	"net/http"
	"time"

	"github.com/cased/cased-go"
)

func main() {
	p := cased.NewPublisher(
		// CASED_PUBLISH_URL=https://publish.cased.com
		cased.WithPublishURL("https://publish.cased.com"),

		// CASED_PUBLISH_KEY=publish_live_1mY8qb355NWIa3uY00H2fk7elpT
		cased.WithPublishKey("publish_live_1mY8qb355NWIa3uY00H2fk7elpT"),

		// CASED_DEBUG=1
		cased.WithDebug(true),

		// CASED_SILENCE=1
		cased.WithSilence(true),

		// You can configure your own client or re-use an existing HTTP client from
		// your application.
		cased.WithHTTPClient(&http.Client{}),

		// You can configure your own transport or re-use an existing HTTP transport
		// from your application.
		cased.WithHTTPTransport(&http.Transport{}),

		// CASED_HTTP_TIMEOUT=10s
		cased.WithHTTPTimeout(10*time.Second),
		cased.WithTransport(cased.NewNoopHTTPTransport()),
	)
	cased.SetPublisher(p)

	// ...
}
```

## Usage

### Publishing events to Cased

There are two ways to publish your first Cased event.

**Manually**

```go
package user

import "github.com/cased/cased-go"

type Organization struct {
	ID   string
	Name string
}

type User struct {
	ID           string
	Username     string
	Organization Organization
}

func (u *User) login(password, passwordConfirmation string) error {
	ae := cased.AuditEvent{
		"action":          "user.login",
		"actor":           cased.NewSensitiveValue(u.Username, "username"),
		"actor_id":        u.ID,
		"organization":    u.Organization.Name,
		"organization_id": u.Organization.ID,
	}
	defer cased.Publish(ae)

	if password != passwordConfirmation {
		ae["action"] = "user.failed_login"
		return errors.New("invalid password")
	}

	return nil
}
```

### Masking & filtering sensitive information

If you are handling sensitive information on behalf of your users you should consider masking or filtering any sensitive information.

```go
package user

import "github.com/cased/cased-go"

type Address struct {
	Street  string
	City    string
	State   string
	ZipCode string
}

type User struct {
	ID       string
	Username string
	Address  *Address
}

func (a *Address) Create() error {
	cased.Publish(cased.AuditEvent{
		"action":   "address.create",
		"actor":    cased.NewSensitiveValue(u.Username, "username"),
		"actor_id": u.ID,
		"location": map[string]interface{}{
			"street":   cased.NewSensitiveValue(u.Address.Street, "street"),
			"city":     cased.NewSensitiveValue(u.Address.City, "city"),
			"state":    cased.NewSensitiveValue(u.Address.State, "state"),
			"zip_code": cased.NewSensitiveValue(u.Address.ZipCode, "zip-code"),
		},
	})

	return nil
}
```

### Disable publishing events

Although rare, there may be times where you wish to disable publishing events to Cased. You can configure it using an environment variable or in the client.

```go
package main

import (
	"net/http"
	"time"

	"github.com/cased/cased-go"
)

func main() {
	p := cased.NewPublisher(
		// CASED_SILENCE=1
		cased.WithSilence(true),
	)
	cased.SetPublisher(p)
	defer cased.Flush(10 * time.Second)

	// This audit event will not get published to Cased.
	cased.Publish(cased.AuditEvent{
		"action":             "user.login",
		"actor":              "dewski",
		"actor_id":           "user_1dsGbftbx1c47iU8c7BzUcKJRcD",
		"organization":       "Cased",
		"organization_id":    "org_1dsGTnNZLzgwb1alwS2szN0KUo5",
		"request_id":         "27d62d1869e5f9826acc6cfd80edca90",
		"request_url":        "https://app.cased.com/saml/consume",
		"request_user_agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_1_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36",
	})
}
```

Or you can configure the entire process to disable publishing events.

```
CASED_SILENCE=1 go run main.go
```

## Contributing

1. Fork it ( https://github.com/cased/cased-go/fork )
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create a new Pull Request
