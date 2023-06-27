package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type AppConfig struct {
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
	log.Infof("Config File Path: %s", configFilePath)

	conf := viper.New()
	conf.SetConfigFile(configFilePath)
	//conf.SetEnvPrefix("APP")
	//replacer := strings.NewReplacer(".", "_")
	//conf.SetEnvKeyReplacer(replacer)
	conf.AutomaticEnv()

	err := conf.ReadInConfig()
	if err != nil {
		log.Errorf("fatal error config file: %v", err)
		return AppConfig{}, err
	}
	var cfg AppConfig

	err = conf.Unmarshal(&cfg)
	if err != nil {
		log.Errorf("configuration unmarshalling failed!. Error: %v", err)
		return cfg, err
	}
	//fmt.Printf("%#v\n", cfg)
	return cfg, nil
}
