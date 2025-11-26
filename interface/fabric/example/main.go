package main

import (
	"fmt"
	"log"

	fabricconfig "github.com/thcrull/fabric-ipfs-interface/interface/fabric/api/config"
	fabricclient "github.com/thcrull/fabric-ipfs-interface/interface/fabric/api/wrapper"
	"github.com/thcrull/fabric-ipfs-interface/shared"
)

// Example use case of the Fabric client wrapper used for testing all functionalities.
func main() {
	//--------------------------
	// INITIALIZE FABRIC CLIENT
	//--------------------------

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

	//-----------------------------
	// PARTICIPANT FUNCTIONALITIES
	//-----------------------------
	fmt.Println("-----Participant Functionalities-----")

	var participantThomas = shared.Participant{
		ParticipantId:                     "thomas",
		EncapsulatedKey:                   "key-thomas",
		HomomorphicSharedKeyCypher:        "key-thomas-homomorphic-shared-key-cypher",
		ParticipantCommunicationKeyCypher: "key-thomas-participant-communication-key-cypher",
		AggregatorCommunicationKeyCypher:  "key-thomas-aggregator-communication-key-cypher",
	}

	var participantMihnea = shared.Participant{
		ParticipantId:                     "mihnea",
		EncapsulatedKey:                   "key-mihnea",
		HomomorphicSharedKeyCypher:        "key-mihnea-homomorphic-shared-key-cypher",
		ParticipantCommunicationKeyCypher: "key-mihnea-participant-communication-key-cypher",
		AggregatorCommunicationKeyCypher:  "key-mihnea-aggregator-communication-key-cypher",
	}

	var participantIlinca = shared.Participant{
		ParticipantId:                     "ilinca",
		EncapsulatedKey:                   "key-ilinca",
		HomomorphicSharedKeyCypher:        "key-ilinca-homomorphic-shared-key-cypher",
		ParticipantCommunicationKeyCypher: "key-ilinca-participant-communication-key-cypher",
		AggregatorCommunicationKeyCypher:  "key-ilinca-aggregator-communication-key-cypher",
	}

	err = service.AddParticipant(&participantThomas)
	if err != nil {
		log.Fatalf("Failed to add participant Thomas: %v", err)
	}

	err = service.AddParticipant(&participantMihnea)
	if err != nil {
		log.Fatalf("Failed to add participant Mihnea: %v", err)
	}

	err = service.AddParticipant(&participantIlinca)
	if err != nil {
		log.Fatalf("Failed to add participant Ilinca: %v", err)
	}

	fmt.Println("Added participants successfully.")

	fetchedParticipant, err := service.GetParticipant("thomas")
	if err != nil {
		log.Fatalf("Failed to get participant Thomas: %v", err)
	}
	fmt.Printf("Fetched participant: %+v\n", fetchedParticipant)

	var participantThomasUpdated = shared.Participant{
		ParticipantId:                     "thomas",
		EncapsulatedKey:                   "key-thomas-updated",
		HomomorphicSharedKeyCypher:        "key-thomas-homomorphic-shared-key-cypher-updated",
		ParticipantCommunicationKeyCypher: "key-thomas-participant-communication-key-cypher-updated",
		AggregatorCommunicationKeyCypher:  "key-thomas-aggregator-communication-key-cypher-updated",
	}

	err = service.UpdateParticipant(&participantThomasUpdated)
	if err != nil {
		log.Fatalf("Failed to update participant Thomas: %v", err)
	}
	fmt.Println("Updated participant Thomas successfully.")

	err = service.DeleteParticipant("ilinca")
	if err != nil {
		log.Fatalf("Failed to delete participant Ilinca: %v", err)
	}
	fmt.Println("Deleted participant Ilinca successfully.")

	exists, err := service.ParticipantExists("ilinca")
	if err != nil {
		log.Fatalf("Failed to check if participant Ilinca exists: %v", err)
	}
	fmt.Printf("Checked participant Ilinca's existance: %t\n", exists)

	participants, err := service.GetAllParticipants()
	if err != nil {
		log.Fatalf("Failed to get all participants: %v", err)
	}
	fmt.Printf("Fetched all participants: %+v\n\n", participants)

	//--------------------------------------------
	// PARTICIPANT MODEL METADATA FUNCTIONALITIES
	//--------------------------------------------
	fmt.Println("-----Participant Model Metadata Functionalities-----")

	var participantThomasModelMetadata1 = shared.ParticipantModelMetadata{
		Epoch:           10,
		ParticipantId:   "thomas",
		ModelHashCid:    "thomas-model-cid",
		HomomorphicHash: "thomas-model-homomorphic-hash",
	}

	var participantThomasModelMetadata2 = shared.ParticipantModelMetadata{
		Epoch:           20,
		ParticipantId:   "thomas",
		ModelHashCid:    "thomas-model-cid",
		HomomorphicHash: "thomas-model-homomorphic-hash",
	}

	var participantMihneaModelMetadata1 = shared.ParticipantModelMetadata{
		Epoch:           10,
		ParticipantId:   "mihnea",
		ModelHashCid:    "mihnea-model-cid",
		HomomorphicHash: "mihnea-model-homomorphic-hash",
	}

	var participantMihneaModelMetadata2 = shared.ParticipantModelMetadata{
		Epoch:           20,
		ParticipantId:   "mihnea",
		ModelHashCid:    "mihnea-model-cid",
		HomomorphicHash: "mihnea-model-homomorphic-hash",
	}

	err = service.AddParticipantModelMetadata(&participantThomasModelMetadata1)
	if err != nil {
		log.Fatalf("Failed to add participant Thomas' first model metadata: %v", err)
	}

	err = service.AddParticipantModelMetadata(&participantThomasModelMetadata2)
	if err != nil {
		log.Fatalf("Failed to add participant Thomas' second model metadata: %v", err)
	}

	err = service.AddParticipantModelMetadata(&participantMihneaModelMetadata1)
	if err != nil {
		log.Fatalf("Failed to add participant Mihnea's first model metadata: %v", err)
	}

	err = service.AddParticipantModelMetadata(&participantMihneaModelMetadata2)
	if err != nil {
		log.Fatalf("Failed to add participant Mihnea's second model metadata: %v", err)
	}

	fmt.Println("Added participant model metadata successfully.")

	fetchedParticipantModelMetadata, err := service.GetParticipantModelMetadata(10, "mihnea")
	if err != nil {
		log.Fatalf("Failed to get participant Mihnea's first model metadata: %v", err)
	}
	fmt.Printf("Fetched participant Mihnea's first model metadata: %+v\n", fetchedParticipantModelMetadata)

	var participantMihneaModelMetadata1Updated = shared.ParticipantModelMetadata{
		Epoch:           10,
		ParticipantId:   "mihnea",
		ModelHashCid:    "mihnea-model-cid-updated",
		HomomorphicHash: "mihnea-model-homomorphic-hash-updated",
	}

	err = service.UpdateParticipantModelMetadata(&participantMihneaModelMetadata1Updated)
	if err != nil {
		log.Fatalf("Failed to update participant Mihnea's first model metadata: %v", err)
	}
	fmt.Println("Updated participant Mihnea's first model metadata successfully.")

	err = service.DeleteParticipantModelMetadata(20, "thomas")
	if err != nil {
		log.Fatalf("Failed to delete participant Thomas's second model metadata: %v", err)
	}
	fmt.Println("Deleted participant Thomas's second model metadata successfully.")

	exists, err = service.ParticipantModelMetadataExists(20, "thomas")
	if err != nil {
		log.Fatalf("Failed to check if participant Thomas's second model metadata exists: %v", err)
	}
	fmt.Printf("Checked participant Thomas's second model metadata's existance: %t\n", exists)

	participantModelMetadataList, err := service.GetAllParticipantModelMetadata()
	if err != nil {
		log.Fatalf("Failed to get all participant model metadata records: %v", err)
	}
	fmt.Printf("Fetched all participant model metadata records: %+v\n\n", participantModelMetadataList)

	participantModelMetadataList, err = service.GetAllParticipantModelMetadataByEpoch(10)
	if err != nil {
		log.Fatalf("Failed to get all participant model metadata records for epoch 10: %v", err)
	}
	fmt.Printf("Fetched all participant model metadata records for epoch 10: %+v\n\n", participantModelMetadataList)

	participantModelMetadataList, err = service.GetAllParticipantModelMetadataByParticipant("mihnea")
	if err != nil {
		log.Fatalf("Failed to get all participant model metadata records for participant Mihnea: %v", err)
	}
	fmt.Printf("Fetched all participant model metadata records for participant Mihnea: %+v\n\n", participantModelMetadataList)

	//-------------------------------------------
	// AGGREGATOR MODEL METADATA FUNCTIONALITIES
	//-------------------------------------------
	fmt.Println("-----Aggregator Model Metadata Functionalities-----")

	var aggregatorModelMetadata1 = shared.AggregatorModelMetadata{
		Epoch:          10,
		ModelHashCid:   "aggregator-model-cid",
		ParticipantIds: []string{"thomas", "mihnea"},
	}

	var aggregatorModelMetadata2 = shared.AggregatorModelMetadata{
		Epoch:          20,
		ModelHashCid:   "aggregator-model-cid",
		ParticipantIds: []string{"mihnea"},
	}

	err = service.AggregationCheck(&aggregatorModelMetadata1)
	if err != nil {
		log.Fatalf("Aggregation check result: %v", err)
	}
	fmt.Printf("Aggregation check succeded.")

	err = service.AddAggregatorModelMetadata(&aggregatorModelMetadata1)
	if err != nil {
		log.Fatalf("Failed to add the first aggregator model metadata: %v", err)
	}

	err = service.AddAggregatorModelMetadata(&aggregatorModelMetadata2)
	if err != nil {
		log.Fatalf("Failed to add the second aggregator model metadata: %v", err)
	}

	fmt.Println("Added aggregator model metadata successfully.")

	fetchedAggregatorModelMetadata, err := service.GetAggregatorModelMetadata(10)
	if err != nil {
		log.Fatalf("Failed to get the first aggregator model metadata: %v", err)
	}
	fmt.Printf("Fetched the first aggregator model metadata: %+v\n", fetchedAggregatorModelMetadata)

	var aggregatorModelMetadataUpdated = shared.AggregatorModelMetadata{
		Epoch:          10,
		ModelHashCid:   "aggregator-model-cid-updated",
		ParticipantIds: []string{"thomas", "mihnea"},
	}

	err = service.UpdateAggregatorModelMetadata(&aggregatorModelMetadataUpdated)
	if err != nil {
		log.Fatalf("Failed to update the first aggregator model metadata: %v", err)
	}
	fmt.Println("Updated the first aggregator model metadata successfully.")

	err = service.DeleteAggregatorModelMetadata(20)
	if err != nil {
		log.Fatalf("Failed to delete the second aggregator model metadata: %v", err)
	}
	fmt.Println("Deleted the second aggregator model metadata successfully.")

	exists, err = service.AggregatorModelMetadataExists(20)
	if err != nil {
		log.Fatalf("Failed to check if the second aggregator model metadata exists: %v", err)
	}
	fmt.Printf("Checked the second aggregator model metadata's existance: %t\n", exists)

	aggregatorModelMetadataList, err := service.GetAllAggregatorModelMetadata()
	if err != nil {
		log.Fatalf("Failed to get all aggregator model metadata records: %v", err)
	}
	fmt.Printf("Fetched all aggregator model metadata records: %+v\n\n", aggregatorModelMetadataList)

	//-----------------------------
	// HISTORY LOG FUNCTIONALITIES
	//-----------------------------
	fmt.Println("-----History Log Functionalities-----")

	logs, err := service.GetAllLogs()
	if err != nil {
		log.Fatalf("Failed to get all logs: %v", err)
	}
	fmt.Printf("Fetched all logs: %+v\n\n", logs)

	logs, err = service.GetAllLogsWithoutCreator()
	if err != nil {
		log.Fatalf("Failed to get all logs without creator information: %v", err)
	}
	fmt.Printf("Fetched all logs without creator information: %+v\n\n", logs)

	// NOTE: THE MSPID AND SERIAL NUMBER WILL BE DIFFERENT FOR YOUR MACHINE!!!
	logs, err = service.GetAllLogsForUser("Org1MSP", 2664457464428616324)
	if err != nil {
		log.Fatalf("Failed to get all logs for user: %v", err)
	}
	fmt.Printf("Fetched all logs for user: %+v\n\n", logs)

	//---------
	// CLEANUP
	//---------
	fmt.Println("-----Cleanup-----")

	err = service.DeleteAllAggregatorModelMetadata()
	if err != nil {
		log.Fatalf("Failed to delete all aggregator model metadata records: %v", err)
	}
	fmt.Println("Deleted all aggregator model metadata records successfully.")

	err = service.DeleteAllParticipantModelMetadata()
	if err != nil {
		log.Fatalf("Failed to delete all participant model metadata records: %v", err)
	}
	fmt.Println("Deleted all participant model metadata records successfully.")

	err = service.DeleteAllParticipants()
	if err != nil {
		log.Fatalf("Failed to delete all participants: %v", err)
	}
	fmt.Println("Deleted all participants successfully.")

	err = service.Close()
	if err != nil {
		log.Fatalf("Failed to close Fabric client: %v", err)
	}
	fmt.Println("Closed Fabric client successfully.")
}
