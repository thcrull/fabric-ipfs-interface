# Hyperledger Fabric and IPFS Interface

This repository, created by **Thomas Crull**, is part of an Auditable Federated Learning research project
lead by **Dr. Roland Kromes** at the **Research Engineering and Infrastructure Team TU Delft**.
It contains a wrapper around the Hyperledger Fabric Gateway, a suite of smart contracts meant
for the Fabric network, and a wrapper around the IPFS (Kubo) RPC API.

### Repository structure and Features
```text
fabric-ipfs-interface/
├── chaincode/                    # Smart contracts for the Fabric network meant for Auditable FL
├── interface/                    # Wrappers around Fabric and IPFS Gateway APIs
│   ├── fabric/
│   │   ├── config/               # Config loader for the Fabric wrapper
│   │   └── wrapper/              # Hyperledger Fabric Gateway wrapper
│   │       ├── fabric_client.go  # General-use Gateway wrapper
│   │       └── metadata.go       # FabricClient wrapper which eases the use of the created chaincode
│   └── ipfs/
│       ├── config/               # Config loader for the IPFS wrapper
│       └── wrapper/              # IPFS RPC API wrapper
├── shared/                       # Shared type definitions
├── weight_pb/                    # Protobuf definition for the models
├── example/                      # Example app using Fabric and IPFS interfaces
├── config/                       # Configuration files for examples and tests
├── testing_utils/                # Test utilities
│   └── generate_model/           # Generates random models in data/ for tests and examples
└── data/                         # Random models used by tests and examples
```
----------------------------------

## Running the example and tests

### Prerequisites

#### Make sure dependencies are in order:
```bash
go mod tidy
```

#### If the **fabric-samples** are not installed, run the following command:
```bash
./install-fabric.sh docker samples binary
```

#### Make sure the fabric-samples basic chaincode uses our chaincode instead of the standard:
```text
1. Go to fabric-samples/asset-transfer-basic/chaincode-go/asset_transfer.go
2. Change the line "assetChaincode, err := contractapi.NewChaincode(&chaincode.SmartContract{})" to "assetChaincode, err := contractapi.NewChaincode(&chaincode.MetadataSmartContract{})"
3. Make sure the chaincode.MetadataSmartContract from step 2 is imported from this repository
```

#### Make sure the Fabric network is running:
```bash
cd fabric-samples/test-network
./network.sh down
./network.sh up createChannel
./network.sh deployCC -ccn basic -ccp ../asset-transfer-basic/chaincode-go -ccl go
```
**Note**: If you run into any Fabric-related errors like "... failed to endorse transaction ...", you can reuse this command to reset the ledger.

#### If the config files have not been added, create the following in **config/**:

*admin.yaml*:
```text
identity:
  cert_path: "/path/to/fabric-ipfs-interface/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/signcerts/Admin@org1.example.com-cert.pem"
  key_path: "/path/to/fabric-ipfs-interface/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/keystore/priv_sk"
  msp_id: "Org1MSP"

network:
  peer_endpoint: "localhost:7051"
  tls_cert_path: "/path/to/fabric-ipfs-interface/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt"
  tls_hostname: "peer0.org1.example.com"
  channel_name: "mychannel"
  chaincode_name: "basic"

ipfs:
  node_path: "http://localhost:5001"
```

*user1.yaml*:
```text
identity:
  cert_path: "/path/to/fabric-ipfs-interface/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/signcerts/User1@org1.example.com-cert.pem"
  key_path: "/path/to/fabric-ipfs-interface/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/keystore/priv_sk"
  msp_id: "Org1MSP"

network:
  peer_endpoint: "localhost:7051"
  tls_cert_path: "/path/to/fabric-ipfs-interface/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt"
  tls_hostname: "peer0.org1.example.com"
  channel_name: "mychannel"
  chaincode_name: "basic"

ipfs:
  node_path: "http://localhost:5001"
```

*user2.yaml*:
```text
identity:
  cert_path: "/path/to/fabric-ipfs-interface/fabric-samples/test-network/organizations/peerOrganizations/org2.example.com/users/User1@org2.example.com/msp/signcerts/User1@org2.example.com-cert.pem"
  key_path: "/path/to/fabric-ipfs-interface/fabric-samples/test-network/organizations/peerOrganizations/org2.example.com/users/User1@org2.example.com/msp/keystore/priv_sk"
  msp_id: "Org2MSP"

network:
  peer_endpoint: "localhost:9051"
  tls_cert_path: "/path/to/fabric-ipfs-interface/fabric-samples/test-network/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt"
  tls_hostname: "peer0.org2.example.com"
  channel_name: "mychannel"
  chaincode_name: "basic"

ipfs:
  node_path: "http://localhost:5001"
```

#### If IPFS Kubo is not installed, run the commands:
```bash
tar -xvzf kubo_v0.38.1_linux-amd64.tar.gz
cd kubo
sudo bash install.sh
ipfs init
```

#### To start the IPFS daemon:
```bash
ipfs daemon
```

#### If the weight models in data/ are not present:
```bash
cd testing_utils/generate_model
go run main.go
```

----------------------------------

### To run the example application

To run the fabric example:
```bash
cd example
go run main.go
```

----------------------------------

### To run the benchmark test

```bash
cd bench
go test ./bench -bench=. -benchmem
```
