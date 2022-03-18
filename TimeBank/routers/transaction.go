package routers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/Jxnuyang/20-software-qkl/chaincode/TimeBank/lib"
	"github.com/Jxnuyang/20-software-qkl/chaincode/TimeBank/utils"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

//老人发起服务
func CreateServicing(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 3 {
		return shim.Error("Expect the right parameters !!!")
	}
	servicingType := args[0]
	servicingOlderID := args[1]
	servicingCount, _ := strconv.Atoi(args[2])
	//servicingValue, _ := strconv.ParseFloat(args[3], 64)
	if servicingType == "" || servicingOlderID == "" || servicingCount == 0 {
		return shim.Error("CreateServicing The parameter has a null value !!!")
	}
	//判断服务是否存在
	serviceRealSituation, err := utils.QueryLedger(stub, lib.JobPriceKey, []string{servicingType})
	if err != nil || len(serviceRealSituation) != 1 {
		return shim.Error("Service does not exist !!!")
	}
	var job lib.JobPrice
	json.Unmarshal(serviceRealSituation[0], &job)

	servicingValue := job.JobUnitCost * float64(servicingCount)

	olderInfo, err := utils.QueryLedger(stub, lib.UserKey, []string{servicingOlderID})
	if err != nil || len(olderInfo) != 1 {
		return shim.Error("Failed to obtain volunteer information !!!")
	}
	var older lib.User
	_ = json.Unmarshal(olderInfo[0], &older)
	//判断老人工分是否足够
	if older.UserAsset < servicingValue {
		return shim.Error("The balance is not enough !!!")
	}

	older.UserAsset -= servicingValue
	_ = utils.WriteLedger(older, stub, lib.UserKey, []string{servicingOlderID})

	servicing := &lib.Servicing{
		ServicingType:    servicingType,
		ServicingOlderID: servicingOlderID,
		ServicingVolID:   "",
		ServicingCount:   servicingCount,
		StartTime:        time.Now().Local().Format("2006-01-02 15:04:05"),
		ServicingValue:   servicingValue,
		ServicingState:   lib.ServiceTradingStatusConstant()["require"],
	}

	err = utils.WriteLedger(servicing, stub, lib.ServicingKey, []string{servicingOlderID, servicingType})
	if err != nil {
		return shim.Error("CreateServicing Writeledger error !!!")
	}

	servicingByte, err := json.Marshal(servicing)
	if err != nil {
		return shim.Error("CreateServicing Marshal error !!!")
	}
	return shim.Success(servicingByte)
}

//志愿者接受服务
func AcceptServicing(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 3 {
		return shim.Error("Expect the right parameters !!!")
	}
	acceptServicingType := args[0]
	accepterID := args[1]
	olderID := args[2]

	if acceptServicingType == "" || accepterID == "" || olderID == "" {
		return shim.Error("AcceptServicing The parameter has a null value !!!")
	}

	if accepterID == olderID {
		return shim.Error("The volunteers and the elderly can't be the same person !!!")
	}

	resultsServicing, err := utils.GetStateByPartialCompositeKeys(stub, lib.ServicingKey, []string{olderID, acceptServicingType})
	if err != nil || len(resultsServicing) != 1 {
		return shim.Error("Failed to obtain information !!!")
	}
	var servicing lib.Servicing
	_ = json.Unmarshal(resultsServicing[0], &servicing)

	if servicing.ServicingState != lib.ServiceTradingStatusConstant()["require"] {
		return shim.Error("The service is not requested !!!")
	}

	olderInfo, err := utils.QueryLedger(stub, lib.UserKey, []string{olderID})
	if err != nil || len(olderInfo) != 1 {
		return shim.Error("Failed to obtain volunteer information !!!")
	}
	var older lib.User
	_ = json.Unmarshal(olderInfo[0], &older)

	// if older.UserAsset < servicing.ServicingValue {
	// 	return shim.Error("The balance is not enough !!!")
	// }

	servicing.ServicingVolID = accepterID
	servicing.ServicingState = lib.ServiceTradingStatusConstant()["accepted"]

	_ = utils.WriteLedger(servicing, stub, lib.ServicingKey, []string{olderID, acceptServicingType})

	servicingByte, err := json.Marshal(servicing)
	if err != nil {
		return shim.Error("Accept Marshal error !!!")
	}
	// serviceTrade := &lib.ServiceTrade{
	// 	TxID:      utils.RandSeq(16),
	// 	ServeID:   accepterID,
	// 	CustID:    olderID,
	// 	TxType:    acceptServicingType,
	// 	WorkHours: servicing.ServicingCount,
	// 	WorkValue: servicing.ServicingValue,
	// }
	// utils.WriteLedger(serviceTrade, stub, lib.ServiceTradeKey, []string{serviceTrade.TxID, olderID, accepterID})

	// serviceTradeByte, err := json.Marshal(serviceTrade)
	// if err != nil {
	// 	return shim.Error("Accept Marshal error !!!")
	// }

	// older.UserAsset -= servicing.ServicingValue
	// utils.WriteLedger(older, stub, lib.UserKey, []string{olderID})

	return shim.Success(servicingByte)
}

//结束服务状态
func DoneServicing(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 4 {
		return shim.Error("Expect correct information input !!!")
	}

	serviceType := args[0]
	elder := args[1]
	accepter := args[2]
	comment := args[3]

	if serviceType == "" || elder == "" {
		return shim.Error("The parameter has a null value !!!")
	}

	if elder == accepter {
		return shim.Error("Volunteers and elderly people cannot be the same person !!!")
	}

	resultsServicing, err := utils.GetStateByPartialCompositeKeys(stub, lib.ServicingKey, []string{elder, serviceType})
	if err != nil || len(resultsServicing) != 1 {
		return shim.Error(fmt.Sprintf("According to %s and %s are error !!!", elder, serviceType))
	}
	var servicing lib.Servicing
	_ = json.Unmarshal(resultsServicing[0], &servicing)

	resultsVolunteer, err := utils.GetStateByPartialCompositeKeys(stub, lib.UserKey, []string{accepter})
	if err != nil {
		return shim.Error("Volunteer query failed !!!")
	}
	var volunteer lib.User
	_ = json.Unmarshal(resultsVolunteer[0], &volunteer)
	// if servicing.ServicingState != lib.ServiceTradingStatusConstant()["require"] {
	// 	return shim.Error("The service has not been accepted !!!")
	// }
	volunteer.UserAsset += servicing.ServicingValue
	volunteer.Comment = append(volunteer.Comment, comment)
	_ = utils.WriteLedger(volunteer, stub, lib.UserKey, []string{accepter})

	servicing.ServicingState = lib.ServiceTradingStatusConstant()["done"]
	_ = utils.WriteLedger(servicing, stub, lib.ServicingKey, []string{elder, serviceType})

	serviceTrade := &lib.ServiceTrade{
		TxID:         utils.RandSeq(16),
		ServeID:      accepter,
		CustID:       elder,
		TxType:       serviceType,
		WorkHours:    servicing.ServicingCount,
		EndTime:      time.Now().Local().Format("2006-01-02 15:04:05"),
		WorkValue:    servicing.ServicingValue,
		ServeComment: comment,
	}

	_ = utils.WriteLedger(serviceTrade, stub, lib.ServiceTradeKey, []string{serviceTrade.TxID, elder, accepter})
	serviceTradeByte, err := json.Marshal(serviceTrade)
	if err != nil {
		return shim.Error("Done Marshal error !!!")
	}

	return shim.Success(serviceTradeByte)
}

//取消服务
func CloseServicing(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("Expect correct information input !!!")
	}

	servicingType := args[0]
	elderID := args[1]

	servicingList, err := utils.GetStateByPartialCompositeKeys(stub, lib.ServicingKey, []string{elderID, servicingType})
	if err != nil || len(servicingList) != 1 {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	var servicing lib.Servicing
	_ = json.Unmarshal(servicingList[0], &servicing)

	elderList, err := utils.QueryLedger(stub, lib.UserKey, []string{elderID})
	if err != nil || len(elderList) != 1 {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	var elder lib.User
	_ = json.Unmarshal(elderList[0], &elder)

	servicing.ServicingState = lib.ServiceTradingStatusConstant()["cancelled"]
	elder.UserAsset += servicing.ServicingValue

	_ = utils.WriteLedger(servicing, stub, lib.ServicingKey, []string{elderID, servicingType})
	_ = utils.WriteLedger(elder, stub, lib.UserKey, []string{elderID})

	return shim.Success([]byte("Service Cancelled successfully ..."))
}

// //供老人查看服务状态，输入老人id以及服务类型id，即可查询服务状态
// func QueryServicingStatus(stub shim.ChaincodeStubInterface, args []string) peer.Response {
// 	if len(args) == 0 {
// 		return shim.Error("Expect correct information !!!")
// 	}
// 	var serviceList []lib.Servicing
// 	results, err := utils.GetStateByPartialCompositeKeys(stub, lib.ServicingKey, args)
// 	if err != nil {
// 		return shim.Error(fmt.Sprintf("%s", err))
// 	}
// 	for _, v := range results {
// 		if v != nil {
// 			var servicing lib.Servicing
// 			_ = json.Unmarshal(v, &servicing)
// 			serviceList = append(serviceList, servicing)
// 		}
// 	}
// 	serviceListByte, err := json.Marshal(serviceList)
// 	if err != nil {
// 		return shim.Error("QueryServiceList Marshal error !!!")
// 	}
// 	return shim.Success(serviceListByte)
// }

// //查询服务记录，输入交易id、交易双方id和交易id，即可查询两人间的所有服务记录
// func QueryServiceTrade(stub shim.ChaincodeStubInterface, args []string) peer.Response {
// 	if len(args) == 0 {
// 		return shim.Error("Expect correct information !!!")
// 	}
// 	var serviceTradeList []lib.ServiceTrade
// 	results, err := utils.GetStateByPartialCompositeKeys(stub, lib.ServiceTradeKey, args)
// 	if err != nil {
// 		return shim.Error(fmt.Sprintf("%s", err))
// 	}
// 	for _, val := range results {
// 		if val != nil {
// 			var serviceTrade lib.ServiceTrade
// 			_ = json.Unmarshal(val, &serviceTrade)
// 			serviceTradeList = append(serviceTradeList, serviceTrade)
// 		}
// 	}
// 	serviceTradeListByte, _ := json.Marshal(serviceTradeList)

// 	return shim.Success(serviceTradeListByte)
// }
