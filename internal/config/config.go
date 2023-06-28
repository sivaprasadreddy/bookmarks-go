package config

import (
	"github.com/spf13/viper"
	"log"
)

type AppConfig struct {
	Environment          string `mapstructure:"ENVIRONMENT"`
	ServerPort           int    `mapstructure:"SERVER_PORT"`
	DbHost               string `mapstructure:"DB_HOST"`
	DbPort               int    `mapstructure:"DB_PORT"`
	DbUserName           string `mapstructure:"DB_USERNAME"`
	DbPassword           string `mapstructure:"DB_PASSWORD"`
	DbDatabase           string `mapstructure:"DB_NAME"`
	DbRunMigrations      bool   `mapstructure:"DB_RUN_MIGRATIONS"`
	DbMigrationsLocation string `mapstructure:"DB_MIGRATIONS_LOCATION"`
}

func GetConfig(configFilePath string) (AppConfig, error) {
	log.Printf("Config File Path: %s\n", configFilePath)

	conf := viper.New()
	conf.SetConfigFile(configFilePath)
	//conf.SetEnvPrefix("APP")
	//replacer := strings.NewReplacer(".", "_")
	//conf.SetEnvKeyReplacer(replacer)
	conf.AutomaticEnv()

	err := conf.ReadInConfig()
	if err != nil {
		log.Printf("fatal error config file: %v\n", err)
		return AppConfig{}, err
	}
	var cfg AppConfig

	err = conf.Unmarshal(&cfg)
	if err != nil {
		log.Printf("configuration unmarshalling failed!. Error: %v\n", err)
		return cfg, err
	}
	//fmt.Printf("%#v\n", cfg)
	return cfg, nil
}
