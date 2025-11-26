package shared

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

// 1. Add a cypher for the key used between aggregator and participants.

// Participant holds participant's information.
// ParticipantId - the participant's id.
// EncapsulatedKey - the participant's encapsulated key.
type Participant struct {
	ParticipantId   string `json:"participant_id"`
	EncapsulatedKey string `json:"encap_key"`
	// 2. add Homomorphic encryption key share cypher (this is used to encrypt the model update)
}

// TxCreatorInfo holds information about the creator of a transaction.
// TxID - the transaction id.
// MSPID - the MSP ID of the creator.
// CommonName - the common name of the creator.
// OrganizationalUnit - the organisational unit of the creator. A set of roles and memberships within the organisation.
// SerialNumber - the serial number of the creator.
type TxCreatorInfo struct {
	TxID               string   `json:"txId"`
	MSPID              string   `json:"mspId"`
	CommonName         string   `json:"commonName"`
	OrganizationalUnit []string `json:"organizationalUnit"`
	SerialNumber       uint64   `json:"serialNumber"`
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
