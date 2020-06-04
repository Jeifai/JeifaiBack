SELECT 
    s.name, 
    COUNT(r.id) 
FROM results r 
LEFT JOIN scrapers s ON(r.scraperid = s.id) 
GROUP BY 1 
ORDER BY 2 DESC;