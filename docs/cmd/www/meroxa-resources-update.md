---
createdAt: 
updatedAt: 
title: "meroxa resources update"
slug: meroxa-resources-update
url: /cli/cmd/meroxa-resources-update/
---
## meroxa resources update

Update a resource

### Synopsis

Use the update command to update various Meroxa resources.

```
meroxa resources update NAME [flags]
```

### Options

```
      --ca-cert string       trusted certificates for verifying resource
      --client-cert string   client certificate for authenticating to the resource
      --client-key string    client private key for authenticating to the resource
  -h, --help                 help for update
  -m, --metadata string      new resource metadata
      --name string          new resource name
      --password string      password
      --ssh-url string       SSH tunneling address
      --ssl                  use SSL
  -u, --url string           new resource url
      --username string      username
```

### Options inherited from parent commands

```
      --cli-config-file string   meroxa configuration file
      --debug                    display any debugging information
      --json                     output json
      --timeout duration         set the duration of the client timeout in seconds (default 10s)
```

### SEE ALSO

* [meroxa resources](/cli/cmd/meroxa-resources/)	 - Manage resources on Meroxa

