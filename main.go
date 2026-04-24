package main

import (
	"fmt"
	"gator/internal/config"
)

func main() {
	// Read config
	// set user to "infernoe" and save to file
	// reread config and print
	fmt.Println("hello, World!")
	FromConfig := config.Read()
	fmt.Println(FromConfig)
}