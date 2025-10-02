package metadata

import (
	"log"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
	"github.com/thcrull/fabirc-interface/fabric-interface/chaincode/metadata/chaincode"
)

// Starts up the metadata chaincode
func main() {
	metadataChaincode, err := contractapi.NewChaincode(&chaincode.MetadataSmartContract{})
	if err != nil {
		log.Panicf("Error creating metadata chaincode: %v", err)
	}

	if err := metadataChaincode.Start(); err != nil {
		log.Panicf("Error starting metadata chaincode: %v", err)
	}
}
