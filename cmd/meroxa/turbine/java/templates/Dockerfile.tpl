FROM quay.io/quarkus/ubi-quarkus-mandrel-builder-image:22.3-java17 as builder
WORKDIR /builder
COPY . .

RUN ./mvnw clean package -Pnative

FROM debian:bullseye
USER nobody
WORKDIR /app

COPY --from=builder /builder/app.json /app
COPY --from=builder /builder/target/*-runner /app/runner

ENTRYPOINT ["/app/runner", "-Dturbine.mode=deploy"]
