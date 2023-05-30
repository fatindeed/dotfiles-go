package cmd

import (
	"github.com/fatindeed/dotfiles-go/app"
	"github.com/spf13/cobra"
)

func init() {
	initFlags(recoverCmd)
	rootCmd.AddCommand(recoverCmd)
}

var recoverCmd = &cobra.Command{
	Use:   "recover",
	Short: "Recover dotfiles",
	RunE: func(cmd *cobra.Command, args []string) error {
		handler := &app.RecoverHandler{}
		return handler.Run(args)
	},
}
