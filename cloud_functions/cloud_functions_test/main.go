package main

import (
	"log"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	_ "github.com/herzs11/domwalk/cloud_functions"
)

func main() {
	if err := funcframework.Start("8080"); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}
