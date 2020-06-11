package main

func mainLocal() {
	DbConnect()
	scraper_name := "Shopify"
	scraper_version := 1
	scraping, err := LastScrapingByNameVersion(scraper_name, scraper_version)
	file_path := GenerateFilePath(scraper_name, scraping)
	fileResponse := GetResponseFromStorage(file_path)
	SaveResponseToFile(fileResponse)
	response, results := Scrape(scraper_name, scraper_version, true)
	RemoveFile()
	_ = response
	_ = results
	_ = err
}
