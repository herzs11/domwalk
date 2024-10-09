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
	Use: "docs",
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
}
