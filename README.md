# Fabric and IPFS Interface

A Go-based interface for integrating Hyperledger Fabric and IPFS into existing applications.

This repository provides:
- A **blockchain interface**: reusable utilities and abstractions to connect to a Fabric network and submit/evaluate transactions.
- **Chaincode packages**: smart contracts specific to the research use case.
- An **IPFS interface**: integration helpers to interact with IPFS for decentralized storage.

```bash
./install-fabric.sh docker samples binary
```

```bash
cd fabric-samples/test-network
./network.sh up createChannel
./network.sh deployCC -ccn basic -ccp ../asset-transfer-basic/chaincode-go -ccl go
```

```bash
go mod tidy
cd interface/fabric/cmd/example
go run main.go
```