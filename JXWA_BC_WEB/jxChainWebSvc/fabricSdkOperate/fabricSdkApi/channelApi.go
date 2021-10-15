/*
 * 通道相关的API
 */

package fabricSdkApi

import (
	"fmt"

	"jxChainWebSvc/fabricSdkOperate/util"

	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"go.uber.org/zap"
)

//创建通道对应的客户端
func CreateChannelClient(sdk *fabsdk.FabricSDK, channelID string, user string, orgName string) (*channel.Client, error) {

	//prepare channel client context using client context
	clientChannelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(user), fabsdk.WithOrg(orgName))
	// Channel client is used to query and execute transactions (Org1 is default org)
	client, err := channel.New(clientChannelContext)

	if err != nil {
		logger.Error("Failed to create chaincode!!!", zap.Error(err))
		panic("Failed to package chaincode!")
	}

	return client, nil
}

//根据channelID以及对应的orderer节点名字创建通道
func CreateChannel(sdk *fabsdk.FabricSDK, resMgmtClient *resmgmt.Client, channelID, ordererEndpoint, OrgName, orgAdmin string) {
	mspClient, err := createMspClient(sdk, OrgName)
	if err != nil {
		logger.Error("mspClient create failed!!!", zap.Error(err))
		panic("mspClient create failed!")
	}

	adminIdentity, err := mspClient.GetSigningIdentity(orgAdmin)
	if err != nil {
		logger.Error("get admin identity failed!!!", zap.Error(err))
		panic("get admin identity failed!")
	}
	req := resmgmt.SaveChannelRequest{ChannelID: channelID,
		ChannelConfigPath: util.GetChannelConfigTxPath(channelID + ".tx"),
		SigningIdentities: []msp.SigningIdentity{adminIdentity}}
	txID, err := resMgmtClient.SaveChannel(req, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(ordererEndpoint))
	fmt.Print("txID is %s", txID.TransactionID)
	if err != nil {
		logger.Error("save channel failed!!!", zap.Error(err))
		panic("save channel failed!")
	}
}

func JoinChannel(orgResMgmt *resmgmt.Client, channelID, ordererEndpoint string) error {
	// Org peers join channel
	err := orgResMgmt.JoinChannel(channelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(ordererEndpoint))
	if err != nil {
		logger.Error("JoinChannel channel failed!!!", zap.Error(err))
		panic("JoinChannel channel failed!")
	}

	return nil
}
