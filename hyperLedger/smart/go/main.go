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

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"log"
)

// cryptoChaincode is allows the following transactions
//
//	"put", "key", val - returns "OK" on success
//	"get", "key" - returns val stored previously

type Sensors struct {
	Type  string
	Value string
}

type DeviceChaincode struct {
	Id            string
	Input1        string
	Input2        string
	Output        string
	Alarm_battery string
	Alarm_power   string
	Sensor_error  string
	Temperatura   string
	Umidade       string
	LastUpdated   string
}

const (
	// AESKeyLength is the default AES key length
	AESKeyLength = 32
	// NonceSize is the default NonceSize
	NonceSize = 24
)

// Init implements chaincode's Init interface
func (t *DeviceChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke implements chaincode's Invoke interface
func (t *DeviceChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Invoke method gets called")
	function, args := stub.GetFunctionAndParameters()
	if function == "invoke" {
		// Make payment of X units from A to B
		return t.invoke(stub, args)
	} else if function == "put" {
		// update an entity from its state
		return t.put(stub, args)
	} else if function == "delete" {
		// Deletes an entity from its state
		return t.delete(stub, args)
	} else if function == "query" {
		// the old "Query" is now implemtned in invoke
		return t.query(stub, args)
	} else if function == "queryAll" {
		// the old "queryAll" is now implemtned in invoke
		return t.queryAll(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"invoke\" \"delete\" \"query\"")
}

// Transaction makes payment of X units from A to B
// Transaction makes payment of X units from A to B
func (t *DeviceChaincode) invoke(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	// Get the args from the transaction proposal
	fmt.Println("Init method gets called")
	if len(args) != 11 {
		return shim.Error("Incorrect arguments. Expecting a key and a value")
	}

	// Set up any variables or assets here by calling stub.PutState()
	//  Aval, err = strconv.Atoi(args[1])
	_id := args[0]

	key, err := stub.CreateCompositeKey("DeviceChaincode", []string{_id})

	isActive, err := stub.GetState(key)

	if err != nil {
		return shim.Error(err.Error())
	}

	println(args[0])
	println("Key: " + key)
	println(fmt.Sprintf("Key (hex): %x", key))
	println("State: " + string(isActive))

	if isActive == nil && err == nil {

		c := DeviceChaincode{
			Id:            args[0],
			Input1:        args[1],
			Input2:        args[2],
			Output:        args[3],
			Alarm_battery: args[4],
			Alarm_power:   args[5],
			Sensor_error:  args[6],
			Temperatura:   args[7],
			Umidade:       args[8],
			LastUpdated:   args[9],
		}

		//the struct will be saved as a Marshalled json
		objJson, _ := json.Marshal(c)
		err = stub.PutState(key, objJson)
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success(nil)
	}

	return shim.Error("SENSOR IS DISABLED or ALREADY ACTIVE!")
}

func (t *DeviceChaincode) put(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	// Get the args from the transaction proposal
	fmt.Println("Init method gets called")
	if len(args) != 11 {
		return shim.Error("Incorrect arguments. Expecting a key and a value")
	}

	// Set up any variables or assets here by calling stub.PutState()
	//  Aval, err = strconv.Atoi(args[1])
	_id := args[0]

	key, err := stub.CreateCompositeKey("DeviceChaincode", []string{_id})

	isActive, err := stub.GetState(key)

	if err != nil {
		return shim.Error(err.Error())
	}

	println("Key: " + key)
	println(fmt.Sprintf("Key (hex): %x", key))
	println("State: " + string(isActive))

	c := DeviceChaincode{
		Id:            args[0],
		Input1:        args[1],
		Input2:        args[2],
		Output:        args[3],
		Alarm_battery: args[4],
		Alarm_power:   args[5],
		Sensor_error:  args[6],
		Temperatura:   args[7],
		Umidade:       args[8],
		LastUpdated:   args[9],
	}

	//the struct will be saved as a Marshalled json
	objJson, _ := json.Marshal(c)
	err = stub.PutState(key, objJson)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)

}

// Deletes an entity from state
func (t *DeviceChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("delete method gets called")
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	A := args[0]

	key, err := stub.CreateCompositeKey("DeviceChaincode", []string{A})

	println("Key: " + key)

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	// Delete the key from the state in ledger
	errDel := stub.DelState(key)
	if errDel != nil {
		return shim.Error("Failed to delete state")
	}

	return shim.Success(nil)
}

// query callback representing the query of a chaincode
func (t *DeviceChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("query method gets called")
	//var A string // Entities
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	A := args[0]

	key, err := stub.CreateCompositeKey("DeviceChaincode", []string{A})

	println("Key: " + key)

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success(Avalbytes)

}

// query callback representing the query of a chaincode
func (t *DeviceChaincode) queryAll(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	log.Printf("GetAssetHistory: ID %v", args[0])

	//var A string // Entities
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	A := args[0]

	key, err := stub.CreateCompositeKey("DeviceChaincode", []string{A})

	println("Key: " + key)

	resultsIterator, err := stub.GetHistoryForKey(key)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return shim.Error(jsonResp)
	}
	defer resultsIterator.Close()

	var records []DeviceChaincode
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
			return shim.Error(jsonResp)
		}

		var device DeviceChaincode
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &device)
			if err != nil {
				jsonResp := "{\"Error\":\"json.Unmarshal " + A + "\"}"
				return shim.Error(jsonResp)
			}

			records = append(records, device)
		}
	}
	fmt.Printf("Query Response:%s\n", records)
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(records)

	//reqBodyBytes.Bytes() // this is the []byte

	return shim.Success(reqBodyBytes.Bytes())
}

func main() {
	err := shim.Start(new(DeviceChaincode))
	if err != nil {
		fmt.Printf("Error starting New key per invoke: %s", err)
	}
}
