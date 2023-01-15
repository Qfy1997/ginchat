FROM golang:1.18-alpine
MAINTAINER qfy
WORKDIR /go/src/docker_go
COPY . .
RUN go env -w GOPROXY=https://goproxy.cn
RUN go get -d -v ./...
RUN go install -v ./...

CMD ["ginchat"]