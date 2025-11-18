package fabricclient

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/thcrull/fabric-ipfs-interface/interface/fabric/api/config"
	"github.com/thcrull/fabric-ipfs-interface/shared"
)

// MetadataService wraps a Fabric client and provides methods
// for interacting with the metadata chaincode.
type MetadataService struct {
	client *FabricClient
}

// NewMetadataService creates a service for metadata transactions
func NewMetadataService(cfg *fabricconfig.FabricConfig) (*MetadataService, error) {
	fabricClient, err := NewFabricClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating metadata service: %w", err)
	}

	return &MetadataService{client: fabricClient}, nil
}

// ---------------------------------------------------------------------------
// THIS SECTION IS FOR THE PARTICIPANT MODEL METADATA RECORDS' FUNCTIONALITIES
// ---------------------------------------------------------------------------

// AddParticipantModelMetadata submits a transaction to add a new metadata record for a participant's model
func (s *MetadataService) AddParticipantModelMetadata(participantModelMetadata *shared.ParticipantModelMetadata) error {
	epochStr := strconv.Itoa(participantModelMetadata.Epoch)

	err := s.client.SubmitTransaction(nil, "AddParticipantModelMetadata", epochStr, participantModelMetadata.ParticipantId, participantModelMetadata.ModelHashCid, participantModelMetadata.HomomorphicHash)
	if err != nil {
		return fmt.Errorf("failed to add participant model metadata record: %w", err)
	}

	return nil
}

// GetParticipantModelMetadata retrieves a participant model metadata record by epoch and participantId
func (s *MetadataService) GetParticipantModelMetadata(epoch int, participantId string) (*shared.ParticipantModelMetadata, error) {
	epochStr := strconv.Itoa(epoch)
	var participantModelMetadata shared.ParticipantModelMetadata

	err := s.client.EvaluateTransaction(&participantModelMetadata, "GetParticipantModelMetadata", epochStr, participantId)
	if err != nil {
		return nil, fmt.Errorf("failed to query the participant model metadata record: %w", err)
	}

	return &participantModelMetadata, nil
}

// ParticipantModelMetadataExists returns true if a participant model metadata record exists
func (s *MetadataService) ParticipantModelMetadataExists(epoch int, participantId string) (bool, error) {
	epochStr := strconv.Itoa(epoch)
	var exists bool

	err := s.client.EvaluateTransaction(&exists, "ParticipantModelMetadataExists", epochStr, participantId)
	if err != nil {
		return false, fmt.Errorf("failed to query if participant model metadata record exists: %w", err)
	}

	return exists, nil
}

// DeleteParticipantModelMetadata deletes a participant model metadata record, returns nil if successful
func (s *MetadataService) DeleteParticipantModelMetadata(epoch int, participantId string) error {
	epochStr := strconv.Itoa(epoch)

	err := s.client.SubmitTransaction(nil, "DeleteParticipantModelMetadata", epochStr, participantId)
	if err != nil {
		return fmt.Errorf("failed to delete participant model metadata record: %w", err)
	}

	return nil
}

// UpdateParticipantModelMetadata updates an existing participant model metadata record
func (s *MetadataService) UpdateParticipantModelMetadata(participantModelMetadata *shared.ParticipantModelMetadata) error {
	epochStr := strconv.Itoa(participantModelMetadata.Epoch)

	err := s.client.SubmitTransaction(nil, "UpdateParticipantModelMetadata", epochStr, participantModelMetadata.ParticipantId, participantModelMetadata.ModelHashCid, participantModelMetadata.HomomorphicHash)
	if err != nil {
		return fmt.Errorf("failed to update participant model metadata record: %w", err)
	}

	return nil
}

// DeleteAllParticipantModelMetadata deletes all participant model metadata records, returns nil if successful
func (s *MetadataService) DeleteAllParticipantModelMetadata() error {
	err := s.client.SubmitTransaction(nil, "DeleteAllParticipantModelMetadata")
	if err != nil {
		return fmt.Errorf("failed to delete all participant model metadata records: %w", err)
	}

	return nil
}

// GetAllParticipantModelMetadata queries all participant model metadata records from the ledger made by any participant on any epoch
func (s *MetadataService) GetAllParticipantModelMetadata() ([]shared.ParticipantModelMetadata, error) {
	var participantModelMetadataList []shared.ParticipantModelMetadata

	err := s.client.EvaluateTransaction(&participantModelMetadataList, "GetAllParticipantModelMetadata")
	if err != nil {
		return nil, fmt.Errorf("failed to query all participant model metadata records: %w", err)
	}
	return participantModelMetadataList, nil
}

// GetAllParticipantModelMetadataByParticipant queries all participant model metadata records from the ledger made by the participant for any epoch
func (s *MetadataService) GetAllParticipantModelMetadataByParticipant(participantId string) ([]shared.ParticipantModelMetadata, error) {
	var participantModelMetadataList []shared.ParticipantModelMetadata

	err := s.client.EvaluateTransaction(&participantModelMetadataList, "GetAllParticipantModelMetadataByParticipant", participantId)
	if err != nil {
		return nil, fmt.Errorf("failed to query the participant model metadata records: %w", err)
	}

	return participantModelMetadataList, nil
}

// GetAllParticipantModelMetadataByEpoch queries all participant model metadata records from the ledger made by any participants for the epoch
func (s *MetadataService) GetAllParticipantModelMetadataByEpoch(epoch int) ([]shared.ParticipantModelMetadata, error) {
	epochStr := strconv.Itoa(epoch)
	var participantModelMetadataList []shared.ParticipantModelMetadata

	err := s.client.EvaluateTransaction(&participantModelMetadataList, "GetAllParticipantModelMetadataByEpoch", epochStr)
	if err != nil {
		return nil, fmt.Errorf("failed to query the participant model metadata records: %w", err)
	}

	return participantModelMetadataList, nil
}

// ---------------------------------------------------------------------------
// THIS SECTION IS FOR THE AGGREGATOR MODEL METADATA RECORDS' FUNCTIONALITIES
// ---------------------------------------------------------------------------

// AddAggregatorModelMetadata submits a transaction to add a new aggregator model metadata record
func (s *MetadataService) AddAggregatorModelMetadata(aggregatorModelMetadata *shared.AggregatorModelMetadata) error {
	epochStr := strconv.Itoa(aggregatorModelMetadata.Epoch)
	var participantIdsJSON, err = json.Marshal(aggregatorModelMetadata.ParticipantIds)
	if err != nil {
		return fmt.Errorf("failed to marshal participant ids JSON: %w", err)
	}

	err = s.client.SubmitTransaction(nil, "AddAggregatorModelMetadata", epochStr, aggregatorModelMetadata.ModelHashCid, string(participantIdsJSON))
	if err != nil {
		return fmt.Errorf("failed to add aggregator model metadata record: %w", err)
	}

	return nil
}

// GetAggregatorModelMetadata retrieves an aggregator model metadata record by epoch
func (s *MetadataService) GetAggregatorModelMetadata(epoch int) (*shared.AggregatorModelMetadata, error) {
	epochStr := strconv.Itoa(epoch)
	var aggregatorModelMetadata shared.AggregatorModelMetadata

	err := s.client.EvaluateTransaction(&aggregatorModelMetadata, "GetAggregatorModelMetadata", epochStr)
	if err != nil {
		return nil, fmt.Errorf("failed to query the aggregator model metadata record: %w", err)
	}

	return &aggregatorModelMetadata, nil
}

// AggregatorModelMetadataExists returns true if an aggregator model metadata record exists
func (s *MetadataService) AggregatorModelMetadataExists(epoch int) (bool, error) {
	epochStr := strconv.Itoa(epoch)
	var exists bool

	err := s.client.EvaluateTransaction(&exists, "AggregatorModelMetadataExists", epochStr)
	if err != nil {
		return false, fmt.Errorf("failed to query if aggregator model metadata record exists: %w", err)
	}

	return exists, nil
}

// DeleteAggregatorModelMetadata deletes an aggregator model metadata record, returns nil if successful
func (s *MetadataService) DeleteAggregatorModelMetadata(epoch int) error {
	epochStr := strconv.Itoa(epoch)

	err := s.client.SubmitTransaction(nil, "DeleteAggregatorModelMetadata", epochStr)
	if err != nil {
		return fmt.Errorf("failed to delete aggregator model metadata record: %w", err)
	}

	return nil
}

// UpdateAggregatorModelMetadata updates an existing aggregator model metadata record
func (s *MetadataService) UpdateAggregatorModelMetadata(aggregatorModelMetadata *shared.AggregatorModelMetadata) error {
	epochStr := strconv.Itoa(aggregatorModelMetadata.Epoch)
	var participantIdsJSON, err = json.Marshal(aggregatorModelMetadata.ParticipantIds)
	if err != nil {
		return fmt.Errorf("failed to marshal participant ids JSON: %w", err)
	}

	err = s.client.SubmitTransaction(nil, "UpdateAggregatorModelMetadata", epochStr, aggregatorModelMetadata.ModelHashCid, string(participantIdsJSON))
	if err != nil {
		return fmt.Errorf("failed to update aggregator model metadata record: %w", err)
	}

	return nil
}

// DeleteAllAggregatorModelMetadata deletes all aggregator model metadata records, returns nil if successful
func (s *MetadataService) DeleteAllAggregatorModelMetadata() error {
	err := s.client.SubmitTransaction(nil, "DeleteAllAggregatorModelMetadata")
	if err != nil {
		return fmt.Errorf("failed to delete all aggregator model metadata records: %w", err)
	}

	return nil
}

// GetAllAggregatorModelMetadata queries the aggregator model metadata records from the ledger for all epochs
func (s *MetadataService) GetAllAggregatorModelMetadata() ([]shared.AggregatorModelMetadata, error) {
	var aggregatorModelMetadataList []shared.AggregatorModelMetadata

	err := s.client.EvaluateTransaction(&aggregatorModelMetadataList, "GetAllAggregatorModelMetadata")
	if err != nil {
		return nil, fmt.Errorf("failed to query all aggregator model metadata records: %w", err)
	}
	return aggregatorModelMetadataList, nil
}

// ---------------------------------------------------------------------------
// THIS SECTION IS FOR THE PARTICIPANT'S FUNCTIONALITIES
// ---------------------------------------------------------------------------

// AddParticipant submits a transaction to add a new participant record
func (s *MetadataService) AddParticipant(participant *shared.Participant) error {
	err := s.client.SubmitTransaction(nil, "AddParticipant", participant.ParticipantId, participant.EncapsulatedKey)
	if err != nil {
		return fmt.Errorf("failed to add participant record: %w", err)
	}

	return nil
}

// GetParticipant retrieves a participant record by id
func (s *MetadataService) GetParticipant(participantId string) (*shared.Participant, error) {
	var participant shared.Participant

	err := s.client.EvaluateTransaction(&participant, "GetParticipant", participantId)
	if err != nil {
		return nil, fmt.Errorf("failed to query the participant record: %w", err)
	}

	return &participant, nil
}

// ParticipantExists returns true if a participant record exists
func (s *MetadataService) ParticipantExists(participantId string) (bool, error) {
	var exists bool

	err := s.client.EvaluateTransaction(&exists, "ParticipantExists", participantId)
	if err != nil {
		return false, fmt.Errorf("failed to query if participant record exists: %w", err)
	}

	return exists, nil
}

// DeleteParticipant deletes a participant record, returns nil if successful
func (s *MetadataService) DeleteParticipant(participantId string) error {
	err := s.client.SubmitTransaction(nil, "DeleteParticipant", participantId)
	if err != nil {
		return fmt.Errorf("failed to delete participant record: %w", err)
	}

	return nil
}

// UpdateParticipant updates a existing participant record
func (s *MetadataService) UpdateParticipant(participant *shared.Participant) error {
	err := s.client.SubmitTransaction(nil, "UpdateParticipant", participant.ParticipantId, participant.EncapsulatedKey)
	if err != nil {
		return fmt.Errorf("failed to update participant record: %w", err)
	}

	return nil
}

// DeleteAllParticipants deletes all participant records, returns nil if successful
func (s *MetadataService) DeleteAllParticipants() error {
	err := s.client.SubmitTransaction(nil, "DeleteAllParticipants")
	if err != nil {
		return fmt.Errorf("failed to delete all participants records: %w", err)
	}

	return nil
}

// GetAllParticipants queries all participant records from the ledger
func (s *MetadataService) GetAllParticipants() ([]shared.Participant, error) {
	var participantsList []shared.Participant

	err := s.client.EvaluateTransaction(&participantsList, "GetAllParticipants")
	if err != nil {
		return nil, fmt.Errorf("failed to query all participant records: %w", err)
	}
	return participantsList, nil
}

// --------------------------------------------
// THIS SECTION DEALS WITH ACCESSING THE LOGS
// --------------------------------------------

// GetAllLogsWithoutCreator returns all logs from the ledger without the creator information.
// This is much faster than GetAllLogs() since it doesn't need to parse the ledger N times for the creators.
func (s *MetadataService) GetAllLogsWithoutCreator() ([]shared.LogEntry, error) {
	var history []shared.LogEntry

	err := s.client.EvaluateTransaction(&history, "GetAllLogs")
	if err != nil {
		return nil, fmt.Errorf("failed to query all logs: %w", err)
	}

	return history, nil
}

// GetAllLogs returns all logs from the ledger with the creator information.
// This might be slow for extremely large ledgers.
func (s *MetadataService) GetAllLogs() ([]shared.LogEntry, error) {
	var history []shared.LogEntry

	err := s.client.EvaluateTransaction(&history, "GetAllLogs")
	if err != nil {
		return nil, fmt.Errorf("failed to query all logs: %w", err)
	}

	for i, log := range history {
		var found, creatorInfo, err = s.client.GetTransactionCreator(context.Background(), log.TxID, 0)
		if err != nil {
			return nil, fmt.Errorf("failed to get transaction creator: %w", err)
		}

		if found {
			log.TxCreator = shared.TxCreatorInfo{
				TxID:     creatorInfo.TxID,
				BlockNum: creatorInfo.BlockNum,
				MSPID:    creatorInfo.MSPID,
				Cert:     creatorInfo.Cert,
			}
			history[i] = log
		}
	}

	return history, nil
}

func (s *MetadataService) Close() error {
	return s.client.Close()
}
