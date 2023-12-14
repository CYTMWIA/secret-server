package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/CYTMWIA/secret-server/crypto"
)

type Config struct {
	Mode           string // "server" or "tencent_sfc" (not implemented yet)
	StorageBackend string // "local" or "tencent_cos" (not implemented yet)
	Addr           string
	ApiKeyList     []string
}

var CONFIG Config = Config{
	Mode:           "server",
	StorageBackend: "local",
	Addr:           "0.0.0.0:8080",
	ApiKeyList:     []string{},
}

func Init(){
	cfg, err := LoadConfig()
	if err != nil {
		fmt.Println(err)
	} else {
		CONFIG = *cfg
	}

	fmt.Printf("Using Config: %#v", CONFIG)
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

func IsVaildUser(api_key string) (bool, error) {
	hkey := crypto.Hash(api_key)
	for _, key := range CONFIG.ApiKeyList {
		if key == hkey {
			return true, nil
		}
	}

	return false, nil
}
