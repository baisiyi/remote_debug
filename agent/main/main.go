package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/siyibai/remote_debug/agent"
)

func main() {
	// 从环境变量或命令行参数获取端口
	port := os.Getenv("AGENT_PORT")
	if port == "" {
		port = "8080" // 默认端口
	}

	// 创建并启动 agent server
	server := agent.NewAgentServer(port)

	// 优雅关闭
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down agent server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5)
		defer cancel()
		if err := server.Stop(ctx); err != nil {
			log.Printf("Error stopping server: %v", err)
		}
	}()

	// 启动服务器
	if err := server.Start(); err != nil && !errors.Is(err, context.Canceled) {
		log.Fatalf("Failed to start agent server: %v", err)
	}
}
