ARG VERSION=develop
ARG GIT_COMMIT=unknown
ARG GO_VERSION=1.23.2

FROM golang:${GO_VERSION}-alpine3.20 AS builder
ARG VERSION
ARG GIT_COMMIT
ARG GO_VERSION

WORKDIR /app

COPY go.mod go.sum .

RUN go mod download

COPY main.go /app/main.go
COPY api/ /app/api/
COPY cmd/ /app/cmd/
COPY internal/ /app/internal/

RUN CGO_ENABLED=0 go build -o /usr/local/bin/goarkitect \
    -ldflags="-s -w -X main.version=${VERSION} -X main.gitCommit=${GIT_COMMIT} -X main.buildTime=$(date -u +'%Y-%m-%dT%H:%M:%SZ') -X main.goVersion=${GO_VERSION} -X main.osArch=$(uname -m)" \
    .

FROM gcr.io/distroless/static AS static

LABEL maintainer="omissis"
# LABEL org.opencontainers.image.created
LABEL org.opencontainers.image.authors="omissis"
LABEL org.opencontainers.image.url="https://github.com/omissis/goarkitect"
LABEL org.opencontainers.image.documentation="https://github.com/omissis/goarkitect"
LABEL org.opencontainers.image.source="https://github.com/omissis/goarkitect"
# LABEL org.opencontainers.image.version
# LABEL org.opencontainers.image.revision
LABEL org.opencontainers.image.vendor="Omissis"
# LABEL org.opencontainers.image.licenses
# LABEL org.opencontainers.image.ref.name
LABEL org.opencontainers.image.title="goarkitect"
LABEL org.opencontainers.image.description="Describe and check architectural constraints of a project using a composable set of rules."
# LABEL org.opencontainers.image.base.digest
# LABEL org.opencontainers.image.base.name

USER nonroot:nonroot

ENTRYPOINT ["/goarkitect"]

FROM static AS dockerbuild

COPY --from=builder /usr/local/bin/goarkitect /goarkitect

FROM static AS goreleaser

COPY goarkitect /goarkitect
