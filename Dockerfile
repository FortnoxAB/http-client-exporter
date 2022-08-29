FROM gcr.io/distroless/static-debian11:nonroot
COPY http-client-exporter /http-client-exporter
USER nonroot
ENTRYPOINT ["/http-client-exporter"]
