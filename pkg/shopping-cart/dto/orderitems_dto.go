package dto

type OrderItemsReqBody struct {
	Products []OrderItemsProductReqBody `json:"products"`
}

type OrderItemsProductReqBody struct {
	ProductID       uint `json:"product_id"`
	ProductQuantity uint `json:"product_quantity"`
}

type OrderItemsFetchRow struct {
	ID              int     `gorm:"column:id"`
	ProductID       uint    `gorm:"column:product_id"`
	ProductQuantity uint    `gorm:"column:product_quantity"`
	UserID          uint    `gorm:"column:user_id"`
	ProductName     string  `gorm:"column:product_name"`
	ProductPrice    float32 `gorm:"column:product_price"`
	CategoryName    string  `gorm:"column:category_name"`
	OrderStatus     string  `gorm:"column:order_status"`
}

type OrderItemsResp struct {
	ID              int     `json:"id"`
	ProductID       uint    `json:"product_id,omitempty"`
	ProductQuantity uint    `json:"product_quantity"`
	UserID          uint    `json:"user_id,omitempty"`
	ProductName     string  `json:"product_name"`
	ProductPrice    float32 `json:"product_price"`
	CategoryName    string  `json:"category_name"`
	OrderStatus     string  `json:"order_status,omitempty"`
}

type OrderItemsRespBody struct {
	OrderItems  []OrderItemsResp `json:"order_items"`
	TotalTax    float32          `json:"total_tax"`
	Total       float32          `json:"total"`
	OrderStatus string           `json:"order_status"`
}

/*

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


*/
