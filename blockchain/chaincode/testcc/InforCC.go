package testcc
import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
)

func (s *InformationCC) init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success([]byte("init InforCC succeed"))
}

// 将对象序列化后保存至账本中
func PutInformation(stub shim.ChaincodeStubInterface, information Information) bool{
	b, err := json.Marshal(information)
	if err != nil {
		fmt.Println(err)
		return false
	}
	err = stub.PutState(information.TraceNo, b)
	if err != nil {
		return false
	}
	return true
}

// 根据指定的Id 查询对应的状态，反序列化后并返回对象
func GetInfoByTraceId(stub shim.ChaincodeStubInterface, traceId string) (Information, bool) {

	var information Information
	b, err := stub.GetState(traceId)
	if err != nil || b == nil {
		return information, false
	} // 有错误 或者 Id不存在[id不存在GetState()返回 nil, nil]

	err = json.Unmarshal(b, &information)
	if err != nil {
		return information, false
	}

	return information, true
}

func (s *InformationCC) addInformation(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	information := Information{}

	err := json.Unmarshal([]byte(args[0]), &information)

	if err != nil {
		return shim.Error("Unmarshal information failed")
	}

	_, exist := GetInfoByTraceId(stub, information.TraceNo)
	if exist {
		return shim.Error("Id specified already exists")
	}

	flag := PutInformation(stub, information)
	if !flag {
		return shim.Error("Add record failed")
	}

	return shim.Success([]byte("Add record succeed"))

}

func (s *InformationCC) queryInformation(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("incorrect nums of  args, expecting 1")
	}

	result, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("query failed according to id")
	}
	if result == nil {
		return shim.Error("get nothing according to id")
	}
	return shim.Success(result)
}