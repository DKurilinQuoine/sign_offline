FROM golang:alpine

ENV PACKAGES make git libc-dev bash gcc linux-headers eudev-dev

WORKDIR /go/src/sign_offline

COPY . .

RUN apk add --no-cache $PACKAGES