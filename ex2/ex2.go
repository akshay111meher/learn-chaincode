
package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	function, _ := stub.GetFunctionAndParameters()
    fmt.Println("Init is running " + function)
    return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
    function, _ := stub.GetFunctionAndParameters()
    fmt.Println("Invoke is running " + function)
	return shim.Success(nil)
}


func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface) pb.Response {
    function, _ := stub.GetFunctionAndParameters()
    fmt.Println("Query is running " + function)
	return shim.Success(nil)
}


func main() {
    	err := shim.Start(new(SimpleChaincode))
    if err != nil {
        fmt.Printf("Error starting Simple chaincode: %s", err)
    }
}
