package routers

import (
	"encoding/json"
	"fmt"

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

//打印所有老人的正在交易的详细信息
func ElderServicingList(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var printelderservicinglist []lib.Servicing
	results, err := utils.QueryLedger(stub, lib.ServicingKey, args)
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	for _, v := range results {
		if v != nil {
			var printelderservicing lib.Servicing
			err := json.Unmarshal(v, &printelderservicing)
			if err != nil {
				return shim.Error((fmt.Sprintf("ElderServiceList unmashal error : %s", err)))
			}
			printelderservicinglist = append(printelderservicinglist, printelderservicing)
		}
	}
	printelderservicingbyte, err := json.Marshal(printelderservicinglist)
	if err != nil {
		return shim.Error(fmt.Sprintf("ElderServiceList marshal error : %s", err))
	}
	return shim.Success(printelderservicingbyte)
}
