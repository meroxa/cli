<img src="https://files.readme.io/29ef4f4-small-meroxaLogo2x.png" alt="Meroxa" width="300">  

We believe that anyone should be empowered to leverage real-time data. Using the Meroxa CLI, you can build data infrastructure in minutes not months.

[Website](https://meroxa.io) |
[Documentation](https://docs.meroxa.com/) |
[Installation Guide](https://docs.meroxa.com/docs/installation-guide) |
[Contribution Guide](CONTRIBUTING.md) |
[Twitter](https://twitter.com/meroxadata)


## Installation Guide

Use the [installation instructions in the Meroxa documentation](http://docs.meroxa.com/).

### Build and Install the Binaries from Source (Advanced Install)

Currently, we provide pre-built Meroxa binaries for macOS (Darwin) Windows, and Linux architectures.

See [Releases](https://github.com/meroxa/cli/releases).

#### Prerequisite Tools

* [Git](https://git-scm.com/)
* [Go](https://golang.org/dl/)

### Vendor

The CLI depends on [meroxa-go](github.com/meroxa/meroxa-go) (which is currently a private repo). To update vendoring the dependency, you'll need to run the following:

```
make gomod
```

#### Build

Build CLI as `meroxa` binary:

```
make build
```

#### Linting

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

### Documentation

Our Meroxa CLI is documented publicly in https://docs.meroxa.com/docs, but on each build we also generate Markdown files for each command, exposing the available commands and help for each one. Check it out at [docs/commands/meroxa](docs/commands/meroxa.md).

## Contributing
For a complete guide to contributing to the Meroxa CLI, see the [Contribution Guide](CONTRIBUTING.md).

## Release

A [goreleaser](https://github.com/goreleaser/goreleaser) GitHub Action is
configured to automatically build the CLI and cut a new release whenever a new
git tag is pushed to the repo.

* Tag - `git tag -a vX.X.X -m "<message goes here>"`
* Push - `git push origin vX.X.X`

