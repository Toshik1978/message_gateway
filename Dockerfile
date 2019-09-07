## Builder
FROM golang:1.13.0 AS build

ARG CI_GOPATH
ARG CI_GOPROXY

ENV GOPATH ${CI_GOPATH:-/.go/}
ENV GOPROXY ${CI_GOPROXY:-https://proxy.golang.org/}

# Go modules
WORKDIR /app
COPY go.mod /app
COPY go.sum /app
COPY Makefile /app
RUN make modules

# Build app
COPY . /app
RUN make build

## Destination image
FROM alpine:latest

RUN apk --no-cache add ca-certificates supervisor tzdata && \
	rm -rf /var/cache/apk/*

COPY --from=build /app/message_gateway /opt/message_gateway/message_gateway
COPY --from=build /app/configs/message_gateway.conf.yaml /opt/message_gateway/configs/message_gateway.conf.yaml
COPY .docker/supervisord.conf /etc/supervisord.conf

VOLUME ["/var/log/supervisor", "/opt/message_gateway/configs"]
WORKDIR /opt/message_gateway/
EXPOSE 8080

CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor.conf"]
