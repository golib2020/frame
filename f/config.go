package f

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

//Config 配置
func Config() *viper.Viper {
	return Instance().GetOrSetFunc("config", func() interface{} {
		return configInit()
	}).(*viper.Viper)
}

func configInit() *viper.Viper {
	v := viper.New()
	v.AutomaticEnv() //自动加载环境变量

	//v.SetConfigType("json")
	v.AddConfigPath("./")
	env := v.GetString("app_env")
	if env == "" {
		v.SetConfigFile("config.json")
	} else {
		v.SetConfigFile(fmt.Sprintf("config.%s.json", env))
	}

	if err := v.ReadInConfig(); err != nil {
		log.Println(err)
	}
	return v
}
