# Sets linux/amd64 in case it's not injected by older Docker versions
ARG BUILDPLATFORM=linux/amd64

ARG ALPINE_VERSION=3.14
ARG GO_VERSION=1.25
ARG XCPUTRANSLATE_VERSION=v0.6.0
ARG GOLANGCI_LINT_VERSION=v1.42.1

FROM alpine:${ALPINE_VERSION} AS alpine
RUN mkdir /files && \
    chown 1000 /files && \
    chmod 700 /files

FROM --platform=${BUILDPLATFORM} qmcgaw/xcputranslate:${XCPUTRANSLATE_VERSION} AS xcputranslate
FROM --platform=${BUILDPLATFORM} qmcgaw/binpot:golangci-lint-${GOLANGCI_LINT_VERSION} AS golangci-lint

FROM --platform=${BUILDPLATFORM} golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS base
ENV CGO_ENABLED=0
WORKDIR /tmp/gobuild
RUN apk --update add git g++
COPY --from=xcputranslate /xcputranslate /usr/local/bin/xcputranslate
COPY --from=golangci-lint /bin /go/bin/golangci-lint
COPY go.mod go.sum ./
RUN go mod download
COPY pkg/ ./pkg/
COPY cmd/ ./cmd/
COPY internal/ ./internal/

FROM base AS test
# Note on the go race detector:
# - we set CGO_ENABLED=1 to have it enabled
# - we installed g++ in the base stage to support the race detector
ENV CGO_ENABLED=1
ENTRYPOINT go test -race -coverpkg=./... -coverprofile=coverage.txt -covermode=atomic ./...

FROM base AS lint
COPY .golangci.yml ./
RUN golangci-lint run --timeout=10m

FROM base AS build
ARG TARGETPLATFORM
RUN GOARCH="$(xcputranslate translate -targetplatform=${TARGETPLATFORM} -field arch)" \
    GOARM="$(xcputranslate translate -targetplatform=${TARGETPLATFORM} -field arm)" \
    go build -trimpath -ldflags="-s -w \
    " -o app cmd/updated/main.go

FROM scratch
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
COPY --from=alpine --chown=1000 /files /files
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
    NAMED_ROOT_MD5=076cfeb40394314adf28b7be79e6ecb1 \
    ROOT_ANCHORS_SHA256=45336725f9126db810a59896ae93819de743c416262f79c4444042c92e520770 \
    SHOUTRRR_SERVICES=
ENTRYPOINT ["/updated"]
#HEALTHCHECK --interval=10s --timeout=5s --start-period=5s --retries=2 CMD ["/updated","healthcheck"]
USER 1000
COPY --from=build --chown=1000 /tmp/gobuild/app /updated
