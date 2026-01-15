# Hyperledger Fabric and IPFS Interface

This repository, created by **Thomas Crull**, is part of an Auditable Federated Learning research project
lead by **Dr. Roland Kromes** at the **Research Engineering and Infrastructure Team TU Delft**.
It contains a wrapper around the Hyperledger Fabric Gateway, a suite of smart contracts meant
for the Fabric network, and a wrapper around the IPFS (Kubo) RPC API.

### Repository structure and Features
- **/chaincode**: Suite of smart contracts for the Fabric network relevant for Auditable Federated Learning
- **/interface**: Wrappers around the Fabric and IPFS Gateway APIs
  - **/fabric**
    - **/config**: Config loader for the Fabric wrapper
    - **/wrapper**: Hyperledger Fabric Gateway wrapper
      - **fabric_client.go**: General use wrapper meant to ease the use of the Gateway
      - **metadata.go**: Specialised wrapper built on top of the FabricClient which enables the use of the created chaincode 
  - **/ipfs**
    - **/config**: Config loader for the IPFS wrapper
    - **/wrapper**: General use wrapper meant to ease the use of the IPFS RPC API
- **/shared**: All types defined within the repository
- **/example**: Basic example application using the Fabric and IPFS interfaces
- **/config**: Directory with configuration files used by the examples and tests
- **/testing_utils**: Utilities used by the tests
  - **/generate_model**: Executable that generates the random models for the **/data** directory which are used by the tests
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

You have to make sure the Fabric network is running:
```bash
cd fabric-samples/test-network
./network.sh down
./network.sh up createChannel
./network.sh deployCC -ccn basic -ccp ../asset-transfer-basic/chaincode-go -ccl go
```
**Note**: If you run into any Fabric-related errors like "... failed to endorse transaction ...", you can reuse this command to reset the ledger.

----------------------------------

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
