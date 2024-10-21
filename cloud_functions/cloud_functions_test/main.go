package main

import (
	"log"

	_ "cloud_functions"
	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
)

func main() {
	if err := funcframework.Start("8080"); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}
