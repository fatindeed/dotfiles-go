package cmd

import (
	"os"

	"github.com/fatindeed/dotfiles-go/app"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile string
	verbose bool

	rootCmd = &cobra.Command{
		Use:          "dotfiles",
		Version:      "v0.0.0-dev",
		Short:        "dotfiles is a dotfiles management application",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}
)

func Execute(version string) {
	if version != "" {
		rootCmd.Version = version
	}
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initCobra)
}

func initCobra() {
	log.SetOutput(os.Stdout)
	if verbose {
		log.SetLevel(log.DebugLevel)
	}

	// Setting Overrides
	viper.Set("name", rootCmd.Use)
	viper.Set("version", rootCmd.Version)

	app.SetConfigFile(cfgFile)
}

func initFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&cfgFile, "config", "c", "", "config file")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	viper.BindPFlag("config", cmd.Flags().Lookup("config"))
	viper.SetDefault("config", "config.yaml")
}
