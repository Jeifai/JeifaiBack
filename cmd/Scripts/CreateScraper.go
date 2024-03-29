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

func main() {
	scraper_name := "Igaming"
	jobs_url := "https://igaming-jobs.personio.de"
	host_url := "https://www.igaming.com"
	scraper := Scraper{scraper_name, jobs_url, host_url}
	scraper.CreateScraper()
}

type Scraper struct {
	Name    string
	JobsUrl string
	HostUrl string
}

func DbConnect() (Db *sql.DB) {
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

func (scraper Scraper) CreateScraper() {
	Db := DbConnect()
	defer Db.Close()
	defer fmt.Println("Successfully created new scraper")

	statement_1 := `INSERT INTO targets (name, url, host, createdat)
                    VALUES ($1, $2, $3, $4) RETURNING id`
	stmt_1, err := Db.Prepare(statement_1)
	if err != nil {
		panic(err.Error())
	}
	defer stmt_1.Close()
	var target_id int
	err = stmt_1.QueryRow(
		scraper.Name, scraper.JobsUrl, scraper.HostUrl, time.Now()).Scan(&target_id)
	if err != nil {
		panic(err.Error())
	}

	statement_2 := `INSERT INTO scrapers (name, targetid, createdat)
                    VALUES ($1, $2, $3) RETURNING id`
	stmt_2, err := Db.Prepare(statement_2)
	if err != nil {
		panic(err.Error())
	}
	defer stmt_2.Close()
	stmt_2.QueryRow(scraper.Name, target_id, time.Now())
}
