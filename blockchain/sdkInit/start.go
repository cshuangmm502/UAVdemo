package sdkInit

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/status"
	contextAPI "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/logging"
	mb "github.com/hyperledger/fabric-protos-go/msp"
	fabAPI "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	contextImpl "github.com/hyperledger/fabric-sdk-go/pkg/context"
	mspctx "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/policydsl"
	"github.com/pkg/errors"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	lcpackager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/lifecycle"
)

const (
	ChaincodeVersion = "1.0"
	lvl              = logging.INFO
)

func SetupLogLevel() {
	logging.SetLevel("fabsdk", lvl)
	logging.SetLevel("fabsdk/common", lvl)
	logging.SetLevel("fabsdk/fab", lvl)
	logging.SetLevel("fabsdk/client", lvl)
}

// 实例化 SDK
func SetupSDK(ConfigFile string, initialized bool) (*fabsdk.FabricSDK, error) {

	if initialized {
		return nil, fmt.Errorf("Fabric SDK 已被实例化")
	}

	sdk, err := fabsdk.New(config.FromFile(ConfigFile))
	if err != nil {
		return nil, fmt.Errorf("实例化Fabric SDK失败: %v", err)
	}

	fmt.Println("Fabric SDK 初始化成功")
	return sdk, nil
}

//--------------------------------------------------------------------
// 创建通道
//--------------------------------------------------------------------
func CreateChannel(sdk *fabsdk.FabricSDK, info *InitInfo) error {
	// channel is exists?
	existChannels, err := ListChannel(sdk, *info)
	for _, v := range existChannels {
		if info.ChannelID == v {
			return errors.New("channel exists")
		}
	}
	// 1. 生成客户端上下文环境， 什么身份--> 组织管理员（哪个组织）
	clientContext := sdk.Context(fabsdk.WithUser(info.OrgAdmin), fabsdk.WithOrg(info.OrgName))
	if clientContext == nil {
		return fmt.Errorf("根据指定的组织管理员创建户端Context失败")
	}

	info.orgAdminClientContext = clientContext

	// 2. 根据上下文环境，创建 resMgmtClient, 用来通道的创建，链码的安装、实例化和升级等
	resMgmtClient, err := resmgmt.New(clientContext)
	if err != nil {
		return fmt.Errorf("根据指定的资源管理客户端Context创建通道管理客户端失败: %v", err)
	}

	info.OrgResMgmt = resMgmtClient

	// 3. mspClient 与证书有关的客户端
	mspClient, err := mspclient.New(sdk.Context(), mspclient.WithOrg(info.OrgName))
	if err != nil {
		return fmt.Errorf("根据指定的 OrgName 创建 Org MSP 客户端实例失败: %v", err)
	}

	adminIdentity, err := mspClient.GetSigningIdentity(info.OrgAdmin)
	if err != nil {
		return fmt.Errorf("获取指定id的签名身份失败: %v", err)
	}

	// 生成创建通道请求
	channelReq := resmgmt.SaveChannelRequest{
		ChannelID:         info.ChannelID,
		ChannelConfigPath: info.ChannelConfig,
		SigningIdentities: []mspctx.SigningIdentity{adminIdentity},
	}
	// RC创建通道
	_, err = resMgmtClient.SaveChannel(channelReq,
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
		resmgmt.WithOrdererEndpoint(info.OrdererName),
	)
	if err != nil {
		return errors.Errorf("创建应用通道失败: %v", err)
	}
	fmt.Printf("成功创建通道\n")
	return nil
}

//--------------------------------------------------------------------
// 加入通道
//--------------------------------------------------------------------
func JoinChannel(sdk *fabsdk.FabricSDK, info *InitInfo) error {

	err := info.OrgResMgmt.JoinChannel(
		info.ChannelID,
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
		resmgmt.WithOrdererEndpoint(info.OrdererName),
	)
	if err != nil {
		return fmt.Errorf("Peers加入通道失败: %v", err)
	}

	fmt.Println("peers 已成功加入通道.")

	return nil
}
//---------------------------------------------------------------------------------------
// 安装链码
//--------------------------------------------------------------------------------------
func CreateCCLifecycle(sdk *fabsdk.FabricSDK,info *InitInfo) error{
	// Package cc
	label, ccPkg := packageCC(info)
	packageID := lcpackager.ComputePackageID(label, ccPkg)

	// Install cc
	err := installCC(label, ccPkg, info)
	if err!=nil{
		return err
	}
	// Get installed cc package
	//getInstalledCCPackage(t, packageID, ccPkg, orgResMgmt)

	//// Query installed cc
	//queryInstalled(t, label, packageID, orgResMgmt)
	//
	// Approve cc
	err = approveCC(packageID, info)
	if err!=nil{
		return err
	}
	//// Query approve cc
	//queryApprovedCC(t, orgResMgmt)
	//
	//// Check commit readiness
	//checkCCCommitReadiness(t, orgResMgmt)
	//
	// Commit cc
	err = commitCC(info)
	if err!=nil{
		return err
	}
	//// Query committed cc
	//queryCommittedCC(t, orgResMgmt)
	//
	// Init cc
	err = initCC(sdk,info)
	if err!=nil{
		return err
	}
	return nil
}

//---------------------------------------------------------------------------------------
// 查询已安装链码
//--------------------------------------------------------------------------------------
func QueryInstalledCC(sdk *fabsdk.FabricSDK, info *InitInfo) {

	resp2, err := info.OrgResMgmt.QueryInstalledChaincodes(resmgmt.WithTargetEndpoints(info.Peer))
	if err != nil {
		fmt.Println("查询已安装的链码失败: ", err)
	}

	fmt.Println("已安装链码包括: ", resp2.GetChaincodes())
}


//---------------------------------------------------------------------------------------
// 打包链码
//--------------------------------------------------------------------------------------
func packageCC(sdkInfo *InitInfo) (string, []byte) {
	fmt.Println("**************Begin to package chaincode****************")
	desc := &lcpackager.Descriptor{
		Path:  sdkInfo.ChaincodeGoPath+sdkInfo.ChaincodePath,
		Type:  pb.ChaincodeSpec_GOLANG,
		Label: sdkInfo.ChaincodeID,
	}
	ccPkg, err := lcpackager.NewCCPackage(desc)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("*******************success to create package chaincode*******************")
	return desc.Label, ccPkg
}

//---------------------------------------------------------------------------------------
// 安装链码
//--------------------------------------------------------------------------------------
func installCC(label string, ccPkg []byte, sdkInfo *InitInfo) error{
	fmt.Println("**************Begin to install chaincode****************")
	installCCReq := resmgmt.LifecycleInstallCCRequest{
		Label:   label,
		Package: ccPkg,
	}

	//peers,err := DiscoverLocalPeers(sdkInfo.OrgAdminClientContext,2)
	//if err!=nil{
	//	return fmt.Errorf("failed to discover the local peers by current org adminclientcontext%v",err)
	//}

	_, err := sdkInfo.OrgResMgmt.LifecycleInstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return fmt.Errorf("failed to install chaincode%v",err)
	}
	fmt.Println("*******************success to install chaincode*******************")
	return nil
}

//---------------------------------------------------------------------------------------
// 批准链码
//--------------------------------------------------------------------------------------
func approveCC(packageID string, sdkInfo *InitInfo) error {
	fmt.Println("**************Begin to request to approve chaincode****************")
	peers,err := DiscoverLocalPeers(sdkInfo.orgAdminClientContext,2)
	if err != nil {
		return fmt.Errorf("failed to discover the local peers by current org adminclientcontext%v",err)
	}

	ccPolicy := policydsl.SignedByNOutOfGivenRole(1, mb.MSPRole_MEMBER, []string{"Org1MSP"})
	approveCCReq := resmgmt.LifecycleApproveCCRequest{
		Name:              sdkInfo.ChaincodeID,
		Version:           "0",
		PackageID:         packageID,
		Sequence:          1,
		EndorsementPlugin: "escc",
		ValidationPlugin:  "vscc",
		SignaturePolicy:   ccPolicy,
		InitRequired:      true,
	}

	txnID, err := sdkInfo.OrgResMgmt.LifecycleApproveCC(sdkInfo.ChannelID, approveCCReq, resmgmt.WithTargets(peers...), resmgmt.WithOrdererEndpoint(sdkInfo.OrdererName), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return fmt.Errorf("faile to approve chaincode%v",err)
	}
	fmt.Println("**************Success to request to approve chaincode****************")
	fmt.Printf("the txnID of approve chaincode is: %s",txnID)
	return nil
}

//---------------------------------------------------------------------------------------
// 提交链码
//--------------------------------------------------------------------------------------
func commitCC(sdkInfo *InitInfo) error {
	fmt.Println("**************Begin to commit chaincode****************")
	ccPolicy := policydsl.SignedByNOutOfGivenRole(1, mb.MSPRole_MEMBER, []string{"Org1MSP"})
	req := resmgmt.LifecycleCommitCCRequest{
		Name:              sdkInfo.ChaincodeID,
		Version:           "0",
		Sequence:          1,
		EndorsementPlugin: "escc",
		ValidationPlugin:  "vscc",
		SignaturePolicy:   ccPolicy,
		InitRequired:      true,
	}
	_, err := sdkInfo.OrgResMgmt.LifecycleCommitCC(sdkInfo.ChannelID, req, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(sdkInfo.OrdererName))
	if err != nil {
		return fmt.Errorf("failed to commit chaincode%v",err)
	}
	fmt.Println("**************Success to commit chaincode****************")
	return nil
}

//----------------------------------------------------------------------------------
// 升级链码
//----------------------------------------------------------------------------------
func initCC(sdk *fabsdk.FabricSDK,sdkInfo *InitInfo) error {
	fmt.Println("**************Begin to init chaincode****************")
	//prepare channel client context using client context
	clientChannelContext := sdk.ChannelContext(sdkInfo.ChannelID, fabsdk.WithUser(sdkInfo.UserName), fabsdk.WithOrg(sdkInfo.OrgName))
	// Channel client is used to query and execute transactions (Org1 is default org)
	client, err := channel.New(clientChannelContext)
	if err != nil {
		return fmt.Errorf("Failed to create new channel client: %s", err)
	}

	// init
	_, err = client.Execute(channel.Request{ChaincodeID: sdkInfo.ChaincodeID, Fcn: "init", Args: [][]byte{[]byte("")}, IsInit: true},
		channel.WithRetry(retry.DefaultChannelOpts))
	if err != nil {
		return fmt.Errorf("Failed to init chaincode%s",err)
	}
	fmt.Println("**************Success to init chaincode****************")
	return nil
}

//----------------------------------------------------------------------------------
// 升级链码
//----------------------------------------------------------------------------------
func UpdataCC(info InitInfo) (fabAPI.TransactionID, error) {

	ccPolicy, err := policydsl.FromString("AND ('Org1MSP.member')")
	if err != nil {
		return "", errors.WithMessage(err, "gen policy from string error")
	}

	req := resmgmt.UpgradeCCRequest{
		Name:    info.ChaincodeID,
		Path:    info.ChaincodePath,
		Version: "2",
		Args:    [][]byte{[]byte("init")},
		Policy:  ccPolicy,
	}

	resp, err := info.OrgResMgmt.UpgradeCC(
		info.ChannelID,
		req,
		resmgmt.WithTargetEndpoints(info.Peer),
	)

	if err != nil {
		return "", errors.WithMessage(err, "failed to upgrade chaincode: %s\n")
	}
	if resp.TransactionID == "" {
		return "",errors.New("Failed to upgrade chaincode")
	}

	fmt.Printf("更新链码 %v 成功", info.ChaincodeID)
	return resp.TransactionID, nil
}

//----------------------------------------------------------------------------------
// 登记用户
//----------------------------------------------------------------------------------
func Register(sdk *fabsdk.FabricSDK, info *InitInfo, newIndentity string) (string, error) {

	mspClient, err := mspclient.New(
		sdk.Context(fabsdk.WithUser(info.OrgAdmin), fabsdk.WithOrg(info.OrgName)),
		mspclient.WithOrg(info.OrgName),
		mspclient.WithCAInstance("ca.org1.dragonwell.com"),
	)
	if err != nil {
		return "", err
	}

	enrollmentSecret, err := mspClient.Register(
		&mspclient.RegistrationRequest{
			Name: newIndentity,
			Type: "client",
			Attributes: []mspclient.Attribute{
				{Name: "user2", Value: "true"},
				{Name: "hf.Revoker", Value: "true", ECert: true},
				{Name: "hf.Registrar.Roles", Value: "*"},
				{Name: "hf.Registrar.Attributes", Value: "*"},
				{Name: "hf.Registrar.Attributes", Value: "*"},
				{Name: "GenCRL", Value: "true"},
			},
		})

	if err != nil {
		return "", errors.Errorf("登记 %v 失败, %s", newIndentity, err)
	}
	return enrollmentSecret, nil
}

//----------------------------------------------------------------------------------
// 注册用户
//----------------------------------------------------------------------------------
func Enroll(sdk *fabsdk.FabricSDK, info *InitInfo, enrollmentSecret string) error {

	mspClient, err := mspclient.New(
		sdk.Context(fabsdk.WithUser(info.OrgAdmin), fabsdk.WithOrg(info.OrgName)),
		mspclient.WithOrg(info.OrgName),
		//mspclient.WithCAInstance("ca.org1.dragonwell.com"),
		)
	if err != nil {
		fmt.Printf("创建 mspClient 失败: %v\n", err)
	}

	err = mspClient.Enroll("user2", mspclient.WithSecret(enrollmentSecret))
	if err != nil {
		return errors.Errorf("注册 %v 失败: %v", "user2", err)
	}

	return nil
}

//----------------------------------------------------------------------------------
// 获取 user 信息
//----------------------------------------------------------------------------------
func GetUserInfo(sdk *fabsdk.FabricSDK, userName string, orgID string) (mspctx.SigningIdentity, error) {
	if userName == "" {
		return nil, errors.Errorf("no username specified")
	}

	mspClient, err := mspclient.New(sdk.Context(), mspclient.WithOrg(orgID))
	user, err := mspClient.GetSigningIdentity(userName)
	if err != nil {
		return nil, errors.Errorf("GetSigningIdentity returned error: %v", err)
	}
	fmt.Printf("Returning user [%s], MSPID [%s]\n", user.Identifier().ID, user.Identifier().MSPID)
	return user, nil
}

//----------------------------------------------------------------------------------
// 获取已有通道
//----------------------------------------------------------------------------------
func ListChannel(sdk *fabsdk.FabricSDK, info InitInfo) ([]string, error) {
	var resultChannels []string

	adminContext := sdk.Context(fabsdk.WithUser(info.OrgAdmin), fabsdk.WithOrg(info.OrgName))
	orgResMgmt, err := resmgmt.New(adminContext)
	if err != nil {
		fmt.Printf(err.Error())
		return resultChannels, err
	}
	response, err := orgResMgmt.QueryChannels(resmgmt.WithTargetEndpoints(info.Peer))
	if err != nil {
		fmt.Printf("failed to query channels: %s\n", err)
	}
	allChannels := response.GetChannels()
	for _, channelId := range allChannels {
		resultChannels = append(resultChannels, channelId.GetChannelId())
		fmt.Println(channelId.GetChannelId())
	}
	return resultChannels, nil
}

// DiscoverLocalPeers queries the local peers for the given MSP context and returns all of the peers. If
// the number of peers does not match the expected number then an error is returned.
func DiscoverLocalPeers(ctxProvider contextAPI.ClientProvider, expectedPeers int) ([]fabAPI.Peer, error) {
	ctx, err := contextImpl.NewLocal(ctxProvider)
	if err != nil {
		return nil, errors.Wrap(err, "error creating local context")
	}

	discoveredPeers, err := retry.NewInvoker(retry.New(retry.TestRetryOpts)).Invoke(
		func() (interface{}, error) {
			peers, serviceErr := ctx.LocalDiscoveryService().GetPeers()
			if serviceErr != nil {
				return nil, errors.Wrapf(serviceErr, "error getting peers for MSP [%s]", ctx.Identifier().MSPID)
			}
			if len(peers) < expectedPeers {
				return nil, status.New(status.TestStatus, status.GenericTransient.ToInt32(), fmt.Sprintf("Expecting %d peers but got %d", expectedPeers, len(peers)), nil)
			}
			return peers, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return discoveredPeers.([]fabAPI.Peer), nil
}
