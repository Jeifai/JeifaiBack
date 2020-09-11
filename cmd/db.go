package cmd

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

type Scraping struct {
	Id        int
	ScraperId int
	CreatedAt time.Time
}

type Scraper struct {
	Id      int
	Name    string
	Version int
}

type Matching struct {
	Id        int
	CreatedAt time.Time
}

type Match struct {
	Id          int
	JobTitle    string
	JobUrl      string
	KeywordText string
	KeywordId   int
	ResultId    int
	CreatedAt   time.Time
	MatchingId  int
}

type Notifier struct {
	Id        int
	UserId    int
	CreatedAt time.Time
}

type Notification struct {
	MatchId     int
	UserId      int
	UserName    string
	UserEmail   string
	CompanyName string
	Title       string
	Url         string
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
		fmt.Println(Bold(Red("Unsuccessfully connected to the database")))
		return
	}
	fmt.Println(Bold(Green("Successfully connected to the database")))
}

func GetScraper(company string) (scraper Scraper) {
	fmt.Println(Gray(8-1, "Starting GetScraper..."))
	err := Db.QueryRow(`SELECT
                            s.name, 
                            MAX(s.version) AS version, 
                            MAX(s.id) AS id 
                        FROM scrapers s
                        WHERE s.name=$1
                        GROUP BY 1;`,
		company).Scan(
		&scraper.Name,
		&scraper.Version,
		&scraper.Id)
	if err != nil {
		panic(err.Error())
	}
	return
}

func GetScrapers() (scrapers []Scraper) {
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

func (scraper *Scraper) StartScrapingSession(count_results int) (scraping Scraping) {
	fmt.Println(Gray(8-1, "Starting StartScrapingSession..."))
	statement := `INSERT INTO scrapings (scraperid, countresults, createdat)
                  VALUES ($1, $2, $3) 
                  RETURNING id, scraperid, createdat`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	err = stmt.QueryRow(scraper.Id, count_results, time.Now()).Scan(
		&scraping.Id, &scraping.ScraperId, &scraping.CreatedAt)
	if err != nil {
		panic(err.Error())
	}
	return
}

func SaveResults(scraper Scraper, scraping Scraping, results []Result) {
	fmt.Println(Gray(8-1, "Starting SaveResults..."))
	valueStrings := []string{}
	valueArgs := []interface{}{}
	timeNow := time.Now() // updatedAt and createdAt will be identical

	var all_urls []string
	for i, elem := range results {
		if !Contains(all_urls, elem.ResultUrl) {
			all_urls = append(all_urls, elem.ResultUrl)
			str1 := "$" + strconv.Itoa(1+i*8) + ","
			str2 := "$" + strconv.Itoa(2+i*8) + ","
			str3 := "$" + strconv.Itoa(3+i*8) + ","
			str4 := "$" + strconv.Itoa(4+i*8) + ","
			str5 := "$" + strconv.Itoa(5+i*8) + ","
			str6 := "$" + strconv.Itoa(6+i*8) + ","
			str7 := "$" + strconv.Itoa(7+i*8) + ","
			str8 := "$" + strconv.Itoa(8+i*8)
			str_n := "(" + str1 + str2 + str3 + str4 + str5 + str6 + str7 + str8 + ")"
			valueStrings = append(valueStrings, str_n)
			valueArgs = append(valueArgs, scraper.Id)
			valueArgs = append(valueArgs, scraping.Id)
			valueArgs = append(valueArgs, elem.Title)
			valueArgs = append(valueArgs, elem.ResultUrl)
			valueArgs = append(valueArgs, elem.Location)
			valueArgs = append(valueArgs, elem.Data)
			valueArgs = append(valueArgs, timeNow)
			valueArgs = append(valueArgs, timeNow)
		}
	}
	smt := `INSERT INTO results (
                scraperid, scrapingid, title, url, location, data, createdat, updatedat)
            VALUES %s ON CONFLICT (url) DO UPDATE
            SET scrapingid = EXCLUDED.scrapingid,
                title = EXCLUDED.title,
                location = EXCLUDED.location,
                updatedat = EXCLUDED.updatedat,
                data = EXCLUDED.data`
	smt = fmt.Sprintf(smt, strings.Join(valueStrings, ","))

	_, err := Db.Exec(smt, valueArgs...)
	if err != nil {
		panic(err.Error())
	}
}

func LastScrapingByNameVersion(
	scraper_name string, scraper_version int) (scraping int) {
	fmt.Println(Gray(8-1, "Starting LastScrapingByNameVersion..."))
	err := Db.QueryRow(`SELECT MAX(s.id)
                       FROM scrapings s 
                       LEFT JOIN scrapers ss ON(s.scraperid = ss.id)
                       LEFT JOIN targets t ON(ss.targetid = t.id)
                       WHERE t.name = $1 AND ss.version = $2;`,
		scraper_name, scraper_version).Scan(&scraping)
	if err != nil {
		panic(err.Error())
	}
	return
}

func ResultsByScraping(scraping int) (results []Result, err error) {
	fmt.Println(Gray(8-1, "Starting ResultsByScraping..."))
	rows, err := Db.Query(`SELECT
                                t.name,
                                r.title, 
                                r.url,
                                r.data
                           FROM results r
                           LEFT JOIN scrapings s ON(r.scrapingid = s.id)
                           LEFT JOIN scrapers ss ON(s.scraperid = ss.id)
                           LEFT JOIN targets t ON(ss.targetid = t.id)
                           WHERE r.scrapingid = $1`, scraping)
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		result := Result{}
		if err = rows.Scan(
			&result.CompanyName,
			&result.Title,
			&result.ResultUrl,
			&result.Data); err != nil {
			return
		}
		results = append(results, result)
	}
	rows.Close()
	return
}

func (matching *Matching) StartMatchingSession(scraper_id int) {
	fmt.Println(Gray(8-1, "Starting StartMatchingSession..."))
	statement := `INSERT INTO matchings (scraperid, createdat)
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
		&matching.Id)
	if err != nil {
		panic(err.Error())
	}
	return
}

func GenerateMatches(matching Matching, scraper_id int) (matches []Match) {
	fmt.Println(Gray(8-1, "Starting GenerateMatches..."))
	rows, err := Db.Query(`
                        WITH 
                            active_keywords AS(
                                SELECT DISTINCT
                                    utk.keywordid,
                                    k.text,
                                    REPLACE(LOWER(k.text), ' ', '') AS lowertext
                                FROM userstargetskeywords utk
                                LEFT JOIN keywords k ON(utk.keywordid = k.id)
                                WHERE utk.deletedat IS NULL),
                            today_results AS(
                                SELECT
                                    r.id AS resultid,
                                    r.title,
                                    r.url,
                                    REPLACE(LOWER(r.title), ' ', '') AS lowertitle
                                FROM results r
                                WHERE r.scraperid = $1
                                AND r.createdat::date = current_timestamp::date)
                        SELECT
                            tr.resultid,
                            ak.keywordid,
                            tr.title,
                            tr.url,
                            ak.text
                        FROM today_results tr
                        INNER JOIN active_keywords ak ON(tr.lowertitle LIKE '%' || ak.lowertext || '%');`, scraper_id)
	for rows.Next() {
		match := Match{}
		if err = rows.Scan(
			&match.ResultId,
			&match.KeywordId,
			&match.JobTitle,
			&match.JobUrl,
			&match.KeywordText,
		); err != nil {
			return
		}
		match.MatchingId = matching.Id
		matches = append(matches, match)
	}
	rows.Close()
	if err != nil {
		panic(err.Error())
	}
	return
}

func SaveMatches(matching Matching, matches []Match) {
	fmt.Println(Gray(8-1, "Starting SaveMatches..."))
	valueStrings := []string{}
	valueArgs := []interface{}{}
	timeNow := time.Now()
	for i, elem := range matches {
		str1 := "$" + strconv.Itoa(1+i*4) + ","
		str2 := "$" + strconv.Itoa(2+i*4) + ","
		str3 := "$" + strconv.Itoa(3+i*4) + ","
		str4 := "$" + strconv.Itoa(4+i*4)
		str_n := "(" + str1 + str2 + str3 + str4 + ")"
		valueStrings = append(valueStrings, str_n)
		valueArgs = append(valueArgs, elem.MatchingId)
		valueArgs = append(valueArgs, elem.ResultId)
		valueArgs = append(valueArgs, elem.KeywordId)
		valueArgs = append(valueArgs, timeNow)
	}
	smt := `INSERT INTO matches (
                matchingid, resultid, keywordid, createdat) 
            VALUES %s ON CONFLICT DO NOTHING` //
	smt = fmt.Sprintf(smt, strings.Join(valueStrings, ","))

	_, err := Db.Exec(smt, valueArgs...)
	if err != nil {
		panic(err.Error())
	}
}

func GetNotifications(scrapers []Scraper) (notifications []Notification) {
	fmt.Println(Gray(8-1, "Starting GetNotifications..."))

	for _, elem := range scrapers {

		scraper_id := elem.Id

		rows, err := Db.Query(`
                            SELECT
                                m.id AS matchid,
                                utk.userid,
                                u.username,
                                u.email,
                                s.name,
                                r.title,
                                r.url
                            FROM userstargetskeywords utk
                            INNER JOIN matches m ON(utk.keywordid = m.keywordid)
                            INNER JOIN results r ON(m.resultid = r.id)
                            INNER JOIN scrapers s ON(r.scraperid = s.id)
                            INNER JOIN users u ON(utk.userid = u.id)
                            LEFT JOIN notifications n ON(m.id = n.matchid)
                            WHERE m.createdat > current_date - interval '0' day
                            AND s.targetid = utk.targetid
                            AND s.id = $1
                            AND n.id IS NULL
                            AND utk.deletedat IS NULL;`, scraper_id)
		counter := 0
		for rows.Next() {
			notification := Notification{}
			rows.Scan(
				&notification.MatchId,
				&notification.UserId,
				&notification.UserName,
				&notification.UserEmail,
				&notification.CompanyName,
				&notification.Title,
				&notification.Url)
			counter++
			notifications = append(notifications, notification)
		}
		rows.Close()

		if err != nil {
			panic(err.Error())
		}
	}
	return
}

func GetUserNotifications(scrapers []Scraper, user string) (notifications []Notification) {
	fmt.Println(Gray(8-1, "Starting GetUserNotifications..."))

	for _, elem := range scrapers {

		scraper_id := elem.Id

		rows, err := Db.Query(`
                            SELECT
                                m.id AS matchid,
                                utk.userid,
                                u.username,
                                u.email,
                                s.name,
                                r.title,
                                r.url
                            FROM userstargetskeywords utk
                            INNER JOIN matches m ON(utk.keywordid = m.keywordid)
                            INNER JOIN results r ON(m.resultid = r.id)
                            INNER JOIN scrapers s ON(r.scraperid = s.id)
                            INNER JOIN users u ON(utk.userid = u.id)
                            LEFT JOIN notifications n ON(m.id = n.matchid)
                            WHERE m.createdat > current_date - interval '0' day
                            AND s.targetid = utk.targetid
                            AND s.id = $1
                            AND utk.userid = $2
                            AND n.id IS NULL
                            AND utk.deletedat IS NULL;`, scraper_id, user)
		counter := 0
		for rows.Next() {
			notification := Notification{}
			rows.Scan(
				&notification.MatchId,
				&notification.UserId,
				&notification.UserName,
				&notification.UserEmail,
				&notification.CompanyName,
				&notification.Title,
				&notification.Url)
			counter++
			notifications = append(notifications, notification)
		}
		rows.Close()

		if err != nil {
			panic(err.Error())
		}
	}
	return
}

func (notifier *Notifier) StartNotifierSession(user_id int) {
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

func SaveUserNotifications(
	notifications []Notification, email Email, notifier Notifier) {
	fmt.Println(Gray(8-1, "Starting SaveUserNotifications..."))
	time_now := time.Now()
	valueStrings := []string{}
	valueArgs := []interface{}{}
	var counter int
	for _, notification := range notifications {
		if notification.UserId == email.UserId {
			str1 := "$" + strconv.Itoa(1+counter*3) + ","
			str2 := "$" + strconv.Itoa(2+counter*3) + ","
			str3 := "$" + strconv.Itoa(3+counter*3)
			str_n := "(" + str1 + str2 + str3 + ")"
			valueStrings = append(valueStrings, str_n)
			valueArgs = append(valueArgs, notifier.Id)
			valueArgs = append(valueArgs, notification.MatchId)
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
	return
}

func SaveEmail(email string, action string) {
	fmt.Println(Gray(8-1, "Starting SaveEmail..."))
	statement := `INSERT INTO sentemails (email, action, sentat)
                  VALUES ($1, $2, $3)`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	stmt.QueryRow(
		email,
		action,
		time.Now(),
	)
}
