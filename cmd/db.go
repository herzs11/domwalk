/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"domwalk/db"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	syncCfg syncConfig
)

// dbCmd represents the db command
var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Pushing and pulling data to and from BigQuery",
	Long: `This command is used to push and pull data to and from BigQuery.

    The --push flag snapshots the current BQ dataset into the domwalk_snapshots dataset, then overwrites the current BQ dataset with the data from the local SQLite database`,
	Example: "domwalk db --push --gorm-db domwalk.db --bq-dataset domwalk",
	PreRun: func(cmd *cobra.Command, args []string) {

		bq_dataset, err := cmd.Flags().GetString("bq-dataset")
		if err != nil || bq_dataset == "" {
			color.Set(color.FgRed)
			cmd.Println("BQ dataset is required (use --bq-dataset or set $GORM_BQ_DATASET)")
			color.Unset()
			os.Exit(1)
		}
		err = db.CreateBigQueryConn()
		if err != nil {
			color.Set(color.FgRed)
			cmd.Println("Error getting bigquery client: ", err)
			color.Unset()
			os.Exit(1)
		}
		d, _ := cmd.Flags().GetBool("domains")
		c, _ := cmd.Flags().GetBool("cert-sans")
		w, _ := cmd.Flags().GetBool("web-redirects")
		dns, _ := cmd.Flags().GetBool("dns")
		s, _ := cmd.Flags().GetBool("sitemaps")
		if !d && !c && !w && !dns && !s {
			d = true
			c = true
			w = true
			dns = true
			s = true
		}
		syncCfg = syncConfig{
			Dataset:            bq_dataset,
			Domains:            d,
			CertSansDomains:    c,
			DNS:                dns,
			WebRedirectDomains: w,
			Sitemaps:           s,
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		backup, _ := cmd.Flags().GetBool("backup")
		if backup {
			color.Yellow("Backing up domains")
			err := backupFile(db.GormDB.DBName)
			if err != nil {
				color.Red("Error backing up domains: ", err)
				os.Exit(1)
			}
		}
		snapshot, _ := cmd.Flags().GetBool("snapshot")
		if snapshot {
			color.Yellow("Snapshotting domains\n")
			err := snapshotDomains()
			if err != nil {
				color.Red("Error snapshotting domains: ", err)
				os.Exit(1)
			}
		}
		if push, _ := cmd.Flags().GetBool("push"); push {
			color.Set(color.FgGreen)
			cmd.Println("Snapshotting domains")
			err := snapshotDomains()
			if err != nil {
				color.Red("Error snapshotting domains: ", err)
				os.Exit(1)
			}
			color.Unset()
			pushToBQ(syncCfg)
			return
		}
		if pull, _ := cmd.Flags().GetBool("pull"); pull {
			color.Yellow("Backing up domains")
			err := backupFile(db.GormDB.DBName)
			if err != nil {
				color.Red("Error backing up domains: ", err)
				os.Exit(1)
			}
			pullFromBQ(syncCfg)
			return
		}
	},
}

func init() {

	dbCmd.Flags().String(
		"bq-dataset", "domwalk",
		"BQ dataset to sync to",
	)

	dbCmd.Flags().Bool("push", false, "Push data to BigQuery (this takes a while)")
	dbCmd.Flags().Bool("pull", false, "Pull data from BigQuery (this takes a while)")
	dbCmd.MarkFlagsMutuallyExclusive("push", "pull")

	dbCmd.Flags().Bool("domains", false, "Sync domains")
	dbCmd.Flags().Bool("cert-sans", false, "Sync cert sans domains")
	dbCmd.Flags().Bool("web-redirects", false, "Sync web redirect domains")
	dbCmd.Flags().Bool("dns", false, "Sync DNS data")
	dbCmd.Flags().Bool("sitemaps", false, "Sync sitemaps")
	dbCmd.Flags().Bool("snapshot", false, "Snapshot domains (automatically done before push)")
	dbCmd.Flags().Bool("backup", false, "Backup local db (automatically done before pull)")
	dbCmd.MarkFlagsOneRequired("push", "pull", "snapshot", "backup")
	rootCmd.AddCommand(dbCmd)
}
