package config

import (
	"fmt"
	"github.com/backend-timedoor/gtimekeeper-framework/container"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
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

	//_, path, name := c.defultConfiguration(config)

	c.vp = viper.New()

	c.vp.SetEnvPrefix("gtimekeeper")
	c.vp.AutomaticEnv()
	c.bindOsEnv()

	if config.Name != "" {
		if _, err := os.Stat(fmt.Sprintf("%s/%s", config.Path, config.Name)); err == nil {
			c.vp.SetConfigType("env")
			c.vp.AddConfigPath(config.Path)
			c.vp.SetConfigName(config.Name)

			if err := c.vp.ReadInConfig(); err != nil {
				log.Fatalf("error reading config file, %s", err)
				os.Exit(0)
			}
		}
	}

	container.Set(ContainerName, c)

	return c
}

func (c *Config) bindOsEnv() {
	for _, s := range os.Environ() {
		c.bind(s)
	}
}

func (c *Config) bind(env string) {
	a := strings.Split(env, "=")
	c.vp.SetDefault(a[0], a[1])
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
