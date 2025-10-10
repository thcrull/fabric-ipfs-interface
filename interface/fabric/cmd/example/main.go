package example

import (
	"fmt"
	"log"

	"github.com/thcrull/fabric-ipfs-interface/interface/fabric/pkg/config"
	"github.com/thcrull/fabric-ipfs-interface/interface/fabric/pkg/wrapper"
)

// Example use case of the Fabric client wrapper used for testing purposes.
func main() {
	// Load Fabric configuration
	cfg, err := fabricconfig.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create the Fabric client
	client, err := metadata.NewFabricClient(cfg)
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

	// Add the metadata records
	newMetadata, err := service.AddMetadata(15, "tom", "encapkey1", "encmodelhash1", "homhash1")
	if err != nil {
		log.Fatalf("Failed to add metadata: %v", err)
	}
	fmt.Printf("Added metadata: %+v\n", newMetadata)

	newMetadata1, err := service.AddMetadata(10, "tom", "encapkey1", "encmodelhash1", "homhash1")
	if err != nil {
		log.Fatalf("Failed to add metadata: %v", err)
	}
	fmt.Printf("Added metadata: %+v\n", newMetadata1)

	newMetadata2, err := service.AddMetadata(15, "tom1", "encapkey1", "encmodelhash1", "homhash1")
	if err != nil {
		log.Fatalf("Failed to add metadata: %v", err)
	}
	fmt.Printf("Added metadata: %+v\n", newMetadata2)

	// Get all metadata records
	allMetadata, err := service.GetAllMetadata()
	if err != nil {
		log.Fatalf("Failed to get all metadata: %v", err)
	}
	fmt.Printf("All metadata:\n")
	for _, m := range allMetadata {
		fmt.Printf("  %+v\n", m)
	}

	// Get all metadata records by participant
	allMetadata1, err := service.GetAllMetadataByParticipant("tom")
	if err != nil {
		log.Fatalf("Failed to get all metadata from participant %s: %v", "tom", err)
	}
	fmt.Printf("All metadata from participant %s:\n", "tom")
	for _, m := range allMetadata1 {
		fmt.Printf("  %+v\n", m)
	}

	// Get all metadata records by epoch
	allMetadata2, err := service.GetAllMetadataByEpoch(15)
	if err != nil {
		log.Fatalf("Failed to get all metadata from epoch %d: %v", 15, err)
	}
	fmt.Printf("All metadata from epoch %d:\n", 15)
	for _, m := range allMetadata2 {
		fmt.Printf("  %+v\n", m)
	}

	// Read the metadata record
	readMetadata, err := service.ReadMetadata(15, "tom")
	if err != nil {
		log.Fatalf("Failed to read metadata: %v", err)
	}
	fmt.Printf("Read metadata: %+v\n", readMetadata)

	// Update the metadata record
	updatedMetadata, err := service.UpdateMetadata(15, "tom", "encapkey2", "encmodelhash2", "homhash2")
	if err != nil {
		log.Fatalf("Failed to update metadata: %v", err)
	}
	fmt.Printf("Updated metadata: %+v\n", updatedMetadata)

	// Read the metadata record
	readMetadata2, err := service.ReadMetadata(15, "tom")
	if err != nil {
		log.Fatalf("Failed to read metadata: %v", err)
	}
	fmt.Printf("Read metadata: %+v\n", readMetadata2)

	// Check if metadata exists
	exists, err := service.MetadataExists(15, "tom")
	if err != nil {
		log.Fatalf("Failed to check metadata existence: %v", err)
	}
	fmt.Printf("Metadata exists: %v\n", exists)

	// Delete the metadata record
	deleted, err := service.DeleteMetadata(15, "tom")
	if err != nil {
		log.Fatalf("Failed to delete metadata: %v", err)
	}
	fmt.Printf("Metadata deleted: %v\n", deleted)

	// Check if metadata exists
	exists2, err := service.MetadataExists(15, "tom")
	if err != nil {
		log.Fatalf("Failed to check metadata existence: %v", err)
	}
	fmt.Printf("Metadata exists: %v\n", exists2)

	// Get all metadata records
	allMetadata3, err := service.GetAllMetadata()
	if err != nil {
		log.Fatalf("Failed to get all metadata: %v", err)
	}
	fmt.Printf("All metadata:\n")
	for _, m := range allMetadata3 {
		fmt.Printf("  %+v\n", m)
	}

	// Delete all metadata records
	deleteAllResponse, err := service.DeleteAllMetadata()
	if err != nil {
		log.Fatalf("Failed to delete all metadata: %v", err)
	}
	fmt.Printf("Deleted all metadata: %t\n", deleteAllResponse)

	// Get all metadata records
	allMetadata4, err := service.GetAllMetadata()
	if err != nil {
		log.Fatalf("Failed to get all metadata: %v", err)
	}
	fmt.Printf("All metadata:\n")
	for _, m := range allMetadata4 {
		fmt.Printf("  %+v\n", m)
	}

	// Close the client
	err = client.Close()
	if err != nil {
		log.Fatalf("Failed to close client: %v", err)
	}

	fmt.Println("Closed the fabric client.")
}
