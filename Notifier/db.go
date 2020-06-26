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

type Scraper struct {
	Id      int
	Name    string
	Version int
}

type Notifier struct {
	Id        int
	ScraperId int
	CreatedAt time.Time
}

type Notification struct {
	UserName    string
	UserEmail   string
	Name        string
	Title       string
	Url         string
}

type Email struct {
    UserName    string
    UserEmail   string
	Company     []Company
}

type Company struct {
	Name string
	Job         []Job
}

type Job struct {
	Title string
	Url   string
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

func GetScrapers() (scrapers []Scraper, err error) {
	fmt.Println("Starting GetScrapers...")
	rows, err := Db.Query(`SELECT
                                s.name, 
                                MAX(s.version) AS version, 
                                MAX(s.id) AS id 
                           FROM scrapers s
                           GROUP BY 1;`)
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		scraper := Scraper{}
		if err = rows.Scan(
			&scraper.Name,
			&scraper.Version,
			&scraper.Id); err != nil {
			return
		}
		scrapers = append(scrapers, scraper)
	}
	rows.Close()
	return
}

func (notifier *Notifier) StartNotifierSession(scraper_id int) (err error) {
	fmt.Println("Starting StartNotifierSession...")
	statement := `INSERT INTO notifiers (scraperid, createdat)
                  VALUES ($1, $2)
                  RETURNING id`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	err = stmt.QueryRow(
		scraper_id,
		time.Now()).Scan(
		&notifier.Id)
	if err != nil {
		panic(err.Error())
	}
	return
}

func GetNotifications(scrapers []Scraper) (notifications []Notification, err error) {
    fmt.Println("Starting GetNotifications...")

    for _, elem := range scrapers {

        //notifier := Notifier{1, 1, time.Now()}
        // notifier.StartNotifierSession(elem.id)

        rows, err := Db.Query(`
                            SELECT DISTINCT
                                u.username,
                                u.email,
                                s.name,
                                r.title,
                                r.url
                            FROM results r
                            INNER JOIN matches m ON(r.id = m.resultid)
                            LEFT JOIN scrapers s ON(r.scraperid = s.id)
                            LEFT JOIN notifications n ON(m.id = n.matchid)
                            LEFT JOIN userstargetskeywords utk ON(m.keywordid = utk.keywordid)
                            LEFT JOIN users u ON(utk.userid = u.id)
                            WHERE m.createdat > current_date - interval '1' day
                            AND s.id = $1
                            AND n.id IS NULL
                            AND u.id = 1
                            ORDER BY 1 DESC;`, elem.Id)
        if err != nil {
		    panic(err.Error())
        }
        for rows.Next() {
            notification := Notification{}
            rows.Scan(
                &notification.UserName,
                &notification.UserEmail,
                &notification.Name,
                &notification.Title,
                &notification.Url)
            notifications = append(notifications, notification)
        }
        rows.Close()
    }
	return
}
