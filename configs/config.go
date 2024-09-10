package configs

import (
	"errors"

	"github.com/spf13/viper"
)

type Cfg struct {
	AuthToken  string `mapstructure:"AUTH_TOKEN" validate:"required"`
	AppName    string `mapstructure:"APP_NAME"`
	DBDir      string `mapstructure:"DB_DIR"`
	Port       int    `mapstructure:"PORT"`
	ShareDir   string `mapstructure:"SHARE_DIR"`
	TimeFormat string `mapstructure:"TIME_FORMAT"`
	TimeZone   string `mapstructure:"TIME_ZONE"`
}

var Conf *Cfg

func LoadConfig() error {
	var cfg Cfg
	vip := viper.New()

	// Setting default values
	vip.SetDefault("PORT", 8000)
	vip.SetDefault("DB_DIR", "/tmp/raspc")
	vip.SetDefault("APP_NAME", "RaspController")
	vip.SetDefault("TIME_FORMAT", "02-Jan-2006")
	vip.SetDefault("TIME_ZONE", "America/Sao_Paulo")

	// Reading the conf.yml configuration file
	vip.SetConfigName("conf")
	vip.SetConfigType("yml")
	vip.AddConfigPath(".")
	vip.AddConfigPath("/opt/raspc")
	vip.AddConfigPath("/etc/raspc")

	// Reading settings from the conf.yml file
	if err := vip.ReadInConfig(); err != nil {
		// Se o arquivo conf.yml não for encontrado, continue sem erro
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return err
		}
	}
	vip.AutomaticEnv()

	// If AUTH_TOKEN is not set, return an error
	if !vip.IsSet("AUTH_TOKEN") {
		return errors.New("AUTH_TOKEN is not set")
	}

	// Assign settings to cfg
	if err := vip.Unmarshal(&cfg); err != nil {
		return err
	}

	// Updates the global variable Conf
	Conf = &cfg

	return nil
}
