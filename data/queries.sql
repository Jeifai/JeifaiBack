/* Delete a whole table */
DROP TABLE author; -- 

/* Delete row from table */
DELETE FROM scrapers WHERE version = 2;

/* Delete column from table */
ALTER TABLE users_targets DROP COLUMN uuid;

/* Make a table empty and reset id */
TRUNCATE scrapers RESTART IDENTITY;

/* Update value in column based on condition */
UPDATE targets SET name = 'Kununu' WHERE id = 45;

/* Extract all the jobs by user */
SELECT t.url, j.created_at, j.title, j.url FROM users_targets ut
LEFT JOIN targets t ON(ut.target_id = t.id)
LEFT JOIN scrapers s ON(ut.target_id = s.target_id)
LEFT JOIN jobs j ON(s.id = j.scraper_id)
WHERE ut.user_id = 13;

ALTER TABLE jobs ADD COLUMN scraping_id integer references scraping(id);

/* New scraper process */
INSERT INTO targets (url, host, created_at, name) VALUES('https://www.babelforce.com/jobs/', 'https://www.babelforce.com', current_timestamp, 'Babelforce');
SELECT id FROM targets WHERE name = 'Babelforce';
INSERT INTO scrapers (name, version, target_id, created_at) VALUES('Babelforce', 1, 86, current_timestamp);