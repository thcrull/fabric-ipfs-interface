package main

import (
	"fmt"
	"log"

	"github.com/thcrull/fabric-ipfs-interface/interface/fabric/pkg/config"
	"github.com/thcrull/fabric-ipfs-interface/interface/fabric/pkg/wrapper"
	"github.com/thcrull/fabric-ipfs-interface/shared"
)

// Example use case of the Fabric client wrapper used for testing purposes.
func main() {
	// Load Fabric configuration
	cfg, err := fabricconfig.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create the Fabric client
	client, err := fabricclient.NewFabricClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create Fabric client: %v", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Printf("Error closing client: %v", err)
		}
		fmt.Println("Closed the fabric client.")
	}()

	// Create the metadata service
	service, err := fabricclient.NewMetadataService(cfg)
	if err != nil {
		log.Fatalf("Failed to create metadata service: %v", err)
	}

	// Add metadata records
	err = service.AddMetadata(&shared.Metadata{Epoch: 15, ParticipantId: "tom", EncapsulatedKey: "encapkey1", EncModelHash: "encmodelhash1", HomomorphicHash: "homhash1"})
	if err != nil {
		log.Fatalf("Failed to add metadata: %v", err)
	}
	fmt.Printf("Added metadata successfully.\n")

	err = service.AddMetadata(&shared.Metadata{Epoch: 10, ParticipantId: "tom", EncapsulatedKey: "encapkey1", EncModelHash: "encmodelhash1", HomomorphicHash: "homhash1"})
	if err != nil {
		log.Fatalf("Failed to add metadata: %v", err)
	}
	fmt.Printf("Added metadata successfully.\n")

	err = service.AddMetadata(&shared.Metadata{Epoch: 15, ParticipantId: "tom1", EncapsulatedKey: "encapkey1", EncModelHash: "encmodelhash1", HomomorphicHash: "homhash1"})
	if err != nil {
		log.Fatalf("Failed to add metadata: %v", err)
	}
	fmt.Printf("Added metadata successfully.\n")

	// Get all metadata records
	allMetadata, err := service.GetAllMetadata()
	if err != nil {
		log.Fatalf("Failed to get all metadata: %v", err)
	}
	fmt.Printf("All metadata:\n")
	for _, metadata := range allMetadata {
		fmt.Printf("  %+v\n", metadata)
	}

	// Get all metadata records by participant
	participantMetadata, err := service.GetAllMetadataByParticipant("tom")
	if err != nil {
		log.Fatalf("Failed to get all metadata from participant %s: %v", "tom", err)
	}
	fmt.Printf("All metadata from participant %s:\n", "tom")
	for _, metadata := range participantMetadata {
		fmt.Printf("  %+v\n", metadata)
	}

	// Get all metadata records by epoch
	epochMetadata, err := service.GetAllMetadataByEpoch(15)
	if err != nil {
		log.Fatalf("Failed to get all metadata from epoch %d: %v", 15, err)
	}
	fmt.Printf("All metadata from epoch %d:\n", 15)
	for _, metadata := range epochMetadata {
		fmt.Printf("  %+v\n", metadata)
	}

	// Read the metadata record
	metadataRecord, err := service.GetMetadata(15, "tom")
	if err != nil {
		log.Fatalf("Failed to read metadata: %v", err)
	}
	fmt.Printf("Read metadata: %+v\n", metadataRecord)

	// Update the metadata record
	err = service.UpdateMetadata(&shared.Metadata{Epoch: 15, ParticipantId: "tom", EncapsulatedKey: "encapkey2", EncModelHash: "encmodelhash2", HomomorphicHash: "homhash2"})
	if err != nil {
		log.Fatalf("Failed to update metadata: %v", err)
	}
	fmt.Printf("Updated metadata successfully.\n")

	// Read the updated metadata record
	updatedMetadataRecord, err := service.GetMetadata(15, "tom")
	if err != nil {
		log.Fatalf("Failed to read metadata: %v", err)
	}
	fmt.Printf("Read metadata: %+v\n", updatedMetadataRecord)

	// Check if metadata exists
	metadataExists, err := service.MetadataExists(15, "tom")
	if err != nil {
		log.Fatalf("Failed to check metadata existence: %v", err)
	}
	fmt.Printf("Metadata exists: %v\n", metadataExists)

	// Delete the metadata record
	err = service.DeleteMetadata(15, "tom")
	if err != nil {
		log.Fatalf("Failed to delete metadata: %v", err)
	}
	fmt.Printf("Metadata deleted successfully.\n")

	// Check if metadata still exists
	metadataStillExists, err := service.MetadataExists(15, "tom")
	if err != nil {
		log.Fatalf("Failed to check metadata existence: %v", err)
	}
	fmt.Printf("Metadata exists: %v\n", metadataStillExists)

	// Get remaining metadata records
	remainingMetadata, err := service.GetAllMetadata()
	if err != nil {
		log.Fatalf("Failed to get all metadata: %v", err)
	}
	fmt.Printf("All metadata:\n")
	for _, metadata := range remainingMetadata {
		fmt.Printf("  %+v\n", metadata)
	}

	// Delete all metadata records
	err = service.DeleteAllMetadata()
	if err != nil {
		log.Fatalf("Failed to delete all metadata: %v", err)
	}
	fmt.Printf("Deleted all metadata successfully.\n")

	// Confirm deletion
	finalMetadataList, err := service.GetAllMetadata()
	if err != nil {
		log.Fatalf("Failed to get all metadata: %v", err)
	}
	fmt.Printf("All metadata:\n")
	for _, metadata := range finalMetadataList {
		fmt.Printf("  %+v\n", metadata)
	}
}
