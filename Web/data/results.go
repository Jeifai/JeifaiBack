package data

import (
	"fmt"
	"time"
)

type Result struct {
	Target    string
	Id        int
	ScraperId int
	Title     string
	Url       string
	CreatedAt time.Time
}

// Get all the results belonging to the targets of a specific user
func (user *User) ResultsByUser() (results []Result, err error) {
	fmt.Println("Starting ResultsByUser...")
	rows, err := Db.Query(`WITH latest_scraping_per_target AS(
                                SELECT
                                    s.targetid,
                                    MAX(ss.id) AS latest_scraping
                                FROM scrapers s
                                LEFT JOIN scrapings ss ON(s.id = ss.scraperid)
                                GROUP BY 1)
                            SELECT
                                t.url,
                                r.createdat,
                                r.title,
                                r.url
                            FROM userstargets ut
                            LEFT JOIN targets t ON(ut.targetid = t.id)
                            LEFT JOIN scrapers s ON(ut.targetid = s.targetid)
                            LEFT JOIN results r ON(s.id = r.scraperid)
                            LEFT JOIN latest_scraping_per_target ls ON(ut.targetid = ls.targetid)
                            WHERE ut.userid = $1
                            AND ut.deletedat IS NULL
                            AND r.scrapingid = ls.latest_scraping`, user.Id)
	if err != nil {
		return
	}
	for rows.Next() {
		result := Result{}
		if err = rows.Scan(&result.Target, &result.CreatedAt, &result.Title, &result.Url); err != nil {
			return
		}
		results = append(results, result)
	}
	rows.Close()
	return
}
