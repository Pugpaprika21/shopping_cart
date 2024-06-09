package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name       string   `gorm:"type:varchar(150);name"`
	Price      float32  `gorm:"price"`
	Qty        uint     `gorm:"qty"`
	Detail     string   `gorm:"detail"`
	CategoryID uint     `gorm:"index"`                         // ForeignKey
	Category   Category `gorm:"constraint:OnDelete:SET NULL;"` // เปลี่ยนเป็น OnDelete:SET NULL หรือเปลี่ยนเป็น OnDelete:RESTRICT ตามที่ต้องการ
}

type ProductAttachment struct {
	gorm.Model
	Filename  string  `gorm:"type:varchar(255);"`
	Path      string  `gorm:"type:text;"`
	Size      int64   `gorm:"type:bigint;"`
	MimeType  string  `gorm:"type:varchar(100);"`
	Extension string  `gorm:"type:varchar(10);"`
	ProductID uint    `gorm:"index"`                        // ForeignKey
	Product   Product `gorm:"constraint:OnDelete:CASCADE;"` // สร้างความสัมพันธ์กับ Product
}
