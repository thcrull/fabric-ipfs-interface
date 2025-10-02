package main

import (
	"fmt"
	"log"

	"github.com/thcrull/fabric-interface/application/pkg/blockchain"
	"github.com/thcrull/fabric-interface/application/pkg/config"
)

func main() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	client, err := blockchain.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create Fabric client: %v", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Printf("Error closing client: %v", err)
		}
	}()

	submitResult, err := client.SubmitTransaction("CreateAsset", "asset2", "green", "500", "Tombus", "3000")
	if err != nil {
		log.Fatalf("Failed to submit transaction: %v", err)
	}
	fmt.Println("Submit result:", string(submitResult))

	queryResult, err := client.EvaluateTransaction("GetAllAssets")
	if err != nil {
		log.Fatalf("Failed to evaluate transaction: %v", err)
	}
	fmt.Println("Query result:", string(queryResult))
}
