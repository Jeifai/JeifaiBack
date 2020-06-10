SELECT DISTINCT
    u.email,
	m.id,
	m.createdat,
	s.name,
	k.text AS keywordText,
	r.title
FROM matches m
LEFT JOIN keywords k ON(m.keywordid = k.id)
LEFT JOIN results r ON(m.resultid = r.id)
LEFT JOIN scrapers s ON(r.scraperid = s.id)
LEFT JOIN userstargetskeywords utk ON(utk.keywordid = k.id)
LEFT JOIN users u ON(utk.userid = u.id)
ORDER BY 1, 2;