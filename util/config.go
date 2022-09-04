package util

import "github.com/spf13/viper"

//Config store all configuration of all application
//The value are read by viper from config file or evironment variables
type Config struct {
	DBDrvier string `mapstructure:"DB_DRIVER"`
	DBSource string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	DBSoureTesting string `mapstructure:"DB_SOURCE_TEST"`
	DBLOCAL bool `mapstructure:"DB_LOCAL"`
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

