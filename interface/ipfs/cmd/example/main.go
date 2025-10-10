package main

import (
	"context"
	"fmt"

	"github.com/thcrull/fabric-ipfs-interface/interface/ipfs/pkg/config"
	"github.com/thcrull/fabric-ipfs-interface/interface/ipfs/pkg/ipfs"
	pb "github.com/thcrull/fabric-ipfs-interface/weightpb"
)

func main() {
	ctx := context.Background()

	// Hardcoded config file path
	configPath := "config.yaml"

	// Load IPFS config from YAML file
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Println("Failed to load config:", err)
		return
	}

	// Initialize IPFS client with loaded config
	ipfsClient, err := client.NewIpfsClient(cfg)
	if err != nil {
		fmt.Println("Failed to create IPFS client:", err)
		return
	}

	// Create protobuf message to add
	weightModel := &pb.WeightModel{
		Values: []float64{1.01, 2.01, 3.01, 4.01, 5.01},
	}

	// Add protobuf data to IPFS (no pinning yet)
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
}
