## Meroxa CLI

## Our commands

All commands start with `meroxa` (the name of the binary).

### Build

Build CLI as `meroxa` binary:

```
make build
```

### Linting

If you want to make sure everything's correct before pushing to GitHub, you'll need to install [`golangci-lint`](https://golangci-lint.run/) and run:

```
$ golangci-lint run
```

Example:

```
❯ golangci-lint run
cmd/display.go:60:6: `appendCell` is unused (deadcode)
func appendCell(cells []*simpletable.Cell, text string) []*simpletable.Cell {
     ^
```

### Vendor

The CLI depends on [meroxa-go](github.com/meroxa/meroxa-go) which is currently
a private repo. To update vendoring the dependency, you'll need to run the following:

```
make gomod
```

### Release

A [goreleaser](https://github.com/goreleaser/goreleaser) GitHub Action is
configured to automatically build the CLI and cut a new release whenever a new
git tag is pushed to the repo.

* Tag - `git tag -a vX.X.X -m "<message goes here>"`
* Push - `git push origin vX.X.X`

### Documentation

Our Meroxa CLI is documented publicly in https://docs.meroxa.com/docs, but on each build we also generate Markdown files for each command, exposing the available commands and help for each one. Check it out at [docs/commands/meroxa](docs/commands/meroxa.md).