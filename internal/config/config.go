package config

type Config struct {
	Port string
}

func (c *Config) Load() {
	c.Port = "3000"
}
