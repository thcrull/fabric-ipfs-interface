package shared

// Metadata represents a participant's federated learning model update
type Metadata struct {
	Epoch           int    `json:"epoch"`
	ParticipantID   string `json:"participant_id"`
	EncapsulatedKey string `json:"encap_key"`
	EncModelHash    string `json:"enc_model_hash"`
	HomomorphicHash string `json:"homomorphic_hash"`
}

// WeightMatrix represents a matrix of weights
type WeightMatrix [][]float64
