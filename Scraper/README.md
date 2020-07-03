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
    * Save response to Cloud Storage

## scrapers.go
All the scrapers are included in a single file called *scrapers.go*.
In order to ensure fully flexibility, *refect* is used to programmatically invoque which target to scrape.

The main idea is to produce scrapers as identical as possibile to each other.
To achive so, Colly is used and each scraper much use it.

Moreover, the structure and the style of each scraper must follow a specific set of rules defined below.

## storage.go
Here all the functionto communicate with Google Cloud Storage' api.
Saving the response to Cloud Storage allows the possibility to run test without having to fetch and disturb the targets. 

## **How to run a scraper?**
	go build
    ./Scraper

At the moment in order to specify which scraper to run, the *main.go* file must be modified.

## **How to create a new scraper?**
The creation of a new scraper is divided in two different part:
* Add in the database the new target and the new scraper (*Scripts/CreateScraper.go*)
    * Name, career's url and host url are necessary
* Create the algorithm to scrape in *scrapers.go*
    * Often it is good practice to build the scraper in a separate folder.

## **How to format and write the code for the scraper?**
Here all the rules created to make *scrapers.go* looking uniformed and standard.
```golang
    func(runtime Runtime) ScraperName(version int, isLocal bool) (response Response, results []Result) {
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

            //Define how to get response, if local or not
            if isLocal {
                c.Visit("localFile.html")
            } else {
                c.Visit(url)
            }
        }
    }
```
