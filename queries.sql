/* Useful drop constraint and create index*/
ALTER TABLE results DROP CONSTRAINT results_pkey;
CREATE INDEX idx_results_url ON results(url);

/* Rename table's name */
ALTER TABLE users_targets RENAME TO userstargets;

/* Rename column's */
ALTER TABLE users RENAME COLUMN name TO user_name;

/* Delete a whole table */
DROP TABLE author; -- 

/* Delete row from table */
DELETE FROM scrapers WHERE version = 2;

/* Delete column from table */
ALTER TABLE users_targets DROP COLUMN uuid;

/* Add column to table */
ALTER TABLE results ADD COLUMN updated_at timestamp NOT NULL DEFAULT current_timestamp;
ALTER TABLE jobs ADD COLUMN scraping_id integer references scraping(id);

/* Make a table empty and reset id */
TRUNCATE scrapers RESTART IDENTITY;

/* Update value in column based on condition */
UPDATE targets SET name = 'Kununu' WHERE id = 45;
UPDATE targets SET url = replace(url, 'https://', '')

/* Extract all the jobs by user */
SELECT t.url, j.created_at, j.title, j.url FROM users_targets ut
LEFT JOIN targets t ON(ut.target_id = t.id)
LEFT JOIN scrapers s ON(ut.target_id = s.target_id)
LEFT JOIN jobs j ON(s.id = j.scraper_id)
WHERE ut.user_id = 13;

/* New scraper process */
INSERT INTO targets (url, host, created_at, name) VALUES('https://www.babelforce.com/jobs/', 'https://www.babelforce.com', current_timestamp, 'Babelforce');
SELECT id FROM targets WHERE name = 'Babelforce';
INSERT INTO scrapers (name, version, target_id, created_at) VALUES('Babelforce', 1, 86, current_timestamp);

/* Count results by last scraping_id*/
SELECT
     s.id,
     COUNT(r.id)
FROM scraping s
LEFT JOIN results r ON(s.id = r.scraping_id)
GROUP BY 1 
ORDER BY s.id
DESC LIMIT 5;

/* Count results by scrapers name*/
SELECT
     ss.name,
     COUNT(r.id)
FROM scraping s
LEFT JOIN scrapers ss ON(s.scraper_id = ss.id)
LEFT JOIN results r ON(s.id = r.scraping_id)
GROUP BY 1 
ORDER BY 2 DESC;



/* Get the latest two extractions */
SELECT s.id FROM scraping s LEFT JOIN scrapers ss ON(s.scraper_id = ss.id) WHERE name = 'Zalando' ORDER BY s.id DESC LIMIT 2;
    /* OLD DATA */
SELECT id, created_at, updated_at, url, title FROM results WHERE scraping_id = 99;
    /* NEW DATA */
SELECT id, created_at, updated_at, url, title FROM results WHERE scraping_id = 114 AND DATE(created_at) = DATE(updated_at);


/* Count results based on scraper name */
SELECT 
    s.name, 
    COUNT(r.id) 
FROM results r 
LEFT JOIN scrapers s ON(r.scraper_id = s.id) 
GROUP BY 1 
ORDER BY 2 DESC;


/* Example of query with Microsoft JSON */
SELECT
    r.id,
    r.title, 
    r.data#>>'{category}' AS category,
    r.data#>>'{country}' AS country
FROM results r
LEFT JOIN scrapers s ON(r.scraper_id = s.id)
WHERE s.name = 'Microsoft'
ORDER BY r.updated_at DESC
LIMIT 10;