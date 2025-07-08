package debug

import (
	"context"
	"fmt"
	debug2 "github.com/siyibai/remote_debug/internal/debug"

	"github.com/spf13/cobra"
)

var DebugCmd = &cobra.Command{
	Use:   "debug",
	Short: "开始调试",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("debug: %s\n", args)
		debug := debug2.NewDebugImpl()
		debug.Debug(context.Background())
	},
}
