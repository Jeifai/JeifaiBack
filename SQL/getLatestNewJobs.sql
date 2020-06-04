SELECT
    r.id,
    r.createdat,
    ss.name,
    r.title,
    r.url
FROM results r
LEFT JOIN scrapings s ON(r.scraperid = s.id)
LEFT JOIN scrapers ss ON(r.scraperid = ss.id)
WHERE r.createdat = r.updatedat
ORDER BY r.createdat DESC
LIMIT 30;