UPDATE
    scrapings
SET
    countresults=subquery.countTotal
FROM (
    WITH
        results_per_day AS ( 
            SELECT
                t.createdat,
                t.countCreated,
                t.countClosed,
                sum(t.countCreated - t.countClosed) over (ORDER BY t.createdat) AS countTotal,
                ROW_NUMBER() OVER () AS rn
            FROM (
                WITH
                    jobs_created AS(
                        SELECT
                            r.createdat::date AS createdat,
                            COUNT(DISTINCT r.id) AS countCreated
                        FROM results r
                        LEFT JOIN scrapers s ON(r.scraperid = s.id)
                        WHERE s.name = 'Dreamingjobs'
                        GROUP BY 1),
                    jobs_closed AS(
                        SELECT
                            r.updatedat::date AS closedat,
                            COUNT(DISTINCT r.id) AS countClosed
                        FROM results r
                        LEFT JOIN scrapers s ON(r.scraperid = s.id)
                        WHERE s.name =  'Dreamingjobs'
                        GROUP BY 1),
                    consecutive_dates AS(
                            SELECT
                                date_trunc('day', dd)::date AS consdate
                            FROM generate_series(
                                (
                                    SELECT
                                        MIN(s.createdat)
                                    FROM scrapers s 
                                    WHERE s.name = 'Dreamingjobs'
                                ),
                                (
                                    SELECT 
                                        MAX(r.updatedat) - INTERVAL '1 DAY' 
                                    FROM scrapers s 
                                    LEFT JOIN results r ON(s.id = r.scraperid) 
                                    WHERE s.name = 'Dreamingjobs'
                                ), 
                                '1 day'::interval) dd)
                SELECT
                    cd.consdate AS createdat,
                    CASE WHEN jcr.countCreated IS NULL THEN 0 ELSE jcr.countCreated END AS countCreated,
                    CASE WHEN jcl.countClosed IS NULL THEN 0 ELSE jcl.countClosed END AS countClosed
                FROM consecutive_dates cd
                LEFT JOIN jobs_created jcr ON(cd.consdate = jcr.createdat)
                LEFT JOIN jobs_closed jcl ON(cd.consdate = jcl.closedat)) AS t)
    SELECT
        rpd.createdat,
        rpd.countTotal,
        x.id
    FROM results_per_day rpd
    LEFT JOIN (
        SELECT
            s.createdat::date AS createddate,
            s.id
        FROM scrapings s
        LEFT JOIN scrapers ss ON(s.scraperid = ss.id)
        WHERE ss.name = 'Dreamingjobs') AS x ON(rpd.createdat = x.createddate)) AS subquery
WHERE scrapings.id=subquery.id;
