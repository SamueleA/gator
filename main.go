package main

import (
	"fmt"
	"os"

	"github.com/samuelea/gator/internal/config"
)

func main() {
	println("Hello, World!")
	gatorConfig, err := config.Read()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading config:%v\n", err)
		os.Exit(1)
	}

	fmt.Println(gatorConfig)

	gatorConfig.SetUser("Sam")


	fmt.Println(gatorConfig)

}