package configuration

import (
	"os"

	"github.com/spf13/viper"
)

type ApplicationConfig struct {
	vip *viper.Viper
}

func BootConfig() *ApplicationConfig {
	app := &ApplicationConfig{}
	app.vip = viper.New()
	app.vip.SetConfigType("env")
	app.vip.SetConfigFile(".env")

	if err := app.vip.ReadInConfig(); err != nil {
		os.Exit(0)
	}

	app.vip.SetEnvPrefix("microdemy")
	app.vip.AutomaticEnv()

	return app
}