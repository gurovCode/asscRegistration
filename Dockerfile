FROM golang:1.13-alpine as alpine-ca

RUN apk add git build-base bash
ADD certs/* /usr/local/share/ca-certificates/
RUN apk add -U --no-cache ca-certificates && update-ca-certificates
RUN GO111MODULE=off go get -tags 'postgres' -u github.com/golang-migrate/migrate/cmd/migrate \
    && go get -u github.com/onsi/ginkgo/ginkgo \
	&& go get -u github.com/onsi/gomega/... \
	&& go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

ENV GOOS=linux
ENV GOARCH=amd64

RUN mkdir /app
WORKDIR /app
