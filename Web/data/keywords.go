package data

import (
	"fmt"
	"strconv"
	"strings"
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
	utk UserTargetKeyword, err error) {
	fmt.Println("Starting SetUserTargetKeyword...")

	valueStrings := []string{}
	valueArgs := []interface{}{}
	for i, elem := range targets {
		str1 := "$" + strconv.Itoa(1+i*4) + ","
		str2 := "$" + strconv.Itoa(2+i*4) + ","
		str3 := "$" + strconv.Itoa(3+i*4) + ","
		str4 := "$" + strconv.Itoa(4+i*4)
		str_n := "(" + str1 + str2 + str3 + str4 + ")"
		valueStrings = append(valueStrings, str_n)
		valueArgs = append(valueArgs, user.Id)
		valueArgs = append(valueArgs, elem.Id)
		valueArgs = append(valueArgs, keyword.Id)
		valueArgs = append(valueArgs, time.Now())
	}

	smt := `INSERT INTO userstargetskeywords (
                userid, targetid, keywordid, createdat)
            VALUES %s`
	smt = fmt.Sprintf(smt, strings.Join(valueStrings, ","))

	_, err = Db.Exec(smt, valueArgs...)
	if err != nil {
		panic(err.Error())
	}

	return
}
