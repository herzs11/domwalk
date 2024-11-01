/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// fileCmd represents the file command
var fileCmd = &cobra.Command{
	Use:   "file",
	Short: "Enrich domains from file",
	Long: `Provide a file with a list of domain names to enrich. The file should be a .txt or .csv file with one domain per line.
`,
	Example: `domwalk file -f domains.txt`,
	PreRun: func(cmd *cobra.Command, args []string) {
		var err error
		file, _ := cmd.Flags().GetString("file")
		if file == "" {
			color.Red("File is required\n")
			os.Exit(1)
		}
		domainsToExecute, err = readDomainsFromFile(file)
		if err != nil {
			color.Red(err.Error())
			os.Exit(1)
		}
		noHeader, _ := cmd.Flags().GetBool("no-header")
		if !noHeader {
			domainsToExecute = domainsToExecute[1:]
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		color.Yellow("Enriching %d domains from file...", len(domainsToExecute))
	},
}

func init() {
	rootCmd.AddCommand(fileCmd)

	fileCmd.Flags().StringP("file", "f", "", "File with domains to enrich")
	fileCmd.Flags().Bool("no-header", false, "The file does not have a header row")
	fileCmd.MarkFlagRequired("f")
}

func readDomainsFromFile(filename string) ([]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}

	// remember to close the file at the end of the program
	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading csv data: %w", err)
	}
	doms := make([]string, 0, len(data))
	for _, row := range data {
		doms = append(doms, row[0])
	}
	return doms, nil
}
