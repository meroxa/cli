import hashlib
import logging
import sys
import pdb

from turbine.src.turbine_app import RecordList, TurbineApp

logging.basicConfig(level=logging.INFO)


def anonymize(records: RecordList) -> RecordList:
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
    async def run(turbine: TurbineApp):
        try:
            # To configure your data stores as resources on the Meroxa Platform
            # use the Meroxa Dashboard, CLI, or Meroxa Terraform Provider.
            # For more details refer to: https://docs.meroxa.com/

            # Identify an upstream data store for your data app
            # with the `resources` function.
            # Replace `source_name` with the resource name the
            # data store was configured with on the Meroxa platform.
            source = await turbine.resources("source_name")

            # Specify which upstream records to pull
            # with the `records` function.
            # Replace `collection_name` with a table, collection,
            # or bucket name in your data store.
            # If you need additional connector configurations, replace '{}'
            # with the key and value, i.e. {"incrementing.field.name": "id"}
            records = await source.records("collection_name", {})

            # Specify which secrets in environment variables should be passed
            # into the Process.
            # Replace 'PWD' with the name of the environment variable.
            #
            # turbine.register_secrets("PWD")

            # Specify what code to execute against upstream records
            # with the `process` function.
            # Replace `anonymize` with the name of your function code.
            anonymized = await turbine.process(records, anonymize)

            # Identify a downstream data store for your data app
            # with the `resources` function.
            # Replace `destination_name` with the resource name the
            # data store was configured with on the Meroxa platform.
            destination_db = await turbine.resources("destination_name")

            # Specify where to write records downstream
            # using the `write` function.
            # Replace `collection_archive` with a table, collection,
            # or bucket name in your data store.
            # If you need additional connector configurations, replace '{}'
            # with the key and value, i.e. {"behavior.on.null.values": "ignore"}
            await destination_db.write(anonymized, "collection_archive", {})
        except Exception as e:
            print(e, file=sys.stderr)
