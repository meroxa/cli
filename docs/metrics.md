## Metrics

Meroxa CLI currently uses [Cased](https://cased.com/) to track user behaviour with the purpose of making informed decisions around sunsetting commands, introducing aliases when we notice our customers  _claim_ for one, and a long etc.

These metrics are pushed automatically for any command that's nested under the root `meroxa` command, and its code configuration lives in [/cmd/meroxa/global/metrics](/cmd/meroxa/global/metrics.go).

We don't log flag values, and only their names so no sensitive information is actually being tracked.

### Configuration

CLI uses some environment variables you can modify in your configuration file (you can check the one you're using by doing `meroxa config`) which will have different effects:

- [`CASED_DEBUG`](/cmd/meroxa/global/metrics.go#L46-L48) - set it to `true` if you want to see confirmation on whether events are actually published.  
- [`CASED_PUBLISH_KEY`](/cmd/meroxa/global/metrics.go#L37) - set the publish key you'd like to use. When developing locally, you should use the API Key with the `publish` type on a testing audit trail in Cased.com. Alternatively, if you don't set this API Key, you would test publishing through a proxy (e.g.: Staging API for example). 
- [`PUBLISH_METRICS`](/cmd/meroxa/global/metrics.go#L42-L44) - set this to `false` if you don't want to publish any metric to cased. When modifying how these metrics could be sent to Cased, you can set that to `stdout` to see the event in standard out. 


