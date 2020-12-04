FROM golang:1.15.5-buster AS builder

ENV GO111MODULE="on"
ENV APP_HOME /usr/src/app
WORKDIR $APP_HOME

COPY go.sum go.mod  ./
RUN go mod download

COPY ./ ./
