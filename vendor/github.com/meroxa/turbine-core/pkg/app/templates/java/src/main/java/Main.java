package org.example;

import io.meroxa.turbine.Runner;
import io.meroxa.turbine.Turbine;
import io.meroxa.turbine.TurbineApp;
import io.meroxa.turbine.TurbineRecord;

import java.util.List;

import static java.time.LocalDateTime.now;

public class MyApp implements TurbineApp {
    public static void main(String[] args) {
        Runner.start(new MyApp());
    }

    @Override
    public void setup(Turbine turbine) {
        turbine
            .resource("test-pg-source")
            .read("user_activity", null)
            .process(this::process)
            .writeTo(turbine.resource("test-mysql-destination"), "user_activity_enriched", null);
    }

    private List<TurbineRecord> process(List<TurbineRecord> records) {
        return records.stream()
            .map(r -> {
                var copy = r.copy();
                copy.setPayload("This was processed at " + now());

                return copy;
            })
            .toList();
    }
}
