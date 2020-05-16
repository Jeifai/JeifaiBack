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
	rows, err := Db.Query(`SELECT t.url, j.created_at, j.title, j.url
                            FROM users_targets ut
                            LEFT JOIN targets t ON(ut.target_id = t.id)
                            LEFT JOIN scrapers s ON(ut.target_id = s.target_id)
                            LEFT JOIN jobs j ON(s.id = j.scraper_id)
                            WHERE ut.user_id = $1`, user.Id)
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
