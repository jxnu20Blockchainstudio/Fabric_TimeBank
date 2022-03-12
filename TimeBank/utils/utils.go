package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"

	"github.com/hyperledger/fabric-chaincode-go/shim"
)

var letters = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
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
