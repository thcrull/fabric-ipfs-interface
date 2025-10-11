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

	//-----------------------------------------------------------------------------------
	// 2. Add a weight model to the IPFS network and the metadata to the Fabric network
	//-----------------------------------------------------------------------------------

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
	metadata := &shared.Metadata{
		Epoch:           1,
		ParticipantId:   "tom",
		EncapsulatedKey: "encapkey",
		EncModelHash:    cid,
		HomomorphicHash: "homhash",
	}

	err = metadataService.AddMetadata(metadata)
	if err != nil {
		log.Fatalf("error adding metadata to Fabric network: %v", err)
		return
	}
	log.Printf("Added metadata successfully.")

	//------------------------------------------------------------------------------------------
	// 3. Fetch the metadata from the Fabric network and the weight model from the IPFS network
	//------------------------------------------------------------------------------------------

	//Fetch the metadata from the Fabric network
	metadata, err = metadataService.GetMetadata(1, "tom")
	if err != nil {
		log.Fatalf("error fetching metadata from Fabric network: %v", err)
		return
	}
	log.Printf("Fetched metadata: %+v", metadata)

	//Fetch weight model from the IPFS network
	var fetchedWeightModel pb.WeightModel
	if ipfsClient.GetFile(context.Background(), metadata.EncModelHash, &fetchedWeightModel) != nil {
		log.Fatalf("error fetching file from IPFS: %v", err)
		return
	}
	log.Printf("Fetched weight model: %+v", fetchedWeightModel.Values)

	//--------------------------------------------------------------------------------
	// 4. Unpin the CID from the IPFS and delete the metadata from the Fabric network
	//--------------------------------------------------------------------------------

	if ipfsClient.UnpinFile(context.Background(), metadata.EncModelHash) != nil {
		log.Fatalf("error unpinning file from IPFS: %v", err)
		return
	}
	log.Printf("Unpinned file with CID: %s", metadata.EncModelHash)

	if metadataService.DeleteMetadata(1, "tom") != nil {
		log.Fatalf("error deleting metadata from Fabric network: %v", err)
		return
	}
	log.Printf("Deleted metadata successfully.")

	//-------------------------------
	// 5. Close the metadata service
	//-------------------------------

	if metadataService.Close() != nil {
		log.Fatalf("error closing the metadata service: %v", err)
		return
	}
	log.Printf("Closed the metadata service.")
}
