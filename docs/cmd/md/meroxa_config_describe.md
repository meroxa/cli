## meroxa config describe

Show Meroxa CLI configuration details

### Synopsis

This command will return the content of your configuration file where you could find your `access_token` and then `refresh_token` for your `meroxa login`. They're stored in the home directory on your machine. On Unix, including MacOS, it's stored in `$HOME`, and on Windows is stored in `%USERPROFILE%`.

```
meroxa config describe [flags]
```

### Examples

```
$ meroxa config describe
Using meroxa config located in "/Users/my-name/Library/Application Support/meroxa/config.env

access_token: c0f928b...c337a0d
actor: user@email.com
actor_uuid: c0f928ba-d40e-40c5-a7fa-cf281c337a0d
refresh_token: c337a0d...c0f928b

```

### Options

```
  -h, --help   help for describe
```

### Options inherited from parent commands

```
      --api-url string           API url
      --cli-config-file string   meroxa configuration file
      --debug                    display any debugging information
      --json                     output json
      --timeout duration         set the duration of the client timeout in seconds (default 10s)
```

### SEE ALSO

* [meroxa config](meroxa_config.md)	 - Manage your Meroxa CLI configuration

