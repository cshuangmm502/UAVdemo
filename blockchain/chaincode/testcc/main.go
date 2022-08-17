package testcc

import (
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
)

type InformationCC struct {
}


func main() {
	// Create a new Smart Contract
	err := shim.Start(new(InformationCC))
	if err != nil {
		fmt.Printf("Error starting Air chaincode: %s", err)
	}
}

// 实现 Init 方法, 实例化账本时使用。
func (s *InformationCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// 实现 Invoke 方法
func (s *InformationCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	// 获取函数名称、参数
	fn, args := stub.GetFunctionAndParameters()

	//调用对应函数
	if fn == "addInformation" {
		return s.addInformation(stub, args)

	} else if fn == "queryInformation" {
		return s.queryInformation(stub, args)

	}

	return shim.Error("Invalid Smart Contract function name.")
}
