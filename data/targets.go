package data

import (
	"time"
)

type Target struct {
	Id          int
	Url         string
	CreatedAt   time.Time
}

type Url struct {
	Id          int
	Url         string
	CreatedAt   time.Time
}

// Add a new target
func (url *Url) CreateTarget() (err error) {
	statement := "insert into targets (id, url, created_at) values ($1, $2, $3) returning id, url, created_at"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(url.Url, time.Now()).Scan(&url.Id, &url.Url, &url.CreatedAt)
	return
}