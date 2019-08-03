智能合约是区块链中一个非常重要的功能和组成部分。

#### 智能合约
hyperledger ctk 中的智能合约包含了一个Chaincode代码和Chaincode管理命令这两部分。

Chaincode 代码是业务的承载体，负责具体的业务逻辑
Chaincode 管理命令负责 智能合约的部署，安装，维护等工作

###### 1.Chaincode代码
hyperledger ctk 的Chaincode是一段运行在容器中的程序。Chaincode是客户端程序和Fabric之间的桥梁。

通过Chaincode客户端程序可以发起交易，查询交易。

Chaincode是运行在Dokcer容器中，因此相对来说安全。

目前支持 java,node，go,go是最稳定的。其他还在完善。

###### 2.Chaincode的管理命令
Chaincode管理命令主要用来对Chaincode进行安装，实例化，调用，打包，签名操作。

Chaincode命令包含在Peer模块中，是peer模块中一个子命令， 该子命令的名称 是chaincode.该子命令是 peer chaincode

####  快速编写和运行一个智能合约
###### 1.创建一个Chaincode代码的目录
首先创建一个目录存放Chaincode的代码。建议放在$GOPATH指定的路径中。

```
mkdir -p $GOPATH/src/github.com/jiqiren2019/ctk/smart_contract/simpledemo
```
###### 2.创建 Chaincode源代码文件并且编写源代码
在创建域代码文件的命令如下：


```
package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"fmt"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type simplechaincode struct {

}

//智能合约初始化,在部署合约的时候会执行
func (t *simplechaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println(" ====== success it is view in docker ======")
	return shim.Success([]byte("init success"))
}

/*
* 智能合约业务逻辑处理
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

在部署之前首先通过下面的命令查看当前的Peer节点已经加入了哪些Channel

```
export set FABRIC_CFG_PATH=/opt/hyperledger/peer
export set CORE_PEER_LOCALMSPID=Org1MSP
export set CORE_PEER_ADDRESS=peer0.org1.ctk.bz:7051
export set CORE_PEER_MSPCONFIGPATH=/opt/hyperledger/fabricconfig/crypto-config/peerOrganizations/org1.ctk.bz/users/Admin@org1.ctk.bz/msp
peer channel list
```

第一步 部署

```
peer chaincode install -n democc -v 1.1 -p github.com/jiqiren2019/ctk/smart_contract/simpledemo
```

第二步 实例化

```
peer chaincode instantiate -o orderer.ctk.bz:7050 -C mychannel -n democc -v 1.1 -c '{"Args":{"init","a","100","b","200"}}' -P "AND('Org1MSP.member','Org1MSP.member','Org1MSP.member')"
```

第三步 调用

```
peer chaincode invoke -o orderer.ctk.bz:7050 -C mychannel -n democc -c '{"Args":["invoke","1","a","b"]}'
```
如果 没有出现错误，那我们就完完成了一个最简单的chaincode的代码编写部署，发布，调用的过程。



