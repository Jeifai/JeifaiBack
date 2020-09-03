/* Export as CSV */
\COPY (SELECT * FROM results r LEFT JOIN scrapers s ON(r.scraperid = s.id) WHERE s.name='Soundcloud') TO '/home/robimalco/output.csv' WITH csv header

SELECT
    s.name AS company_name, 
    TO_CHAR(r.createdat, 'YYYY-MM-DD')AS job_created_at, 
    COUNT(r.id) AS new_jobs
FROM results r 
LEFT JOIN scrapers s ON(r.scraperid = s.id) 
WHERE s.name IN('Twitter', 'Google', 'Zalando') 
GROUP BY 1, 2 
ORDER BY 1, 2 DESC;

/* HERE A LIST OF EXAMPLE QUERIES */

/* Make a table empty and reset id */
TRUNCATE scrapers RESTART IDENTITY;
TRUNCATE userstargetskeywords, keywords RESTART IDENTITY;

/* Useful drop constraint and create index*/
ALTER TABLE results DROP CONSTRAINT resultspkey;
CREATE INDEX idxresultsurl ON results(url);

/* Add constraint */
ALTER TABLE targets ADD UNIQUE (name);
ALTER TABLE results ADD CONSTRAINT id_unique UNIQUE (id);
ALTER TABLE userstargetskeywords ADD UNIQUE (userid, targetid, keywordid);

/* Rename table's name */
ALTER TABLE userstargets RENAME TO userstargets;

/* Rename table's column */
ALTER TABLE users RENAME COLUMN name TO username;

/* Delete a whole table */
DROP TABLE author; -- 

/* Delete row from table */
DELETE FROM scrapers WHERE version = 2;

/* Delete column from table */
ALTER TABLE userstargets DROP COLUMN uuid;

/* Add column to table */
ALTER TABLE results ADD COLUMN updatedat timestamp NOT NULL DEFAULT current_timestamp;
ALTER TABLE userstargetskeywords ADD COLUMN updatedat timestamp NOT NULL DEFAULT current_timestamp;
ALTER TABLE jobs ADD COLUMN scrapingid integer references scraping(id);
ALTER TABLE targets ALTER COLUMN name SET NOT NULL;

/* Update value in column based on condition */
UPDATE targets SET name = 'Kununu' WHERE id = 45;
UPDATE targets SET url = replace(url, 'https://', '')

/* New scraper process */
INSERT INTO targets (url, host, createdat, name) VALUES('https://www.babelforce.com/jobs/', 'https://www.babelforce.com', currenttimestamp, 'Babelforce');
SELECT id FROM targets WHERE name = 'Babelforce';
INSERT INTO scrapers (name, version, targetid, createdat) VALUES('Babelforce', 1, 86, current_timestamp);


/* Get the latest two extractions */
SELECT s.id FROM scraping s LEFT JOIN scrapers ss ON(s.scraperid = ss.id) WHERE name = 'Zalando' ORDER BY s.id DESC LIMIT 2;
    /* OLD DATA */
SELECT id, createdat, updatedat, url, title FROM results WHERE scrapingid = 99;
    /* NEW DATA */
SELECT id, createdat, updatedat, url, title FROM results WHERE scrapingid = 114 AND DATE(createdat) = DATE(updatedat);

/* Example of query with Microsoft JSON */
SELECT
    r.id,
    r.title, 
    r.data#>>'{category}' AS category,
    r.data#>>'{country}' AS country
FROM results r
LEFT JOIN scrapers s ON(r.scraperid = s.id)
WHERE s.name = 'Microsoft'
ORDER BY r.updatedat DESC
LIMIT 10;

/* Update targets CASCADE */
INSERT INTO targets (url, host, createdat, name) VALUES('https://boards.greenhouse.io/urbansportsclub/', 'https://urbansportsclub.com', current_timestamp, 'Urbansport');
UPDATE targets SET id = 7 WHERE url = 'https://boards.greenhouse.io/urbansportsclub/';
UPDATE scrapers SET targetid = 7 WHERE targetid = 241;
DELETE FROM targets WHERE id = 241;

/* Check if column is unique */
SELECT CASE WHEN count(distinct id)= count(id)
THEN 'column values are unique' ELSE 'column values are NOT unique' END
FROM results;


/* Massive insert into */
INSERT INTO userskeywords(userid, keywordid, createdat)
SELECT DISTINCT
	utk.userid,
	utk.keywordid,
	MAX(utk.createdat)
FROM userstargetskeywords utk
WHERE utk.deletedat IS NULL
GROUP BY 1, 2
ORDER BY 1, 2, 3;