# Turbine
...

### Fixtures

Fixtures are JSON-formatted samples of data records you can use while locally developing your Turbine app. Whether CDC or non-CDC-formatted data records, fixtures adhere to the following structure:

```json
{
  "collection_name": [
    {
      "key": "1",
      "value": {
		  "schema": {
			  //...
		  },
		  "payload": {
			  //...
		  }
		}
	}
  ]
```

- `collection_name` — Identifies the name of the records or events you are streaming to your data app.
- `key` — Denotes one or more sample records within a fixture file. `key` is always a string.
- `value` — Holds the `schema` and `payload` of the sample data record.
- `schema` — Comes as part of your sample data record. `schema` describes the record or event structure.
- `payload` — Comes as part of your sample data record. `payload` describes what about the record or event changed.

Your newly created data app should have a `demo-cdc.json` and `demo-non-cdc.json` in the `/fixtures` directory as examples to follow.

### Testing

Testing should follow standard NodeJS development practices. Included in the repo is a sample jest configuration, but you can use any testing framework you would use to test Node apps.

## Documentation && Reference

The most comprehensive documentation for Turbine and how to work with Turbine apps is on the Meroxa site: [https://docs.meroxa.com/](https://docs.meroxa.com)

## Contributing

Check out the [/docs/](./docs/) folder for more information on how to contribute.
