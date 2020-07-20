package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	. "github.com/logrusorgru/aurora"
)

var matchCompany string

var matchCmd = &cobra.Command{
	Use:   "match",
	Short: "Run the matcher",
	Long:  `Run the matcher for specific targets or for all of them.`,
	Run: func(cmd *cobra.Command, args []string) {
		TempMatch(matchCompany)
	},
}

func init() {
	rootCmd.AddCommand(matchCmd)
	matchCmd.Flags().StringVarP(&matchCompany, "match", "m", "", "Specify a company or all of them")
}

func TempMatch(company string) {
	DbConnect()
	defer Db.Close()

	if company == "all" {
		scrapers := GetScrapers()
		for _, elem := range scrapers {
			RunMatcher(elem)
		}
	} else {
		scraper := GetScraper(company)
		RunMatcher(scraper)
	}
}

func RunMatcher(scraper Scraper) {
	fmt.Println(Blue("Running --> "), Bold(Blue(scraper.Name)))

	matching := Matching{}
	matching.StartMatchingSession(scraper.Id)
	matches := GetMatches(matching, scraper.Id)

	for _, elem := range matches {
		fmt.Println(
			Bold(Green("\tNew Match -->")),
			Faint(Green(elem.KeywordText)),
			Bold(Green(elem.JobTitle)),
			Faint(Underline(BrightGreen(elem.JobUrl))))
	}

	if len(matches) > 0 {
		SaveMatches(matching, matches)
	}
}
