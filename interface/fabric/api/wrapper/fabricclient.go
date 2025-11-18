package fabricclient

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"

	"github.com/thcrull/fabric-ipfs-interface/interface/fabric/utils"
	"github.com/thcrull/fabric-ipfs-interface/shared"
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

// GetTransactionCreator retrieves the creator identity of a given transaction
// by scanning committed blocks starting at the provided block number.
// TxID - the transaction ID to look for.
// StartBlock - the block number to start scanning from. Leave 0 to scan the whole ledger. Can be used for faster searching.
func (c *FabricClient) GetTransactionCreator(ctx context.Context, txID string, startBlock uint64) (bool, *shared.TxCreatorInfo, error) {
	// Fast scan the ledger for the block that contains the transaction.
	found, targetBlock, err := c.findTxBlock(ctx, txID, startBlock)
	if err != nil {
		return false, nil, err
	}

	if !found {
		return false, nil, nil
	}

	// Make the context cancellable so we can stop listening to the filtered block events when we find the target block.
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	blocks, err := c.Network.BlockEvents(ctx, client.WithStartBlock(targetBlock))
	if err != nil {
		return false, nil, fmt.Errorf("failed to subscribe to block events: %w", err)
	}

	block := <-blocks

	if block == nil || block.Data == nil {
		return false, nil, fmt.Errorf("block %d is corrupted", targetBlock)
	}

	for _, envBytes := range block.Data.Data {
		// We unmarshal the data into an Envelope which is the payload wrapped around with the signature.
		env := &common.Envelope{}
		if err := proto.Unmarshal(envBytes, env); err != nil {
			// Skip malformed envelope
			continue
		}

		// We unmarshal the Payload from the Envelope and get its Header.
		payload := &common.Payload{}
		if err := proto.Unmarshal(env.Payload, payload); err != nil {
			continue
		}

		// This contains the metadata about the transaction, such as the creator.
		hdr := payload.GetHeader()
		if hdr == nil {
			continue
		}

		// We unmarshal the ChannelHeader from the Header of the Payload.
		// We do this to check if the transaction ID matches the one we're looking for.
		ch := &common.ChannelHeader{}
		if err := proto.Unmarshal(hdr.GetChannelHeader(), ch); err != nil {
			continue
		}

		if ch.GetTxId() != txID {
			continue
		}

		// If we are in the right transaction, we can extract the creator identity from the SignatureHeader.
		sigHdrBytes := hdr.GetSignatureHeader()
		if sigHdrBytes == nil {
			return false, nil, fmt.Errorf("signature header missing for transaction %s", txID)
		}

		sigHdr := &common.SignatureHeader{}
		if err := proto.Unmarshal(sigHdrBytes, sigHdr); err != nil {
			return false, nil, fmt.Errorf("failed to unmarshal SignatureHeader for tx %s: %w", txID, err)
		}

		creatorBytes := sigHdr.GetCreator()
		if creatorBytes == nil {
			return false, nil, fmt.Errorf("creator identity missing in SignatureHeader for %s", txID)
		}

		// Parse the creator identity bytes into a SerializedIdentity.
		sid := &msp.SerializedIdentity{}
		if err := proto.Unmarshal(creatorBytes, sid); err != nil {
			return false, nil, fmt.Errorf("failed to unmarshal SerializedIdentity: %w", err)
		}

		// Parse the certificate bytes into a certificate.
		cert, err := parseCert(sid.IdBytes)
		if err != nil {
			return false, nil, fmt.Errorf("invalid creator certificate: %w", err)
		}

		// Return creator information
		return true, &shared.TxCreatorInfo{
			TxID:     txID,
			BlockNum: block.Header.Number,
			MSPID:    sid.Mspid,
			Cert:     cert,
		}, nil
	}

	// This is an error because it should be impossible to reach.
	return false, nil, fmt.Errorf("transaction %s not found, but present in the filtered blocks", txID)
}

// findTxBlock scans the ledger for the block that contains the given transaction.
// txId - the transaction ID to look for.
// startBlock - the block number to start scanning from. Leave 0 to scan the whole ledger. Can be used for faster searching.
func (c *FabricClient) findTxBlock(ctx context.Context, txID string, startBlock uint64) (bool, uint64, error) {
	// Make the context cancellable so we can stop listening to the filtered block events when we find the target block.
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Fetch the chain info for the height.
	bytes, err := c.Network.GetContract("qscc").EvaluateTransaction("GetChainInfo", c.Network.Name())
	if err != nil {
		return false, 0, fmt.Errorf("failed to get chain info: %w", err)
	}

	var info common.BlockchainInfo
	if err := proto.Unmarshal(bytes, &info); err != nil {
		return false, 0, fmt.Errorf("failed to unmarshal chain info: %w", err)
	}

	chainHeight := info.Height

	// We scan filtered blocks and not normal ones because it is much faster and efficient.
	// The filtered block and filtered transactions only contain metadata, not the entire payload.
	filteredBlocks, err := c.Network.FilteredBlockEvents(ctx, client.WithStartBlock(startBlock))
	if err != nil {
		return false, 0, fmt.Errorf("failed to subscribe to filtered block events: %w", err)
	}

	var targetBlock uint64 = 0
	var found = false
	for filteredBlock := range filteredBlocks {
		if filteredBlock == nil {
			continue
		}

		for _, filteredTransaction := range filteredBlock.FilteredTransactions {
			if filteredTransaction.Txid == txID {
				targetBlock = filteredBlock.Number
				found = true
				break
			}
		}

		// Stop if found, or we reached the end of the ledger.
		// Needed because the FilteredBlockEvents function subscribes us to the latest block and does not stop until cancelled.
		if found || filteredBlock.Number >= chainHeight-1 {
			break
		}
	}

	return found, targetBlock, nil
}

// parseCert is a helper function that parses a certificate from its PEM representation.
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
