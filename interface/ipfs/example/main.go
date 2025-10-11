package main

import (
	"context"
	"fmt"

	"github.com/thcrull/fabric-ipfs-interface/interface/ipfs/api/config"
	"github.com/thcrull/fabric-ipfs-interface/interface/ipfs/api/wrapper"
	pb "github.com/thcrull/fabric-ipfs-interface/weightpb"
)

// Example use case of the IPFS client wrapper used for testing purposes.
func main() {
	ctx := context.Background()
	configPath := "config.yaml"

	// Load IPFS config from the YAML file
	cfg, err := ipfsconfig.LoadConfig(configPath)
	if err != nil {
		fmt.Println("Failed to load config:", err)
		return
	}

	// Initialise the IPFS client with loaded config
	ipfsClient, err := ipfsclient.NewIpfsClient(cfg)
	if err != nil {
		fmt.Println("Failed to create IPFS client:", err)
		return
	}

	// Create the protobuf message to add
	weightModel := &pb.WeightModel{
		Values: []float64{1.01, 2.01, 3.01, 4.01, 5.01},
	}

	// Add protobuf data to IPFS
	cid, err := ipfsClient.AddFile(ctx, weightModel)
	if err != nil {
		fmt.Println("Failed to add file to IPFS:", err)
		return
	}

	fmt.Println("Added file with CID:", cid)

	// Pin the CID explicitly
	if err := ipfsClient.PinFile(ctx, cid); err != nil {
		fmt.Println("Failed to pin CID:", err)
		return
	}

	// Retrieve the protobuf message back from IPFS
	var loadedModel pb.WeightModel
	if err := ipfsClient.GetFile(ctx, cid, &loadedModel); err != nil {
		fmt.Println("Failed to retrieve file from IPFS:", err)
		return
	}

	// Print the loaded values
	fmt.Println("Loaded WeightModel values:", loadedModel.Values)

	// Unpin the CID
	if err := ipfsClient.UnpinFile(ctx, cid); err != nil {
		fmt.Println("Failed to unpin CID:", err)
		return
	}

	fmt.Println("Successfully unpinned CID:", cid)

	// Add and pin a new file
	weightModel2 := &pb.WeightModel{
		Values: []float64{10.01, 20.01, 30.01, 40.01, 50.01},
	}

	cid2, err := ipfsClient.AddAndPinFile(context.Background(), weightModel2)
	if err != nil {
		fmt.Println("Failed to add and pin file to IPFS:", err)
		return
	}

	fmt.Println("Added and pinned file with CID:", cid2)

	// Retrieve the protobuf message back from IPFS
	var loadedModel2 pb.WeightModel
	if err := ipfsClient.GetFile(ctx, cid2, &loadedModel2); err != nil {
		fmt.Println("Failed to retrieve file from IPFS:", err)
		return
	}

	// Print the loaded values
	fmt.Println("Loaded WeightModel values:", loadedModel2.Values)

	// Unpin the CID
	if err := ipfsClient.UnpinFile(ctx, cid2); err != nil {
		fmt.Println("Failed to unpin CID:", err)
		return
	}

	fmt.Println("Successfully unpinned CID:", cid2)
}
