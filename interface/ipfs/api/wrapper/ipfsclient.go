package ipfsclient

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/ipfs/boxo/files"
	"github.com/ipfs/boxo/path"
	"github.com/ipfs/kubo/client/rpc"
	"github.com/thcrull/fabric-ipfs-interface/interface/ipfs/api/config"
)

// IpfsClient is a wrapper around the IPFS node HTTP API. It provides
// convenient methods for interacting with the IPFS node.
type IpfsClient struct {
	httpClient  *http.Client
	NodeHttpApi *rpc.HttpApi
}

// NewIpfsClient creates a new IpfsClient instance.
func NewIpfsClient(configPath string) (*IpfsClient, error) {
	cfg, err := ipfsconfig.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("error loading IPFS config: %w", err)
	}

	httpClient := &http.Client{}

	nodeHttpApi, err := rpc.NewURLApiWithClient(cfg.Ipfs.NodePath, httpClient)
	if err != nil {
		return nil, fmt.Errorf("error creating an IPFS Node HTTP API: %w", err)
	}

	return &IpfsClient{
		httpClient:  httpClient,
		NodeHttpApi: nodeHttpApi,
	}, nil
}

// AddFile adds a protobuf message to IPFS and returns its CID.
func (c *IpfsClient) AddFile(ctx context.Context, msg proto.Message) (string, error) {
	data, err := proto.Marshal(msg)
	if err != nil {
		return "", fmt.Errorf("failed to marshal protobuf message: %w", err)
	}

	file := files.NewReaderFile(bytes.NewReader(data))

	cid, err := c.NodeHttpApi.Unixfs().Add(ctx, file)
	if err != nil {
		return "", fmt.Errorf("failed to add file to IPFS: %w", err)
	}

	return cid.String(), nil
}

// GetFile retrieves a protobuf message from IPFS, unmarshels it and leaves the result in msg.
func (c *IpfsClient) GetFile(ctx context.Context, cid string, msg proto.Message) error {
	ipfsPath, err := path.NewPath(cid)
	if err != nil {
		return fmt.Errorf("invalid CID path: %w", err)
	}

	node, err := c.NodeHttpApi.Unixfs().Get(ctx, ipfsPath)
	if err != nil {
		return fmt.Errorf("failed to get file from IPFS: %w", err)
	}

	file, ok := node.(files.File)
	if !ok {
		return fmt.Errorf("unexpected node type: %T", node)
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read IPFS file: %w", err)
	}

	if err := proto.Unmarshal(data, msg); err != nil {
		return fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	return nil
}

// PinFile pins a CID to the local IPFS node. Without pinning, the data related
// to the CID will be stored in the IPFS but will get deleted by the garbage collector later on.
// Pinning the CID will prevent this from happening.
func (c *IpfsClient) PinFile(ctx context.Context, cid string) error {
	ipfsPath, err := path.NewPath(cid)
	if err != nil {
		return fmt.Errorf("invalid CID path: %w", err)
	}

	err = c.NodeHttpApi.Pin().Add(ctx, ipfsPath)
	if err != nil {
		return fmt.Errorf("failed to pin CID: %w", err)
	}

	return nil
}

// UnpinFile removes a pin for a CID from the local IPFS node, letting the garbage collector
// delete the data related to the CID.
func (c *IpfsClient) UnpinFile(ctx context.Context, cid string) error {
	ipfsPath, err := path.NewPath(cid)
	if err != nil {
		return fmt.Errorf("invalid CID path: %w", err)
	}

	err = c.NodeHttpApi.Pin().Rm(ctx, ipfsPath)
	if err != nil {
		return fmt.Errorf("failed to unpin CID from IPFS node: %w", err)
	}

	return nil
}

// AddAndPinFile adds a protobuf message to IPFS and pins it.
func (c *IpfsClient) AddAndPinFile(ctx context.Context, msg proto.Message) (string, error) {
	cid, err := c.AddFile(ctx, msg)
	if err != nil {
		return "", fmt.Errorf("failed to add file to IPFS: %w", err)
	}

	err = c.PinFile(ctx, cid)
	if err != nil {
		return "", fmt.Errorf("failed to pin CID: %w", err)
	}

	return cid, nil
}
