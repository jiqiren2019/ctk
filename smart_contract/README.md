Intelligent contract is a very important function and component of blockchain.

#### smart contract
The smart contract in the hyperledger CTK contains a Chaincode code and a Chaincode management command.

Chaincode Code is the carrier of business, responsible for specific business logic
Chaincode The administrative command is responsible for the deployment, installation, maintenance and other work of intelligent contracts

###### 1.Chaincode code

Hyperledger CTK's Chaincode is a program that runs in a container.Chaincode is a bridge between the client program and Fabric.

A transaction can be initiated and interrogated through a Chaincode client program.

Chaincode is run in the Dokcer container and is therefore relatively safe.

Currently support Java,node,go,go is the most stable.Others are still being refined.

###### 2.Chaincode admin command

The Chaincode management command is primarily used to install, instantiate, invoke, package, and sign Chaincode.

The Chaincode command is contained in the Peer module and is a subcommand in the Peer module. The name of the subcommand is Chaincode

####  Write and run an intelligent contract quickly

###### 1.Create a directory of Chaincode code

First create a directory to hold the code for Chaincode.Put it in the path specified by $GOPATH.

```
mkdir -p $GOPATH/src/github.com/jiqiren2019/ctk/smart_contract/simpledemo
```
###### 2.Create the Chaincode source file and write the source code
The command to create the domain code file is as follows:


```
package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"fmt"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type simplechaincode struct {

}

//Intelligent contract initialization, which is performed when the contract is deployed
func (t *simplechaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println(" ====== success it is view in docker ======")
	return shim.Success([]byte("init success"))
}

/*
* Intelligent contract business logic processing
 */
func (t *simplechaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println(" ====== success it is view in docker ======")
	return shim.Success([]byte("invoke success"))
}

func main() {
	var chaincode = new(simplechaincode)
	err := shim.Start(chaincode)
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}

```

Before deploying, first see which channels the current Peer node has joined by using the following command

```
export set FABRIC_CFG_PATH=/opt/hyperledger/peer
export set CORE_PEER_LOCALMSPID=Org1MSP
export set CORE_PEER_ADDRESS=peer0.org1.ctk.bz:7051
export set CORE_PEER_MSPCONFIGPATH=/opt/hyperledger/fabricconfig/crypto-config/peerOrganizations/org1.ctk.bz/users/Admin@org1.ctk.bz/msp
peer channel list
```

First step deployment

```
peer chaincode install -n democc -v 1.1 -p github.com/jiqiren2019/ctk/smart_contract/simpledemo
```

Step 2 instantiate

```
peer chaincode instantiate -o orderer.ctk.bz:7050 -C mychannel -n democc -v 1.1 -c '{"Args":{"init","a","100","b","200"}}' -P "AND('Org1MSP.member','Org1MSP.member','Org1MSP.member')"
```

Step 3 call

```
peer chaincode invoke -o orderer.ctk.bz:7050 -C mychannel -n democc -c '{"Args":["invoke","1","a","b"]}'
```
If there are no errors, we are done with the simplest of procedures of writing, deploying, distributing, and calling a chaincode.



