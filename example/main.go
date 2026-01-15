package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/thcrull/fabric-ipfs-interface/interface/fabric/wrapper"
	"github.com/thcrull/fabric-ipfs-interface/interface/ipfs/wrapper"
	pb "github.com/thcrull/fabric-ipfs-interface/weight_pb"
)

// readVectorFromFile reads int64 values from a binary file
func readVectorFromFile(filename string) ([]int64, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if len(data)%8 != 0 {
		return nil, fmt.Errorf("file size not multiple of 8")
	}

	n := len(data) / 8
	vec := make([]int64, n)

	for i := 0; i < n; i++ {
		bits := binary.LittleEndian.Uint64(data[i*8 : i*8+8])
		vec[i] = int64(bits)
	}
	return vec, nil
}

func main() {
	startTime := time.Now() // total runtime start

	//-----------------------------------
	// 1. Create Fabric and IPFS clients
	//-----------------------------------
	metadataService, err := fabric_client.NewMetadataService("../config/admin.yaml")
	if err != nil {
		log.Fatalf("error creating fabric client: %v", err)
	}

	ipfsClient, err := ipfs_client.NewIpfsClient("../config/admin.yaml")
	if err != nil {
		log.Fatalf("error creating ipfs client: %v", err)
	}

	//---------------------------------------------
	// 2. Add a participant to the Fabric network
	//---------------------------------------------
	participantId := 10
	err = metadataService.AddParticipant(
		participantId,
		"encapsulated-key",
		"homomorphic-shared-key",
		"participant-comm-key",
	)
	if err != nil {
		log.Fatalf("error adding participant: %v", err)
	}
	log.Printf("Added participant %d successfully.", participantId)

	//---------------------------------------------
	// 3. Add an aggregator to the Fabric network
	//---------------------------------------------
	aggregatorId := 20
	err = metadataService.AddAggregator(aggregatorId, map[string]string{
		strconv.Itoa(aggregatorId): "participant-comm-key",
	})
	if err != nil {
		log.Fatalf("error adding aggregator: %v", err)
	}
	log.Printf("Added aggregator %d successfully.", aggregatorId)

	//-----------------------------------------------------
	// 4. Read weight model from a file and add it to IPFS
	//-----------------------------------------------------
	step4Start := time.Now() // start timer for Step 4→6

	vec, err := readVectorFromFile("../data/data_100000000.bin")
	if err != nil {
		log.Fatalf("failed to read vector: %v", err)
	}
	log.Printf("Read %d values from file.", len(vec))

	weightModel := &pb.WeightModel{Values: vec}

	cid, err := ipfsClient.AddAndPinFile(context.Background(), weightModel)
	if err != nil {
		log.Fatalf("failed to add weight model to IPFS: %v", err)
	}
	log.Printf("Pinned weight model to IPFS with CID: %s", cid)

	err = metadataService.AddParticipantModelMetadata(participantId, 1, cid, "homomorphic-hash-placeholder")
	if err != nil {
		log.Fatalf("failed to add participant model metadata: %v", err)
	}
	log.Printf("Added participant model metadata successfully.")

	//-------------------------------------------------
	// 5. Fetch participant model metadata and model
	//-------------------------------------------------
	modelMeta, err := metadataService.GetParticipantModelMetadata(participantId, 1)
	if err != nil {
		log.Fatalf("failed to fetch participant model metadata: %v", err)
	}
	log.Printf("Fetched participant model metadata: %+v", modelMeta)

	var fetchedModel pb.WeightModel
	err = ipfsClient.GetFile(context.Background(), modelMeta.ModelHashCid, &fetchedModel)
	if err != nil {
		log.Fatalf("failed to fetch weight model from IPFS: %v", err)
	}

	//--------------------------------------------------
	// 6. Add an aggregated weight model (same vector)
	//--------------------------------------------------
	aggregatedModel := &pb.WeightModel{Values: vec}
	cid, err = ipfsClient.AddAndPinFile(context.Background(), aggregatedModel)
	if err != nil {
		log.Fatalf("failed to add aggregated model to IPFS: %v", err)
	}
	log.Printf("Pinned aggregated model to IPFS with CID: %s", cid)

	err = metadataService.AddAggregatorModelMetadata(aggregatorId, 1, cid, []int{participantId})
	if err != nil {
		log.Fatalf("failed to add aggregator model metadata: %v", err)
	}
	log.Printf("Added aggregator model metadata successfully.")

	step4to6Elapsed := time.Since(step4Start)
	log.Printf("Elapsed time for reading, fetching and adding: %s", step4to6Elapsed)

	//-------------------------------------------------------------
	// 7. Clean up: unpin a participant model and delete metadata
	//-------------------------------------------------------------
	err = ipfsClient.UnpinFile(context.Background(), modelMeta.ModelHashCid)
	if err != nil {
		log.Fatalf("failed to unpin participant model: %v", err)
	}
	log.Printf("Unpinned participant model with CID: %s", modelMeta.ModelHashCid)

	err = metadataService.DeleteParticipantModelMetadata(participantId, 1)
	if err != nil {
		log.Fatalf("failed to delete participant model metadata: %v", err)
	}
	log.Printf("Deleted participant model metadata successfully.")

	//---------------------------
	// 8. Total execution time
	//---------------------------
	totalElapsed := time.Since(startTime)
	log.Printf("Total execution time: %s", totalElapsed)

	//---------------------
	// Optional: Teardown
	//---------------------
	if err = metadataService.DeleteParticipant(participantId); err != nil {
		log.Fatalf("failed to delete participant: %v", err)
	}

	if err = metadataService.DeleteAggregatorModelMetadata(aggregatorId, 1); err != nil {
		log.Fatalf("failed to delete aggregator model metadata: %v", err)
	}

	if err = metadataService.DeleteAggregator(aggregatorId); err != nil {
		log.Fatalf("failed to delete aggregator: %v", err)
	}

	//--------------------------------
	// 9. Close the metadata service
	//--------------------------------
	if err := metadataService.Close(); err != nil {
		log.Fatalf("failed to close metadata service: %v", err)
	}
	log.Printf("Closed metadata service.")

}
