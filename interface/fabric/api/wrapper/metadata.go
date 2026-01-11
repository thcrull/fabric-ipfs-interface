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

// AddParticipant submits a transaction to add a new participant record. The participant record will be bound to the caller's identity,
// thus changes made to the record can only be done by the creator or an admin.
func (s *MetadataService) AddParticipant(participantId int, encapsulatedKey string, homomorphicSharedKeyCypher string, communicationKeyCypher string) error {
	participantIdStr := strconv.Itoa(participantId)

	err := s.client.SubmitTransaction(nil, "AddParticipant", participantIdStr, encapsulatedKey, homomorphicSharedKeyCypher, communicationKeyCypher)
	if err != nil {
		return fmt.Errorf("failed to add participant record for id %d: %w", participantId, err)
	}

	return nil
}

// GetParticipant retrieves a participant record by id. Can be done by anyone.
func (s *MetadataService) GetParticipant(participantId int) (*shared.Participant, error) {
	participantIdStr := strconv.Itoa(participantId)
	var participant shared.Participant

	err := s.client.EvaluateTransaction(&participant, "GetParticipant", participantIdStr)
	if err != nil {
		return nil, fmt.Errorf("failed to query the participant record for id %d: %w", participantId, err)
	}

	return &participant, nil
}

// ParticipantExists returns true if a participant record exists. Can be done by anyone.
func (s *MetadataService) ParticipantExists(participantId int) (bool, error) {
	participantIdStr := strconv.Itoa(participantId)
	var exists bool

	err := s.client.EvaluateTransaction(&exists, "ParticipantExists", participantIdStr)
	if err != nil {
		return false, fmt.Errorf("failed to query if participant record exists for id %d: %w", participantId, err)
	}

	return exists, nil
}

// DeleteParticipant deletes the participant record, returns nil if successful. Can only be done by the participant's creator or an admin.
func (s *MetadataService) DeleteParticipant(participantId int) error {
	participantIdStr := strconv.Itoa(participantId)
	err := s.client.SubmitTransaction(nil, "DeleteParticipant", participantIdStr)
	if err != nil {
		return fmt.Errorf("failed to delete participant record for id %d: %w", participantId, err)
	}

	return nil
}

// UpdateParticipant updates the caller's participant record. Can only be done by the participant's creator or an admin.
func (s *MetadataService) UpdateParticipant(participantId int, encapsulatedKey string, homomorphicSharedKeyCypher string, communicationKeyCypher string) error {
	participantIdStr := strconv.Itoa(participantId)

	err := s.client.SubmitTransaction(nil, "UpdateParticipant", participantIdStr, encapsulatedKey, homomorphicSharedKeyCypher, communicationKeyCypher)
	if err != nil {
		return fmt.Errorf("failed to update participant record for id %d: %w", participantId, err)
	}

	return nil
}

// DeleteAllParticipants deletes all participant records, returns nil if successful. Only the admin can delete all participant records.
func (s *MetadataService) DeleteAllParticipants() error {
	err := s.client.SubmitTransaction(nil, "DeleteAllParticipants")
	if err != nil {
		return fmt.Errorf("failed to delete all participants records: %w", err)
	}

	return nil
}

// GetAllParticipants queries all participant records from the ledger. Can be done by anyone.
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

// AddAggregator submits a transaction to add a new aggregator record. The aggregator record will be bound to the caller's identity,
// thus changes made to the record can only be done by the creator or an admin.
func (s *MetadataService) AddAggregator(aggregatorId int, communicationKeysCyphers map[string]string) error {
	aggregatorIdStr := strconv.Itoa(aggregatorId)

	var communicationKeysCyphersJSON, err = json.Marshal(communicationKeysCyphers)
	if err != nil {
		return fmt.Errorf("failed to marshal communication keys cyphers JSON: %w", err)
	}

	err = s.client.SubmitTransaction(nil, "AddAggregator", aggregatorIdStr, string(communicationKeysCyphersJSON))
	if err != nil {
		return fmt.Errorf("failed to add aggregator record for id %d: %w", aggregatorId, err)
	}

	return nil
}

// GetAggregator retrieves an aggregator record by id.Can be done by anyone.
func (s *MetadataService) GetAggregator(aggregatorId int) (*shared.Aggregator, error) {
	aggregatorIdStr := strconv.Itoa(aggregatorId)
	var aggregator shared.Aggregator

	err := s.client.EvaluateTransaction(&aggregator, "GetAggregator", aggregatorIdStr)
	if err != nil {
		return nil, fmt.Errorf("failed to query the aggregator record for id %d: %w", aggregatorId, err)
	}

	return &aggregator, nil
}

// AggregatorExists returns true if an aggregator record exists. Can be done by anyone.
func (s *MetadataService) AggregatorExists(aggregatorId int) (bool, error) {
	aggregatorIdStr := strconv.Itoa(aggregatorId)
	var exists bool

	err := s.client.EvaluateTransaction(&exists, "AggregatorExists", aggregatorIdStr)
	if err != nil {
		return false, fmt.Errorf("failed to query if aggregator record exists for id %d: %w", aggregatorId, err)
	}

	return exists, nil
}

// DeleteAggregator deletes the aggregator record, returns nil if successful. Can be done only by the aggregator's creator or an admin.
func (s *MetadataService) DeleteAggregator(aggregatorId int) error {
	aggregatorIdStr := strconv.Itoa(aggregatorId)

	err := s.client.SubmitTransaction(nil, "DeleteAggregator", aggregatorIdStr)
	if err != nil {
		return fmt.Errorf("failed to delete aggregator record for id %d: %w", aggregatorId, err)
	}

	return nil
}

// UpdateAggregator updates the caller's aggregator record. Can only be done by the aggregator's creator or an admin.
func (s *MetadataService) UpdateAggregator(aggregatorId int, communicationKeysCyphers map[string]string) error {
	aggregatorIdStr := strconv.Itoa(aggregatorId)
	var communicationKeysCyphersJSON, err = json.Marshal(communicationKeysCyphers)
	if err != nil {
		return fmt.Errorf("failed to marshal communication keys cyphers JSON: %w", err)
	}

	err = s.client.SubmitTransaction(nil, "UpdateAggregator", aggregatorIdStr, string(communicationKeysCyphersJSON))
	if err != nil {
		return fmt.Errorf("failed to update aggregator record: %w", err)
	}

	return nil
}

// DeleteAllAggregators deletes all aggregator records, returns nil if successful. Only the admin can delete all aggregator records.
func (s *MetadataService) DeleteAllAggregators() error {
	err := s.client.SubmitTransaction(nil, "DeleteAllAggregators")
	if err != nil {
		return fmt.Errorf("failed to delete all aggregators records: %w", err)
	}

	return nil
}

// GetAllAggregators queries all aggregator records from the ledger. Can be done by anyone.
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

// AddParticipantModelMetadata submits a transaction to add a new metadata record for a participant's model.
// Only the owner of the participant record or an admin can add a new metadata record for the participant's id.
// The metadata record will be bound to the caller's identity, thus changes made to the record can only be done by the creator or an admin.
func (s *MetadataService) AddParticipantModelMetadata(participantId int, epoch int, modelHashCid string, homomorphicHash string) error {
	participantIdStr := strconv.Itoa(participantId)
	epochStr := strconv.Itoa(epoch)

	err := s.client.SubmitTransaction(nil, "AddParticipantModelMetadata", participantIdStr, epochStr, modelHashCid, homomorphicHash)
	if err != nil {
		return fmt.Errorf("failed to add participant model metadata record for participant id %d and epoch %d: %w", participantId, epoch, err)
	}

	return nil
}

// GetParticipantModelMetadata retrieves a participant model metadata record by participantId and epoch. Can be done by anyone.
func (s *MetadataService) GetParticipantModelMetadata(participantId int, epoch int) (*shared.ParticipantModelMetadata, error) {
	participantIdStr := strconv.Itoa(participantId)
	epochStr := strconv.Itoa(epoch)
	var participantModelMetadata shared.ParticipantModelMetadata

	err := s.client.EvaluateTransaction(&participantModelMetadata, "GetParticipantModelMetadata", participantIdStr, epochStr)
	if err != nil {
		return nil, fmt.Errorf("failed to query the participant model metadata record for participant id %d and epoch %d: %w", participantId, epoch, err)
	}

	return &participantModelMetadata, nil
}

// ParticipantModelMetadataExists returns true if a participant model metadata record exists. Can be done by anyone.
func (s *MetadataService) ParticipantModelMetadataExists(participantId int, epoch int) (bool, error) {
	participantIdStr := strconv.Itoa(participantId)
	epochStr := strconv.Itoa(epoch)
	var exists bool

	err := s.client.EvaluateTransaction(&exists, "ParticipantModelMetadataExists", participantIdStr, epochStr)
	if err != nil {
		return false, fmt.Errorf("failed to query if participant model metadata record exists for participant id %d and epoch %d: %w", participantId, epoch, err)
	}

	return exists, nil
}

// DeleteParticipantModelMetadata deletes a participant model metadata record, returns nil if successful. Can be done only by the record's owner or an admin.
func (s *MetadataService) DeleteParticipantModelMetadata(participantId int, epoch int) error {
	participantIdStr := strconv.Itoa(participantId)
	epochStr := strconv.Itoa(epoch)

	err := s.client.SubmitTransaction(nil, "DeleteParticipantModelMetadata", participantIdStr, epochStr)
	if err != nil {
		return fmt.Errorf("failed to delete participant model metadata record for participant id %d and epoch %d: %w", participantId, epoch, err)
	}

	return nil
}

// UpdateParticipantModelMetadata updates an existing participant model metadata record. Can be done only by the record's owner or an admin.
func (s *MetadataService) UpdateParticipantModelMetadata(participantId int, epoch int, modelHashCid string, homomorphicHash string) error {
	participantIdStr := strconv.Itoa(participantId)
	epochStr := strconv.Itoa(epoch)

	err := s.client.SubmitTransaction(nil, "UpdateParticipantModelMetadata", participantIdStr, epochStr, modelHashCid, homomorphicHash)
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

// GetAllParticipantModelMetadata queries all participant model metadata records from the ledger made by any participant on any epoch. Can be done by anyone.
func (s *MetadataService) GetAllParticipantModelMetadata() ([]shared.ParticipantModelMetadata, error) {
	var participantModelMetadataList []shared.ParticipantModelMetadata

	err := s.client.EvaluateTransaction(&participantModelMetadataList, "GetAllParticipantModelMetadata")
	if err != nil {
		return nil, fmt.Errorf("failed to query all participant model metadata records: %w", err)
	}
	return participantModelMetadataList, nil
}

// GetAllParticipantModelMetadataByParticipant queries all participant model metadata records from the ledger made by the participant for any epoch. Can be done by anyone.
func (s *MetadataService) GetAllParticipantModelMetadataByParticipant(participantId int) ([]shared.ParticipantModelMetadata, error) {
	participantIdStr := strconv.Itoa(participantId)
	var participantModelMetadataList []shared.ParticipantModelMetadata

	err := s.client.EvaluateTransaction(&participantModelMetadataList, "GetAllParticipantModelMetadataByParticipant", participantIdStr)
	if err != nil {
		return nil, fmt.Errorf("failed to query the participant model metadata records by participant id %d: %w", participantId, err)
	}

	return participantModelMetadataList, nil
}

// GetAllParticipantModelMetadataByEpoch queries all participant model metadata records from the ledger made by any participants for the epoch. Can be done by anyone.
func (s *MetadataService) GetAllParticipantModelMetadataByEpoch(epoch int) ([]shared.ParticipantModelMetadata, error) {
	epochStr := strconv.Itoa(epoch)
	var participantModelMetadataList []shared.ParticipantModelMetadata

	err := s.client.EvaluateTransaction(&participantModelMetadataList, "GetAllParticipantModelMetadataByEpoch", epochStr)
	if err != nil {
		return nil, fmt.Errorf("failed to query the participant model metadata records by epoch %d: %w", epoch, err)
	}

	return participantModelMetadataList, nil
}

// ---------------------------------------------------------------------------
// THIS SECTION IS FOR THE AGGREGATOR MODEL METADATA RECORDS' FUNCTIONALITIES
// ---------------------------------------------------------------------------

// AddAggregatorModelMetadata submits a transaction to add a new aggregator model metadata record.
// Only the owner of the aggregator record or an admin can add a new metadata record for the aggregator's id.
// The metadata record will be bound to the caller's identity, thus changes made to the record can only be done by the creator or an admin.
func (s *MetadataService) AddAggregatorModelMetadata(aggregatorId int, epoch int, modelHashCid string, participantIds []int) error {
	aggregatorIdStr := strconv.Itoa(aggregatorId)
	epochStr := strconv.Itoa(epoch)
	var participantIdsJSON, err = json.Marshal(participantIds)
	if err != nil {
		return fmt.Errorf("failed to marshal participant ids JSON: %w", err)
	}

	err = s.client.SubmitTransaction(nil, "AddAggregatorModelMetadata", aggregatorIdStr, epochStr, modelHashCid, string(participantIdsJSON))
	if err != nil {
		return fmt.Errorf("failed to add aggregator model metadata record for aggregator id %d and epoch %d: %w", aggregatorId, epoch, err)
	}

	return nil
}

// GetAggregatorModelMetadata retrieves an aggregator model metadata record by aggregatorId and epoch. Can be done by anyone.
func (s *MetadataService) GetAggregatorModelMetadata(aggregatorId int, epoch int) (*shared.AggregatorModelMetadata, error) {
	aggregatorIdStr := strconv.Itoa(aggregatorId)
	epochStr := strconv.Itoa(epoch)
	var aggregatorModelMetadata shared.AggregatorModelMetadata

	err := s.client.EvaluateTransaction(&aggregatorModelMetadata, "GetAggregatorModelMetadata", aggregatorIdStr, epochStr)
	if err != nil {
		return nil, fmt.Errorf("failed to query the aggregator model metadata record for aggregator id %d and epoch %d: %w", aggregatorId, epoch, err)
	}

	return &aggregatorModelMetadata, nil
}

// AggregatorModelMetadataExists returns true if an aggregator model metadata record exists. Can be done by anyone.
func (s *MetadataService) AggregatorModelMetadataExists(aggregatorId int, epoch int) (bool, error) {
	aggregatorIdStr := strconv.Itoa(aggregatorId)
	epochStr := strconv.Itoa(epoch)
	var exists bool

	err := s.client.EvaluateTransaction(&exists, "AggregatorModelMetadataExists", aggregatorIdStr, epochStr)
	if err != nil {
		return false, fmt.Errorf("failed to query if aggregator model metadata record exists for aggregator id %d and epoch %d: %w", aggregatorId, epoch, err)
	}

	return exists, nil
}

// DeleteAggregatorModelMetadata deletes an aggregator model metadata record, returns nil if successful. Can be done only by the record's owner or an admin.
func (s *MetadataService) DeleteAggregatorModelMetadata(aggregatorId int, epoch int) error {
	aggregatorIdStr := strconv.Itoa(aggregatorId)
	epochStr := strconv.Itoa(epoch)

	err := s.client.SubmitTransaction(nil, "DeleteAggregatorModelMetadata", aggregatorIdStr, epochStr)
	if err != nil {
		return fmt.Errorf("failed to delete aggregator model metadata record for aggregator id %d and epoch %d: %w", aggregatorId, epoch, err)
	}

	return nil
}

// UpdateAggregatorModelMetadata updates an existing aggregator model metadata record. Can be done only by the record's owner or an admin.
func (s *MetadataService) UpdateAggregatorModelMetadata(aggregatorId int, epoch int, modelHashCid string, participantIds []int) error {
	aggregatorIdStr := strconv.Itoa(aggregatorId)
	epochStr := strconv.Itoa(epoch)
	var participantIdsJSON, err = json.Marshal(participantIds)
	if err != nil {
		return fmt.Errorf("failed to marshal participant ids JSON: %w", err)
	}

	err = s.client.SubmitTransaction(nil, "UpdateAggregatorModelMetadata", aggregatorIdStr, epochStr, modelHashCid, string(participantIdsJSON))
	if err != nil {
		return fmt.Errorf("failed to update aggregator model metadata record for aggregator id %d and epoch %d: %w", aggregatorId, epoch, err)
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

// GetAllAggregatorModelMetadata queries the aggregators' model metadata records from the ledger for all epochs. Can be done by anyone.
func (s *MetadataService) GetAllAggregatorModelMetadata() ([]shared.AggregatorModelMetadata, error) {
	var aggregatorModelMetadataList []shared.AggregatorModelMetadata

	err := s.client.EvaluateTransaction(&aggregatorModelMetadataList, "GetAllAggregatorModelMetadata")
	if err != nil {
		return nil, fmt.Errorf("failed to query all aggregator model metadata records: %w", err)
	}
	return aggregatorModelMetadataList, nil
}

// GetAllAggregatorModelMetadataByAggregator queries the aggregator model metadata records from the ledger for all epochs. Can be done by anyone.
func (s *MetadataService) GetAllAggregatorModelMetadataByAggregator(aggregatorId int) ([]shared.AggregatorModelMetadata, error) {
	aggregatorIdStr := strconv.Itoa(aggregatorId)
	var aggregatorModelMetadataList []shared.AggregatorModelMetadata

	err := s.client.EvaluateTransaction(&aggregatorModelMetadataList, "GetAllAggregatorModelMetadata", aggregatorIdStr)
	if err != nil {
		return nil, fmt.Errorf("failed to query all aggregator model metadata records by aggregator id %d: %w", aggregatorId, err)
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
		var found, creatorInfo, err = s.client.GetTransactionCreator(context.Background(), log.TxId, 0)
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
