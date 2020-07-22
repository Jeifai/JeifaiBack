package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	. "github.com/logrusorgru/aurora"
)

var (
	scrapeCompany   string
    runLocally      string
	runSavers       string
)

var scrapeCmd = &cobra.Command{
	Use:   "scrape",
	Short: "Run the scraper",
	Long:  `Run the scraper for specific targets or for all of them.`,
	Run: func(cmd *cobra.Command, args []string) {
		Scrape(scrapeCompany, runLocally, runSavers)
	},
}

func init() {
	rootCmd.AddCommand(scrapeCmd)
	scrapeCmd.Flags().StringVarP(&scrapeCompany, "scrape", "s", "", "Specify a company or all of them")
    scrapeCmd.Flags().StringVarP(&runLocally, "islocal", "i", "", "Specify if the scraper will run locally or not")
	scrapeCmd.Flags().StringVarP(&runSavers, "results", "r", "", "Specify to save results or not")
}

func Scrape(company string, runLocally string, runSavers string) {
	DbConnect()
	defer Db.Close()
	if company == "all" {
		scrapers := GetScrapers()
		for _, elem := range scrapers {
			RunScraper(elem, runLocally, runSavers)
		}
	} else {
		scraper := GetScraper(company)
		RunScraper(scraper, runLocally, runSavers)
	}
}

func RunScraper(scraper Scraper, runLocally string, runSavers string) {
	if runLocally == "false" {
		fmt.Println(BrightBlue("Scraping -->"), Bold(BrightBlue(scraper.Name)))
        response, results := Extract(scraper.Name, scraper.Version, false)
	    if runSavers == "true" {
            n_results := len(results)
            if n_results > 0 {
                fmt.Println(Green("Number of results scraped: "), Bold(Green(n_results)))
                scraping := scraper.StartScrapingSession()
                file_path := GenerateFilePath(scraper.Name, scraping.Id)
                SaveResults(scraper, scraping, results)
                SaveResponseToStorage(response, file_path)
            } else {
                fmt.Println(Bold(Red("DANGER, NO RESULTS FOUND")))
            }
        }
	} else if runLocally == "true" {
		scraping := LastScrapingByNameVersion(scraper.Name, scraper.Version)
		file_path := GenerateFilePath(scraper.Name, scraping)
		fileResponse := GetResponseFromStorage(file_path)
		SaveResponseToFile(fileResponse)
		_, results := Extract(scraper.Name, scraper.Version, true)
		n_results := len(results)
		if n_results > 0 {
			fmt.Println(Green("Number of results scraped: "), Bold(Green(n_results)))
		}
		RemoveFile()
	}
}
