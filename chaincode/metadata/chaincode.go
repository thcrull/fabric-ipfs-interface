package metadata

import (
	"encoding/json"
	"fmt"

	"github.com/thcrull/fabric-ipfs-interface/shared"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// MetadataSmartContract provides functions for managing metadata records in the world state of a Fabric network.
type MetadataSmartContract struct {
	contractapi.Contract
}

// ----------------------------------------------------------
// THIS SECTION DEALS WITH PARTICIPANT MODEL UPDATE METADATA
// ----------------------------------------------------------

// AddParticipantModelMetadata issues a new participant's model update metadata record to the world state with the given details.
func (s *MetadataSmartContract) AddParticipantModelMetadata(
	ctx contractapi.TransactionContextInterface,
	epoch int,
	participantId string,
	modelHashCid string,
	homomorphicHash string,
) error {
	compositeKey, err := ctx.GetStub().CreateCompositeKey("participant_model_metadata", []string{participantId, fmt.Sprintf("%d", epoch)})
	if err != nil {
		return fmt.Errorf("failed creating composite key: %v", err)
	}

	exists, err := s.ParticipantModelMetadataExists(ctx, epoch, participantId)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the participant model metadata record for epoch %d from participant %s already exists", epoch, participantId)
	}

	participantModelMetadata := shared.ParticipantModelMetadata{
		Epoch:           epoch,
		ParticipantId:   participantId,
		ModelHashCid:    modelHashCid,
		HomomorphicHash: homomorphicHash,
	}
	metadataJSON, err := json.Marshal(participantModelMetadata)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(compositeKey, metadataJSON)
}

// GetParticipantModelMetadata returns the participant's model update metadata record stored in the world state for the given epoch and participant id.
func (s *MetadataSmartContract) GetParticipantModelMetadata(ctx contractapi.TransactionContextInterface, epoch int, participantId string) (*shared.ParticipantModelMetadata, error) {
	compositeKey, err := ctx.GetStub().CreateCompositeKey("participant_model_metadata", []string{participantId, fmt.Sprintf("%d", epoch)})
	if err != nil {
		return nil, fmt.Errorf("failed creating composite key: %v", err)
	}

	metadataJSON, err := ctx.GetStub().GetState(compositeKey)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if metadataJSON == nil {
		return nil, fmt.Errorf("the participant model metadata record for epoch %d from participant %s does not exist", epoch, participantId)
	}

	var participantModelMetadata shared.ParticipantModelMetadata
	err = json.Unmarshal(metadataJSON, &participantModelMetadata)
	if err != nil {
		return nil, err
	}

	return &participantModelMetadata, nil
}

// ParticipantModelMetadataExists returns true when a participant model metadata record for the given epoch and participantId exists in the world state
func (s *MetadataSmartContract) ParticipantModelMetadataExists(ctx contractapi.TransactionContextInterface, epoch int, participantId string) (bool, error) {
	compositeKey, err := ctx.GetStub().CreateCompositeKey("participant_model_metadata", []string{participantId, fmt.Sprintf("%d", epoch)})
	if err != nil {
		return false, fmt.Errorf("failed creating composite key: %v", err)
	}

	metadataJSON, err := ctx.GetStub().GetState(compositeKey)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return metadataJSON != nil, nil
}

// DeleteParticipantModelMetadata deletes a given participant model metadata record from the world state.
func (s *MetadataSmartContract) DeleteParticipantModelMetadata(ctx contractapi.TransactionContextInterface, epoch int, participantId string) error {
	exists, err := s.ParticipantModelMetadataExists(ctx, epoch, participantId)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the participant model metadata record for epoch %d from participant %s does not exist", epoch, participantId)
	}

	compositeKey, err := ctx.GetStub().CreateCompositeKey("participant_model_metadata", []string{participantId, fmt.Sprintf("%d", epoch)})
	if err != nil {
		return fmt.Errorf("failed creating composite key: %v", err)
	}

	return ctx.GetStub().DelState(compositeKey)
}

// UpdateParticipantModelMetadata updates an existing participant model metadata record in the world state with provided parameters.
func (s *MetadataSmartContract) UpdateParticipantModelMetadata(
	ctx contractapi.TransactionContextInterface,
	epoch int,
	participantId string,
	modelHashCid string,
	homomorphicHash string,
) error {
	exists, err := s.ParticipantModelMetadataExists(ctx, epoch, participantId)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the participant model metadata record for epoch %d from participant %s does not exist", epoch, participantId)
	}

	// overwriting original metadata with new metadata
	participantModelMetadata := shared.ParticipantModelMetadata{
		Epoch:           epoch,
		ParticipantId:   participantId,
		ModelHashCid:    modelHashCid,
		HomomorphicHash: homomorphicHash,
	}
	metadataJSON, err := json.Marshal(participantModelMetadata)
	if err != nil {
		return err
	}

	compositeKey, err := ctx.GetStub().CreateCompositeKey("participant_model_metadata", []string{participantId, fmt.Sprintf("%d", epoch)})
	if err != nil {
		return fmt.Errorf("failed creating composite key: %v", err)
	}

	return ctx.GetStub().PutState(compositeKey, metadataJSON)
}

// DeleteAllParticipantModelMetadata deletes all participant model metadata records from the world state.
func (s *MetadataSmartContract) DeleteAllParticipantModelMetadata(ctx contractapi.TransactionContextInterface) error {
	participantModelMetadataBlocks, err := s.GetAllParticipantModelMetadata(ctx)
	if err != nil {
		return fmt.Errorf("error getting all participant model metadata records for deletion: %v", err)
	}

	for _, participantModelMetadataBlock := range participantModelMetadataBlocks {
		err := s.DeleteParticipantModelMetadata(ctx, participantModelMetadataBlock.Epoch, participantModelMetadataBlock.ParticipantId)
		if err != nil {
			return fmt.Errorf("error deleting participant model metadata record: %v", err)
		}
	}

	return nil
}

// GetAllParticipantModelMetadata returns all participant model metadata records found in the world state.
func (s *MetadataSmartContract) GetAllParticipantModelMetadata(ctx contractapi.TransactionContextInterface) ([]*shared.ParticipantModelMetadata, error) {
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("participant_model_metadata", []string{})
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var participantModelMetadataBlocks []*shared.ParticipantModelMetadata
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var participantModelMetadata shared.ParticipantModelMetadata
		if err := json.Unmarshal(queryResponse.Value, &participantModelMetadata); err != nil {
			return nil, err
		}
		participantModelMetadataBlocks = append(participantModelMetadataBlocks, &participantModelMetadata)
	}

	return participantModelMetadataBlocks, nil
}

// GetAllParticipantModelMetadataByParticipant returns all participant model metadata records found in the world state created by the participant.
func (s *MetadataSmartContract) GetAllParticipantModelMetadataByParticipant(ctx contractapi.TransactionContextInterface, participantId string) ([]*shared.ParticipantModelMetadata, error) {
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("participant_model_metadata", []string{participantId})
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var participantModelMetadataBlocks []*shared.ParticipantModelMetadata
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var participantModelMetadata shared.ParticipantModelMetadata
		if err := json.Unmarshal(queryResponse.Value, &participantModelMetadata); err != nil {
			return nil, err
		}
		participantModelMetadataBlocks = append(participantModelMetadataBlocks, &participantModelMetadata)
	}

	return participantModelMetadataBlocks, nil
}

// GetAllParticipantModelMetadataByEpoch returns all participant model metadata records found in the world state for the given epoch.
func (s *MetadataSmartContract) GetAllParticipantModelMetadataByEpoch(ctx contractapi.TransactionContextInterface, epoch int) ([]*shared.ParticipantModelMetadata, error) {
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("participant_model_metadata", []string{})
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var participantModelMetadataBlocks []*shared.ParticipantModelMetadata
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var metadata shared.ParticipantModelMetadata
		if err := json.Unmarshal(queryResponse.Value, &metadata); err != nil {
			return nil, err
		}

		if metadata.Epoch == epoch {
			participantModelMetadataBlocks = append(participantModelMetadataBlocks, &metadata)
		}
	}

	return participantModelMetadataBlocks, nil
}

// ---------------------------------------------------
// THIS SECTION DEALS WITH AGGREGATOR MODEL METADATA
// ---------------------------------------------------

// AddAggregatorModelMetadata issues a new aggregator's model aggregation metadata record to the world state with the given details.
func (s *MetadataSmartContract) AddAggregatorModelMetadata(
	ctx contractapi.TransactionContextInterface,
	epoch int,
	participantIds []string,
	modelHashCid string,
) error {
	compositeKey, err := ctx.GetStub().CreateCompositeKey("aggregator_model_metadata", []string{fmt.Sprintf("%d", epoch)})
	if err != nil {
		return fmt.Errorf("failed creating composite key: %v", err)
	}

	exists, err := s.AggregatorModelMetadataExists(ctx, epoch)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the aggregator model metadata record for epoch %d already exists", epoch)
	}

	aggregatorModelMetadata := shared.AggregatorModelMetadata{
		Epoch:          epoch,
		ParticipantIds: participantIds,
		ModelHashCid:   modelHashCid,
	}
	metadataJSON, err := json.Marshal(aggregatorModelMetadata)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(compositeKey, metadataJSON)
}

// GetAggregatorModelMetadata returns the aggregator's model aggregation metadata record stored in the world state for the given epoch.
func (s *MetadataSmartContract) GetAggregatorModelMetadata(ctx contractapi.TransactionContextInterface, epoch int) (*shared.AggregatorModelMetadata, error) {
	compositeKey, err := ctx.GetStub().CreateCompositeKey("aggregator_model_metadata", []string{fmt.Sprintf("%d", epoch)})
	if err != nil {
		return nil, fmt.Errorf("failed creating composite key: %v", err)
	}

	metadataJSON, err := ctx.GetStub().GetState(compositeKey)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if metadataJSON == nil {
		return nil, fmt.Errorf("the aggregator model metadata record for epoch %d does not exist", epoch)
	}

	var aggregatorModelMetadata shared.AggregatorModelMetadata
	err = json.Unmarshal(metadataJSON, &aggregatorModelMetadata)
	if err != nil {
		return nil, err
	}

	return &aggregatorModelMetadata, nil
}

// AggregatorModelMetadataExists returns true when an aggregator model metadata record for the given epoch exists in the world state.
func (s *MetadataSmartContract) AggregatorModelMetadataExists(ctx contractapi.TransactionContextInterface, epoch int) (bool, error) {
	compositeKey, err := ctx.GetStub().CreateCompositeKey("aggregator_model_metadata", []string{fmt.Sprintf("%d", epoch)})
	if err != nil {
		return false, fmt.Errorf("failed creating composite key: %v", err)
	}

	metadataJSON, err := ctx.GetStub().GetState(compositeKey)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return metadataJSON != nil, nil
}

// DeleteAggregatorModelMetadata deletes a given aggregator model metadata record from the world state.
func (s *MetadataSmartContract) DeleteAggregatorModelMetadata(ctx contractapi.TransactionContextInterface, epoch int) error {
	exists, err := s.AggregatorModelMetadataExists(ctx, epoch)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the aggregator model metadata record for epoch %d does not exist", epoch)
	}

	compositeKey, err := ctx.GetStub().CreateCompositeKey("aggregator_model_metadata", []string{fmt.Sprintf("%d", epoch)})
	if err != nil {
		return fmt.Errorf("failed creating composite key: %v", err)
	}

	return ctx.GetStub().DelState(compositeKey)
}

// UpdateAggregatorModelMetadata updates an existing aggregator model metadata record in the world state with provided parameters.
func (s *MetadataSmartContract) UpdateAggregatorModelMetadata(
	ctx contractapi.TransactionContextInterface,
	epoch int,
	participantIds []string,
	modelHashCid string,
) error {
	exists, err := s.AggregatorModelMetadataExists(ctx, epoch)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the aggregator model metadata record for epoch %d does not exist", epoch)
	}

	// overwriting original metadata with new metadata
	aggregatorModelMetadata := shared.AggregatorModelMetadata{
		Epoch:          epoch,
		ParticipantIds: participantIds,
		ModelHashCid:   modelHashCid,
	}
	metadataJSON, err := json.Marshal(aggregatorModelMetadata)
	if err != nil {
		return err
	}

	compositeKey, err := ctx.GetStub().CreateCompositeKey("aggregator_model_metadata", []string{fmt.Sprintf("%d", epoch)})
	if err != nil {
		return fmt.Errorf("failed creating composite key: %v", err)
	}

	return ctx.GetStub().PutState(compositeKey, metadataJSON)
}

// DeleteAllAggregatorModelMetadata deletes all aggregator model metadata records from the world state.
func (s *MetadataSmartContract) DeleteAllAggregatorModelMetadata(ctx contractapi.TransactionContextInterface) error {
	aggregatorModelMetadataBlocks, err := s.GetAllAggregatorModelMetadata(ctx)
	if err != nil {
		return fmt.Errorf("error getting all aggregator model metadata records for deletion: %v", err)
	}

	for _, aggregatorModelMetadataBlock := range aggregatorModelMetadataBlocks {
		err := s.DeleteAggregatorModelMetadata(ctx, aggregatorModelMetadataBlock.Epoch)
		if err != nil {
			return fmt.Errorf("error deleting aggregator model metadata record: %v", err)
		}
	}

	return nil
}

// GetAllAggregatorModelMetadata returns all aggregator model metadata records found in the world state.
func (s *MetadataSmartContract) GetAllAggregatorModelMetadata(ctx contractapi.TransactionContextInterface) ([]*shared.AggregatorModelMetadata, error) {
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("aggregator_model_metadata", []string{})
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var aggregatorModelMetadataBlocks []*shared.AggregatorModelMetadata
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var aggregatorModelMetadata shared.AggregatorModelMetadata
		if err := json.Unmarshal(queryResponse.Value, &aggregatorModelMetadata); err != nil {
			return nil, err
		}
		aggregatorModelMetadataBlocks = append(aggregatorModelMetadataBlocks, &aggregatorModelMetadata)
	}

	return aggregatorModelMetadataBlocks, nil
}

// ---------------------------------------------------
// THIS SECTION DEALS WITH PARTICIPANT INFORMATION
// ---------------------------------------------------

// AddParticipant issues a participant record to the world state with the given details.
func (s *MetadataSmartContract) AddParticipant(
	ctx contractapi.TransactionContextInterface,
	participantId string,
	encapsulatedKey string,
) error {
	compositeKey, err := ctx.GetStub().CreateCompositeKey("participant", []string{participantId})
	if err != nil {
		return fmt.Errorf("failed creating composite key: %v", err)
	}

	exists, err := s.ParticipantExists(ctx, participantId)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the participant record for id %s already exists", participantId)
	}

	participant := shared.Participant{
		ParticipantId:   participantId,
		EncapsulatedKey: encapsulatedKey,
	}
	participantJSON, err := json.Marshal(participant)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(compositeKey, participantJSON)
}

// GetParticipant returns the participant record stored in the world state for the given id.
func (s *MetadataSmartContract) GetParticipant(ctx contractapi.TransactionContextInterface, participantId string) (*shared.Participant, error) {
	compositeKey, err := ctx.GetStub().CreateCompositeKey("participant", []string{participantId})
	if err != nil {
		return nil, fmt.Errorf("failed creating composite key: %v", err)
	}

	participantJSON, err := ctx.GetStub().GetState(compositeKey)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if participantJSON == nil {
		return nil, fmt.Errorf("the participant record for id %s does not exist", participantId)
	}

	var participant shared.Participant
	err = json.Unmarshal(participantJSON, &participant)
	if err != nil {
		return nil, err
	}

	return &participant, nil
}

// ParticipantExists returns true when a participant record for the given id exists in the world state.
func (s *MetadataSmartContract) ParticipantExists(ctx contractapi.TransactionContextInterface, participantId string) (bool, error) {
	compositeKey, err := ctx.GetStub().CreateCompositeKey("participant", []string{participantId})
	if err != nil {
		return false, fmt.Errorf("failed creating composite key: %v", err)
	}

	participantJSON, err := ctx.GetStub().GetState(compositeKey)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return participantJSON != nil, nil
}

// DeleteParticipant deletes a given participant record from the world state.
func (s *MetadataSmartContract) DeleteParticipant(ctx contractapi.TransactionContextInterface, participantId string) error {
	exists, err := s.ParticipantExists(ctx, participantId)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the participant record for id %s does not exist", participantId)
	}

	compositeKey, err := ctx.GetStub().CreateCompositeKey("participant", []string{participantId})
	if err != nil {
		return fmt.Errorf("failed creating composite key: %v", err)
	}

	return ctx.GetStub().DelState(compositeKey)
}

// UpdateParticipant updates a participant record in the world state with provided parameters.
func (s *MetadataSmartContract) UpdateParticipant(
	ctx contractapi.TransactionContextInterface,
	participantId string,
	encapsulatedKey string,
) error {
	exists, err := s.ParticipantExists(ctx, participantId)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the participant record for id %s does not exist", participantId)
	}

	// overwriting original participant with new participant
	participant := shared.Participant{
		ParticipantId:   participantId,
		EncapsulatedKey: encapsulatedKey,
	}
	participantJSON, err := json.Marshal(participant)
	if err != nil {
		return err
	}

	compositeKey, err := ctx.GetStub().CreateCompositeKey("participant", []string{participantId})
	if err != nil {
		return fmt.Errorf("failed creating composite key: %v", err)
	}

	return ctx.GetStub().PutState(compositeKey, participantJSON)
}

// DeleteAllParticipants deletes all participant records from the world state.
func (s *MetadataSmartContract) DeleteAllParticipants(ctx contractapi.TransactionContextInterface) error {
	participants, err := s.GetAllParticipants(ctx)
	if err != nil {
		return fmt.Errorf("error getting all participant records for deletion: %v", err)
	}

	for _, participant := range participants {
		err := s.DeleteParticipant(ctx, participant.ParticipantId)
		if err != nil {
			return fmt.Errorf("error deleting participant record: %v", err)
		}
	}

	return nil
}

// GetAllParticipants returns all participant records found in the world state.
func (s *MetadataSmartContract) GetAllParticipants(ctx contractapi.TransactionContextInterface) ([]*shared.Participant, error) {
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("participant", []string{})
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var participants []*shared.Participant
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var participant shared.Participant
		if err := json.Unmarshal(queryResponse.Value, &participant); err != nil {
			return nil, err
		}
		participants = append(participants, &participant)
	}

	return participants, nil
}
