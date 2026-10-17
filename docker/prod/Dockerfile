# syntax=docker/dockerfile:1
#
# File Browser — production image, built entirely from source.
# Build from the repo root:  docker build -f docker/prod/Dockerfile .
# Or use:                    docker compose -f docker/prod/compose.yaml up -d --build

ARG NODE_VERSION=24
ARG GO_VERSION=1.25
ARG ALPINE_VERSION=3.23

# ---- 1. Frontend assets -------------------------------------------------
FROM node:${NODE_VERSION}-alpine AS frontend
ENV COREPACK_ENABLE_DOWNLOAD_PROMPT=0
RUN corepack enable
WORKDIR /src/frontend
COPY frontend/package.json frontend/pnpm-lock.yaml ./
RUN --mount=type=cache,target=/root/.local/share/pnpm/store \
    pnpm install --frozen-lockfile
COPY frontend/ ./
RUN pnpm exec vite build

# ---- 2. Go binary (frontend embedded) -----------------------------------
FROM golang:${GO_VERSION}-alpine AS backend
WORKDIR /src
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY . .
COPY --from=frontend /src/frontend/dist ./frontend/dist
ARG VERSION=dev
ARG COMMIT_SHA=unknown
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -trimpath \
      -ldflags="-s -w \
        -X github.com/filebrowser/filebrowser/v2/version.Version=${VERSION} \
        -X github.com/filebrowser/filebrowser/v2/version.CommitSHA=${COMMIT_SHA}" \
      -o /out/filebrowser .

# ---- 3. Runtime helpers (certs, mime types, tini, JSON.sh) --------------
FROM alpine:${ALPINE_VERSION} AS fetcher
RUN apk add --no-cache ca-certificates mailcap tini-static && \
    wget -qO /JSON.sh https://raw.githubusercontent.com/dominictarr/JSON.sh/0d5e5c77365f63809bf6e77ef44a1f34b0e05840/JSON.sh

# ---- 4. Final image -----------------------------------------------------
FROM busybox:1.37.0-musl

LABEL org.opencontainers.image.title="File Browser" \
      org.opencontainers.image.source="https://github.com/filebrowser/filebrowser" \
      org.opencontainers.image.licenses="Apache-2.0"

ARG UID=1000
ARG GID=1000

COPY --from=fetcher /sbin/tini-static /bin/tini
COPY --from=fetcher /JSON.sh /JSON.sh
COPY --from=fetcher /etc/ca-certificates.conf /etc/ca-certificates.conf
COPY --from=fetcher /etc/ca-certificates /etc/ca-certificates
COPY --from=fetcher /etc/ssl /etc/ssl
COPY --from=fetcher /etc/mime.types /etc/mime.types
COPY --from=backend /out/filebrowser /bin/filebrowser
COPY docker/prod/rootfs/ /

RUN addgroup -g "$GID" user && \
    adduser -D -H -u "$UID" -G user user && \
    mkdir -p /config /database /srv && \
    chown user:user /config /database /srv && \
    chmod 755 /init.sh /healthcheck.sh

USER user
VOLUME ["/srv", "/config", "/database"]
EXPOSE 80

HEALTHCHECK --start-period=5s --interval=30s --timeout=3s --retries=3 CMD ["/healthcheck.sh"]
ENTRYPOINT ["tini", "--", "/init.sh"]
