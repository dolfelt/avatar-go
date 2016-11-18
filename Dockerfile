FROM golang:1.7
MAINTAINER Daniel Olfelt "https://github.com/dolfelt"

ENV APP_PATH=/go/src/github.com/dolfelt/avatar-go
ENV GOOS linux
ENV GOARCH amd64
ENV CGO_ENABLED 0

RUN wget https://github.com/Masterminds/glide/releases/download/v0.12.3/glide-v0.12.3-linux-amd64.tar.gz
RUN tar -zxvf glide-v0.12.3-linux-amd64.tar.gz
RUN cd linux-amd64 && mv glide /usr/bin

WORKDIR $APP_PATH
