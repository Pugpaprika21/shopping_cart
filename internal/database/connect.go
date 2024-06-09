package database

import (
	"context"
	"fmt"
	"time"

	"github.com/Pugpaprika21/configs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// connect, err := database.NewOpenConnect().MongoGetConfig(db)
// connect, err := database.NewOpenConnect().Select("pgsql").GetConfig(db)
type connect struct {
	driver string
	dsn    string
	config *configs.DatabaseConfiguration
}

func NewOpenConnect() *connect {
	return &connect{}
}

func (o *connect) Select(driver string) *connect {
	o.driver = driver
	return o
}

func (o *connect) GetConfig(c *configs.DatabaseConfiguration) (*gorm.DB, error) {
	o.config = c
	o.drivers()

	db, err := gorm.Open(postgres.Open(o.dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, err
	}

	return db, nil
}

func (o *connect) drivers() {
	switch o.driver {
	case "mysql":
		o.dsn = ""
	case "pgsql":
		o.dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
			o.config.Host, o.config.Username, o.config.Password, o.config.DBname, o.config.Port, o.config.SSLMode, o.config.TimeZone)
	default:
		o.dsn = ""
	}
}

func (o *connect) MongoGetConfig(c *configs.DatabaseConfiguration) (*mongo.Client, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(c.MongoDNS).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("could not connect to MongoDB: %v", err)
	}

	return client, nil
}
