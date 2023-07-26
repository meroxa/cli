# frozen_string_literal: true

require "rubygems"
require "bundler/setup"
require "turbine_rb"

class MyApp
  def call(app)
    # To configure resources for your production datastores
    # on Meroxa, use the Dashboard, CLI, or Terraform Provider
    # For more details refer to: http://docs.meroxa.com/
    #
    # Identify the upstream datastore with the `resource` function
    # Replace `demopg` with the resource name configured on Meroxa
    database = app.resource(name: "demopg")

    # Specify which upstream records to pull
    # with the `records` function
    # Replace `collection_name` with a table, collection,
    # or bucket name in your data store.
    # If a configuration is needed for your source,
    # you can pass it as a second argument to the `records` function. For example:
    # database.records(collection: "collection_name", configs: {"incrementing.column.name" => "id"})
    records = database.records(collection: "collection_name")

    # Register secrets to be available in the function:
    # app.register_secrets("MY_ENV_TEST")

    # Register several secrets at once:
    # app.register_secrets(["MY_ENV_TEST", "MY_OTHER_ENV_TEST"])

    # Specify the code to execute against `records` with the `process` function.
    # Replace `Passthrough` with your desired function.
    # Ensure desired function matches `Passthrough`'s' function signature.
    processed_records = app.process(records: records, process: Passthrough.new)

    # Specify where to write records using the `write` function.
    # Replace `collection_archive` with whatever data organisation method
    # is relevant to the datastore (e.g., table, bucket, collection, etc.)
    # If additional connector configs are needed, provided another argument. For example:
    # database.write(
    #   records: processed_records,
    #   collection: "collection_archive",
    #   configs: {"behavior.on.null.values" => "ignore"})
    database.write(records: processed_records, collection: "collection_archive")
  end
end

class Passthrough < TurbineRb::Process
  def call(records:)
    puts "got records: #{records}"
    # To get the value of unformatted records, use record .value getter method
    # records.map { |r| puts r.value }
    #
    # To transform unformatted records, use record .value setter method
    # records.map { |r| r.value = "newdata" }
    #
    # To get the value of json formatted records, use record .get method
    # records.map { |r| puts r.get("message") }
    #
    # To transform json formatted records, use record .set methods
    # records.map { |r| r.set('message', 'goodbye') }
    records
  end
end

TurbineRb.register(MyApp.new)
