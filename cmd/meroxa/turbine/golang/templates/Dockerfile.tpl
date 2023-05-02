FROM golang:{{.GoVersion}} as builder
WORKDIR /builder
COPY . .
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -tags server -o {{.AppName}} ./...

FROM gcr.io/distroless/static
USER nobody
WORKDIR /app
COPY --from=builder /builder/app.json /app
COPY --from=builder /builder/{{.AppName}} /app

ENTRYPOINT ["/app/{{.AppName}}", "server", "-serve-func"]