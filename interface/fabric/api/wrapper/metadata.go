package fabricclient

import (
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

// AddMetadata submits a transaction to add a new metadata record
func (s *MetadataService) AddMetadata(metadata *shared.Metadata) error {
	epochStr := strconv.Itoa(metadata.Epoch)

	err := s.client.SubmitTransaction(nil, "AddMetadata", epochStr, metadata.ParticipantId, metadata.EncapsulatedKey, metadata.EncModelHash, metadata.HomomorphicHash)
	if err != nil {
		return fmt.Errorf("failed to add metadata record: %w", err)
	}

	return nil
}

// GetMetadata retrieves a metadata record by epoch and participantId
func (s *MetadataService) GetMetadata(epoch int, participantId string) (*shared.Metadata, error) {
	epochStr := strconv.Itoa(epoch)
	var metadata shared.Metadata

	err := s.client.EvaluateTransaction(&metadata, "GetMetadata", epochStr, participantId)
	if err != nil {
		return nil, fmt.Errorf("failed to query the metadata record: %w", err)
	}

	return &metadata, nil
}

// MetadataExists returns true if a metadata record exists
func (s *MetadataService) MetadataExists(epoch int, participantId string) (bool, error) {
	epochStr := strconv.Itoa(epoch)
	var exists bool

	err := s.client.EvaluateTransaction(&exists, "MetadataExists", epochStr, participantId)
	if err != nil {
		return false, fmt.Errorf("failed to query if metadata record exists: %w", err)
	}

	return exists, nil
}

// DeleteMetadata deletes a metadata record, returns nil if successful
func (s *MetadataService) DeleteMetadata(epoch int, participantId string) error {
	epochStr := strconv.Itoa(epoch)

	err := s.client.SubmitTransaction(nil, "DeleteMetadata", epochStr, participantId)
	if err != nil {
		return fmt.Errorf("failed to delete metadata record: %w", err)
	}

	return nil
}

// UpdateMetadata updates an existing metadata record
func (s *MetadataService) UpdateMetadata(metadata *shared.Metadata) error {
	epochStr := strconv.Itoa(metadata.Epoch)

	err := s.client.SubmitTransaction(nil, "UpdateMetadata", epochStr, metadata.ParticipantId, metadata.EncapsulatedKey, metadata.EncModelHash, metadata.HomomorphicHash)
	if err != nil {
		return fmt.Errorf("failed to update metadata record: %w", err)
	}

	return nil
}

// DeleteAllMetadata deletes all metadata records, returns nil if successful
func (s *MetadataService) DeleteAllMetadata() error {
	err := s.client.SubmitTransaction(nil, "DeleteAllMetadata")
	if err != nil {
		return fmt.Errorf("failed to delete metadata record: %w", err)
	}

	return nil
}

// GetAllMetadata queries all metadata records from the ledger
func (s *MetadataService) GetAllMetadata() ([]shared.Metadata, error) {
	var metadataList []shared.Metadata

	err := s.client.EvaluateTransaction(&metadataList, "GetAllMetadata")
	if err != nil {
		return nil, fmt.Errorf("failed to query all metadata records: %w", err)
	}
	return metadataList, nil
}

// GetAllMetadataByParticipant queries all metadata records from the ledger made by the participant
func (s *MetadataService) GetAllMetadataByParticipant(participantId string) ([]shared.Metadata, error) {
	var metadataList []shared.Metadata

	err := s.client.EvaluateTransaction(&metadataList, "GetAllMetadataByParticipant", participantId)
	if err != nil {
		return nil, fmt.Errorf("failed to query the metadata records: %w", err)
	}

	return metadataList, nil
}

// GetAllMetadataByEpoch queries all metadata records from the ledger from the epoch
func (s *MetadataService) GetAllMetadataByEpoch(epoch int) ([]shared.Metadata, error) {
	epochStr := strconv.Itoa(epoch)
	var metadataList []shared.Metadata

	err := s.client.EvaluateTransaction(&metadataList, "GetAllMetadataByEpoch", epochStr)
	if err != nil {
		return nil, fmt.Errorf("failed to query the metadata records: %w", err)
	}

	return metadataList, nil
}

func (s *MetadataService) Close() error {
	return s.client.Close()
}
