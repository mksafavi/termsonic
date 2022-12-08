package main

import (
	"flag"
	"fmt"
	"os"

	"git.sixfoisneuf.fr/termsonic/src"
)

var (
	configFile = flag.String("config", "", "Path to the configuration file")
)

func main() {
	flag.Parse()

	var cfg *src.Config
	var err error
	if *configFile == "" {
		cfg, err = src.LoadDefaultConfig()
		if err != nil {
			fmt.Printf("Could not start termsonic: %v", err)
			os.Exit(1)
		}
	} else {
		f, err := os.Open(*configFile)
		if err != nil {
			fmt.Printf("Could not read configuration file: %v", err)
			os.Exit(1)
		}
		f.Close()

		cfg, err = src.LoadConfigFromFile(*configFile)
		if err != nil {
			fmt.Printf("Error loading configuration file: %v", err)
			os.Exit(1)
		}
	}

	src.Run(cfg)
}
