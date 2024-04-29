package configs

import (
	"errors"
	"github.com/spf13/viper"
)

type Cfg struct {
	JwtExpiresIn int    `mapstructure:"JWT_EXPIRES_IN"`
	JWTSecret    string `mapstructure:"JWT_SECRET" validate:"required"`
	AppName      string `mapstructure:"APP_NAME"`
	DBDir        string `mapstructure:"DB_DIR"`
	Port         int    `mapstructure:"PORT"`
}

var Conf *Cfg

func LoadConfig() error {
	var cfg Cfg
	vip := viper.New()

	// Definindo valores padrão
	vip.SetDefault("PORT", 8000)
	vip.SetDefault("DB_DIR", "/tmp/rosedb")
	vip.SetDefault("JWT_EXPIRES_IN", 3600)
	vip.SetDefault("APP_NAME", "RaspController")

	// Lendo o arquivo de configuração conf.yml
	vip.SetConfigName("conf")
	vip.SetConfigType("yml")
	vip.AddConfigPath(".")

	// Lendo as configurações do arquivo conf.yml
	if err := vip.ReadInConfig(); err != nil {
		// Se o arquivo conf.yml não for encontrado, continue sem erro
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return err
		}
	}

	// Se JWT_SECRET não estiver definido, retorne um erro
	if !vip.IsSet("JWT_SECRET") {
		return errors.New("JWT_SECRET is not set")
	}

	// Atribua as configurações ao cfg
	if err := vip.Unmarshal(&cfg); err != nil {
		return err
	}

	// Atualiza a variável global Conf
	Conf = &cfg

	return nil
}
