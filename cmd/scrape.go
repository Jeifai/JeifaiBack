package cmd

import (
    "github.com/spf13/cobra"
    scraper "github.com/Jeifai/JeifaiBack/cmd/scraper"

)

var company string

var scrapeCmd = &cobra.Command{
	Use:   "scrape",
	Short: "Run the scraper",
    Long: `Run the scraper for specific targets or for all of them.`,
	Run: func(cmd *cobra.Command, args []string) {
        scraper.Scrape(company)
	},
}

func init() {
    rootCmd.AddCommand(scrapeCmd)
	scrapeCmd.Flags().StringVarP(&company, "scrapecompany", "s", "", "Specify a company or all of them")
}
