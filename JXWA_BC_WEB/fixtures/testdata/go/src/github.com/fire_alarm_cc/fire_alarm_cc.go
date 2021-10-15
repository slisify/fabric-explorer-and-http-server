/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type FireAlarmData struct {
	AlertorID      string `json:"alertorID"`        //传感器ID
	AlertTime      string `json:"alertTime"`      //发出警报的时间
	HouseNumber    string `json:"houseNumber"`    //住宅编号或者房建号
	HouseOwnerName string `json:"houseOwnerName"` //住户的名字
}

// Init ...
// func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
// 	txID := stub.GetTxID()
// 	_, args := stub.GetFunctionAndParameters()

// 	err := t.reset(stub, txID, args)
// 	if err != nil {
// 		return shim.Error(err.Error())
// 	}

// 	if transientMap, err := stub.GetTransient(); err == nil {
// 		if transientData, ok := transientMap["result"]; ok {
// 			return shim.Success(transientData)
// 		}
// 	}
// 	return shim.Success(nil)

// }

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	//不需要初始化什么具体内容
	return shim.Success(nil)
}

// func (t *SimpleChaincode) resetCC(stub shim.ChaincodeStubInterface, args []string) pb.Response {
// 	// Deletes an entity from its state
// 	if err := t.reset(stub, stub.GetTxID(), args); err != nil {
// 		return shim.Error(err.Error())
// 	}
// 	return shim.Success(nil)
// }

//暂时不需要reset的相关操作
// func (t *SimpleChaincode) reset(stub shim.ChaincodeStubInterface, txID string, args []string) error {
// 	var A, B string    // Entities
// 	var Aval, Bval int // Asset holdings
// 	var err error

// 	if len(args) != 4 {
// 		return errors.New("Incorrect number of arguments. Expecting 4")
// 	}

// 	// Initialize the chaincode
// 	A = args[0]
// 	Aval, err = strconv.Atoi(args[1])
// 	if err != nil {
// 		return errors.New("Expecting integer value for asset holding")
// 	}
// 	B = args[2]
// 	Bval, err = strconv.Atoi(args[3])
// 	if err != nil {
// 		return errors.New("Expecting integer value for asset holding")
// 	}

// 	// Write the state to the ledger
// 	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
// 	if err != nil {
// 		return err
// 	}

// 	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// Query ...
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Error("Unknown supported call")
}

//set sets given key-value in state
func (t *SimpleChaincode) set(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	if len(args) < 5 {
		return shim.Error("Incorrect number of arguments. Expecting a key and a value")
	}

	// Initialize the chaincode

	//TODO:该用什么ID
	eventID := "testEvent"
	if len(args) >= 6 {
		eventID = args[5]
	}

	alertorID := args[1]

	fireAlarmData := FireAlarmData{
		AlertorID:      args[1],
		AlertTime:      args[2],
		HouseNumber:    args[3],
		HouseOwnerName: args[4],
	}

	// Write the state to the ledger
	fireAlarmDataJson, err := json.Marshal(fireAlarmData)
	err = stub.PutState(alertorID, []byte(fireAlarmDataJson))
	if err != nil {
		fmt.Printf("Failed to set value for key[%s] : ", alertorID, err)
		return shim.Error(err.Error())
	}

	err = stub.SetEvent(eventID, []byte("Test Payload"))
	if err != nil {
		fmt.Printf("Failed to set event for key[%s] : ", alertorID, err)
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// Invoke ...
// Transaction makes payment of X units from A to B
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()

	if function == "invokecc" {
		return t.invokeCC(stub, args)
	}

	//目前不提供重置功能
	// if function == "reset" {
	// 	return t.resetCC(stub, args)
	// }

	if function != "invoke" {
		return shim.Error("Unknown function call")
	}

	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting at least 2")
	}

	if args[0] == "delete" {
		// Deletes an entity from its state
		return t.delete(stub, args)
	}

	if args[0] == "query" {
		// queries an entity state
		return t.query(stub, args)
	}

	if args[0] == "set" {
		// setting an entity state
		return t.set(stub, args)
	}

	// if args[0] == "move" {
	// 	eventID := "testEvent"
	// 	if len(args) >= 5 {
	// 		eventID = args[4]
	// 	}
	// 	if err := stub.SetEvent(eventID, []byte("Test Payload")); err != nil {
	// 		return shim.Error("Unable to set CC event: testEvent. Aborting transaction ...")
	// 	}
	// 	return t.move(stub, args)
	// }

	return shim.Error("Unknown action, check the first argument, must be one of 'delete', 'query', or 'move'")
}


// Deletes an entity from state
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	AlertorID := args[1]

	// Delete the key from the state in ledger
	err := stub.DelState(AlertorID)
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	return shim.Success(nil)
}

// Query callback representing the query of a chaincode
func (t *SimpleChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var AlertorID string // Entities
	var err error

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	AlertorID = args[1]

	// Get the state from the ledger
	fireAlarmDataJson, err := stub.GetState(AlertorID)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + AlertorID + "\"}"
		return shim.Error(jsonResp)
	}

	if fireAlarmDataJson == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + AlertorID + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(fireAlarmDataJson)
}

type argStruct struct {
	Args []string `json:"Args"`
}

func asBytes(args []string) [][]byte {
	bytes := make([][]byte, len(args))
	for i, arg := range args {
		bytes[i] = []byte(arg)
	}
	return bytes
}

// invokeCC invokes another chaincode
// arg0: ID of chaincode to invoke
// arg1: Chaincode arguments in the form: {"Args": ["arg0", "arg1",...]}
func (t *SimpleChaincode) invokeCC(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting ID of chaincode to invoke and args")
	}

	ccID := args[0]
	invokeArgsJSON := args[1]

	argStruct := argStruct{}
	if err := json.Unmarshal([]byte(invokeArgsJSON), &argStruct); err != nil {
		return shim.Error(fmt.Sprintf("Invalid invoke args: %s", err))
	}

	if err := stub.PutState(stub.GetTxID()+"_invokedcc", []byte(ccID)); err != nil {
		return shim.Error(fmt.Sprintf("Error putting state: %s", err))
	}

	return stub.InvokeChaincode(ccID, asBytes(argStruct.Args), "")
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
