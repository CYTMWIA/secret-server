package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/CYTMWIA/secret-server/backend"
)

type Config struct {
	Mode string // "server" or "function"

	StorageBackend string // "local" or "s3"
	BackendLocal   backend.LocalBackendConfig
	BackendS3      backend.S3BackendConfig

	Addr string

	ApiKeyList []string
}

func (cfg *Config) Print() {
	fmt.Printf("%#v\n", *cfg)
}

func LoadConfig() (*Config, error) {
	paths := []string{
		"config.json",
		"config/config.json",
		"config/config-example.json", // for development
	}
	for _, path := range paths {
		fmt.Printf("Loading %s ", path)
		content, err := os.ReadFile(path)
		if err == nil {
			cfg, err := load_config(content)
			if err == nil {
				fmt.Println("OK")
				return cfg, nil
			}
		}
		fmt.Println("FAIL")
	}
	return nil, errors.New("unable to load config")
}

func load_config(data []byte) (*Config, error) {
	var config Config
	err := json.Unmarshal(data, &config)
	return &config, err
}
