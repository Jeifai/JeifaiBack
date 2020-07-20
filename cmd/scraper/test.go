package cmd

func Runtest(company string) string {
    if company == "all" {
        return "running scraper for all the company"
    } else {
        return "running scraper for single company"
    }
}