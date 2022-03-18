package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
)

func RandSqe(lenNum int) string {
	var chars = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}
	str := strings.Builder{}
	length := len(chars)
	rand.Seed(time.Now().UnixNano()) //重新播种，否则值不会变
	for i := 0; i < lenNum; i++ {
		str.WriteString(chars[rand.Intn(length)])
	}
	return str.String()
}

func WriteLedger(obj interface{}, stub shim.ChaincodeStubInterface, objectType string, args []string) error {

	_ = fmt.Sprintf("%s start WriteLedger ...", objectType)

	var key string
	if val, err := stub.CreateCompositeKey(objectType, args); err != nil {
		return errors.New(fmt.Sprintf("WriteLedger CreateCompositeKey error !!!"))
	} else {
		key = val
	}

	bytes, err := json.Marshal(obj)
	if err != nil {
		return errors.New(fmt.Sprintf("WriteLedger Marshal error!!!"))
	}
	stub.PutState(key, bytes)

	return nil

}

func DeleteLedger(stub shim.ChaincodeStubInterface, objectType string, args []string) error {

	_ = fmt.Sprintf("%s start DeleteLedger ...", objectType)

	var key string

	if val, err := stub.CreateCompositeKey(objectType, args); err != nil {
		return errors.New(fmt.Sprintf("DeleteLedger CreateCompositeKey error !!!"))
	} else {
		key = val
	}

	if err := stub.DelState(key); err != nil {
		return errors.New(fmt.Sprintf("DeleteLedger DeleteState error !!!"))
	}

	return nil
}

//指定查询
func QueryLedger(stub shim.ChaincodeStubInterface, objectType string, args []string) ([][]byte, error) {

	_ = fmt.Sprintf("%s start QueryLedger ...", objectType)

	var results [][]byte
	if len(args) == 0 {

		val, err := stub.GetStateByPartialCompositeKey(objectType, args)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("QueryLedger GetStateByPartialCompositeKey error !!!"))
		}
		defer val.Close()

		for val.HasNext() {
			value, err := val.Next()
			if err != nil {
				return nil, errors.New(fmt.Sprintf("QueryLegder ergodic error !!!"))
			}

			results = append(results, value.GetValue())
		}
	} else {
		for _, v := range args {
			key, err := stub.CreateCompositeKey(objectType, []string{v})
			if err != nil {
				return nil, errors.New(fmt.Sprintf("QueryLedger CreateCompositeKey error !!!"))
			}

			bytes, err := stub.GetState(key)

			if err != nil {
				return nil, errors.New(fmt.Sprintf("QueryLedger GetState error !!!"))
			}
			if bytes != nil {
				results = append(results, bytes)
			}
		}
	}
	return results, nil
}

//模糊查询
func GetStateByPartialCompositeKeys(stub shim.ChaincodeStubInterface, objectType string, args []string) ([][]byte, error) {
	resultIterator, err := stub.GetStateByPartialCompositeKey(objectType, args)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s-Error obtaining data: %s", objectType, err))
	}
	defer resultIterator.Close()

	var results [][]byte
	for resultIterator.HasNext() {
		val, err := resultIterator.Next()
		if err != nil {
			return nil, errors.New(fmt.Sprintf("%s-Error returning data: %s", objectType, err))
		}

		results = append(results, val.GetValue())
	}
	return results, nil
}
