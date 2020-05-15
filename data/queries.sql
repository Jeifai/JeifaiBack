/* Delete a whole table */
DROP TABLE author; -- 

/* Make a table empty and reset id */
TRUNCATE scrapers RESTART IDENTITY;

/* Update value in column based on condition */
UPDATE targets SET name = 'Kununu' WHERE id = 45;

/* Insert in table */
INSERT INTO scrapers (name, created_at) VALUES('Kununu', current_timestamp);


INSERT INTO targets_scrapers (scraper_id, target_id, created_at) VALUES(1, 45, current_timestamp);