package fabricclient

import (
	"os"
	"testing"

	fabricconfig "github.com/thcrull/fabric-ipfs-interface/interface/fabric/api/config"
)

var testMetadataServiceClient *MetadataService
var testMetadataServiceAdmin *MetadataService

// TestMain runs before all tests
func TestMain(m *testing.M) {
	// Set up the metadata service clients
	setup()

	// Run the tests
	code := m.Run()

	// Tear down the metadata service clients
	teardown()

	os.Exit(code)
}

func setup() {
	// Create the metadata service client with the Client role
	cfg, err := fabricconfig.LoadConfig("config.yaml")
	if err != nil {
		panic("failed to fetch config file: " + err.Error())
	}

	testMetadataServiceClient, err = NewMetadataService(cfg)
	if err != nil {
		panic("failed to create metadata service: " + err.Error())
	}

	// Create the metadata service client with the Admin role
	cfg, err = fabricconfig.LoadConfig("config.yaml")
	if err != nil {
		panic("failed to fetch config file: " + err.Error())
	}

	testMetadataServiceAdmin, err = NewMetadataService(cfg)
	if err != nil {
		panic("failed to create metadata service: " + err.Error())
	}

	return
}

func teardown() {
	// Delete all records from the ledger
	err := testMetadataServiceAdmin.DeleteAllParticipantModelMetadata()
	if err != nil {
		panic("failed to delete all participant model metadata records: " + err.Error())
	}

	err = testMetadataServiceAdmin.DeleteAllAggregatorModelMetadata()
	if err != nil {
		panic("failed to delete all aggregator model metadata records: " + err.Error())
	}

	err = testMetadataServiceAdmin.DeleteAllParticipants()
	if err != nil {
		panic("failed to delete all participants records: " + err.Error())
	}

	err = testMetadataServiceAdmin.DeleteAllAggregators()
	if err != nil {
		panic("failed to delete all aggregators records: " + err.Error())
	}

	// Close the metadata service clients
	err = testMetadataServiceClient.Close()
	if err != nil {
		panic("failed to close client's metadata service: " + err.Error())
	}

	err = testMetadataServiceAdmin.Close()
	if err != nil {
		panic("failed to close admin's metadata service: " + err.Error())
	}

	return
}

func TestAddParticipant(t *testing.T) {
	
}
