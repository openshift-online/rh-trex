FROM registry.access.redhat.com/ubi9/ubi-minimal:9.2-750.1697534106

RUN \
    microdnf install -y \
    util-linux \
    && \
    microdnf clean all

COPY \
    ocm-example-service \
    /usr/local/bin/

EXPOSE 8000

ENTRYPOINT ["/usr/local/bin/ocm-example-service", "serve"]

LABEL name="ocm-example-service" \
      vendor="Red Hat" \
      version="0.0.1" \
      summary="OCM Example Service API" \
      description="OCM Example Service API"
