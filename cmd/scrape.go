package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	. "github.com/logrusorgru/aurora"
)

var (
	scrapeCompany     string
	runSavers         string
	excludedCompanies []string
)

var scrapeCmd = &cobra.Command{
	Use:   "scrape",
	Short: "Run the scraper",
	Long:  `Run the scraper for specific targets or for all of them.`,
	Run: func(cmd *cobra.Command, args []string) {
		Scrape(scrapeCompany, runSavers, excludedCompanies)
	},
}

func init() {
	rootCmd.AddCommand(scrapeCmd)
	scrapeCmd.Flags().StringVarP(&scrapeCompany, "scrape", "s", "", "Specify a company or all of them")
	scrapeCmd.Flags().StringVarP(&runSavers, "results", "r", "", "Specify to save results or not")
	scrapeCmd.Flags().StringSliceVarP(&excludedCompanies, "excluded", "e", []string{}, "Specify which companies to exclude from scraping")
}

func Scrape(company string, runSavers string, excludedCompanies []string) {
	start := time.Now()
	DbConnect()
	defer Db.Close()
	if company == "all" {
		scrapers := GetScrapers()
		for _, elem := range scrapers {
			if !Contains(excludedCompanies, elem.Name) {
				RunScraper(elem, runSavers)
			}
		}
	} else {
		scraper := GetScraper(company)
		RunScraper(scraper, runSavers)
	}
	elapsed := time.Since(start)
	fmt.Println(BrightMagenta("Total Execution Time -->"), Bold(BrightMagenta(elapsed)))
}

func RunScraper(scraper Scraper, runSavers string) {
	fmt.Println(BrightBlue("Scraping -->"), Bold(BrightBlue(scraper.Name)))
	results := Extract(scraper.Name)
	if runSavers == "true" {
		count_results := len(results)
		if count_results > 0 {
			fmt.Println(Green("Number of results scraped: "), Bold(Green(count_results)))
			scraping := scraper.StartScrapingSession(count_results)
			SaveResults(scraper, scraping, results)
		} else {
			fmt.Println(Bold(Red("DANGER, NO RESULTS FOUND")))
		}
	}
}
