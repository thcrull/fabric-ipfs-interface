package metadata

import (
	"fmt"
	"strconv"

	"github.com/thcrull/fabric-ipfs-interface/shared"
)

// MetadataService wraps a Fabric client and provides convenient methods
// for interacting with metadata chaincode.
type MetadataService struct {
	client *Client
}

// NewMetadataService creates a service for metadata operations
func NewMetadataService(client *Client) *MetadataService {
	return &MetadataService{client: client}
}

// AddMetadata submits a transaction to add a new metadata record
func (s *MetadataService) AddMetadata(
	epoch int,
	participantID string,
	encapsulatedKey string,
	encModelHash string,
	homomorphicHash string,
) (bool, error) {
	epochStr := strconv.Itoa(epoch)

	err := s.client.SubmitTransaction(nil, "AddMetadata", epochStr, participantID, encapsulatedKey, encModelHash, homomorphicHash)
	if err != nil {
		return false, fmt.Errorf("failed to add metadata record: %w", err)
	}

	return true, nil
}

// ReadMetadata retrieves a metadata record by epoch and participantID
func (s *MetadataService) ReadMetadata(epoch int, participantID string) (*shared.Metadata, error) {
	epochStr := strconv.Itoa(epoch)
	var metadata shared.Metadata

	err := s.client.EvaluateTransaction(&metadata, "ReadMetadata", epochStr, participantID)
	if err != nil {
		return nil, fmt.Errorf("failed to query the metadata record: %w", err)
	}

	return &metadata, nil
}

// MetadataExists returns true if a metadata record exists
func (s *MetadataService) MetadataExists(epoch int, participantID string) (bool, error) {
	epochStr := strconv.Itoa(epoch)
	var exists bool

	err := s.client.EvaluateTransaction(&exists, "MetadataExists", epochStr, participantID)
	if err != nil {
		return false, fmt.Errorf("failed to query if metadata record exists: %w", err)
	}

	return exists, nil
}

// DeleteMetadata deletes a metadata record, returns true if successful
func (s *MetadataService) DeleteMetadata(epoch int, participantID string) (bool, error) {
	epochStr := strconv.Itoa(epoch)

	err := s.client.SubmitTransaction(nil, "DeleteMetadata", epochStr, participantID)
	if err != nil {
		return false, fmt.Errorf("failed to delete metadata record: %w", err)
	}

	return true, nil
}

// UpdateMetadata updates an existing metadata record
func (s *MetadataService) UpdateMetadata(
	epoch int,
	participantID string,
	encapsulatedKey string,
	encModelHash string,
	homomorphicHash string,
) (bool, error) {
	epochStr := strconv.Itoa(epoch)

	err := s.client.SubmitTransaction(nil, "UpdateMetadata", epochStr, participantID, encapsulatedKey, encModelHash, homomorphicHash)
	if err != nil {
		return false, fmt.Errorf("failed to update metadata record: %w", err)
	}

	return true, nil
}

// GetAllMetadata queries all metadata entries from the ledger
func (s *MetadataService) GetAllMetadata() ([]shared.Metadata, error) {
	var metadataList []shared.Metadata

	err := s.client.EvaluateTransaction(&metadataList, "GetAllMetadata")
	if err != nil {
		return nil, fmt.Errorf("failed to query all metadata records: %w", err)
	}
	return metadataList, nil
}
