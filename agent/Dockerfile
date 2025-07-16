# docker build -t remote_debug:latest -f agent/Dockerfile .
FROM golang:1.24 AS builder

WORKDIR /app

COPY agent/go.mod agent/go.sum ./
RUN go mod download

COPY agent/ ./
WORKDIR /app/main
RUN go build -o /app/agent

FROM golang:1.24

WORKDIR /app

COPY --from=builder /app/agent /app/agent

RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN apt-get update && apt-get install -y lsof

EXPOSE 8081

ENTRYPOINT ["/app/agent"]