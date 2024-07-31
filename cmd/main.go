package main

import (
	"ep-golang-caching/configs"
	"fmt"
	"log"
)

func main() {
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	fmt.Printf("Database Host: %s\n", config.Database.Host)
	fmt.Printf("Cache Type: %s\n", config.Cache.Type)
	fmt.Printf("HTTP Address: %s\n", config.HTTP.Address)
}
