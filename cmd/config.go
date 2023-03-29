package cmd

import (
	"os"

	"github.com/pelletier/go-toml/v2"
	"github.com/quokkamail/quokka/config"
	"github.com/spf13/cobra"
)

func NewConfigCmd() *cobra.Command {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration file",
	}

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create new configuration file",
		RunE: func(cmd *cobra.Command, args []string) error {
			tomlBytes, err := toml.Marshal(config.Default)
			if err != nil {
				return err
			}

			if err := os.WriteFile("config.toml", tomlBytes, 0666); err != nil {
				return err
			}

			return nil
		},
	}

	configCmd.AddCommand(createCmd)

	return configCmd
}
