package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/ipfs/boxo/files"
	"github.com/ipfs/kubo/client/rpc"
	"google.golang.org/protobuf/proto"

	pb "github.com/thcrull/fabric-ipfs-interface/weight.pb.go"
)

func main() {
	httpClient := http.Client{}

	// Connect to the IPFS node
	nodeHttpClient, err := rpc.NewURLApiWithClient("http://localhost:5001", &httpClient)
	if err != nil {
		fmt.Println("Error creating a HTTP API client: ", err)
		return
	}

	// Create a file from string content
	weightModel := []float64{1.0, 2.0, 3.0, 4.0, 5.0}

	message := &WeightModel{
		Values: weightModel,
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		fmt.Println("Error marshalling protobuf:", err)
		return
	}

	file := files.NewReaderFile(bytes.NewReader(data))

	cid, err := nodeHttpClient.Unixfs().Add(context.Background(), file)
	if err != nil {
		fmt.Println("Error adding file to IPFS:", err)
		return
	}

	fmt.Println("Added file with CID:", cid)
}
