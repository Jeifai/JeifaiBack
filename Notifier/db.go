package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	. "github.com/logrusorgru/aurora"
)

type Scraper struct {
	Id      int
	Name    string
	Version int
}

type Notifier struct {
	Id        int
	UserId    int
	CreatedAt time.Time
}

type Notification struct {
	MatchId   int
	UserId    int
	UserName  string
	UserEmail string
	Name      string
	Title     string
	Url       string
}

type Email struct {
	MatchId   int
	UserId    int
	UserName  string
	UserEmail string
	Company   []Company
}

type Company struct {
	Name string
	Job  []Job
}

type Job struct {
	MatchId int
	Title   string
	Url     string
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
		fmt.Println(Bold(Red("Unsuccessfully connected to the database")))
		return
	}
	fmt.Println(Bold(Green("Successfully connected to the database")))
}

func GetScrapers() (scrapers []Scraper, err error) {
	fmt.Println(Gray(8-1, "Starting GetScrapers..."))
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

func PrepareNotifications(scrapers []Scraper) (notifications []Notification, err error) {
	fmt.Println(Gray(8-1, "Starting GetNotifications..."))

	for _, elem := range scrapers {

		scraper_id := elem.Id

		rows, err := Db.Query(`
                            SELECT DISTINCT
                                m.id,
                                u.id,
                                u.username,
                                u.email,
                                s.name,
                                r.title,
                                r.url
                            FROM results r
                            LEFT JOIN matches m ON(r.id = m.resultid)
                            LEFT JOIN scrapers s ON(r.scraperid = s.id)
                            LEFT JOIN userstargetskeywords utk ON(m.keywordid = utk.keywordid)
                            LEFT JOIN users u ON(utk.userid = u.id)
                            LEFT JOIN notifications n ON(m.id = n.matchid)
                            WHERE m.createdat > current_date - interval '0' day
                            AND s.id = $1
                            AND n.id IS NULL
                            AND u.id = 1
                            ORDER BY 1 DESC;`, scraper_id)
		if err != nil {
			panic(err.Error())
		}

		counter := 0
		for rows.Next() {
			notification := Notification{}
			rows.Scan(
				&notification.MatchId,
				&notification.UserId,
				&notification.UserName,
				&notification.UserEmail,
				&notification.Name,
				&notification.Title,
				&notification.Url)
			counter++
			notifications = append(notifications, notification)
		}
		rows.Close()
	}
	return
}

func (notifier *Notifier) StartNotifierSession(user_id int) (err error) {
	fmt.Println(Gray(8-1, "Starting StartNotifierSession..."))
	statement := `INSERT INTO notifiers (userid, createdat)
                  VALUES ($1, $2)
                  RETURNING id, userid, createdat`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	err = stmt.QueryRow(
		user_id,
		time.Now()).Scan(
		&notifier.Id, &notifier.UserId, &notifier.CreatedAt)
	if err != nil {
		panic(err.Error())
	}
	return
}

func SaveNotification(notifier Notifier, email Email) {
	fmt.Println(Gray(8-1, "Starting SaveNotification..."))
	time_now := time.Now()
	valueStrings := []string{}
	valueArgs := []interface{}{}
	var counter int
	for _, company := range email.Company {
		for _, job := range company.Job {
			str1 := "$" + strconv.Itoa(1+counter*3) + ","
			str2 := "$" + strconv.Itoa(2+counter*3) + ","
			str3 := "$" + strconv.Itoa(3+counter*3)
			str_n := "(" + str1 + str2 + str3 + ")"
			valueStrings = append(valueStrings, str_n)
			valueArgs = append(valueArgs, notifier.Id)
			valueArgs = append(valueArgs, job.MatchId)
			valueArgs = append(valueArgs, time_now)
			counter++
		}
	}
	smt := `INSERT INTO notifications(notifierid, matchid, createdat) VALUES %s`
	smt = fmt.Sprintf(smt, strings.Join(valueStrings, ","))

	_, err := Db.Exec(smt, valueArgs...)
	if err != nil {
		panic(err.Error())
	}
}
