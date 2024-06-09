SELECT
	p.id,
	p.name AS product_name,
	p.price AS product_price,
	p.qty AS product_quantity,
	p.detail AS product_detail,
	p.category_id,
	c.name AS category_name,
	CONCAT(pa.filename, pa.extension) AS product_attachment
FROM
	products AS p
INNER JOIN product_attachments AS pa ON p.id = pa.product_id
INNER JOIN categories AS c ON p.category_id = c.id
ORDER BY
	pa.id DESC;