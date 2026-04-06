package blockchain

import (
	"sync"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

// 配置信息
var (
	sdk           *fabsdk.FabricSDK                              // Fabric SDK
	channelClient *channel.Client                               // 通道客户端(复用)
	clientMutex   sync.Mutex                                    // 客户端锁
	configPath    = "config.yaml"                                // 配置文件路径
	channelName   = "appchannel"                                 // 通道名称
	user          = "Admin"                                      // 用户
	chainCodeName = "fabric-mims"                                // 链码名称
	endpoints     = []string{"peer0.jd.com", "peer0.taobao.com"} // 要发送交易的节点
)

// Init 初始化
func Init() {
	var err error
	// 通过配置文件初始化SDK
	sdk, err = fabsdk.New(config.FromFile(configPath))
	if err != nil {
		panic(err)
	}
}

// getChannelClient 获取或创建通道客户端(单例模式)
func getChannelClient() (*channel.Client, error) {
	clientMutex.Lock()
	defer clientMutex.Unlock()

	// 如果客户端已存在，直接返回
	if channelClient != nil {
		return channelClient, nil
	}

	// 创建新的通道客户端
	ctx := sdk.ChannelContext(channelName, fabsdk.WithUser(user))
	cli, err := channel.New(ctx)
	if err != nil {
		return nil, err
	}
	channelClient = cli
	return channelClient, nil
}

// resetChannelClient 重置通道客户端(用于连接失败时重试)
func resetChannelClient() {
	clientMutex.Lock()
	defer clientMutex.Unlock()
	channelClient = nil
}

// ChannelExecute 区块链交互
// fcn: 要调用的链码函数名
// args: 传递给链码函数的参数列表
func ChannelExecute(fcn string, args [][]byte) (channel.Response, error) {
	cli, err := getChannelClient()
	if err != nil {
		// 如果获取客户端失败，重置客户端并重试一次
		resetChannelClient()
		cli, err = getChannelClient()
		if err != nil {
			return channel.Response{}, err
		}
	}

	// 对区块链账本的写操作 - 不指定endpoints让SDK自动选择背书节点，避免txid冲突
	resp, err := cli.Execute(channel.Request{
		ChaincodeID: chainCodeName,
		Fcn:         fcn,
		Args:        args,
	})
	if err != nil {
		return channel.Response{}, err
	}
	return resp, nil
}

// ChannelQuery 区块链查询
// fcn: 要查询的链码函数名
// args: 传递给链码函数的参数列表
func ChannelQuery(fcn string, args [][]byte) (channel.Response, error) {
	cli, err := getChannelClient()
	if err != nil {
		// 如果获取客户端失败，重置客户端并重试一次
		resetChannelClient()
		cli, err = getChannelClient()
		if err != nil {
			return channel.Response{}, err
		}
	}

	// 对区块链账本查询的操作 - 不指定endpoints避免txid冲突
	resp, err := cli.Query(channel.Request{
		ChaincodeID: chainCodeName,
		Fcn:         fcn,
		Args:        args,
	})
	if err != nil {
		return channel.Response{}, err
	}
	return resp, nil
}
