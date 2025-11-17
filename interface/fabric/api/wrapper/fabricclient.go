package fabricclient

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"

	"github.com/thcrull/fabric-ipfs-interface/interface/fabric/utils"
	"google.golang.org/grpc"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-protos-go-apiv2/common"
	"github.com/hyperledger/fabric-protos-go-apiv2/msp"

	"github.com/thcrull/fabric-ipfs-interface/interface/fabric/api/config"
)

// FabricClient is a wrapper around the Fabric Gateway client. It provides
// convenient methods for interacting with the Fabric network.
type FabricClient struct {
	Gateway  *client.Gateway
	Network  *client.Network
	Contract *client.Contract
	conn     *grpc.ClientConn
}

// NewFabricClient creates a new FabricClient instance by connecting to the Fabric Gateway.
// It sets up the gRPC connection, loads the client identity and signer,
// and prepares the network and contract for interaction.
// Returns an error if any of these steps fail.
func NewFabricClient(cfg *fabricconfig.FabricConfig) (*FabricClient, error) {
	conn, err := fabricutils.NewGrpcConnection(cfg)
	if err != nil {
		return nil, err
	}

	id, err := fabricutils.NewIdentity(cfg)
	if err != nil {
		conn.Close()
		return nil, err
	}

	sign, err := fabricutils.NewSign(cfg)
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

// EvaluateTransaction evaluates a transaction without modifying the ledger state. Used for querying the ledger.
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

func parseCert(raw []byte) (*x509.Certificate, error) {
	if cert, err := x509.ParseCertificate(raw); err == nil {
		return cert, nil
	}
	block, _ := pem.Decode(raw)
	if block == nil {
		return nil, fmt.Errorf("failed to parse certificate bytes")
	}
	return x509.ParseCertificate(block.Bytes)
}

type TxCreatorInfo struct {
	TxID     string
	BlockNum uint64
	MSPID    string
	Cert     *x509.Certificate
}

// GetTransactionCreator retrieves the creator identity of a given transaction
// by scanning committed blocks starting at the provided block number.
//
// txID:       the transaction ID to search for
// startBlock: usually 0 unless you want faster lookups
//
// Returns: MSP ID + parsed X.509 certificate + block number.
func (c *FabricClient) GetTransactionCreator(ctx context.Context, txID string, startBlock uint64) (*TxCreatorInfo, error) {
	blocks, err := c.Network.BlockEvents(ctx, client.WithStartBlock(startBlock))
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to block events: %w", err)
	}

	for block := range blocks {
		if block == nil || block.Data == nil {
			continue
		}

		for _, envBytes := range block.Data.Data {
			// Envelope
			env := &common.Envelope{}
			if err := proto.Unmarshal(envBytes, env); err != nil {
				// skip malformed envelope
				continue
			}

			// Payload
			payload := &common.Payload{}
			if err := proto.Unmarshal(env.Payload, payload); err != nil {
				continue
			}

			hdr := payload.GetHeader()
			if hdr == nil {
				continue
			}

			// Channel Header → txID
			ch := &common.ChannelHeader{}
			if err := proto.Unmarshal(hdr.GetChannelHeader(), ch); err != nil {
				continue
			}

			if ch.GetTxId() != txID {
				continue
			}

			// FOUND transaction → extract creator
			sigHdrBytes := hdr.GetSignatureHeader()
			if sigHdrBytes == nil {
				return nil, fmt.Errorf("signature header missing for transaction %s", txID)
			}

			sigHdr := &common.SignatureHeader{}
			if err := proto.Unmarshal(sigHdrBytes, sigHdr); err != nil {
				return nil, fmt.Errorf("failed to unmarshal SignatureHeader for tx %s: %w", txID, err)
			}

			creatorBytes := sigHdr.GetCreator()
			if creatorBytes == nil {
				return nil, fmt.Errorf("creator identity missing in SignatureHeader for %s", txID)
			}

			sid := &msp.SerializedIdentity{}
			if err := proto.Unmarshal(creatorBytes, sid); err != nil {
				return nil, fmt.Errorf("failed to unmarshal SerializedIdentity: %w", err)
			}

			cert, err := parseCert(sid.IdBytes)
			if err != nil {
				return nil, fmt.Errorf("invalid creator certificate: %w", err)
			}

			// Return creator information
			return &TxCreatorInfo{
				TxID:     txID,
				BlockNum: block.Header.Number,
				MSPID:    sid.Mspid,
				Cert:     cert,
			}, nil
		}
	}

	return nil, fmt.Errorf("transaction %s not found", txID)
}
