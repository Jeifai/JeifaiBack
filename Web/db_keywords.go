package main

import (
	"fmt"
	"time"
)

type Keyword struct {
	Id        int
	Text      string `validate:"required,max=30,min=3"`
	CreatedAt time.Time
}

type UserTargetKeyword struct {
	Id          int
	UserId      int
	TargetId    int
	KeywordId   int
	CreatedAt   time.Time
	CreatedDate string
	UpdatedAt   time.Time
	KeywordText string
	TargetUrl   string
	TargetName  string
}

func (keyword *Keyword) CreateKeyword() (err error) {
	fmt.Println("Starting CreateKeyword...")
	statement := `INSERT INTO keywords (text, createdat)
                  VALUES ($1, current_timestamp)
                  RETURNING id, createdat`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	err = stmt.QueryRow(
		keyword.Text,
	).Scan(
		&keyword.Id,
		&keyword.CreatedAt,
	)
	return err
}

func (keyword *Keyword) KeywordByText() (err error) {
	fmt.Println("Starting KeywordByText...")
	err = Db.QueryRow(`SELECT
                         k.id
                       FROM keywords k
                       WHERE k.text=$1`, keyword.Text).Scan(&keyword.Id)
	return
}

func (user *User) GetUserTargetKeyword() (
	utks []UserTargetKeyword, err error) {
	fmt.Println("Starting GetUserTargetKeyword...")

	rows, err := Db.Query(`SELECT
                                utk.id,
                                utk.userid,
                                utk.targetid,
                                utk.keywordid,
                                utk.createdat,
                                TO_CHAR(utk.createdat, 'YYYY-MM-DD'),
                                utk.updatedat,
                                k.text,
                                t.name
                            FROM userstargetskeywords utk
                            LEFT JOIN keywords k ON(utk.keywordid = k.id)
                            LEFT JOIN targets t ON(utk.targetid = t.id)
                            WHERE utk.userid = $1
                            AND utk.deletedat IS NULL
                            ORDER BY utk.updatedat DESC`, user.Id)
	if err != nil {
		return
	}
	for rows.Next() {
		utk := UserTargetKeyword{}
		if err = rows.Scan(
			&utk.Id,
			&utk.UserId,
			&utk.TargetId,
			&utk.KeywordId,
			&utk.CreatedAt,
			&utk.CreatedDate,
			&utk.UpdatedAt,
			&utk.KeywordText,
			&utk.TargetName); err != nil {
			return
		}
		utks = append(utks, utk)
	}
	rows.Close()
	return
}

func SetUserTargetKeyword(
	user User, targets []Target, keyword Keyword) (err error) {
	fmt.Println("Starting SetUserTargetKeyword...")

	for _, elem := range targets {
		statement := `INSERT INTO userstargetskeywords (
                        userid, targetid, keywordid, createdat)
                        VALUES ($1, $2, $3, $4)
                        ON CONFLICT (userid, targetid, keywordid) 
                        DO UPDATE SET deletedat = NULL, updatedat = current_timestamp`
		stmt, err := Db.Prepare(statement)
		if err != nil {
			panic(err.Error())
		}
		defer stmt.Close()
		stmt.QueryRow(
			user.Id,
			elem.Id,
			keyword.Id,
			time.Now(),
		)
	}
	return
}

func (utk *UserTargetKeyword) SetDeletedAtIntUserTargetKeyword() (err error) {
	fmt.Println("Starting SetDeletedAtIntUserTargetKeyword...")
	statement := `UPDATE userstargetskeywords
                  SET deletedat = current_timestamp
                  WHERE userid = $1
                  AND targetid = $2
                  AND keywordid = $3;`

	stmt, err := Db.Prepare(statement)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	_, err = stmt.Exec(utk.UserId, utk.TargetId, utk.KeywordId)
	return
}
