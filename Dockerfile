FROM registry.digitalocean.com/francisco/golang-base:1.18 AS build
ARG ARCH="amd64"
ARG OS="linux"
ARG PROJECT="token-gateway"

WORKDIR $GOPATH/src/github.com
COPY ./$PROJECT ./$PROJECT

ENV GO111MODULE=on

WORKDIR $GOPATH/src/github.com/$PROJECT
RUN go mod tidy && sh build.sh

## release
FROM alpine:3.14
ARG PROJECT="token-gateway"
COPY --from=build /go/src/github.com/$PROJECT/output /app

WORKDIR /app
EXPOSE 8080
CMD ["./bootstrap.sh"]
