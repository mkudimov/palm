/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	tokenhandler "github.com/mkudimov/palm/token_handler"
	"github.com/mkudimov/palm/token_handler/readers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:   "read",
	Short: "A command for reading and handling tokens from a file",
	Long: `Read is a CLI command for reading and handling tokens from a file he required number of tokens. 
	After handling non-unique tokens will be dumped to a file. Example of usage:

read -i '~/home/user/desktop/input-sample.txt' -o '~/home/user/desktop/output-sample.txt'.`,
	Run: func(cmd *cobra.Command, args []string) {
		in, _ := cmd.Flags().GetString("in")
		out, _ := cmd.Flags().GetString("out")
		read(in, out)
	},
}

func init() {
	readCmd.Flags().StringP("in", "i", "", "Use to pass input file path with tokens for handling.")
	readCmd.Flags().StringP("out", "o", "", "Use to pass output file path for non-unique tokens.")
	readCmd.MarkFlagRequired("in")
	readCmd.MarkFlagRequired("out")
	rootCmd.AddCommand(readCmd)
}

func read(in, out string) error {
	reader := readers.NewFastReader()
	if connection := viper.GetString("postgresql"); connection != "" {
		hook, err := readers.NewPostgreSQLHook(connection)
		if err != nil {
			panic(err)
		}
		reader.AppendHook(hook)
	}
	handler := tokenhandler.NewHandler(reader)
	return handler.Handle(in, out)
}
