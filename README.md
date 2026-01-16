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

To make sure dependencies are in order:
```bash
go mod tidy
```

If the **fabric-samples** are not installed, run the following command:
```bash
./install-fabric.sh docker samples binary
```

Make sure the fabric-samples basic chaincode uses our chaincode instead of the standard:
```text
1. Go to fabric-samples/asset-transfer-basic/chaincode-go/asset_transfer.go
2. Change the line "assetChaincode, err := contractapi.NewChaincode(&chaincode.SmartContract{})" to "assetChaincode, err := contractapi.NewChaincode(&chaincode.MetadataSmartContract{})"
3. Make sure the chaincode.MetadataSmartContract from step 2 is imported from this repository
```

You have to make sure the Fabric network is running:
```bash
cd fabric-samples/test-network
./network.sh down
./network.sh up createChannel
./network.sh deployCC -ccn basic -ccp ../asset-transfer-basic/chaincode-go -ccl go
```
**Note**: If you run into any Fabric-related errors like "... failed to endorse transaction ...", you can reuse this command to reset the ledger.


If IPFS Kubo is not installed, run the commands:
```bash
tar -xvzf kubo_v0.38.1_linux-amd64.tar.gz
cd kubo
sudo bash install.sh
ipfs init
```

To start the IPFS daemon:
```bash
ipfs daemon
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
