// Copyright 2023 Quokka Contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

			if err := os.WriteFile("config.toml", tomlBytes, 0600); err != nil {
				return err
			}

			return nil
		},
	}

	configCmd.AddCommand(createCmd)

	return configCmd
}
