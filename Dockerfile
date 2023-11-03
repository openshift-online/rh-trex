FROM registry.access.redhat.com/ubi9/ubi-minimal:9.2-750.1697534106

RUN \
    microdnf install -y \
    util-linux \
    && \
    microdnf clean all

COPY \
    rh-trex \
    /usr/local/bin/

EXPOSE 8000

ENTRYPOINT ["/usr/local/bin/rh-trex", "serve"]

LABEL name="rh-trex" \
      vendor="Red Hat" \
      version="0.0.1" \
      summary="RH TREX API" \
      description="RH TREX API"
