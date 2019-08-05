package main
import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/util/decimal"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/ethereum/go-ethereum/crypto"
	"time"
	"errors"
	"strings"
	"encoding/hex"
	"bytes"
	"encoding/gob"
	"strconv"
	"regexp"
	"math/rand"
	"encoding/base64"
	"crypto/md5"
	"sort"
	"github.com/bitly/go-simplejson"
)
type TokenConfig struct {
	originatorAccount string
	name              string
	title             string
	logo              string
	gross             decimal.Decimal
	output            decimal.Decimal
	mineral           decimal.Decimal
	precision         int64
	rate              decimal.Decimal
	dayAward          decimal.Decimal
	site              string
	email             string
}
type TokenInfo struct {
	creator string
	cc                string
	name              string
	desc              string
	logo              string
	total             string
	award             string
	balance           string
	decimal           string
	mineral           string
	url               string
	email             string
	publishTime       string
}
type TokenChaincode struct {
	rand                     *rand.Rand
	dateZone                 *time.Location
	chargeLowLimit           decimal.Decimal
	chargeDestroy            decimal.Decimal
	chargeDestroyUpperLimit  decimal.Decimal
	tokenCharge              decimal.Decimal
	exchangeCharge           decimal.Decimal
	superMinerTokenAwardAll  decimal.Decimal
	normalMinerTokenAwardAll decimal.Decimal
	tokenChargeDestroy       decimal.Decimal
	mainToken 					 TokenConfig
	cid                        string
	blackDestroyDateKey        string
	minersFeeDateKey		string
	addressPrefix     string
	mortgageChaincode string
	candyChaincode    string
	maxPrecision      int64
	maxCurrency       string
	configKey         string
	mineralKey        string
	blackKey          string
	grossKey          string
	outputKey         string
	rateKey           string
	dayRateKey        string
	precisionKey      string
	tokenAwardKey     string
	tokenTransPrefix  string
	mortgPrefix       string
	candyPrefix       string
	agentPayPrefix	  string
	feeAccount        float64
	storageChainCode  string
	storageFeeAccount string
	tokenMortAmount   string
	tokenMortPre      string
	tokenMortDays     string
	tokenMortRate     string
	tokenSpareNodeAccount string
	tokenFoundAccount string
	minerFeeAccount	   string
}
func (self *TokenChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return self.InitMainMoney(stub)
}
func (self *TokenChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "transfer" {
		return self.transfer(stub, args)
	} else if function == "batTransfer" {
		return self.batTransfer(stub, args)
	} else if function == "publish" {
		return self.publish(stub, args)
	} else if function == "balance" {
		return self.balance(stub, args)
	} else if function == "info" {
		return self.info(stub, args)
	} else if function == "award" {
		return self.award(stub, args)
	} else if function == "mineral" {
		return self.getMineral(stub, args)
	} else if function == "awardPreview" {
		return self.awardPreview(stub, args)
	} else if function == "awardList" {
		return self.awardList(stub, args)
	} else if function == "mortgagePut" {
		return self.mortgagePut(stub, args)
	} else if function == "mortgageGet" {
		return self.mortgageGet(stub, args)
	} else if function == "mortgageGetView" {
		return self.mortgageGetView(stub, args)
	} else if function == "candyPut" {
		return self.candyPut(stub, args)
	} else if function == "widthDrawCandy" {
		return self.candyWidthDraw(stub, args)
	} else if function == "candyMortgageGet" {
		return self.candyMortgageGet(stub, args)
	} else if function == "moneyDestroy" {
		return self.moneyDestroy(stub, args)
	} else if function == "moneyDestroyList" {
		return self.moneyDestroyList(stub, args)
	} else if function == "minersFeeInfo" {
		return self.minersFeeInfo(stub, args)
	} else if function == "minersFeeGive" {
		return self.minersFeeGive(stub, args)
	} else if function == "applySpareNode"{
		return self.applySpareNode(stub,args)
	}else if function == "unitFoundTransfer"{
		return self.applyUnitFoundTransfer(stub,args)
	}else if function == "nodeMortgReturn"{
		return self.nodeMortReturn(stub,args)
	}
	return shim.Error("Method is invalid")
}
func (self *TokenChaincode) InitMainMoney(stub shim.ChaincodeStubInterface) pb.Response {
	token := self.mainToken
	mainMoneyConfigbytes, err := stub.GetState(token.name + self.configKey)
	if err != nil {
		return shim.Error("Check the main money configuration for errors.")
	}
	if mainMoneyConfigbytes != nil {
		return shim.Success([]byte("Automatic skip has been initialized"))
	}
	err = self.updateAccountBalance(stub, token.name, token.originatorAccount, token.output, 6)
	if err != nil {
		return shim.Error("Error 1 writing token data")
	}
	mainMoneyMineral := token.gross.Sub(token.output)
	err = self.updateAccountBalance(stub, token.name, self.mineralKey, mainMoneyMineral, token.precision)
	if err != nil {
		return shim.Error("Error 2 writing token data")
	}
	err = stub.PutState(token.name+self.grossKey, []byte( token.gross.String() ))
	if err != nil {
		return shim.Error("Error 3 writing token data")
	}
	err = stub.PutState(token.name+self.outputKey, []byte( token.output.String() ))
	if err != nil {
		return shim.Error("Error 4 writing token data")
	}
	err = stub.PutState(token.name+self.rateKey, []byte( token.rate.String() ))
	if err != nil {
		return shim.Error("Error 5 writing token data")
	}
	err = stub.PutState(token.name+self.dayRateKey, []byte( token.dayAward.String() ))
	if err != nil {
		return shim.Error("Error 6 writing token data")
	}
	err = stub.PutState(token.name+self.precisionKey, []byte( strconv.FormatInt(token.precision, 10) ))
	if err != nil {
		return shim.Error("Error 7 writing token data")
	}
	awardKey := strings.ToLower(token.name + self.tokenAwardKey + token.originatorAccount)
	var stateString = "0,0," + self.getCurrentDateString()
	err = stub.PutState(awardKey, []byte( stateString ));
	if err != nil {
		return shim.Error("Error saving account time")
	}
	var config [13] string
	config[0] = token.originatorAccount
	config[1] = token.name
	config[2] = token.name
	config[3] = token.gross.String()
	config[4] = token.output.String()
	config[5] = strconv.FormatInt(token.precision, 10)
	config[6] = token.site
	config[7] = token.logo
	config[8] = token.email
	config[9] = token.rate.String()
	config[10] = self.formatNumber(self.tokenCharge, token.precision)
	config[11] = "0"
	config[12] = self.getCurrentDateString()
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	err = enc.Encode(config)
	if err != nil {
		return shim.Error("The token information encoding failed")
	}
	err = stub.PutState(token.name+self.configKey, network.Bytes())
	if err != nil {
		return shim.Error("The token information failed to be saved to the state")
	}
	return shim.Success([]byte("init mainmoney success"))
}
func (self *TokenChaincode) moneyDestroyList(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	historyQuery,err := stub.GetHistoryForKey(strings.ToLower(self.blackDestroyDateKey))
	if err != nil{
		return shim.Error( err.Error() )
	}
	var content = ""
	for historyQuery.HasNext(){
		keyModification,err := historyQuery.Next();if err != nil{
			return shim.Error( err.Error() )
		}
		value := string(keyModification.GetValue())
		resultArray := strings.Split(value, ";")
		if len(resultArray) >= 3{
			if content != "" {
				content += ";"
			}
			content += resultArray[2]
		}
	}
	return shim.Success([]byte( content ))
}
func (self *TokenChaincode) minersFeeInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	minersFeeValbytes, err := stub.GetState(strings.ToLower(self.minersFeeDateKey))
	if err != nil {
		return shim.Error("get " + self.minersFeeDateKey + " state error")
	}
	return shim.Success(minersFeeValbytes)
}
func (self *TokenChaincode) minersFeeGive(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("minersFeeGive£¨£© args require 1 £¨ [1] blockHeight  £©")
	}
	if args[0] == ""  {
		return shim.Error("The parameter cannot be empty.")
	}
	blockHeight, err := strconv.Atoi(args[0]);
	if err != nil {
		return shim.Error("args blockHeight formate error" + err.Error())
	}
	var beginBlockNumber uint64 = 1
	var endBlockNumber uint64 = uint64(blockHeight)
	if beginBlockNumber >= endBlockNumber {
		return shim.Error("endBlockNumber must be greater than beginBlockNumber")
	}
	minersFeeValbytes, err := stub.GetState(strings.ToLower(self.minersFeeDateKey))
	if err != nil {
		return shim.Error("get " + self.minersFeeDateKey + " state error")
	}
	if minersFeeValbytes != nil {
		minersFeeInfo := string(minersFeeValbytes)
		minersFeePayload := strings.Split(minersFeeInfo, ";")
		minersFeeResult := strings.Split(minersFeePayload[2], ",")
		_lastEndBlockNum, err := strconv.Atoi(minersFeeResult[2]);
		if err != nil {
			return shim.Error("lastEndBlockNum format error")
		}
		lastEndBlockNum := uint64(_lastEndBlockNum)
		beginBlockNumber = lastEndBlockNum+1
		if lastEndBlockNum >= endBlockNumber  {
			return shim.Error("last time endBlockNum is " + minersFeeResult[2] + " This time it needs to be bigger than that")
		}
	}else{
		if beginBlockNumber != 1{
			return shim.Error("first time beginBlockNum You have to start at 1")
		}
	}
	myargs := [][]byte{[]byte("GetServiceCharge"), []byte(self.cid), []byte(strconv.FormatUint(beginBlockNumber,10)), []byte(strconv.FormatUint(endBlockNumber,10))}
	response := stub.InvokeChaincode("qscc", myargs, self.cid)
	if response.Status != 200 {
		return shim.Error(" invoke qscc GetServiceCharge error:" + response.Message)
	}
	accountList, err := simplejson.NewJson(response.Payload)
	if err != nil {
		return shim.Error("NewJson:" + err.Error())
	}
	if len(accountList.MustMap()) == 0{
		return shim.Error("No data is automatically skipped")
	}
	for account,_award := range accountList.MustMap(){
		balance,err := self.getAccountBalance(stub,self.mainToken.name, account)
		if err != nil {
			fmt.Println("error:", err.Error())
			return shim.Error("get Account("+account+") Balance Error: " + err.Error())
		}
		award, err := decimal.NewDecimalFromString(_award.(string))
		if err != nil {
			return shim.Error(" NewDecimalFromString error" + err.Error())
		}
		balance  = balance.Add(award)
		self.updateAccountBalance(stub,self.mainToken.name,account,balance,self.mainToken.precision)
	}
	txid := stub.GetTxID()
	listEncode := base64.StdEncoding.EncodeToString(response.Payload)
	saveMinersFeeInfo := self.mainToken.name+";minersFee;" + self.getCurrentDateString() + "," + strconv.FormatUint(beginBlockNumber,10) + "," + strconv.FormatUint(endBlockNumber,10) + "," + txid + "," + self.minerFeeAccount + ","+ listEncode
	err = stub.PutState(strings.ToLower(self.minersFeeDateKey), []byte( saveMinersFeeInfo ))
	if err != nil {
		return shim.Error("PutState MinersFeeInfo error:" + err.Error())
	}
	return shim.Success([]byte(saveMinersFeeInfo))
}
func (self *TokenChaincode) moneyDestroy(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("moneyDestroy£¨£© args require 1 £¨ [1] blockHeight  £©")
	}
	if args[0] == ""  {
		return shim.Error("The parameter cannot be empty.")
	}
	blockHeight, err := strconv.Atoi(args[0]);
	if err != nil {
		return shim.Error("args blockHeight formate error" + err.Error())
	}
	var beginBlockNumber uint64 = 1
	var endBlockNumber uint64 = uint64(blockHeight)
	if beginBlockNumber >= endBlockNumber {
		return shim.Error("endBlockNumber must be greater than beginBlockNumber")
	}
	blackDestroyValbytes, err := stub.GetState(strings.ToLower(self.blackDestroyDateKey))
	if err != nil {
		return shim.Error("get " + self.blackDestroyDateKey + " state error")
	}
	if blackDestroyValbytes != nil {
		blackDestroyInfo := string(blackDestroyValbytes)
		fmt.Println("blackDestroyInfo",blackDestroyInfo)
		blackDestroyPayload := strings.Split(blackDestroyInfo, ";")
		blackDestroyResult := strings.Split(blackDestroyPayload[2], ",")
		_lastEndBlockNum, err := strconv.Atoi(blackDestroyResult[2]);
		if err != nil {
			return shim.Error("lastEndBlockNum format error")
		}
		lastEndBlockNum := uint64(_lastEndBlockNum)
		beginBlockNumber = lastEndBlockNum+1
		if lastEndBlockNum >= endBlockNumber  {
			return shim.Error("last time endBlockNum is " + blackDestroyResult[2] + " This time it needs to be bigger than that")
		}
	}else{
		if beginBlockNumber != 1{
			return shim.Error("first time beginBlockNum You have to start at 1")
		}
	}
	myargs := [][]byte{[]byte("GetBlockMoney"), []byte(self.cid), []byte(strconv.FormatUint(beginBlockNumber,10)), []byte(strconv.FormatUint(endBlockNumber,10))}
	response := stub.InvokeChaincode("qscc", myargs, self.cid)
	if response.Status != 200 {
		return shim.Error(" invoke qscc GetBlockMoney failed error:" + response.Message)
	}
	moneyDestroy, err := decimal.NewDecimalFromString(string(response.Payload))
	if err != nil {
		return shim.Error(" NewDecimalFromString error" + err.Error())
	}
	blackAccountVal, err := self.getAccountBalance(stub,self.mainToken.name, self.blackKey)
	if err != nil {
		fmt.Println("error:", err.Error())
		return shim.Error("error: " + err.Error())
	}
	blackAccountVal = blackAccountVal.Add(moneyDestroy)
	err = self.updateAccountBalance(stub,self.mainToken.name,self.blackKey,blackAccountVal,self.mainToken.precision)
	if err != nil {
		return shim.Error("PutState blackAccount error:" + err.Error())
	}
	saveBlackDestroyInfo := self.mainToken.name+";moneyDestroy;" + self.getCurrentDateString() + "," + strconv.FormatUint(beginBlockNumber,10) + "," + strconv.FormatUint(endBlockNumber,10) + "," + moneyDestroy.String() + "," + blackAccountVal.String()
	err = stub.PutState(strings.ToLower(self.blackDestroyDateKey), []byte( saveBlackDestroyInfo ));
	if err != nil {
		return shim.Error("PutState BlackDestroyInfo error:" + err.Error())
	}
	return shim.Success([]byte(saveBlackDestroyInfo))
}
func (self *TokenChaincode) batTransfer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var outAccount, token, inAccount, remark, outMoney, minerFee,goRatio string
	var realBlackMoney, superMinnerAward, normalMinnerAward, poolMinnerAward ,minerFeeVal decimal.Decimal
	var normalMinerAccount, superMinerAccount,poolMinerAccount string
	var err error
	var checkExists bool = false
	var zeroDecimal = decimal.NewFromFloat(0)
	if len(args) != 8 {
		for i := 0; i < len(args); i++ {
			fmt.Println("args[", i, "] = ", args[i])
		}
		fmt.Println("batTransfer£¨£© parameter require 8 £¨ [1] token [2] inAcount [3] outAccount [4] number [5] sign [6] minerFee [7] remark [8] goRatio£©")
		return shim.Error("Missing parameter")
	}
	if args[0] == "" || args[1] == "" || args[2] == "" || args[3] == "" || args[4] == "" || args[5] == "" || args[7] == "" {
		return shim.Error("Parameters cannot be empty")
	}
	token = strings.ToLower(args[0])
	inAccount = strings.ToLower(args[1])
	outAccount = strings.ToLower(args[2])
	outMoney = strings.ToLower(args[3])
	minerFee = strings.ToLower(args[5])
	remark = args[6]
	goRatio = args[7]
	superMinerAccount = strings.ToLower(stub.GetSuperMiner())
	normalMinerAccount = strings.ToLower(stub.GetNormalMiner())
	poolMinerAccount = strings.ToLower(stub.GetPoolMiner())
	goRatioVal, err := decimal.NewDecimalFromString(goRatio)
	if err != nil ||  goRatioVal.LessThanOrEqual(decimal.NewFromFloat(0)){
		return shim.Error("GoRatio format error")
	}
	var message = strings.ToLower(stub.GetTxID() + outAccount)
	var messageTransfer = hex.EncodeToString([]byte(message))
	var checkStatus = self.checkSign(outAccount, messageTransfer, args[4])
	if !checkStatus {
		return shim.Error("Attestation of failure")
	}
	if !self.checkAccountFormat(outAccount) {
		return shim.Error("Error format of outgoing account")
	}
	if !self.checkAccountFormat(superMinerAccount) {
		return shim.Error("Node s packing error")
	}
	if superMinerAccount == normalMinerAccount {
		return shim.Error("Node packing exception")
	}
	if normalMinerAccount != "" {
		if !self.checkAccountFormat(normalMinerAccount) {
			return shim.Error("Node n packing error")
		}
	}
	if self.checkAccountDisabled(outAccount) {
		return shim.Error("The account [" + outAccount + "] has been disabled")
	}
	if !self.checkRemarkFormat(remark) {
		return shim.Error("Remarks up to 32 characters [2 characters in Chinese]")
	}
	inAccountArray := strings.Split(inAccount, "|")
	outMoneyArray := strings.Split(outMoney, "|")
	agentPayArray := outMoneyArray
	outTokenArray := strings.Split(token, "|")
	tokenPrecisionMaps := make(map[string]int64)
	if len(inAccountArray) <= 1 || len(outMoneyArray) <= 1 {
		return shim.Error("The number of transactions has to be greater than 1")
	}
	if len(inAccountArray) != len(outMoneyArray) {
		return shim.Error("inAccount and money do not match")
	}
	if len(inAccountArray) != len(outTokenArray) {
		return shim.Error("inAccount and token do not match")
	}
	minerFeeVal, err = decimal.NewDecimalFromString(minerFee)
	if err != nil {
		return shim.Error("Limit amount conversion failed")
	}
	inAccountNumber := decimal.NewFromFloat(float64(len(inAccountArray)))
	avgMinerFeeVal := minerFeeVal.Div(inAccountNumber)
	chargeLowLimit,_ := self.getServiceCharge(goRatioVal)
	if avgMinerFeeVal.LessThan(chargeLowLimit) {
		return shim.Error("The miner's fee is too low")
	}
	realBlackMoney = decimal.NewFromFloat(0)
	superMinnerAward = decimal.NewFromFloat(0)
	normalMinnerAward = decimal.NewFromFloat(0)
	poolMinnerAward = decimal.NewFromFloat(0)
	accountBalanceMaps := make(map[string]decimal.Decimal)
	var agentPayTimes float64 = 0
	var precision int64
	for i := 0; i < len(inAccountArray); i++ {
		_inAccount := strings.ToLower(inAccountArray[i])
		_outToken := strings.ToLower(outTokenArray[i])
		if !self.checkAccountFormat(_inAccount) {
			return shim.Error("Wrong format to transfer to account")
		}
		if self.checkAccountDisabled(_inAccount) {
			return shim.Error("The account [" + _inAccount + "] has been disabled.")
		}
		if _inAccount == outAccount {
			return shim.Error("You can't trade with yourself")
		}
		accountBalanceMaps[ _outToken+"_"+_inAccount ], err = self.getAccountBalance(stub, _outToken, _inAccount);
		if err != nil {
			return shim.Error(err.Error())
		}
		_, checkExists = accountBalanceMaps[ self.mainToken.name+"_"+_inAccount ]
		if !checkExists {
			accountBalanceMaps[ self.mainToken.name+"_"+_inAccount ], err = self.getAccountBalance(stub, self.mainToken.name, _inAccount);
			if err != nil {
				return shim.Error(err.Error())
			}
		}
		agentPayAddr := self.getAgentPayAddr(_outToken)
		agentPayVal,err := self.getAccountBalance(stub, self.mainToken.name, agentPayAddr)
		if err != nil {
			return shim.Error(err.Error())
		}
		accountBalanceMaps[self.mainToken.name + "_"+agentPayAddr] = agentPayVal
	}
	for i := 0; i < len(outTokenArray); i++ {
		_outToken := strings.ToLower(outTokenArray[i])
		_, checkExists = accountBalanceMaps[ _outToken+"_"+outAccount ]
		if !checkExists {
			accountBalanceMaps[ _outToken+"_"+outAccount ], err = self.getAccountBalance(stub, _outToken, outAccount)
			if err != nil {
				return shim.Error(err.Error())
			}
		}
		precision, checkExists = tokenPrecisionMaps[ _outToken ]
		if !checkExists {
			precision, err := self.getTokenPrecision(stub,_outToken)
			if err != nil {
				return shim.Error(err.Error())
			}
			tokenPrecisionMaps[ _outToken ] = precision
		}
	}
	mainPrecision, err := self.getTokenPrecision(stub,self.mainToken.name)
	if err != nil {
		return shim.Error(err.Error())
	}
	tokenPrecisionMaps[ self.mainToken.name ] = mainPrecision
	_, checkExists = accountBalanceMaps[ self.mainToken.name+"_"+outAccount ]
	if !checkExists {
		accountBalanceMaps[ self.mainToken.name+"_"+outAccount ], err = self.getAccountBalance(stub, self.mainToken.name, outAccount)
		if err != nil {
			return shim.Error(err.Error())
		}
	}
	for i := 0; i < len(inAccountArray); i++ {
		_inAccount := strings.ToLower(inAccountArray[i])
		_token := strings.ToLower(outTokenArray[i])
		outVal, err := decimal.NewDecimalFromString(outMoneyArray[i])
		if err != nil {
			return shim.Error("Transaction quantity format error")
		}
		precision = tokenPrecisionMaps[ _token ]
		precisionLen := self.getPrecisionLen(outMoneyArray[i])
		if precisionLen > precision {
			return shim.Error("The number of transactions is too many decimal places")
		}
		if outVal.LessThanOrEqual(decimal.NewFromFloat(0)) {
			return shim.Error("The number of transactions must be greater than 0")
		}
		_inAccountBal := accountBalanceMaps[ _token+"_"+_inAccount ]
		_outAccountBal := accountBalanceMaps[ _token+"_"+outAccount ]
		accountBalanceMaps[ _token+"_"+_inAccount ] = _inAccountBal.Add(outVal)
		accountBalanceMaps[ _token+"_"+outAccount ] = _outAccountBal.Sub(outVal)
		if accountBalanceMaps[ _token+"_"+outAccount ].LessThan(zeroDecimal) {
			return shim.Error("The transfer account balance is insufficient")
		}
		agentPayAddr := self.getAgentPayAddr(_token)
		fmt.Println("token:",_token,",agentPayAccount:",agentPayAddr)
		agentPayKey := self.mainToken.name+"_" + agentPayAddr
		if accountBalanceMaps[agentPayKey].GreaterThanOrEqual(avgMinerFeeVal){
			accountBalanceMaps[agentPayKey] = accountBalanceMaps[agentPayKey].Sub(avgMinerFeeVal)
			agentPayArray[i] = agentPayAddr
			agentPayTimes++
		}else{
			agentPayArray[i] = ""
		}
	}
	realMinerFeeVal := minerFeeVal.Sub(avgMinerFeeVal.Mul(decimal.NewFromFloat(agentPayTimes)))
	accountBalanceMaps[ self.mainToken.name+"_"+outAccount ] = accountBalanceMaps[ self.mainToken.name+"_"+outAccount ].Sub(realMinerFeeVal)
	if accountBalanceMaps[ self.mainToken.name+"_"+outAccount ].LessThan(zeroDecimal) {
		return shim.Error("The transfer account balance is insufficient")
	}
	minnerMoneyCalData, err := self.minerMoneyCalculate(stub, superMinerAccount, normalMinerAccount,poolMinerAccount, minerFeeVal)
	if err != nil {
		return shim.Error(err.Error())
	}
	normalMinnerAward = minnerMoneyCalData["normalMinerAward"]
	superMinnerAward = minnerMoneyCalData["superMinerAward"]
	poolMinnerAward = minnerMoneyCalData["poolMinerAward"]
	realBlackMoney = minnerMoneyCalData["blackMoney"]
	for key, balance := range accountBalanceMaps {
		_accountStr := strings.Split(key, "_")
		_token := strings.ToLower(_accountStr[0])
		fmt.Println("token:",_token," precision:",tokenPrecisionMaps[ _token ])
		var err = stub.PutState(key, []byte( balance.StringFixed(int32(tokenPrecisionMaps[ _token ])) ))
		if err != nil {
			return shim.Error("Error saving the balance of transfer account")
		}
	}
	var agentPayAccountStr = strings.Join(agentPayArray,"|")
	var returnString = "token;batTransfer;" + token + "," + self.formatNumber(realBlackMoney, self.mainToken.precision) + "," + superMinerAccount + "," + normalMinerAccount + "," + self.formatNumber(minerFeeVal, self.mainToken.precision) + "," + self.formatNumber(superMinnerAward, self.mainToken.precision) + "," + self.formatNumber(normalMinnerAward, self.mainToken.precision) + "," + outAccount + "," + inAccount + "," + outMoney + "," + remark + "," + stub.GetTxID()+ "," + self.base64Encode(self.formatNumber(avgMinerFeeVal, self.mainToken.precision)) + "," + poolMinerAccount + "," + self.formatNumber(poolMinnerAward, self.mainToken.precision) + "," + agentPayAccountStr + "," + self.formatNumber(avgMinerFeeVal, self.mainToken.precision)
	return shim.Success([]byte(returnString))
}
func (self *TokenChaincode) makeTradingContractId(account string) string {
	nowTime := time.Now().In(self.dateZone).Unix()
	return strings.ToLower(self.md5string(account + strconv.FormatInt(nowTime, 10)))
}
func (self *TokenChaincode) publish(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	if len(args) != 12 {
		fmt.Println("publish£¨£© parameter require 12 £¨ [1] originatorAccount [2] name [3] title [4] gross [5] output [6] precision [7] site [8] logo [9] email [10] rate [11] sign [12] goRatio£©, current " + strconv.Itoa(len(args)))
		return shim.Error("Missing parameter")
	}
	if args[0] == "" || args[1] == "" || args[2] == "" || args[3] == "" || args[4] == "" || args[5] == "" || args[7] == "" || args[8] == "" || args[9] == "" || args[10] == ""|| args[11] == "" {
		return shim.Error("Parameters cannot be empty")
	}
	var token TokenConfig = TokenConfig{}
	token.originatorAccount = strings.ToLower(args[0])
	token.name = strings.ToLower(args[1])
	token.title = args[2]
	token.site = args[6]
	token.logo = args[7]
	token.email = args[8]
	token.gross, err = decimal.NewDecimalFromString(args[3])
	var precisionStr = args[5]
	var sign string = args[10]
	var goRatio = args[11]
	superMinerAccount := stub.GetSuperMiner()
	normalMinerAccount := stub.GetNormalMiner()
	poolMinerAccount := stub.GetPoolMiner()
	if err != nil {
		return shim.Error("Total asset format error")
	}
	if !self.checkAccountFormat(token.originatorAccount) {
		return shim.Error("Founder account format error")
	}
	if len(token.title) < 2 || len(token.title) > 64 {
		return shim.Error("The token name must be between 2 and 64 bytes long")
	}
	if !self.checkTokenFormat(token.name) {
		return shim.Error("Token name abbreviation format error[Case English number length between 2 and 20 bits]")
	}
	if token.name == self.mainToken.name {
		return shim.Error("Token existing")
	}
	if len(token.logo) >= 20480 {
		return shim.Error("The token image cannot be exceeded 10 k")
	}
	maxCurrencyValue, _ := decimal.NewDecimalFromString(self.maxCurrency)
	if token.gross.GreaterThan(maxCurrencyValue) {
		return shim.Error("Total assets cannot exceed the system's upper limit of one trillion yuan")
	}
	token.output, err = decimal.NewDecimalFromString(args[4])
	if err != nil {
		return shim.Error("An error occurred in the pre-dig conversion format")
	}
	if token.output.GreaterThan(token.gross) {
		return shim.Error("The amount of pre-dig shall not exceed the total amount of money")
	}
	token.precision, err = strconv.ParseInt(precisionStr, 10, 32)
	if token.precision > self.maxPrecision {
		return shim.Error("A maximum of 10 decimal places is supported")
	}
	if err != nil {
		return shim.Error("Token Decimal precision format error")
	}
	token.rate, err = decimal.NewDecimalFromString(args[9])
	if err != nil {
		return shim.Error("Error formatting annualized returns")
	}
	if token.rate.Equal(decimal.NewFromFloat(0)) {
		if token.output.LessThan(token.gross) {
			return shim.Error("The amount of pre-dig and the amount remaining do not match the total amount of money")
		}
	}
	token.dayAward = token.rate.Div(decimal.NewFromFloat(100)).Div(decimal.NewFromFloat(365))
	tokenConfigbytes, err := stub.GetState(token.name + self.configKey)
	if err != nil {
		return shim.Error("Check the token configuration for errors.")
	}
	if tokenConfigbytes != nil {
		return shim.Error("Token existing")
	}
	if self.checkAccountDisabled(token.originatorAccount) {
		return shim.Error("The account [" + token.originatorAccount + "] has been disabled")
	}
	goRatioVal, err := decimal.NewDecimalFromString(goRatio)
	if err != nil ||  goRatioVal.LessThanOrEqual(decimal.NewFromFloat(0)){
		return shim.Error("GoRatio format error")
	}
	messageBody := hex.EncodeToString([]byte(strings.ToLower(stub.GetTxID() + token.name + token.originatorAccount)))
	var checkStatus = self.checkSign(token.originatorAccount, messageBody, sign)
	if !checkStatus {
		return shim.Error("Attestation of failure")
	}
	outAccountMainVal, err := self.getAccountBalance(stub, self.mainToken.name, token.originatorAccount)
	if err != nil {
		fmt.Println("error:", err.Error())
		return shim.Error(err.Error())
	}
	normalMinerAwardPayload := "0"
	superMinerAwardPayload :="0"
	poolMinerAwardPayload :="0"
	blackMoneyPayload :="0"
	tokenChargePayload :=""
	toAcc:=""
	tokenMortAmount:="0"
	_,tokenCharge := self.getServiceCharge(goRatioVal)
	minerMoneyCalData, err := self.minerMoneyCalculate(stub,superMinerAccount, normalMinerAccount,poolMinerAccount, tokenCharge)
	if err != nil {
		return shim.Error(err.Error())
	}
	outAccountMainVal = outAccountMainVal.Sub(tokenCharge)
	err = self.updateAccountBalance(stub, self.mainToken.name, token.originatorAccount,outAccountMainVal, self.mainToken.precision)
	if err != nil {
		return shim.Error("Error saving the balance of transfer account")
	}
	poolMinerAwardPayload = self.formatNumber(minerMoneyCalData["poolMinerAward"], self.mainToken.precision)
	normalMinerAwardPayload = self.formatNumber(minerMoneyCalData["normalMinerAward"], self.mainToken.precision)
	superMinerAwardPayload = self.formatNumber(minerMoneyCalData["superMinerAward"], self.mainToken.precision)
	blackMoneyPayload = self.formatNumber(minerMoneyCalData["blackMoney"], self.mainToken.precision)
	tokenChargePayload = self.formatNumber(self.tokenCharge, self.mainToken.precision)
	fmt.Println("normalMinerAward", normalMinerAwardPayload)
	fmt.Println("superMinerAward", superMinerAwardPayload)
	err = self.updateAccountBalance(stub, token.name, token.originatorAccount, token.output, token.precision)
	if err != nil {
		return shim.Error("Error 1 writing token data")
	}
	token.mineral = token.gross.Sub(token.output)
	err = self.updateAccountBalance(stub, token.name, self.mineralKey, token.mineral, token.precision)
	if err != nil {
		return shim.Error("Error 2 writing token data")
	}
	err = stub.PutState(token.name+self.grossKey, []byte( token.gross.String() ))
	if err != nil {
		return shim.Error("Error 3 writing token data")
	}
	err = stub.PutState(token.name+self.outputKey, []byte( token.output.String() ))
	if err != nil {
		return shim.Error("Error 4 writing token data")
	}
	err = stub.PutState(token.name+self.rateKey, []byte( token.rate.String() ))
	if err != nil {
		return shim.Error("Error 5 writing token data")
	}
	err = stub.PutState(token.name+self.dayRateKey, []byte( token.dayAward.String() ))
	if err != nil {
		return shim.Error("Error 6 writing token data")
	}
	err = stub.PutState(token.name+self.precisionKey, []byte( precisionStr ))
	if err != nil {
		return shim.Error("Error 7 writing token data")
	}
	awardKey := strings.ToLower(token.name + self.tokenAwardKey + token.originatorAccount)
	var stateString = "0,0," + self.getCurrentDateString()
	err = stub.PutState(awardKey, []byte( stateString ));
	if err != nil {
		return shim.Error("Error saving account time")
	}
	var config [13] string
	config[0] = token.originatorAccount
	config[1] = token.name
	config[2] = token.title
	config[3] = token.gross.String()
	config[4] = token.output.String()
	config[5] = strconv.FormatInt(token.precision, 10)
	config[6] = token.site
	config[7] = token.logo
	config[8] = token.email
	config[9] = token.rate.String()
	config[10] = self.formatNumber(self.tokenCharge, token.precision)
	config[11] = "0"
	config[12] = self.getCurrentDateString()
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	err = enc.Encode(config)
	if err != nil {
		return shim.Error("The token information encoding failed")
	}
	err = stub.PutState(token.name+self.configKey, network.Bytes())
	if err != nil {
		return shim.Error("The token information failed to be saved to the state")
	}
	var returnString = "token;publish;" + token.name + "," + blackMoneyPayload + "," + superMinerAccount + "," + normalMinerAccount + "," + tokenChargePayload + "," + superMinerAwardPayload + "," + normalMinerAwardPayload + "," + token.originatorAccount + ","+toAcc+","+tokenMortAmount+"," + self.base64Encode(token.name) + "," + stub.GetTxID()+","+poolMinerAccount+","+poolMinerAwardPayload
	return shim.Success([]byte(returnString))
}
func (self *TokenChaincode) minerMoneyCalculate(stub shim.ChaincodeStubInterface, superMinerAccount, normalMinerAccount, poolMinnerAccount string, minerFeeVal decimal.Decimal) (map[string]decimal.Decimal, error) {
	var err error
	var blackAccountVal decimal.Decimal
	var returnMap = make(map[string]decimal.Decimal)
	returnMap["blackMoney"] = decimal.NewFromFloat(0)
	returnMap["normalMinerAward"] = decimal.NewFromFloat(0)
	returnMap["superMinerAward"] = decimal.NewFromFloat(0)
	returnMap["poolMinerAward"] = decimal.NewFromFloat(0)
	var chargeDestroy = minerFeeVal.Div(decimal.NewFromFloat(10))
	var realBlackMoney = decimal.NewFromFloat(0)
	var superMinnerAward decimal.Decimal
	var normalMinnerAward decimal.Decimal = decimal.NewFromFloat(0)
	var poolMinerAward decimal.Decimal = decimal.NewFromFloat(0)
	/************************** ÊÖÐø·Ñ×ÔÏú»Ù¹¦ÄÜ ¿ªÊ¼  **************************/
	blackAccountVal, err = self.getAccountBalance(stub, self.mainToken.name, self.blackKey)
	fmt.Println("blackAccountVal:", blackAccountVal.String())
	if err != nil {
		fmt.Println("error:", err.Error())
		return returnMap, errors.New(err.Error())
	}
	blackAccountVal = blackAccountVal.Add(chargeDestroy)
	superminerChargeAwardAll := minerFeeVal.Mul(decimal.RequireFromString("0.4"))
	normalminerChargeAwardAll := minerFeeVal.Mul(decimal.RequireFromString("0.5"))
	if blackAccountVal.LessThan(self.chargeDestroyUpperLimit) {
		realBlackMoney = chargeDestroy
		if normalMinerAccount == "" {
			superMinnerAward = minerFeeVal.Sub(chargeDestroy)
		} else {
			if poolMinnerAccount == ""{
				normalMinnerAward = normalminerChargeAwardAll
			}else{
				normalMinnerAward = normalminerChargeAwardAll.Mul(decimal.RequireFromString("0.9"))
				poolMinerAward = normalminerChargeAwardAll.Mul(decimal.RequireFromString("0.1"))
			}
			superMinnerAward = superminerChargeAwardAll
		}
	} else {
		if normalMinerAccount == "" {
			superMinnerAward = minerFeeVal
		} else {
			if poolMinnerAccount == ""{
				normalMinnerAward = normalminerChargeAwardAll
			}else{
				normalMinnerAward = normalminerChargeAwardAll.Mul(decimal.RequireFromString("0.9"))
				poolMinerAward = normalminerChargeAwardAll.Mul(decimal.RequireFromString("0.1"))
			}
			superMinnerAward = superminerChargeAwardAll
		}
	}
	/************************** ÊÖÐø·Ñ×ÔÏú»Ù¹¦ÄÜ ½áÊø  **************************/
	returnMap["blackMoney"] = realBlackMoney
	returnMap["normalMinerAward"] = normalMinnerAward
	returnMap["superMinerAward"] = superMinnerAward
	returnMap["poolMinerAward"] = poolMinerAward
	return returnMap, nil
}
func (self *TokenChaincode) accountMoneyCalculate(stub shim.ChaincodeStubInterface, method int, precision int64, normalMinerAccount, superMinerAccount, poolMinerAccount, outAccount, inAccount, tokenId string, outVal, minerFeeVal,goRatio decimal.Decimal) (map[string]decimal.Decimal, error) {
	var err error
	var outAccountMainVal, inAccountMainVal, outAccountTokenVal, inAccountTokenVal decimal.Decimal
	var returnMap = make(map[string]decimal.Decimal)
	returnMap["poolMinerAward"] = decimal.NewFromFloat(0)
	returnMap["normalMinerAward"] = decimal.NewFromFloat(0)
	returnMap["superMinerAward"] = decimal.NewFromFloat(0)
	returnMap["blackMoney"] = decimal.NewFromFloat(0)
	returnMap["agentPayMoney"] = decimal.NewFromFloat(0)
	tokenId = strings.ToLower(tokenId)
	if outVal.LessThanOrEqual(decimal.NewFromFloat(0)) {
		return returnMap, errors.New("The number of transactions must be greater than 0")
	}
	chargeLowLimit,_ := self.getServiceCharge(goRatio)
	if minerFeeVal.LessThan(chargeLowLimit) {
		fmt.Println("minerFeeVal:",minerFeeVal)
		fmt.Println("chargeLowLimit:",chargeLowLimit)
		return returnMap, errors.New("The miner's fee is too low")
	}
	outAccountMainVal, err = self.getAccountBalance(stub, self.mainToken.name, outAccount)
	if err != nil {
		fmt.Println("error:", err.Error())
		return returnMap, errors.New(err.Error())
	}
	inAccountMainVal, err = self.getAccountBalance(stub, self.mainToken.name, inAccount)
	if err != nil {
		fmt.Println("error:", err.Error())
		return returnMap, errors.New(err.Error())
	}
	inAccountTokenVal, err = self.getAccountBalance(stub, tokenId, inAccount)
	if err != nil {
		fmt.Println("error:", err.Error())
		return returnMap, errors.New(err.Error())
	}
	if tokenId == self.mainToken.name {
		if outAccountMainVal.LessThan(outVal.Add(minerFeeVal)) {
			return returnMap, errors.New("The transfer account balance is insufficient")
		}
	} else {
		outAccountTokenVal, err = self.getAccountBalance(stub, tokenId, outAccount)
		if err != nil {
			fmt.Println("error:", err.Error())
			return returnMap, errors.New(err.Error())
		}
		if outAccountTokenVal.LessThan(outVal) {
			return returnMap, errors.New("The transfer account balance is insufficient")
		}
		if outAccountMainVal.LessThan(minerFeeVal) {
			return returnMap, errors.New("The transfer account balance is insufficient")
		}
	}
	minnerMoneyCalData, err := self.minerMoneyCalculate(stub, superMinerAccount, normalMinerAccount,poolMinerAccount, minerFeeVal)
	if err != nil {
		return returnMap, err
	}
	returnMap["poolMinerAward"] = minnerMoneyCalData["poolMinerAward"]
	returnMap["normalMinerAward"] = minnerMoneyCalData["normalMinerAward"]
	returnMap["superMinerAward"] = minnerMoneyCalData["superMinerAward"]
	returnMap["blackMoney"] = minnerMoneyCalData["blackMoney"]
	returnMap["superMinerAward"] = minnerMoneyCalData["superMinerAward"]
	fmt.Println("normalMinerAward", returnMap["normalMinerAward"])
	fmt.Println("superMinerAward", returnMap["superMinerAward"])
	outAccountMainVal = outAccountMainVal.Sub(minerFeeVal)
	if tokenId == self.mainToken.name {
		inAccountMainVal = inAccountMainVal.Add(outVal)
		outAccountMainVal = outAccountMainVal.Sub(outVal)
		err = self.updateAccountBalance(stub, self.mainToken.name, outAccount, outAccountMainVal, self.mainToken.precision)
		if err != nil {
			return returnMap, errors.New("Error saving the balance of transfer account")
		}
		err = self.updateAccountBalance(stub, self.mainToken.name, inAccount, inAccountMainVal, self.mainToken.precision)
		if err != nil {
			return returnMap, errors.New("Error saving into account balance")
		}
	} else {
		agentPayAddr := self.getAgentPayAddr(tokenId)
		agentPayVal, err := self.getAccountBalance(stub, self.mainToken.name, agentPayAddr)
		if err != nil {
			return  returnMap,err
		}
		if agentPayVal.GreaterThanOrEqual(minerFeeVal){
			agentPayVal = agentPayVal.Sub(minerFeeVal)
			outAccountMainVal = outAccountMainVal.Add(minerFeeVal)
			err = self.updateAccountBalance(stub, self.mainToken.name, agentPayAddr, agentPayVal, self.mainToken.precision)
			if err != nil {
				return returnMap, errors.New("Error saving the balance of transfer account")
			}
			returnMap["agentPayMoney"] = minerFeeVal
		}
		err = self.updateAccountBalance(stub, self.mainToken.name, outAccount, outAccountMainVal, self.mainToken.precision)
		if err != nil {
			return returnMap, errors.New("Error saving the balance of transfer account")
		}
		outAccountTokenVal, err = self.getAccountBalance(stub, tokenId, outAccount)
		if err != nil {
			fmt.Println("error:", err.Error())
			return returnMap, errors.New(err.Error())
		}
		inAccountTokenVal, err = self.getAccountBalance(stub, tokenId, inAccount)
		if err != nil {
			fmt.Println("error:", err.Error())
			return returnMap, errors.New(err.Error())
		}
		inAccountTokenVal = inAccountTokenVal.Add(outVal)
		outAccountTokenVal = outAccountTokenVal.Sub(outVal)
		err = self.updateAccountBalance(stub, tokenId, outAccount, outAccountTokenVal, precision)
		if err != nil {
			return returnMap, errors.New("Error saving into account balance")
		}
		err = self.updateAccountBalance(stub, tokenId, inAccount, inAccountTokenVal, precision)
		if err != nil {
			return returnMap, errors.New("Error saving into account balance")
		}
	}
	return returnMap, nil
}
func (self *TokenChaincode) transfer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var token, outAccount, inAccount, minerFee, remark,goRatio string
	var outVal decimal.Decimal
	var err error
	if len(args) != 8 {
		fmt.Println("transfer£¨£© parameter require 8 £¨ [1] token [2] inAcount [3] outAccount [4] number [5] sign [6] minerFee [7] remark [8] goRatio£©")
		return shim.Error("Missing parameter")
	}
	if args[0] == "" || args[1] == "" || args[2] == "" || args[3] == "" || args[4] == "" || args[5] == "" || args[7] == "" {
		return shim.Error("Parameters cannot be empty")
	}
	superMinerAccount := stub.GetSuperMiner()
	normalMinerAccount := stub.GetNormalMiner()
	poolMinerAccount := stub.GetPoolMiner()
	fmt.Println("super:",superMinerAccount," normal:",normalMinerAccount," pool:",poolMinerAccount)
	token = strings.ToLower(args[0])
	inAccount = strings.ToLower(args[1])
	outAccount = strings.ToLower(args[2])
	remark = args[6]
	minerFee = args[5]
	goRatio = args[7]
	if !self.checkAccountFormat(inAccount) {
		return shim.Error("Wrong format to transfer to account")
	}
	/*if !self.checkAccountFormat(outAccount) {
		return shim.Error("Error format of outgoing account")
	}*/
	if inAccount == outAccount {
		return shim.Error("You can't trade with yourself")
	}
	if self.checkAccountDisabled(outAccount) {
		return shim.Error("The account [" + outAccount + "] has been disabled")
	}
	if !self.checkRemarkFormat(remark) {
		return shim.Error("Remarks up to 32 characters [2 characters in Chinese]")
	}
	if !self.checkTokenFormat(token) {
		return shim.Error("Token name abbreviation format error[Case English number length between 2 and 20 bits].")
	}
	minerFeeVal, err := decimal.NewDecimalFromString(minerFee)
	if err != nil {
		return shim.Error("Miner fee format error")
	}
	outVal, err = decimal.NewDecimalFromString(args[3])
	if err != nil {
		return shim.Error("Transaction quantity format error")
	}
	if outVal.LessThan(decimal.NewFromFloat(0)) {
		return shim.Error("The number of transactions must be greater than 0")
	}
	goRatioVal, err := decimal.NewDecimalFromString(goRatio)
	if err != nil ||  goRatioVal.LessThanOrEqual(decimal.NewFromFloat(0)){
		return shim.Error("GoRatio format error")
	}
	var signHex string = args[4]
	var message = strings.ToLower(stub.GetTxID() + token + outAccount + inAccount + args[3])
	var messageTransfer = hex.EncodeToString([]byte(message))
	var checkStatus = self.checkSign(outAccount, messageTransfer, signHex)
	if !checkStatus {
		return shim.Error("Attestation of failure")
	}
	precision, err := self.getTokenPrecision(stub,token)
	if err != nil {
		return shim.Error(err.Error())
	}
	precisionLen := self.getPrecisionLen(args[3])
	if precisionLen > precision {
		return shim.Error("The number of transactions is too many decimal places")
	}
	/********************* »ñÈ¡×ªÈëÕËºÅµÄÓà¶î **************/
	key := strings.ToLower(token + "_" + inAccount)
	keyValbytes, err := stub.GetState(strings.ToLower(key))
	if err != nil {
		return shim.Error("Error getting transfer balance")
	}
	if keyValbytes == nil {
		awardKey := strings.ToLower(token + self.tokenAwardKey + inAccount)
		nowDate := self.getCurrentDateString()
		var stateString = "0,0," + nowDate
		err = stub.PutState(awardKey, []byte( stateString ))
		if err != nil {
			return shim.Error("Error saving transfer time to account")
		}
	}
	/****************************************************/
	calData, err := self.accountMoneyCalculate(stub, 0, precision, normalMinerAccount, superMinerAccount, poolMinerAccount,outAccount, inAccount, token, outVal, minerFeeVal,goRatioVal)
	if err != nil {
		return shim.Error(err.Error())
	}
	poolMinerAward := calData["poolMinerAward"]
	agentPayAccount := ""
	if calData["agentPayMoney"].GreaterThan(decimal.NewFromFloat(0)){
		agentPayAccount = self.getAgentPayAddr(token)
	}
	returnString := "token;transfer;" + token + "," + self.formatNumber(calData["blackMoney"], precision) + "," + superMinerAccount + "," + normalMinerAccount + "," + self.formatNumber(minerFeeVal, precision) + "," + self.formatNumber(calData["superMinerAward"], precision) + "," + self.formatNumber(calData["normalMinerAward"], precision) + "," + outAccount + "," + inAccount + "," + self.formatNumber(outVal, precision) + "," + remark + "," + stub.GetTxID()+"," + poolMinerAccount + "," + self.formatNumber(poolMinerAward, precision) + "," + agentPayAccount
	return shim.Success([]byte(returnString))
}
func (self *TokenChaincode) balance(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		fmt.Println("balance£¨£© parameter require 2 £¨ [1] token [2] account £©")
		return shim.Error("Missing parameter")
	}
	if args[0] == "" || args[1] == "" {
		return shim.Error("Parameters cannot be empty")
	}
	token := strings.ToLower(args[0])
	account := strings.ToLower(args[1])
	if !self.checkTokenFormat(token) {
		return shim.Error("Token format error")
	}
	if !self.checkAccountFormat(account) {
		return shim.Error("Account format error.")
	}
	accountVal, err := self.getAccountBalance(stub, token, account)
	if err != nil {
		fmt.Println("error:", err.Error())
		return shim.Error(err.Error())
	}
	precision, err := self.getTokenPrecision(stub,token)
	if err != nil {
		return shim.Error(err.Error())
	}
	formatedNumber := self.formatNumber(accountVal, precision)
	return shim.Success([]byte(formatedNumber))
}
func (self *TokenChaincode) awardPreview(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var money decimal.Decimal
	if len(args) != 2 {
		fmt.Println("awardPreview£¨£© args require 2 £¨ [1] token [2] account £©")
		return shim.Error("Missing parameter")
	}
	if args[0] == "" || args[1] == "" {
		return shim.Error("Parameters cannot be empty")
	}
	var token = strings.ToLower(args[0])
	var account = strings.ToLower(args[1])
	money, _, err := self.doAwardPreview(stub, account, token)
	if err != nil {
		return shim.Error(err.Error())
	}
	precision, err := self.getTokenPrecision(stub,token)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte( self.formatNumber(money, precision) ))
}
type AccountBalance struct {
	money string
	time  int64
}
type AccountBalanceList []*AccountBalance
type BalanceReport struct {
	money     decimal.Decimal
	time      decimal.Decimal
	beginDate string
	endDate   string
}
func (self *TokenChaincode) countMoneyReport(list AccountBalanceList, endTime int64, beginTime int64, hourRate decimal.Decimal, isDebug bool) []BalanceReport {
	var lastBalance *AccountBalance
	var reportList []BalanceReport
	if isDebug {
		beginDate := time.Unix(beginTime, 0)
		endDate := time.Unix(endTime, 0)
		fmt.Println("beginDate:", beginDate.In(self.dateZone).Format("2006-01-02 15:04:05"), "endDate:", endDate.In(self.dateZone).Format("2006-01-02 15:04:05"))
		fmt.Println("-------------------------------------")
	}
	for j := len(list); j >= 1; j-- {
		var reporttime int64
		var _beinTime int64
		var _endTime int64
		balance := list[j-1]
		if isDebug {
			date := time.Unix(balance.time, 0)
			fmt.Println(" balance.time:", date.In(self.dateZone).Format("2006-01-02 15:04:05"))
		}
		if balance.time > beginTime && balance.time < endTime {
			report := BalanceReport{}
			report.money, _ = decimal.NewDecimalFromString(balance.money)
			if lastBalance != nil {
				if lastBalance.time > endTime {
					_beinTime = balance.time
					_endTime = endTime
				} else {
					_beinTime = balance.time
					_endTime = lastBalance.time
				}
			} else {
				_beinTime = balance.time
				_endTime = endTime
			}
			if _endTime == _beinTime {
				continue
			}
			reporttime = _endTime - _beinTime
			report.beginDate = time.Unix(_beinTime, 0).In(self.dateZone).Format("01/02 15:04:05")
			report.endDate = time.Unix(_endTime, 0).In(self.dateZone).Format("01/02 15:04:05")
			report.time, _ = decimal.NewDecimalFromString(strconv.FormatInt(reporttime, 10))
			reportList = append(reportList, report)
		} else if (balance.time < beginTime && lastBalance != nil && lastBalance.time > beginTime) ||
			(lastBalance == nil && balance.time < beginTime) ||
			(balance.time < endTime && (beginTime-balance.time) < 86400 && (balance.time-endTime) < 86400) {
			report := BalanceReport{}
			report.money, _ = decimal.NewDecimalFromString(balance.money)
			if lastBalance != nil {
				if lastBalance.time < beginTime {
					continue
				}
				if lastBalance.time > endTime {
					_beinTime = beginTime
					_endTime = endTime
				} else {
					_beinTime = beginTime
					_endTime = lastBalance.time
				}
			} else {
				_beinTime = beginTime
				_endTime = endTime
			}
			if _endTime == _beinTime {
				continue
			}
			reporttime = _endTime - _beinTime
			report.beginDate = time.Unix(_beinTime, 0).In(self.dateZone).Format("01/02 15:04:05")
			report.endDate = time.Unix(_endTime, 0).In(self.dateZone).Format("01/02 15:04:05")
			report.time, _ = decimal.NewDecimalFromString(strconv.FormatInt(reporttime, 10))
			reportList = append(reportList, report)
		}
		lastBalance = balance
	}
	return reportList
}
func (self *TokenChaincode) sumReportMoney(reportList []BalanceReport, hourRate decimal.Decimal, isDebug bool) decimal.Decimal {
	var hourDecimal = decimal.NewFromFloat(3600)
	var sumAwardMoney decimal.Decimal = decimal.NewFromFloat(0)
	for i := 0; i < len(reportList); i++ {
		report := reportList[i]
		countHour := report.time.Div(hourDecimal)
		awardMoney := countHour.Mul(hourRate).Mul(report.money)
		sumAwardMoney = sumAwardMoney.Add(awardMoney)
	}
	return sumAwardMoney
}
func (self *TokenChaincode) doAwardPreview(stub shim.ChaincodeStubInterface, account, token string) (decimal.Decimal, string, error) {
	var sumMoney = decimal.NewFromFloat(0)
	var nowDate = self.getCurrentDateString()
	var endTime int64 = 0
	var beginTime int64 = 0
	awardKey := strings.ToLower(token + self.tokenAwardKey + account)
	dayRate, err := self.getStateDecimalValue(stub, token+self.dayRateKey)
	if err != nil {
		return sumMoney, nowDate, errors.New(err.Error())
	}
	if dayRate.Equal(decimal.NewFromFloat(0)) {
		return sumMoney, nowDate, nil
	}
	var hourRate decimal.Decimal = dayRate.Div(decimal.NewFromFloat(24))
	accountAwardValbytes, err := stub.GetState(awardKey)
	if err != nil {
		return sumMoney, nowDate, errors.New("Last earnings call error")
	}
	if accountAwardValbytes == nil {
		return sumMoney, nowDate, errors.New("This token is not obtained")
	}
	resultArray := strings.Split(string(accountAwardValbytes), ",")
	beginTime, err = self.getUnixtimeFromString(string(resultArray[2]));
	if err != nil {
		return sumMoney, nowDate, errors.New("Last earnings time conversion error")
	}
	nowZeroDatestr := time.Now().In(self.dateZone).Format("2006-01-02")
	nowZeroTimestamp, err := time.ParseInLocation("2006-01-02", nowZeroDatestr, self.dateZone)
	if err != nil {
		return sumMoney, nowDate, errors.New("Current time conversion error")
	}
	endTime = nowZeroTimestamp.Unix()
	fmt.Println(" lastDate:", string(accountAwardValbytes), " beginTime:", beginTime, " endTime:", endTime)
	if beginTime == endTime {
		return sumMoney, nowDate, nil
	}
	moneyHistory, err := self.getAccountMoneyHistory(stub, account, token, false)
	sumMoney = self.sumReportMoney(self.countMoneyReport(moneyHistory, endTime, beginTime, hourRate, false), hourRate, false)
	return sumMoney, nowDate, nil
}
func (self *TokenChaincode) award(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		fmt.Println("award() args require 3 £¨ [1] token [2] acount [3] sign £©")
		return shim.Error("Missing parameter")
	}
	if args[0] == "" || args[1] == "" || args[2] == "" {
		return shim.Error("Parameters cannot be empty")
	}
	var token = strings.ToLower(args[0])
	var account = strings.ToLower(args[1])
	var signHex string = args[2]
	var returnString string
	if !self.checkTokenFormat(token) {
		return shim.Error("Token Name abbreviation format error.")
	}
	if !self.checkAccountFormat(account) {
		return shim.Error("Account format error.")
	}
	if self.checkAccountDisabled(account) {
		return shim.Error("The account [" + account + "] has been disabled")
	}
	var message = strings.ToLower(stub.GetTxID() + token + account)
	var messageTransfer = hex.EncodeToString([]byte(message))
	var checkStatus = self.checkSign(account, messageTransfer, signHex)
	if !checkStatus {
		return shim.Error("Attestation of failure")
	}
	precision, err := self.getTokenPrecision(stub,token)
	if err != nil {
		return shim.Error(err.Error())
	}
	accountVal, err := self.getAccountBalance(stub, token, account)
	if err != nil {
		return shim.Error(err.Error())
	}
	money, nowDate, err := self.doAwardPreview(stub, account, token)
	if err != nil {
		return shim.Error(err.Error())
	}
	if money.LessThanOrEqual(decimal.NewFromFloat(0)) {
		return shim.Error("There is no revenue at the moment")
	}
	mineralVal, err := self.getAccountBalance(stub, token, self.mineralKey)
	if err != nil {
		return shim.Error("Error in obtaining mine allowance")
	}
	accountVal = accountVal.Add(money)
	if mineralVal.LessThan(money) {
		return shim.Error("The mine is short of surplus")
	}
	awardKey := strings.ToLower(token + self.tokenAwardKey + account)
	var stateString = "1," + money.StringFixed(10) + "," + nowDate
	err = stub.PutState(awardKey, []byte( stateString ));
	if err != nil {
		return shim.Error("The data of saving mining income is wrong")
	}
	err = self.updateAccountBalance(stub, token, account, accountVal, precision)
	if err != nil {
		return shim.Error(err.Error())
	}
	mineralVal = mineralVal.Sub(money)
	err = self.updateAccountBalance(stub, token, self.mineralKey, mineralVal, precision)
	if err != nil {
		return shim.Error(err.Error())
	}
	returnString = "token;award;" + token + ",0," + stub.GetSuperMiner() + ",,0,0,0," + self.mineralKey + "," + account + "," + self.formatNumber(money, precision) + ",pos," + stub.GetTxID()
	return shim.Success([]byte(returnString))
}
func (self *TokenChaincode) awardList(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var content string
	if args == nil || len(args) != 2 {
		fmt.Println("awardList£¨£© parameter require 2 £¨ [1] token [2] account £©")
		return shim.Error("Missing parameter")
	}
	if args[0] == "" || args[1] == "" {
		return shim.Error("Parameters cannot be empty")
	}
	content = ""
	token := strings.ToLower(args[0])
	account := strings.ToLower(args[1])
	awardKey := strings.ToLower(token + self.tokenAwardKey + account)
	historyQuery, err := stub.GetHistoryForKey(awardKey)
	if err != nil {
		return shim.Error("Error getting list data")
	}
	for historyQuery.HasNext() {
		keyModification, err := historyQuery.Next();
		if err != nil {
			return shim.Error(err.Error())
		}
		value := string(keyModification.GetValue())
		resultArray := strings.Split(value, ",")
		if len(resultArray) == 3 && resultArray[0] == "1" {
			if content != "" {
				content += "|"
			}
			content += resultArray[1] + "," + resultArray[2]
		}
	}
	return shim.Success([]byte(content))
}
func (self *TokenChaincode) saveMortgPutLog(stub shim.ChaincodeStubInterface, account, mortgageId, detailId, money string) error {
	key := mortgageId + "_data"
	reldata := account + "," + detailId + "," + money
	err := stub.PutState(key, []byte(reldata))
	if err != nil {
		return errors.New("mortgage data saving error")
	}
	return nil
}
func (self *TokenChaincode) getCurrentDateString() string {
	return time.Now().In(self.dateZone).Format("2006-01-02")
}
func (self *TokenChaincode) getUnixtimeFromString(dateStr string) (int64, error) {
	dateTime, err := time.ParseInLocation("2006-01-02", dateStr, self.dateZone)
	return dateTime.Unix(), err
}
func (self *TokenChaincode) getMineral(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		fmt.Println("mineral£¨£© args require 1 £¨ [1] token£©")
		return shim.Error("Missing parameter")
	}
	if args[0] == "" {
		return shim.Error("Parameters cannot be empty")
	}
	var token = strings.ToLower(args[0])
	mineralVal, err := self.getAccountBalance(stub, token, self.mineralKey)
	if err != nil {
		fmt.Println("error:", err.Error())
		return shim.Error("error: " + err.Error())
	}
	return shim.Success([]byte(mineralVal.String()))
}
func (self *TokenChaincode) info(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		fmt.Println("args require 1 £¨ [1] token£©")
		return shim.Error("Missing parameter")
	}
	if args[0] == "" {
		return shim.Error("Parameters cannot be empty")
	}
	var token = strings.ToLower(args[0])
	if !self.checkTokenFormat(token) {
		return shim.Error("Token Name abbreviation format error")
	}
	configbytes, err := stub.GetState(token + self.configKey)
	if err != nil {
		return shim.Error("Get token data state error.")
	}
	if configbytes == nil {
		return shim.Error("The token data obtained is empty")
	}
	return shim.Success(configbytes)
}
func (self *TokenChaincode) mortgageGetView(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if args == nil || len(args) != 2 {
		fmt.Println("mortgageGetView£¨£© parameter require 2 £¨ [1] token [2] mortgageId £©")
		return shim.Error("Missing parameter")
	}
	if args[0] == "" || args[1] == "" {
		for i := 0; i < len(args); i++ {
			fmt.Println("args[", strconv.Itoa(i), "] = ", args[i])
		}
		return shim.Error("Parameters cannot be empty")
	}
	var token, account, mortgageId string
	var money decimal.Decimal
	var err error
	token = strings.ToLower(args[0])
	mortgageId = strings.ToLower(args[1])
	mortgageAddr := self.mortgPrefix + mortgageId
	mortgageKey := token + "_" + mortgageAddr
	mortgageValbytes, err := stub.GetState(mortgageKey)
	if err != nil {
		return shim.Error("Failed to obtain mortgage data")
	}
	if mortgageValbytes == nil {
		return shim.Error("Mortgage data is empty")
	}
	key := mortgageId + "_data"
	historyQuery, err := stub.GetHistoryForKey(key)
	if err != nil {
		return shim.Error(err.Error())
	}
	accountMoney := make(map[string]decimal.Decimal)
	for historyQuery.HasNext() {
		keyModification, err := historyQuery.Next();
		if err != nil {
			return shim.Error(err.Error())
		}
		reldata := string(keyModification.GetValue())
		resultArray := strings.Split(reldata, ",")
		if len(resultArray) == 3 {
			account = resultArray[0]
			money, _ = decimal.NewDecimalFromString(resultArray[2])
			_, isExist := accountMoney[account]
			if !isExist {
				accountMoney[account] = money
			} else {
				accountMoney[account] = money.Add(accountMoney[account])
			}
		}
	}
	var invokeData = ""
	if len(accountMoney) == 0 {
		return shim.Error("No amount is redeemable")
	}
	var keys []string
	for _addr, _ := range accountMoney {
		keys = append(keys, _addr)
	}
	sort.Strings(keys)
	for _, _addr := range keys {
		_money := accountMoney[_addr]
		fmt.Println("%s->%s\n", _addr, _money)
		if invokeData != "" {
			invokeData += ","
		}
		invokeData += _addr + ":" + _money.String()
	}
	return shim.Success([]byte(invokeData))
}
func (self *TokenChaincode) mortgagePut(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if args == nil || len(args) != 9 {
		fmt.Println("mortgagePut£¨£© parameter require 9 £¨ [1] token [2] account [3] money [4] mortgageId [5] detailId [6] sign [7] method [8] mortgRate  [9] mortgDay £©")
		return shim.Error("Missing parameter")
	}
	if args[0] == "" || args[1] == "" || args[2] == "" || args[3] == "" || args[4] == "" || args[5] == "" || args[6] == "" {
		for i := 0; i < len(args); i++ {
			fmt.Println("args[", strconv.Itoa(i), "] = ", args[i])
		}
		return shim.Error("Parameters cannot be empty")
	}
	var token, account, sign, method, mortgageId, mortgageDetailId,mortgRate,mortgDays string
	var accountVal, mortgageAccountVal, mortgageMoney decimal.Decimal
	var err error
	token = strings.ToLower(args[0])
	account = strings.ToLower(args[1])
	mortgageMoney, err = decimal.NewDecimalFromString(args[2])
	mortgageId = strings.ToLower(args[3])
	mortgageDetailId = strings.ToLower(args[4])
	sign = args[5]
	method = args[6]
	mortgRate = args[7]
	mortgDays = args[8]
	if err != nil {
		return shim.Error("Incorrect mortgage amount format")
	}
	if mortgRate != "" && mortgDays != ""{
		mortgRateDecimal,err:=decimal.NewDecimalFromString(mortgRate)
		if err != nil {
			return shim.Error("Incorrect mortgRate format")
		}
		if mortgRateDecimal.GreaterThan(decimal.RequireFromString("1")) || mortgRateDecimal.LessThanOrEqual(decimal.RequireFromString("0"))  {
			return shim.Error("Incorrect mortgRate format")
		}
		_,err = decimal.NewDecimalFromString(mortgDays)
		if err != nil{
			return shim.Error("Incorrect mortgDays format")
		}
		mortgRedeemKey:= "mortgRedeemKey_"+mortgageId
		err = stub.PutState(mortgRedeemKey,[]byte(mortgRate+"_"+mortgDays))
	}
	if self.checkAccountDisabled(account) {
		return shim.Error("The account [" + account + "] has been disabled")
	}
	messageBody := hex.EncodeToString([]byte(strings.ToLower(stub.GetTxID() + account)))
	var checkStatus = self.checkSign(account, messageBody, sign)
	if !checkStatus {
		return shim.Error("Attestation of failure")
	}
	if method != "put" && method != "new" {
		return shim.Error("mortg method invalid")
	}
	mortgageAddr := self.mortgPrefix + mortgageId
	mortgageKey := token + "_" + mortgageAddr
	mortgageValbytes, err := stub.GetState(mortgageKey)
	if err != nil {
		return shim.Error("Error getting frozen balance of mortgage account")
	}
	if mortgageValbytes == nil {
		mortgageAccountVal = decimal.NewFromFloat(0)
	} else {
		mortgageAccountVal, _ = decimal.NewDecimalFromString(string(mortgageValbytes))
	}
	accountVal, err = self.getAccountBalance(stub, token, account);
	if err != nil {
		return shim.Error(err.Error())
	}
	if accountVal.LessThan(mortgageMoney) {
		return shim.Error("Account [" + account + "] has insufficient balance")
	}
	accountVal = accountVal.Sub(mortgageMoney)
	precision, err := self.getTokenPrecision(stub,token)
	if err != nil {
		return shim.Error(err.Error())
	}
	precision32 := int32(precision)
	if err != nil {
		return shim.Error("Token [" + token + "] decimal precision format error")
	}
	err = self.updateAccountBalance(stub, token, account, accountVal, precision)
	if err != nil {
		return shim.Error(err.Error())
	}
	mortgageAccountVal = mortgageAccountVal.Add(mortgageMoney)
	err = stub.PutState(mortgageKey, []byte(mortgageAccountVal.StringFixed(precision32)))
	if err != nil {
		return shim.Error("Account mortgage amount saved error")
	}
	self.saveMortgPutLog(stub, account, mortgageId, mortgageDetailId, self.formatNumber(mortgageMoney, precision))
	superMinerAccount := stub.GetSuperMiner()
	normalMinerAccount := stub.GetNormalMiner()
	var returnString = "mortg;" + method + ";" + token + ",0," + superMinerAccount + "," + normalMinerAccount + ",0,0,0," + account + "," + mortgageAddr + "," + self.formatNumber(mortgageMoney, precision) + "," + self.base64Encode("mortgage put") + "," + stub.GetTxID() + "," + self.base64Encode(mortgageId+":"+mortgageDetailId)
	return shim.Success([]byte(returnString))
}
func (self *TokenChaincode) mortgageGet(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if args == nil || len(args) != 2 {
		fmt.Println("mortgageGet£¨£© parameter require 2 £¨ [1] token [2] mortgageId £©")
		return shim.Error("Missing parameter")
	}
	if args[0] == "" || args[1] == "" {
		for i := 0; i < len(args); i++ {
			fmt.Println("args[", strconv.Itoa(i), "] = ", args[i])
		}
		return shim.Error("Parameters cannot be empty")
	}
	var token, account, accountList, moneyList, mortgageId string
	var money, accountVal, mortgageAccountVal, getMoney, getMoneyTotal decimal.Decimal
	var err error
	token = strings.ToLower(args[0])
	mortgageId = strings.ToLower(args[1])
	accountList = ""
	moneyList = ""
	getMoneyTotal = decimal.NewFromFloat(0)
	trans := [][]byte{[]byte("checkProjectExpire"), []byte(mortgageId)}
	response := stub.InvokeChaincode(self.mortgageChaincode, trans, self.cid)
	response_payload := string(response.Payload)
	if !strings.Contains(response_payload,"expired"){
		return shim.Error("The mortgage is not over yet")
	}
	mortgRedeemKey :="mortgRedeemKey_"+mortgageId
	mortgRedeemByte,err:=stub.GetState(mortgRedeemKey)
	if err != nil {
		return shim.Error("get mortgRedeemInfo error")
	}
	mortgStatus:=false
	mortgRateDecimal:=decimal.RequireFromString("0")
	if mortgRedeemByte != nil{
		mortgRedeemArray:= strings.Split(string(mortgRedeemByte),"_")
		mortgRateString:=mortgRedeemArray[0]
		mortgRateDecimal,err = decimal.NewDecimalFromString(mortgRateString)
		if err != nil{
			return shim.Error("get mortgrate error")
		}
		responsePayloadArray:=strings.Split(response_payload,":")
		mortgDaysString:=mortgRedeemArray[1]
		if !self.checkFitTimeLimit(stub,mortgageId,responsePayloadArray[1],mortgDaysString){
			return shim.Error("It's not time yet")
		}
		mortgStatus = true
	}
	mortgageAddr := self.mortgPrefix + mortgageId
	mortgageKey := token + "_" + mortgageAddr
	mortgageValbytes, err := stub.GetState(mortgageKey)
	if err != nil {
		return shim.Error("Failed to obtain mortgage data")
	}
	if mortgageValbytes == nil {
		return shim.Error("Mortgage data is empty")
	}
	mortgageAccountVal, err = decimal.NewDecimalFromString(string(mortgageValbytes))
	if err != nil {
		return shim.Error("Wrong format of available mortgage amount")
	}
	precision, err := self.getTokenPrecision(stub,token)
	if err != nil {
		return shim.Error(err.Error())
	}
	precision32 := int32(precision)
	if err != nil {
		return shim.Error("Token [" + token + "] decimal precision format error")
	}
	key := mortgageId + "_data"
	historyQuery, err := stub.GetHistoryForKey(key)
	if err != nil {
		return shim.Error(err.Error())
	}
	accountMoney := make(map[string]decimal.Decimal)
	for historyQuery.HasNext() {
		keyModification, err := historyQuery.Next();
		if err != nil {
			return shim.Error(err.Error())
		}
		reldata := string(keyModification.GetValue())
		resultArray := strings.Split(reldata, ",")
		if len(resultArray) == 3 {
			account = resultArray[0]
			money, _ = decimal.NewDecimalFromString(resultArray[2])
			_, isExist := accountMoney[account]
			if !isExist {
				accountMoney[account] = money
			} else {
				accountMoney[account] = money.Add(accountMoney[account])
			}
		}
	}
	if len(accountMoney) == 0 {
		return shim.Error("No amount is redeemable")
	}
	var keys []string
	for _addr, _ := range accountMoney {
		keys = append(keys, _addr)
	}
	sort.Strings(keys)
	for _, _addr := range keys {
		_money := accountMoney[_addr]
		fmt.Println("addr:", _addr, "money:", accountMoney[_addr])
		if mortgStatus{
			_money = _money.Mul(mortgRateDecimal)
		}
		getMoneyTotal = getMoneyTotal.Add(_money)
		if self.checkAccountDisabled(_addr) {
			return shim.Error("The account [" + _addr + "] has been disabled")
		}
		accountVal, err = self.getAccountBalance(stub, token, _addr);
		if err != nil {
			return shim.Error(err.Error())
		}
		mortgageAccountVal = mortgageAccountVal.Sub(_money)
		accountVal = accountVal.Add(_money)
		err = self.updateAccountBalance(stub, token, _addr, accountVal, precision)
		if err != nil {
			return shim.Error(err.Error())
		}
		if accountList != "" {
			accountList = accountList + "|"
		}
		if moneyList != "" {
			moneyList = moneyList + "|"
		}
		accountList = accountList + _addr
		moneyList = moneyList + self.formatNumber(_money, precision)
	}
	if mortgageAccountVal.LessThan(getMoney) {
		fmt.Println("mortgage project real balance:", mortgageAccountVal.String(), " -- ", "getMoneyTotal:", getMoneyTotal.String())
		return shim.Error("The available mortgage amount is insufficient")
	}
	err = stub.PutState(mortgageKey, []byte(mortgageAccountVal.StringFixed(precision32)))
	if err != nil {
		return shim.Error("Account mortgage amount saved error")
	}
	superMinerAccount := stub.GetSuperMiner()
	normalMinerAccount := stub.GetNormalMiner()
	if mortgStatus{
		timeByte:=[]byte(time.Now().In(self.dateZone).Format("2006-01-02 15:04:05"))
		err =stub.PutState("mortgRedeemTime_"+mortgageId,timeByte)
		if err != nil{
			return shim.Error("save mortgRedeem history error")
		}
	}
	var returnString = "mortg;get;" + token + ",0," + superMinerAccount + "," + normalMinerAccount + ",0,0,0," + mortgageAddr + "," + accountList + "," + moneyList + "," + self.base64Encode("mortgage return") + "," + stub.GetTxID()
	return shim.Success([]byte(returnString))
}
func (self *TokenChaincode) applySpareNode(stub shim.ChaincodeStubInterface,args []string) pb.Response{
	if args == nil || len(args) != 8 {
		fmt.Println("directmortgagePut£¨£© parameter require 8 £¨ [1] token [2] account [3] money [4] sign,[5] mortgageId,[6]mortgProjectId [7]mortgReturnRate,[8]mortgReturnDays£©")
		return shim.Error("Missing parameter")
	}
	if args[0] == "" || args[1] == "" || args[2] == "" || args[3] == "" || args[4] == "" || args[5] == ""|| args[6] == ""|| args[7] == "" {
		for i := 0; i < len(args); i++ {
			fmt.Println("args[", strconv.Itoa(i), "] = ", args[i])
		}
		return shim.Error("Parameters cannot be empty")
	}
	var token, account, sign,mortgid,mortgProjectId,mortgReturnRate,mortgReturnDays string
	var nodeAmount,accountVal,mortgAmount decimal.Decimal
	var err error
	token = strings.ToLower(args[0])
	account = strings.ToLower(args[1])
	nodeAmount, err = decimal.NewDecimalFromString(args[2])
	sign = args[3]
	mortgid = args[4]
	mortgProjectId = args[5]
	mortgReturnRate = args[6]
	mortgReturnDays = args[7]
	if err != nil {
		return shim.Error("Incorrect mortgage amount format")
	}
	if self.checkAccountDisabled(account) {
		return shim.Error("The account [" + account + "] has been disabled")
	}
	messageBody := hex.EncodeToString([]byte(strings.ToLower(stub.GetTxID() + account)))
	var checkStatus = self.checkSign(account, messageBody, sign)
	if !checkStatus {
		return shim.Error("Attestation of failure")
	}
	accountVal, err = self.getAccountBalance(stub, token, account)
	if err != nil {
		return shim.Error(err.Error())
	}
	if accountVal.LessThan(nodeAmount) {
		return shim.Error("Account [" + account + "] has insufficient balance")
	}
	accountVal = accountVal.Sub(nodeAmount)
	precision, err := self.getTokenPrecision(stub,token)
	if err != nil {
		return shim.Error(err.Error())
	}
	/*precision32 := int32(precision)
	if err != nil {
		return shim.Error("Token [" + token + "] decimal precision format error")
	}*/
	err = self.updateAccountBalance(stub, token, account, accountVal, precision)
	if err != nil {
		return shim.Error(err.Error())
	}
	superMinerAccount := stub.GetSuperMiner()
	normalMinerAccount := stub.GetNormalMiner()
	mortgAccount:=self.mortgPrefix+mortgProjectId
	mortgAmount, err = self.getAccountBalance(stub, token, mortgAccount)
	if err != nil {
		return shim.Error(err.Error())
	}
	mortgRedeemKey :="mortgRedeemKey_"+mortgid
	stub.PutState(mortgRedeemKey,[]byte(mortgReturnRate+"_"+mortgReturnDays))
	self.saveMortgPutLog(stub, account, mortgProjectId, mortgid, self.formatNumber(nodeAmount, precision))
	mortgAmount = mortgAmount.Add(nodeAmount)
	err = self.updateAccountBalance(stub, token, mortgAccount, mortgAmount, precision)
	if err != nil {
		return shim.Error(err.Error())
	}
	var returnString = "node;applyNode;" + token + ",0," + superMinerAccount + "," + normalMinerAccount + ",0,0,0," + account + "," + mortgAccount+ "," + self.formatNumber(nodeAmount,precision)+ ","+self.base64Encode(mortgid)+"," + stub.GetTxID()
	return shim.Success([]byte(returnString))
}
func (self *TokenChaincode) nodeMortReturn(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if args == nil || len(args) != 2 {
		fmt.Println("tokenMortgageGet£¨£© parameter require 2 £¨ [1] account [2] nodeMortId£©")
		return shim.Error("Missing parameter")
	}
	if args[0] == "" || args[1] == "" {
		for i := 0; i < len(args); i++ {
			fmt.Println("args[", strconv.Itoa(i), "] = ", args[i])
		}
		return shim.Error("Parameters cannot be empty")
	}
	superMinerAccount := stub.GetSuperMiner()
	normalMinerAccount := stub.GetNormalMiner()
	account := strings.ToLower(args[0])
	nodeMortId := strings.ToLower(args[1])
	trans:= [][]byte{[]byte("nodeMortgReturn"), []byte(account), []byte(nodeMortId)}
	response := stub.InvokeChaincode(self.mortgageChaincode, trans, self.cid)
	if response.Status != 200 {
		return shim.Error("Chain code call failed"+response.Message)
	}
	nodeMortReturnByte:=response.Payload
	nodeMortArray:=strings.Split(string(nodeMortReturnByte),",")
	nodeMortReturnDecimal:=decimal.RequireFromString(nodeMortArray[0])
	token:= nodeMortArray[1]
	precision, err := self.getTokenPrecision(stub,token)
	if err != nil {
		return shim.Error(err.Error())
	}
	directProjectId :=nodeMortArray[2]
	tokenMortBalanceDecimal,err:=self.getAccountBalance(stub,token,self.mortgPrefix+directProjectId)
	if err != nil{
		return shim.Error("get tokenMort balance error")
	}
	if tokenMortBalanceDecimal.LessThanOrEqual(decimal.RequireFromString("0")) {
		return  shim.Error("tokenMort has Redeemed")
	}
	if tokenMortBalanceDecimal.LessThan(nodeMortReturnDecimal){
		nodeMortReturnDecimal = tokenMortBalanceDecimal
	}
	tokenMortBalanceDecimal = tokenMortBalanceDecimal.Sub(nodeMortReturnDecimal)
	err = self.updateAccountBalance(stub,token,self.tokenMortPre+directProjectId,tokenMortBalanceDecimal,precision)
	if err != nil{
		return shim.Error("update tokenMort balance error ")
	}
	accountBalanceDecimal,err:=self.getAccountBalance(stub,self.mainToken.name,account)
	if err != nil{
		return shim.Error("get tokenMort balance error")
	}
	accountBalanceDecimal = accountBalanceDecimal.Add(nodeMortReturnDecimal)
	err = self.updateAccountBalance(stub,token,account,accountBalanceDecimal,precision)
	if err != nil{
		return shim.Error("update tokenMort balance error ")
	}
	var returnString = "token;transfer;ctk,0," + superMinerAccount + "," + normalMinerAccount + ",0,0,0," + self.tokenMortPre+nodeMortId + "," + account + "," + self.formatNumber(nodeMortReturnDecimal, precision) + ",," + stub.GetTxID()
	return shim.Success([]byte(returnString))
}
func (self *TokenChaincode) applyUnitFoundTransfer(stub shim.ChaincodeStubInterface,args []string) pb.Response{
	if args == nil || len(args) != 5 {
		fmt.Println("directmortgagePut£¨£© parameter require 5 £¨ [1] token [2] account [3] money [4] foundToken [5] sign£©")
		return shim.Error("Missing parameter")
	}
	if args[0] == "" || args[1] == "" || args[2] == "" || args[3] == "" || args[4] == ""  {
		for i := 0; i < len(args); i++ {
			fmt.Println("args[", strconv.Itoa(i), "] = ", args[i])
		}
		return shim.Error("Parameters cannot be empty")
	}
	var token, account, sign string
	var freeFound,accountVal,agentPayAccountAmount decimal.Decimal
	var err error
	token = strings.ToLower(args[0])
	account = strings.ToLower(args[1])
	freeFound, err = decimal.NewDecimalFromString(args[2])
	sign = args[4]
	if err != nil {
		return shim.Error("Incorrect mortgage amount format")
	}
	if self.checkAccountDisabled(account) {
		return shim.Error("The account [" + account + "] has been disabled")
	}
	messageBody := hex.EncodeToString([]byte(strings.ToLower(stub.GetTxID() + account+args[2])))
	var checkStatus = self.checkSign(account, messageBody, sign)
	if !checkStatus {
		return shim.Error("Attestation of failure")
	}
	accountVal, err = self.getAccountBalance(stub, token, account)
	if err != nil {
		return shim.Error(err.Error())
	}
	if accountVal.LessThan(freeFound) {
		return shim.Error("Account [" + account + "] has insufficient balance")
	}
	accountVal = accountVal.Sub(freeFound)
	precision, err := self.getTokenPrecision(stub,token)
	if err != nil {
		return shim.Error(err.Error())
	}
	/*precision32 := int32(precision)
	if err != nil {
		return shim.Error("Token [" + token + "] decimal precision format error")
	}*/
	err = self.updateAccountBalance(stub, token, account, accountVal, precision)
	if err != nil {
		return shim.Error(err.Error())
	}
	superMinerAccount := stub.GetSuperMiner()
	normalMinerAccount := stub.GetNormalMiner()
	agentPayAddr := self.getAgentPayAddr(token)
	agentPayAccountAmount, err = self.getAccountBalance(stub, token,agentPayAddr)
	if err != nil {
		return shim.Error(err.Error())
	}
	agentPayAccountAmount = agentPayAccountAmount.Add(freeFound)
	err = self.updateAccountBalance(stub, token, agentPayAddr, agentPayAccountAmount, precision)
	if err != nil {
		return shim.Error(err.Error())
	}
	var returnString = "freeFound;transfer;" + token + ",0," + superMinerAccount + "," + normalMinerAccount + ",0,0,0," + account + "," + self.tokenFoundAccount+ "," + self.formatNumber(freeFound,precision)+ "," + "" + "," + stub.GetTxID()
	return shim.Success([]byte(returnString))
}
func (self *TokenChaincode) candyPut(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if args == nil || len(args) != 6 {
		fmt.Println("candyPut£¨£© parameter require 6 £¨ [1] candyToken [2] account [3] money [4] candyId [5] sign [6] token£©")
		return shim.Error("Missing parameter")
	}
	if args[0] == "" || args[1] == "" || args[2] == "" || args[3] == "" || args[4] == "" {
		for i := 0; i < len(args); i++ {
			fmt.Println("args[", strconv.Itoa(i), "] = ", args[i])
		}
		return shim.Error("Parameters cannot be empty")
	}
	var candyToken, account, sign, candyId string
	var candyMortgageMoney, accountVal decimal.Decimal
	var err error
	candyToken = strings.ToLower(args[0])
	account = strings.ToLower(args[1])
	candyMortgageMoney, err = decimal.NewDecimalFromString(args[2])
	candyId = strings.ToLower(args[3])
	sign = args[4]
	limitToken := args[5]
	if err != nil {
		return shim.Error("Incorrect mortgage amount format")
	}
	if candyMortgageMoney.LessThanOrEqual(decimal.NewFromFloat(0)) {
		return shim.Error("The total number of candy can not be less than zero.")
	}
	if self.checkAccountDisabled(account) {
		return shim.Error("The account [" + account + "] has been disabled")
	}
	messageBody := hex.EncodeToString([]byte(strings.ToLower(stub.GetTxID() + account)))
	var checkStatus = self.checkSign(account, messageBody, sign)
	if !checkStatus {
		return shim.Error("Attestation of failure")
	}
	candyAddr := self.candyPrefix + candyId
	candyKey := candyToken + "_" + candyAddr
	candyValbytes, err := stub.GetState(candyKey)
	if err != nil {
		return shim.Error("Error getting candy balance of candy account")
	}
	if candyValbytes != nil {
		return shim.Error("The candy box project already exists")
	}
	if limitToken != self.mainToken.name && limitToken != "" {
		precisionValbytesToken, err := stub.GetState(limitToken + self.precisionKey)
		if err != nil || precisionValbytesToken == nil {
			return shim.Error("Token [" + limitToken + "] was not found")
		}
	}
	precision, err := self.getTokenPrecision(stub,candyToken)
	if err != nil {
		return shim.Error(err.Error())
	}
	accountVal, err = self.getAccountBalance(stub, candyToken, account)
	if err != nil {
		return shim.Error(err.Error())
	}
	if accountVal.LessThan(candyMortgageMoney) {
		return shim.Error("Account [" + account + "] has insufficient balance")
	}
	accountVal = accountVal.Sub(candyMortgageMoney)
	err = self.updateAccountBalance(stub, candyToken, account, accountVal, precision)
	if err != nil {
		return shim.Error("candyPut Error saving account [" + account + "] balance")
	}
	precision32 := int32(precision)
	if err != nil {
		return shim.Error("Token [" + candyToken + "] decimal precision format error")
	}
	err = stub.PutState(candyKey, []byte(candyMortgageMoney.StringFixed(precision32)))
	if err != nil {
		return shim.Error("Account mortgage amount saved error")
	}
	candyUserKey := candyKey + "_account"
	err = stub.PutState(candyUserKey, []byte(account))
	if err != nil {
		return shim.Error("Failed to save the association")
	}
	superMinerAccount := stub.GetSuperMiner()
	normalMinerAccount := stub.GetNormalMiner()
	var returnString = "candy;put;" + candyToken + ",0," + superMinerAccount + "," + normalMinerAccount + ",0,0,0," + account + "," + candyAddr + "," + self.formatNumber(candyMortgageMoney, precision) + "," + self.base64Encode(candyId) + "," + stub.GetTxID()
	return shim.Success([]byte(returnString))
}
func (self *TokenChaincode) candyWidthDraw(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 7 {
		fmt.Println("candyWidthDraw£¨£© parameter require 7 £¨ [1] candyToken [2] candyId [3] account [4] sign [5] token [6] drawLimit [7] candyRate  £©")
		return shim.Error("Missing parameter")
	}
	if args[0] == "" || args[1] == "" || args[2] == "" || args[3] == "" || args[5] == "" {
		return shim.Error("Parameters cannot be empty")
	}
	candyToken := args[0]
	candyId := args[1]
	account := args[2]
	sign := args[3]
	token := args[4]
	drawLimitString := args[5]
	candyRate := args[6]
	drawlimit, err := decimal.NewDecimalFromString(drawLimitString)
	if err != nil {
		return shim.Error("Limit amount conversion failed")
	}
	if self.checkAccountDisabled(account) {
		return shim.Error("The account [" + account + "] has been disabled")
	}
	candyAddr := self.candyPrefix + candyId
	candyKey := candyToken + "_" + candyAddr
	candyAmountByte, err := stub.GetState(strings.ToLower(candyKey))
	if err != nil {
		return shim.Error("get candy amount err")
	}
	accountBalance, err := decimal.NewDecimalFromString("0")
	if err != nil {
		return shim.Error("Failed to initialize account balance")
	}
	accountBalance, err = self.getAccountBalance(stub, token, account)
	if err != nil {
		return shim.Error("get token balance err")
	}
	candyTotalDecimal, err := decimal.NewDecimalFromString(string(candyAmountByte))
	behooveAmount := decimal.RequireFromString("0")
	if candyRate != "" {
		candyRateDecimal, err := decimal.NewDecimalFromString(candyRate)
		if err != nil {
			return shim.Error("Candy box conversion failure")
		}
		calculateAmount := accountBalance.Mul(candyRateDecimal)
		if calculateAmount.GreaterThanOrEqual(drawlimit) {
			behooveAmount = drawlimit
		} else {
			behooveAmount = calculateAmount
		}
	} else {
		behooveAmount = drawlimit
	}
	if behooveAmount.GreaterThan(candyTotalDecimal) {
		behooveAmount = candyTotalDecimal
	}
	if behooveAmount.LessThanOrEqual(decimal.RequireFromString("0")) {
		widthDrawFail := "candy;widthDraw;" + candyId + "," + stub.GetTxID()
		return shim.Error(widthDrawFail)
	}
	precision, err := self.getTokenPrecision(stub,candyToken)
	if err != nil {
		return shim.Error(err.Error())
	}
	messageBody := hex.EncodeToString([]byte(strings.ToLower(stub.GetTxID() + account)))
	var checkStatus = self.checkSign(account, messageBody, sign)
	if !checkStatus {
		return shim.Error("Attestation of failure")
	}
	userTokenDecimal, err := self.getAccountBalance(stub, candyToken, account)
	if err != nil {
		return shim.Error("query account balance err")
	}
	userTokenDecimal = userTokenDecimal.Add(behooveAmount)
	err = self.updateAccountBalance(stub, candyToken, account, userTokenDecimal, precision)
	if err != nil {
		return shim.Error("update candy user balance err")
	}
	candyTotalDecimal = candyTotalDecimal.Sub(behooveAmount)
	err = stub.PutState(candyKey, []byte(candyTotalDecimal.StringFixed(int32(precision))))
	if err != nil {
		return shim.Error("update candy project balance err")
	}
	superMinerAccount := stub.GetSuperMiner()
	normalMinerAccount := stub.GetNormalMiner()
	returnString := "candy;widthDraw;" + candyToken + ",0," + superMinerAccount + "," + normalMinerAccount + ",0,0,0," + candyAddr + "," + account + "," + self.formatNumber(behooveAmount, precision) + "," + self.base64Encode(candyId) + "," + stub.GetTxID()
	return shim.Success([]byte(returnString))
}
func (self *TokenChaincode) candyMortgageGet(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if args == nil || len(args) != 3 {
		fmt.Println("candyMortgageGet£¨£© parameter require 3 £¨ [1] account [2] candyId [3] token£©")
		return shim.Error("Missing parameter")
	}
	if args[0] == "" || args[1] == "" || args[2] == "" {
		for i := 0; i < len(args); i++ {
			fmt.Println("args[", strconv.Itoa(i), "] = ", args[i])
		}
		return shim.Error("Parameters cannot be empty")
	}
	var account, candyId, token string
	var accountVal, candyMortgageAccountVal decimal.Decimal
	var err error
	account = strings.ToLower(args[0])
	candyId = strings.ToLower(args[1])
	token = args[2]
	if self.checkAccountDisabled(account) {
		return shim.Error("The account [" + account + "] has been disabled")
	}
	precision, err := self.getTokenPrecision(stub,token)
	if err != nil {
		return shim.Error(err.Error())
	}
	precision32 := int32(precision)
	if err != nil {
		return shim.Error("Token [" + token + "] decimal precision format error")
	}
	candyAddr := self.candyPrefix + candyId
	candyMortgageKey := token + "_" + candyAddr
	candyUserKey := candyMortgageKey + "_account"
	candyUserRelationByte, err := stub.GetState(candyUserKey)
	if err != nil {
		return shim.Error("Failed to obtain candy box correlation")
	}
	if candyUserRelationByte == nil || len(candyUserRelationByte) == 0 {
		return shim.Error("Failed to obtain candy box correlation")
	}
	if account != string(candyUserRelationByte) {
		return shim.Error("Have no right to redeem")
	}
	candyMortgageValbytes, err := stub.GetState(candyMortgageKey)
	if err != nil {
		return shim.Error("Error getting frozen balance of candyMortgage account")
	}
	if candyMortgageValbytes == nil {
		return shim.Error("There is no candy box mortgage data")
	}
	candyMortgageAccountVal, err = decimal.NewDecimalFromString(string(candyMortgageValbytes))
	if err != nil {
		return shim.Error("The amount of mortgage can be used in wrong format")
	}
	if candyMortgageAccountVal.LessThanOrEqual(decimal.RequireFromString("0")) {
		returnString := "candy;get;" + candyId + "," + stub.GetTxID()
		return shim.Error(returnString)
	}
	accountVal, err = self.getAccountBalance(stub, token, account)
	if err != nil {
		return shim.Error(err.Error())
	}
	accountVal = accountVal.Add(candyMortgageAccountVal)
	err = self.updateAccountBalance(stub, token, account, accountVal, precision)
	if err != nil {
		return shim.Error("Error saving account [" + account + "] balance")
	}
	err = stub.PutState(candyMortgageKey, []byte(decimal.RequireFromString("0").StringFixed(precision32) ))
	if err != nil {
		return shim.Error("Account mortgage amount saved error")
	}
	superMinerAccount := stub.GetSuperMiner()
	normalMinerAccount := stub.GetNormalMiner()
	returnString := "candy;get;" + token + ",0," + superMinerAccount + "," + normalMinerAccount + ",0,0,0," + candyAddr + "," + account + "," + candyMortgageAccountVal.StringFixed(precision32) + "," + self.base64Encode(candyId) + "," + stub.GetTxID()
	return shim.Success([]byte(returnString))
}
func (self *TokenChaincode) checkSign(pubaddress string, message string, signedAddress string) bool {
	self.printDateStr()
	var err error
	var addr = strings.ToLower(pubaddress)
	/*if addrpre == "0x" {
		addr = string([]byte(addr)[2:])
	}*/
	addr = string([]byte(addr)[2:])
	var msg = crypto.Keccak256([]byte(message))
	sign, err := hex.DecodeString(strings.ToLower(signedAddress))
	if err != nil {
		fmt.Printf("\n sign decode error: %s", err)
		return false
	}
	recoveredPub, err := crypto.Ecrecover(msg, sign)
	if err != nil {
		fmt.Printf("\n ECRecover error: %s", err)
		return false
	}
	pubKey := crypto.ToECDSAPub(recoveredPub)
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	recoveredAddrstr := hex.EncodeToString(recoveredAddr[:])
	if addr != recoveredAddrstr {
		fmt.Println("\n sign of account addres is", "0x"+recoveredAddrstr, "not is", "0x"+addr)
		return false
	} else {
		return true
	}
}
func (self *TokenChaincode) getAccountMoneyHistory(stub shim.ChaincodeStubInterface, account, token string, isDebug bool) (AccountBalanceList, error) {
	var list = AccountBalanceList{}
	var key = strings.ToLower(token + "_" + account)
	historyQuery, err := stub.GetHistoryForKey(key)
	if err != nil {
		return list, err
	}
	var i = 0
	for historyQuery.HasNext() {
		keyModification, err := historyQuery.Next();
		if err != nil {
			return list, err
		}
		if isDebug {
			fmt.Println("historyQuery.HasNext()", "Timestamp:", keyModification.Timestamp.Seconds, "value:", string(keyModification.GetValue()))
		}
		list = append(list, &AccountBalance{string(keyModification.GetValue()), keyModification.Timestamp.Seconds})
		i++
	}
	return list, nil
}
func (self *TokenChaincode) updateAccountBalance(stub shim.ChaincodeStubInterface, token, account string, balance decimal.Decimal, precision int64) error {
	var key = token + "_" + account
	var err = stub.PutState(key, []byte( balance.StringFixed(int32(precision)) ))
	if err != nil {
		return errors.New("Error saving account [" + account + "] balance")
	}
	return nil
}
func (self *TokenChaincode) getAccountBalance(stub shim.ChaincodeStubInterface, token string, account string) (decimal.Decimal, error) {
	key := token + "_" + account
	return self.getStateDecimalValue(stub, key)
}
func (self *TokenChaincode) getTokenPrecision(stub shim.ChaincodeStubInterface, token string) (precision int64,err error) {
	precisionValbytes, err := stub.GetState(token + self.precisionKey)
	if err != nil || precisionValbytes == nil {
		return precision,errors.New("Token [" + token + "] was not found")
	}
	precision, err = strconv.ParseInt(string(precisionValbytes), 10, 32)
	if err != nil {
		return precision,errors.New("Token [" + token + "] precision format error")
	}
	return precision,err
}
func (self *TokenChaincode) checkAccountDisabled(account string) bool {
	if account == self.blackKey || account == self.mineralKey {
		return true
	}
	accountLen := len(account)
	if accountLen != 42 {
		return true
	}
	accountPrefix := string(account[0:10])
	if accountPrefix == self.tokenTransPrefix || accountPrefix == self.mortgPrefix || accountPrefix == self.candyPrefix {
		return true
	}
	return false
}
func (self *TokenChaincode) getStateDecimalValue(stub shim.ChaincodeStubInterface, key string) (decimal.Decimal, error) {
	var err error
	var keyVal decimal.Decimal
	var keyValbytes []byte
	var ling = decimal.NewFromFloat(0)
	keyValbytes, err = stub.GetState(strings.ToLower(key))
	if err != nil {
		return ling, errors.New("Failed to get state data [" + key + "]")
	}
	if keyValbytes == nil {
		keyVal = decimal.NewFromFloat(0)
		return keyVal, nil
	} else {
		keyVal, err = decimal.NewDecimalFromString(string(keyValbytes))
		if err != nil {
			fmt.Println("NewDecimalFromString error , key: " + key + " val: " + string(keyValbytes))
			return ling, errors.New("Failed to convert state data")
		}
		return keyVal, nil
	}
}
func (self *TokenChaincode) foundUnitFree(stub shim.ChaincodeStubInterface,foundToken,outAccount string,minerFeeVal decimal.Decimal)(bool,string){
	foundTokenKey:=self.mainToken.name+"_"+foundToken+"_"+self.tokenFoundAccount
	foundTokenKey = strings.ToLower(foundTokenKey)
	foundTokenValDecimal,err:= self.getStateDecimalValue(stub,foundTokenKey)
	if err != nil {
		return  false,"get founTokenval err"
	}
	if foundTokenValDecimal.GreaterThanOrEqual(minerFeeVal){
		foundTokenValDecimal = foundTokenValDecimal.Sub(minerFeeVal)
		outAccountMainVal, err := self.getAccountBalance(stub, self.mainToken.name, outAccount)
		if err != nil {
			fmt.Println("error:", err.Error())
			return false,"foundTokenUnit get outaccount balance error"
		}
		outAccountMainVal = outAccountMainVal.Add(minerFeeVal)
		err = self.updateAccountBalance(stub,self.mainToken.name,outAccount,outAccountMainVal,self.mainToken.precision)
		if err != nil {
			return false,"save foundTokenUnit outaccount error"
		}
		err = self.updateAccountBalance(stub,self.mainToken.name,foundToken+"_"+self.tokenFoundAccount,foundTokenValDecimal,self.mainToken.precision)
		if err != nil {
			return false,"save foundTokenUnit outaccount error"
		}
		outAccount = outAccount +"|"+self.tokenFoundAccount
		return true,outAccount
	}else{
		return true,""
	}
}
func (self *TokenChaincode) getPrecisionLen(number string) int64 {
	numLen := 0
	if strings.Contains(number, ".") {
		numString1 := strings.Split(number, ".")
		numLen = len(numString1[1])
	}
	return int64(numLen)
}
func (t *TokenChaincode) makeDetailId(account,mortgageId string) string {
	nowTime := time.Now().In(t.dateZone).Unix()
	return strings.ToLower(t.md5string(account + mortgageId +  strconv.FormatInt(nowTime,10)))
}
func (self *TokenChaincode) checkRemarkFormat(remark string) bool {
	match, err := regexp.MatchString(`^[A-Za-z0-9+/=]{0,64}$`, remark)
	if err != nil {
		return false
	}
	return match
}
func (self *TokenChaincode) checkAccountFormat(account string) bool {
	match, err := regexp.MatchString(`^`+self.addressPrefix+`[A-Za-z0-9]{40}$`, account)
	if err != nil {
		return false
	}
	return match
}
func (self *TokenChaincode) checkTokenFormat(token string) bool {
	match, err := regexp.MatchString(`^[a-z0-9-]{2,20}$`, token)
	if err != nil {
		return false
	}
	return match
}
func  (self *TokenChaincode) checkFitTimeLimit(stub shim.ChaincodeStubInterface,mortgageId,endTime,mortgDaysString string)bool{
	mortgRedeemKey:="mortgRedeemTime_"+mortgageId
	mortgRedeemTimeByte,err:= stub.GetState(mortgRedeemKey)
	if err != nil{
		return false
	}
	mortgDays,err:=strconv.Atoi(mortgDaysString)
	if err != nil {
		return false
	}
	nowTime:=time.Now().In(self.dateZone)
	days := 0
	if mortgRedeemTimeByte != nil{
		lastRedeemTime,err:= time.Parse("2006-01-02 15:04:05",string(mortgRedeemTimeByte))
		if err != nil{
			return false
		}
		days = self.calSubDays(nowTime,lastRedeemTime)
	}else{
		endTimeInt64,err:=strconv.ParseInt(endTime,10,64)
		if err != nil{
			return false
		}
		endTime:=time.Unix(endTimeInt64,0).In(self.dateZone)
		days = self.calSubDays(nowTime,endTime)
	}
	if days == 0{
		return false
	}
	if days%mortgDays == 0{
		return true
	}else{
		return false
	}
}
func (self *TokenChaincode) getAgentPayAddr(token string) (string){
	return strings.ToLower(self.agentPayPrefix + self.md5string(token))
}
func (self *TokenChaincode) getServiceCharge(goRatio decimal.Decimal) (chargeLowLimit,tokenCharge decimal.Decimal){
	chargeLowLimit = self.chargeLowLimit.Mul(goRatio)
	tokenCharge = self.tokenCharge.Mul(goRatio)
	return chargeLowLimit,tokenCharge
}
func (self *TokenChaincode) calSubDays(t1,t2 time.Time)int{
	t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, self.dateZone)
	t2 = time.Date(t2.Year(), t2.Month(), t2.Day(), 0, 0, 0, 0, self.dateZone)
	sd:=t1.Sub(t2).Hours()
	days:=int(sd/(24))
	return days
}
func (self *TokenChaincode) base64Encode(content string) string {
	var encoded bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &encoded)
	encoder.Write([]byte(content))
	encoder.Close()
	return string(encoded.Bytes())
}
func (self *TokenChaincode) formatNumber(number decimal.Decimal, precision int64) string {
	return number.StringFixed(int32(precision))
}
func (self *TokenChaincode) md5string(str string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(str))
	cipherStr := md5Ctx.Sum(nil)
	return strings.ToLower(hex.EncodeToString(cipherStr))
}
func (self *TokenChaincode) printDateStr() {
	shijian := time.Now().In(self.dateZone).Format("2006-01-02 15:04:05")
	fmt.Printf("\n%s ", shijian)
}
func (self *TokenChaincode) printDebugStr(str string) {
	shijian := time.Now().In(self.dateZone).Format("2006-01-02 15:04:05")
	fmt.Printf("\n%s %s", shijian, str)
}
func main() {
	var err error
	var token = new(TokenChaincode)
	token.rand = rand.New(rand.NewSource(time.Now().Unix()))
	var tryGetDecimal = func(number string) decimal.Decimal {
		newDecimal, err := decimal.NewDecimalFromString(number);
		if err != nil {
			panic("Error NewDecimalFromString:" + number)
		}
		return newDecimal
	}
	mainConfig := TokenConfig{}
	mainConfig.name = "ctk"
	mainConfig.title = "ctk"
	mainConfig.originatorAccount = "0x8fc2e72e0532d50addd67ca3bf120f645d9f9239"
	mainConfig.output = tryGetDecimal("24498550000")
	mainConfig.gross = tryGetDecimal("24498550000")
	mainConfig.mineral = tryGetDecimal("0")
	mainConfig.precision = 6
	mainConfig.rate = tryGetDecimal("0")
	mainConfig.dayAward = tryGetDecimal("0")
	mainConfig.site = ""
	mainConfig.logo = "ffd8ffe000104a46494600010101004800480000ffdb0043000604040405040605050609060506090b080606080b0c0a0a0b0a0a0c100c0c0c0c0c0c100c0e0f100f0e0c1313141413131c1b1b1b1c20202020202020202020ffdb0043010707070d0c0d181010181a1511151a20202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020ffc00011080028002803011100021101031101ffc4001a000002030101000000000000000000000006070405080103ffc40033100002010302040404030900000000000001020304051100060712131421223141233252710881b1162425425153627291ffc4001a010002030101000000000000000000000000040203050106ffc40028110002010302050403010000000000000000010203041112311321415161053271a1222343b1ffda000c03010002110311003f00d27b8b73da36fd09abb8c8c078f4e1891a59a423d9234058fe83df5194923b18e4495ebf14f574f7311d0eda6ec51bcfdec8d0cee3fc54210bf9e74bbb9f032adbc87bb4f8efb1771dc68adb1ccf475d704069e2a8007c6ce1a9c9048127d3ece3e524f86af8d44ca6549a09ed5bc286e9b8ae167a18de78eda8bdcdc53069c4e4e1a9f3fdc5183e19d0a69bc7638e0d24fb97dcda99023cf23c70c92451f5654462918c02c40c85c9f4c9d00660e28f0bf884529eed7edc14f75b957cff000ad923c88907379e4e93366358e11e04851ff4e90b9fc16a931ea0f53c450bf4d93b8a2e8565446b1d075b0d2c72069311927a91f2ff0029e43ca7d7df483b98f4f70ec28b6f9ec326c9c5bdc967b8536d8dbb476ab6c52f4a692b2f3298e390b228662e1d0282f9fa989d3d6d5b2b1115b9a387991a4adf2d6bd0d3bd7471c558d1a9a88e0732441f1e6e47210b2e7d09035a06710ef97ca1b25aaa2e95c5bb6a600b08d79dd8b3044555f76666006a32928acbd8ec6397842937deeb9aeb73b7d6d1c51d04d4c9352d32d5f467925ee0a160b130922571d31ca724fe9ac6adea597fad6dd59ab42c57f47bf6012f17faaae92919960a74a74e9d674be1c72b22f4fb82aec7cff572f97d73a52bd65539e9c3ff004d0b7b7e175ca033728335251dc6a20eb24552d1f4cb14eb52bf881cff0030048f03f9eacb6694f1e3ec85d272867cfd1aaf85535a64d81687b3d65455dacc67b5ef195e7814310699dd42e7a0c0a7d86bd0c763cf4b72daff006d175b2d65bc901aa23c44cde824521e327d7d1d46a3521aa2d77084b4bcf613b7ba1bd5aa7c5759a4ef3d29a78d7b881bdfcb3202478fb633af3cfd3ea45e3a1b2af20f9f52c2e3c31bc57f0ce3a4b6d2a457e7ac3587b92b148caa8d14272e0f261394f4dbfa9cf9b5b76d423c2d335beff265dc569713541edb7c0afbcf03f8c75702c95b47de2f384edd2b21771cde920505630148f124eaea76d461ed8e0a6a5c579fba59349ec3b75f2d9b46db457e9526bcc517efd2c5ca54c9fec8918638c65b1927dcfaea67095bd2df7f9f6cdda9b6fcfdade67a774b65516e411cc7e562d86c63eda005b1db1c63e41dbd6d4d3c41631594d35f3b896a489149e84fda0ed7ca1b270739c6803d7f6238a46faa56fb3a589e480b196b44d5f040e90f791c539841f1747c118c8c7e401d9768f16239a8fa777aaa9b7b451b5ce845d8c530a9ccaac60aa78246e905e99e43ea7eda0038d8564dcb6fa2aea6bed74971fe235325baa6797ad2f60c476e923f2a79c0ce7401ffd9"
	mainConfig.email = ""
	token.mainToken = mainConfig
	token.addressPrefix = "0x"
	token.cid = "ctk-main"
	token.blackDestroyDateKey = "black_destroy_info"
	token.minersFeeDateKey = "miners_fee_info"
	token.mortgageChaincode = "ctk-mortg"
	token.maxPrecision = 10
	token.maxCurrency = "1000000000000"
	token.configKey = "_config"
	token.mineralKey = token.addressPrefix+"0000000000000000000000000000000000000001"
	token.blackKey = token.addressPrefix+"0000000000000000000000000000000000000000"
	token.minerFeeAccount =  token.addressPrefix + "0000000000000000000000000000000000000002"
	token.tokenChargeDestroy = tryGetDecimal("30")
	token.chargeDestroyUpperLimit = tryGetDecimal("24477550000")
	token.chargeLowLimit = tryGetDecimal("1")
	token.tokenCharge = tryGetDecimal("3000")
	token.tokenMortAmount = "10000"
	token.tokenMortPre = "0x00000004"
	token.tokenMortDays = "30"
	token.tokenMortRate = "0.1"
	token.normalMinerTokenAwardAll = tryGetDecimal("150")
	token.superMinerTokenAwardAll = tryGetDecimal("150")
	token.grossKey = "_gross"
	token.outputKey = "_output"
	token.rateKey = "_rate"
	token.dayRateKey = "_dayrate"
	token.precisionKey = "_precision"
	token.tokenAwardKey = "_award_"
	token.tokenTransPrefix = token.addressPrefix+"00000001"
	token.mortgPrefix = token.addressPrefix+"00000002"
	token.candyPrefix = token.addressPrefix+"00000003"
	token.agentPayPrefix = token.addressPrefix+"00000005"
	token.feeAccount = 0.01
	token.tokenSpareNodeAccount = "0xa867b34ff1b5c80a7658fcbd1012e568c66af10a"
	token.tokenFoundAccount = "0xa867b34ff1b5c80a7658fcbd1012e568c66af101"
	token.storageFeeAccount = token.addressPrefix+"E57603c1c96F6190763F9fd1cf3Dcbf977efC8E6"
	token.dateZone = time.FixedZone("CST", 28800)
	err = shim.Start(token)
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}
