package data

import (
	"fmt"
	"time"

	"github.com/lib/pq"
)

type Target struct {
	Id          int
	Url         string `validate:"url"`
	Host        string
	CreatedAt   time.Time
	CreatedDate string
}

// Add a new target
func (target *Target) CreateTarget() (err error) {
	fmt.Println("Starting CreateTarget...")
	statement := `INSERT INTO targets (url, host, createdat)
                  VALUES ($1, $2, $3)
                  RETURNING id, url, host, createdat`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		panic(err.Error())
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
	statement := `INSERT INTO userstargets (userid, targetid, createdat) 
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
	rows, err := Db.Query(`SELECT
                            t.id,
                            t.url,
                            t.createdat,
                            TO_CHAR(t.createdat, 'YYYY-MM-DD')
                           FROM users u
                           INNER JOIN userstargets ut ON(u.id = ut.userid) 
                           INNER JOIN targets t ON(ut.targetid = t.id)
                           WHERE ut.deletedat IS NULL
                           AND u.id=$1
                           ORDER BY t.createdat DESC`, user.Id)
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		target := Target{}
		if err = rows.Scan(
			&target.Id,
			&target.Url,
			&target.CreatedAt,
			&target.CreatedDate); err != nil {
			panic(err.Error())
		}
		targets = append(targets, target)
	}
	rows.Close()
	return
}

// Get the target for a specific url
func (target *Target) TargetByUrl() (err error) {
	fmt.Println("Starting TargetByUrl...")
	err = Db.QueryRow(`SELECT
                         t.id
                       FROM targets t
                       WHERE t.url=$1`, target.Url).Scan(&target.Id)
	return
}

// Get all the targets based on a list of urls
func TargetsByUrls(targetUrls []string) (targets []Target, err error) {
	fmt.Println("Starting TargetsByUrls...")

	rows, err := Db.Query(`SELECT
                                t.id,
                                t.url
                            FROM targets t
                            WHERE t.url LIKE ANY($1)`, pq.Array(targetUrls))
	if err != nil {
		return
	}
	for rows.Next() {
		target := Target{}
		if err = rows.Scan(&target.Id, &target.Url); err != nil {
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
	err = Db.QueryRow(`SELECT
                         t.id, 
                         t.url, 
                         t.createdat 
                       FROM users u
                       INNER JOIN userstargets ut ON(u.id = ut.userid) 
                       INNER JOIN targets t ON(ut.targetid = t.id)
                       WHERE u.id=$1
                       AND t.url=$2
                       AND ut.deletedat IS NULL`, user.Id, url).Scan(&target.Id, &target.Url, &target.CreatedAt)
	return
}

// Update userstargets in column deletedat
func (target *Target) SetDeletedAtInUsersTargetsByUserAndTarget(
	user User) (err error) {
	fmt.Println("Starting SetDeletedAtInUserTargetsByUserAndTarget...")
	statement := `UPDATE userstargets
                  SET deletedat = current_timestamp
                  WHERE userid = $1
                  AND targetid = $2;`

	stmt, err := Db.Prepare(statement)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	_, err = stmt.Exec(user.Id, target.Id)
	return
}
