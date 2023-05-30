package cmd

import (
	"github.com/fatindeed/dotfiles-go/app"
	"github.com/spf13/cobra"
)

func init() {
	initFlags(viewCmd)
	rootCmd.AddCommand(viewCmd)
}

var viewCmd = &cobra.Command{
	Use:   "view [flags] filename (- for stdin)",
	Short: "View encrypted contents",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		handler := &app.ViewHandler{}
		return handler.Run(args)
	},
}
