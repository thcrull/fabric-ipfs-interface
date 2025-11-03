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

// Participant holds participant's information.
// ParticipantId - the participant's id.
// EncapsulatedKey - the participant's encapsulated key.
type Participant struct {
	ParticipantId   string `json:"participant_id"`
	EncapsulatedKey string `json:"encap_key"`
}
