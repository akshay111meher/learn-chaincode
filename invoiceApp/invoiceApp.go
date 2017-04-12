
package main

import (
	"fmt"
	"strings"
	"strconv"
	"errors"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Error struct{
	Err string
}
type EmployeeWorkDetails interface{
	HoursPerWeek()
	HoursPerDay()
}
type Employee struct{
	Id string
	Name string
}

type Project struct{
	Id string
  Name string
	TxList []Transaction
}



func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte,error) {
    fmt.Println("Init is running " + function)
    return nil,nil
}
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte,error) {
    fmt.Println("Invoke is running " + function)

		if function =="initEmployee"{
			return t.initEmployee(stub,args)
		}else if function == "submitEfforts" {
			return t.submitEfforts(stub.args)
		}else if function =="createEmployee" {
			return t.createEmployee(stub.args)
		}else if function == "createProject"{
			return t.createProject(stub,args)
		}
	error := Error{"function invocation error, "+function+" doesnt exist"}
	errorMarshal,_:= json.Marshal(error)
	stub.SetEvent("functionError",errorMarshal)
	return nil,errors.New("Received unknown function invocation")
}
