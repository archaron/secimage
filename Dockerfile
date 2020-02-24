# Builder
FROM golang:1.13-alpine AS builder

RUN apk add --no-cache gcc musl-dev

COPY .  /go/src/github.com/archaron/secimage

WORKDIR /go/src/github.com/archaron/secimage

ARG REPO=main
ARG VERSION=dev

ADD https://github.com/chai2010/webp/archive/master.zip webp.zip

RUN unzip webp.zip \
  && cp -r ./webp-master/internal ./vendor/github.com/chai2010/webp/
RUN ls ./vendor/github.com/chai2010/webp/
RUN set -x \
    && export BUILD=$(date -u +%s%N) \
    && export CGO_ENABLED=1 \
    && export LDFLAGS="-w -s -X ${REPO}/misc.Version=${VERSION} -X ${REPO}/misc.BuildTime=${BUILD} -extldflags \"-static\"" \
    && go build -v -ldflags "${LDFLAGS}" -o /go/bin/secimage .

    #  -mod=vendor

# Executable image
FROM scratch

WORKDIR /

COPY --from=builder /go/bin/secimage /bin/secimage
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/archaron/secimage/config.yml /config.yml

EXPOSE 8300 8301 8302

CMD ["/bin/secimage"]
