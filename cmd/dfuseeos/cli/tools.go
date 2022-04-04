package cli

import "github.com/zhongshuwen/histnew/tools"

func init() {
	RootCmd.AddCommand(tools.Cmd)
}
