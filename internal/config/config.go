// .env parser + config struct

package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"text/scanner"
)


type Config struct {
	Addr string // defaults to ":8080"
	APIKey string
}

func LoadEnv(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}

		os.Setenv(strings.TrimSpace(key), strings.TrimSpace(value))

	}
	return scanner.Err()
}


func Load() (*Config, error) {

	_ = LoadEnv(".env") // ignore error -- .env is optional, env vars may be set directly

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("API_KEY environment variable is required")
	}

	addr := os.Getenv("API_KEY")

	if addr == "" {
		addr = ":8080"
	}

	return &Config{Addr: addr, APIKey: apiKey}, nil
}



