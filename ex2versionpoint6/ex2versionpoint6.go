
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


type Employee struct {
	Name string `json:"name"`
	EmployeeId int `json:"employeeId"`
	Project string `json:"project"`
}

type Customer struct {
	Name string `json:"name"`
	CustomerId int `json:"customerId"`
}
type Project struct {
	Name string `json:"name"`
	ProjectId int `json:"projectId"`
	CustomerOf string `json:"customerOf"`
	StartTime string `json:"startDate"`
	EndTime string `json:"endDate"`
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte,error) {
    fmt.Println("Init is running " + function)
    return nil,nil
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte,error) {
    
    fmt.Println("Invoke is running " + function)
    
    if function =="initEmployee"{
		return t.initEmployee(stub,args)
	}else if function =="getEmployee"{
		return t.getEmployee(stub,args)
	}else if function =="initCustomer"{
		return t.initCustomer(stub,args)
	}else if function =="getCustomer"{
		return t.getCustomer(stub,args)
	}else if function =="initProject"{
		return t.initProject(stub,args)
	}else if function =="getProject"{
		return t.getProject(stub,args)
	}
	return nil,errors.New("Received unknown function invocation")
}

func (t *SimpleChaincode) initProject(stub shim.ChaincodeStubInterface, args []string) ([]byte,error){
	
	if len(args) != 3 {
		return nil,errors.New("Incorrect number of arguments. Expecting 3")
	}
	// ==== Input sanitation ====
	fmt.Println("- start initProject")
	if len(args[0]) <= 0 {
		return nil,errors.New("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return nil,errors.New("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return nil,errors.New("3rd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return nil,errors.New("4th argument must be a non-empty string")
	}
	if len(args[4]) <= 0 {
		return nil,errors.New("5th argument must be a non-empty string")
	}
	
	projectName := args[0]
	customerOf := strings.ToLower(args[2])
	projectId, err := strconv.Atoi(args[1])
	projectIdAsString := args[1]
	startDate:= args[3]
	endDate:= args[4]
	if err != nil {
		return nil,errors.New("2nd argument must be a numeric string")
	}
	
	projectAsBytes, err := stub.GetState(projectIdAsString)
	if err != nil {
		return nil,err
	} else if projectAsBytes != nil {
		fmt.Println("This project already exists: " + projectIdAsString)
		return nil, errors.New("This project already exists "+projectIdAsString)
	}
	
	project:= Project{projectName,projectId,customerOf,startDate,endDate}
	
	projectJSONasBytes, err := json.Marshal(project)
	
	if err != nil {
		return nil,err
	}
	err = stub.PutState(projectIdAsString, projectJSONasBytes)
	
	if err != nil {
		return nil,err
	}
	
	// composite key creation for searching projects according to customer
	/*
	indexName := "customerOf"
	customerOfIndexKey, err := stub.CreateCompositeKey(indexName, []string{customerOf, projectIdAsString})
	if err != nil {
		return nil,err
	}

	value := []byte{0x00}
	stub.PutState(customerOfIndexKey, value)
    */
	fmt.Println("- end initProject")
	return nil,nil
}

func (t *SimpleChaincode) getProject(stub shim.ChaincodeStubInterface, args []string) ([]byte,error){
	if len(args) !=1{
		return nil,errors.New("Incorrect number of arguments. Expecting 1")
	}
	projectId := args[0]
	project, err := stub.GetState(projectId)
	if err != nil {
		return nil,err
	}
	
	return project,nil
}


func (t *SimpleChaincode) getCustomer(stub shim.ChaincodeStubInterface, args []string) ([]byte,error){
	if len(args) !=1{
		return nil,errors.New("Incorrect number of arguments. Expecting 1")
	}
	customerId := args[0]
	customer, err := stub.GetState(customerId)
	if err != nil {
		return nil,err
	}
	
	return customer,nil
}

func (t *SimpleChaincode) initCustomer(stub shim.ChaincodeStubInterface,args []string) ([]byte,error) {
		if len(args) != 2{
			return nil,errors.New("Incorrect number of arguments. Expecting 2")
		}	
			// ==== Input sanitation ====
		fmt.Println("- start initCustomer")
		if len(args[0]) <= 0 {
			return nil,errors.New("1st argument must be a non-empty string")
		}
		if len(args[1]) <= 0 {
			return nil,errors.New("2nd argument must be a non-empty string")
		}
		
		customerName := args[0]
		customerId, err := strconv.Atoi(args[1])
		customerIdAsString := args[1]
		if err != nil {
			return nil,errors.New("2nd argument must be a numeric string")
		}
		customerAsBytes, err := stub.GetState(customerIdAsString)
		if err != nil {
			 return nil,err
		} else if customerAsBytes != nil {
			fmt.Println("This customer already exists: " + customerIdAsString)
			return nil,errors.New("This customer already exists: "+customerIdAsString)
		}

		customer:= Customer{customerName,customerId}
		
		customerJSONasBytes, err := json.Marshal(customer)
		
		if err != nil {
			return nil,err
		}
		err = stub.PutState(customerIdAsString, customerJSONasBytes)
		
		if err != nil {
			return nil,err
		}
		
		fmt.Println("- end initCustomer")
		return nil,nil
}

func (t *SimpleChaincode) initEmployee(stub shim.ChaincodeStubInterface, args []string) ([]byte,error){
	
	if len(args) != 3 {
		return nil,errors.New("Incorrect number of arguments. Expecting 3")
	}
	// ==== Input sanitation ====
	fmt.Println("- start initEmployee")
	if len(args[0]) <= 0 {
		return nil,errors.New("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return nil,errors.New("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return nil,errors.New("3rd argument must be a non-empty string")
	}
	
	employeeName := args[0]
	project := strings.ToLower(args[2])
	employeeId, err := strconv.Atoi(args[1])
	employeeIdAsString := args[1]
	if err != nil {
		return nil,errors.New("2nd argument must be a numeric string")
	}
	
	employeeAsBytes, err := stub.GetState(employeeIdAsString)
	if err != nil {
		return nil,err
	} else if employeeAsBytes != nil {
		fmt.Println("This employee already exists: " + employeeIdAsString)
		return byte[]("duplicate"),errors.New("This employee already exists: "+employeeIdAsString)
	}
	
	employee:= Employee{employeeName,employeeId,project}
	fmt.Println(employee)
	employeeJSONasBytes, err := json.Marshal(employee)
	fmt.Println(employeeJSONasBytes)
	if err != nil {
		return nil,err
	}
	err = stub.PutState(employeeIdAsString, employeeJSONasBytes)
	
	if err != nil {
		return nil,err
	}
	employee_temp,_ := stub.GetState(employeeIdAsString)
	fmt.Println(employee_temp)
	// composite key to get employees by project
	
	/* 
	indexName := "project"
	projectIndexKey, err := stub.CreateCompositeKey(indexName, []string{project, employeeIdAsString})
	if err != nil {
		return nil,err
	}

	value := []byte{0x00}
	stub.PutState(projectIndexKey, value) 
	*/
	fmt.Println("- end initEmployee")
	return nil,nil
}

func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte,error) {
    fmt.Println("Query is running " + function)
    
    if function =="getEmployee"{
		return t.getEmployee(stub,args)
	}else if function =="getCustomer"{
		return t.getCustomer(stub,args)
	}else if function =="getProject"{
		return t.getProject(stub,args)
	}
	
	return nil,errors.New("Received unknown function query")
}

func (t *SimpleChaincode) getEmployee(stub shim.ChaincodeStubInterface, args []string) ([]byte,error){
	fmt.Printf("getEmployee called")
	if len(args) !=1{
		return nil,errors.New("Incorrect number of arguments. Expecting 1")
	}
	employeeId := args[0]
	employee, err := stub.GetState(employeeId)
	fmt.Println(employee)
	if err != nil {
		return nil,err
	}
	
	return employee,nil
}


func main() {
    	err := shim.Start(new(SimpleChaincode))
    if err != nil {
        fmt.Printf("Error starting Simple chaincode: %s", err)
    }
}
