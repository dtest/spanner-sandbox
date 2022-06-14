package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Configurations exported
type Config struct {
	Server  ServerConfig
	Spanner SpannerConfig
}

// ServerConfigurations exported
type ServerConfig struct {
	Host string
	Port int
}

// DatabaseConfigurations exported
type SpannerConfig struct {
	Project_id      string
	Instance_id     string
	Database_id     string
	CredentialsFile string
}

func NewConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err.Error())
	}

	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.port", 8083)

	var c Config

	err := viper.Unmarshal(&c)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}

	return &c, nil
}

func (c *SpannerConfig) DB() string {
	return fmt.Sprintf(
		"projects/%s/instances/%s/databases/%s",
		c.Project_id,
		c.Instance_id,
		c.Database_id,
	)
}

func (c *ServerConfig) URL() string {
	return fmt.Sprintf(
		"%s:%d",
		c.Host,
		c.Port,
	)
}
