package shared

import (
	"crypto/x509"
)

// ParticipantModelMetadata represents a participant's model update's metadata.
// Epoch - the epoch of the model update.
// ParticipantId - the participant's id.
// ModelHashCid - the IPFS CID of the model update.
// HomomorphicHash - the homomorphic hash of the model update.
type ParticipantModelMetadata struct {
	Epoch           int    `json:"epoch"`
	ParticipantId   string `json:"participant_id"`
	ModelHashCid    string `json:"model_hash_cid"`
	HomomorphicHash string `json:"homomorphic_hash"`
}

// AggregatorModelMetadata represents an aggregator's global model update's metadata.
// Epoch - the epoch of the model update.
// ParticipantIds - the participants' ids that contributed to the global model update.
// ModelHashCid - the IPFS CID of the global model update.
type AggregatorModelMetadata struct {
	Epoch          int      `json:"epoch"`
	ParticipantIds []string `json:"participant_ids"`
	ModelHashCid   string   `json:"model_hash_cid"`
}

// Participant holds participant's information.
// ParticipantId - the participant's id.
// EncapsulatedKey - the participant's encapsulated key.
type Participant struct {
	ParticipantId   string `json:"participant_id"`
	EncapsulatedKey string `json:"encap_key"`
}

// TxCreatorInfo holds information about the creator of a transaction.
// TxID - the transaction id.
// BlockNum - the block number where the transaction was committed.
// MSPID - the MSP id of the creator.
// Cert - the creator's certificate.'
type TxCreatorInfo struct {
	TxID     string           `json:"txId"`
	BlockNum uint64           `json:"blockNum"`
	MSPID    string           `json:"mspId"`
	Cert     x509.Certificate `json:"cert"`
}

// LogEntry holds a transaction log entry.
// TxID - the transaction id.
// TxCreator - information about the creator of the transaction.
// Timestamp - the timestamp of the transaction.
// IsDelete - whether the entry is a delete operation.
// Value - the value of the entry. This is usually the changes made by the transaction.
type LogEntry struct {
	TxID      string        `json:"txId"`
	Timestamp string        `json:"timestamp"`
	IsDelete  bool          `json:"isDelete"`
	Value     interface{}   `json:"value"`
	TxCreator TxCreatorInfo `json:"txCreator"`
}
