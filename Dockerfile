# Builder
FROM golang:1-alpine AS builder

COPY .  /go/src/github.com/archaron/secimage

WORKDIR /go/src/github.com/archaron/secimage

ARG REPO=main
ARG VERSION=dev

RUN set -x \
    && export BUILD=$(date -u +%s%N) \
    && export CGO_ENABLED=0 \
    && export LDFLAGS="-w -s -X ${REPO}/misc.Version=${VERSION} -X ${REPO}/misc.BuildTime=${BUILD} -extldflags \"-static\"" \
    && go build -v -mod=vendor -ldflags "${LDFLAGS}" -o /go/bin/secimage .

# Executable image
FROM scratch

WORKDIR /

COPY --from=builder /go/bin/secimage /bin/secimage
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/archaron/secimage/config.yml /config.yml

EXPOSE 8300 8301 8302

CMD ["/bin/secimage"]
