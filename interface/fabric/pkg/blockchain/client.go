package blockchain

import (
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"google.golang.org/grpc"

	"github.com/thcrull/fabric-interface/application/internal/bcutils"
	"github.com/thcrull/fabric-interface/application/pkg/config"
)

// Client is a wrapper around a Fabric Gateway connection, providing
// access to a specific network and smart contract for submitting and querying transactions.
type Client struct {
	Gateway  *client.Gateway
	Network  *client.Network
	Contract *client.Contract
	conn     *grpc.ClientConn
}

// NewClient creates a new Client instance by connecting to the Fabric Gateway.
// It sets up the gRPC connection, loads the client identity and signer,
// and prepares the network and contract for interaction.
// Returns an error if any of these steps fail.
func NewClient(cfg *config.Config) (*Client, error) {
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

	return &Client{
		Gateway:  gw,
		Network:  network,
		Contract: contract,
		conn:     conn,
	}, nil
}

// SubmitTransaction submits a transaction that modifies the ledger state.
// `name` is the chaincode function name and `args` are its parameters.
// Returns the transaction result or an error.
func (c *Client) SubmitTransaction(name string, args ...string) ([]byte, error) {
	return c.Contract.SubmitTransaction(name, args...)
}

// EvaluateTransaction evaluates (queries) a transaction without modifying the ledger state.
// `name` is the chaincode function name and `args` are its parameters.
// Returns the query result or an error.
func (c *Client) EvaluateTransaction(name string, args ...string) ([]byte, error) {
	return c.Contract.EvaluateTransaction(name, args...)
}

// Close cleans up the Client by closing the Gateway and gRPC connection.
// Returns an error if closing any resource fails.
func (c *Client) Close() error {
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
