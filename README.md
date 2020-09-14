# JaifaiBack

Here everything happening behind the scenes of Jaifai.

There are currently 4 programs available:

* *scrape*: extract data from the web
* *match*: create relations between job offers and keywords
* *notify*: whenever a match happens, notify all the interested users 
* *check*: identitfy whenever a scrape returns bad results

## Scrape

This is Jeifai's core. Here happens the extraction of data from career's pages.

How the program flows:

* Start main function
* Connect to the database
* Get the scraper
* Start scraping session
* Scrape
* Save results to the database

### scrapers.go
All the scrapers are included in a single file called *scrapers.go*. The most important aspect is to produce scrapers as identical as possibile to each other. The two main tools used to scrape are:
* *[Colly](http://go-colly.org/)*: Golang scraping best library
* *[Chromedp](https://github.com/chromedp/chromedp)*: Run an headless Google Chrome instance

Not all the career pages are identical:
* *HTML*: data are stored into the HTML response, Colly is used.
* *API*: data are returned after an API call, Colly is used.
* *Javascript*: data are stored into the HTML, but after Javascript renders the page, Colly and Chromedp are used.
* *API_POST*: data are returned after an initial API call to get the cookies, Colly and Chromedp are used.
* *Pagination*: if any of the category above presents pagination, it is necessary to implement a logic for it.

### **How to create a new scraper?**
The creation of a new scraper is divided in two different part:
* Add in the database the new target and the new scraper (*Scripts/CreateScraper.go*)
    * Name, career's url and host url are necessary
```golang
func main() {
    scraper_name := "Google"
    jobs_url := "https://www.google.com/careers"
    host_url := "https://www.google.com"
    scraper := Scraper{scraper_name, jobs_url, host_url}
    scraper.CreateScraper()
}
```

* Create the algorithm to scrape in *scrapers.go*
    * Often it is good practice to build and test the scraper in a separate folder.
    * Here an example fo scraper which extract the information directly from the HTML.
```golang
func (runtime Runtime) Morressier() (results Results) {
    c := colly.NewCollector()
    start_url := "https://morressier-jobs.personio.de/"
    type Job struct {
        Url      string
        Title    string
        Location string
        Type     string
    }
    c.OnHTML("a", func(e *colly.HTMLElement) {
        if strings.Contains(e.Attr("class"), "job-box-link") {
            result_title := e.ChildText(".jb-title")
            result_url := e.Attr("href")
            result_description := e.ChildTexts("span")[0]
            result_location := e.ChildTexts("span")[2]
            results.Add(
                runtime.Name,
                result_title,
                result_url,
                result_location,
                Job{
                    result_url,
                    result_title,
                    result_location,
                    result_description,
                },
            )
        }
    })
    c.OnRequest(func(r *colly.Request) {
        fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
    })
    c.OnError(func(r *colly.Response, err error) {
        fmt.Println(Red("Request URL:"), Red(r.Request.URL))
    })
    c.Visit(start_url)
    return
}
```

### **How to run a scraper?**
```bash
go build
./JeifaiBack scrape -s=[scraper_name] -r=[true/false]
```
* -s select any scraper name
* -r true if results need to be saved, false otherwise (might be useful for testing purposes)