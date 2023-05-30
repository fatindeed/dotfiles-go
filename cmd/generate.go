package cmd

import (
	"github.com/fatindeed/dotfiles-go/app"
	"github.com/spf13/cobra"
)

var encryptConfig = new(app.Encrypt)

func init() {
	initFlags(generateCommand)
	generateCommand.Flags().StringVar(&encryptConfig.Template, "key-template", "AES256_GCM", "key type")
	generateCommand.Flags().StringVar(&encryptConfig.KeyPath, "key-path", "keyset.json", "key path")
	generateCommand.Flags().StringVar(&encryptConfig.MasterKeyURI, "master-key-uri", "", "master key uri")
	rootCmd.AddCommand(generateCommand)
}

var generateCommand = &cobra.Command{
	Use:   "generate",
	Short: "Generate new config.yaml and keyset",
	RunE: func(cmd *cobra.Command, args []string) error {
		handler := &app.GenerateHandler{Encrypt: encryptConfig}
		return handler.Run(args)
	},
}
