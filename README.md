# 说明

## 功能

本项目主要目的是实现两个功能：

1. 以http的形式远程提供数据上链，查询以及访问的操作，还包括链码的生命周期管理，包括链码的打包，安装等步骤。
2. 实现web形式实时展现区块链的一些特征，包括当前的区块hash，高度，链码信息等等。

## 目录结构

主要包含以下三个大的文件夹：

**blockchain-explorer**

主要是提供浏览器的工程，

**JXWA_BC_WEB**

主要是提供HTTP服务的

**test-network**

负责启动fabric网络。

## 具体用法

1. 进入test-network目录：

拉起基本的网络：

```
./network up
```

创建通道：

```
./network createChannel -c mychannel
```

这里会产生通道信息，存在在../channel-artifacts目录下，这里需要拷贝到JXWA_BC_WEB/fixtures/fabric/v2.2/channel目录下以及blockchain-explorer/organizations目录下

把test-network下的organizations 的所有文件拷贝到JXWA_BC_WEB/fixtures/fabric/v1/crypto-config/目录下。

2. 进入JXWA_BC_WEB

```
cd JXWA_BC_WEB/jxChainWebSvc
go run jxchainwebsvc.go -f etc/jxchainwebsvc-api.yaml
```

web框架是通过go zero实现的，具体原因需要参考。

3. 进入blockchain-explorer

```
docker-compose up -d
```
