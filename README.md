## Meroxa CLI

All commands start with `meroxa` (the name of the binary).

### Examples

* Create Resource:
    ```
    meroxa create resource postgres --name mypg --url postgres://user:secret@localhost:5432/db
    ```
* Create Connection:
    ```
    meroxa create connection mypg --config '{"table.whitelist":"public.purchases"}'
    ```
* List Resources:
    ```
    meroxa list resources
    ```

### Build

Build CLI as `meroxa` binary:
```
go build -i -o meroxa .
```

### Vendor

The CLI depends on [meroxa-go](github.com/meroxa/meroxa-go) which is currently
a private repo. In order to build the CLI you must vendor that package with the
Go mod sum DB disabled:

```
GOSUMDB=off go mod vendor
```

### Update Dependencies

In order to update `meroxa-go` (while it is still private), you'll need to run
the following:

```
GOSUMDB=off go get -u github.com/meroxa/meroxa-go
```

### Release

A [goreleaser](https://github.com/goreleaser/goreleaser) Github Action is
configured to automatically build the CLI and cut a new release whenever a new
git tag is pushed to the repo.

* Tag - `git tag -a vX.X.X -m "<message goes here>"`
* Push - `git push origin vX.X.X`


### Test

_TODO_


