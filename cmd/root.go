/*
Copyright © 2024 Seth Herz sherz@unum.com
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/herzs11/go-doms/domain"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

var ENRICH_DOMAIN_CF_URL string

var (
	domainsToExecute []string
	processConfig    ProcessConfig
	rParams          RequestParams
	token            *oauth2.Token
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "domwalk",
	Short: "CLI tool to find and store domain relationships",
	Long: `domwalk is a CLI tool to find and store domain relationships.
	It is written in Go and acts as a client for a domain enrichment Cloud Function.

	Currently, the tool can enrich domains with the following relationships:
	- Certificate Subject Alternative Names (SANs)
	- Web Redirects
	- SitemapLoc Web Domains
	- SitemapLoc Contact Page Domains

	The tool can also enrich domains with DNS data. In a future version, this dns data will be used to form additional domain relationships
	`,
	Example: `domwalk domains -d unum.com,coloniallife.com --workers 20 --cert-sans --web-redirects --sitemaps --dns`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		ENRICH_DOMAIN_CF_URL = os.Getenv("ENRICH_DOMAIN_CF_URL")
		if ENRICH_DOMAIN_CF_URL == "" {
			color.Red("ENRICH_DOMAIN_CF_URL environment variable must be set\n")
			os.Exit(1)
		}
		var err error
		token, err = getCredentials()
		if err != nil {
			color.Red(err.Error())
		}
		domainsToExecute = args
		cs, _ := cmd.Flags().GetBool("cert-sans")
		wr, _ := cmd.Flags().GetBool("web-redirects")
		sm, _ := cmd.Flags().GetBool("sitemaps")
		dns, _ := cmd.Flags().GetBool("dns")
		workers, _ := cmd.Flags().GetInt("workers")
		if workers < 1 {
			color.Red("Workers must be greater than 0\n")
			os.Exit(1)
		}
		minDate, _ := cmd.Flags().GetString("min-freshness")
		staleDate, err := time.Parse(time.DateOnly, minDate)
		if err != nil {
			color.Red("Invalid date format for min-freshness: (YYYY-MM-DD)\n")
			os.Exit(1)
		}
		if !cs && !wr && !sm && !dns {
			cs = true
			wr = true
			sm = true
			dns = true
		}
		processConfig = ProcessConfig{
			Workers: workers,
			EnrichmentConfig: domain.EnrichmentConfig{
				CertSans:         cs,
				DNS:              dns,
				Sitemap:          sm,
				WebRedirect:      wr,
				MinFreshnessDate: staleDate,
			},
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		noReturn, _ := cmd.Flags().GetBool("no-return")
		onlyMatched, _ := cmd.Flags().GetBool("only-matched")
		rParams = RequestParams{
			DomainNames:        domainsToExecute,
			ProcessConfig:      processConfig,
			NoResponse:         noReturn,
			OnlyMatchedDomains: onlyMatched,
		}
		reqBody, err := json.Marshal(rParams)
		if err != nil {
			color.Red("Error creating request parameters: %s\n", err.Error())
			os.Exit(1)
		}
		req, err := http.NewRequest("POST", ENRICH_DOMAIN_CF_URL, bytes.NewReader(reqBody))
		if err != nil {
			color.Red("Error creating request: %s\n", err.Error())
			os.Exit(1)
		}
		if token.AccessToken != "" {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
		}
		// Create client with no timeout
		// This is a CLI tool and we want to wait for the response
		// If the request hangs, the user can kill the process
		client := http.DefaultClient
		client.Timeout = 0
		proxyUrl, err := http.ProxyFromEnvironment(req)
		if err != nil {
			fmt.Println("Error getting proxy URL: ", err)
		}
		if proxyUrl != nil {
			transport := &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			}
			client.Transport = transport
		}
		resp, err := client.Do(req)
		if err != nil {
			color.Red("Error sending request: %s\n", err.Error())
			os.Exit(1)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			color.Red("Error processing request: %s\n", resp.Status)
			os.Exit(1)
		}
		if noReturn {
			color.Green("Enriched domains\n")
			return
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			color.Red("Error reading response: %s\n", err.Error())
			os.Exit(1)
		}
		var data bytes.Buffer
		err = json.Indent(&data, body, "", "  ")
		if err != nil {
			color.Red("Error formatting response: %s\n", err.Error())
			os.Exit(1)
		}
		outputFile, _ := cmd.Flags().GetString("output")
		if outputFile != "" {
			err = os.WriteFile(outputFile, data.Bytes(), 0644)
			if err != nil {
				color.Red("Error writing output file: %s\n", err.Error())
				os.Exit(1)
			}
			color.Green("Results written to %s\n", outputFile)
			return
		}
		color.Magenta("Results:\n%s\n", data.String())
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String(
		"min-freshness", "0001-01-01", "Minimum date to refresh relationships, (YYYY-MM-DD)",
	)
	rootCmd.PersistentFlags().Bool("cert-sans", false, "Enrich domains with cert SANs")
	rootCmd.PersistentFlags().Bool("web-redirects", false, "Enrich domains with web redirects")
	rootCmd.PersistentFlags().Bool("sitemaps", false, "Enrich domains with sitemap web domains")
	rootCmd.PersistentFlags().Bool("dns", false, "Enrich domains with dns data")
	rootCmd.PersistentFlags().IntP("workers", "w", 15, "Number of concurrent workers to use")
	rootCmd.PersistentFlags().BoolP("no-return", "q", false, "Do not return results")
	rootCmd.PersistentFlags().BoolP("only-matched", "m", false, "Only return matched domains")
	rootCmd.PersistentFlags().StringP(
		"output", "o", "", "Output JSON file for results, cannot be used with --no-return",
	)
	rootCmd.MarkFlagsMutuallyExclusive("no-return", "output")
	rootCmd.MarkFlagsMutuallyExclusive("no-return", "only-matched")

}
