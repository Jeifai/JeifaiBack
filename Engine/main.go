package main

import (
	"fmt"
    "os"
)

func main() {
    DbConnect()
	// scrape("Babelforce")
	test("IMusician", 1)
}
func scrape(scraper_name string) {
	scrapers, err := Scrapers()
	if err != nil {
		return
	}
	for _, elem := range scrapers {
		if elem.Name == scraper_name {
			scraping, err := elem.Scraping()
			if err != nil {
				return
			}
			response, results := runner(elem.Name, elem.Version, false)
			SaveResponseToStorage(elem, scraping, response)
			SaveResults(elem, scraping, results)
		}
	}
}

type Test struct {
	Name     string
	Version  int
	FilePath string
	Scraping int
}

func test(scraper_name string, scraper_version int) {
	test := Test{Name: scraper_name, Version: scraper_version}
    fileResponse := test.GetResponseFromStorage()
	storedResults, err := test.ResultsByScraping()
	if err != nil {
		fmt.Println(err)
	}
	save_response_to_file(fileResponse)
	httpResponse, newResults := runner(scraper_name, scraper_version, true)
    remove_file()
    _ = httpResponse

    var bool_array []bool
	for _, stored_element := range storedResults {
        for _, new_element := range newResults {
            if stored_element.ResultUrl == new_element.ResultUrl {
                bool_url := new_element.ResultUrl == stored_element.ResultUrl
                bool_title := new_element.Title == stored_element.Title
                bool_comparison := bool_url == bool_title
                bool_array = append(bool_array, bool_comparison)
            }
        }
    }
    if len(storedResults) == len(bool_array) {
        fmt.Println("TEST STATUS: OK")
    } else {
         fmt.Println("TEST STATUS: KO")       
    }
}

func save_response_to_file(response string) {
    dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	f, err := os.Create(dir + "/response.html")
	if err != nil {
		fmt.Println(err)
	}
    defer f.Close()
	f.WriteString(response)
}

func remove_file() {
    dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	err_2 := os.Remove(dir + "/response.html")
	if err_2 != nil {
		fmt.Println(err_2)
	}
}
