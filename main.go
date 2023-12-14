package main

import (
	"fmt"
	"runtime/debug"

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

	err := config.Init()
	if err != nil {
		panic(err)
	}

	server.Serve()
}
