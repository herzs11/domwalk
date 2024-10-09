/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// docsCmd represents the docs command
var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		err := doc.GenMarkdownTree(rootCmd, "/Users/sherz/GolandProjects/domwalk/docs/")
		if err != nil {
			panic(err)
		}
	},
	Hidden: true,
}

func init() {
	rootCmd.AddCommand(docsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// docsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// docsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
