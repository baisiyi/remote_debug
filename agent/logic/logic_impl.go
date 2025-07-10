package logic

import (
	"encoding/json"
	"fmt"
	"github.com/siyibai/remote_debug/agent/model"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
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

	// 执行命令
	response := model.CommandResponse{}
	output, err := a.executeCommand(req)
	if err != nil {
		response.Success = false
		response.Error = err.Error()
	} else {
		response.Success = true
		response.Output = output
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}

// executeCommand 执行命令
func (a *Impl) executeCommand(req model.CommandRequest) (string, error) {
	cmd := exec.Command(req.Command, req.Args...)

	// 设置工作目录
	if req.WorkDir != "" {
		cmd.Dir = req.WorkDir
	}

	// 设置环境变量
	if len(req.Env) > 0 {
		cmd.Env = append(os.Environ(), req.Env...)
	}

	// 执行命令并获取输出
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("command execution failed: %w", err)
	}

	return string(output), nil
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

	// 创建上传目录
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		http.Error(w, "Failed to create upload directory", http.StatusInternalServerError)
		return
	}

	// 保存文件
	filePath := filepath.Join(uploadDir, header.Filename)
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
