package config

import (
	"log"
	"os"

	"github.com/backend-timedoor/gtimekeeper-framework/container"
	"github.com/spf13/viper"
)

const ContainerName string = "config"

type Config struct {
	vp *viper.Viper
}

type Configuration struct {
	Path string
	Name string
}

func New(config *Configuration) *Config {
	c := &Config{}

	_, path, name := c.defultConfiguration(config)

	c.vp = viper.New()
	c.vp.SetConfigType("env")
	c.vp.AddConfigPath(path)
	c.vp.SetConfigName(name)

	if err := c.vp.ReadInConfig(); err != nil {
		log.Fatalf("error reading config file, %s", err)
		os.Exit(0)
	}

	c.vp.SetEnvPrefix("gtimekeeper")
	c.vp.AutomaticEnv()

	container.Set(ContainerName, c)

	return c
}

func (c *Config) defultConfiguration(config *Configuration) (string, string, string) {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("error getting working directory env, %s", err)
		os.Exit(0)
	}

	path := config.Path
	name := config.Name

	if path == "" {
		path = pwd
	}

	if name == "" {
		name = ".env"
	}

	return pwd, path, name

}
