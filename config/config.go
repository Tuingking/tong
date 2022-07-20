package config

import (
	"github.com/spf13/viper"
	"github.com/tuingking/tong/pkg/gcs"
	"github.com/tuingking/tong/pkg/kafka"
	"github.com/tuingking/tong/pkg/mysql"
)

type Config struct {
	Version string
	Env     string
	Mysql   mysql.Option
	Kafka   kafka.Option
	GCS     gcs.Option
}

type option struct {
	configPath string
	configFile string
	configType string
}

func defaultOption() option {
	return option{
		configPath: "$HOME/.tong/",
		configFile: "config",
		configType: "yaml",
	}
}

func Init() Config {
	opts := []option{defaultOption(), {"./config/", "config", "yaml"}}

	vip := viper.New()
	for _, opt := range opts {
		vip.AddConfigPath(opt.configPath)
		vip.SetConfigName(opt.configFile)
		vip.SetConfigType(opt.configType)

		if err := vip.ReadInConfig(); err != nil {
			continue
		}
	}

	var config Config
	if err := vip.Unmarshal(&config); err != nil {
		panic(err)
	}

	// log.Printf("[config] config used: %v\n", vip.ConfigFileUsed())

	return config
}
