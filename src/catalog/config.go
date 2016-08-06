package catalog

import "fmt"

const (
	DefaultGrpcPort   = 18857
	DefaultHealthPort = 8080
)

type Config struct {
	Debug         bool
	GrpcPort      int
	HealthPort    int
	TlsCertFile   string
	TlsKeyFile    string
	ClientID      string
	ClientSecret  string
	AuthDiscovery string
}

func (c *Config) GetGrpcPortString() string {
	return fmt.Sprintf(":%d", c.GrpcPort)
}

func (c *Config) GetHealthPortString() string {
	return fmt.Sprintf(":%d", c.HealthPort)
}

func NewConfig() *Config {
	return &Config{GrpcPort: DefaultGrpcPort, HealthPort: DefaultHealthPort}
}
