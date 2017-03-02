
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
type EmployeeAndProject struct{
	Emp EmployeeEfforts
	Prj ProjectEfforts
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

type EmployeeEfforts struct{
	Efforts string
	Project string
}

type ProjectEfforts struct{
	Efforts string
	Employee string
}

type EmployeeInvoice struct{
	List []EmployeeEfforts
}
type ProjectInvoice struct{
	List []ProjectEfforts
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
	}else if function =="changeProject"{
		return t.changeProject(stub,args)
	}else if function =="submitEfforts"{
		return t.submitEfforts(stub,args)
	}
	stub.SetEvent("functionError",[]byte("function invocation error, "+function+" doesnt exist"))
	return nil,errors.New("Received unknown function invocation")
}

func (t *SimpleChaincode) submitEfforts(stub shim.ChaincodeStubInterface, args []string)([]byte,error){
	if len(args)!=4{
		stub.SetEvent("submitEffortsError",[]byte("Incorrect Number of arguments. Expecting 4"))
		return nil,errors.New("Incorrect Number of arguments. Expecting 4")
	}
	fmt.Println("- start submitEfforts -")
	employeeId := args[0]
	timestamp := args[1]
	employeeIdAndTimestamp := args[2]
	efforts := args[3]
	employeeJSONasBytes,err := stub.GetState(employeeId)

	if err != nil {
		stub.SetEvent("submitEffortsError",[]byte("unable to find employee or fetching employee error"))
		return nil,err
	}
	if len(employeeJSONasBytes)==0{
		stub.SetEvent("submitEffortsError",[]byte("employee with ID "+employeeId+" does not exists"))
		return nil,errors.New("employee with ID "+employeeId+" does not exists")
	}
	var e Employee
	json.Unmarshal(employeeJSONasBytes,&e)

	var ei EmployeeInvoice
	var pi ProjectInvoice

	employeeInvoiceAsBytes,err := stub.GetState(employeeIdAndTimestamp)

	if err != nil{
		stub.SetEvent("submitEffortsError",[]byte("unable to find employee with timestamp or fetching employee and timestamp error"))
		return nil,err
	}
	json.Unmarshal(employeeInvoiceAsBytes,&ei)

  var empAndPrj EmployeeAndProject
	var ef EmployeeEfforts
	ef = EmployeeEfforts{efforts,e.Project}
	// toSend1,err := json.Marshal(ef)
	// stub.SetEvent("notifySubmitEfforts",toSend1)
	ei.List = append(ei.List,ef)

	employeeInvoiceAsBytes,err = json.Marshal(ei)
	if err!= nil{
		stub.SetEvent("submitEffortsError",[]byte("employee invoice marshalling error"))
		return nil,err
	}
	err = stub.PutState(employeeIdAndTimestamp,employeeInvoiceAsBytes)
	if err !=nil{
		return nil,err
	}

	var projectIdAndTimestamp string
	projectIdAndTimestamp = e.Project+""+timestamp

	projectInvoiceAsBytes,err := stub.GetState(projectIdAndTimestamp)

	json.Unmarshal(projectInvoiceAsBytes,&pi)

	var pf ProjectEfforts
	pf = ProjectEfforts{efforts,employeeId}
	pi.List = append(pi.List,pf)
	projectInvoiceAsBytes,err = json.Marshal(pi)

	if err!= nil{
		stub.SetEvent("submitEffortsError",[]byte("unable to find project or fetching project error"))
		return nil,err
	}

	err = stub.PutState(projectIdAndTimestamp,projectInvoiceAsBytes)

	if err!= nil{
		stub.SetEvent("submitEffortsError",[]byte("unable to find project or fetching project time error"))
		return nil,err
	}
	empAndPrj = EmployeeAndProject{Emp:ef,Prj:pf}
	toSend1,err := json.Marshal(empAndPrj)
	stub.SetEvent("notifySubmitEfforts",toSend1)
	return nil,nil
}
func (t *SimpleChaincode) changeProject(stub shim.ChaincodeStubInterface, args []string)([]byte,error){
	if len(args)!=2 {
		stub.SetEvent("changeProjectError",[]byte("Incorrect Number of arguments. Expecting 2"))
		return nil,errors.New("Incorrect Number of arguments. Expecting 2")
	}
	// ==== Input sanitation ====
	fmt.Println("- start changeProject")
	if len(args[0]) <= 0 {
		stub.SetEvent("changeProjectError",[]byte("1st argument must be a non-empty string"))
		return nil,errors.New("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		stub.SetEvent("changeProjectError",[]byte("2nd argument must be a non-empty string"))
		return nil,errors.New("2nd argument must be a non-empty string")
	}
	employeeId := args[0]
	project := strings.ToLower(args[1])
	_, err := strconv.Atoi(project)

	if err != nil {
		stub.SetEvent("changeProjectError",[]byte("2nd argument must be a numeric string"))
		return nil,errors.New("2nd argument must be a numeric string")
	}
	_, err = strconv.Atoi(employeeId)

	if err != nil {
		stub.SetEvent("changeProjectError",[]byte("1st argument must be a numeric string"))
		return nil,errors.New("1st argument must be a numeric string")
	}

	_, err = stub.GetState(project)
	if err != nil {
		return nil,err
	}

	employee, err := stub.GetState(employeeId)
	if err != nil {
		return nil,err
	}
	if len(employee)==0{
		stub.SetEvent("changeProjectError",[]byte("employee with ID "+employeeId+" doesnot exists"))
		return nil,errors.New("employee with ID "+employeeId+" doesnot exists")
	}

	var e Employee

	json.Unmarshal(employee,&e)
	e.Project = project
	fmt.Println(e)

	employeeJSONasBytes, err := json.Marshal(e)
	fmt.Println(employeeJSONasBytes)
	if err != nil {
		return nil,err
	}
	err = stub.PutState(employeeId, employeeJSONasBytes)

	if err != nil {
		return nil,err
	}
	stub.SetEvent("notifyChangeProject",employeeJSONasBytes)
	return nil,nil
}
func (t *SimpleChaincode) initProject(stub shim.ChaincodeStubInterface, args []string) ([]byte,error){

	if len(args) != 5 {
		stub.SetEvent("initProjectError",[]byte("Incorrect number of arguments. Expecting 5"))
		return nil,errors.New("Incorrect number of arguments. Expecting 5")
	}
	// ==== Input sanitation ====
	fmt.Println("- start initProject")
	if len(args[0]) <= 0 {
		stub.SetEvent("initProjectError",[]byte("1st argument must be a non-empty string"))
		return nil,errors.New("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		stub.SetEvent("initProjectError",[]byte("2nd argument must be a non-empty string"))
		return nil,errors.New("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		stub.SetEvent("initProjectError",[]byte("3th argument must be a non-empty string"))
		return nil,errors.New("3rd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		stub.SetEvent("initProjectError",[]byte("4th argument must be a non-empty string"))
		return nil,errors.New("4th argument must be a non-empty string")
	}
	if len(args[4]) <= 0 {
		stub.SetEvent("initProjectError",[]byte("5th argument must be a non-empty string"))
		return nil,errors.New("5th argument must be a non-empty string")
	}

	projectName := args[0]
	customerOf := strings.ToLower(args[2])
	projectId, err := strconv.Atoi(args[1])
	projectIdAsString := args[1]
	startDate:= args[3]
	endDate:= args[4]
	if err != nil {
		stub.SetEvent("initProjectError",[]byte("2nd argument must be a numeric string"))
		return nil,errors.New("2nd argument must be a numeric string")
	}
	customerJSONasBytes,err := stub.GetState(customerOf)
	if err!=nil {
		stub.SetEvent("notifyInitProject",[]byte("This customer doesnot exists: "+customerOf))
		return nil,errors.New("This customer doesnot exists: "+customerOf)
	}
	if len(customerJSONasBytes)==0 {
		stub.SetEvent("notifyInitProject",[]byte("This customer doesnot exists: "+customerOf))
		return nil,errors.New("This project doesnot exists: "+customerOf)
	}

	projectAsBytes, err := stub.GetState(projectIdAsString)
	if err != nil {
		return nil,err
	} else if projectAsBytes != nil {
		stub.SetEvent("duplicateProject",[]byte("This project already exists: " + projectIdAsString))
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
	stub.SetEvent("notifyInitProject",projectJSONasBytes)
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
			stub.SetEvent("initCustomerError",[]byte("Incorrect number of arguments. Expecting 2"))
			return nil,errors.New("Incorrect number of arguments. Expecting 2")
		}
			// ==== Input sanitation ====
		fmt.Println("- start initCustomer")
		if len(args[0]) <= 0 {
			stub.SetEvent("initCustomerError",[]byte("1st argument must be a non-empty string"))
			return nil,errors.New("1st argument must be a non-empty string")
		}
		if len(args[1]) <= 0 {
			stub.SetEvent("initCustomerError",[]byte("2nd argument must be a non-empty string"))
			return nil,errors.New("2nd argument must be a non-empty string")
		}

		customerName := args[0]
		customerId, err := strconv.Atoi(args[1])
		customerIdAsString := args[1]
		if err != nil {
			stub.SetEvent("initCustomerError",[]byte("2nd argument must be a numeric string"))
			return nil,errors.New("2nd argument must be a numeric string")
		}
		customerAsBytes, err := stub.GetState(customerIdAsString)
		if err != nil {
			 return nil,err
		} else if customerAsBytes != nil {
			fmt.Println("This customer already exists: " + customerIdAsString)
			stub.SetEvent("duplicateCustomer",[]byte("This customer already exists: " + customerIdAsString))
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
		stub.SetEvent("notifyInitCustomer",customerJSONasBytes)
		fmt.Println("- end initCustomer")
		return nil,nil
}

func (t *SimpleChaincode) initEmployee(stub shim.ChaincodeStubInterface, args []string) ([]byte,error){

	if len(args) != 3 {
		stub.SetEvent("initEmployeeError",[]byte("Incorrect number of arguments. Expecting 3"))
		return nil,errors.New("Incorrect number of arguments. Expecting 3")
	}
	// ==== Input sanitation ====
	fmt.Println("- start initEmployee")
	if len(args[0]) <= 0 {
		stub.SetEvent("initEmployeeError",[]byte("1st argument must be a non-empty string"))
		return nil,errors.New("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		stub.SetEvent("initEmployeeError",[]byte("2nd argument must be a non-empty string"))
		return nil,errors.New("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		stub.SetEvent("initEmployeeError",[]byte("3rd argument must be a non-empty string"))
		return nil,errors.New("3rd argument must be a non-empty string")
	}

	employeeName := args[0]
	project := strings.ToLower(args[2])
	employeeId, err := strconv.Atoi(args[1])
	employeeIdAsString := args[1]

	if err != nil {
		stub.SetEvent("initEmployeeError",[]byte("2nd argument must be a numeric string"))
		return nil,errors.New("2nd argument must be a numeric string")
	}
	projectJSONasBytes,err := stub.GetState(project)
	if err!=nil {
		stub.SetEvent("notifyInitEmployee",[]byte("This project doesnot exists: "+project))
		return nil,errors.New("This project doesnot exists: "+project)
	}
	if len(projectJSONasBytes)==0 {
		stub.SetEvent("notifyInitEmployee",[]byte("This project doesnot exists: "+project))
		return nil,errors.New("This project doesnot exists: "+project)
	}

	employeeAsBytes, err := stub.GetState(employeeIdAsString)
	if err != nil {
		return nil,err
	} else if employeeAsBytes != nil {
		stub.SetEvent("duplicateEmployee",[]byte("This employee already exists: " + employeeIdAsString))
		fmt.Println("This employee already exists: " + employeeIdAsString)
		return []byte("duplicate"),errors.New("This employee already exists: "+employeeIdAsString)
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
	stub.SetEvent("notifyInitEmployee",employeeJSONasBytes)
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
	}else if function == "getEfforts"{
		return t.getEfforts(stub,args)
	}

	return nil,errors.New("Received unknown function query")
}
func (t *SimpleChaincode) getEfforts(stub shim.ChaincodeStubInterface, args []string) ([]byte,error){
	fmt.Printf("getEfforts called")
	if len(args) !=1{
		return nil,errors.New("Incorrect number of arguments. Expecting 1")
	}
	idAndTimestamp := args[0]
	invoice, err := stub.GetState(idAndTimestamp)
	fmt.Println(invoice)
	if err != nil {
		return nil,err
	}

	return invoice,nil
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
