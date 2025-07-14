package logic

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/siyibai/remote_debug/agent/model"
)

type Impl struct {
}

func New() *Impl {
	return &Impl{}
}

// HealthHandler 健康检查
func (a *Impl) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	if err != nil {
		return
	}
}

// CommandHandler 命令执行处理
func (a *Impl) CommandHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req model.CommandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 异步执行命令
	go func(req model.CommandRequest) {
		_, _ = a.executeCommand(req)
	}(req)

	// 立即返回响应
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "command started",
	})
}

// executeCommand 执行命令
func (a *Impl) executeCommand(req model.CommandRequest) ([]string, error) {
	var ret []string
	for _, command := range req.PipeLine {
		cmd := exec.Command(command.Root, command.Args...)

		// 设置工作目录
		if command.WorkDir != "" {
			cmd.Dir = command.WorkDir
		}

		// 设置环境变量
		if len(command.Env) > 0 {
			cmd.Env = append(os.Environ(), command.Env...)
		}

		// 执行命令并获取输出
		output, err := cmd.CombinedOutput()
		if err != nil {
			return ret, fmt.Errorf("command execution failed: %w", err)
		}
		ret = append(ret, string(output))
	}
	return ret, nil
}

// UploadHandler 文件上传处理
func (a *Impl) UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 解析 multipart 表单
	if err := r.ParseMultipartForm(32 << 20); err != nil { // 32MB max
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 获取目标路径，如果没有提供则使用默认路径
	destPath := r.FormValue("dest_path")
	if destPath == "" {
		destPath = "./bin" // 默认路径
	}

	// 创建上传目录
	if err := os.MkdirAll(destPath, 0755); err != nil {
		http.Error(w, "Failed to create upload directory", http.StatusInternalServerError)
		return
	}

	// 保存文件
	filePath := filepath.Join(destPath, header.Filename)
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Failed to create file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	// 返回成功响应
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("File %s uploaded successfully", header.Filename),
		"path":    filePath,
	})
	if err != nil {
		return
	}
}
