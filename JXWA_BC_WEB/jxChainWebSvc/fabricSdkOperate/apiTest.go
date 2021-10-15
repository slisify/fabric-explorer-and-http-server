package main

import (
	"fabricSdkOperate/fabricSdkApi"
	"fabricSdkOperate/metadata"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/cryptosuite"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	pkgmsp "github.com/hyperledger/fabric-sdk-go/pkg/msp"
)

const (
	channelID      = "mychannel"
	orgName        = "Org1"
	orgAdmin       = "Admin"
	ordererOrgName = "OrdererOrg"
	peer1          = "peer0.org1.example.com"
	ordererDomain  = "orderer.example.com"
	ccLable        = "example_cc_fabtest_e2e_0"
)

var (
	ccID = "example_cc_fabtest_e2e" + metadata.TestRunID
)

func main() {
	fabricSdkApi.InitLogger()
	fabricSdkApi.SetupAndRun(true)
	sdk, _ := fabricSdkApi.GetSdk()
	if sdk != nil {
		cleanupUserData(sdk)
		defer cleanupUserData(sdk)
		e2e(sdk)
	}
	fabricSdkApi.CloseSdk()
}

// type SDKFunc func(sdk *fabsdk.FabricSDK)

// func setupAndRun(isCreateChannel bool, configOpt core.ConfigProvider, operation SDKFunc, sdkOpts ...fabsdk.Option) {
// 	sdk, err := fabricSdkApi.CreateSdk(configOpt, sdkOpts...)
// 	if err != nil {
// 		fmt.Println("Failed to create new SDK: %s", err)
// 	}
// 	defer fabricSdkApi.CloseSdk(sdk)

// 	// Delete all private keys from the crypto suite store
// 	// and users from the user store at the end
// 	cleanupUserData(sdk)
// 	defer cleanupUserData(sdk)

// 	//创建资源管理客户端实体
// 	resMgmtClient, err := fabricSdkApi.CreateResmgmt(sdk, orgAdmin, ordererOrgName)
// 	if err != nil {
// 		fmt.Println("resMgmtClient create failed!!!")
// 	}

// 	sdk, err = fabricSdkApi.GetSdk()

// 	if isCreateChannel {
// 		//创建系统通道
// 		fabricSdkApi.CreateChannel(sdk, resMgmtClient, channelID, ordererDomain, orgName, orgAdmin)

// 		//创建peer客户端实体
// 		orgResMgmt, err := fabricSdkApi.CreateResmgmt(sdk, orgAdmin, orgName)
// 		if err != nil {
// 			fmt.Println("Failed to create new resource management client: %s", err)
// 		}
// 		//peer加入系统通道
// 		err = fabricSdkApi.JoinChannel(orgResMgmt, channelID, ordererDomain)
// 		if err != nil {
// 			fmt.Println("join channel failed!!!")
// 		}

// 		//创建channel的客户端，方便后续调用，这是用户的通道
// 		client, err := fabricSdkApi.CreateChannelClient(sdk, channelID, "User1", orgName)

// 		//链码的生命周期管理，包括打包，安装，同意，提交等
// 		if metadata.CCMode == "lscc" {
// 			ccPath := util.GetDeployPath()
// 			fabricSdkApi.CreateCC(orgResMgmt, channelID, ccID, ccPath, []string{"Org1MSP"})
// 		} else {
// 			fabricSdkApi.CreateCCLifecycle(orgResMgmt, sdk, client, ccLable, ccID, channelID, ordererDomain, peer1)
// 		}
// 	}

// 	//实际链码功能测试
// 	operation(sdk)
// }

func cleanupUserData(sdk *fabsdk.FabricSDK) {
	var keyStorePath, credentialStorePath string

	configBackend, err := sdk.Config()
	if err != nil {
		// if an error is returned from Config, it means configBackend was nil, in this case simply hard code
		// the keyStorePath and credentialStorePath to the default values
		// This case is mostly happening due to configless test that is not passing a ConfigProvider to the SDK
		// which makes configBackend = nil.
		// Since configless test uses the same config values as the default ones (config_test.yaml), it's safe to
		// hard code these paths here
		keyStorePath = "/tmp/msp/keystore"
		credentialStorePath = "/tmp/state-store"
	} else {
		cryptoSuiteConfig := cryptosuite.ConfigFromBackend(configBackend)
		identityConfig, err := pkgmsp.ConfigFromBackend(configBackend)
		if err != nil {
			fmt.Println(err)
		}

		keyStorePath = cryptoSuiteConfig.KeyStorePath()
		credentialStorePath = identityConfig.CredentialStorePath()
	}

	cleanupPath(keyStorePath)
	cleanupPath(credentialStorePath)
}

// CleanupTestPath removes the contents of a state store.
func cleanupPath(storePath string) {
	err := os.RemoveAll(storePath)
	if err != nil {
		fmt.Println("Cleaning up directory '%s' failed: %v", storePath, err)
	}
}

//测试点对点的交易
func e2e(sdk *fabsdk.FabricSDK) {
	//prepare channel client context using client context
	clientChannelContext := sdk.ChannelContext(channelID, fabsdk.WithUser("User1"), fabsdk.WithOrg(orgName))
	// Channel client is used to query and execute transactions (Org1 is default org)
	client, err := channel.New(clientChannelContext)

	// client, err := CreateChannelClient(sdk, channelID, "User1", orgName)
	if err != nil {
		fmt.Println("Failed to create new channel client: %s", err)
	}

	existingValue := fabricSdkApi.QueryCC(client, ccID)
	ccEvent := moveFunds(client)

	// Verify move funds transaction result on the same peer where the event came from.
	verifyFundsIsMoved(client, existingValue, ccEvent)
	fmt.Println("Chaincode invoke SUCCESSS!!!!!!!!!!!!!!!!!!!!!!!!")
}

func moveFunds(client *channel.Client) *fab.CCEvent {

	eventID := "test([a-zA-Z]+)"

	// Register chaincode event (pass in channel which receives event details when the event is complete)
	reg, notifier, err := client.RegisterChaincodeEvent(ccID, eventID)
	if err != nil {
		fmt.Println("Failed to register cc event: %s", err)
	}
	defer client.UnregisterChaincodeEvent(reg)

	// Move funds
	fabricSdkApi.ExecuteCC(client, ccID)

	var ccEvent *fab.CCEvent
	select {
	case ccEvent = <-notifier:
		fmt.Println("Received CC event: %#v\n", ccEvent)
	case <-time.After(time.Second * 20):
		fmt.Println("Did NOT receive CC event for eventId(%s)\n", eventID)
	}

	return ccEvent
}

func verifyFundsIsMoved(client *channel.Client, value []byte, ccEvent *fab.CCEvent) {
	newValue := fabricSdkApi.QueryCC(client, ccID, ccEvent.SourceURL)
	valueInt, err := strconv.Atoi(string(value))
	if err != nil {
		fmt.Println(err.Error())
	}
	valueAfterInvokeInt, err := strconv.Atoi(string(newValue))
	if err != nil {
		fmt.Println(err.Error())
	}
	if valueInt+1 != valueAfterInvokeInt {
		fmt.Println("Execute failed. Before: %s, after: %s", value, newValue)
	}
}
