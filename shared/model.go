package shared

// UserInfo holds information about a user of the ledger.
// TxID - the transaction id.
// MSPID - the MSP ID of the user.
// SerialNumber - the serial number of the user.
// CommonName - the common name of the user.
// OrganizationalUnit - the organisational unit of the user. A set of roles and memberships within the organisation.
type UserInfo struct {
	MSPID              string   `json:"mspId"`
	SerialNumber       string   `json:"serialNumber"`
	CommonName         string   `json:"commonName"`
	OrganizationalUnit []string `json:"organizationalUnit"`
}

// Participant holds participant's information.
// ParticipantId - the participant's id.
// EncapsulatedKey - the participant's encapsulated key.
// HomomorphicSharedKeyCypher - the homomorphic shared key cypher. This key can be retrieved using the decapsulated key, and it is used to encrypt the model update.
// CommunicationKeyCypher - the participant's communication key cypher. This key is used to encrypt the messages exchanged between the participant and the aggregator.
type Participant struct {
	ParticipantId              string `json:"participant_id"`
	EncapsulatedKey            string `json:"encap_key"`
	HomomorphicSharedKeyCypher string `json:"homomorphic_key_cypher"`
	CommunicationKeyCypher     string `json:"comm_key_cypher"`
}

// Aggregator holds aggregator's information.
// AggregatorId - the aggregator's id.
// CommunicationKeyCyphers - a map of communication key cyphers, where the key is the participant's id. These keys are used to encrypt the messages exchanged between the aggregator and the participant.
type Aggregator struct {
	AggregatorId             string            `json:"aggregator_id"`
	CommunicationKeysCyphers map[string]string `json:"comm_keys_cyphers"`
}

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
// AggregatorId - the aggregator's id.
// Epoch - the epoch of the model update.
// ParticipantIds - the participants' ids that contributed to the global model update.
// ModelHashCid - the IPFS CID of the global model update.
type AggregatorModelMetadata struct {
	Epoch          int      `json:"epoch"`
	AggregatorId   string   `json:"aggregator_id"`
	ParticipantIds []string `json:"participant_ids"`
	ModelHashCid   string   `json:"model_hash_cid"`
}

// LogEntry holds a transaction log entry.
// TxID - the transaction id.
// TxCreator - information about the creator of the transaction.
// Timestamp - the timestamp of the transaction.
// IsDelete - whether the entry is a delete operation.
// Value - the value of the entry. This is usually the changes made by the transaction.
type LogEntry struct {
	TxID      string      `json:"txId"`
	TxCreator UserInfo    `json:"txCreator"`
	Timestamp string      `json:"timestamp"`
	IsDelete  bool        `json:"isDelete"`
	Changes   interface{} `json:"value"`
}
