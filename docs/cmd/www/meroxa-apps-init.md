---
createdAt: 
updatedAt: 
title: "meroxa apps init"
slug: meroxa-apps-init
url: /cli/cmd/meroxa-apps-init/
---
## meroxa apps init

Initialize a Meroxa Data Application

```
meroxa apps init [APP_NAME] [--path pwd] --lang js|go [flags]
```

### Examples

```
meroxa apps init my-app --path ~/code --lang jsmeroxa apps init my-app --lang go # will be initialized in a dir called my-app in the current directorymeroxa apps init my-app --lang go --path $GOPATH/src/github.com/my.orgmeroxa apps init my-app --lang go --skip-mod-init # will not initialize the new go modulemeroxa apps init my-app --lang go --mod-vendor # will initialize the new go module and download dependencies to the vendor directory
```

### Options

```
  -h, --help            help for init
  -l, --lang string     language to use (js|go) (required)
      --mod-vendor      whether to download modules to vendor or globally while initializing a Go application
      --path string     path where application will be initialized (current directory as default)
      --skip-mod-init   whether to run 'go mod init' while initializing a Go application
```

### Options inherited from parent commands

```
      --cli-config-file string   meroxa configuration file
      --debug                    display any debugging information
      --json                     output json
      --timeout duration         set the duration of the client timeout in seconds (default 10s)
```

### SEE ALSO

* [meroxa apps](/cli/cmd/meroxa-apps/)	 - Manage Meroxa Data Applications

