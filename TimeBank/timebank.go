package main

import (
	"fmt"
	"time"

	"github.com/Jxnuyang/20-software-qkl/chaincode/TimeBank/lib"
	"github.com/Jxnuyang/20-software-qkl/chaincode/TimeBank/routers"
	"github.com/Jxnuyang/20-software-qkl/chaincode/TimeBank/utils"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

type BlockChainTimeBank struct {
}

// Init 链码初始化
func (t *BlockChainTimeBank) Init(stub shim.ChaincodeStubInterface) peer.Response {
	//fmt.Sprintf("Init chaincode start ...")
	timeLocal, err := time.LoadLocation("Asia/Chongqing")
	if err != nil {
		return shim.Error(fmt.Sprintf("时区设置失败%s", err))
	}
	time.Local = timeLocal
	//初始化默认数据
	var managerID = "666666"
	var managerAsset float64 = 100
	manager := &lib.Manager{
		ManagerID:    managerID,
		ManagerAsset: managerAsset,
	}
	// 写入账本
	utils.WriteLedger(manager, stub, lib.ManagerKey, []string{managerID})
	var orgids = [3]string{
		"330027", //青山湖校区邮政编码
		"330022", //瑶湖校区邮政编码
		"332020", //共青城校区邮政编码
	}
	var orgNames = [3]string{"青山湖校区", "瑶湖校区", "共青城校区"}
	var userNums = [3]int{0, 0, 0}
	//var empty []string
	//初始化账号数据
	for i, val := range orgids {
		oragnization := &lib.Organization{
			OrgID:      val,
			OrgName:    orgNames[i],
			UserSum:    userNums[i],
			HaveUserID: []string{},
		}
		// 写入账本
		utils.WriteLedger(oragnization, stub, lib.OrganizationKey, []string{val})
	}

	//杨大爷
	userwang := &lib.User{
		UserID:             "01",
		UserName:           "杨大爷",
		UserIdentification: "110101199003076675",
		Sex:                "男",
		Birthday:           "1990-3-7",
		Address:            "江西师范大学青山湖校区",
		Postcode:           "330027",
		Ability:            []string{"做饭", "教学高数"},
		StarSign:           3,
		UserAsset:          float64(1000),
		Comment:            []string{"很细心", "Good!!!", "杨老师教的高数就是好"},
		RecommenderID:      "",
	}
	// var orgList1 lib.Organization
	// org1, err := utils.QueryLedger(stub, lib.OrganizationKey, []string{userwang.Postcode})
	// if err != nil {
	// 	return shim.Error("Failed to find the organization !!!")
	// }
	// _ = json.Unmarshal(org1[0], &orgList1)

	// orgList1.UserSum++
	// orgList1.HaveUserID = append(orgList1.HaveUserID, userwang.UserID)

	// _ = utils.WriteLedger(orgList1, stub, lib.OrganizationKey, []string{userwang.Postcode})
	//吕婆婆

	// var orgList2 lib.Organization
	userlv := &lib.User{
		UserID:             "02",
		UserName:           "吕婆婆",
		UserIdentification: "110101199006068985",
		Sex:                "女",
		Birthday:           "1990-6-6",
		Address:            "江西师范大学瑶湖校区",
		Postcode:           "330022",
		Ability:            []string{"做饭", "洗衣服", "教学化学"},
		StarSign:           5,
		UserAsset:          float64(6000),
		Comment:            []string{"很细心", "Very Good!!!", "太可爱了"},
		RecommenderID:      "",
	}
	// org2, err := utils.QueryLedger(stub, lib.OrganizationKey, []string{userlv.Postcode})
	// if err != nil {
	// 	return shim.Error("Failed to find the organization !!!")
	// }
	// //var orgList lib.Organization
	// _ = json.Unmarshal(org2[0], &orgList2)

	// orgList2.UserSum++
	// orgList2.HaveUserID = append(orgList2.HaveUserID, userlv.UserID)

	// _ = utils.WriteLedger(orgList2, stub, lib.OrganizationKey, []string{userwang.Postcode})
	_ = utils.WriteLedger(userwang, stub, lib.UserKey, []string{userwang.UserID})
	_ = utils.WriteLedger(userlv, stub, lib.UserKey, []string{userlv.UserID})
	return shim.Success(nil)
}

// Invoke 实现Invoke接口调用智能合约
func (t *BlockChainTimeBank) Invoke(stub shim.ChaincodeStubInterface) peer.Response {

	funcName, args := stub.GetFunctionAndParameters()
	//fmt.Sprintf("Invoke %s start ...", funcName)
	if funcName == "CreateUser" { //创建用户
		return routers.CreateUser(stub, args)
	} else if funcName == "CreateService" { //创建服务
		return routers.CreateService(stub, args)
	} else if funcName == "CreateOrg" { //创建组织
		return routers.CreateOrg(stub, args)
	} else if funcName == "UserList" { //打印用户信息
		return routers.UserList(stub, args)
	} else if funcName == "OrgList" { //打印组织信息
		return routers.OrgList(stub, args)
	} else if funcName == "ManagerList" { //打印管理员信息
		return routers.ManagerList(stub, args)
	} else if funcName == "CreateServicing" { //发起服务
		return routers.CreateServicing(stub, args)
	} else if funcName == "AcceptServicing" { //接受服务
		return routers.AcceptServicing(stub, args)
	} else if funcName == "DoneServicing" { //结束服务
		return routers.DoneServicing(stub, args)
	} else if funcName == "CloseServicing" { //取消服务
		return routers.CloseServicing(stub, args)
	} else if funcName == "QueryServicingStatus" { //查询订单状态
		return routers.QueryServicingStatus(stub, args)
	} else if funcName == "QueryServiceTrade" { //查询服务记录
		return routers.QueryServiceTrade(stub, args)
	} else if funcName == "TransferAsset" { //转移资产
		return routers.TransferAsset(stub, args)
	} else if funcName == "InheritAsset" { //继承资产
		return routers.InheritAsset(stub, args)
	} else if funcName == "RechargeAsset" { //充值资产
		return routers.RechargeAsset(stub, args)
	} else if funcName == "SpecialTradeList" { //打印特殊交易
		return routers.SpecialTradeList(stub, args)
	} else if funcName == "ElderServicingList" { //打印所有老人服务状态
		return routers.ElderServicingList(stub, args)
	} else if funcName == "ServiceList" { //打印可提供的服务
		return routers.ServiceList(stub, args)
	} else if funcName == "UpdateUserInfo" {
		return routers.UpdateUserInfo(stub, args)
	} else if funcName == "UpdateServiceInfo" {
		return routers.UpdateServiceInfo(stub, args)
	} else if funcName == "GetUpdateHistory" {
		return routers.GetUpdateHistory(stub, args)
	} else if funcName == "RechargeList" {
		return routers.RechargeList(stub, args)
	} else {
		return shim.Error("Invoke funcName error !!!")
	}
	//return shim.Success([]byte(""))
}

//启动并进入链码
func main() {
	err := shim.Start(new(BlockChainTimeBank))
	if err != nil {
		fmt.Printf("Chaincode start error %s", err)
	}
}
