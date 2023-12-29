package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

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

func secret_text(origin string, display_length int) string {
	if display_length > len(origin) {
		return origin
	} else {
		return origin[:display_length] + "******"
	}
}

func (cfg *Config) Print() {
	var be_printed = *cfg
	be_printed.BackendS3.Id = secret_text(be_printed.BackendS3.Id, 4)
	be_printed.BackendS3.Secret = secret_text(be_printed.BackendS3.Secret, 4)

	bytes, err := json.MarshalIndent(be_printed, "", "  ")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(bytes))
	}
}

func LoadConfig() (*Config, error) {
	var cfg Config

	load_config_file(&cfg)

	load_config_env(&cfg)

	return &cfg, nil
}

func load_config_file(cfg *Config) error {
	paths := []string{
		"config.json",
		"config/config.json",
		"config/config-example.json", // for development
	}
	for _, path := range paths {
		fmt.Printf("Loading %s ", path)
		content, err := os.ReadFile(path)
		if err == nil {
			err := json.Unmarshal(content, cfg)
			if err == nil {
				fmt.Println("OK")
				return nil
			}
		}
		fmt.Println("FAIL")
	}
	return errors.New("unable to load config")
}

func load_config_env(cfg *Config) {
	load_config_env_into_object(reflect.ValueOf(cfg), "SECRET_SERVER")
}

func load_config_env_into_object(obj reflect.Value, env_name string) {
	switch obj.Kind() {
	case reflect.Struct:
		for i := 0; i < obj.NumField(); i++ {
			field_name := obj.Type().Field(i).Name
			new_env_name := env_name + "_" + strings.ToUpper(field_name)
			load_config_env_into_object(obj.Field(i), new_env_name)
		}
	case reflect.String:
		env_value := os.Getenv(env_name)
		if env_value != "" && obj.CanSet() {
			fmt.Println("Read env:", env_name)
			obj.SetString(env_value)
		}
	case reflect.Pointer:
		load_config_env_into_object(obj.Elem(), env_name)
	default:
		return
	}
}
