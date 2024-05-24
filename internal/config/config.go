package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
}

func (c *Config) Load() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file", err)
	}
	c.Port = "3000"
}
