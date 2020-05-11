package data

import (
    "time"
    "fmt"
)

type Target struct {
	Id          int
	Url         string
	CreatedAt   time.Time
}

// Add a new target
func (target *Target) CreateTarget() (err error) {
    fmt.Println("Starting CreateTarget...")
	statement := "insert into targets (url, created_at) values ($1, $2) returning id, url, created_at"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
    }
    defer stmt.Close()
    err = stmt.QueryRow(target.Url, time.Now()).Scan(&target.Id, &target.Url, &target.CreatedAt)
    fmt.Println("Closing CreateTarget...")
	return err
}

// Add a new relation user <--> target
func (target *Target) CreateUserTarget(user User) (err error) {
    fmt.Println("Starting CreateUserTarget...")
	statement := "insert into users_targets (uuid, user_id, target_id, created_at) values ($1, $2, $3, $4) returning id, uuid, user_id, target_id, created_at"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
    }
    defer stmt.Close()
    err = stmt.QueryRow(createUUID(), user.Id, target.Id, time.Now()).Scan()
    fmt.Println("Closing CreateUserTarget...")
	return
}