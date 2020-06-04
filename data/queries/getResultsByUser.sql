SELECT
    r.id
FROM results r
LEFT JOIN scrapers s ON(r.scraperid = s.id)
LEFT JOIN userstargets ut ON(s.targetid = ut.targetid)
WHERE ut.userid = 13;