package commands

import (
	"github.com/spf13/cobra"
	"github.com/wisnuragaprawida/project/bootstrap"
	"github.com/wisnuragaprawida/project/pkg/log"
)

func init() {
	registerCommand(startTryService)
}

func startTryService(dep *bootstrap.Dependency) *cobra.Command {
	return &cobra.Command{
		Use:   "try",
		Short: "try service",
		Run: func(cmd *cobra.Command, args []string) {
			log.Info("try service")
		},
	}
}
