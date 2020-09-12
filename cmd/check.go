package cmd

import (
	"os"

	"github.com/jedib0t/go-pretty/table"
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
										WITH
										    max_results_per_day AS(
										        SELECT
										            s.createdat::date,
										            s.scraperid,
										            MAX(s.countresults) AS countresults,
										            MAX(s.id) AS id
										        FROM scrapings s
										        GROUP BY 1, 2)
										SELECT
										    m.id,
										    m.scraperid,
										    m.countresults,
										    rank() OVER(
										        PARTITION BY m.scraperid
										        ORDER BY m.id DESC
										    ) AS rank
										FROM max_results_per_day m) t
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

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Name", "Three_day_ago", "One_day_ago", "delta_perc"})

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

		t.AppendRow([]interface{}{name, three_day_ago, one_day_ago, delta_perc})

	}

	t.Render()

	rows.Close()
}
