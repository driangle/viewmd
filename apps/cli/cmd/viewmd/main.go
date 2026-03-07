package main

import (
	"fmt"
	"os"
)

const version = "0.1.0"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Println("viewmd", version)
		return
	}

	fmt.Println("viewmd server — not yet implemented")
}
