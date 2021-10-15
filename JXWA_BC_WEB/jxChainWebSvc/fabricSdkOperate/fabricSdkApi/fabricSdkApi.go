package fabricSdkApi

import (
	"errors"
	"jxChainWebSvc/fabricSdkOperate/metadata"
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

// const (
// 	channelID      = "mychannel"
// 	orgName        = "Org1"
// 	orgAdmin       = "Admin"
// 	ordererOrgName = "OrdererOrg"
// 	peer1          = "peer0.org1.example.com"
// )

// var (
// 	ccID = "example_cc_fabtest_e2e" + metadata.TestRunID
// )

/*
 *sdk相关操作API
 */

var g_sdk *fabsdk.FabricSDK

func CreateSdk(configOpt core.ConfigProvider, sdkOpts ...fabsdk.Option) (*fabsdk.FabricSDK, error) {
	var err error
	if g_sdk != nil {
		CloseSdk()
		g_sdk, err = fabsdk.New(configOpt, sdkOpts...)
		if err != nil {
			fmt.Println("Create sdk failed!!")
		}
	} else {
		g_sdk, err = fabsdk.New(configOpt, sdkOpts...)
		if err != nil {
			fmt.Println("Create sdk failed!!")
		}
	}

	return g_sdk, err
}

func CloseSdk() {
	if g_sdk == nil {
		panic("close sdk failed!!!, SDK is nil!!!")
	}
	g_sdk.Close()
}

//该接口用于外部获取统一的sdk,sdk的初始化以及后续的web服务调用是分开的
func GetSdk() (*fabsdk.FabricSDK, error) {
	if g_sdk == nil {
		error := errors.New("sdk is nil!!!!")
		return nil, error
	}
	return g_sdk, nil
}

func createChannelAndCC(sdk *fabsdk.FabricSDK, resMgmtClient *resmgmt.Client, client *channel.Client, adminName, orgName, ccID, channelID, targetEndpoint string) {
	//创建系统通道
	CreateChannel(sdk, resMgmtClient, channelID, "orderer.example.com", "Org1", "Admin")

	//创建peer客户端实体
	orgResMgmt, err := CreateResmgmt(sdk, adminName, orgName)
	if err != nil {
		fmt.Println("Failed to create new resource management client: %s", err)
	}
	//peer加入系统通道
	err = JoinChannel(orgResMgmt, channelID, "orderer.example.com")
	if err != nil {
		fmt.Println("join channel failed!!!")
	}

	// Create chaincode package for example cc
	if metadata.CCMode == "lscc" {
		CreateCC(orgResMgmt, channelID, ccID, "github.com/example_cc", []string{"Org1MSP"})
	} else {
		CreateCCLifecycle(orgResMgmt, sdk, client, "example_cc_fabtest_e2e_0", ccID, channelID, "orderer.example.com", targetEndpoint)
	}
}
