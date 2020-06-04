package main

// import "fmt"

func main() {
	DbConnect()

	scraper_name := "Google"

	// Select all new results based on scraper name
	scraper, err := GetScraperByScraperName(scraper_name)
	if err != nil {
		panic(err.Error())
	}

	scraping, err := GetLastScrapingByScraperId(scraper)
	if err != nil {
		panic(err.Error())
	}

	results, err := GetNewResultsByScrapingId(scraping)
	if err != nil {
		panic(err.Error())
	}
	_ = results

	/**
	  keywords, err := GetKeywordsByScraperId(scraper)
	  if err != nil {
	      panic(err.Error())
	  }
	*/

	// Select all the keywords based on scraper name

	// Match
	defer Db.Close()
}

// Voglio matchare i risulati dell ultima onda di google con le kws degli utenti che hanno google
