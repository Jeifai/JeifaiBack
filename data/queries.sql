/* Delete a whole table */
DROP TABLE author; -- 

/* Make a table empty and reset id */
TRUNCATE scrapers RESTART IDENTITY;

/* Update value in column based on condition */
UPDATE targets SET name = 'Kununu' WHERE id = 45;

/* Insert in table */
INSERT INTO scrapers (name, version, target_id, created_at) VALUES('Kununu', 1, current_timestamp);

/* Extract all the jobs by user */
SELECT t.url, j.created_at, j.title, j.url FROM users_targets ut
LEFT JOIN targets t ON(ut.target_id = t.id)
LEFT JOIN scrapers s ON(ut.target_id = s.target_id)
LEFT JOIN jobs j ON(s.id = j.scraper_id)
WHERE ut.user_id = 13;

ALTER TABLE jobs ADD COLUMN scraping_id integer references scraping(id);