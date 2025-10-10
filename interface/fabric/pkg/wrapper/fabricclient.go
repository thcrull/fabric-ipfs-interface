package metadata

import (
	"encoding/json"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/thcrull/fabric-ipfs-interface/interface/fabric/internal/bcutils"
	"google.golang.org/grpc"

	"github.com/thcrull/fabric-ipfs-interface/interface/fabric/pkg/config"
)

// Client is a wrapper around a Fabric Gateway connection, providing
// access to a specific network and smart contract for submitting and querying transactions.
type FabricClient struct {
	Gateway  *client.Gateway
	Network  *client.Network
	Contract *client.Contract
	conn     *grpc.ClientConn
}

// NewClient creates a new Client instance by connecting to the Fabric Gateway.
// It sets up the gRPC connection, loads the client identity and signer,
// and prepares the network and contract for interaction.
// Returns an error if any of these steps fail.
func NewClient(cfg *fabricconfig.FabricConfig) (*FabricClient, error) {
	conn, err := bcutils.NewGrpcConnection(cfg)
	if err != nil {
		return nil, err
	}

	id, err := bcutils.NewIdentity(cfg)
	if err != nil {
		conn.Close()
		return nil, err
	}

	sign, err := bcutils.NewSign(cfg)
	if err != nil {
		conn.Close()
		return nil, err
	}

	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(conn),
	)
	if err != nil {
		conn.Close()
		return nil, err
	}

	network := gw.GetNetwork(cfg.Network.ChannelName)
	contract := network.GetContract(cfg.Network.ChaincodeName)

	return &FabricClient{
		Gateway:  gw,
		Network:  network,
		Contract: contract,
		conn:     conn,
	}, nil
}

// SubmitTransaction submits a transaction that modifies the ledger state.
// Name is the chaincode function name, args are its parameters, and out is the output address.
// Returns the transaction result or an error.
func (c *FabricClient) SubmitTransaction(out interface{}, name string, args ...string) error {
	res, err := c.Contract.SubmitTransaction(name, args...)
	if err != nil {
		return err
	}

	if res == nil {
		return nil
	}

	return json.Unmarshal(res, out)
}

// EvaluateTransaction evaluates (queries) a transaction without modifying the ledger state.
// Name is the chaincode function name, args are its parameters, and out is the output address.
// Returns the query result or an error.
func (c *FabricClient) EvaluateTransaction(out interface{}, name string, args ...string) error {
	res, err := c.Contract.EvaluateTransaction(name, args...)
	if err != nil {
		return err
	}

	if res == nil {
		return nil
	}

	return json.Unmarshal(res, out)
}

// Close cleans up the Client by closing the Gateway and gRPC connection.
// Returns an error if closing any resource fails.
func (c *FabricClient) Close() error {
	err := c.Gateway.Close()
	if err != nil {
		return err
	}

	err = c.conn.Close()
	if err != nil {
		return err
	}

	return nil
}
