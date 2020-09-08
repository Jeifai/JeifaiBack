# Scrapers

This is Jeifai's core. Here happens the extraction of data from career's pages.

How the program flows:

* Start main function
* Connect to the database
* Extract all the targets
* For each target
    * Start scraping session
    * Scrape
    * Save results to the database

## scrapers.go
All the scrapers are included in a single file called *scrapers.go*.
In order to ensure fully flexibility, *refect* is used to programmatically invoque which target to scrape.

The main idea is to produce scrapers as identical as possibile to each other.
To achive so, Colly is used and each scraper much use it.

When a JavaScript render is required, chromedp is used, sometimes also to produce action on a headless browser instance.

Moreover, the structure and the style of each scraper must follow a specific set of rules defined below.

## **How to run a scraper?**
	sudo go build
    ./Jeifaiback scrape -s=[scraper_name] -r=[true/false]

-s]select any scraper name
-r] true if results need to be saved, false otherwise (might be useful for testing purposes)

## **How to create a new scraper?**
The creation of a new scraper is divided in two different part:
* Add in the database the new target and the new scraper (*Scripts/CreateScraper.go*)
    * Name, career's url and host url are necessary
* Create the algorithm to scrape in *scrapers.go*
    * Often it is good practice to build and test the scraper in a separate folder.

## **How to format and write the code for the scraper?**
Most of the scrapers follow these rules, created to make *scrapers.go* looking uniformed and standard.
```golang
    func(runtime Runtime) ScraperName(version int) (results []Result) {
        switch version {
        case n:
            c := colly.NewCollector()   // define Collector
            url := "https://robimalco.github.io/dreamingjobs.github.io/" // define url
            tag_department := "li[class=department]" // define variables
            type Job struct {       // define Struct
                Title      string
                Url        string
                ...
            }
            // Write scraper logic
            c.OnHTML() {}
            c.OnResponse() {}
            c.OnRequest() {}
            c.OnScraped() {}
            c.OnError() {}

            c.Visit(url)
        }
    }
```
