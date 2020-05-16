package main

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strconv"
    "time"
	"strings"
)

type Scraper struct {
	Id      int
	Version int
	Name    string
}

var Db *sql.DB

func init() {

	// Load Environmental Variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbhost := os.Getenv("DBHOST")
	dbuser := os.Getenv("DBUSER")
	dbport, err := strconv.ParseInt(os.Getenv("DBPORT"), 10, 64)
	dbname := os.Getenv("DBNAME")
	dbpassword := os.Getenv("DBPASSWORD")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbhost, dbport, dbuser, dbpassword, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	Db = db
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to the database")
}

func createUUID() (uuid string) {
	u := new([16]byte)
	_, err := rand.Read(u[:])
	if err != nil {
		return
	}
	u[8] = (u[8] | 0x40) & 0x7F
	u[6] = (u[6] & 0xF) | (0x4 << 4)
	uuid = fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
	return
}

// Get all the scraper that need to be executed
func Scrapers() (scrapers []Scraper, err error) {
	fmt.Println("Starting Scrapers...")
	rows, err := Db.Query(`SELECT s.id, s.version, s.name FROM scrapers s`)
	if err != nil {
		return
	}
	for rows.Next() {
		scraper := Scraper{}
		if err = rows.Scan(&scraper.Id, &scraper.Version, &scraper.Name); err != nil {
			return
		}
		scrapers = append(scrapers, scraper)
	}
	rows.Close()
	return
}

// Save all the jobs extracted
func SaveJobs(scraper Scraper, jobs []Job) {
	fmt.Println("Starting SaveJobs...")
	for _, elem := range jobs {
		fmt.Println(elem.JobTitle)
        statement := `INSERT INTO jobs (uuid, scraper_id, job_title, job_url, created_at) 
                      VALUES ($1, $2, $3, $4, $5)`
		_, err := Db.Exec(
            statement,
            createUUID(),
            scraper.Id,
            elem.JobTitle,
            elem.JobUrl,
            time.Now())
		if strings.Contains(err.Error(), "duplicate key value violates") {
            fmt.Println(err)
		} else {
            panic(err)
        }
	}
}