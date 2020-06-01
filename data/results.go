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
                                    s.target_id,
                                    MAX(ss.id) AS latest_scraping
                                FROM scrapers s
                                LEFT JOIN scraping ss ON(s.id = ss.scraper_id)
                                GROUP BY 1)
                            SELECT
                                t.url,
                                r.created_at,
                                r.title,
                                r.url
                            FROM users_targets ut
                            LEFT JOIN targets t ON(ut.target_id = t.id)
                            LEFT JOIN scrapers s ON(ut.target_id = s.target_id)
                            LEFT JOIN results r ON(s.id = r.scraper_id)
                            LEFT JOIN latest_scraping_per_target ls ON(ut.target_id = ls.target_id)
                            WHERE ut.user_id = $1
                            AND ut.deleted_at IS NULL
                            AND r.scraping_id = ls.latest_scraping`, user.Id)
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
