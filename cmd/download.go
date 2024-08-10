/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"net/http"

	"github.com/cavaliergopher/grab/v3"
	"github.com/kivson/cnpj-import/downloader"
	"github.com/spf13/cobra"
)

var dest, url string

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		downloader := downloader.Downloader{
			HttpClient: &http.Client{},
			Graber:     grab.NewClient(),
			Dest:       dest,
			BaseUrl:    url,
		}
		downloader.DownloadAll()
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().StringVarP(&url, "url", "u", "https://dadosabertos.rfb.gov.br/CNPJ/", "Url where the zips are listed")
	downloadCmd.Flags().StringVarP(&dest, "folder", "f", "./zips", "Folder destination")
}
