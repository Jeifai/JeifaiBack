package data

import (
	"fmt"
)

type Result struct {
	Target    string
	Id        int
	ScraperId int
	Title     string
	Url       string
	CreatedAt string
}

// Get all the results belonging to the targets of a specific user
func (user *User) ResultsByUser() (results []Result, err error) {
	fmt.Println("Starting ResultsByUser...")
	rows, err := Db.Query(`SELECT DISTINCT
                                s.name,
                                TO_CHAR(m.createdat, 'DD/MM/YYYY'),
                                r.title,
                                r.url
                            FROM matches m
                            INNER JOIN keywords k ON(m.keywordid = k.id)
                            INNER JOIN results r ON(m.resultid = r.id)
                            INNER JOIN scrapers s ON(r.scraperid = s.id)
                            INNER JOIN userstargetskeywords utk ON(k.id = utk.keywordid)
                            WHERE m.createdat > current_date - interval '3' day
                            AND utk.userid = $1;`, user.Id)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for rows.Next() {
		result := Result{}
		if err = rows.Scan(&result.Target, &result.CreatedAt, &result.Title, &result.Url); err != nil {
			return
		}
		results = append(results, result)
	}
	rows.Close()
	return
}
