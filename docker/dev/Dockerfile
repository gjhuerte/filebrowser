# syntax=docker/dockerfile:1
#
# File Browser — development image (Node + Go + Air).
# Source code is bind-mounted at runtime; use docker/dev/compose.yaml.

ARG NODE_VERSION=24
ARG GO_VERSION=1.25

FROM golang:${GO_VERSION}-alpine AS go

FROM node:${NODE_VERSION}-alpine
COPY --from=go /usr/local/go /usr/local/go

ENV GOPATH=/go \
    PATH=/go/bin:/usr/local/go/bin:$PATH \
    CGO_ENABLED=0 \
    GOFLAGS=-buildvcs=false \
    COREPACK_ENABLE_DOWNLOAD_PROMPT=0

RUN apk add --no-cache git && \
    go install github.com/air-verse/air@latest && \
    corepack enable && \
    mkdir -p /srv /database

# Pre-fetch the exact pnpm version pinned in frontend/package.json
COPY frontend/package.json /tmp/pkg/package.json
RUN cd /tmp/pkg && corepack install && rm -rf /tmp/pkg

WORKDIR /app
