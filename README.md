<p>   
    <a href="https://meroxa.io/" target="_blank">
      <img alt="Meroxa logo" width="300" src="https://meroxa-public-assets.s3.us-east-2.amazonaws.com/meroxa-logo-sticker.svg">
    </a>
    <br>
</p>

We believe that anyone should be empowered to leverage real-time data. Using the Meroxa CLI, you can build data infrastructure in minutes not months.

[Website](https://meroxa.io) |
[Documentation](https://docs.meroxa.com/) |
[Installation Guide](https://docs.meroxa.com/cli/installation-guide) |
[Contribution Guide](CONTRIBUTING.md) |
[Twitter](https://twitter.com/meroxadata)  

[![Support Server](https://img.shields.io/discord/828680256877363200.svg?label=Meroxa%20Community&logo=Discord&colorB=7289da&style=for-the-badge)](https://discord.meroxa.com)


## Documentation

Meroxa is documented publicly in https://docs.meroxa.com/, but on each build we also generate Markdown files for each command, exposing the available commands and help for each one. Check it out at [docs/cmd/md/meroxa](docs/cmd/md/meroxa.md).

## Contributing

For a complete guide to contributing to the Meroxa CLI, see the [Contribution Guide](CONTRIBUTING.md).

## Installation Guide

Please follow the [installation instructions in the Meroxa Documentation](https://docs.meroxa.com/cli/installation-guide).

### Build and Install the Binaries from Source (Advanced Install)

Currently, we provide pre-built Meroxa binaries for macOS (Darwin) Windows, and Linux architectures.

See [Releases](https://github.com/meroxa/cli/releases).

If you run into any issues during compiling the binaries, checkout the [troubleshooting guide](#troubleshooting).

Prerequisite Tools:

* [Git](https://git-scm.com/)
* [Go](https://golang.org/dl/)

To build from source:

1. The CLI depends on [meroxa-go](github.com/meroxa/meroxa-go). To update vendoring the dependency, you'll need to run the following:

```
make gomod
```

2. Build CLI as `meroxa` binary:

```
make build
```

### Development Configuration
The [meroxa-dev CLI](https://github.com/meroxa/meroxa-dev#install) makes it convenient to set the necessary env vars for different dev environments.

## Release

A [goreleaser](https://github.com/goreleaser/goreleaser) GitHub Action is
configured to automatically build the CLI and cut a new release whenever a new
git tag is pushed to the repo.

* Tag - `git tag -a vX.X.X -m "<message goes here>"`
* Push - `git push origin vX.X.X`

With every release, a new Homebrew formula will be automatically updated on [meroxa/homebrew-taps](https://github.com/meroxa/homebrew-taps).

## Linting

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

## Tests

To run the test suite:

```
make test
```

## Shell Completion

If you want to enable shell completion manually, you can generate your own using our `meroxa completion` command.

Type `meroxa help completion` for more information, or simply take a look at our [documentation](docs/cmd/md/meroxa_completion.md).

## Development with local TurbineJS dependency

The CLI offers a mode, in which the local version of `turbine-js` is loaded by default which can be useful for developing the CLI and the JS dependency locally when contributing to this project.

You can run the cli in this mode, by adding the following environment variable into your `.envrc` at the root of this project's directory:

```bash
MEROXA_USE_LOCAL_TURBINE_JS=true
```

If you run into any `executable not found` issues, please ensure that you have cloned [turbine-js](https://github.com/meroxa/turbine-js) locally and
installed it globally via running `npm install -g .` in the turbine-js directory.

If you'd like to see this development mode extended with other dependencies or features, you can send us a feature request in our Github issues: https://github.com/meroxa/cli/issues

## Troubleshooting

### make gomod

#### Error: fatal: could not read Username for 'https://github.com': terminal prompts disabled

If you have already setup your access to `meroxa` repositories via [SSH](https://docs.github.com/en/github/authenticating-to-github/connecting-to-github-with-ssh)
you can enforce the download of the required dependencies listed in the [go.mod file](https://github.com/meroxa/cli/blob/master/go.mod) via `ssh` instead
of `https` to circumvent this error.

To do so, setup your `gitconfig`, by running:

```
git config --global url."git@github.com:".insteadOf "https://github.com"
```

Verify the correct setup, by running ` cat ~/.gitconfig`. You should see the following output:

```
[user]
	name = Jane Doe
	email = janedoe@gmail.com
[url "git@github.com:"]
	insteadOf = https://github.com
```

Run the `make gomod` command again, and you should see all dependencies being downloaded successfully.
