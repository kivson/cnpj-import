/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/kivson/cnpj-import/importer"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var path, dbType, dbDsn string

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import records from ZIP files into database",
	Long:  `Import records from ZIP files into database`,
	Run: func(cmd *cobra.Command, args []string) {
		bar := progressbar.Default(-1, "Processed Records")

		imp := importer.Importer{
			DbType:     importer.DBType(dbType),
			DbDsn:      dbDsn,
			ProgressFn: bar.Add,
		}
		imp.ImportZipFolder(path)
	},
}

func init() {
	rootCmd.AddCommand(importCmd)

	importCmd.Flags().StringVarP(&path, "path", "p", "./zips", "Path to zip folder with files to be imported")
	importCmd.Flags().StringVarP(&dbType, "dbtype", "t", "sqlite", "Database type (sqlite)")
	importCmd.Flags().StringVarP(&dbDsn, "dbdsn", "d", "./database.db?_journal_mode=MEMORY&_sync=OFF", "Database connection dsn")
}
