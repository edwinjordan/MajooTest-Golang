# syntax=docker/dockerfile:1.7

# Arguments with default value (for build).
ARG PLATFORM=linux/amd64
ARG NODE_VERSION=22
ARG GO_VERSION=1.25

FROM busybox:1.37-glibc as glibc

# -----------------------------------------------------------------------------
# Base image for building the Golang app.
# -----------------------------------------------------------------------------
FROM --platform=${PLATFORM} golang:${GO_VERSION}-bookworm AS base_go
ENV MOON_TOOLCHAIN_FORCE_GLOBALS=1 LEFTHOOK=0 CI=true
WORKDIR /srv

# Install system dependencies and Moon CLI (via npm for reliability)
RUN apt-get update && \
    apt-get -yqq --no-install-recommends install curl npm tini jq ca-certificates git && \
    npm install -g @moonrepo/cli && \
    which moon && moon --version && \
    apt-get -yqq autoremove && apt-get -yqq clean && rm -rf /var/lib/apt/lists/*

# -----------------------------------------------------------------------------
# Builder: install deps and build the application.
# -----------------------------------------------------------------------------
FROM base_go AS builder
ENV CGO_ENABLED=1

# Enable caching for Go modules
RUN go env -w GOCACHE=/go-cache
RUN go env -w GOMODCACHE=/gomod-cache

# Copy project sources (including .moon if you already have it)
COPY --link . .

# Build the binary
RUN go build -o /srv/apps/go-app/build/go-app .

# -----------------------------------------------------------------------------
# Runner: minimal production image
# -----------------------------------------------------------------------------
FROM --platform=${PLATFORM} gcr.io/distroless/cc-debian12 AS runner

# ----- Read application environment variables --------------------------------
ARG  DATABASE_URL SMTP_HOST SMTP_PORT SMTP_USERNAME SMTP_PASSWORD SMTP_EMAIL_FROM

# Copy the build output files from builder stage
COPY --from=builder --chown=nonroot:nonroot /srv/apps/go-app/build/go-app /srv/go-app

# Copy some necessary system utilities
COPY --from=base_go /usr/bin/tini /usr/bin/tini
COPY --from=glibc /usr/bin/env /usr/bin/env
COPY --from=glibc /bin/clear /bin/clear
COPY --from=glibc /bin/mkdir /bin/mkdir
COPY --from=glibc /bin/which /bin/which
COPY --from=glibc /bin/cat /bin/cat
COPY --from=glibc /bin/ls /bin/ls
COPY --from=glibc /bin/sh /bin/sh

# Define the host and port to listen on.
ARG SERVER_ENV=production HOST=0.0.0.0 PORT=8000
ENV SERVER_ENV=$SERVER_ENV TINI_SUBREAPER=true
ENV HOST=$HOST PORT=$PORT

WORKDIR /srv
USER nonroot:nonroot
EXPOSE $PORT

ENTRYPOINT ["/usr/bin/tini", "--"]
CMD ["/srv/go-app"]