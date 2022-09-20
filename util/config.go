package util

import (
	"time"

	"github.com/spf13/viper"
)

//Config store all configuration of all application
//The value are read by viper from config file or evironment variables
type Config struct {
	DBDrvier string `mapstructure:"DB_DRIVER"`
	DBSource string `mapstructure:"DB_SOURCE"`
	HttpServerAddress string `mapstructure:"HTTP_SERVER_ADDRESS"`
	GrpcServerAddress string `mapstructure:"GRPC_SERVER_ADDRESS"`
	DBSoureTesting string `mapstructure:"DB_SOURCE_TEST"`
	DBLOCAL bool `mapstructure:"DB_LOCAL"`
	TokenSymmetricKey string `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCES_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env") //json or the other

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

