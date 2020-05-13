package data

import (
	"fmt"
	"time"
)

type Target struct {
	Id        int
	Url       string
	CreatedAt time.Time
}

// Add a new target
func (target *Target) CreateTarget() (err error) {
	fmt.Println("Starting CreateTarget...")
    statement := `INSERT INTO targets (url, created_at)
                  VALUES ($1, $2) RETURNING id, url, created_at`
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
    statement := `INSERT INTO users_targets (uuid, user_id, target_id, created_at) 
                  VALUES ($1, $2, $3, $4) RETURNING id, uuid, user_id, target_id, created_at`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(createUUID(), user.Id, target.Id, time.Now()).Scan()
	fmt.Println("Closing CreateUserTarget...")
	return
}

// Get all the targets for a specific user
func (user *User) UsersTargetsByUser() (targets []Target, err error) {
	fmt.Println("Starting UsersTargetsByUser...")
	rows, err := Db.Query(`SELECT t.id, t.url, t.created_at 
                           FROM users u
                           INNER JOIN users_targets ut ON(u.id = ut.user_id) 
                           INNER JOIN targets t ON(ut.target_id = t.id)
                           WHERE u.id=$1`, user.Id)
	if err != nil {
		return
	}
	for rows.Next() {
		target := Target{}
		if err = rows.Scan(&target.Id, &target.Url, &target.CreatedAt); err != nil {
			return
		}
		targets = append(targets, target)
	}
	rows.Close()
	return
}

// Get the target for a specific user and url, must return a unique value
func (user *User) UsersTargetsByUserAndUrl(url string) (target Target, err error) {
    fmt.Println("Starting UsersTargetsByUserAndUrl...")
	err = Db.QueryRow(`SELECT t.id, t.url, t.created_at 
                       FROM users u
                       INNER JOIN users_targets ut ON(u.id = ut.user_id) 
                       INNER JOIN targets t ON(ut.target_id = t.id)
                       WHERE u.id=$1
                       AND t.url=$2`, user.Id, url).Scan(&target.Id, &target.Url, &target.CreatedAt)
	return
}

// Delete a relation user <--> target
func (target *Target) DeleteUserTargetByUserAndTarget(user User) (err error) {
    fmt.Println("Starting DeleteUserTargetByUserAndTarget...")
    statement := `DELETE FROM users_targets 
                  WHERE user_id = $1
                  AND target_id = $2;`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
    _, err = stmt.Exec(user.Id, target.Id)
	return
}