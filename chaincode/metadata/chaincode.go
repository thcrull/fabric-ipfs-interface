package metadata

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// AddMetadata issues a new metadata block to the world state with given details.
func (s *MetadataSmartContract) AddMetadata(
	ctx contractapi.TransactionContextInterface,
	epoch int,
	participantId string,
	encapsulatedKey string,
	encModelHash string,
	homomorphicHash string,
) error {
	compositeKey, err := ctx.GetStub().CreateCompositeKey("metadata", []string{participantId, fmt.Sprintf("%d", epoch)})

	exists, err := s.MetadataExists(ctx, epoch, participantId)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the metadata block for epoch %d from participant %s already exists", epoch, participantId)
	}

	metadata := Metadata{
		Epoch:           epoch,
		ParticipantID:   participantId,
		EncapsulatedKey: encapsulatedKey,
		EncModelHash:    encModelHash,
		HomomorphicHash: homomorphicHash,
	}
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(compositeKey, metadataJSON)
}

// ReadMetadata returns the metadata block stored in the world state with the given epoch and participant id.
func (s *MetadataSmartContract) ReadMetadata(ctx contractapi.TransactionContextInterface, epoch int, participantId string) (*Metadata, error) {
	compositeKey, err := ctx.GetStub().CreateCompositeKey("metadata", []string{participantId, fmt.Sprintf("%d", epoch)})

	metadataJSON, err := ctx.GetStub().GetState(compositeKey)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if metadataJSON == nil {
		return nil, fmt.Errorf("the metadata block for epoch %d from participant %s does not exist", epoch, participantId)
	}

	var metadata Metadata
	err = json.Unmarshal(metadataJSON, &metadata)
	if err != nil {
		return nil, err
	}

	return &metadata, nil
}

// MetadataExists returns true when a metadata block with the given epoch and participantId exists in the world state
func (s *MetadataSmartContract) MetadataExists(ctx contractapi.TransactionContextInterface, epoch int, participantId string) (bool, error) {
	compositeKey, err := ctx.GetStub().CreateCompositeKey("metadata", []string{participantId, fmt.Sprintf("%d", epoch)})

	metadataJSON, err := ctx.GetStub().GetState(compositeKey)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return metadataJSON != nil, nil
}

// GetAllMetadata returns all metadata blocks found in the world state
func (s *MetadataSmartContract) GetAllMetadata(ctx contractapi.TransactionContextInterface) ([]*Metadata, error) {
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("metadata", []string{})
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var metadataBlocks []*Metadata
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var metadata Metadata
		if err := json.Unmarshal(queryResponse.Value, &metadata); err != nil {
			return nil, err
		}
		metadataBlocks = append(metadataBlocks, &metadata)
	}

	return metadataBlocks, nil
}
