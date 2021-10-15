/*
 * 资源管理实体相关的API
 */

package fabricSdkApi

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"

	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"go.uber.org/zap"
)

//创建资源管理客户端，只有管理员账号才有权限创建
func CreateResmgmt(sdk *fabsdk.FabricSDK, user string, orgName string) (*resmgmt.Client, error) {
	//clientContext allows creation of transactions using the supplied identity as the credential.
	clientContext := sdk.Context(fabsdk.WithUser(user), fabsdk.WithOrg(orgName))

	// Resource management client is responsible for managing channels (create/update channel)
	// Supply user that has privileges to create channel (in this case orderer admin)
	resMgmtClient, err := resmgmt.New(clientContext)
	if err != nil {
		logger.Error("Failed to create Resmgmt!!!", zap.Error(err))
		panic("Failed to Create Resmgmt!")
	}

	return resMgmtClient, nil
}

func createMspClient(sdk *fabsdk.FabricSDK, organizationName string) (*mspclient.Client, error) {
	mspClient, err := mspclient.New(sdk.Context(), mspclient.WithOrg(organizationName))
	if err != nil {
		logger.Error("Failed to create msp client!!!", zap.Error(err))
		panic("Failed to Create msp client!")
	}
	return mspClient, nil
}
