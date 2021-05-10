---
createdAt: 
updatedAt: 
title: "meroxa update resource"
slug: meroxa-update-resource
url: /cli/meroxa-update-resource/
---
## meroxa update resource

Update a resource

### Synopsis

Use the update command to update various Meroxa resources.

```
meroxa update resource NAME [flags]
```

### Options

```
      --ca-cert string       trusted certificates for verifying resource
      --client-cert string   client certificate for authenticating to the resource
      --client-key string    client private key for authenticating to the resource
  -h, --help                 help for resource
  -m, --metadata string      new resource metadata
      --name string          new resource name
      --password string      password
      --ssl                  use SSL
  -u, --url string           new resource url
      --username string      username
```

### Options inherited from parent commands

```
      --config string      config file
      --debug              display any debugging information
      --json               output json
      --timeout duration   set the client timeout (default 10s)
```

### SEE ALSO

* [meroxa update](/cli/meroxa-update/)	 - Update a component

