## meroxa secrets create

Create a Turbine Secret

### Synopsis

This command will create a secret as promted by the user.'
After successful creation, the secret can be used in a connector. 


```
meroxa secrets create NAME --data '{}' [flags]
```

### Examples

```
meroxa secret create NAME
		          meroxa secret create NAME --data '{}'
		
```

### Options

```
      --data string   Secret's data, passed as a JSON string
  -h, --help          help for create
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

* [meroxa secrets](meroxa_secrets.md)	 - Manage Turbine Data Applications

