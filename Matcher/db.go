package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Match struct {
	Title string
	Text  string
}

var Db *sql.DB

func DbConnect() {
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
}

func GetMatches(scraper_name string) (matches []Match, err error) {
	fmt.Println("Starting GetMatches...")
	rows, err := Db.Query(`WITH latest_scraper AS(
                            SELECT
                                MAX(s.id) AS id
                            FROM scrapers ss
                            LEFT JOIN scrapings s ON(ss.id = s.scraperid)
                            WHERE ss.name = $1)
                        SELECT
                            r.title,
                            k.text
                        FROM targets t
                        INNER JOIN scrapers s ON(t.id = s.targetid)
                        INNER JOIN results r ON(s.id = r.scraperid)
                        INNER JOIN latest_scraper ls ON(r.scrapingid = ls.id)
                        LEFT JOIN userstargetskeywords utk ON(t.id = utk.targetid)
                        LEFT JOIN keywords k ON(utk.keywordid = k.id)
                        WHERE r.createdat = r.updatedat
                        AND REPLACE(LOWER(r.title), ' ', '') LIKE '%' || REPLACE(LOWER(k.text), ' ', '') || '%'`,
		scraper_name)
	if err != nil {
		return
	}
	for rows.Next() {
		match := Match{}
		if err = rows.Scan(
			&match.Title, &match.Text); err != nil {
			return
		}
		matches = append(matches, match)
	}
	rows.Close()
	return
}
