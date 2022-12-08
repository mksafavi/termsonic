package main

import (
	"flag"
	"fmt"
	"os"

	"git.sixfoisneuf.fr/termsonic/src"
)

var (
	baseURL  = flag.String("url", "", "URL to your Subsonic server")
	username = flag.String("username", "", "Subsonic username")
	password = flag.String("password", "", "Subsonic password")
)

func main() {
	flag.Parse()

	cfg, err := src.LoadDefaultConfig()
	if err != nil {
		fmt.Printf("Could not start termsonic: %v", err)
		os.Exit(1)
	}

	if *baseURL != "" {
		cfg.BaseURL = *baseURL
	}

	if *username != "" {
		cfg.Username = *username
	}

	if *password != "" {
		cfg.Password = *password
	}

	src.Run(cfg)
}
