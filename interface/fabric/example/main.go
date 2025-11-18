package main

import (
	"context"
	"fmt"
	"log"

	fabricconfig "github.com/thcrull/fabric-ipfs-interface/interface/fabric/api/config"
	fabricclient "github.com/thcrull/fabric-ipfs-interface/interface/fabric/api/wrapper"
)

// Example use case of the Fabric client wrapper used for testing all functionalities.
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

	//// Create the metadata service
	//service, err := fabricclient.NewMetadataService(cfg)
	//if err != nil {
	//	log.Fatalf("Failed to create metadata service: %v", err)
	//}
	//
	//// -------------------------------------------------------------
	//// SECTION 1: PARTICIPANT FUNCTIONALITIES
	//// -------------------------------------------------------------
	//
	//fmt.Println("---- PARTICIPANT FUNCTIONALITIES ----")
	//
	//// Add participant
	//err = service.AddParticipant(&shared.Participant{ParticipantId: "tom", EncapsulatedKey: "encapkey1"})
	//if err != nil {
	//	log.Fatalf("Failed to add participant: %v", err)
	//}
	//fmt.Println("Added participant successfully.")
	//
	//// Get participant
	//participant, err := service.GetParticipant("tom")
	//if err != nil {
	//	log.Fatalf("Failed to get participant: %v", err)
	//}
	//fmt.Printf("Retrieved participant: %+v\n", participant)
	//
	//// Update participant
	//err = service.UpdateParticipant(&shared.Participant{ParticipantId: "tom", EncapsulatedKey: "encapkey2"})
	//if err != nil {
	//	log.Fatalf("Failed to update participant: %v", err)
	//}
	//fmt.Println("Updated participant successfully.")
	//
	//// Check participant existence
	//exists, err := service.ParticipantExists("tom")
	//if err != nil {
	//	log.Fatalf("Failed to check participant existence: %v", err)
	//}
	//fmt.Printf("Participant exists: %v\n", exists)
	//
	//// Get all participants
	//allParticipants, err := service.GetAllParticipants()
	//if err != nil {
	//	log.Fatalf("Failed to get all participants: %v", err)
	//}
	//fmt.Println("All participants:")
	//for _, p := range allParticipants {
	//	fmt.Printf("  %+v\n", p)
	//}
	//
	//// -------------------------------------------------------------
	//// SECTION 2: PARTICIPANT MODEL METADATA FUNCTIONALITIES
	//// -------------------------------------------------------------
	//
	//fmt.Println("---- PARTICIPANT MODEL METADATA FUNCTIONALITIES ----")
	//
	//// Add participant model metadata
	//err = service.AddParticipantModelMetadata(&shared.ParticipantModelMetadata{
	//	Epoch:           15,
	//	ParticipantId:   "tom",
	//	ModelHashCid:    "modelhash1",
	//	HomomorphicHash: "homhash1",
	//})
	//if err != nil {
	//	log.Fatalf("Failed to add participant model metadata: %v", err)
	//}
	//fmt.Println("Added participant model metadata successfully.")
	//
	//// Get participant model metadata
	//pm, err := service.GetParticipantModelMetadata(15, "tom")
	//if err != nil {
	//	log.Fatalf("Failed to get participant model metadata: %v", err)
	//}
	//fmt.Printf("Retrieved participant model metadata: %+v\n", pm)
	//
	//// Update participant model metadata
	//err = service.UpdateParticipantModelMetadata(&shared.ParticipantModelMetadata{
	//	Epoch:           15,
	//	ParticipantId:   "tom",
	//	ModelHashCid:    "modelhash2",
	//	HomomorphicHash: "homhash2",
	//})
	//if err != nil {
	//	log.Fatalf("Failed to update participant model metadata: %v", err)
	//}
	//fmt.Println("Updated participant model metadata successfully.")
	//
	//// Add a second participant model metadata
	//err = service.AddParticipantModelMetadata(&shared.ParticipantModelMetadata{
	//	Epoch:           15,
	//	ParticipantId:   "tom1",
	//	ModelHashCid:    "modelhash2",
	//	HomomorphicHash: "homhash2",
	//})
	//if err != nil {
	//	log.Fatalf("Failed to add a second participant model metadata: %v", err)
	//}
	//fmt.Println("Added a second participant model metadata successfully.")
	//
	//// Add a third participant model metadata
	//err = service.AddParticipantModelMetadata(&shared.ParticipantModelMetadata{
	//	Epoch:           25,
	//	ParticipantId:   "tom",
	//	ModelHashCid:    "modelhash3",
	//	HomomorphicHash: "homhash3",
	//})
	//if err != nil {
	//	log.Fatalf("Failed to add a third participant model metadata: %v", err)
	//}
	//fmt.Println("Added a third participant model metadata successfully.")
	//
	//// Get all participant model metadata
	//allPM, err := service.GetAllParticipantModelMetadata()
	//if err != nil {
	//	log.Fatalf("Failed to get all participant model metadata: %v", err)
	//}
	//fmt.Println("All participant model metadata:")
	//for _, p := range allPM {
	//	fmt.Printf("  %+v\n", p)
	//}
	//
	//// Get participant models metadata by participant
	//pmByParticipant, err := service.GetAllParticipantModelMetadataByParticipant("tom")
	//if err != nil {
	//	log.Fatalf("Failed to get participant model metadata by participant: %v", err)
	//}
	//fmt.Println("Participant model metadata by participant:")
	//for _, p := range pmByParticipant {
	//	fmt.Printf("  %+v\n", p)
	//}
	//
	//// Get participant models metadata by epoch
	//pmByEpoch, err := service.GetAllParticipantModelMetadataByEpoch(15)
	//if err != nil {
	//	log.Fatalf("Failed to get participant model metadata by epoch: %v", err)
	//}
	//fmt.Println("Participant model metadata by epoch:")
	//for _, p := range pmByEpoch {
	//	fmt.Printf("  %+v\n", p)
	//}
	//
	//// Check if a record exists
	//pmExists, err := service.ParticipantModelMetadataExists(15, "tom")
	//if err != nil {
	//	log.Fatalf("Failed to check if participant model metadata exists: %v", err)
	//}
	//fmt.Printf("Participant model metadata exists: %v\n", pmExists)
	//
	//// Delete participant model metadata
	//err = service.DeleteParticipantModelMetadata(15, "tom")
	//if err != nil {
	//	log.Fatalf("Failed to delete participant model metadata: %v", err)
	//}
	//fmt.Println("Deleted participant model metadata successfully.")
	//
	//// -------------------------------------------------------------
	//// SECTION 3: AGGREGATOR MODEL METADATA FUNCTIONALITIES
	//// -------------------------------------------------------------
	//
	//fmt.Println("---- AGGREGATOR MODEL METADATA FUNCTIONALITIES ----")
	//
	//// Add aggregator model metadata
	//err = service.AddAggregatorModelMetadata(&shared.AggregatorModelMetadata{
	//	Epoch:          20,
	//	ModelHashCid:   "aggmodelhash1",
	//	ParticipantIds: []string{"tom", "tom1"},
	//})
	//if err != nil {
	//	log.Fatalf("Failed to add aggregator model metadata: %v", err)
	//}
	//fmt.Println("Added aggregator model metadata successfully.")
	//
	//// Get aggregator model metadata
	//am, err := service.GetAggregatorModelMetadata(20)
	//if err != nil {
	//	log.Fatalf("Failed to get aggregator model metadata: %v", err)
	//}
	//fmt.Printf("Retrieved aggregator model metadata: %+v\n", am)
	//
	//// Update aggregator model metadata
	//err = service.UpdateAggregatorModelMetadata(&shared.AggregatorModelMetadata{
	//	Epoch:          20,
	//	ModelHashCid:   "aggmodelhash2",
	//	ParticipantIds: []string{"tom", "tom1", "alice"},
	//})
	//if err != nil {
	//	log.Fatalf("Failed to update aggregator model metadata: %v", err)
	//}
	//fmt.Println("Updated aggregator model metadata successfully.")
	//
	//// Check existence
	//amExists, err := service.AggregatorModelMetadataExists(20)
	//if err != nil {
	//	log.Fatalf("Failed to check aggregator model metadata existence: %v", err)
	//}
	//fmt.Printf("Aggregator model metadata exists: %v\n", amExists)
	//
	//// Add a second aggregator model metadata
	//err = service.AddAggregatorModelMetadata(&shared.AggregatorModelMetadata{
	//	Epoch:          30,
	//	ModelHashCid:   "aggmodelhash2",
	//	ParticipantIds: []string{"tom2", "tom3"},
	//})
	//if err != nil {
	//	log.Fatalf("Failed to add a second aggregator model metadata: %v", err)
	//}
	//fmt.Println("Added a second aggregator model metadata successfully.")
	//
	//// Get all aggregator model metadata
	//allAM, err := service.GetAllAggregatorModelMetadata()
	//if err != nil {
	//	log.Fatalf("Failed to get all aggregator model metadata: %v", err)
	//}
	//fmt.Println("All aggregator model metadata:")
	//for _, a := range allAM {
	//	fmt.Printf("  %+v\n", a)
	//}
	//
	//// Delete aggregator model metadata
	//err = service.DeleteAggregatorModelMetadata(20)
	//if err != nil {
	//	log.Fatalf("Failed to delete aggregator model metadata: %v", err)
	//}
	//fmt.Println("Deleted aggregator model metadata successfully.")
	//
	//// -------------------------------------------------------------
	//// SECTION 4: LOGS
	//// -------------------------------------------------------------
	//
	//fmt.Println("---- LEDGER LOGS ----")
	//
	//logs, err := service.GetAllLogs()
	//if err != nil {
	//	log.Fatalf("Failed to get all logs: %v", err)
	//}
	//fmt.Println("All logs:")
	//for _, logEntry := range logs {
	//	fmt.Printf("  %+v\n", logEntry)
	//}
	//
	//// -------------------------------------------------------------
	//// CLEANUP
	//// -------------------------------------------------------------
	//
	//fmt.Println("---- CLEANUP ----")
	//
	//// Delete all participants
	//err = service.DeleteAllParticipants()
	//if err != nil {
	//	log.Fatalf("Failed to delete all participants: %v", err)
	//}
	//fmt.Println("Deleted all participants successfully.")
	//
	//// Delete all participant model metadata
	//err = service.DeleteAllParticipantModelMetadata()
	//if err != nil {
	//	log.Fatalf("Failed to delete all participant model metadata: %v", err)
	//}
	//fmt.Println("Deleted all participant model metadata successfully.")
	//
	//// Delete all aggregator model metadata
	//err = service.DeleteAllAggregatorModelMetadata()
	//if err != nil {
	//	log.Fatalf("Failed to delete all aggregator model metadata: %v", err)
	//}
	//fmt.Println("Deleted all aggregator model metadata successfully.")
	//
	//fmt.Println("All tests completed successfully.")

	ctx := context.Background()
	txID := "20f293db6a4d0f265dcd3814c2f3d736a7f263b0172299069c2ee694b30ed610"

	creator, err := client.GetTransactionCreator(ctx, txID, 0)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	fmt.Printf("Creator MSPID: %s\n", creator.MSPID)
	fmt.Printf("Block Number: %d\n", creator.BlockNum)
	fmt.Printf("Subject CN: %s\n", creator.Cert.Subject.CommonName)
	fmt.Printf("Issuer CN: %s\n", creator.Cert.Issuer.CommonName)
	fmt.Printf("Full Cert:\n%+v\n", creator.Cert)
}
