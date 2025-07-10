package agent

import (
	"context"
	"github.com/siyibai/remote_debug/agent/logic"
	"log"
	"net/http"
)

// Server agent 服务器
type Server struct {
	server *http.Server
	port   string
	impl   *logic.Impl
}

// NewAgentServer 创建新的 agent 服务器
func NewAgentServer(port string) *Server {
	return &Server{
		impl: logic.New(),
		port: port,
	}
}

// Start 启动服务器
func (a *Server) Start() error {
	mux := http.NewServeMux()

	// 注册路由
	mux.HandleFunc("/health", a.impl.HealthHandler)
	mux.HandleFunc("/command", a.impl.CommandHandler)
	mux.HandleFunc("/upload", a.impl.UploadHandler)

	a.server = &http.Server{
		Addr:    ":" + a.port,
		Handler: mux,
	}

	log.Printf("Agent server starting on port %s", a.port)
	return a.server.ListenAndServe()
}

// Stop 停止服务器
func (a *Server) Stop(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}
