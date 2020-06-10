package main

// "fmt"

func main() {
	DbConnect()
	defer Db.Close()

	scraper_name := "Blinkist"
	scrapers, err := GetScrapers()

	for _, elem := range scrapers {
		if elem.Name == scraper_name {
			isLocal := false
			response, results := Scrape(elem.Name, elem.Version, isLocal)
			if len(results) > 0 {
				scraping, err := elem.StartScrapingSession()
				file_path := GenerateFilePath(elem.Name, scraping.Id)
				SaveResults(elem, scraping, results)
				SaveResponseToStorage(response, file_path)
				_ = err
			}
		}
	}
	_ = err
}
