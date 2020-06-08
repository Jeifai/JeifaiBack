package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Match struct {
	CreatedAt   time.Time
	CompanyName string
	JobTitle    string
	JobUrl      string
	Keyword     string
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

func GetMatches() (matches []Match, err error) {
	fmt.Println("Starting GetMatches...")
	rows, err := Db.Query(`WITH latest_scraper AS(
                            SELECT
                                ss.name,
                                MAX(s.id) AS id
                            FROM scrapers ss
                            LEFT JOIN scrapings s ON(ss.id = s.scraperid)
                            GROUP BY 1)
                        SELECT
                            r.createdat AS created_at,
                            ls.name AS company,
                            r.title AS job_title,
                            r.url AS job_url,
                            k.text AS keyword_text
                        FROM targets t
                        INNER JOIN scrapers s ON(t.id = s.targetid)
                        INNER JOIN results r ON(s.id = r.scraperid)
                        INNER JOIN latest_scraper ls ON(r.scrapingid = ls.id)
                        LEFT JOIN userstargetskeywords utk ON(t.id = utk.targetid)
                        LEFT JOIN keywords k ON(utk.keywordid = k.id)
                        WHERE r.createdat = r.updatedat
                        AND REPLACE(LOWER(r.title), ' ', '') LIKE '%' || REPLACE(LOWER(k.text), ' ', '') || '%'`)
	if err != nil {
		return
	}
	for rows.Next() {
		match := Match{}
		if err = rows.Scan(
			&match.CreatedAt,
			&match.CompanyName,
			&match.JobTitle,
			&match.JobUrl,
			&match.Keyword); err != nil {
			return
		}
		matches = append(matches, match)
	}
	rows.Close()
	return
}
