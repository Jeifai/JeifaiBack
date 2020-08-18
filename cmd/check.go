package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check that last scraping session made good results",
	Long: `This command checks the last scraping session and the previous one in order
    to establish if the number of results extracted do not differ too much`,
	Run: func(cmd *cobra.Command, args []string) {
		Check()
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}

func Check() {
	DbConnect()
	defer Db.Close()
	RunCheckQuery()
}

func RunCheckQuery() {
	rows, err := Db.Query(`WITH
                            pivot_data AS (
                                SELECT * FROM crosstab(
                                    'SELECT
                                        s.name::text AS name,
                                        t.rank::text,
                                        t.countresults AS countresults
                                    FROM (
                                        SELECT
                                            s.id,
                                            s.scraperid,
                                            s.countresults,
                                            rank() OVER(
                                                PARTITION BY s.scraperid
                                                ORDER BY s.id DESC
                                            ) AS rank
                                        FROM scrapings s) t
                                    LEFT JOIN scrapers s ON(t.scraperid = s.id)
                                    LEFT JOIN scrapings ss ON(t.id = ss.id)
                                    WHERE t.rank < 4
                                    ORDER BY 1, 2') AS x (name text, one_day_ago int, two_day_ago int, three_day_ago int))
                        SELECT
                            name,
                            CASE WHEN three_day_ago IS NULL THEN 0 ELSE three_day_ago END,
                            CASE WHEN one_day_ago IS NULL THEN 0 ELSE one_day_ago END,
                            CASE WHEN one_day_ago - three_day_ago IS NULL THEN 0 ELSE
                                to_char(100.0 * (one_day_ago - three_day_ago) / three_day_ago,'999D99')::float END
                        FROM pivot_data
                        ORDER BY 4 DESC;`)
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		var name string
		var three_day_ago int
		var one_day_ago int
		var delta_perc float32
		if err = rows.Scan(
			&name,
			&three_day_ago,
			&one_day_ago,
			&delta_perc); err != nil {
			panic(err.Error())
		}
		fmt.Println(name, three_day_ago, one_day_ago, delta_perc)
	}
	rows.Close()
}
