/* Delete a whole table */
DROP TABLE author; -- 

/* Make a table empty and reset id */
TRUNCATE scrapers RESTART IDENTITY;

/* Update value in column based on condition */
UPDATE targets SET name = 'Kununu' WHERE id = 45;

/* Insert in table */
INSERT INTO scrapers (name, version, created_at) VALUES('Kununu', 1, current_timestamp);