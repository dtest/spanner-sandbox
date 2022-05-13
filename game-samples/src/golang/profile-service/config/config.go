package config

import "fmt"

// Configurations exported
type Configurations struct {
	Server  ServerConfigurations
	Spanner SpannerConfigurations
}

// ServerConfigurations exported
type ServerConfigurations struct {
	Host string
	Port int
}

// DatabaseConfigurations exported
type SpannerConfigurations struct {
	Project_id      string
	Instance_id     string
	Database_id     string
	CredentialsFile string
}

func (c *SpannerConfigurations) URL() string {
	return fmt.Sprintf(
		"projects/%s/instances/%s/databases/%s",
		c.Project_id,
		c.Instance_id,
		c.Database_id,
	)
}

func (c *ServerConfigurations) URL() string {
	return fmt.Sprintf(
		"%s:%d",
		c.Host,
		c.Port,
	)
}
