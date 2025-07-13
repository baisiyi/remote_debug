package debug

import (
	"context"
	debug2 "github.com/siyibai/remote_debug/internal/debug"

	"github.com/spf13/cobra"
)

var DebugCmd = &cobra.Command{
	Use:   "debug",
	Short: "开始调试",
	Run: func(cmd *cobra.Command, args []string) {
		debug := debug2.NewDebugImpl(args)
		debug.Debug(context.Background())
	},
}
