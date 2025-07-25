package model

// CommandRequest 命令请求结构
type CommandRequest struct {
	PipeLine []Command
}

type Command struct {
	Root    string   `json:"root"`               // 主命令，如 "ls"、"go"、"dlv"
	Args    []string `json:"args"`               // 参数数组，如 ["-l", "/tmp"]
	WorkDir string   `json:"work_dir,omitempty"` // 可选，指定工作目录
	Env     []string `json:"env,omitempty"`      // 可选，环境变量
}

// CommandResponse 命令响应结构
type CommandResponse struct {
	Success bool     `json:"success"`
	Output  []string `json:"output,omitempty"`
	Error   string   `json:"error,omitempty"`
}

// RunTplParams 用于渲染 run.sh 模板的参数
// App: 可执行文件名
// DestPath: 远程目标路径
// ServerDebugPort: 调试端口
// RunCmdArgs: 运行参数
type RunTplParams struct {
	App             string
	DestPath        string
	ServerDebugPort string
	RunCmdArgs      string
}

/*
{
  "command": "dlv",
  "args": ["run", "main.go"],
  "work_dir": "/app",
  "env": ["GO111MODULE=on"]
}
import (
    "os/exec"
)

func RunCommand(req CommandRequest) ([]byte, error) {
    cmd := exec.Command(req.Command, req.Args...)
    if req.WorkDir != "" {
        cmd.Dir = req.WorkDir
    }
    if len(req.Env) > 0 {
        cmd.Env = append(cmd.Env, req.Env...)
    }
    return cmd.CombinedOutput()
}
*/
