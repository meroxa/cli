
# Meroxa Data App

## Getting Started

The CLI will create a new folder with the name provided located in the directory where the command was issued. If you want to initialize the app somewhere else, you can append the `--path` flag to the command (`meroxa apps init testapp --lang py --path ~/anotherdir`). Once you enter the `testapp` directory, the contents will look like this:

```bash
$ tree testapp/
testapp
├── README.md
├── main.py
├── app.json
├── __init__.py
└── fixtures
    ├── demo-cdc.json
    └── demo-no-cdc.json
```

This will be a full-fledged Turbine app that can run. You can even run the tests using the command `meroxa apps run` in the root of the app directory. It provides just enough to show you what you need to get started.

### `main.py`

This configuration file is where you begin your Turbine journey. Any time a Turbine app runs, this is the entry point for the entire application. When the project is created, the file will look like this:

```python
# Dependencies of the example data app
import hashlib
import logging
import sys
import typing as t

from turbine.runtime import Record, Runtime

logging.basicConfig(level=logging.INFO)


def anonymize(records: t.List[Record]) -> t.List[Record]:
    logging.info(f"processing {len(records)} record(s)")
    for record in records:
        logging.info(f"input: {record}")
        try:
            payload = record.value["payload"]["after"]

            # Hash the email
            payload["email"] = hashlib.sha256(
                payload["email"].encode("utf-8")
            ).hexdigest()

            logging.info(f"output: {record}")
        except Exception as e:
            print("Error occurred while parsing records: " + str(e))
            logging.info(f"output: {record}")
    return records


class App:
    @staticmethod
    async def run(turbine: Turbine):
      try:
        source = await turbine.resources("source_name")

        records = await source.records("collection_name")

        anonymized = await turbine.process(records, anonymize)

        destination_db = await turbine.resources("destination_name")

        await destination_db.write(anonymized, "collection_archive")
      except Exception as e:
          print(e, file=sys.stderr)
```

Let's talk about the important parts of this code. Turbine apps have five functions that comprise the entire DSL. Outside of these functions, you can write whatever code you want to accomplish your tasks:

```python
async def run(turbine: Turbine):
```

`run` is the main entry point for the application. This is where you can initialize the Turbine framework. This is also the place where, when you deploy your Turbine app to Meroxa, Meroxa will use this as the place to boot up the application.

```python
source = await turbine.resources("source_name")
```

The `resources` function identifies the upstream or downstream system that you want your code to work with. The `source_name` is the string identifier of the particular system. The string should map to an associated identifier in your `app.json` to configure what's being connected to. For more details, see the `app.json` section.

```python
records = await source.records("collection_name")
```

Once you've got `resources` set up, you can now stream records from it, but you need to identify what records you want. The `records` function identifies the records or events that you want to stream into your data app.

```python
anonymized = await turbine.process(records, anonymize)
```

The `process` function is Turbine's way of saying, for the records that are coming in, I want you to process these records against a function. Once your app is deployed on Meroxa, Meroxa will do the work to take each record or event that does get streamed to your app and then run your code against it. This allows Meroxa to scale out your processing relative to the velocity of the records streaming in.

```python
await destination_db.write(anonymized, "collection_archive")
```

The `write` function is optional. It takes any records given to it and streams them to the downstream system. In many cases, you might not need to stream data to another system, but this gives you an easy way to do so.


### `app.json`

This file contains all of the options for configuring a Turbine app. Upon initialization of an app, the CLI will scaffold the file for you with available options:

```json
{
  "name": "testapp",
  "language": "python",
  "environment": "common",
  "resources": {
    "source_name": "fixtures/demo-cdc.json"
  }
}
```

* `name` - The name of your application. This should not change after app initialization.
* `language` - Tells Meroxa what language the app is upon deployment.
* `environment` - "common" is the only available environment. Meroxa does have the ability to create isolated environments but this feature is currently in beta.
* `resources` - These are the named integrations that you'll use in your application. The `source_name` needs to match the name of the resource that you'll set up in Meroxa using the `meroxa resources create` command or via the Dashboard. You can point to the path in the fixtures that'll be used to mock the resource when you run `meroxa apps run`.

### Fixtures

Fixtures are JSON-formatted samples of data records you can use while locally developing your Turbine app. Whether CDC or non-CDC-formatted data records, fixtures adhere to the following structure:

```json
{
  "collection_name": [
    {
      "key": "1",
      "value": {
  		  "schema": {
  			  /* ... */
  		  },
  		  "payload": {
  			  /* ... */
  		  }
      }
    }
  ]
}
```
* `collection_name` — Identifies the name of the records or events you are streaming to your data app.
* `key` — Denotes one or more sample records within a fixture file. `key` is always a string.
* `value` — Holds the `schema` and `payload` of the sample data record.
* `schema` — Comes as part of your sample data record. `schema` describes the record or event structure.
* `payload` — Comes as part of your sample data record. `payload` describes what about the record or event changed.

Your newly created data app should have a `demo-cdc.json` and `demo-non-cdc.json` in the `/fixtures` directory as examples to follow.

### Testing

Testing should follow standard Python development practices.

## Documentation && Reference

The most comprehensive documentation for Turbine and how to work with Turbine apps is on the Meroxa site: [https://docs.meroxa.com/](https://docs.meroxa.com)

## Example apps

[See what a sample python data app looks like using our framework](https://github.com/meroxa/turbine-py-examples)
