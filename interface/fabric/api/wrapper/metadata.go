package fabricclient

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/thcrull/fabric-ipfs-interface/shared"
)

// MetadataService wraps a Fabric client and provides methods
// for interacting with the metadata chaincode.
type MetadataService struct {
	client *FabricClient
}

// NewMetadataService creates a service for metadata transactions
func NewMetadataService(configPath string) (*MetadataService, error) {
	fabricClient, err := NewFabricClient(configPath)
	if err != nil {
		return nil, fmt.Errorf("error creating metadata service: %w", err)
	}

	return &MetadataService{client: fabricClient}, nil
}

// ---------------------------------------------------------------------------
// THIS SECTION IS FOR THE PARTICIPANT'S FUNCTIONALITIES
// ---------------------------------------------------------------------------

// AddParticipant submits a transaction to add a new participant record. The new participant will be bound to the caller of the transaction. It returns the participant id generated from the MSPID and SerialNumber of the caller.
func (s *MetadataService) AddParticipant(encapsulatedKey string, homomorphicSharedKeyCypher string, communicationKeyCypher string) (string, error) {
	var participantId string

	err := s.client.SubmitTransaction(&participantId, "AddParticipant", encapsulatedKey, homomorphicSharedKeyCypher, communicationKeyCypher)
	if err != nil {
		return "", fmt.Errorf("failed to add participant record: %w", err)
	}

	return participantId, nil
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

// DeleteParticipant deletes the caller's participant record, returns nil if successful
func (s *MetadataService) DeleteParticipant() error {
	err := s.client.SubmitTransaction(nil, "DeleteParticipant")
	if err != nil {
		return fmt.Errorf("failed to delete participant record: %w", err)
	}

	return nil
}

// UpdateParticipant updates the caller's participant record
func (s *MetadataService) UpdateParticipant(encapsulatedKey string, homomorphicSharedKeyCypher string, communicationKeyCypher string) error {
	err := s.client.SubmitTransaction(nil, "UpdateParticipant", encapsulatedKey, homomorphicSharedKeyCypher, communicationKeyCypher)
	if err != nil {
		return fmt.Errorf("failed to update participant record: %w", err)
	}

	return nil
}

// DeleteAllParticipants deletes all participant records, returns nil if successful. Only the admin can delete all participants records.
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

// -----------------------------------------------------
// THIS SECTION IS FOR THE AGGREGATOR'S FUNCTIONALITIES
// -----------------------------------------------------

// AddAggregator submits a transaction to add a new aggregator record. The new aggregator will be bound to the caller of the transaction. It returns the aggregator id generated from the MSPID and SerialNumber of the caller.
func (s *MetadataService) AddAggregator(communicationKeysCyphers map[string]string) (string, error) {
	var communicationKeysCyphersJSON, err = json.Marshal(communicationKeysCyphers)
	if err != nil {
		return "", fmt.Errorf("failed to marshal communication keys cyphers JSON: %w", err)
	}

	var aggregatorId string

	err = s.client.SubmitTransaction(&aggregatorId, "AddAggregator", string(communicationKeysCyphersJSON))
	if err != nil {
		return "", fmt.Errorf("failed to add aggregator record: %w", err)
	}

	return aggregatorId, nil
}

// GetAggregator retrieves an aggregator record by id
func (s *MetadataService) GetAggregator(aggregatorId string) (*shared.Aggregator, error) {
	var aggregator shared.Aggregator

	err := s.client.EvaluateTransaction(&aggregator, "GetAggregator", aggregatorId)
	if err != nil {
		return nil, fmt.Errorf("failed to query the aggregator record: %w", err)
	}

	return &aggregator, nil
}

// AggregatorExists returns true if an aggregator record exists
func (s *MetadataService) AggregatorExists(aggregatorId string) (bool, error) {
	var exists bool

	err := s.client.EvaluateTransaction(&exists, "AggregatorExists", aggregatorId)
	if err != nil {
		return false, fmt.Errorf("failed to query if aggregator record exists: %w", err)
	}

	return exists, nil
}

// DeleteAggregator deletes the caller's aggregator record, returns nil if successful
func (s *MetadataService) DeleteAggregator() error {
	err := s.client.SubmitTransaction(nil, "DeleteAggregator")
	if err != nil {
		return fmt.Errorf("failed to delete aggregator record: %w", err)
	}

	return nil
}

// UpdateAggregator updates the caller's aggregator record
func (s *MetadataService) UpdateAggregator(communicationKeysCyphers map[string]string) error {
	var communicationKeysCyphersJSON, err = json.Marshal(communicationKeysCyphers)
	if err != nil {
		return fmt.Errorf("failed to marshal communication keys cyphers JSON: %w", err)
	}

	err = s.client.SubmitTransaction(nil, "UpdateAggregator", string(communicationKeysCyphersJSON))
	if err != nil {
		return fmt.Errorf("failed to update aggregator record: %w", err)
	}

	return nil
}

// DeleteAllAggregators deletes all aggregator records, returns nil if successful. Only the admin can delete all aggregators records.
func (s *MetadataService) DeleteAllAggregators() error {
	err := s.client.SubmitTransaction(nil, "DeleteAllAggregators")
	if err != nil {
		return fmt.Errorf("failed to delete all aggregators records: %w", err)
	}

	return nil
}

// GetAllAggregators queries all aggregator records from the ledger
func (s *MetadataService) GetAllAggregators() ([]shared.Aggregator, error) {
	var aggregatorsList []shared.Aggregator

	err := s.client.EvaluateTransaction(&aggregatorsList, "GetAllAggregators")
	if err != nil {
		return nil, fmt.Errorf("failed to query all aggregators records: %w", err)
	}
	return aggregatorsList, nil
}

// ---------------------------------------------------------------------------
// THIS SECTION IS FOR THE PARTICIPANT MODEL METADATA RECORDS' FUNCTIONALITIES
// ---------------------------------------------------------------------------

// AddParticipantModelMetadata submits a transaction to add a new metadata record for a participant's model. The model can only be for the caller's participant.
func (s *MetadataService) AddParticipantModelMetadata(epoch int, modelHashCid string, homomorphicHash string) error {
	epochStr := strconv.Itoa(epoch)

	err := s.client.SubmitTransaction(nil, "AddParticipantModelMetadata", epochStr, modelHashCid, homomorphicHash)
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

// DeleteParticipantModelMetadata deletes a participant model metadata record, returns nil if successful. Can only delete records made by the caller's participant.
func (s *MetadataService) DeleteParticipantModelMetadata(epoch int) error {
	epochStr := strconv.Itoa(epoch)

	err := s.client.SubmitTransaction(nil, "DeleteParticipantModelMetadata", epochStr)
	if err != nil {
		return fmt.Errorf("failed to delete participant model metadata record: %w", err)
	}

	return nil
}

// UpdateParticipantModelMetadata updates an existing participant model metadata record. Can only update records made by the caller's participant.
func (s *MetadataService) UpdateParticipantModelMetadata(epoch int, modelHashCid string, homomorphicHash string) error {
	epochStr := strconv.Itoa(epoch)

	err := s.client.SubmitTransaction(nil, "UpdateParticipantModelMetadata", epochStr, modelHashCid, homomorphicHash)
	if err != nil {
		return fmt.Errorf("failed to update participant model metadata record: %w", err)
	}

	return nil
}

// DeleteAllParticipantModelMetadata deletes all participant model metadata records, returns nil if successful. Only the admin can delete all participant model metadata records.
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

// AddAggregatorModelMetadata submits a transaction to add a new aggregator model metadata record. The model can only be for the caller's aggregator.
func (s *MetadataService) AddAggregatorModelMetadata(epoch int, modelHashCid string, participantIds []string) error {
	epochStr := strconv.Itoa(epoch)
	var participantIdsJSON, err = json.Marshal(participantIds)
	if err != nil {
		return fmt.Errorf("failed to marshal participant ids JSON: %w", err)
	}

	err = s.client.SubmitTransaction(nil, "AddAggregatorModelMetadata", epochStr, modelHashCid, string(participantIdsJSON))
	if err != nil {
		return fmt.Errorf("failed to add aggregator model metadata record: %w", err)
	}

	return nil
}

// GetAggregatorModelMetadata retrieves an aggregator model metadata record by epoch
func (s *MetadataService) GetAggregatorModelMetadata(aggregatorId string, epoch int) (*shared.AggregatorModelMetadata, error) {
	epochStr := strconv.Itoa(epoch)
	var aggregatorModelMetadata shared.AggregatorModelMetadata

	err := s.client.EvaluateTransaction(&aggregatorModelMetadata, "GetAggregatorModelMetadata", epochStr, aggregatorId)
	if err != nil {
		return nil, fmt.Errorf("failed to query the aggregator model metadata record: %w", err)
	}

	return &aggregatorModelMetadata, nil
}

// AggregatorModelMetadataExists returns true if an aggregator model metadata record exists
func (s *MetadataService) AggregatorModelMetadataExists(epoch int, aggregatorId string) (bool, error) {
	epochStr := strconv.Itoa(epoch)
	var exists bool

	err := s.client.EvaluateTransaction(&exists, "AggregatorModelMetadataExists", epochStr, aggregatorId)
	if err != nil {
		return false, fmt.Errorf("failed to query if aggregator model metadata record exists: %w", err)
	}

	return exists, nil
}

// DeleteAggregatorModelMetadata deletes an aggregator model metadata record, returns nil if successful. Can only delete records made by the caller's aggregator.
func (s *MetadataService) DeleteAggregatorModelMetadata(epoch int) error {
	epochStr := strconv.Itoa(epoch)

	err := s.client.SubmitTransaction(nil, "DeleteAggregatorModelMetadata", epochStr)
	if err != nil {
		return fmt.Errorf("failed to delete aggregator model metadata record: %w", err)
	}

	return nil
}

// UpdateAggregatorModelMetadata updates an existing aggregator model metadata record. Can only update records made by the caller's aggregator.
func (s *MetadataService) UpdateAggregatorModelMetadata(epoch int, modelHashCid string, participantIds []string) error {
	epochStr := strconv.Itoa(epoch)
	var participantIdsJSON, err = json.Marshal(participantIds)
	if err != nil {
		return fmt.Errorf("failed to marshal participant ids JSON: %w", err)
	}

	err = s.client.SubmitTransaction(nil, "UpdateAggregatorModelMetadata", epochStr, modelHashCid, string(participantIdsJSON))
	if err != nil {
		return fmt.Errorf("failed to update aggregator model metadata record: %w", err)
	}

	return nil
}

// DeleteAllAggregatorModelMetadata deletes all aggregator model metadata records, returns nil if successful. Only the admin can delete all aggregator model metadata records.
func (s *MetadataService) DeleteAllAggregatorModelMetadata() error {
	err := s.client.SubmitTransaction(nil, "DeleteAllAggregatorModelMetadata")
	if err != nil {
		return fmt.Errorf("failed to delete all aggregator model metadata records: %w", err)
	}

	return nil
}

// GetAllAggregatorModelMetadata queries the aggregators' model metadata records from the ledger for all epochs
func (s *MetadataService) GetAllAggregatorModelMetadata() ([]shared.AggregatorModelMetadata, error) {
	var aggregatorModelMetadataList []shared.AggregatorModelMetadata

	err := s.client.EvaluateTransaction(&aggregatorModelMetadataList, "GetAllAggregatorModelMetadata")
	if err != nil {
		return nil, fmt.Errorf("failed to query all aggregator model metadata records: %w", err)
	}
	return aggregatorModelMetadataList, nil
}

// GetAllAggregatorModelMetadataByAggregator queries the aggregator model metadata records from the ledger for all epochs
func (s *MetadataService) GetAllAggregatorModelMetadataByAggregator(aggregatorId string) ([]shared.AggregatorModelMetadata, error) {
	var aggregatorModelMetadataList []shared.AggregatorModelMetadata

	err := s.client.EvaluateTransaction(&aggregatorModelMetadataList, "GetAllAggregatorModelMetadata", aggregatorId)
	if err != nil {
		return nil, fmt.Errorf("failed to query all aggregator model metadata records: %w", err)
	}
	return aggregatorModelMetadataList, nil
}

// --------------------------------------------
// THIS SECTION DEALS WITH ACCESSING THE LOGS
// --------------------------------------------

// GetAllLogsWithoutCreator returns all logs from the ledger without the creator information.
// This is much faster than GetAllLogs() since it doesn't need to parse the ledger N times for the creators.
// Only admins can use this function.
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
// Only admins can use this function.
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
			log.TxCreator = *creatorInfo
			history[i] = log
		}
	}

	return history, nil
}

// GetAllLogsForUser returns all logs from the ledger tied to the user.
// Only admins can use this function.
// MSPID - the MSP ID of the user
// SerialNumber - the serial number of the user
func (s *MetadataService) GetAllLogsForUser(MSPID string, SerialNumber string) ([]shared.LogEntry, error) {
	var allLogs, err = s.GetAllLogs()
	if err != nil {
		return nil, fmt.Errorf("failed to query all logs: %w", err)
	}

	var history []shared.LogEntry
	for _, log := range allLogs {
		if log.TxCreator.MSPID == MSPID && log.TxCreator.SerialNumber == SerialNumber {
			history = append(history, log)
		}
	}

	return history, nil
}

// Close closes the Fabric client
func (s *MetadataService) Close() error {
	return s.client.Close()
}
