package main

import (
	"fmt"
	"time"

	"github.com/lib/pq"
)

type Target struct {
	Id          int
	Url         string
	Host        string
	Name        string `validate:"required,max=30,min=3"`
	CreatedAt   time.Time
	CreatedDate string
}

// Add a new target
func (target *Target) CreateTarget() (err error) {
	fmt.Println("Starting CreateTarget...")
	statement := `INSERT INTO targets (name, createdat)
                  VALUES ($1, $2)
                  RETURNING id, name, createdat`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	err = stmt.QueryRow(
		target.Name,
		time.Now(),
	).Scan(
		&target.Id,
		&target.Name,
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
                            t.name,
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
			&target.Name,
			&target.CreatedAt,
			&target.CreatedDate); err != nil {
			panic(err.Error())
		}
		targets = append(targets, target)
	}
	rows.Close()
	return
}

// Get the target for a specific name
func (target *Target) TargetByName() (err error) {
	fmt.Println("Starting TargetByName...")
	err = Db.QueryRow(`SELECT
                         t.id
                       FROM targets t
                       WHERE t.name=$1`, target.Name).Scan(&target.Id)
	return
}

// Get all the targets based on a list of names
func TargetsByNames(targetNames []string) (targets []Target, err error) {
	fmt.Println("Starting TargetsByNames...")

	rows, err := Db.Query(`SELECT
                                t.id,
                                t.name
                            FROM targets t
                            WHERE t.name LIKE ANY($1)`, pq.Array(targetNames))
	if err != nil {
		return
	}
	for rows.Next() {
		target := Target{}
		if err = rows.Scan(&target.Id, &target.Name); err != nil {
			return
		}
		targets = append(targets, target)
	}
	rows.Close()
	return
}

// Get the target for a specific user and name, must return a unique value
func (user *User) UsersTargetsByUserAndName(name string) (target Target, err error) {
	fmt.Println("Starting UsersTargetsByUserAndName...")
	err = Db.QueryRow(`SELECT
                         t.id, 
                         t.name, 
                         t.createdat 
                       FROM users u
                       INNER JOIN userstargets ut ON(u.id = ut.userid) 
                       INNER JOIN targets t ON(ut.targetid = t.id)
                       WHERE u.id=$1
                       AND t.name=$2
                       AND ut.deletedat IS NULL`, user.Id, name).Scan(&target.Id, &target.Name, &target.CreatedAt)
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
