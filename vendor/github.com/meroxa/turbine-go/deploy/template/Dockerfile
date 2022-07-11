FROM gcr.io/distroless/static
USER nobody
WORKDIR /app
COPY app.json /app
COPY {{.AppName}}.cross {{.AppName}}
ENTRYPOINT ["/app/{{.AppName}}", "--serve"]