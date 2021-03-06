package main

import (
	// "encoding/base64"
	"errors"
	"fmt"
	"strings"
	"strconv"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/crypto/primitives"
)

type AssetManagementChaincode struct {
}
type Error struct{
  Err string
}
type Circle struct{
  Id string
  Owner string
  Radius string
}
func (t *AssetManagementChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
  fmt.Println("Init Chaincode...")
	err := stub.CreateTable("AssetsOwnership", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "Asset", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "Owner", Type: shim.ColumnDefinition_BYTES, Key: false},
	})
	if err != nil {
		return nil, errors.New("Failed creating AssetsOnwership table.")
	}
  return nil,nil
}

func (t *AssetManagementChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error){
  fmt.Println("Invoke is invoking function "+function)

  if function == "createCircle"{
    return t.createCircle(stub,args)
  }
	// else if function == "createCircle"{
  //   return t.createCircle(stub,args)
  // }
  return nil,nil
}

func (t *AssetManagementChaincode) transferCircle(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
	if len(args) != 2 {
    error := Error{"Incorrect number of arguments. Expecting 2"}
    errorMarshal, _ := json.Marshal(error)
    stub.SetEvent("transferCircleError", errorMarshal)
    return nil, errors.New("Incorrect number of arguments. Expecting 2")
  }

	id:= args[0]
	assestAsJson,err := stub.GetState(id)

	if len(assestAsJson)==0{
		fmt.Println("Asset doesnt exists")
		error := Error{"asset doesnt exists"}
		errorMarshal, _ := json.Marshal(error)
		stub.SetEvent("transferCircleError", errorMarshal)
		return nil, errors.New("asset doesnt exists")
	}

	var cir Circle
	json.Unmarshal(assestAsJson,&cir)

	callerCert, err := stub.GetCallerCertificate()
	if err != nil {
		fmt.Println("Failed getting metadata")
		return nil, errors.New("Failed getting metadata.")
	}
	cc := convert(callerCert[:])

	if cc== cir.Owner{
		cir.Owner = args[1]
		assestAsJson,err = json.Marshal(cir)
		stub.PutState(id,assestAsJson)
	}else{
		fmt.Println("you are not the rightful owner of the asset")
		error := Error{"you are not the rightful owner of the asset"}
		errorMarshal, _ := json.Marshal(error)
		stub.SetEvent("transferCircleError", errorMarshal)
		return nil, errors.New("you are not the rightful owner of the asset")
	}
	return nil,nil
}

func (t *AssetManagementChaincode) createCircle(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
  if len(args) != 3 {
    error := Error{"Incorrect number of arguments. Expecting 3"}
    errorMarshal, _ := json.Marshal(error)
    stub.SetEvent("createCircleError", errorMarshal)
    return nil, errors.New("Incorrect number of arguments. Expecting 3")
  }
  // ==== Input sanitation ====
  fmt.Println("- start createCircle")
  if len(args[0]) <= 0 {
    error := Error{"1st argument must be a non-empty string"}
    errorMarshal, _ := json.Marshal(error)
    stub.SetEvent("createCircleError", errorMarshal)
    return nil, errors.New("1st argument must be a non-empty string")
  }
  if len(args[1]) <= 0 {
    error := Error{"2nd argument must be a non-empty string"}
    errorMarshal, _ := json.Marshal(error)
    stub.SetEvent("createCircleError", errorMarshal)
    return nil, errors.New("2nd argument must be a non-empty string")
  }
  if len(args[2]) <= 0 {
    error := Error{"3rd argument must be a non-empty string"}
    errorMarshal, _ := json.Marshal(error)
    stub.SetEvent("createCircleError", errorMarshal)
    return nil, errors.New("3rd argument must be a non-empty string")
  }
	callerCert, err := stub.GetCallerCertificate()
	if err != nil {
		fmt.Println("Failed getting metadata")
		return nil, errors.New("Failed getting metadata.")
	}
	payload, err := stub.GetPayload()
	if err != nil {
		return nil, errors.New("Failed getting payload")
	}
	binding, err := stub.GetBinding()
	if err != nil {
		return nil, errors.New("Failed getting binding")
	}
  cc := convert(callerCert[:])
	p := convert(payload[:])
	b := convert(binding[:])

	if len(callerCert) == 0 {
		fmt.Println("Invalid caller certificate. Empty.")
		fmt.Println(p)
		fmt.Println(b)
		return nil, errors.New(cc)
	}

	fmt.Println("The caller is [%x]", callerCert)

	id:= args[0]
	owner:= cc
	radius:= args[2]

	assestAsJson,err := stub.GetState(id)

	if len(assestAsJson)>0{
		fmt.Println("Asset already exists")
		return nil, errors.New("Cant create asset already exists")
	}

	var c Circle
	c = Circle{id,owner,radius}
	assestAsJson,err = json.Marshal(c)
	stub.PutState(id,assestAsJson)
	return nil,nil
}


func (t *AssetManagementChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte,error) {
    fmt.Println("Query is running " + function)

    if function =="getCircle"{
			return t.getCircle(stub,args)
		}else if function == "getCertificate" {
			return t.getCertificate(stub,args)
		}
	return nil,errors.New("Received unknown function query")
}

func (t *AssetManagementChaincode) getCircle(stub shim.ChaincodeStubInterface, args []string) ([]byte,error) {
	fmt.Println("getCircle called")
	if len(args) !=1{
		return nil,errors.New("Incorrect number of arguments. Expecting 1")
	}
	Id := args[0]
	circle, err := stub.GetState(Id)
	if err != nil {
		return nil,err
	}

	return circle,nil
}
func (t *AssetManagementChaincode) getCertificate(stub shim.ChaincodeStubInterface, args []string) ([]byte,error) {
	fmt.Println("getCertificate called")
	if len(args) !=1{
		return nil,errors.New("Incorrect number of arguments. Expecting 1")
	}
	callerCert,err := stub.GetCallerCertificate()
	if err != nil{
		return callerCert,nil
	}else{
		return nil,err
	}
}
func convert( b []byte ) string {
    s := make([]string,len(b))
    for i := range b {
        s[i] = strconv.Itoa(int(b[i]))
    }
    return strings.Join(s,",")
}
func main() {
	primitives.SetSecurityLevel("SHA3", 256)
	err := shim.Start(new(AssetManagementChaincode))
	if err != nil {
		fmt.Printf("Error starting AssetManagementChaincode: %s", err)
	}
}
