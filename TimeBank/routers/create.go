package routers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Jxnuyang/20-software-qkl/chaincode/TimeBank/lib"
	"github.com/Jxnuyang/20-software-qkl/chaincode/TimeBank/utils"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

//创建用户
func CreateUser(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	//_ = fmt.Sprintf("User %s is creating ...", args[0])

	if len(args) != 8 {
		return shim.Error("Expect to have the right information to create a user !!!")
	}
	var userability []string
	userid := utils.RandSeq(11)
	username := args[0]
	useridentification := args[1]
	usersex := args[2]
	userbirthday := args[3]
	useraddress := args[4]
	userpostcode := args[5]
	userability = append(userability, args[6])
	//userstarsign

	//userasset := float64(0)
	//comment
	recommondid := args[7]
	CtUser := &lib.User{
		UserID:             userid,
		UserName:           username,
		UserIdentification: useridentification,
		Sex:                usersex,
		Birthday:           userbirthday,
		Address:            useraddress,
		Postcode:           userpostcode,
		Ability:            userability,
		UserAsset:          float64(0),
		RecommenderID:      recommondid,
	}
	//推荐有奖
	//count:=0
	if recommondid != "" {
		//直接推荐人
		//count++
		recommonduser1, err := utils.QueryLedger(stub, lib.UserKey, []string{recommondid})
		if err != nil {
			return shim.Error(fmt.Sprintf("CreateUser QueryLedger error : %s", err))
		}
		var recommonduser1list lib.User
		json.Unmarshal(recommonduser1[0], &recommonduser1list)
		recommonduser1list.UserAsset = recommonduser1list.UserAsset + 50
		_ = utils.WriteLedger(recommonduser1list, stub, lib.UserKey, []string{recommonduser1list.UserID})
		recommondid1 := recommonduser1list.RecommenderID
		//上级推荐人
		if recommondid1 != "" {
			//count++
			recommonduser2, err := utils.QueryLedger(stub, lib.UserKey, []string{recommondid1})
			if err != nil {
				return shim.Error(fmt.Sprintf("CreateUser QueryLedger error : %s", err))
			}
			var recommonduser2list lib.User
			json.Unmarshal(recommonduser2[0], &recommonduser2list)
			recommonduser2list.UserAsset = recommonduser2list.UserAsset + 30
			_ = utils.WriteLedger(recommonduser2list, stub, lib.UserKey, []string{recommonduser2list.UserID})
			recommondid2 := recommonduser2list.RecommenderID
			//上上级推荐人
			if recommondid2 != "" {
				//count++
				recommonduser3, err := utils.QueryLedger(stub, lib.UserKey, []string{recommondid2})
				if err != nil {
					return shim.Error(fmt.Sprintf("CreateUser QueryLedger error : %s", err))
				}
				var recommonduser3list lib.User
				json.Unmarshal(recommonduser3[0], &recommonduser3list)
				recommonduser3list.UserAsset = recommonduser3list.UserAsset + 20
				_ = utils.WriteLedger(recommonduser3list, stub, lib.UserKey, []string{recommonduser3list.UserID})
			}
		}
	}

	_ = utils.WriteLedger(CtUser, stub, lib.UserKey, []string{userid})

	//将新用户加入所属组织
	org, err := utils.QueryLedger(stub, lib.OrganizationKey, []string{CtUser.Postcode})
	if err != nil {
		return shim.Error("Failed to find the organization !!!")
	}
	var orgList lib.Organization
	_ = json.Unmarshal(org[0], &orgList)

	orgList.UserSum++
	orgList.HaveUserID = append(orgList.HaveUserID, CtUser.UserID)

	_ = utils.WriteLedger(orgList, stub, lib.OrganizationKey, []string{CtUser.Postcode})
	return shim.Success([]byte(fmt.Sprintf("User %s was Created successfully ...", userid)))
}

//创建服务
func CreateService(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	//_ = fmt.Sprintf("%s service is creating ...", args[1])
	if len(args) != 3 {
		return shim.Error("Expect to have the right information to create a service !!!")
	}

	serid := utils.RandSeq(6)
	sername := args[0]
	serunitcost, _ := strconv.ParseFloat(args[1], 64)
	serdescription := args[2]

	service := &lib.JobPrice{
		JobID:               serid,
		JobName:             sername,
		JobUnitCost:         serunitcost,
		DetailedDescription: serdescription,
	}

	_ = utils.WriteLedger(service, stub, lib.JobPriceKey, []string{serid})

	return shim.Success([]byte(fmt.Sprintf("%s service onlysign %s was created successfully", args[0], serid)))
}

//创建组织
func CreateOrg(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	//_ = fmt.Sprintf("Org %s is creating ...", args[0])
	if len(args) != 2 {
		return shim.Error("Expect to have the right information to create a org !!!")
	}

	orgid := args[0]
	orgname := args[1]
	orgnum := 0

	org := &lib.Organization{
		OrgID:      orgid,
		OrgName:    orgname,
		UserSum:    orgnum,
		HaveUserID: []string{},
	}

	_ = utils.WriteLedger(org, stub, lib.OrganizationKey, []string{orgid})

	return shim.Success([]byte(fmt.Sprintf("%s org postcode %s was created successfully", args[1], args[0])))
}
