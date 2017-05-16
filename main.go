package main

import (
	"log"
	"os"

	web "github.com/micro/go-web"
)

var (
	apiKey string
	addr   string
)

func main() {
	service := web.NewService(
		web.Name("demo.translate"),
		web.Version("0.1"),
		web.Address(addr),
	)

	service.HandleFunc("/translate", translateHandler)

	if err := service.Init(); err != nil {
		log.Fatal(err)
	}

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	if apiKey = os.Getenv("API_KEY"); apiKey == "" {
		log.Fatal("Environment variable 'API_KEY' must be set")
	}

	// Set default address if not set
	if addr = os.Getenv("ADDR"); addr == "" {
		addr = ":8080"
	}
}
