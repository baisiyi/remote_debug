package debug

import (
	"context"
	"fmt"
	"github.com/siyibai/remote_debug/internal/client"
	"github.com/siyibai/remote_debug/internal/model"
	"os"
	"os/exec"

	"github.com/siyibai/remote_debug/config"
)

const (
	ObjectName   = "debug"
	ObjectInPath = "/app"
)

type DebugImpl struct {
	serApi *client.SerApi
}

func NewDebugImpl() *DebugImpl {
	return &DebugImpl{
		serApi: client.NewSerApi(),
	}
}

func (d *DebugImpl) Debug(ctx context.Context) {

	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	// 编译项目
	d.buildProject(ctx, cfg.ProjectPath)
	// 上传项目可执行文件
	d.syncProject(ctx)
	// 运行可执行文件
	d.runProject(ctx)
}

func (d *DebugImpl) buildProject(ctx context.Context, projectPath string) {
	buildCmd := exec.Command("go", "build", "-gcflags", "all=-N -l", "-o", ObjectName, projectPath)
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	if err := buildCmd.Run(); err != nil {
		fmt.Println(err)
		return
	}
}

func (d *DebugImpl) syncProject(ctx context.Context) {
	err := d.serApi.UploadFile(ctx, fmt.Sprintf("%s/%s", ObjectInPath, ObjectName))
	if err != nil {
		fmt.Println(err)
	}
}

func (d *DebugImpl) runProject(ctx context.Context) {
	runCmd := model.Command{
		Main:    "dlv",
		Args:    []string{"exec", "./myapp", "--port=8080", "--", "-config", "/config.yaml"},
		WorkDir: "/app",
		Env:     []string{},
	}
	err := d.serApi.RunCommand(ctx, &model.RunCommandReq{
		PipeLine: []model.Command{runCmd},
	})
	if err != nil {
		fmt.Println(err)
	}
}
