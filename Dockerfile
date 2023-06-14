FROM golang:alpine AS builder

WORKDIR /app
COPY . ./
ENV GOBIN=/app CGO_ENABLED=0 GOOS=linux
RUN apk update
RUN apk add --no-cache git curl tar
RUN go mod download
RUN go build -o ./forum
RUN go install github.com/DarthSim/overmind/v2

ENV RQLITE_VERSION=7.20.2
RUN curl -L https://github.com/rqlite/rqlite/releases/download/v${RQLITE_VERSION}/rqlite-v${RQLITE_VERSION}-linux-amd64-musl.tar.gz -o rqlite-v${RQLITE_VERSION}-linux-amd64-musl.tar.gz && \
    tar xvfz rqlite-v${RQLITE_VERSION}-linux-amd64-musl.tar.gz && \
    cp rqlite-v${RQLITE_VERSION}-linux-amd64-musl/rqlited /app && \
    cp rqlite-v${RQLITE_VERSION}-linux-amd64-musl/rqlite /app

FROM alpine:latest

RUN apk update
RUN apk add --no-cache tmux

EXPOSE 8080

WORKDIR /app
COPY --from=builder /app/forum ./
COPY --from=builder /app/overmind ./
COPY --from=builder /app/rqlited ./
COPY --from=builder /app/rqlite ./

ADD Procfile ./

ENTRYPOINT ["/app/overmind", "start", "--procfile", "/app/Procfile", "--auto-restart", "server,db"]
