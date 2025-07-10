package model

// CommandRequest 命令请求结构
type CommandRequest struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
	WorkDir string   `json:"work_dir,omitempty"`
	Env     []string `json:"env,omitempty"`
}

// CommandResponse 命令响应结构
type CommandResponse struct {
	Success bool   `json:"success"`
	Output  string `json:"output,omitempty"`
	Error   string `json:"error,omitempty"`
}
