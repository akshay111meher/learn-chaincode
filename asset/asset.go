
package main

import (
	"fmt"
	// "strings"
	// "strconv"
	"errors"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Square struct {
	Id string
  Side string
  Owner string
}

type Circle struct{
  Id string
  Radius string
  Owner string
}

type Error struct{
	Err string
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte,error) {
    fmt.Println("Init is running " + function)
    return nil,nil
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte,error) {

    fmt.Println("Invoke is running " + function)

    if function =="createSquare"{
		    return t.createSquare(stub,args)
	  }else if function =="createCircle"{
      return t.createCircle(stub,args)
    }else if function=="transferCircle"{
      return t.transferCircle(stub,args)
    }else if function=="transferSquare"{
      return t.transferSquare(stub,args)
    }
	error := Error{"function invocation error, "+function+" doesnt exist"}
	errorMarshal,_:= json.Marshal(error)
	stub.SetEvent("functionError",errorMarshal)
	return nil,errors.New("Received unknown function invocation")
}
func (t *SimpleChaincode) transferCircle(stub shim.ChaincodeStubInterface, args []string) ([]byte,error){
  if len(args) != 3 {
		error := Error{"Incorrect number of arguments. Expecting 3"}
		errorMarshal,_:= json.Marshal(error)
		stub.SetEvent("transferCircleError",errorMarshal)
		return nil,errors.New("Incorrect number of arguments. Expecting 3")
	}
  id:= args[0]
  newOwner:= args[1]
  preOwner:= args[2]

  circleAsBytes, err := stub.GetState(id)
	if err != nil {
		return nil,err
	} else if circleAsBytes != nil {
    var c Circle
    json.Unmarshal(circleAsBytes,&c)
    if c.Owner==preOwner{
      c.Owner= newOwner
      circleAsBytes,_ = json.Marshal(c)
      err = stub.PutState(id, circleAsBytes)

      if err != nil {
        return nil,err
      }
      stub.SetEvent("notifyTransferCircle",circleAsBytes)
      fmt.Println("- end transferCircle")
      return nil,nil
    }else{
      error := Error{"This circle doesnot belong to : "+preOwner}
      errorMarshal,_:= json.Marshal(error)
      stub.SetEvent("transferCircleError",errorMarshal)
  		fmt.Println("This circle doesnot belong to : "+preOwner)
  		return []byte("duplicate"),errors.New("This circle doesnot belong to : "+preOwner)
    }
	}else{
    error := Error{"This circle does not exists: " + id}
		errorMarshal,_:= json.Marshal(error)
		stub.SetEvent("transferCircleError",errorMarshal)
		fmt.Println("This circle does not exists: " + id)
		return []byte("duplicate"),errors.New("This circle doesnot exists: "+id)
  }
}
func (t *SimpleChaincode) transferSquare(stub shim.ChaincodeStubInterface, args []string) ([]byte,error){
  if len(args) != 3 {
    error := Error{"Incorrect number of arguments. Expecting 3"}
    errorMarshal,_:= json.Marshal(error)
    stub.SetEvent("transferCircleError",errorMarshal)
    return nil,errors.New("Incorrect number of arguments. Expecting 3")
  }
  id:= args[0]
  newOwner:= args[1]
  preOwner:= args[2]

  squareAsBytes, err := stub.GetState(id)
  if err != nil {
    return nil,err
  } else if squareAsBytes != nil {
    var s Square
    json.Unmarshal(squareAsBytes,&s)
    if s.Owner==preOwner{
      s.Owner= newOwner
      squareAsBytes,_ = json.Marshal(s)
      err = stub.PutState(id, squareAsBytes)

      if err != nil {
        return nil,err
      }
      stub.SetEvent("notifyTransferSquare",squareAsBytes)
      fmt.Println("- end transferSquare")
      return nil,nil
    }else{
      error := Error{"This square doesnot belong to : "+preOwner}
      errorMarshal,_:= json.Marshal(error)
      stub.SetEvent("transferSquareError",errorMarshal)
      fmt.Println("This square doesnot belong to : "+preOwner)
      return []byte("duplicate"),errors.New("This square doesnot belong to : "+preOwner)
    }
  }else{
    error := Error{"This square does not exists: " + id}
    errorMarshal,_:= json.Marshal(error)
    stub.SetEvent("transferSquareError",errorMarshal)
    fmt.Println("This square does not exists: " + id)
    return []byte("duplicate"),errors.New("This square doesnot exists: "+id)
  }
}
func (t *SimpleChaincode) createSquare(stub shim.ChaincodeStubInterface, args []string) ([]byte,error){
  if len(args) != 3 {
		error := Error{"Incorrect number of arguments. Expecting 3"}
		errorMarshal,_:= json.Marshal(error)
		stub.SetEvent("createSquareError",errorMarshal)
		return nil,errors.New("Incorrect number of arguments. Expecting 3")
	}
	// ==== Input sanitation ====
	fmt.Println("- start createSquare")
	if len(args[0]) <= 0 {
		error := Error{"1st argument must be a non-empty string"}
		errorMarshal,_:= json.Marshal(error)
		stub.SetEvent("createSquareError",errorMarshal)
		return nil,errors.New("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		error := Error{"2nd argument must be a non-empty string"}
		errorMarshal,_:= json.Marshal(error)
		stub.SetEvent("createSquareError",errorMarshal)
		return nil,errors.New("2nd argument must be a non-empty string")
	}
  if len(args[2]) <= 0 {
    error := Error{"3rd argument must be a non-empty string"}
    errorMarshal,_:= json.Marshal(error)
    stub.SetEvent("createSquareError",errorMarshal)
    return nil,errors.New("3rd argument must be a non-empty string")
  }

  id:= args[0]
  side:=args[1]
  owner:=args[2]

	squareAsBytes, err := stub.GetState(id)
	if err != nil {
		return nil,err
	} else if squareAsBytes != nil {
		error := Error{"This square already exists: " + id}
		errorMarshal,_:= json.Marshal(error)
		stub.SetEvent("duplicateSquare",errorMarshal)
		fmt.Println("This square already exists: " + id)
		return []byte("duplicate"),errors.New("This square already exists: "+id)
	}

	square:= Square{id,side,owner}
	fmt.Println(square)
	squareAsBytes, err = json.Marshal(square)
	fmt.Println(squareAsBytes)
	if err != nil {
		return nil,err
	}
	err = stub.PutState(id, squareAsBytes)

	if err != nil {
		return nil,err
	}

	stub.SetEvent("notifyCreateSquare",squareAsBytes)
	fmt.Println("- end createSquare")
	return nil,nil
}
func (t *SimpleChaincode) createCircle(stub shim.ChaincodeStubInterface, args []string) ([]byte,error){
  if len(args) != 3 {
		error := Error{"Incorrect number of arguments. Expecting 3"}
		errorMarshal,_:= json.Marshal(error)
		stub.SetEvent("createCircleError",errorMarshal)
		return nil,errors.New("Incorrect number of arguments. Expecting 3")
	}
	// ==== Input sanitation ====
	fmt.Println("- start createCircle")
	if len(args[0]) <= 0 {
		error := Error{"1st argument must be a non-empty string"}
		errorMarshal,_:= json.Marshal(error)
		stub.SetEvent("createCircleError",errorMarshal)
		return nil,errors.New("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		error := Error{"2nd argument must be a non-empty string"}
		errorMarshal,_:= json.Marshal(error)
		stub.SetEvent("createCircleError",errorMarshal)
		return nil,errors.New("2nd argument must be a non-empty string")
	}
  if len(args[2]) <= 0 {
    error := Error{"3rd argument must be a non-empty string"}
    errorMarshal,_:= json.Marshal(error)
    stub.SetEvent("createCircleError",errorMarshal)
    return nil,errors.New("3rd argument must be a non-empty string")
  }

  id:= args[0]
  radius:=args[1]
  owner:=args[2]

	circleAsBytes, err := stub.GetState(id)
	if err != nil {
		return nil,err
	} else if circleAsBytes != nil {
		error := Error{"This circle already exists: " + id}
		errorMarshal,_:= json.Marshal(error)
		stub.SetEvent("duplicateCircle",errorMarshal)
		fmt.Println("This circle already exists: " + id)
		return []byte("duplicate"),errors.New("This circle already exists: "+id)
	}

	circle:= Circle{id,radius,owner}
	fmt.Println(circle)
	circleAsBytes, err = json.Marshal(circle)
	fmt.Println(circleAsBytes)
	if err != nil {
		return nil,err
	}
	err = stub.PutState(id, circleAsBytes)

	if err != nil {
		return nil,err
	}

	stub.SetEvent("notifyCreateCircle",circleAsBytes)
	fmt.Println("- end createCircle")
	return nil,nil
}
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte,error) {
    fmt.Println("Query is running " + function)

    return nil,nil
}
func main() {
    err := shim.Start(new(SimpleChaincode))
  if err != nil {
      fmt.Printf("Error starting Simple chaincode: %s", err)
  }
}
