package fabricclient

import (
	"os"
	"strconv"
	"testing"
)

var testMetadataServiceUser1 *MetadataService
var testMetadataServiceUser2 *MetadataService
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
	var err error

	// Create the metadata service clients with the Client role
	testMetadataServiceUser1, err = NewMetadataService("../../../../config/user1.yaml")
	if err != nil {
		panic("failed to create metadata service: " + err.Error())
	}

	testMetadataServiceUser2, err = NewMetadataService("../../../../config/user2.yaml")
	if err != nil {
		panic("failed to create metadata service: " + err.Error())
	}

	// Create the metadata service client with the Admin role
	testMetadataServiceAdmin, err = NewMetadataService("../../../../config/admin.yaml")
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
	err = testMetadataServiceUser1.Close()
	if err != nil {
		panic("failed to close user one's metadata service: " + err.Error())
	}

	err = testMetadataServiceUser2.Close()
	if err != nil {
		panic("failed to close user two's metadata service: " + err.Error())
	}

	err = testMetadataServiceAdmin.Close()
	if err != nil {
		panic("failed to close admin's metadata service: " + err.Error())
	}

	return
}

func TestGeneral(t *testing.T) {
	// -----------------------------
	// PARTICIPANT FUNCTIONALITIES
	// -----------------------------
	t.Log("-----Participant Functionalities-----")

	// User1 = Thomas
	thomasId := 1
	err := testMetadataServiceUser1.AddParticipant(
		thomasId,
		"key-thomas",
		"key-thomas-homomorphic-shared-key-cypher",
		"key-thomas-participant-communication-key-cypher",
	)
	if err != nil {
		t.Fatalf("User1 failed to add participant Thomas: %v", err)
	}
	t.Logf("Added participant Thomas with id: %d", thomasId)

	// User2 = Mihnea
	mihneaId := 2
	err = testMetadataServiceUser2.AddParticipant(
		mihneaId,
		"key-mihnea",
		"key-mihnea-homomorphic-shared-key-cypher",
		"key-mihnea-participant-communication-key-cypher",
	)
	if err != nil {
		t.Fatalf("User2 failed to add participant Mihnea: %v", err)
	}
	t.Logf("Added participant Mihnea with id: %d", mihneaId)

	// Admin = Ilinca
	ilincaId := 3
	err = testMetadataServiceAdmin.AddParticipant(
		ilincaId,
		"key-ilinca",
		"key-ilinca-homomorphic-shared-key-cypher",
		"key-ilinca-participant-communication-key-cypher",
	)
	if err != nil {
		t.Fatalf("Admin failed to add participant Ilinca: %v", err)
	}
	t.Logf("Added participant Ilinca with id: %d", ilincaId)

	t.Log("Participants added successfully.")

	// Fetch participant Thomas
	thomas, err := testMetadataServiceUser1.GetParticipant(thomasId)
	if err != nil {
		t.Fatalf("Failed to get participant Thomas: %v", err)
	}
	t.Logf("Fetched participant Thomas: %+v", thomas)

	// Update participant Thomas
	err = testMetadataServiceUser1.UpdateParticipant(
		1,
		"key-thomas-updated",
		"key-thomas-homomorphic-shared-key-cypher-updated",
		"key-thomas-participant-communication-key-cypher-updated",
	)
	if err != nil {
		t.Fatalf("Failed to update participant Thomas: %v", err)
	}
	t.Log("Updated participant Thomas successfully.")

	// Delete participant Ilinca (admin)
	err = testMetadataServiceAdmin.DeleteParticipant(ilincaId)
	if err != nil {
		t.Fatalf("Failed to delete participant Ilinca: %v", err)
	}
	t.Log("Deleted participant Ilinca successfully.")

	// Check participant Ilinca existence
	exists, err := testMetadataServiceAdmin.ParticipantExists(ilincaId)
	if err != nil {
		t.Fatalf("Failed to check participant Ilinca existence: %v", err)
	}
	t.Logf("Participant Ilinca exists? %t", exists)

	// Fetch all participants
	allParticipants, err := testMetadataServiceAdmin.GetAllParticipants()
	if err != nil {
		t.Fatalf("Failed to fetch all participants: %v", err)
	}
	t.Logf("All participants: %+v", allParticipants)

	// ----------------------------
	// AGGREGATOR FUNCTIONALITIES
	// ----------------------------
	t.Log("-----Aggregator Functionalities-----")

	// Admin = Aggregator
	aggregatorId := 4
	err = testMetadataServiceAdmin.AddAggregator(aggregatorId, map[string]string{
		strconv.Itoa(thomasId): "key-thomas-communication-key-cypher",
		strconv.Itoa(mihneaId): "key-mihnea-communication-key-cypher",
	})
	if err != nil {
		t.Fatalf("Failed to add aggregator Aggregator: %v", err)
	}
	t.Logf("Added aggregator Aggregator with id: %d", aggregatorId)

	// User1 = BadAggregator
	badAggregatorId := 5
	err = testMetadataServiceUser1.AddAggregator(badAggregatorId, map[string]string{
		strconv.Itoa(ilincaId): "key-ilinca-communication-key-cypher",
	})
	if err != nil {
		t.Fatalf("Failed to add aggregator BadAggregator: %v", err)
	}
	t.Logf("Added aggregator BadAggregator with id: %d", badAggregatorId)

	// Fetch aggregator Aggregator
	aggregator, err := testMetadataServiceAdmin.GetAggregator(aggregatorId)
	if err != nil {
		t.Fatalf("Failed to fetch aggregator Aggregator: %v", err)
	}
	t.Logf("Fetched aggregator Aggregator: %+v", aggregator)

	// Update aggregator Aggregator
	err = testMetadataServiceAdmin.UpdateAggregator(aggregatorId, map[string]string{
		strconv.Itoa(thomasId): "key-thomas-communication-key-cypher-updated",
		strconv.Itoa(mihneaId): "key-mihnea-communication-key-cypher-updated",
	})
	if err != nil {
		t.Fatalf("Failed to update aggregator Aggregator: %v", err)
	}
	t.Log("Updated aggregator Aggregator successfully.")

	// Delete aggregator BadAggregator
	err = testMetadataServiceUser1.DeleteAggregator(badAggregatorId)
	if err != nil {
		t.Fatalf("Failed to delete aggregator BadAggregator: %v", err)
	}
	t.Log("Deleted aggregator BadAggregator successfully.")

	// Check aggregator BadAggregator existence
	exists, err = testMetadataServiceUser1.AggregatorExists(badAggregatorId)
	if err != nil {
		t.Fatalf("Failed to check aggregator BadAggregator existence: %v", err)
	}
	t.Logf("Aggregator BadAggregator exists? %t", exists)

	// Fetch all aggregators
	allAggregators, err := testMetadataServiceAdmin.GetAllAggregators()
	if err != nil {
		t.Fatalf("Failed to fetch all aggregators: %v", err)
	}
	t.Logf("All aggregators: %+v", allAggregators)

	// --------------------------------------------
	// PARTICIPANT MODEL METADATA FUNCTIONALITIES
	// --------------------------------------------
	t.Log("-----Participant Model Metadata Functionalities-----")

	// Add participant model metadata
	err = testMetadataServiceUser1.AddParticipantModelMetadata(thomasId, 10, "thomas-model-cid", "thomas-model-homomorphic-hash")
	if err != nil {
		t.Fatalf("Failed to add Thomas model metadata epoch 10: %v", err)
	}

	err = testMetadataServiceUser1.AddParticipantModelMetadata(thomasId, 20, "thomas-model-cid", "thomas-model-homomorphic-hash")
	if err != nil {
		t.Fatalf("Failed to add Thomas model metadata epoch 20: %v", err)
	}

	err = testMetadataServiceUser2.AddParticipantModelMetadata(mihneaId, 10, "mihnea-model-cid", "mihnea-model-homomorphic-hash")
	if err != nil {
		t.Fatalf("Failed to add Mihnea model metadata epoch 10: %v", err)
	}

	err = testMetadataServiceUser2.AddParticipantModelMetadata(mihneaId, 20, "mihnea-model-cid", "mihnea-model-homomorphic-hash")
	if err != nil {
		t.Fatalf("Failed to add Mihnea model metadata epoch 20: %v", err)
	}

	t.Log("Added participant model metadata successfully.")

	// Fetch participant model metadata
	modelMeta, err := testMetadataServiceUser2.GetParticipantModelMetadata(mihneaId, 10)
	if err != nil {
		t.Fatalf("Failed to fetch Mihnea model metadata epoch 10: %v", err)
	}
	t.Logf("Fetched Mihnea model metadata: %+v", modelMeta)

	// Update Mihnea model metadata
	err = testMetadataServiceUser2.UpdateParticipantModelMetadata(mihneaId, 10, "mihnea-model-cid-updated", "mihnea-model-homomorphic-hash-updated")
	if err != nil {
		t.Fatalf("Failed to update Mihnea model metadata epoch 10: %v", err)
	}
	t.Log("Updated Mihnea model metadata successfully.")

	// Delete Thomas model metadata epoch 20
	err = testMetadataServiceUser1.DeleteParticipantModelMetadata(thomasId, 20)
	if err != nil {
		t.Fatalf("Failed to delete Thomas model metadata epoch 20: %v", err)
	}
	t.Log("Deleted Thomas model metadata epoch 20 successfully.")

	// Check Thomas model metadata existence epoch 20
	exists, err = testMetadataServiceUser1.ParticipantModelMetadataExists(thomasId, 20)
	if err != nil {
		t.Fatalf("Failed to check Thomas model metadata existence epoch 20: %v", err)
	}
	t.Logf("Thomas model metadata epoch 20 exists? %t", exists)

	// Fetch all participant model metadata for epoch 10
	modelsByEpoch, err := testMetadataServiceAdmin.GetAllParticipantModelMetadataByEpoch(10)
	if err != nil {
		t.Fatalf("Failed to fetch all model metadata for epoch 10: %v", err)
	}
	t.Logf("All model metadata for epoch 10: %+v", modelsByEpoch)

	// Fetch all participant model metadata for Mihnea
	modelsByParticipant, err := testMetadataServiceUser2.GetAllParticipantModelMetadataByParticipant(mihneaId)
	if err != nil {
		t.Fatalf("Failed to fetch all model metadata for Mihnea: %v", err)
	}
	t.Logf("All model metadata for Mihnea: %+v", modelsByParticipant)

	// -------------------------------------------
	// AGGREGATOR MODEL METADATA FUNCTIONALITIES
	// -------------------------------------------
	t.Log("-----Aggregator Model Metadata Functionalities-----")

	// Add aggregator model metadata (Admin)
	err = testMetadataServiceAdmin.AddAggregatorModelMetadata(aggregatorId, 10, "aggregator-model-cid", []int{thomasId, mihneaId})
	if err != nil {
		t.Fatalf("Failed to add aggregator model metadata epoch 10: %v", err)
	}

	err = testMetadataServiceAdmin.AddAggregatorModelMetadata(aggregatorId, 20, "aggregator-model-cid", []int{mihneaId})
	if err != nil {
		t.Fatalf("Failed to add aggregator model metadata epoch 20: %v", err)
	}

	t.Log("Added aggregator model metadata successfully.")

	// Fetch aggregator model metadata epoch 10
	aggMeta, err := testMetadataServiceAdmin.GetAggregatorModelMetadata(aggregatorId, 10)
	if err != nil {
		t.Fatalf("Failed to fetch aggregator model metadata epoch 10: %v", err)
	}
	t.Logf("Fetched aggregator model metadata epoch 10: %+v", aggMeta)

	// Update aggregator model metadata epoch 10
	err = testMetadataServiceAdmin.UpdateAggregatorModelMetadata(aggregatorId, 10, "aggregator-model-cid-updated", []int{thomasId, mihneaId})
	if err != nil {
		t.Fatalf("Failed to update aggregator model metadata epoch 10: %v", err)
	}
	t.Log("Updated aggregator model metadata successfully.")

	// Delete aggregator model metadata epoch 20
	err = testMetadataServiceAdmin.DeleteAggregatorModelMetadata(aggregatorId, 20)
	if err != nil {
		t.Fatalf("Failed to delete aggregator model metadata epoch 20: %v", err)
	}
	t.Log("Deleted aggregator model metadata epoch 20 successfully.")

	// Check aggregator model metadata existence epoch 20
	exists, err = testMetadataServiceAdmin.AggregatorModelMetadataExists(20, aggregatorId)
	if err != nil {
		t.Fatalf("Failed to check aggregator model metadata existence epoch 20: %v", err)
	}
	t.Logf("Aggregator model metadata epoch 20 exists? %t", exists)

	// Fetch all aggregator model metadata
	allAggMeta, err := testMetadataServiceAdmin.GetAllAggregatorModelMetadata()
	if err != nil {
		t.Fatalf("Failed to fetch all aggregator model metadata: %v", err)
	}
	t.Logf("All aggregator model metadata: %+v", allAggMeta)
}

func TestDeleteAllAdminOnly(t *testing.T) {
	// 1. Try as USER1 → SHOULD FAIL
	err := testMetadataServiceUser1.DeleteAllParticipants()
	if err == nil {
		t.Fatalf("User1 was able to call DeleteAllParticipants but should NOT have permission")
	}
	t.Log("Correctly blocked User1 from DeleteAllParticipants")

	err = testMetadataServiceUser1.DeleteAllAggregators()
	if err == nil {
		t.Fatalf("User1 was able to call DeleteAllAggregators but should NOT have permission")
	}
	t.Log("Correctly blocked User1 from DeleteAllAggregators")

	err = testMetadataServiceUser1.DeleteAllParticipantModelMetadata()
	if err == nil {
		t.Fatalf("User1 was able to call DeleteAllParticipantModelMetadata but should NOT have permission")
	}
	t.Log("Correctly blocked User1 from DeleteAllParticipantModelMetadata")

	err = testMetadataServiceUser1.DeleteAllAggregatorModelMetadata()
	if err == nil {
		t.Fatalf("User1 was able to call DeleteAllAggregatorModelMetadata but should NOT have permission")
	}
	t.Log("Correctly blocked User1 from DeleteAllAggregatorModelMetadata")

	// 2. Try as ADMIN → SHOULD PASS
	if err := testMetadataServiceAdmin.DeleteAllParticipants(); err != nil {
		t.Fatalf("Admin failed DeleteAllParticipants: %v", err)
	}
	t.Log("Admin successfully called DeleteAllParticipants")

	if err := testMetadataServiceAdmin.DeleteAllAggregators(); err != nil {
		t.Fatalf("Admin failed DeleteAllAggregators: %v", err)
	}
	t.Log("Admin successfully called DeleteAllAggregators")

	if err := testMetadataServiceAdmin.DeleteAllParticipantModelMetadata(); err != nil {
		t.Fatalf("Admin failed DeleteAllParticipantModelMetadata: %v", err)
	}
	t.Log("Admin successfully called DeleteAllParticipantModelMetadata")

	if err := testMetadataServiceAdmin.DeleteAllAggregatorModelMetadata(); err != nil {
		t.Fatalf("Admin failed DeleteAllAggregatorModelMetadata: %v", err)
	}
	t.Log("Admin successfully called DeleteAllAggregatorModelMetadata")
}

func TestLoggingAdminOnly(t *testing.T) {
	// 1. USER1 should NOT be able to read logs
	_, err := testMetadataServiceUser1.GetAllLogsWithoutCreator()
	if err == nil {
		t.Fatalf("User1 was able to call GetAllLogsWithoutCreator but should NOT have permission")
	}
	t.Log("Correctly blocked User1 from GetAllLogsWithoutCreator")

	_, err = testMetadataServiceUser1.GetAllLogs()
	if err == nil {
		t.Fatalf("User1 was able to call GetAllLogs but should NOT have permission")
	}
	t.Log("Correctly blocked User1 from GetAllLogs")

	_, err = testMetadataServiceUser1.GetAllLogsForUser("Org1MSP", "user1-serial")
	if err == nil {
		t.Fatalf("User1 was able to call GetAllLogsForUser but should NOT have permission")
	}
	t.Logf("Correctly blocked User1 from GetAllLogsForUser")

	// 3. ADMIN should access all logs successfully
	adminLogs, err := testMetadataServiceAdmin.GetAllLogsWithoutCreator()
	if err != nil {
		t.Fatalf("Admin failed GetAllLogsWithoutCreator: %v", err)
	}
	t.Logf("Admin successfully fetched logs without creator: %d entries", len(adminLogs))

	adminLogsFull, err := testMetadataServiceAdmin.GetAllLogs()
	if err != nil {
		t.Fatalf("Admin failed GetAllLogs: %v", err)
	}
	t.Logf("Admin successfully fetched logs with creator info: %d entries", len(adminLogsFull))

	adminLogsForUser, err := testMetadataServiceAdmin.GetAllLogsForUser("Org1MSP", "user1-serial")
	if err != nil {
		t.Fatalf("Admin failed GetAllLogsForUser: %v", err)
	}
	t.Logf("Admin successfully fetched logs for user (Org1MSP, user1-serial): %d entries", len(adminLogsForUser))
}
