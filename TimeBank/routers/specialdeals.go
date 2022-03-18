package routers

import (
	"encoding/json"
	"fmt"

	"github.com/Jxnuyang/20-software-qkl/chaincode/TimeBank/lib"
	"github.com/Jxnuyang/20-software-qkl/chaincode/TimeBank/utils"

	//"math/rand"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

// var letters = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// func randSeq(n int) string {
// 	b := make([]rune, n)
// 	for i := range b {
// 		b[i] = letters[rand.Intn(len(letters))]
// 	}
// 	return string(b)
// }

//转移工分
func TransferAsset(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	//_ = fmt.Sprintf("TransferAsset starting ...")
	//fmt.Println("TransferAsset starting ...")
	var SenderID string   //发送方ID
	var ReceiverID string //接收方ID
	var Amount float64    //工分交易金额

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting right number of information!!!")
	}
	SenderID = args[0]
	ReceiverID = args[1]
	Amount, _ = strconv.ParseFloat(args[2], 64)
	// if val, err := strconv.ParseFloat(amount, 64); err != nil {
	// 	return shim.Error(fmt.Sprintf("amount参数格式转换出错: %s", err))
	// } else {
	// 	Amount = val
	// }
	sender, err := utils.QueryLedger(stub, lib.UserKey, []string{SenderID})
	if err != nil {
		return shim.Error(fmt.Sprintf("发送方ID有误%s", err))
	}
	var senderlist lib.User
	_ = json.Unmarshal(sender[0], &senderlist)

	receiver, err := utils.QueryLedger(stub, lib.UserKey, []string{ReceiverID})
	if err != nil {
		return shim.Error(fmt.Sprintf("接受方ID有误%s", err))
	}
	var receiverlist lib.User
	_ = json.Unmarshal(receiver[0], &receiverlist)

	//验证发生方工分是否足够
	if senderlist.UserAsset < Amount {
		//fmt.Println("Insufficient transaction failed...")
		return shim.Error("UserAsset is less than the  Amount...")
	}

	senderlist.UserAsset -= Amount
	receiverlist.UserAsset += Amount

	_ = utils.WriteLedger(senderlist, stub, lib.UserKey, []string{SenderID})
	_ = utils.WriteLedger(receiverlist, stub, lib.UserKey, []string{ReceiverID})
	// jsontest, _ := json.Marshal(sender)
	// jsontest1, _ := json.Marshal(receiver)
	// stub.PutState(sender.UserID, jsontest)
	// stub.PutState(receiver.UserID, jsontest1)
	txdt := &lib.TransferAsset{
		TransferID:   utils.RandSeq(16),
		FromAsset:    args[0],
		ToAsset:      args[1],
		AssetValue:   Amount,
		TransferTime: time.Now().Format("2006-01-02 15:04:05"),
	}

	jsonTxDT, _ := json.Marshal(txdt)
	// stub.PutState(txdt.TransferID, jsonTxDT)
	_ = utils.WriteLedger(txdt, stub, lib.TransferAssetKey, []string{txdt.TransferID})
	return shim.Success(jsonTxDT)
}

//继承工分
func InheritAsset(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	//_ = fmt.Sprintf("InheritAsset starting ...")
	// fmt.Println("InheritAsset starting ...")
	var SenderID string   //发送方ID
	var ReceiverID string //接收方ID
	var Amount float64    //工分交易金额
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting right number of information!!!")
	}
	SenderID = args[0]
	ReceiverID = args[1]

	sender, err := utils.QueryLedger(stub, lib.UserKey, []string{SenderID})
	if err != nil {
		return shim.Error(fmt.Sprintf("发送方ID有误%s", err))
	}
	var senderlist lib.User
	_ = json.Unmarshal(sender[0], &senderlist)

	receiver, err := utils.QueryLedger(stub, lib.UserKey, []string{ReceiverID})
	if err != nil {
		return shim.Error(fmt.Sprintf("接受方ID有误%s", err))
	}
	var receiverlist lib.User
	_ = json.Unmarshal(receiver[0], &receiverlist)

	Amount = senderlist.UserAsset
	receiverlist.UserAsset += Amount
	senderlist.UserAsset -= Amount

	_ = utils.WriteLedger(senderlist, stub, lib.UserKey, []string{SenderID})
	_ = utils.WriteLedger(receiverlist, stub, lib.UserKey, []string{ReceiverID})
	//sender.UserAsset = 0
	// jsontest, _ := json.Marshal(sender)
	// jsontest1, _ := json.Marshal(receiver)
	// stub.PutState(sender.UserID, jsontest)
	// stub.PutState(receiver.UserID, jsontest1)

	txdt := &lib.TransferAsset{
		TransferID:   utils.RandSeq(16),
		FromAsset:    args[0],
		ToAsset:      args[1],
		AssetValue:   Amount,
		TransferTime: time.Now().Format("2006-01-02 15:04:05"),
	}
	jsonTxDT, _ := json.Marshal(txdt)
	// stub.PutState(txdt.TransferID, jsonTxDT)
	_ = utils.WriteLedger(txdt, stub, lib.TransferAssetKey, []string{txdt.TransferID})
	return shim.Success(jsonTxDT)
}

//充值工分
func RechargeAsset(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	//fmt.Println("RechargeAsset starting ...")

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting right number of information!!!")
	}
	var rechargerID string //充值方ID
	var Amount float64     //工分交易金额
	rechargerID = args[0]
	Amount, _ = strconv.ParseFloat(args[1], 64)

	recharger, err := utils.QueryLedger(stub, lib.UserKey, []string{rechargerID})
	if err != nil {
		return shim.Error(fmt.Sprintf("充值方ID有误%s", err))
	}
	var rechargerlist lib.User
	_ = json.Unmarshal(recharger[0], &rechargerlist)

	// if val, err := strconv.ParseFloat(amount, 64); err != nil {
	// 	return shim.Error(fmt.Sprintf("amount参数格式转换出错: %s", err))
	// } else {
	// 	Amount = val
	// }
	rechargerlist.UserAsset += Amount
	_ = utils.WriteLedger(rechargerlist, stub, lib.UserKey, []string{rechargerID})
	// jsontest, _ := json.Marshal(sender)
	// stub.PutState(sender.UserID, jsontest)

	txdt := &lib.RechargeSystem{
		RechargeID:    utils.RandSeq(16),
		ToUserID:      args[0],
		RechargeValue: Amount,
		RechargeTime:  time.Now().Format("2006-01-02 15:04:05"),
	}
	jsonTxDT, _ := json.Marshal(txdt)
	_ = utils.WriteLedger(txdt, stub, lib.RechargeSystemKey, []string{txdt.RechargeID})
	//stub.PutState(txdt.RechargeID, jsonTxDT)

	var managerID = "666666"
	manager, err := utils.QueryLedger(stub, lib.ManagerKey, []string{managerID})
	if err != nil {
		return shim.Error(fmt.Sprintf("Recharge QueryLedger manager error : %s", err))
	}
	var managerlist lib.Manager
	_ = json.Unmarshal(manager[0], &managerlist)
	managerlist.ManagerAsset += Amount * 1 //一分一元
	_ = utils.WriteLedger(managerlist, stub, lib.ManagerKey, []string{managerID})
	// if val, err := utils.QuerrManager(stub, managerID); err != nil {
	// 	return shim.Error(fmt.Sprintf("%s", err))
	// } else {
	// 	val.ManagerAsset += Amount
	// 	jsontest, _ := json.Marshal(managerID)
	// 	stub.PutState(val.ManagerID, jsontest)
	// }

	return shim.Success(jsonTxDT)
}

// //特殊交易查询
// func SpecialTradeList(stub shim.ChaincodeStubInterface, args []string) peer.Response {
// 	if len(args) != 1 {
// 		return shim.Error("Expect correct information !!!")
// 	}
// 	var specialTradeList []lib.TransferAsset
// 	results, err := utils.GetStateByPartialCompositeKeys(stub, lib.TransferAssetKey, args)
// 	if err != nil {
// 		return shim.Error(fmt.Sprintf("%s", err))
// 	}
// 	for _, val := range results {
// 		if val != nil {
// 			var specialTrade lib.TransferAsset
// 			_ = json.Unmarshal(val, &specialTrade)
// 			specialTradeList = append(specialTradeList, specialTrade)
// 		}
// 	}
// 	specialTradeListByte, _ := json.Marshal(specialTradeList)

// 	return shim.Success(specialTradeListByte)
// }

// //充值查询,可根据txid查询指定记录，也可以查询所有记录
// func RechargeList(stub shim.ChaincodeStubInterface, args []string) peer.Response {
// 	var rechargeList []lib.RechargeSystem
// 	results, err := utils.GetStateByPartialCompositeKeys(stub, lib.RechargeSystemKey, args)
// 	if err != nil {
// 		return shim.Error(fmt.Sprintf("%s", err))
// 	}
// 	for _, val := range results {
// 		if val != nil {
// 			var recharge lib.RechargeSystem
// 			_ = json.Unmarshal(val, &recharge)
// 			rechargeList = append(rechargeList, recharge)
// 		}
// 	}
// 	rechargeListByte, _ := json.Marshal(rechargeList)

// 	return shim.Success(rechargeListByte)
// }
