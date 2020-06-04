/* Make a table empty and reset id */
TRUNCATE scrapers RESTART IDENTITY;

/* Useful drop constraint and create index*/
ALTER TABLE results DROP CONSTRAINT resultspkey;
CREATE INDEX idxresultsurl ON results(url);

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
ALTER TABLE results ADD COLUMN updatedat timestamp NOT NULL DEFAULT currenttimestamp;
ALTER TABLE jobs ADD COLUMN scrapingid integer references scraping(id);

/* Update value in column based on condition */
UPDATE targets SET name = 'Kununu' WHERE id = 45;
UPDATE targets SET url = replace(url, 'https://', '')

/* New scraper process */
INSERT INTO targets (url, host, createdat, name) VALUES('https://www.babelforce.com/jobs/', 'https://www.babelforce.com', currenttimestamp, 'Babelforce');
SELECT id FROM targets WHERE name = 'Babelforce';
INSERT INTO scrapers (name, version, targetid, createdat) VALUES('Babelforce', 1, 86, currenttimestamp);


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
