// Pull in fixture data from the fixtures file
const demo = require("./fixtures/demo-cdc.json");
// Pull in the Anonymize function from the base app
const { App } = require("./index.js");
const { RecordsArray } = require("@meroxa/turbine-js-framework");

// If you using convenience functions like `.get` or `.set` on your records...
// You can use this helper in your tests to wrap your fixture data
function recordsHelper(collectionName) {
  const records = new RecordsArray();

  demo[collectionName].forEach((record) => {
    records.pushRecord(record);
  });

  return records;
}

// This example unit test was built using QUnit, a JavaScript testing framework
// However, you may use any testing framework of your choice
// To learn more about how to use this testing framework
// Refer to the QUnit documentation https://qunitjs.com/intro
// To run this example unit test, use `npm test`
QUnit.module("My data app", () => {
  QUnit.test("anonymize function works on `customer_email`", (assert) => {
    const app = new App();
    // Take records from the fixture with the defined collection name
    const records = recordsHelper("collection_name");

    // Apply function against the wrapped fixture data
    const results = app.anonymize(records);

    // Test the raw output of the function applied to your fixture data
    assert.strictEqual(
      results[0]._rawValue.payload.customer_email,
      "~~~1072928786~~~"
    );
  });
});
