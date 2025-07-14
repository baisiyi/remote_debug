package debug

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/siyibai/remote_debug/internal/client"
	"github.com/siyibai/remote_debug/internal/model"

	"github.com/siyibai/remote_debug/config"
)

const (
	ObjectName = "debug"
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
	fmt.Println("启动dlv调试...")
	d.runProject(ctx)
	fmt.Println("启动成功，开始调试")
}

func (d *DebugImpl) buildProject(ctx context.Context, projectPath string) (err error) {
	// 编译
	err = d.runCommandWithDir(ctx, projectPath, "go", "build", "-gcflags", "all=-N -l", "-o", ObjectName)
	if err != nil {
		return
	}
	return
}

func (d *DebugImpl) syncProject(ctx context.Context, projectPath string) (err error) {
	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Printf("获取配置失败：%v", err)
		return
	}

	err = d.serApi.UploadFile(ctx, fmt.Sprintf("%s/%s", projectPath, ObjectName), cfg.RemoteAddress.DestPath)
	if err != nil {
		fmt.Printf("可执行文件同步失败：%v", err)
		return
	}
	return
}

func (d *DebugImpl) runProject(ctx context.Context) {
	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Printf("获取配置失败：%v", err)
		return
	}

	runCmd := model.Command{
		Root:    "dlv",
		Args:    []string{"exec", ObjectName, "--port=8080", "--", "-config", "/config.yaml"},
		WorkDir: cfg.RemoteAddress.DestPath,
		Env:     []string{},
	}
	rsp, err := d.serApi.RunCommand(ctx, &model.CommandRequest{
		PipeLine: []model.Command{runCmd},
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(rsp)
}

func (d *DebugImpl) runCommandWithDir(ctx context.Context, dir, name string, args ...string) (err error) {
	buildCmd := exec.Command(name, args...)
	buildCmd.Dir = dir
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	if err = buildCmd.Run(); err != nil {
		fmt.Printf("编译失败：err:%v", err)
		return
	}
	return
}
