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
	tokengenerator "github.com/mkudimov/palm/token_generator"
	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "A command for generating tokens and dumping them to a file.",
	Long: `Generate is a CLI command for generating he required number of tokens. 
	After generating tokens they will be dumped to a file. Example of usage:

generate -n 10000 -l 7 -o '~/home/user/desktop/output-sample.txt'.`,
	Run: func(cmd *cobra.Command, args []string) {
		tokenLength, _ := cmd.Flags().GetInt("length")
		tokenNumber, _ := cmd.Flags().GetInt64("number")
		out, _ := cmd.Flags().GetString("out")
		generate(tokenLength, tokenNumber, out)
	},
}

func init() {
	generateCmd.Flags().IntP("length", "l", 7, "Use to pass the length of tokens.")
	generateCmd.Flags().Int64P("number", "n", 1000, "Use to pass the number of tokens required.")
	generateCmd.Flags().StringP("out", "o", "", "Use to pass output file path for generated tokens.")
	generateCmd.MarkFlagRequired("out")
	rootCmd.AddCommand(generateCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func generate(length int, number int64, out string) error {
	return tokengenerator.Generate(length, number, out)
}
