/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// domainsCmd represents the domains command
var domainsCmd = &cobra.Command{
	Use:   "domains",
	Short: "Enrich domains from a list of domain names",
	Long:  `Provide a list of domain names to enrich.`,
	Example: `domwalk domains -d "unum.com,coloniallife.com"
domwalk domains -d unum.com -d coloniallife.com`,
	PreRun: func(cmd *cobra.Command, args []string) {
		domainsToExecute, _ = cmd.Flags().GetStringSlice("domains")
	},
	Run: func(cmd *cobra.Command, args []string) {
		color.Yellow("Enriching %d domains...", len(domainsToExecute))
	},
}

func init() {
	rootCmd.AddCommand(domainsCmd)
	domainsCmd.Flags().StringSliceP("domains", "d", []string{}, "Domains to enrich")
	domainsCmd.MarkFlagRequired("domains")
}
