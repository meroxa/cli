FROM gcr.io/distroless/static
USER nonroot:nonroot
WORKDIR /app
COPY app.json /app
COPY {{.AppName}}.cross {{.AppName}}
ENTRYPOINT ["/app/{{.AppName}}", "--serve"]