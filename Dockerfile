# Create the intermediate builder image.
FROM golang:1.15 as builder

# Docker is copying directory contents so we need to copy them in same directories.
WORKDIR /go/src/github.com/dstdfx/bookish-spork
COPY cmd cmd
COPY internal internal
COPY vendor vendor
COPY .git .git

# Other files can be copied into the WORKDIR.
COPY ["go*", "./"]

# Build the static application binary.
RUN BUILD_GIT_COMMIT=$(git rev-parse HEAD) \
    BUILD_GIT_TAG=$(git describe --abbrev=0) \
    BUILD_DATE=$(date +%Y%m%d) \
    GO111MODULE=on CGO_ENABLED=0 GOOS=linux \
    go build -mod=vendor -a -installsuffix cgo \
    -ldflags \
    "-X github.com/dstdfx/bookish-spork/cmd/bookish-spork/app.buildGitCommit=${BUILD_GIT_COMMIT} \
    -X github.com/dstdfx/bookish-spork/cmd/bookish-spork/app.buildGitTag=${BUILD_GIT_TAG} \
    -X github.com/dstdfx/bookish-spork/cmd/bookish-spork/app.buildDate=${BUILD_DATE}" \
    -o bookish-spork ./cmd/bookish-spork/bookishspork.go

# Create the final small image.
FROM alpine:3.12

RUN apk update && apk upgrade && \
    apk add --no-cache \
    ca-certificates wget && \
    rm -rf /var/cache/apk/*

COPY --from=builder /go/src/github.com/dstdfx/bookish-spork/bookish-spork /usr/bin/bookish-spork

EXPOSE 63100 63101

ENTRYPOINT ["/usr/bin/bookish-spork"]
