package data

import (
	"fmt"
	"time"
)

type Job struct {
	Target    string
	Id        int
	ScraperId int
	Title     string
	Url       string
	CreatedAt time.Time
}

// Get all the jobs belonging to the targets of a specific user
func (user *User) JobsByUser() (jobs []Job, err error) {
	fmt.Println("Starting JobsByUser...")
	rows, err := Db.Query(`WITH latest_scraping_per_target AS(
                                SELECT
                                    s.target_id,
                                    MAX(ss.id) AS latest_scraping
                                FROM scrapers s
                                LEFT JOIN scraping ss ON(s.id = ss.scraper_id)
                                GROUP BY 1)
                            SELECT t.url, j.created_at, j.title, j.url
                            FROM users_targets ut
                            LEFT JOIN targets t ON(ut.target_id = t.id)
                            LEFT JOIN scrapers s ON(ut.target_id = s.target_id)
                            LEFT JOIN jobs j ON(s.id = j.scraper_id)
                            LEFT JOIN latest_scraping_per_target ls ON(ut.target_id = ls.target_id)
                            WHERE ut.user_id = $1
                            AND j.scraping_id = ls.latest_scraping`, user.Id)
	if err != nil {
		return
	}
	for rows.Next() {
		job := Job{}
		if err = rows.Scan(&job.Target, &job.CreatedAt, &job.Title, &job.Url); err != nil {
			return
		}
		jobs = append(jobs, job)
	}
	rows.Close()
	return
}
