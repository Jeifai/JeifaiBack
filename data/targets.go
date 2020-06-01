package data

import (
	"fmt"
	"time"
)

type Target struct {
	Id        int
	Url       string
	Host      string
	CreatedAt time.Time
}

// Add a new target
func (target *Target) CreateTarget() (err error) {
	fmt.Println("Starting CreateTarget...")
	statement := `INSERT INTO targets (url, host, created_at)
                  VALUES ($1, $2, $3) RETURNING id, url, host, created_at`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		fmt.Println("Error on CreateTarget")
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(
		target.Url,
		target.Host,
		time.Now(),
	).Scan(
		&target.Id,
		&target.Url,
		&target.Host,
		&target.CreatedAt,
	)
	return err
}

// Add a new relation user <--> target
func (target *Target) CreateUserTarget(user User) {
	fmt.Println("Starting CreateUserTarget...")
	statement := `INSERT INTO users_targets (user_id, target_id, created_at) 
                  VALUES ($1, $2, $3)`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	stmt.QueryRow(user.Id, target.Id, time.Now())
}

// Get all the targets for a specific user
func (user *User) UsersTargetsByUser() (targets []Target, err error) {
	fmt.Println("Starting UsersTargetsByUser...")
	rows, err := Db.Query(`SELECT t.id, t.url, t.created_at 
                           FROM users u
                           INNER JOIN users_targets ut ON(u.id = ut.user_id) 
                           INNER JOIN targets t ON(ut.target_id = t.id)
                           WHERE ut.deleted_at IS NULL
                           AND u.id=$1`, user.Id)
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		target := Target{}
		if err = rows.Scan(&target.Id, &target.Url, &target.CreatedAt); err != nil {
			panic(err.Error())
		}
		targets = append(targets, target)
	}
	rows.Close()
	return
}

// Get all the targets for a specific url
func (target *Target) TargetsByUrl() (err error) {
	fmt.Println("Starting TargetsByUrl...")
	err = Db.QueryRow(`SELECT t.id FROM targets t WHERE t.url=$1`, target.Url).Scan(&target.Id)
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
                       AND t.url=$2
                       AND ut.deleted_at IS NULL`, user.Id, url).Scan(&target.Id, &target.Url, &target.CreatedAt)
	return
}

// Update users_targets in column deleted_at
func (target *Target) SetDeletedAtInUsersTargetsByUserAndTarget(
	user User) (err error) {
	fmt.Println("Starting SetDeletedAtInUserTargetsByUserAndTarget...")
	statement := `UPDATE users_targets
                  SET deleted_at = current_timestamp
                  WHERE user_id = $1
                  AND target_id = $2;`

	stmt, err := Db.Prepare(statement)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	_, err = stmt.Exec(user.Id, target.Id)
	return
}
