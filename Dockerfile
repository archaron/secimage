# Builder
FROM golang:alpine AS builder
COPY .  /go/src/github.com/archaron/secimage
WORKDIR /go/src/github.com/archaron/secimage
RUN \
    export VERSION=$(git rev-parse --verify HEAD) && \
    export BUILD=$(date -u +%s%N) && \
    export CGO_ENABLED=0 && \
    export LDFLAGS="-w -s -X main.Version=${VERSION} -X main.BuildTime=${BUILD} -extldflags \"-static\"" && \
    go build -v -ldflags "${LDFLAGS}" -o /go/bin/secimage .

# Executable image
FROM scratch
WORKDIR /
COPY --from=builder /go/bin/secimage /secimage
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/archaron/secimage/config.yml /config.yml
EXPOSE 8300 8301 8302
CMD ["/secimage"]
