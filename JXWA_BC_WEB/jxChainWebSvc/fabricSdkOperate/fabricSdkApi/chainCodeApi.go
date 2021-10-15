/*
 *链码相关API
 */
package fabricSdkApi

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"

	"jxChainWebSvc/fabricSdkOperate/util"

	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/policydsl"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"

	pb "github.com/hyperledger/fabric-protos-go/peer"
	packager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	lcpackager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/lifecycle"
	"go.uber.org/zap"
)

/*
 * 链码相关操作
 */
// CreateCCLifecycle package cc, install cc, get installed cc package, query installed cc
// approve cc, query approve cc, check commit readiness, commit cc, query committed cc
func CreateCCLifecycle(orgResMgmt *resmgmt.Client, sdk *fabsdk.FabricSDK, client *channel.Client, labelName, ccID, channelID, ordererEndpoint, peerDomain string) {
	// Package cc
	label, ccPkg := packageCC(labelName)
	packageID := lcpackager.ComputePackageID(label, ccPkg)

	// Install cc
	InstallCC(label, ccPkg, orgResMgmt)

	// Get installed cc package
	GetInstalledCCPackage(packageID, ccPkg, orgResMgmt, peerDomain)

	// Query installed cc
	QueryInstalled(label, packageID, orgResMgmt, peerDomain)

	// Approve cc
	ApproveCC(packageID, orgResMgmt, ccID, channelID, peerDomain, ordererEndpoint, []string{"Org1MSP","Org2MSP"})

	// Query approve cc
	QueryApprovedCC(orgResMgmt, ccID, channelID, peerDomain)

	// Check commit readiness
	CheckCCCommitReadiness(orgResMgmt, ccID, channelID, peerDomain, []string{"Org1MSP","Org2MSP"})

	// Commit cc
	CommitCC(orgResMgmt, ccID, channelID, peerDomain, ordererEndpoint, []string{"Org1MSP","Org2MSP"})

	// Query committed cc
	QueryCommittedCC(orgResMgmt, ccID, channelID, peerDomain)

	// Init cc
	InitCC(sdk, client, ccID)

	logger.Info("Failed to create chaincode!!!")
}

func CreateCC(orgResMgmt *resmgmt.Client, channelID, ccID, chaincodePath string, signedOrgMSPName []string) {
	ccPkg, err := packager.NewCCPackage(chaincodePath, util.GetDeployPath())
	if err != nil {
		logger.Error("Failed to package chaincode!!!", zap.Error(err))
		panic("Failed to package chaincode!")
	}
	// Install example cc to org peers
	installCCReq := resmgmt.InstallCCRequest{Name: ccID, Path: chaincodePath, Version: "0", Package: ccPkg}
	_, err = orgResMgmt.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		logger.Error("Failed to install chaincode!!!", zap.Error(err))
		panic("Failed to install chaincode!")
	}
	// Set up chaincode policy
	ccPolicy := policydsl.SignedByAnyMember(signedOrgMSPName)
	// Org resource manager will instantiate 'example_cc' on channel
	_, err = orgResMgmt.InstantiateCC(
		channelID,
		resmgmt.InstantiateCCRequest{Name: ccID, Path: chaincodePath, Version: "0", Args: util.ExampleCCInitArgs(), Policy: ccPolicy},
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
	)
	//fmt.Print("resp is %s", resp.TransactionID)
	if err != nil {
		logger.Error("Failed to instantiate chaincode!!!", zap.Error(err))
		panic("Failed to instantiate chaincode!")
	}
}

//链码打包
func packageCC(labelName string) (string, []byte) {
	desc := &lcpackager.Descriptor{
		Path:  util.GetLcDeployPath(),
		Type:  pb.ChaincodeSpec_GOLANG,
		Label: labelName,
	}
	ccPkg, err := lcpackager.NewCCPackage(desc)
	if err != nil {
		logger.Error("Failed to package chaincode!!!", zap.Error(err))
		panic("Failed to package chaincode!")
	}
	return desc.Label, ccPkg
}

//安装链码
func InstallCC(label string, ccPkg []byte, orgResMgmt *resmgmt.Client) {
	installCCReq := resmgmt.LifecycleInstallCCRequest{
		Label:   label,
		Package: ccPkg,
	}

	packageID := lcpackager.ComputePackageID(installCCReq.Label, installCCReq.Package)

	resp, err := orgResMgmt.LifecycleInstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		//logger.Error("Failed to install chaincode!!!", zap.Error(err))
		//panic("Failed to install chaincode!")
		fmt.Println("err is %s", err)
	}

	if packageID != resp[0].PackageID {
		panic("error: packageID not equal!!!!")
	}
}

//根据链码包ID查询是否已安装
func GetInstalledCCPackage(packageID string, ccPkg []byte, orgResMgmt *resmgmt.Client, targetEndpoint string) {
	_, err := orgResMgmt.LifecycleGetInstalledCCPackage(packageID, resmgmt.WithTargetEndpoints(targetEndpoint), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		logger.Error("Failed to GetInstalled chaincode!!!", zap.Error(err))
		panic("Failed to GetInstalled chaincode!")
	}
}

//查询已安装的链码
func QueryInstalled(label string, packageID string, orgResMgmt *resmgmt.Client, targetEndpoint string) {
	resp, err := orgResMgmt.LifecycleQueryInstalledCC(resmgmt.WithTargetEndpoints(targetEndpoint), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		logger.Error("Failed to query chaincode!!!", zap.Error(err))
		panic("Failed to query chaincode!")
	}
	if packageID != resp[0].PackageID {
		panic("error: query installed chaincode package ID not equal !!!!")

	}
	if label != resp[0].Label {
		panic("error: query installed chaincode label not equal !!!!")
	}
}

//批注链码
func ApproveCC(packageID string, orgResMgmt *resmgmt.Client, ccID, channelID, targetEndpoint, ordererEndpoint string, signedOrgMSPName []string) {
	ccPolicy := policydsl.SignedByAnyMember(signedOrgMSPName)
	approveCCReq := resmgmt.LifecycleApproveCCRequest{
		Name:              ccID,
		Version:           "0",
		PackageID:         packageID,
		Sequence:          1,
		EndorsementPlugin: "escc",
		ValidationPlugin:  "vscc",
		SignaturePolicy:   ccPolicy,
		InitRequired:      true,
	}

	txnID, err := orgResMgmt.LifecycleApproveCC(channelID, approveCCReq, resmgmt.WithTargetEndpoints(targetEndpoint), resmgmt.WithOrdererEndpoint(ordererEndpoint), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		logger.Error("Failed to approve chaincode!!!", zap.Error(err))
		panic("Failed to approve chaincode!")
	}

	if txnID == "" {
		panic("error: ApproveCC failed !!!!")
	}
	//fmt.Printf("txnid is %s \n", txnID)
}

//查询已批准的链码
func QueryApprovedCC(orgResMgmt *resmgmt.Client, ccID, channelID, targetEndpoint string) {
	queryApprovedCCReq := resmgmt.LifecycleQueryApprovedCCRequest{
		Name:     ccID,
		Sequence: 1,
	}
	resp, err := orgResMgmt.LifecycleQueryApprovedCC(channelID, queryApprovedCCReq, resmgmt.WithTargetEndpoints(targetEndpoint), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		logger.Error("Failed to query approved chaincode!!!", zap.Error(err))
		panic("Failed to query approved chaincode!")
	}
	if resp.Name == "" && resp.Version == "" {
		panic("error: ApproveCC failed !!!!")
	}
}

// Check commit readiness
func CheckCCCommitReadiness(orgResMgmt *resmgmt.Client, ccID, channelID, targetEndpoints string, signedOrgMSPName []string) {
	ccPolicy := policydsl.SignedByAnyMember(signedOrgMSPName)
	req := resmgmt.LifecycleCheckCCCommitReadinessRequest{
		Name:              ccID,
		Version:           "0",
		EndorsementPlugin: "escc",
		ValidationPlugin:  "vscc",
		SignaturePolicy:   ccPolicy,
		Sequence:          1,
		InitRequired:      true,
	}
	resp, err := orgResMgmt.LifecycleCheckCCCommitReadiness(channelID, req, resmgmt.WithTargetEndpoints(targetEndpoints), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		logger.Error("Failed to check chaincode commit readiness!!!", zap.Error(err))
		panic("Failed to check chaincode commit readiness!")
	}

	if len(resp.Approvals) == 0 {
		panic("error: CheckCCCommitReadiness failed !!!!")
	}
	//fmt.Printf("resp is %s \n", resp)
}

func CommitCC(orgResMgmt *resmgmt.Client, ccID, channelID, targetEndpoints, ordererEndpoint string, signedOrgMSPName []string) {
	ccPolicy := policydsl.SignedByAnyMember(signedOrgMSPName)
	req := resmgmt.LifecycleCommitCCRequest{
		Name:              ccID,
		Version:           "0",
		Sequence:          1,
		EndorsementPlugin: "escc",
		ValidationPlugin:  "vscc",
		SignaturePolicy:   ccPolicy,
		InitRequired:      true,
	}
	txnID, err := orgResMgmt.LifecycleCommitCC(channelID, req, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithTargetEndpoints(targetEndpoints), resmgmt.WithOrdererEndpoint(ordererEndpoint))
	if err != nil {
		logger.Error("Failed to commit chaincode!!!", zap.Error(err))
		panic("Failed to commit chaincode!")
	}
	if txnID == "" {
		panic("error: CommitCC failed !!!!")
	}
	//fmt.Printf("txnID is %s \n", txnID)
}

func QueryCommittedCC(orgResMgmt *resmgmt.Client, ccID, channelID, targetEndpoints string) {
	req := resmgmt.LifecycleQueryCommittedCCRequest{
		Name: ccID,
	}
	resp, err := orgResMgmt.LifecycleQueryCommittedCC(channelID, req, resmgmt.WithTargetEndpoints(targetEndpoints), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		logger.Error("Failed to query commit chaincode!!!", zap.Error(err))
		panic("Failed to query commit chaincode!")
	}

	if ccID != resp[0].Name {
		panic("error: QueryCommittedCC failed, ccID not equal!!")
	}
}

func InitCC(sdk *fabsdk.FabricSDK, client *channel.Client, ccID string) {
	//prepare channel client context using client context
	// clientChannelContext := sdk.ChannelContext(channelID, fabsdk.WithUser("User1"), fabsdk.WithOrg(orgName))
	// Channel client is used to query and execute transactions (Org1 is default org)
	// client, err := channel.New(clientChannelContext)
	// if err != nil {
	// 	fmt.Println("Failed to create new channel client: %s", err)
	// }

	// init
	_, err := client.Execute(channel.Request{ChaincodeID: ccID, Fcn: "init", Args: util.ExampleCCInitArgsLc(), IsInit: true},
		channel.WithRetry(retry.DefaultChannelOpts))
	if err != nil {
		logger.Error("Failed to init chaincode!!!", zap.Error(err))
		panic("Failed to init chaincode!")
	}
}

//链码执行以及链码查询
func QueryCC(client *channel.Client, chainCodeID string, args [][]byte, targetEndpoints ...string) []byte {
	response, err := client.Query(channel.Request{ChaincodeID: chainCodeID, Fcn: "invoke", Args: args},
		channel.WithRetry(retry.DefaultChannelOpts),
		channel.WithTargetEndpoints(targetEndpoints...),
	)
	if err != nil {
		logger.Error("Failed to query chaincode!!!", zap.Error(err))
		panic("Failed to query chaincode!")
	}
	return response.Payload
}

func ExecuteCC(client *channel.Client, ccID string, args [][]byte) (string, error) {
	resp, err := client.Execute(channel.Request{ChaincodeID: ccID, Fcn: "invoke", Args: args},
		channel.WithRetry(retry.DefaultChannelOpts))
	if err != nil {
		logger.Error("Failed to execute chaincode!!!", zap.Error(err))
		panic("Failed to execute chaincode!")
	}
	result := string(resp.TransactionID)
	//这里主要是返回交易ID
	return result, err
}
