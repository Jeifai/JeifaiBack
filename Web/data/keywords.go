package data

import (
	"fmt"
	"time"
)

type Keyword struct {
	Id        int
	Text      string
	CreatedAt time.Time
}

type UserTargetKeyword struct {
	Id          int
	UserId      int
	TargetId    int
	KeywordId   int
	CreatedAt   time.Time
	KeywordText string
	TargetUrl   string
}

// Add a new keyword
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
                                k.text,
                                t.url
                            FROM userstargetskeywords utk
                            LEFT JOIN keywords k ON(utk.keywordid = k.id)
                            LEFT JOIN targets t ON(utk.targetid = t.id)
                            WHERE utk.userid = $1
                            AND utk.deletedat IS NULL`, user.Id)
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
			&utk.KeywordText,
			&utk.TargetUrl); err != nil {
			return
		}
		utks = append(utks, utk)
	}
	rows.Close()
	return
}

func SetUserTargetKeyword(
	user User, targets []Target, keyword Keyword) (
	utks []UserTargetKeyword, err error) {
	fmt.Println("Starting SetUserTargetKeyword...")

	for _, elem := range targets {
		statement := `INSERT INTO userstargetskeywords (
                        userid, targetid, keywordid, createdat)
                        VALUES ($1, $2, $3, $4)
                        RETURNING id, userid, targetid, keywordid, createdat`
		stmt, err := Db.Prepare(statement)
		if err != nil {
			panic(err.Error())
		}
		defer stmt.Close()
		utk := UserTargetKeyword{}
		err = stmt.QueryRow(
			user.Id,
			elem.Id,
			keyword.Id,
			time.Now(),
		).Scan(
			&utk.Id,
			&utk.UserId,
			&utk.TargetId,
			&utk.KeywordId,
			&utk.CreatedAt,
		)
		utk.KeywordText = keyword.Text
		utk.TargetUrl = elem.Url
		utks = append(utks, utk)
	}
	return
}
