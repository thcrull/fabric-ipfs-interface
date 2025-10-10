package ipfsclient

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/ipfs/boxo/files"
	"github.com/ipfs/boxo/path"
	"github.com/ipfs/kubo/client/rpc"
	"github.com/thcrull/fabric-ipfs-interface/interface/ipfs/pkg/config"
	"google.golang.org/protobuf/proto"
)

type IpfsClient struct {
	httpClient  *http.Client
	NodeHttpApi *rpc.HttpApi
}

func NewIpfsClient(cfg *config.IpfsConfig) (*IpfsClient, error) {
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

// AddFile uploads a protobuf message to IPFS and returns the CID (not pinned).
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

// GetFile retrieves a file from IPFS by CID and unmarshals it into a protobuf message.
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

// PinFile explicitly pins a CID to the local IPFS node.
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

// UnpinFile removes a pin for a CID from the local IPFS node.
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
