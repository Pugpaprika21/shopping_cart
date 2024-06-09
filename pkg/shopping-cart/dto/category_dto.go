package dto

type CategoryReqBody struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CategoryFetchRow struct {
	ID                  uint   `gorm:"column:id"`
	CategoryName        string `gorm:"column:category_name"`
	CategoryDescription string `gorm:"column:category_description"`
}

type CategoryResp struct {
	ID                  uint   `json:"id"`
	CategoryName        string `json:"category_name"`
	CategoryDescription string `json:"category_description"`
}
