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
)

type Scraping struct {
	Id        int
	ScraperId int
	Uuid      string
	CreatedAt time.Time
}

type Scraper struct {
	Id      int
	Name    string
	Version int
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

	// 0x40 is reserved variant from RFC 4122
	u[8] = (u[8] | 0x40) & 0x7F
	// Set the four most significant bits (bits 12 through 15) of the
	// time_hi_and_version field to the 4-bit version number.
	u[6] = (u[6] & 0xF) | (0x4 << 4)
	uuid = fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
	return
}

// Get all the scraper that need to be executed
func Scrapers() (scrapers []Scraper, err error) {
	fmt.Println("Starting Scrapers...")
	rows, err := Db.Query(`SELECT s.name, MAX(s.version) AS version, MAX(s.id) AS id 
                           FROM scrapers s GROUP BY 1;`)
	if err != nil {
		return
	}
	for rows.Next() {
		scraper := Scraper{}
		if err = rows.Scan(&scraper.Name, &scraper.Version, &scraper.Id); err != nil {
			return
		}
		scrapers = append(scrapers, scraper)
	}
	rows.Close()
	return
}

// Get all the information about a scraper based on its name
func (scraper *Scraper) ScraperByName() (err error) {
	fmt.Println("Starting ScraperByName...")
	err = Db.QueryRow(`SELECT s.id
                       FROM scrapers s WHERE s.name=$1`, scraper.Name).Scan(&scraper.Id)
	fmt.Println("Closing ScraperByName...")
	return
}

// Save in the database a new Scraping session and return its values
func (scraper *Scraper) Scraping() (scraping Scraping, err error) {
	statement := "insert into scraping (uuid, scraper_id, created_at) values ($1, $2, $3) returning id, uuid, scraper_id, created_at"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	// use QueryRow to return a row and scan the returned id into the Session struct
	err = stmt.QueryRow(createUUID(), scraper.Id, time.Now()).Scan(&scraping.Id, &scraping.Uuid, &scraping.ScraperId, &scraping.CreatedAt)
	return
}

// Save all the results extracted
func SaveResults(scraper Scraper, scraping Scraping, results []Result) {
	fmt.Println("Starting SaveResults...")
	for _, elem := range results {
		fmt.Println(scraper.Name)
		fmt.Println("\t", elem.Title)
		fmt.Println("\t\t", elem.Title)
		statement := "INSERT INTO results (uuid, scraper_id, scraping_id, title, url, created_at) VALUES ($1, $2, $3, $4, $5, $6)"
		_, err := Db.Exec(statement, createUUID(), scraper.Id, scraping.Id, elem.Title, elem.ResultUrl, time.Now())
		if err != nil {
			return
		}
	}
}
