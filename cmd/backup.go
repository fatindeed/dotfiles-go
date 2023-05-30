package cmd

import (
	"github.com/fatindeed/dotfiles-go/app"
	"github.com/spf13/cobra"
)

func init() {
	initFlags(backupCmd)
	rootCmd.AddCommand(backupCmd)
}

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup dotfiles",
	RunE: func(cmd *cobra.Command, args []string) error {
		handler := &app.BackupHandler{}
		return handler.Run(args)
	},
}
