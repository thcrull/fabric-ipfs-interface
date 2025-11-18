package main

import (
	"context"
	"log"

	"github.com/thcrull/fabric-ipfs-interface/interface/fabric/api/config"
	"github.com/thcrull/fabric-ipfs-interface/interface/fabric/api/wrapper"
	"github.com/thcrull/fabric-ipfs-interface/interface/ipfs/api/config"
	"github.com/thcrull/fabric-ipfs-interface/interface/ipfs/api/wrapper"
	"github.com/thcrull/fabric-ipfs-interface/shared"
	pb "github.com/thcrull/fabric-ipfs-interface/weightpb"
)

// This serves as an example high-level use case of both interfaces.
func main() {
	//-----------------------------------
	// 1. Create Fabric and IPFS clients
	//-----------------------------------

	//Load config files
	fabricConfig, err := fabricconfig.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("error loading fabric config: %v", err)
		return
	}

	ipfsConfig, err := ipfsconfig.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("error loading ipfs config: %v", err)
		return
	}

	//Initialise the metadata service and IPFS client
	metadataService, err := fabricclient.NewMetadataService(fabricConfig)
	if err != nil {
		log.Fatalf("error creating fabric client: %v", err)
		return
	}

	ipfsClient, err := ipfsclient.NewIpfsClient(ipfsConfig)
	if err != nil {
		log.Fatalf("error creating ipfs client: %v", err)
		return
	}

	//---------------------------------------------
	// 2. Add a participant to the Fabric network
	//---------------------------------------------

	participant := &shared.Participant{
		ParticipantId:   "tom",
		EncapsulatedKey: "encapkey",
	}

	err = metadataService.AddParticipant(participant)
	if err != nil {
		log.Fatalf("error adding participant to Fabric network: %v", err)
		return
	}
	log.Printf("Added participant successfully.")

	//---------------------------------------------------------------------------------------------------
	// 3. Add a weight model to the IPFS network and the metadata to the Fabric network as a Participant
	//---------------------------------------------------------------------------------------------------

	//Serialize the weight model as a protobuf message and add it to the IPFS network
	weightModel := &pb.WeightModel{
		Values: []float64{1.25, 2.50, 3.75, 4.25, 5.50},
	}

	cid, err := ipfsClient.AddAndPinFile(context.Background(), weightModel)
	if err != nil {
		log.Fatalf("error adding and pinning file to IPFS: %v", err)
		return
	}
	log.Printf("Added and pinned file with CID: %s", cid)

	//Add the metadata of the weight model to the Fabric network
	participantModelMetadata := &shared.ParticipantModelMetadata{
		Epoch:           1,
		ParticipantId:   "tom",
		ModelHashCid:    cid,
		HomomorphicHash: "homhash",
	}

	err = metadataService.AddParticipantModelMetadata(participantModelMetadata)
	if err != nil {
		log.Fatalf("error adding participant model metadata to Fabric network: %v", err)
		return
	}
	log.Printf("Added participant model metadata successfully.")

	//---------------------------------------------------------------------------------------------------------------
	// 4. Fetch the participant's model metadata from the Fabric network and the weight model from the IPFS network
	//---------------------------------------------------------------------------------------------------------------

	//Fetch the metadata from the Fabric network
	participantModelMetadata, err = metadataService.GetParticipantModelMetadata(1, "tom")
	if err != nil {
		log.Fatalf("error fetching metadata from Fabric network: %v", err)
		return
	}
	log.Printf("Fetched metadata: %+v", participantModelMetadata)

	//Fetch weight model from the IPFS network
	var fetchedWeightModel pb.WeightModel
	if ipfsClient.GetFile(context.Background(), participantModelMetadata.ModelHashCid, &fetchedWeightModel) != nil {
		log.Fatalf("error fetching file from IPFS: %v", err)
		return
	}
	log.Printf("Fetched weight model: %+v", fetchedWeightModel.Values)

	//--------------------------------------------------------------------------------------------------
	// 5. Unpin the CID from the IPFS and delete the participant model metadata from the Fabric network
	//--------------------------------------------------------------------------------------------------

	if ipfsClient.UnpinFile(context.Background(), participantModelMetadata.ModelHashCid) != nil {
		log.Fatalf("error unpinning file from IPFS: %v", err)
		return
	}
	log.Printf("Unpinned file with CID: %s", participantModelMetadata.ModelHashCid)

	if metadataService.DeleteParticipantModelMetadata(1, "tom") != nil {
		log.Fatalf("error deleting metadata from Fabric network: %v", err)
		return
	}
	log.Printf("Deleted metadata successfully.")

	//------------------------------------------------------------------------------------------------------------
	// 6. Add an aggregated weight model to the IPFS network and the metadata to the Fabric network as an Aggregator
	//------------------------------------------------------------------------------------------------------------

	aggregatedWeightModel := &pb.WeightModel{
		Values: []float64{6.25, 7.50, 8.75, 9.25, 10.50},
	}

	cid, err = ipfsClient.AddAndPinFile(context.Background(), aggregatedWeightModel)
	if err != nil {
		log.Fatalf("error adding and pinning aggregated file to IPFS: %v", err)
		return
	}
	log.Printf("Added and pinned file with CID: %s", cid)

	aggregatorModelMetadata := &shared.AggregatorModelMetadata{
		Epoch:          1,
		ParticipantIds: []string{"tom"},
		ModelHashCid:   cid,
	}
	err = metadataService.AddAggregatorModelMetadata(aggregatorModelMetadata)
	if err != nil {
		log.Fatalf("error adding aggregated model metadata to Fabric network: %v", err)
		return
	}
	log.Printf("Added aggregated model metadata successfully.")

	// Deleting and unpinning the aggregated model metadata from
	// the Fabric and IPFS networks is similar to the steps from 5.

	//--------------------------------------------
	// 7. Fetch the logs from the Fabric network
	//--------------------------------------------

	history, err := metadataService.GetAllLogs()
	if err != nil {
		log.Fatalf("error fetching logs from Fabric network: %v", err)
		return
	}
	log.Printf("Fetched logs: %+v", history)

	//-------------------------------
	// 8. Close the metadata service
	//-------------------------------

	if metadataService.Close() != nil {
		log.Fatalf("error closing the metadata service: %v", err)
		return
	}
	log.Printf("Closed the metadata service.")
}
