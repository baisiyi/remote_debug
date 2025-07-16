package debug

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/siyibai/remote_debug/internal/client"
	"github.com/siyibai/remote_debug/internal/model"

	"path/filepath"
	"text/template"

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
		fmt.Println(err)
		return
	}
	fmt.Println("编译完成")
	// 上传项目可执行文件
	fmt.Println("开始同步可执行文件...")
	err = d.syncProject(ctx, cfg.ProjectPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("同步完成")
	fmt.Println("准备启动调试...")
	err = d.buildShell(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("准备完毕")
	// 运行可执行文件
	fmt.Println("启动dlv调试...")
	d.runProject(ctx)
	fmt.Println("启动成功，开始调试")
}

func (d *DebugImpl) buildProject(ctx context.Context, projectPath string) (err error) {

	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Printf("获取配置失败：%v", err)
		return
	}
	buildCmd := exec.Command("sh", "-c", fmt.Sprintf(cfg.BuildCmdFmt, ObjectName))
	buildCmd.Dir = projectPath
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	// 设置交叉编译环境变量
	if cfg.CrossCompileCmd != "" {
		cross := strings.Split(cfg.CrossCompileCmd, " ")
		buildCmd.Env = append(os.Environ(), cross...)
	}
	if err = buildCmd.Run(); err != nil {
		fmt.Printf("编译失败：err:%v", err)
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

func (d *DebugImpl) buildShell(ctx context.Context) (err error) {
	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Printf("获取配置失败：%v", err)
		return
	}

	// 1. 渲染 run.sh 脚本
	tplParams := model.RunTplParams{
		App:             ObjectName,
		DestPath:        cfg.RemoteAddress.DestPath,
		ServerDebugPort: "2345", // 可根据需要从cfg获取
		RunCmdArgs:      cfg.RunCmdArgs,
	}
	tpl, err := template.New("run").Parse(model.RunTpl)
	if err != nil {
		fmt.Printf("run.sh 模板解析失败: %v", err)
		return
	}
	tmpDir := os.TempDir()
	runShPath := filepath.Join(tmpDir, "run.sh")
	runShFile, err := os.Create(runShPath)
	if err != nil {
		fmt.Printf("run.sh 文件创建失败: %v", err)
		return
	}
	defer runShFile.Close()
	if err = tpl.Execute(runShFile, tplParams); err != nil {
		fmt.Printf("run.sh 模板渲染失败: %v", err)
		return
	}
	// 确保脚本有可执行权限
	err = runShFile.Chmod(0755)
	if err != nil {
		return
	}

	// 2. 上传 run.sh 到远程
	err = d.serApi.UploadFile(ctx, runShPath, cfg.RemoteAddress.DestPath)
	if err != nil {
		fmt.Printf("run.sh 脚本上传失败：%v", err)
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
		Root:    "sh",
		Args:    []string{"-c", "./run.sh"},
		WorkDir: cfg.RemoteAddress.DestPath,
		Env:     []string{},
	}
	_, err = d.serApi.RunCommand(ctx, &model.CommandRequest{
		PipeLine: []model.Command{runCmd},
	})
	if err != nil {
		fmt.Println(err)
	}
}
