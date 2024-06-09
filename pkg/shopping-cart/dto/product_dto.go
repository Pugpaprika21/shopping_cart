package dto

type ProductReqBody struct {
	Name       string  `form:"name"`
	Price      float32 `form:"price"`
	Qty        uint    `form:"qty"`
	Detail     string  `form:"detail"`
	CategoryID uint    `form:"category_id"`
}

type ProductFetchRow struct {
	ID                int     `gorm:"column:id"`
	ProductName       string  `gorm:"column:product_name"`
	ProductPrice      float64 `gorm:"column:product_price"`
	ProductQuantity   int     `gorm:"column:product_quantity"`
	ProductDetail     string  `gorm:"column:product_detail"`
	CategoryID        uint    `gorm:"column:category_id"`
	CategoryName      string  `gorm:"column:category_name"`
	ProductAttachment string  `gorm:"column:product_attachment"`
}

type ProductResp struct {
	ID                int     `json:"id"`
	ProductName       string  `json:"product_name"`
	ProductPrice      float64 `json:"product_price"`
	ProductQuantity   int     `json:"product_quantity"`
	ProductDetail     string  `json:"product_detail"`
	CategoryID        uint    `json:"category_id,omitempty"`
	CategoryName      string  `json:"category_name"`
	ProductAttachment string  `json:"product_attachment"`
}
