package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	EthJsonRpcURL string `mapstructure:"ETH_JSON_RPC_URL"`
}

var AppConfig *Config

func init() {
	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AutomaticEnv()
	failOnError(v.BindEnv("ETH_JSON_RPC_URL"), "fail on Bind ETH_JSON_RPC_URL")
	err := v.ReadInConfig()
	if err != nil {
		log.Println("load from environment variable")
	}
	err = v.Unmarshal(&AppConfig)
	if err != nil {
		failOnError(err, "Failed to read enivroment")
	}

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
