package main

import (
	"UAVdemo/blockchain/sdkInit"
	"fmt"
	"os"
)

const (
	configFile  = "./sdkconf.yaml"
	initialized = false
	testcc = "testcc"
)

func main() {
	sdkInit.SetupLogLevel()

	initInfo := &sdkInit.InitInfo{
		ChannelID:     "mychannel",
		ChannelConfig: os.Getenv("GOPATH") + "/src/github.com/hauturier.com/UAVdemo/blockchain/channel-artifacts/channel.tx",

		OrgAdmin: "Admin",
		UserName: "User1",
		OrgName:  "Org1",

		OrdererName: "orderer1.hauturier.com",
		Peer:        "peer0.org1.hauturier.com",

		ChaincodeID:     testcc,
		ChaincodeGoPath: os.Getenv("GOPATH"),
		ChaincodePath:   "/src/github.com/hauturier.com/UAVdemo/blockchain/chaincode/testcc",
	}
	//-----------------------------------------
	//----------------实例化 sdk---------------
	//-----------------------------------------
	fmt.Println("----------------实例化 sdk---------------")
	sdk, err := sdkInit.SetupSDK(configFile, initialized)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	defer sdk.Close()
	//-----------------------------------------
	//------------------创建通道-----------------
	//-----------------------------------------
	fmt.Println("----------------创建通道---------------")
	err = sdkInit.CreateChannel(sdk, initInfo)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//-----------------------------------------
	//------------------加入通道-----------------
	//-----------------------------------------
	fmt.Println("----------------加入通道---------------")
	err = sdkInit.JoinChannel(sdk, initInfo)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = sdkInit.CreateCCLifecycle(sdk,initInfo)
	if err != nil{
		fmt.Println(err)
	}

}
