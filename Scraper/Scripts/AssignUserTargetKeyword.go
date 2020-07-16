package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var Db *sql.DB

func main() {
	users := []int{1, 2, 17}
	targets := []int{180}
	CreateUserTarget(users, targets)
	CreateUserTargetKeywords(users, targets)
}

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}
	dbhost := os.Getenv("DBHOST")
	dbuser := os.Getenv("DBUSER")
	dbport, err := strconv.ParseInt(os.Getenv("DBPORT"), 10, 64)
	if err != nil {
		panic(err.Error())
	}
	dbname := os.Getenv("DBNAME")
	dbpassword := os.Getenv("DBPASSWORD")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbhost, dbport, dbuser, dbpassword, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err.Error())
	}
	Db = db

	if err = Db.Ping(); err != nil {
		Db.Close()
		fmt.Println("Unsuccessfully connected to the database")
		return
	}
	fmt.Println("Successfully connected to the database")
	return
}

func CreateUserTarget(users []int, targets []int) {
	for _, userid := range users {
		for _, targetid := range targets {
			statement := `INSERT INTO userstargets (userid, targetid, createdat) 
                            VALUES ($1, $2, current_timestamp);`
			stmt, err := Db.Prepare(statement)
			if err != nil {
				panic(err.Error())
			}
			defer stmt.Close()
			stmt.QueryRow(userid, targetid)
			fmt.Println("CreateUserTarget", userid, targetid, " --> DONE")
		}
	}
}

func CreateUserTargetKeywords(users []int, targets []int) {
	for _, userid := range users {
		keywords := GetUserTargetKeyword(userid)
		for _, keywordid := range keywords {
			for _, targetid := range targets {
				statement := `INSERT INTO userstargetskeywords (userid, targetid, keywordid, createdat, updatedat) 
                                VALUES ($1, $2, $3, current_timestamp, current_timestamp);`
				stmt, err := Db.Prepare(statement)
				if err != nil {
					panic(err.Error())
				}
				defer stmt.Close()
				stmt.QueryRow(userid, targetid, keywordid)
				fmt.Println("CreateUserTargetKeywords", userid, targetid, keywordid, " --> DONE")
			}
		}
	}
}

func GetUserTargetKeyword(userid int) (keywords []int) {
	rows, err := Db.Query(`SELECT
                                DISTINCT utk.keywordid
                            FROM userstargetskeywords utk
                            WHERE utk.userid = $1
                            AND utk.deletedat IS NULL`, userid)
	if err != nil {
		return
	}
	for rows.Next() {
		var keyword int
		if err = rows.Scan(
			&keyword); err != nil {
			return
		}
		keywords = append(keywords, keyword)
	}
	rows.Close()
	return
}
