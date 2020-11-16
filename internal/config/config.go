package config

import (
	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"
)

var Config appConfig

type appConfig struct {
	REDISPOOL    *redis.Pool
	REDISPOOLERR error
	RESTAPIPort  int    `mapstructure:"rest_api_port"`
	REDISURL     string `mapstructure:"redis_url"`
	DIAGPORT     int    `mapstructure:"diag_port"`
}

func LoadConfig(configPaths ...string) error {
	v := viper.New()
	v.SetConfigName("server")
	v.SetConfigType("yaml")
	v.SetEnvPrefix("restful")

	v.AutomaticEnv()

	//Config.RESTAPIPort = v.Get("8080").(string)
	//Config. = v.Get("API_KEY").(string)
	v.SetDefault("rest_api_port", 8080)
	v.SetDefault("diag_port", 8081)
	//v.SetDefault("db_url", "tcp://192.168.114.145:9000?debug=true")
	v.SetDefault("redis_url", "192.168.114.145:6379")

	//for _, path := range configPaths {
	//	v.AddConfigPath(path)
	//}
	//
	//err := v.ReadInConfig()
	//if err != nil {
	//	panic("Couldn't load configuration, cannot start. Terminating. Error: " + err.Error())
	//}
	//
	//log.Println("veper", v.AllKeys())
	//
	//if err := v.ReadInConfig(); err != nil {
	//	return fmt.Errorf("failed to read the configuration file: %s", err)
	//}

	return v.Unmarshal(&Config)
}
