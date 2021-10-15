package main

import (
	"flag"
	"fmt"

	"jxChainWebSvc/fabricSdkOperate/fabricSdkApi"
	_ "jxChainWebSvc/fabricSdkOperate/util"
	"jxChainWebSvc/internal/config"
	"jxChainWebSvc/internal/handler"
	"jxChainWebSvc/internal/svc"

	"github.com/tal-tech/go-zero/core/conf"
	"github.com/tal-tech/go-zero/rest"
)

var configFile = flag.String("f", "etc/jxChainWebSvc-api.yaml", "the config file")

func main() {
	//第二阶段，先查询是否存在未关闭的SDK，有则先关掉，然后再打开
	fabricSdkApi.SetupAndRun(true)
	defer fabricSdkApi.CloseSdk()
	fmt.Printf("sdk start ok!!!!!!")

	//第一阶段把web服务起起来
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	ctx := svc.NewServiceContext(c)
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
