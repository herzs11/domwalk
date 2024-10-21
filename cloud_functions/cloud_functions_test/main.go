package main

import (
	"log"

	_ "dev.azure.com/Unum/Mkt_Analytics/_git/domwalk/cloud_functions"
	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
)

func main() {
	if err := funcframework.Start("8080"); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}
