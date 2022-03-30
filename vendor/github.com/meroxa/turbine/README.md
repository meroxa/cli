## Turbine - Meroxa Data App Framework

## Build

Building the project `go build ./...` automatically compiles the binary using the `local` runner. The local runner
leverages local fixtures for sample data.

In order to build for the `platform`, you must use the build tag `platform`:
```shell
go build -tags platform ./...
```

### Fixtures

By default, the local runner looks for sample data in a directory named `fixtures`. Sample data is JSON formatted, with
multiple tables nested.

Example:
```json
{
  "users": {
    "100": {
      "id": "100",
      "username": "alice",
      "email": "alice@example.com"
    }
  },
  "user_activity": {
    "1": {
      "id": "1",
      "user_id": "100",
      "email": "user@example.com",
      "action": "registered"
    },
    "2": {
      "id": "2",
      "user_id": "100",
      "email": "user@example.com",
      "action": "logged in"
    }
  }
}
```

### Usage

When built for the `platform` runtime, a number of CLI flags are introduced.

* `function <name> <args>` - Triggers the function named and passes the provided args.
* `serve <name>` - Wraps the function name in a gRPC service (designed for use by funtime).
* `listfunctions` - Returns the list of functions registered.
* `deploy` - Deploys the data app on the Meroxa Data Platform.
