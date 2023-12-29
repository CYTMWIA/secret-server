package main

import (
	"fmt"
	"runtime/debug"

	"github.com/CYTMWIA/secret-server/backend"
	"github.com/CYTMWIA/secret-server/config"
	"github.com/CYTMWIA/secret-server/server"
)

func print_build_info() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		fmt.Println("Unable to read build info.")
	} else {
		fmt.Print(info)
	}
	fmt.Println("====================")
}

func main() {
	print_build_info()

	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	cfg.Print()

	var backend_ backend.StorageBackend
	switch cfg.StorageBackend {
	case "local":
		backend_ = backend.NewLocalBackend(&cfg.BackendLocal)
	case "s3":
		backend_ = backend.NewS3Backend(&cfg.BackendS3)
	default:
		fmt.Println("Unkown backend", cfg.StorageBackend)
		return
	}

	err = server.Serve(cfg.Addr, cfg.Mode, backend_, cfg.ApiKeyList)
	if err != nil {
		fmt.Println(err)
		return
	}
}
