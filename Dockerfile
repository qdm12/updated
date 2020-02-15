ARG ALPINE_VERSION=3.11
ARG GO_VERSION=1.13.7

FROM alpine:${ALPINE_VERSION} AS alpine
RUN apk --update add ca-certificates tzdata

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder
ENV CGO_ENABLED=0
RUN apk --update add git
WORKDIR /tmp/gobuild
COPY go.mod go.sum ./
RUN go mod download 2>&1
COPY cmd/updated/main.go cmd/app/main.go
COPY internal ./internal
COPY pkg ./pkg
RUN go build -ldflags="-s -w" -o app cmd/app/main.go

FROM alpine:${ALPINE_VERSION}
ARG BUILD_DATE
ARG VCS_REF
LABEL \
    org.opencontainers.image.authors="quentin.mcgaw@gmail.com" \
    org.opencontainers.image.created=$BUILD_DATE \
    org.opencontainers.image.version="" \
    org.opencontainers.image.revision=$VCS_REF \
    org.opencontainers.image.url="https://github.com/qdm12/updated" \
    org.opencontainers.image.documentation="https://github.com/qdm12/updated/blob/master/README.md" \
    org.opencontainers.image.source="https://github.com/qdm12/updated" \
    org.opencontainers.image.title="updated" \
    org.opencontainers.image.description="Updated updates periodically files locally or to a Git repository"
COPY --from=alpine --chown=1000 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=alpine --chown=1000 /usr/share/zoneinfo /usr/share/zoneinfo
COPY --chown=1000 known_hosts /known_hosts
ENV \
    OUTPUT_DIR=./files \
    PERIOD=24h \
    RESOLVE_HOSTNAMES=no \
    HTTP_TIMEOUT=3s \
    LOG_ENCODING=console \
    LOG_LEVEL=info \
    TZ=America/Montreal \
    GIT=no \
    GIT_URL= \
    SSH_KEY=./key \
    SSH_KEY_PASSPHRASE= \
    SSH_KNOWN_HOSTS=./known_hosts \
    NAMED_ROOT_MD5=5ff5afb5c009f6198b5af9a8a0fa51e5 \
    ROOT_ANCHORS_SHA256=45336725f9126db810a59896ae93819de743c416262f79c4444042c92e520770 \
    GOTIFY_URL= \
    GOTIFY_TOKEN= \
    NODE_ID=0
ENTRYPOINT ["/updated"]
#HEALTHCHECK --interval=10s --timeout=5s --start-period=5s --retries=2 CMD ["/updated","healthcheck"]
USER 1000
COPY --from=builder --chown=1000 /tmp/gobuild/app /updated
