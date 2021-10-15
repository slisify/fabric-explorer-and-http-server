package fabricSdkApi

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"jxChainWebSvc/fabricSdkOperate/metadata"
	"jxChainWebSvc/fabricSdkOperate/util"

	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"go.uber.org/zap"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

const (
	channelID      = "mychannel"
	orgName        = "Org1"
	orgAdmin       = "Admin"
	ordererOrgName = "OrdererOrg"
	peer1          = "peer0.org1.example.com"
	ordererDomain  = "orderer.example.com"
	ccLable        = "example_cc_fabtest_e2e_0"
	//congfigName    = "config_e2e.yaml"
	congfigName    = "config_test.yaml"
)

var (
	ccID = "example_cc_fabtest_e2e" + metadata.TestRunID
)

//对外的功能应该主要包含以下：
//1. 创建SDK
//func SetupAndRun(isCreateChannel bool, sdkOpts ...fabsdk.Option)
//2. 关闭SDK
//3. 创建通道，这里应该不需要暴露，
//4. 加入通道，应该也不需要暴露
//5. 上传链码
//6. 部署链码
//7. 调用链码
//8. 查询链码
//9. 初始化

//API1:SDK的初始化以及相关的初始通道创建
func SetupAndRun(isCreateChannel bool, sdkOpts ...fabsdk.Option) {
	//配置环境路径等
	//获得当天的文件目录
	str, _ := os.Getwd()
	//web服务主程序的主要目录
	webPathStr := str + "../../"
	os.Setenv("FABRIC_SDK_GO_PROJECT_PATH", webPathStr)
	//配置环境变量，设置主目录，方便后面进行定位

	//记得优先初始化记录日志
	InitLogger()

	configPath := util.GetConfigPath(congfigName)
	configOpt := config.FromFile(configPath)

	sdk, err := CreateSdk(configOpt, sdkOpts...)
	if err != nil {
		fmt.Println("Failed to create new SDK: %s", err)
	}

	//创建资源管理客户端实体
	resMgmtClient, err := CreateResmgmt(sdk, orgAdmin, ordererOrgName)
	if err != nil {
		fmt.Println("resMgmtClient create failed!!!")
	}

	//sdk, err = GetSdk()

	if isCreateChannel {
		//创建系统通道
		CreateChannel(sdk, resMgmtClient, channelID, ordererDomain, orgName, orgAdmin)

		//创建peer客户端实体，并加入通道
		orgResMgmt, err := PeerJoinChannel(orgAdmin, orgName, channelID, ordererDomain)
		if err != nil {
			panic("create orgResMgmt failed!!!!!!!")
		}

		//创建channel的客户端，方便后续调用，这是用户的通道
		client, err := CreateChannelClient(sdk, channelID, "User1", orgName)
		if err != nil {
			panic("create client failed!!!!!!!")
		}

		//链码的生命周期管理，包括打包，安装，同意，提交等
		// if metadata.CCMode == "lscc" {
		// 	ccPath := util.GetDeployPath()
		// 	fabricSdkApi.CreateCC(orgResMgmt, channelID, ccID, ccPath, []string{"Org1MSP"})
		// } else {
		// 	fabricSdkApi.CreateCCLifecycle(orgResMgmt, sdk, client, ccLable, ccID, channelID, ordererDomain, peer1)
		// }
		CreateCCLifecycle(orgResMgmt, sdk, client, ccLable, ccID, channelID, ordererDomain, peer1)
	}
}

//API2:节点加入对应的通道
func PeerJoinChannel(orgAdmin, orgName, channelID, ordererDomain string) (*resmgmt.Client, error) {
	sdk, err := GetSdk()
	if err != nil {
		panic("error:get sdk failed!")
	}
	//创建peer客户端实体
	orgResMgmt, err := CreateResmgmt(sdk, orgAdmin, orgName)
	if err != nil {
		fmt.Println("Failed to create new resource management client: %s", err)
	}
	//peer加入系统通道
	err = JoinChannel(orgResMgmt, channelID, ordererDomain)
	if err != nil {
		fmt.Println("join channel failed!!!")
	}

	return orgResMgmt, err
}

//API3:调用链码
func ExecuteChainCode(QueryData [][]byte) (string, error) {
	sdk, err := GetSdk()
	if sdk == nil {
		panic("error sdk is not exist!")
	}
	client, err := CreateChannelClient(sdk, channelID, "User1", orgName)
	// client, err := CreateChannelClient(sdk, channelID, "User1", orgName)
	if err != nil {
		//panic("Failed to create new channel client: %s", err)
		panic("Failed to create new channel client")
	}

	//拼接入参以及默认的链码执行操作
	executeData := jointCCArgs(util.FireAlarmArgs, QueryData)

	fmt.Println("now print the data")
	// for i := range executeData {
	// 	//fmt.Println("%s", val)
	// 	logger.Error("invoke args!!", zap.String("str:", string(executeData[i][:])))
	// }

	//查询是否存在对应的数据
	//existingValue := QueryCC(client, ccID, QueryData)

	eventID := "test([a-zA-Z]+)"

	// Register chaincode event (pass in channel which receives event details when the event is complete)
	reg, notifier, err := client.RegisterChaincodeEvent(ccID, eventID)
	if err != nil {
		panic("Failed to register cc event")
	}
	defer client.UnregisterChaincodeEvent(reg)
	var ccEvent *fab.CCEvent

	//执行相关的链码，应该要有相关的操作把web传来的数据解析成对应的格式
	transactionID, err := ExecuteCC(client, ccID, executeData)
	if err != nil {
		panic("error: execute failed!!")
	}

	select {
	case ccEvent = <-notifier:
		fmt.Println("Received CC event: %#v\n", ccEvent)
	case <-time.After(time.Second * 20):
		fmt.Println("Did NOT receive CC event for eventId(%s)\n", eventID)
	}

	//检查执行的结果是否符合预期
	// newValue := QueryCC(client, ccID, QueryData, ccEvent.SourceURL)
	// valueInt, err := strconv.Atoi(string(existingValue))
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// valueAfterInvokeInt, err := strconv.Atoi(string(newValue))
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// if valueInt+1 != valueAfterInvokeInt {
	// 	fmt.Println("Execute failed. Before: %s, after: %s", existingValue, newValue)
	// }

	return transactionID, nil
}

//API4:查询数据
func QueryData(queryData [][]byte) []byte {
	sdk, err := GetSdk()
	if sdk == nil {
		panic("error sdk is not exist!")
	}
	client, err := CreateChannelClient(sdk, channelID, "User1", orgName)
	// client, err := CreateChannelClient(sdk, channelID, "User1", orgName)
	if err != nil {
		fmt.Println("Failed to create new channel client: %s", err)
	}

	//查询是否存在对应的数据
	executeData := jointCCArgs(util.QueryAlarmArgs, queryData)
	existingValue := QueryCC(client, ccID, executeData)

	return existingValue
}

//API5:删除数据
func DeleteData(deleteData [][]byte) string {
	sdk, err := GetSdk()
	if sdk == nil {
		panic("error sdk is not exist!")
	}
	client, err := CreateChannelClient(sdk, channelID, "User1", orgName)
	if err != nil {
		fmt.Println("Failed to create new channel client: %s", err)
	}

	//查询是否存在对应的数据
	executeData := jointCCArgs(util.DeleteAlarmArgs, deleteData)
	deleteResult := QueryCC(client, ccID, executeData)
	if deleteResult == nil {
		return "ok"
	} else {
		return "failed"
	}
}


//内部函数，用于连接http的参数以及对应的链码操作
func jointCCArgs(funcArgs []byte, ccArgs [][]byte) [][]byte {
	onChhainDataLen := len(ccArgs)
	logger.Error("input len!!", zap.String("str:", strconv.Itoa(onChhainDataLen)))
	//invokeArgsLen := len(util.FireAlarmArgs)
	//var totolLen int
	j := 0
	//totolLen = onChhainDataLen + invokeArgsLen
	executeData := make([][]byte, onChhainDataLen+1)
	for i := range executeData {
		if i == 0 {
			len := len(funcArgs)
			executeData[i] = make([]byte, len)
			copy(executeData[i], funcArgs)
			logger.Error("invoke args!!", zap.String("str:", string(executeData[i][:])))
		} else {
			len := len(ccArgs[j])
			logger.Error("invoke args len!!", zap.Int("len:", len))
			logger.Error("invoke args onchaindata!!", zap.String("data:", string(ccArgs[j][:])))
			executeData[i] = make([]byte, len)
			copy(executeData[i], ccArgs[j])
			j++
			logger.Error("invoke args!!", zap.String("str:", string(executeData[i][:])))
		}
	}

	return executeData
}
