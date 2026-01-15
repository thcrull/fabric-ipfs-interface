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

## Running the examples

It is recommended to reset the fabric network before running an example.
You can use the following command to reset it:

```bash
cd fabric-samples/test-network
./network.sh down
```

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