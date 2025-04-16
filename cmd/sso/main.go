package main

import (
	"fmt"
	"love-signal-sso/internal/config"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println("Hello from sso app!")
	fmt.Printf("with config: %v", cfg)
}
