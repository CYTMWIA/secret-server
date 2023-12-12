package main

import (
	"fmt"
	"runtime/debug"
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

	Serve()
}