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
make build
```

### Vendor

The CLI depends on [meroxa-go](github.com/meroxa/meroxa-go) which is currently
a private repo. To update vendoring the dependency, you'll need to run the following:

```
make gomod
```

### Release

A [goreleaser](https://github.com/goreleaser/goreleaser) Github Action is
configured to automatically build the CLI and cut a new release whenever a new
git tag is pushed to the repo.

* Tag - `git tag -a vX.X.X -m "<message goes here>"`
* Push - `git push origin vX.X.X`


### Test

_TODO_


