package routers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Jxnuyang/20-software-qkl/chaincode/TimeBank/lib"
	"github.com/Jxnuyang/20-software-qkl/chaincode/TimeBank/utils"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

//打印用户信息
func UserList(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var printuserlist []lib.User
	results, err := utils.QueryLedger(stub, lib.UserKey, args)
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	for _, v := range results {
		if v != nil {
			var printuser lib.User
			err := json.Unmarshal(v, &printuser)
			if err != nil {
				return shim.Error(fmt.Sprintf("UserList unmarshal error : %s", err))
			}
			printuserlist = append(printuserlist, printuser)
		}
	}
	printuserbyte, err := json.Marshal(printuserlist)
	if err != nil {
		return shim.Error(fmt.Sprintf("UserList Marshal error : %s", err))
	}
	return shim.Success(printuserbyte)
}

//打印组织信息
func OrgList(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var printorglist []lib.Organization
	results, err := utils.QueryLedger(stub, lib.OrganizationKey, args)
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	for _, v := range results {
		if v != nil {
			var printorg lib.Organization
			err := json.Unmarshal(v, &printorg)
			if err != nil {
				return shim.Error((fmt.Sprintf("OrgList unmashal error : %s", err)))
			}
			printorglist = append(printorglist, printorg)
		}
	}
	printorgbyte, err := json.Marshal(printorglist)
	if err != nil {
		return shim.Error(fmt.Sprintf("OrgList marshal error : %s", err))
	}
	return shim.Success(printorgbyte)
}

//打印管理员信息
func ManagerList(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var printmanagerlsit []lib.Manager
	results, err := utils.QueryLedger(stub, lib.ManagerKey, args)
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	for _, v := range results {
		if v != nil {
			var printmanager lib.Manager
			err := json.Unmarshal(v, &printmanager)
			if err != nil {
				return shim.Error((fmt.Sprintf("ManagerList unmashal error : %s", err)))
			}
			printmanagerlsit = append(printmanagerlsit, printmanager)
		}
	}
	printmanagerbyte, err := json.Marshal(printmanagerlsit)
	if err != nil {
		return shim.Error(fmt.Sprintf("ManagerList marshal error : %s", err))
	}
	return shim.Success(printmanagerbyte)
}

//打印所有可选择服务
func ServiceList(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var printservicelist []lib.JobPrice
	results, err := utils.QueryLedger(stub, lib.JobPriceKey, args)
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	for _, v := range results {
		if v != nil {
			var printservice lib.JobPrice
			err := json.Unmarshal(v, &printservice)
			if err != nil {
				return shim.Error((fmt.Sprintf("ServiceList unmashal error : %s", err)))
			}
			printservicelist = append(printservicelist, printservice)
		}
	}
	printservicebyte, err := json.Marshal(printservicelist)
	if err != nil {
		return shim.Error(fmt.Sprintf("ServiceList marshal error : %s", err))
	}
	return shim.Success(printservicebyte)
}

//打印服务状态，可打印所有服务状态以及指定服务状态
//输入老人id以及服务类型id，即可查询指定服务状态
func QueryServicingStatus(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	// if len(args) == 0 {
	// 	return shim.Error("Expect correct information !!!")
	// }
	var serviceList []lib.Servicing
	results, err := utils.GetStateByPartialCompositeKeys(stub, lib.ServicingKey, args)
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	for _, v := range results {
		if v != nil {
			var servicing lib.Servicing
			_ = json.Unmarshal(v, &servicing)
			serviceList = append(serviceList, servicing)
		}
	}
	serviceListByte, err := json.Marshal(serviceList)
	if err != nil {
		return shim.Error("QueryServiceList Marshal error !!!")
	}
	return shim.Success(serviceListByte)
}

//打印服务记录，可打印所有服务记录以及指定服务记录
//输入交易id、交易双方id和交易id，即可查询两人间的所有服务记录
func QueryServiceTrade(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// if len(args) == 0 {
	// 	return shim.Error("Expect correct information !!!")
	// }
	var serviceTradeList []lib.ServiceTrade
	results, err := utils.GetStateByPartialCompositeKeys(stub, lib.ServiceTradeKey, args)
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	for _, val := range results {
		if val != nil {
			var serviceTrade lib.ServiceTrade
			_ = json.Unmarshal(val, &serviceTrade)
			serviceTradeList = append(serviceTradeList, serviceTrade)
		}
	}
	serviceTradeListByte, _ := json.Marshal(serviceTradeList)

	return shim.Success(serviceTradeListByte)
}

//特殊交易查询
func SpecialTradeList(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Expect correct information !!!")
	}
	var specialTradeList []lib.TransferAsset
	results, err := utils.GetStateByPartialCompositeKeys(stub, lib.TransferAssetKey, args)
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	for _, val := range results {
		if val != nil {
			var specialTrade lib.TransferAsset
			_ = json.Unmarshal(val, &specialTrade)
			specialTradeList = append(specialTradeList, specialTrade)
		}
	}
	specialTradeListByte, _ := json.Marshal(specialTradeList)

	return shim.Success(specialTradeListByte)
}

//充值查询,可根据txid查询指定记录，也可以查询所有记录
func RechargeList(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var rechargeList []lib.RechargeSystem
	results, err := utils.GetStateByPartialCompositeKeys(stub, lib.RechargeSystemKey, args)
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	for _, val := range results {
		if val != nil {
			var recharge lib.RechargeSystem
			_ = json.Unmarshal(val, &recharge)
			rechargeList = append(rechargeList, recharge)
		}
	}
	rechargeListByte, _ := json.Marshal(rechargeList)

	return shim.Success(rechargeListByte)
}

//打印更新历史，输入名称以及id
func GetUpdateHistory(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 2 {
		return shim.Error("Expect correct information input !!!")
	}

	var UpdateValue string
	UpdateName := args[0]
	UpdateID := args[1]

	if UpdateName == "user" {
		if v, err := stub.CreateCompositeKey(lib.UserKey, []string{UpdateID}); err != nil {
			return shim.Error("GetUpdateHistory CreateCompositeKey error !!!")
		} else {
			UpdateValue = v
		}
	}

	if UpdateName == "service" {
		if v, err := stub.CreateCompositeKey(lib.JobPriceKey, []string{UpdateID}); err != nil {
			return shim.Error("GetUpdateHistory CreateCompositeKey error !!!")
		} else {
			UpdateValue = v
		}
	}

	keyInter, err := stub.GetHistoryForKey(UpdateValue)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer keyInter.Close()
	var hisList []string
	for keyInter.HasNext() {
		response, interErr := keyInter.Next()
		if interErr != nil {
			return shim.Error(interErr.Error())
		}
		id := response.TxId
		value := response.Value
		status := response.IsDelete
		timeStamp := response.Timestamp
		tm := time.Unix(timeStamp.Seconds, 0)
		timeString := tm.Format("2006-01-02 03:04:05 PM")
		all := fmt.Sprintf("%s,%s,%t,%s", id, string(value), status, timeString)
		//fmt.Println(id, string(value), status, timeString)
		hisList = append(hisList, all)
	}
	jsonText, _ := json.Marshal(hisList)
	return shim.Success(jsonText)
}

// //打印所有老人的正在交易的详细信息
// func ElderServicingList(stub shim.ChaincodeStubInterface, args []string) peer.Response {
// 	var printelderservicinglist []lib.Servicing
// 	results, err := utils.QueryLedger(stub, lib.ServicingKey, args)
// 	if err != nil {
// 		return shim.Error(fmt.Sprintf("%s", err))
// 	}
// 	for _, v := range results {
// 		if v != nil {
// 			var printelderservicing lib.Servicing
// 			err := json.Unmarshal(v, &printelderservicing)
// 			if err != nil {
// 				return shim.Error((fmt.Sprintf("ElderServiceList unmashal error : %s", err)))
// 			}
// 			printelderservicinglist = append(printelderservicinglist, printelderservicing)
// 		}
// 	}
// 	printelderservicingbyte, err := json.Marshal(printelderservicinglist)
// 	if err != nil {
// 		return shim.Error(fmt.Sprintf("ElderServiceList marshal error : %s", err))
// 	}
// 	return shim.Success(printelderservicingbyte)
// }
