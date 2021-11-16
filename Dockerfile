FROM scratch

COPY http-client-exporter /http-client-exporter

ENTRYPOINT ["/http-client-exporter"]
