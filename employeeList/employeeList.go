package main

import (
	"fmt"
	"errors"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type EmployeeList struct {
	NumberList []string
}
type CustomerList struct {
	NumberList []string
}
type ProjectList struct {
	NumberList []string
}
type SimpleChaincode struct {
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte,error) {
    fmt.Println("Init is running " + function)
    var array []string
	list1 := EmployeeList{NumberList:array}
	list2 := CustomerList{NumberList:array}
	list3 := ProjectList{NumberList:array}
	
	employeeListMarshal, err := json.Marshal(list1)	
	if err != nil {
		return nil,err
	}
	err = stub.PutState("WiproEmployees", employeeListMarshal)
	if err != nil {
		return nil,err
	}
	
	customerListMarshal, err := json.Marshal(list2)	
	if err != nil {
		return nil,err
	}
	err = stub.PutState("WiproCustomers", customerListMarshal)
	if err != nil {
		return nil,err
	}
	
	projectListMarshal, err := json.Marshal(list3)	
	if err != nil {
		return nil,err
	}
	err = stub.PutState("WiproProjects", projectListMarshal)
	if err != nil {
		return nil,err
	}
    return nil,nil
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface , function string, args []string) ([]byte,error){
	fmt.Println("Invoke is running "+function)
	if function =="initEmployee"{
		return t.initEmployee(stub,args)
	}else if function == "initCustomer"{
		return t.initCustomer(stub,args)
	}else if function == "initProject"{
		return t.initProject(stub,args)
	}
	
	return nil,errors.New("Received unknown function invoke")
}

func (t *SimpleChaincode) initEmployee(stub shim.ChaincodeStubInterface, args []string) ([]byte,error){
	if len(args)!=1{
		return nil,errors.New("Incorrect number of arguments. Expecting 1")
	}
	employeeListMarshal, err := stub.GetState("WiproEmployees")
	if err != nil {
		return nil,err
	}
	var e EmployeeList
	json.Unmarshal(employeeListMarshal,&e)
	e.NumberList = append(e.NumberList,args[0])
	
	employeeListMarshal, err = json.Marshal(e)
	if err!=nil{
		return nil,err
	}
	
	err = stub.PutState("WiproEmployees", employeeListMarshal)
	if err!=nil{
		return nil,err
	}
	return nil,nil
}

func (t *SimpleChaincode) initCustomer(stub shim.ChaincodeStubInterface, args []string) ([]byte,error){
	if len(args)!=1{
		return nil,errors.New("Incorrect number of arguments. Expecting 1")
	}
	customerListMarshal, err := stub.GetState("WiproCustomers")
	if err != nil {
		return nil,err
	}
	var c CustomerList
	json.Unmarshal(customerListMarshal,&c)
	c.NumberList = append(c.NumberList,args[0])
	
	customerListMarshal, err = json.Marshal(c)
	if err!=nil{
		return nil,err
	}
	
	err = stub.PutState("WiproCustomers", customerListMarshal)
	if err!=nil{
		return nil,err
	}
	return nil,nil
}
func (t *SimpleChaincode) initProject(stub shim.ChaincodeStubInterface, args []string) ([]byte,error){
	if len(args)!=1{
		return nil,errors.New("Incorrect number of arguments. Expecting 1")
	}
	projectListMarshal, err := stub.GetState("WiproProjects")
	if err != nil {
		return nil,err
	}
	var p ProjectList
	json.Unmarshal(projectListMarshal,&p)
	p.NumberList = append(p.NumberList,args[0])
	
	projectListMarshal, err = json.Marshal(p)
	if err!=nil{
		return nil,err
	}
	
	err = stub.PutState("WiproProjects", projectListMarshal)
	if err!=nil{
		return nil,err
	}
	return nil,nil
}

func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface , function string, args []string) ([]byte,error){
	fmt.Println("Query is running "+function)
	if function=="allEmployees"{
		return t.allEmployees(stub,args)
	}else if function=="allCustomers"{
		return t.allCustomers(stub,args)
	}else if function=="allProjects"{
		return t.allProjects(stub,args)
	}
	return nil,errors.New("Received unknown function query "+function)
}
func (t *SimpleChaincode) allCustomers(stub shim.ChaincodeStubInterface, args []string)([]byte, error){
	if len(args)!=0{
		return nil,errors.New("Received incorrect number of arguments. Expected 0")
	}
	customerListMarshal, err := stub.GetState("WiproCustomers")
	if err != nil {
		return nil,err
	}
	return customerListMarshal,nil
}
func (t *SimpleChaincode) allProjects(stub shim.ChaincodeStubInterface, args []string)([]byte, error){
	if len(args)!=0{
		return nil,errors.New("Received incorrect number of arguments. Expected 0")
	}
	projectsListMarshal, err := stub.GetState("WiproProjects")
	if err != nil {
		return nil,err
	}
	return projectsListMarshal,nil
}
func (t *SimpleChaincode) allEmployees(stub shim.ChaincodeStubInterface, args []string)([]byte, error){
	if len(args)!=0{
		return nil,errors.New("Received incorrect number of arguments. Expected 0")
	}
	employeeListMarshal, err := stub.GetState("WiproEmployees")
	if err != nil {
		return nil,err
	}
	return employeeListMarshal,nil
}
func main() {
    	err := shim.Start(new(SimpleChaincode))
    if err != nil {
        fmt.Printf("Error starting Simple chaincode: %s", err)
    }
}
