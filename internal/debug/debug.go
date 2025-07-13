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
	args   []string
}

func NewDebugImpl(args []string) *DebugImpl {
	return &DebugImpl{
		serApi: client.NewSerApi(),
		args:   args,
	}
}

func (d *DebugImpl) Debug(ctx context.Context) {

	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	// 编译项目
	fmt.Println("开始编译...")
	err = d.buildProject(ctx, cfg.ProjectPath)
	if err != nil {
		return
	}
	fmt.Println("编译完成")
	// 上传项目可执行文件
	fmt.Println("开始同步可执行文件...")
	err = d.syncProject(ctx, cfg.ProjectPath)
	if err != nil {
		return
	}
	fmt.Println("同步完成")
	// 运行可执行文件
	d.runProject(ctx)
}

func (d *DebugImpl) buildProject(ctx context.Context, projectPath string) (err error) {
	// cd
	if err = d.runCommand(ctx, "cd", projectPath); err != nil {
		return err
	}

	// 编译
	err = d.runCommand(ctx, "go", "build", "-gcflags", "all=-N -l", "-o", ObjectName)
	if err != nil {
		return
	}
	return
}

func (d *DebugImpl) syncProject(ctx context.Context, projectPath string) (err error) {
	err = d.serApi.UploadFile(ctx, fmt.Sprintf("%s/%s", projectPath, ObjectName))
	if err != nil {
		fmt.Printf("可执行文件同步失败：%v", err)
		return
	}
	return
}

func (d *DebugImpl) runProject(ctx context.Context) {
	runCmd := model.Command{
		Root:    "dlv",
		Args:    []string{"exec", "./myapp", "--port=8080", "--", "-config", "/config.yaml"},
		WorkDir: "/app",
		Env:     []string{},
	}
	err := d.serApi.RunCommand(ctx, &model.CommandRequest{
		PipeLine: []model.Command{runCmd},
	})
	if err != nil {
		fmt.Println(err)
	}
}

func (d *DebugImpl) runCommand(ctx context.Context, name string, args ...string) (err error) {
	buildCmd := exec.Command(name, args...)
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	if err = buildCmd.Run(); err != nil {
		fmt.Printf("编译失败：err:%v", err)
		return
	}
	return
}
