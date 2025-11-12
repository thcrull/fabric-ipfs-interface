# Fabric and IPFS Interface

A Go-based interface for integrating Hyperledger Fabric and IPFS into existing applications.

This repository provides:
- A **blockchain interface**: reusable utilities and abstractions to connect to a Fabric network and submit/evaluate transactions.
- **Chaincode packages**: smart contracts specific to the research use case.
- An **IPFS interface**: integration helpers to interact with IPFS for decentralised storage.

----------------------------------

### To run the Fabric interface example
All commands should be run from the root of this repository.

If the fabric samples are not installed, run the following command:
```bash
./install-fabric.sh docker samples binary
```

To start the fabric network and deploy the chaincode:
```bash
cd fabric-samples/test-network
./network.sh up createChannel
./network.sh deployCC -ccn basic -ccp ../asset-transfer-basic/chaincode-go -ccl go
```

To run the fabric example:
```bash
go mod tidy
cd interface/fabric/example
go run main.go
```
----------------------------------

### To run the IPFS interface example
All commands should be run from the root of this repository.

If the IPFS has not been instantiated, run the commands:
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

To run the example:
```bash
cd interface/ipfs/example
go run main.go
```

----------------------------------

### To run the example application
All commands should be run from the root of this repository,
and both the fabric and IPFS networks should be running (see above).

```bash
go mod tidy
cd interface/example
go run main.go
```