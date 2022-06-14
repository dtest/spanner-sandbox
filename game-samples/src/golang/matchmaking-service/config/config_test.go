package config

import (
	"bytes"
	"regexp"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func readConfig(yml []byte) (*Config, error) {
	viper.SetConfigType("yaml")

	// Read the config
	err := viper.ReadConfig(bytes.NewBuffer(yml))
	if err != nil {
		return &Config{}, err
	}

	var c Config
	err = viper.Unmarshal(&c)
	if err != nil {
		return &Config{}, err
	}

	return &c, nil
}

func TestServerURL(t *testing.T) {
	c, err := NewConfig()
	assert.Nil(t, err)

	assert.Regexp(t, regexp.MustCompile(`^[A-Za-z0-9.]*:8081$`), c.Server.URL())
}

func TestSpannerDB(t *testing.T) {
	cfgExample := []byte(`
server:
  host: localhost
  port: 8081
spanner:
  project_id: test-123
  instance_id: game-test-1
  database_id: game-db-1
`)

	c, err := readConfig(cfgExample)
	assert.Nil(t, err)

	assert.Regexp(t, regexp.MustCompile(`^projects/[a-z0-9-]*/instances/[a-z0-9-]*/databases/[a-z0-9-]*$`), c.Spanner.DB())
}
