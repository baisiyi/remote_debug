FROM golang:1.22

WORKDIR /app

# 安装 dlv
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN cd ./agent/main go build -o

CMD ["tail", "-f", "/dev/null"]