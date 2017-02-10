
package main

import (
	"fmt"
	"strings"
	"strconv"
	"encoding/json"
	"encoding/base64"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}


type Employee struct {
	name string 
	employeeId int
	project string
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	function, _ := stub.GetFunctionAndParameters()
    fmt.Println("Init is running " + function)
    return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
    function, args := stub.GetFunctionAndParameters()
    fmt.Println("Invoke is running " + function)
    
    if function =="initEmployee"{
		return t.initEmployee(stub,args)
	}
	return shim.Error("Received unknown function invocation")
}

func (t *SimpleChaincode) initEmployee(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	// ==== Input sanitation ====
	fmt.Println("- start initEmployee")
	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return shim.Error("3rd argument must be a non-empty string")
	}
	
	employeeName := args[0]
	project := strings.ToLower(args[2])
	employeeId, err := strconv.Atoi(args[1])
	employeeIdAsString := args[1]
	if err != nil {
		return shim.Error("2rd argument must be a numeric string")
	}
	
	employeeAsBytes, err := stub.GetState(employeeIdAsString)
	if err != nil {
		return shim.Error("Failed to get employee: " + err.Error())
	} else if employeeAsBytes != nil {
		fmt.Println("This employee already exists: " + employeeIdAsString)
		return shim.Error("This marble already exists: " + employeeIdAsString)
	}
	
	employee:= &Employee{employeeName,employeeId,project}
	
	employeeJSONasBytes, err := json.Marshal(employee)
	
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(employeeIdAsString, employeeJSONasBytes)
	
	if err != nil {
		return shim.Error(err.Error())
	}
	
	fmt.Println("- end initEmployee")
	return shim.Success(nil)
}

func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface) pb.Response {
    function, args := stub.GetFunctionAndParameters()
    fmt.Println("Query is running " + function)
    
    if function =="getEmployee"{
		return t.getEmployee(stub,args)
	}
	
	return shim.Error("Received unknown function query")
}

func (t *SimpleChaincode) getEmployee(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	if len(args) !=1{
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	employeeId := args[0]
	employee, err := stub.GetState(employeeId)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed retrieving employee [%s]: [%s]", employeeId, err))
	}
	
	return shim.Success([]byte(base64.StdEncoding.EncodeToString(employee)))
}


func main() {
    	err := shim.Start(new(SimpleChaincode))
    if err != nil {
        fmt.Printf("Error starting Simple chaincode: %s", err)
    }
}
