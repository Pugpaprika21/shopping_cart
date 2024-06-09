package configs

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type (
	Configuration struct {
		configFilePath string
		viper          *viper.Viper
		app            ApplicationConfiguration
		db             DatabaseConfiguration
	}

	ApplicationConfiguration struct {
		Port    string
		BaseURL string
	}

	DatabaseConfiguration struct {
		Gorm     *gorm.DB
		Driver   string
		Host     string
		Username string
		Password string
		DBname   string
		Port     string
		SSLMode  string
		TimeZone string
		MongoDNS string
	}
)

func NewConfiguration() *Configuration {
	return &Configuration{
		configFilePath: "./configs/",
		viper:          viper.New(),
	}
}

func (c *Configuration) LoadEnvironment(configType string) *viper.Viper {
	return c.loadConfigEnvironment(configType)
}

func (c *Configuration) LoadAppConfig() *ApplicationConfiguration {
	return &c.app
}

func (c *Configuration) LoadDBConfig() *DatabaseConfiguration {
	return &c.db
}

func (c *Configuration) loadConfigEnvironment(configType string) *viper.Viper {
	c.configFilePath = c.configFilePath + configType

	if err := c.checkFileExists(); err != nil {
		log.Fatal(err.Error())
	}

	c.readConfigFile()
	c.setEnvVariables()

	return c.viper
}

func (c *Configuration) checkFileExists() error {
	if _, err := os.Stat(c.configFilePath); os.IsNotExist(err) {
		return err
	}
	return nil
}

func (c *Configuration) readConfigFile() {
	c.viper.SetConfigFile(c.configFilePath)
	c.viper.SetConfigType("yaml")

	if err := c.viper.ReadInConfig(); err != nil {
		log.Fatal(err.Error())
	}
}

func (c *Configuration) setEnvVariables() {
	c.viper.AutomaticEnv()
	c.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	c.app.Port = c.viper.GetString("port")
	c.app.BaseURL = c.viper.GetString("base_url")

	c.db.Driver = c.viper.GetString("database.driver")
	c.db.Host = c.viper.GetString("database.host")
	c.db.Username = c.viper.GetString("database.username")
	c.db.Password = c.viper.GetString("database.password")
	c.db.DBname = c.viper.GetString("database.name")
	c.db.Port = c.viper.GetString("database.port")
	c.db.SSLMode = c.viper.GetString("database.sslmode")
	c.db.TimeZone = c.viper.GetString("database.timezone")

	c.db.MongoDNS = c.viper.GetString("mongo.dns")
}
