package util

import (
	"github.com/spf13/viper"
	"log"
	"time"
)

var EnableRedis bool
var JWTConfig jwtConfigS

type jwtConfigS struct {
	Issuer string
	Secret string
	Expire time.Duration
}

func ReadConfig() (err error) {
	vp := viper.New()
	vp.SetConfigFile("config.yaml")

	err = vp.ReadInConfig()
	if err != nil {
		log.Println("Can not find config file. Generating default config.")
		EnableRedis = false
		JWTConfig = jwtConfigS{
			Issuer: "go-message-pusher",
			Secret: "GentleComet",
			Expire: 7200 * time.Second,
		}

		vp.Set("EnableRedis", false)
		vp.Set("JWT", JWTConfig)
		err = vp.SafeWriteConfigAs("config.yaml")
		return err
	}

	err = vp.UnmarshalKey("JWT", &JWTConfig)
	if err != nil {
		return err
	}
	JWTConfig.Expire *= time.Second

	EnableRedis = vp.GetBool("EnableRedis")

	return nil
}
