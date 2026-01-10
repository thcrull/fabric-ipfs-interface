package chaincode

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

// ---------------------------------------------------
// THIS SECTION DEALS WITH PARTICIPANT INFORMATION
// ---------------------------------------------------

// AddParticipant issues a participant record to the world state with the given details. The record is bound to the client's identity.
func (s *MetadataSmartContract) AddParticipant(
	ctx contractapi.TransactionContextInterface,
	participantId int,
	encapsulatedKey string,
	homomorphicSharedKeyCypher string,
	communicationKeyCypher string,
) error {
	exists, err := s.ParticipantExists(ctx, participantId)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the participant record for id %d already exists", participantId)
	}

	MSPID, serialNumber, err := getCreatorInfo(ctx)
	if err != nil {
		return fmt.Errorf("failed getting creator info: %v", err)
	}

	participant := shared.Participant{
		ParticipantId:              participantId,
		EncapsulatedKey:            encapsulatedKey,
		HomomorphicSharedKeyCypher: homomorphicSharedKeyCypher,
		CommunicationKeyCypher:     communicationKeyCypher,
		MSPID:                      MSPID,
		SerialNumber:               serialNumber,
	}
	participantJSON, err := json.Marshal(participant)
	if err != nil {
		return err
	}

	compositeKey, err := ctx.GetStub().CreateCompositeKey("participant", []string{fmt.Sprintf("%d", participantId)})
	if err != nil {
		return fmt.Errorf("failed creating composite key: %v", err)
	}

	return ctx.GetStub().PutState(compositeKey, participantJSON)
}

// GetParticipant returns the participant record stored in the world state for the given id.
func (s *MetadataSmartContract) GetParticipant(ctx contractapi.TransactionContextInterface, participantId int) (*shared.Participant, error) {
	compositeKey, err := ctx.GetStub().CreateCompositeKey("participant", []string{fmt.Sprintf("%d", participantId)})
	if err != nil {
		return nil, fmt.Errorf("failed creating composite key: %v", err)
	}

	participantJSON, err := ctx.GetStub().GetState(compositeKey)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if participantJSON == nil {
		return nil, fmt.Errorf("the participant record for id %d does not exist", participantId)
	}

	var participant shared.Participant
	err = json.Unmarshal(participantJSON, &participant)
	if err != nil {
		return nil, err
	}

	return &participant, nil
}

// ParticipantExists returns true when a participant record for the given id exists in the world state.
func (s *MetadataSmartContract) ParticipantExists(ctx contractapi.TransactionContextInterface, participantId int) (bool, error) {
	compositeKey, err := ctx.GetStub().CreateCompositeKey("participant", []string{fmt.Sprintf("%d", participantId)})
	if err != nil {
		return false, fmt.Errorf("failed creating composite key: %v", err)
	}

	participantJSON, err := ctx.GetStub().GetState(compositeKey)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return participantJSON != nil, nil
}

// DeleteParticipant deletes a given participant record from the world state. Can only be done by the owner of the participant or an admin.
func (s *MetadataSmartContract) DeleteParticipant(ctx contractapi.TransactionContextInterface, participantId int) error {
	errAdminCheck := adminCheck(ctx)
	errParticipantCheck := s.ownerCheckParticipant(ctx, participantId)
	if errAdminCheck != nil && errParticipantCheck != nil {
		return fmt.Errorf("permission denied: client is not an admin or owner of participant %d", participantId)
	}

	compositeKey, err := ctx.GetStub().CreateCompositeKey("participant", []string{fmt.Sprintf("%d", participantId)})
	if err != nil {
		return fmt.Errorf("failed creating composite key: %v", err)
	}

	return ctx.GetStub().DelState(compositeKey)
}

// UpdateParticipant updates a participant record in the world state with provided parameters. Can only be done by the owner of the participant or an admin.
func (s *MetadataSmartContract) UpdateParticipant(
	ctx contractapi.TransactionContextInterface,
	participantId int,
	encapsulatedKey string,
	homomorphicSharedKeyCypher string,
	communicationKeyCypher string,
) error {
	errAdminCheck := adminCheck(ctx)
	errParticipantCheck := s.ownerCheckParticipant(ctx, participantId)
	if errAdminCheck != nil && errParticipantCheck != nil {
		return fmt.Errorf("permission denied: client is not an admin or owner of participant %d", participantId)
	}

	participant, err := s.GetParticipant(ctx, participantId)
	if err != nil {
		return err
	}

	// overwriting original participant with new participant
	participant.EncapsulatedKey = encapsulatedKey
	participant.HomomorphicSharedKeyCypher = homomorphicSharedKeyCypher
	participant.CommunicationKeyCypher = communicationKeyCypher

	participantJSON, err := json.Marshal(participant)
	if err != nil {
		return err
	}

	compositeKey, err := ctx.GetStub().CreateCompositeKey("participant", []string{fmt.Sprintf("%d", participantId)})
	if err != nil {
		return fmt.Errorf("failed creating composite key: %v", err)
	}

	return ctx.GetStub().PutState(compositeKey, participantJSON)
}

// DeleteAllParticipants deletes all participant records from the world state. Can only be done by an admin.
func (s *MetadataSmartContract) DeleteAllParticipants(ctx contractapi.TransactionContextInterface) error {
	err := adminCheck(ctx)
	if err != nil {
		return fmt.Errorf("permission denied: %v", err)
	}

	participants, err := s.GetAllParticipants(ctx)
	if err != nil {
		return fmt.Errorf("error getting all participant records for deletion: %v", err)
	}

	for _, participant := range participants {
		compositeKey, err := ctx.GetStub().CreateCompositeKey("participant", []string{fmt.Sprintf("%d", participant.ParticipantId)})
		if err != nil {
			return fmt.Errorf("failed creating composite key: %v", err)
		}

		err = ctx.GetStub().DelState(compositeKey)
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

// ------------------------------------------------
// THIS SECTION DEALS WITH AGGREGATOR INFORMATION
// ------------------------------------------------

// AddAggregator issues a new aggregator record to the world state with the given details. The record is bound to the client's identity.
func (s *MetadataSmartContract) AddAggregator(ctx contractapi.TransactionContextInterface, aggregatorId int, communicationKeysCyphersJSON string) error {
	exists, err := s.AggregatorExists(ctx, aggregatorId)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the aggregator record for id %d already exists", aggregatorId)
	}

	var communicationKeysCyphers map[string]string
	err = json.Unmarshal([]byte(communicationKeysCyphersJSON), &communicationKeysCyphers)
	if err != nil {
		return fmt.Errorf("failed unmarshaling communication keys cyphers: %v", err)
	}

	MSPID, serialNumber, err := getCreatorInfo(ctx)
	if err != nil {
		return fmt.Errorf("failed getting creator info: %v", err)
	}

	aggregator := shared.Aggregator{
		AggregatorId:             aggregatorId,
		CommunicationKeysCyphers: communicationKeysCyphers,
		MSPID:                    MSPID,
		SerialNumber:             serialNumber,
	}
	aggregatorJSON, err := json.Marshal(aggregator)
	if err != nil {
		return err
	}

	compositeKey, err := ctx.GetStub().CreateCompositeKey("aggregator", []string{fmt.Sprintf("%d", aggregatorId)})
	if err != nil {
		return fmt.Errorf("failed creating composite key: %v", err)
	}

	return ctx.GetStub().PutState(compositeKey, aggregatorJSON)
}

// GetAggregator returns the aggregator record stored in the world state for the given aggregatorId.
func (s *MetadataSmartContract) GetAggregator(ctx contractapi.TransactionContextInterface, aggregatorId int) (*shared.Aggregator, error) {
	compositeKey, err := ctx.GetStub().CreateCompositeKey("aggregator", []string{fmt.Sprintf("%d", aggregatorId)})
	if err != nil {
		return nil, fmt.Errorf("failed creating composite key: %v", err)
	}

	aggregatorJSON, err := ctx.GetStub().GetState(compositeKey)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if aggregatorJSON == nil {
		return nil, fmt.Errorf("the aggregator record for id %d does not exist", aggregatorId)
	}

	var aggregator shared.Aggregator
	err = json.Unmarshal(aggregatorJSON, &aggregator)
	if err != nil {
		return nil, err
	}

	return &aggregator, nil
}

// AggregatorExists returns true when an aggregator record for the given aggregatorId exists in the world state
func (s *MetadataSmartContract) AggregatorExists(ctx contractapi.TransactionContextInterface, aggregatorId int) (bool, error) {
	compositeKey, err := ctx.GetStub().CreateCompositeKey("aggregator", []string{fmt.Sprintf("%d", aggregatorId)})
	if err != nil {
		return false, fmt.Errorf("failed creating composite key: %v", err)
	}

	aggregatorJSON, err := ctx.GetStub().GetState(compositeKey)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return aggregatorJSON != nil, nil
}

// DeleteAggregator deletes a given aggregator record from the world state. Can only be done by the owner of the aggregator or an admin.
func (s *MetadataSmartContract) DeleteAggregator(ctx contractapi.TransactionContextInterface, aggregatorId int) error {
	errAdminCheck := adminCheck(ctx)
	errAggregatorCheck := s.ownerCheckAggregator(ctx, aggregatorId)
	if errAdminCheck != nil && errAggregatorCheck != nil {
		return fmt.Errorf("permission denied: client is not an admin or owner of aggregator %d", aggregatorId)
	}

	compositeKey, err := ctx.GetStub().CreateCompositeKey("aggregator", []string{fmt.Sprintf("%d", aggregatorId)})
	if err != nil {
		return fmt.Errorf("failed creating composite key: %v", err)
	}

	return ctx.GetStub().DelState(compositeKey)
}

// UpdateAggregator updates an existing aggregator record in the world state with provided parameters. Can only be done by the owner of the aggregator or an admin.
func (s *MetadataSmartContract) UpdateAggregator(ctx contractapi.TransactionContextInterface, aggregatorId int, communicationKeysCyphersJSON string) error {
	errAdminCheck := adminCheck(ctx)
	errAggregatorCheck := s.ownerCheckAggregator(ctx, aggregatorId)
	if errAdminCheck != nil && errAggregatorCheck != nil {
		return fmt.Errorf("permission denied: client is not an admin or owner of aggregator %d", aggregatorId)
	}

	// overwriting original metadata with new metadata
	aggregator, err := s.GetAggregator(ctx, aggregatorId)
	if err != nil {
		return err
	}

	var communicationKeysCyphers map[string]string
	err = json.Unmarshal([]byte(communicationKeysCyphersJSON), &communicationKeysCyphers)
	if err != nil {
		return fmt.Errorf("failed unmarshaling communication keys cyphers: %v", err)
	}

	aggregator.CommunicationKeysCyphers = communicationKeysCyphers

	aggregatorJSON, err := json.Marshal(aggregator)
	if err != nil {
		return err
	}

	compositeKey, err := ctx.GetStub().CreateCompositeKey("aggregator", []string{fmt.Sprintf("%d", aggregatorId)})
	if err != nil {
		return fmt.Errorf("failed creating composite key: %v", err)
	}

	return ctx.GetStub().PutState(compositeKey, aggregatorJSON)
}

// DeleteAllAggregators deletes all aggregator records from the world state. Can only be done by an admin.
func (s *MetadataSmartContract) DeleteAllAggregators(ctx contractapi.TransactionContextInterface) error {
	err := adminCheck(ctx)
	if err != nil {
		return fmt.Errorf("permission denied: %v", err)
	}

	aggregators, err := s.GetAllAggregators(ctx)
	if err != nil {
		return fmt.Errorf("error getting all aggregator records for deletion: %v", err)
	}

	for _, aggregator := range aggregators {
		compositeKey, err := ctx.GetStub().CreateCompositeKey("aggregator", []string{fmt.Sprintf("%d", aggregator.AggregatorId)})
		if err != nil {
			return fmt.Errorf("failed creating composite key: %v", err)
		}

		err = ctx.GetStub().DelState(compositeKey)
		if err != nil {
			return fmt.Errorf("error deleting aggregator record: %v", err)
		}
	}

	return nil
}

// GetAllAggregators returns all aggregator records found in the world state.
func (s *MetadataSmartContract) GetAllAggregators(ctx contractapi.TransactionContextInterface) ([]*shared.Aggregator, error) {
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("aggregator", []string{})
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var aggregators []*shared.Aggregator
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var aggregator shared.Aggregator
		if err := json.Unmarshal(queryResponse.Value, &aggregator); err != nil {
			return nil, err
		}
		aggregators = append(aggregators, &aggregator)
	}

	return aggregators, nil
}

// ----------------------------------------------------------
// THIS SECTION DEALS WITH PARTICIPANT MODEL UPDATE METADATA
// ----------------------------------------------------------

// AddParticipantModelMetadata issues a new participant's model update metadata record to the world state with the given details. Can only be done by the owner of the participant or an admin.
func (s *MetadataSmartContract) AddParticipantModelMetadata(
	ctx contractapi.TransactionContextInterface,
	participantId int,
	epoch int,
	modelHashCid string,
	homomorphicHash string,
) error {
	modelExists, err := s.ParticipantModelMetadataExists(ctx, participantId, epoch)
	if err != nil {
		return err
	}
	if modelExists {
		return fmt.Errorf("the participant model metadata record from participant %d for epoch %d already exists", participantId, epoch)
	}

	errAdminCheck := adminCheck(ctx)
	errParticipantCheck := s.ownerCheckParticipant(ctx, participantId)
	if errAdminCheck != nil && errParticipantCheck != nil {
		return fmt.Errorf("permission denied: client is not an admin or owner of participant %d", participantId)
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

	compositeKey, err := ctx.GetStub().CreateCompositeKey("participant_model_metadata", []string{fmt.Sprintf("%d", participantId), fmt.Sprintf("%d", epoch)})
	if err != nil {
		return fmt.Errorf("failed creating composite key: %v", err)
	}

	return ctx.GetStub().PutState(compositeKey, metadataJSON)
}

// GetParticipantModelMetadata returns the participant's model update metadata record stored in the world state for the given participantId and epoch.
func (s *MetadataSmartContract) GetParticipantModelMetadata(ctx contractapi.TransactionContextInterface, participantId int, epoch int) (*shared.ParticipantModelMetadata, error) {
	compositeKey, err := ctx.GetStub().CreateCompositeKey("participant_model_metadata", []string{fmt.Sprintf("%d", participantId), fmt.Sprintf("%d", epoch)})
	if err != nil {
		return nil, fmt.Errorf("failed creating composite key: %v", err)
	}

	metadataJSON, err := ctx.GetStub().GetState(compositeKey)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if metadataJSON == nil {
		return nil, fmt.Errorf("the participant model metadata record from participant %d for epoch %d does not exist", participantId, epoch)
	}

	var participantModelMetadata shared.ParticipantModelMetadata
	err = json.Unmarshal(metadataJSON, &participantModelMetadata)
	if err != nil {
		return nil, err
	}

	return &participantModelMetadata, nil
}

// ParticipantModelMetadataExists returns true when a participant model metadata record for the given participantId and epoch exists in the world state.
func (s *MetadataSmartContract) ParticipantModelMetadataExists(ctx contractapi.TransactionContextInterface, participantId int, epoch int) (bool, error) {
	compositeKey, err := ctx.GetStub().CreateCompositeKey("participant_model_metadata", []string{fmt.Sprintf("%d", participantId), fmt.Sprintf("%d", epoch)})
	if err != nil {
		return false, fmt.Errorf("failed creating composite key: %v", err)
	}

	metadataJSON, err := ctx.GetStub().GetState(compositeKey)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return metadataJSON != nil, nil
}

// DeleteParticipantModelMetadata deletes a given participant model metadata record from the world state. Can only be done by the owner of the participant or an admin.
func (s *MetadataSmartContract) DeleteParticipantModelMetadata(ctx contractapi.TransactionContextInterface, participantId int, epoch int) error {
	modelExists, err := s.ParticipantModelMetadataExists(ctx, participantId, epoch)
	if err != nil {
		return err
	}
	if modelExists {
		return fmt.Errorf("the participant model metadata record from participant %d for epoch %d already exists", participantId, epoch)
	}

	errAdminCheck := adminCheck(ctx)
	errParticipantCheck := s.ownerCheckParticipant(ctx, participantId)
	if errAdminCheck != nil && errParticipantCheck != nil {
		return fmt.Errorf("permission denied: client is not an admin or owner of participant %d", participantId)
	}

	compositeKey, err := ctx.GetStub().CreateCompositeKey("participant_model_metadata", []string{fmt.Sprintf("%d", participantId), fmt.Sprintf("%d", epoch)})
	if err != nil {
		return fmt.Errorf("failed creating composite key: %v", err)
	}

	return ctx.GetStub().DelState(compositeKey)
}

// UpdateParticipantModelMetadata updates an existing participant model metadata record in the world state with provided parameters. Can only be done by the owner of the participant or an admin.
func (s *MetadataSmartContract) UpdateParticipantModelMetadata(
	ctx contractapi.TransactionContextInterface,
	participantId int,
	epoch int,
	modelHashCid string,
	homomorphicHash string,
) error {
	modelExists, err := s.ParticipantModelMetadataExists(ctx, participantId, epoch)
	if err != nil {
		return err
	}
	if modelExists {
		return fmt.Errorf("the participant model metadata record from participant %d for epoch %d already exists", participantId, epoch)
	}

	errAdminCheck := adminCheck(ctx)
	errParticipantCheck := s.ownerCheckParticipant(ctx, participantId)
	if errAdminCheck != nil && errParticipantCheck != nil {
		return fmt.Errorf("permission denied: client is not an admin or owner of participant %d", participantId)
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

	compositeKey, err := ctx.GetStub().CreateCompositeKey("participant_model_metadata", []string{fmt.Sprintf("%d", participantId), fmt.Sprintf("%d", epoch)})
	if err != nil {
		return fmt.Errorf("failed creating composite key: %v", err)
	}

	return ctx.GetStub().PutState(compositeKey, metadataJSON)
}

// DeleteAllParticipantModelMetadata deletes all participant model metadata records from the world state. Can only be done by an admin.
func (s *MetadataSmartContract) DeleteAllParticipantModelMetadata(ctx contractapi.TransactionContextInterface) error {
	err := adminCheck(ctx)
	if err != nil {
		return fmt.Errorf("permission denied: %v", err)
	}

	participantModelMetadataBlocks, err := s.GetAllParticipantModelMetadata(ctx)
	if err != nil {
		return fmt.Errorf("error getting all participant model metadata records for deletion: %v", err)
	}

	for _, participantModelMetadataBlock := range participantModelMetadataBlocks {
		compositeKey, err := ctx.GetStub().CreateCompositeKey("participant_model_metadata", []string{fmt.Sprintf("%d", participantModelMetadataBlock.ParticipantId), fmt.Sprintf("%d", participantModelMetadataBlock.Epoch)})
		if err != nil {
			return fmt.Errorf("failed creating composite key: %v", err)
		}

		err = ctx.GetStub().DelState(compositeKey)
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
func (s *MetadataSmartContract) GetAllParticipantModelMetadataByParticipant(ctx contractapi.TransactionContextInterface, participantId int) ([]*shared.ParticipantModelMetadata, error) {
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("participant_model_metadata", []string{fmt.Sprintf("%d", participantId)})
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

// aggregationCheck checks if all participants have submitted their model metadata records for the epoch that the aggregator is responsible for.
func (s *MetadataSmartContract) aggregationCheck(ctx contractapi.TransactionContextInterface, aggregatorModelMetadata *shared.AggregatorModelMetadata) error {
	// TODO: Maybe we should skip this check for epoch 1 since it might be the starting model?
	var missing []int
	for _, participantId := range aggregatorModelMetadata.ParticipantIds {
		key, err := ctx.GetStub().CreateCompositeKey("participant_model_metadata",
			[]string{fmt.Sprintf("%d", participantId), fmt.Sprintf("%d", aggregatorModelMetadata.Epoch)})
		if err != nil {
			return err
		}

		state, err := ctx.GetStub().GetState(key)
		if err != nil {
			return err
		}
		if state == nil {
			missing = append(missing, participantId)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("aggregation denied, missing model metadata from: %v", missing)
	}
	return nil
}

// AddAggregatorModelMetadata issues a new aggregator's model aggregation metadata record to the world state with the given details. Can only be done by the owner of the aggregator or an admin.
func (s *MetadataSmartContract) AddAggregatorModelMetadata(
	ctx contractapi.TransactionContextInterface,
	aggregatorId int,
	epoch int,
	modelHashCid string,
	participantIdsJSON string,
) error {
	modelExists, err := s.AggregatorModelMetadataExists(ctx, epoch, aggregatorId)
	if err != nil {
		return err
	}
	if modelExists {
		return fmt.Errorf("the aggregator model metadata record from aggregator %d for epoch %d already exists", aggregatorId, epoch)
	}

	errAdminCheck := adminCheck(ctx)
	errAggregatorCheck := s.ownerCheckAggregator(ctx, aggregatorId)
	if errAdminCheck != nil && errAggregatorCheck != nil {
		return fmt.Errorf("permission denied: client is not an admin or owner of aggregator %d", aggregatorId)
	}

	var participantIds []int
	err = json.Unmarshal([]byte(participantIdsJSON), &participantIds)
	if err != nil {
		return fmt.Errorf("failed to unmarshal participant Ids JSON: %v", err)
	}

	aggregatorModelMetadata := shared.AggregatorModelMetadata{
		AggregatorId:   aggregatorId,
		Epoch:          epoch,
		ParticipantIds: participantIds,
		ModelHashCid:   modelHashCid,
	}

	err = s.aggregationCheck(ctx, &aggregatorModelMetadata)
	if err != nil {
		return err
	}

	metadataJSON, err := json.Marshal(aggregatorModelMetadata)
	if err != nil {
		return err
	}

	compositeKey, err := ctx.GetStub().CreateCompositeKey("aggregator_model_metadata", []string{fmt.Sprintf("%d", aggregatorId), fmt.Sprintf("%d", epoch)})
	if err != nil {
		return fmt.Errorf("failed creating composite key: %v", err)
	}

	return ctx.GetStub().PutState(compositeKey, metadataJSON)
}

// GetAggregatorModelMetadata returns the aggregator's model aggregation metadata record stored in the world state for the given epoch.
func (s *MetadataSmartContract) GetAggregatorModelMetadata(ctx contractapi.TransactionContextInterface, aggregatorId int, epoch int) (*shared.AggregatorModelMetadata, error) {
	compositeKey, err := ctx.GetStub().CreateCompositeKey("aggregator_model_metadata", []string{fmt.Sprintf("%d", aggregatorId), fmt.Sprintf("%d", epoch)})
	if err != nil {
		return nil, fmt.Errorf("failed creating composite key: %v", err)
	}

	metadataJSON, err := ctx.GetStub().GetState(compositeKey)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if metadataJSON == nil {
		return nil, fmt.Errorf("the aggregator model metadata record from aggregator %d for epoch %d does not exist", aggregatorId, epoch)
	}

	var aggregatorModelMetadata shared.AggregatorModelMetadata
	err = json.Unmarshal(metadataJSON, &aggregatorModelMetadata)
	if err != nil {
		return nil, err
	}

	return &aggregatorModelMetadata, nil
}

// AggregatorModelMetadataExists returns true when an aggregator model metadata record for the given epoch exists in the world state.
func (s *MetadataSmartContract) AggregatorModelMetadataExists(ctx contractapi.TransactionContextInterface, aggregatorId int, epoch int) (bool, error) {
	compositeKey, err := ctx.GetStub().CreateCompositeKey("aggregator_model_metadata", []string{fmt.Sprintf("%d", aggregatorId), fmt.Sprintf("%d", epoch)})
	if err != nil {
		return false, fmt.Errorf("failed creating composite key: %v", err)
	}

	metadataJSON, err := ctx.GetStub().GetState(compositeKey)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return metadataJSON != nil, nil
}

// DeleteAggregatorModelMetadata deletes a given aggregator model metadata record from the world state. Can only be done by the owner of the aggregator or an admin.
func (s *MetadataSmartContract) DeleteAggregatorModelMetadata(ctx contractapi.TransactionContextInterface, aggregatorId int, epoch int) error {
	exists, err := s.AggregatorModelMetadataExists(ctx, epoch, aggregatorId)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the aggregator model metadata record for epoch %d does not exist", epoch)
	}

	errAdminCheck := adminCheck(ctx)
	errAggregatorCheck := s.ownerCheckAggregator(ctx, aggregatorId)
	if errAdminCheck != nil && errAggregatorCheck != nil {
		return fmt.Errorf("permission denied: client is not an admin or owner of aggregator %d", aggregatorId)
	}

	compositeKey, err := ctx.GetStub().CreateCompositeKey("aggregator_model_metadata", []string{fmt.Sprintf("%d", aggregatorId), fmt.Sprintf("%d", epoch)})
	if err != nil {
		return fmt.Errorf("failed creating composite key: %v", err)
	}

	return ctx.GetStub().DelState(compositeKey)
}

// UpdateAggregatorModelMetadata updates an existing aggregator model metadata record in the world state with provided parameters. Can only be done by the owner of the aggregator or an admin.
func (s *MetadataSmartContract) UpdateAggregatorModelMetadata(
	ctx contractapi.TransactionContextInterface,
	aggregatorId int,
	epoch int,
	modelHashCid string,
	participantIdsJSON string,
) error {
	exists, err := s.AggregatorModelMetadataExists(ctx, epoch, aggregatorId)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the aggregator model metadata record from %d for epoch %d does not exist", aggregatorId, epoch)
	}

	errAdminCheck := adminCheck(ctx)
	errAggregatorCheck := s.ownerCheckAggregator(ctx, aggregatorId)
	if errAdminCheck != nil && errAggregatorCheck != nil {
		return fmt.Errorf("permission denied: client is not an admin or owner of aggregator %d", aggregatorId)
	}

	var participantIds []int
	err = json.Unmarshal([]byte(participantIdsJSON), &participantIds)
	if err != nil {
		return fmt.Errorf("failed to unmarshal participant Ids JSON: %v", err)
	}

	// overwriting original metadata with new metadata
	aggregatorModelMetadata := shared.AggregatorModelMetadata{
		Epoch:          epoch,
		AggregatorId:   aggregatorId,
		ParticipantIds: participantIds,
		ModelHashCid:   modelHashCid,
	}

	err = s.aggregationCheck(ctx, &aggregatorModelMetadata)
	if err != nil {
		return err
	}

	metadataJSON, err := json.Marshal(aggregatorModelMetadata)
	if err != nil {
		return err
	}

	compositeKey, err := ctx.GetStub().CreateCompositeKey("aggregator_model_metadata", []string{fmt.Sprintf("%d", aggregatorId), fmt.Sprintf("%d", epoch)})
	if err != nil {
		return fmt.Errorf("failed creating composite key: %v", err)
	}

	return ctx.GetStub().PutState(compositeKey, metadataJSON)
}

// DeleteAllAggregatorModelMetadata deletes all aggregator model metadata records from the world state. Can only be done by an admin.
func (s *MetadataSmartContract) DeleteAllAggregatorModelMetadata(ctx contractapi.TransactionContextInterface) error {
	err := adminCheck(ctx)
	if err != nil {
		return fmt.Errorf("permission denied: %v", err)
	}

	aggregatorModelMetadataBlocks, err := s.GetAllAggregatorModelMetadata(ctx)
	if err != nil {
		return fmt.Errorf("error getting all aggregator model metadata records for deletion: %v", err)
	}

	for _, aggregatorModelMetadataBlock := range aggregatorModelMetadataBlocks {
		compositeKey, err := ctx.GetStub().CreateCompositeKey("aggregator_model_metadata", []string{fmt.Sprintf("%d", aggregatorModelMetadataBlock.AggregatorId), fmt.Sprintf("%d", aggregatorModelMetadataBlock.Epoch)})
		if err != nil {
			return fmt.Errorf("failed creating composite key: %v", err)
		}

		err = ctx.GetStub().DelState(compositeKey)
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

// GetAllAggregatorModelMetadataByAggregator returns all aggregator model metadata records found in the world state.
func (s *MetadataSmartContract) GetAllAggregatorModelMetadataByAggregator(ctx contractapi.TransactionContextInterface, aggregatorId int) ([]*shared.AggregatorModelMetadata, error) {
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

		if aggregatorModelMetadata.AggregatorId == aggregatorId {
			aggregatorModelMetadataBlocks = append(aggregatorModelMetadataBlocks, &aggregatorModelMetadata)
		}
	}

	return aggregatorModelMetadataBlocks, nil
}

// --------------------------------------------
// THIS SECTION DEALS WITH ACCESSING THE LOGS
// --------------------------------------------

// GetAllLogs returns the history for all objects in the world state. Can only be done by the admin.
func (s *MetadataSmartContract) GetAllLogs(ctx contractapi.TransactionContextInterface) ([]shared.LogEntry, error) {
	err := adminCheck(ctx)
	if err != nil {
		return nil, fmt.Errorf("permission denied: %v", err)
	}

	keys, err := s.getAllKeys(ctx)
	if err != nil {
		return nil, err
	}

	var history []shared.LogEntry
	for _, key := range keys {
		resultsIterator, err := ctx.GetStub().GetHistoryForKey(key)
		if err != nil {
			resultsIterator.Close()
			return nil, fmt.Errorf("failed to get log history: %v", err)
		}

		for resultsIterator.HasNext() {
			modification, err := resultsIterator.Next()
			if err != nil {
				resultsIterator.Close()
				return nil, err
			}

			var record interface{}
			if modification.Value != nil {
				err := json.Unmarshal(modification.Value, &record)
				if err != nil {
					// fallback to raw value
					record = string(modification.Value)
				}
			}

			entry := shared.LogEntry{
				TxId:      modification.TxId,
				Timestamp: modification.Timestamp.String(),
				IsDelete:  modification.IsDelete,
				Changes:   record,
			}
			history = append(history, entry)
		}
		resultsIterator.Close()
	}

	return history, nil
}

// getAllKeys returns all keys of the objects in the world state. Can only be done by the admin.
func (s *MetadataSmartContract) getAllKeys(ctx contractapi.TransactionContextInterface) ([]string, error) {
	err := adminCheck(ctx)
	if err != nil {
		return nil, fmt.Errorf("permission denied: %v", err)
	}

	iterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer iterator.Close()

	var keys []string
	for iterator.HasNext() {
		kv, err := iterator.Next()
		if err != nil {
			return nil, err
		}

		keys = append(keys, kv.Key)
	}

	return keys, nil
}

// -------------------------------------------
// THIS SECTION DEALS WITH UTILITY FUNCTIONS
// -------------------------------------------

// getCreatorInfo returns the MSPID and serial number of the transaction creator.
func getCreatorInfo(ctx contractapi.TransactionContextInterface) (string, string, error) {
	MSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", "", fmt.Errorf("failed to get MSPID: %v", err)
	}

	certificate, err := ctx.GetClientIdentity().GetX509Certificate()
	if err != nil {
		return "", "", fmt.Errorf("failed to get certificate: %v", err)
	}
	serialNumber := certificate.SerialNumber.String()

	return MSPID, serialNumber, nil
}

// adminCheck checks if the transaction creator is an admin.
func adminCheck(ctx contractapi.TransactionContextInterface) error {
	certificate, err := ctx.GetClientIdentity().GetX509Certificate()
	if err != nil {
		return fmt.Errorf("failed getting client certificate: %v", err)
	}
	roles := certificate.Subject.OrganizationalUnit
	for _, role := range roles {
		if role == "admin" {
			return nil
		}
	}
	return fmt.Errorf("client is not an admin")
}

// ownerCheckParticipant checks if the creator of the participant is the client.
func (s *MetadataSmartContract) ownerCheckParticipant(ctx contractapi.TransactionContextInterface, participantId int) error {
	MSPID, serialNumber, err := getCreatorInfo(ctx)
	if err != nil {
		return fmt.Errorf("failed getting creator info: %v", err)
	}

	participant, err := s.GetParticipant(ctx, participantId)
	if err != nil {
		return err
	}

	if participant.MSPID != MSPID || participant.SerialNumber != serialNumber {
		return fmt.Errorf("client is not the owner of the participant: %d", participantId)
	}
	return nil
}

// ownerCheckAggregator checks if the creator of the aggregator is the client.
func (s *MetadataSmartContract) ownerCheckAggregator(ctx contractapi.TransactionContextInterface, aggregatorId int) error {
	MSPID, serialNumber, err := getCreatorInfo(ctx)
	if err != nil {
		return fmt.Errorf("failed getting creator info: %v", err)
	}

	aggregator, err := s.GetAggregator(ctx, aggregatorId)
	if err != nil {
		return err
	}

	if aggregator.MSPID != MSPID || aggregator.SerialNumber != serialNumber {
		return fmt.Errorf("client is not the owner of the aggregator: %d", aggregatorId)
	}
	return nil
}
