package shared

// UserInfo holds information about a user of the ledger.
// MSPID - the MSP id of the user.
// SerialNumber - the serial number of the user.
// CommonName - the common name of the user.
// OrganizationalUnit - the organisational unit of the user. A set of roles and memberships within the organisation.
type UserInfo struct {
	MSPID              string   `json:"msp_id"`
	SerialNumber       string   `json:"serial_number"`
	CommonName         string   `json:"common_name"`
	OrganizationalUnit []string `json:"organizational_unit"`
}

// Participant holds participant's information.
// ParticipantId - the participant's id.
// EncapsulatedKey - the participant's encapsulated key.
// HomomorphicSharedKeyCypher - the homomorphic shared key cypher. This key can be retrieved using the decapsulated key, and it is used to encrypt the model update.
// CommunicationKeyCypher - the participant's communication key cypher. This key is used to encrypt the messages exchanged between the participant and the aggregator.
// MSPID - the MSP id of the user who created the participant.
// SerialNumber - the serial number of the user who created the participant.
type Participant struct {
	ParticipantId              int    `json:"participant_id"`
	EncapsulatedKey            string `json:"encap_key"`              // STEP 2
	HomomorphicSharedKeyCypher string `json:"homomorphic_key_cypher"` //STEP 4
	CommunicationKeyCypher     string `json:"comm_key_cypher"`
	MSPID                      string `json:"msp_id"`
	SerialNumber               string `json:"serial_number"`
}

// Aggregator holds aggregator's information.
// AggregatorId - the aggregator's id.
// CommunicationKeysCyphers - a map of communication key cyphers, where the key is the participant's id. These keys are used to encrypt the messages exchanged between the aggregator and the participant.
// MSPID - the MSP id of the user who created the aggregator.
// SerialNumber - the serial number of the user who created the aggregator.
type Aggregator struct {
	AggregatorId             int            `json:"aggregator_id"`
	CommunicationKeysCyphers map[int]string `json:"comm_keys_cyphers"`
	MSPID                    string         `json:"msp_id"`
	SerialNumber             string         `json:"serial_number"`
}

// ParticipantModelMetadata represents a participant's model update's metadata.
// Epoch - the epoch of the model update.
// ParticipantId - the participant's id.
// ModelHashCid - the IPFS CID of the model update.
// HomomorphicHash - the homomorphic hash of the model update.
type ParticipantModelMetadata struct {
	Epoch           int    `json:"epoch"`
	ParticipantId   int    `json:"participant_id"`
	ModelHashCid    string `json:"model_hash_cid"`
	HomomorphicHash string `json:"homomorphic_hash"`
}

// AggregatorModelMetadata represents an aggregator's global model update's metadata.
// AggregatorId - the aggregator's id.
// Epoch - the epoch of the model update.
// ParticipantIds - the participants' ids that contributed to the global model update.
// ModelHashCid - the IPFS CID of the global model update.
type AggregatorModelMetadata struct {
	Epoch          int    `json:"epoch"`
	AggregatorId   int    `json:"aggregator_id"`
	ParticipantIds []int  `json:"participant_ids"`
	ModelHashCid   string `json:"model_hash_cid"`
}

// LogEntry holds a transaction log entry.
// TxId - the transaction id.
// TxCreator - information about the creator of the transaction.
// Timestamp - the timestamp of the transaction.
// IsDelete - whether the entry is a delete operation.
// Changes - the changes made by the transaction.
type LogEntry struct {
	TxId      string      `json:"tx_id"`
	TxCreator UserInfo    `json:"tx_creator"`
	Timestamp string      `json:"timestamp"`
	IsDelete  bool        `json:"is_delete"`
	Changes   interface{} `json:"changes"`
}
