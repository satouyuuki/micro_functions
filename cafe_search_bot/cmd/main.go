package main

import (
	"log"
	"os"

	_ "github.com/satouyuuki/micro_functions/cafe_search_bot"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
)

func main() {

	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}
	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}
