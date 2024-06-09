package migrations

import (
	"log"

	"github.com/Pugpaprika21/pkg/shopping-cart/models"
	"gorm.io/gorm"
)

type migration struct {
	query *gorm.DB
}

func NewMigration() *migration {
	return &migration{}
}

func (m *migration) GetConnect(connect *gorm.DB) *migration {
	m.query = connect
	return m
}

func (m *migration) Run(dst ...interface{}) error {
	if len(dst) == 0 {
		dst = m.runbatch()
	}

	for _, model := range dst {
		if err := m.query.AutoMigrate(model); err != nil {
			log.Printf("Failed to migrate model: %v, error: %v\n", model, err)
			return err
		}
	}
	return nil
}

// register models here!!
func (m *migration) runbatch() []interface{} {
	return []interface{}{
		&models.Product{},
		&models.ProductAttachment{},
		&models.Category{},
		&models.OrderItems{},
	}
}
