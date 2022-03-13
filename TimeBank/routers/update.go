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

//更新用户信息
func UpdateUserInfo(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 8 {
		return shim.Error("Expect to have the right information to update a user !!!")
	}

	olderInfo, err := utils.QueryLedger(stub, lib.UserKey, []string{args[0]})
	if err != nil {
		return shim.Error("Query error !!!")
	}
	var userlist lib.User

	_ = json.Unmarshal(olderInfo[0], &userlist)

	//results, err := utils.QueryLedger(stub, lib.UserKey, args[0])
	//userid := utils.RandSeq(11)
	username := args[1]
	useridentification := args[2]
	usersex := args[3]
	userbirthday := args[4]
	useraddress := args[5]
	userpostcode := args[6]
	userlist.Ability = append(userlist.Ability, args[7])
	//userstarsign

	// //userasset := float64(0)
	// //comment
	// //recommondid := args[7]
	CtUser := &lib.User{
		UserID:             userlist.UserID,
		UserName:           username,
		UserIdentification: useridentification,
		Sex:                usersex,
		Birthday:           userbirthday,
		Address:            useraddress,
		Postcode:           userpostcode,
		Ability:            userlist.Ability,
		UserAsset:          userlist.UserAsset,
		RecommenderID:      userlist.RecommenderID,
	}
	_ = utils.WriteLedger(CtUser, stub, lib.UserKey, []string{args[0]})
	return shim.Success([]byte(fmt.Sprintf("user %s update information successfully ...", args[0])))
}

//更新服务信息
func UpdateServiceInfo(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 4 {
		return shim.Error("Expect to have the right information to update a service !!!")
	}
	olderInfo, err := utils.QueryLedger(stub, lib.JobPriceKey, []string{args[0]})
	if err != nil {
		return shim.Error("service not exise !!!")
	}
	var servicelist lib.JobPrice
	_ = json.Unmarshal(olderInfo[0], &servicelist)
	serviceName := args[1]
	serviceUnitCost, _ := strconv.ParseFloat(args[2], 64)
	serviceDetailedDescription := args[3]

	CTService := &lib.JobPrice{
		JobID:               args[0],
		JobName:             serviceName,
		JobUnitCost:         serviceUnitCost,
		DetailedDescription: serviceDetailedDescription,
	}
	_ = utils.WriteLedger(CTService, stub, lib.JobPriceKey, []string{args[0]})
	return shim.Success([]byte(fmt.Sprintf("service %s update information successfully ...", args[0])))
}

//查询更新历史
func GetUpdateHistory(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	keyInter, err := stub.GetHistoryForKey(args[0])
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
		all := fmt.Sprintf("%s,%s,%s,%s", id, string(value), status, timeString)
		//fmt.Println(id, string(value), status, timeString)
		hisList = append(hisList, all)
	}
	jsonText, _ := json.Marshal(hisList)
	return shim.Success(jsonText)
}
