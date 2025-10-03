package main

import (
	"fmt"
	"log"

	"github.com/thcrull/fabric-ipfs-interface/interface/fabric/pkg/blockchain"
	"github.com/thcrull/fabric-ipfs-interface/interface/fabric/pkg/config"
)

func main() {
	// Load Fabric configuration
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create the Fabric client
	client, err := metadata.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create Fabric client: %v", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Printf("Error closing client: %v", err)
		}
	}()

	// Create the metadata service
	service := metadata.NewMetadataService(client)

	// ----------------------
	// Add a metadata record
	// ----------------------
	newMetadata, err := service.AddMetadata(15, "tom", "encapkey1", "encmodelhash1", "homhash1")
	if err != nil {
		log.Fatalf("Failed to add metadata: %v", err)
	}
	fmt.Printf("Added metadata: %+v\n", newMetadata)

	// ----------------------
	// Get all metadata records
	// ----------------------
	allMetadata, err := service.GetAllMetadata()
	if err != nil {
		log.Fatalf("Failed to get all metadata: %v", err)
	}
	fmt.Printf("All metadata:\n")
	for _, m := range allMetadata {
		fmt.Printf("  %+v\n", m)
	}

	// ----------------------
	// Read the metadata record
	// ----------------------
	readMetadata, err := service.ReadMetadata(15, "tom")
	if err != nil {
		log.Fatalf("Failed to read metadata: %v", err)
	}
	fmt.Printf("Read metadata: %+v\n", readMetadata)

	// ----------------------
	// Update the metadata record
	// ----------------------
	updatedMetadata, err := service.UpdateMetadata(15, "tom", "encapkey2", "encmodelhash2", "homhash2")
	if err != nil {
		log.Fatalf("Failed to update metadata: %v", err)
	}
	fmt.Printf("Updated metadata: %+v\n", updatedMetadata)

	// ----------------------
	// Read the metadata record
	// ----------------------
	readMetadata2, err := service.ReadMetadata(15, "tom")
	if err != nil {
		log.Fatalf("Failed to read metadata: %v", err)
	}
	fmt.Printf("Read metadata: %+v\n", readMetadata2)

	// ----------------------
	// Check if metadata exists
	// ----------------------
	exists, err := service.MetadataExists(15, "tom")
	if err != nil {
		log.Fatalf("Failed to check metadata existence: %v", err)
	}
	fmt.Printf("Metadata exists: %v\n", exists)

	// ----------------------
	// Delete the metadata record
	// ----------------------
	deleted, err := service.DeleteMetadata(15, "tom")
	if err != nil {
		log.Fatalf("Failed to delete metadata: %v", err)
	}
	fmt.Printf("Metadata deleted: %v\n", deleted)

	// ----------------------
	// Check if metadata exists
	// ----------------------
	exists2, err := service.MetadataExists(15, "tom")
	if err != nil {
		log.Fatalf("Failed to check metadata existence: %v", err)
	}
	fmt.Printf("Metadata exists: %v\n", exists2)

	// ----------------------
	// Get all metadata records
	// ----------------------
	allMetadata2, err := service.GetAllMetadata()
	if err != nil {
		log.Fatalf("Failed to get all metadata: %v", err)
	}
	fmt.Printf("All metadata:\n")
	for _, m := range allMetadata2 {
		fmt.Printf("  %+v\n", m)
	}
}
