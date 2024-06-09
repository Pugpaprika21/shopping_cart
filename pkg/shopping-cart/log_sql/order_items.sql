SELECT 
    oi.id,
    oi.product_id,
    oi.product_quantity, 
    oi.user_id, 
    p.name AS product_name,
    p.price AS product_price,
    c.name AS category_name,
    oi.status AS order_status
FROM 
    order_items AS oi
INNER JOIN 
    products AS p ON oi.product_id = p.id
INNER JOIN 
    categories AS c ON c.id = p.category_id
WHERE 
	oi.status = 'pending' AND oi.user_id = '1'
ORDER BY 
     oi.id DESC;
