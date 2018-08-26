# Builder
  FROM golang:alpine AS builder

    COPY . /go/src/github.com/archaron/secimage

    WORKDIR /go/src/github.com/archaron/secimage

    RUN set -x \
  && export CGO_ENABLED="0" \
  && export CGO_CFLAGS="-g -O9" \
  && export CGO_CXXFLAGS="-g -O9" \
  && export CGO_FFLAGS="-g -O9" \
  && export CGO_LDFLAGS="-g -O9" \
  && go build \
    -ldflags "-w -s -extldflags \"-static\"" \
    -gcflags '-m' \
    -gccgoflags '-O9' \
    -v \
    -o /go/bin/secimage ./

  # Executable image
    FROM scratch

    WORKDIR /

    COPY --from=builder /go/bin/secimage /secimage

    EXPOSE 8300 8301 8302

    CMD ["/secimage"]
