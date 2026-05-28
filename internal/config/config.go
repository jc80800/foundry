package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port int
}

func Load() (Config, error) {
	port := 8080
	if v := os.Getenv("PORT"); v != "" {
		p, err := strconv.Atoi(v)
		if err != nil {
			return Config{}, fmt.Errorf("invalid PORT %q: %w", v, err)
		}
		port = p
	}
	return Config{Port: port}, nil
}
