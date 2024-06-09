package server

import (
	"github.com/Pugpaprika21/configs"
	"github.com/Pugpaprika21/internal/database"
	"github.com/Pugpaprika21/pkg/shopping-cart/migrations"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type (
	echoServer struct {
		Echo   *echo.Echo
		Server *EchoServerEnvironment
	}

	EchoServerEnvironment struct {
		Env     *viper.Viper
		App     *configs.ApplicationConfiguration
		Connect *EchoDB
	}

	EchoDB struct {
		ORM   *gorm.DB      `orm:"omitempty"`
		Mongo *mongo.Client `mongo:"omitempty"`
	}
)

func NewEchoServer() *echoServer {
	config := configs.NewConfiguration()
	env := config.LoadEnvironment("config.yaml")
	app := config.LoadAppConfig()
	db := config.LoadDBConfig()

	orm, err := database.NewOpenConnect().Select("pgsql").GetConfig(db)
	if err != nil {
		panic(err.Error())
	}

	err = migrations.NewMigration().GetConnect(orm).Run()
	if err != nil {
		panic(err.Error())
	}

	return &echoServer{
		Echo: echo.New(),
		Server: &EchoServerEnvironment{
			Env: env,
			App: app,
			Connect: &EchoDB{
				ORM: orm,
			},
		},
	}
}
