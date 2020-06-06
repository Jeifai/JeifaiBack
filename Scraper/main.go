package main

// "fmt"

func main() {
	DbConnect()
	scraper_name := "Shopify"
	scrapers, err := GetScrapers()
	if err != nil {
		panic(err.Error())
	}
	for _, elem := range scrapers {
		if elem.Name == scraper_name {
			scraping, err := elem.StartScrapingSession()
			if err != nil {
				panic(err.Error())
			}
			isLocal := false
			response, results := Scrape(elem.Name, elem.Version, isLocal)
			file_path := GenerateFilePath(elem.Name, scraping.Id)
			SaveResults(elem, scraping, results)
			SaveResponseToStorage(response, file_path)
		}
	}
	defer Db.Close()
}
