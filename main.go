package main

import (
	"fmt"
	"gator/internal/config"
)

func main() {
	cfg := config.ReadConfig()
	cfg.SetUser("MrInfernoe")
	cfg = config.ReadConfig()
	fmt.Printf("config struct: %v\n", cfg)
}