WITH
	titles_words AS(
		SELECT 
			unnest(
				regexp_split_to_array(
					string_agg(
						trim(
							regexp_replace(
								regexp_replace(r.title, '[^a-zA-Z]', ' ', 'g'),
							'\s+', ' ', 'g')
						),
					' '),
				' ')
			) AS word
		FROM results r
		LEFT JOIN scrapers s ON(r.scraperid = s.id)
		WHERE s.name='Google')
SELECT
	word,
	COUNT(word)
FROM titles_words
WHERE word IS NOT NULL
GROUP BY 1
ORDER BY 2 DESC
LIMIT 20;
