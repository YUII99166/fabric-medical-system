# Hyperledger Fabric SDK 使用指南

## 一、SDK 基本概念

Fabric SDK 是用于与 Hyperledger Fabric 区块链网络交互的客户端库。本项目使用 **Go SDK**。

## 二、项目中的 SDK 配置

### 1. 配置文件位置
```
application/server/config-local-dev.yaml  # 本地开发配置
application/server/config.yaml            # 生产环境配置
```

### 2. SDK 初始化代码
文件：`application/server/blockchain/sdk.go`

```go
package blockchain

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

var (
	sdk           *fabsdk.FabricSDK
	channelName   = "appchannel"      // 通道名称
	user          = "Admin"            // 用户身份
	chainCodeName = "fabric-mims"      // 链码名称
	endpoints     = []string{"peer0.jd.com", "peer0.taobao.com"} // 目标节点
	configPath    = "./config-local-dev.yaml"
)

// Init 初始化 SDK
func Init() {
	var err error
	sdk, err = fabsdk.New(config.FromFile(configPath))
	if err != nil {
		panic(err)
	}
}
```

## 三、SDK 使用方法

### 1. 查询区块链（只读操作）

**函数签名：**
```go
func ChannelQuery(fcn string, args [][]byte) (channel.Response, error)
```

**参数说明：**
- `fcn`: 链码函数名（如 "queryAccountV2List"）
- `args`: 函数参数列表（字节数组的数组）

**示例1：查询所有用户**
```go
package main

import (
	bc "application/blockchain"
	"encoding/json"
	"fmt"
)

func QueryAllUsers() {
	// 调用链码的 queryAccountV2List 函数
	// 不需要参数，传空数组
	resp, err := bc.ChannelQuery("queryAccountV2List", [][]byte{})
	if err != nil {
		fmt.Printf("查询失败: %v\n", err)
		return
	}
	
	// 解析返回的数据
	var userList []map[string]interface{}
	err = json.Unmarshal(resp.Payload, &userList)
	if err != nil {
		fmt.Printf("解析数据失败: %v\n", err)
		return
	}
	
	fmt.Printf("查询到 %d 个用户\n", len(userList))
	for _, user := range userList {
		fmt.Printf("用户: %v\n", user)
	}
}
```

**示例2：查询指定病人的病历**
```go
func QueryPrescriptionsByPatient(patientID string) {
	// 准备参数
	var args [][]byte
	args = append(args, []byte(patientID))
	
	// 调用链码函数
	resp, err := bc.ChannelQuery("queryPrescription", args)
	if err != nil {
		fmt.Printf("查询失败: %v\n", err)
		return
	}
	
	// 解析结果
	var prescriptions []map[string]interface{}
	json.Unmarshal(resp.Payload, &prescriptions)
	
	fmt.Printf("病人 %s 有 %d 条病历\n", patientID, len(prescriptions))
}
```

### 2. 写入区块链（修改操作）

**函数签名：**
```go
func ChannelExecute(fcn string, args [][]byte) (channel.Response, error)
```

**示例1：创建用户**
```go
func CreateUser(accountData map[string]interface{}) {
	// 将数据序列化为 JSON
	jsonData, err := json.Marshal(accountData)
	if err != nil {
		fmt.Printf("序列化失败: %v\n", err)
		return
	}
	
	// 准备参数
	var args [][]byte
	args = append(args, jsonData)
	
	// 调用链码函数
	resp, err := bc.ChannelExecute("createAccountV2", args)
	if err != nil {
		fmt.Printf("创建用户失败: %v\n", err)
		return
	}
	
	fmt.Printf("创建成功，交易ID: %s\n", resp.TransactionID)
}
```

**示例2：创建病历**
```go
func CreatePrescription(prescription map[string]interface{}) {
	// 序列化数据
	jsonData, _ := json.Marshal(prescription)
	
	// 调用链码
	var args [][]byte
	args = append(args, jsonData)
	
	resp, err := bc.ChannelExecute("createPrescription", args)
	if err != nil {
		fmt.Printf("创建病历失败: %v\n", err)
		return
	}
	
	fmt.Printf("病历创建成功\n")
	fmt.Printf("交易ID: %s\n", resp.TransactionID)
	fmt.Printf("区块高度: %d\n", resp.BlockNumber)
}
```

## 四、在 API 中使用 SDK

### 示例：用户查询 API
文件：`application/server/api/v2/account.go`

```go
package v2

import (
	bc "application/blockchain"
	"application/pkg/app"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

// QueryAccountList 查询用户列表
func QueryAccountList(c *gin.Context) {
	appG := app.Gin{C: c}
	
	// 1. 调用区块链查询
	var bodyBytes [][]byte
	resp, err := bc.ChannelQuery("queryAccountV2List", bodyBytes)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	
	// 2. 解析返回数据
	var accountList []map[string]interface{}
	err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &accountList)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	
	// 3. 返回给前端
	appG.Response(http.StatusOK, "成功", accountList)
}
```

### 示例：创建病历 API
```go
func CreatePrescription(c *gin.Context) {
	appG := app.Gin{C: c}
	
	// 1. 接收前端数据
	var prescription map[string]interface{}
	if err := c.ShouldBindJSON(&prescription); err != nil {
		appG.Response(http.StatusBadRequest, "参数错误", err.Error())
		return
	}
	
	// 2. 序列化数据
	jsonData, _ := json.Marshal(prescription)
	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, jsonData)
	
	// 3. 调用区块链
	resp, err := bc.ChannelExecute("createPrescription", bodyBytes)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	
	// 4. 返回结果
	appG.Response(http.StatusOK, "成功", map[string]interface{}{
		"transaction_id": resp.TransactionID,
		"message": "病历创建成功",
	})
}
```

## 五、完整的调用流程

```
┌─────────────┐
│   前端 Vue   │
└──────┬──────┘
       │ HTTP POST /api/v2/createPrescription
       │ Body: { patient_id: "xxx", diagnosis: "xxx" }
       ↓
┌─────────────────┐
│  后端 Gin API    │
│  (Go Server)    │
└──────┬──────────┘
       │ bc.ChannelExecute("createPrescription", args)
       ↓
┌─────────────────┐
│  Fabric SDK     │
│  (Go SDK)       │
└──────┬──────────┘
       │ gRPC 调用
       ↓
┌─────────────────┐
│  Peer 节点       │
│  (peer0.jd.com) │
│  端口: 27051     │
└──────┬──────────┘
       │ 执行链码
       ↓
┌─────────────────┐
│  Chaincode      │
│  (fabric-mims)  │
└──────┬──────────┘
       │ 写入账本
       ↓
┌─────────────────┐
│  区块链账本      │
│  (Ledger)       │
└─────────────────┘
```

## 六、常用链码函数列表

### 用户管理
- `queryAccountV2List` - 查询所有用户
- `createAccountV2` - 创建用户
- `register` - 用户注册

### 病历管理
- `queryPrescription` - 查询病历
- `createPrescription` - 创建病历
- `queryFullMedicalHistory` - 查询完整病历历史

### 补充记录
- `addSupplementRecord` - 添加补充记录
- `querySupplementRecords` - 查询补充记录

### 授权管理
- `requestAccess` - 请求访问授权
- `approveAccess` - 批准授权
- `queryAccessRequests` - 查询授权请求

### 药品订单
- `createDrugOrder` - 创建药品订单
- `queryDrugOrder` - 查询药品订单

## 七、错误处理

```go
resp, err := bc.ChannelQuery("queryAccountV2List", [][]byte{})
if err != nil {
	// 可能的错误：
	// 1. 网络连接失败
	// 2. 节点不可用
	// 3. 链码函数不存在
	// 4. 参数错误
	// 5. 权限不足
	fmt.Printf("区块链调用失败: %v\n", err)
	return
}

// 检查返回状态
if resp.ChaincodeStatus != 200 {
	fmt.Printf("链码执行失败: %s\n", string(resp.Payload))
	return
}
```

## 八、调试技巧

### 1. 查看 SDK 日志
在配置文件中设置日志级别：
```yaml
client:
  logging:
    level: debug  # info, debug, warning, error
```

### 2. 查看节点日志
```bash
# 查看 peer 节点日志
docker logs peer0.jd.com

# 实时查看日志
docker logs -f peer0.jd.com
```

### 3. 测试链码函数
```bash
# 进入 CLI 容器
docker exec -it cli bash

# 查询测试
peer chaincode query \
  -C appchannel \
  -n fabric-mims \
  -c '{"Args":["queryAccountV2List"]}'

# 调用测试
peer chaincode invoke \
  -C appchannel \
  -n fabric-mims \
  -c '{"Args":["createAccountV2", "{...}"]}'
```

## 九、性能优化

### 1. 使用多个 Peer 节点
```go
endpoints := []string{
	"peer0.jd.com",
	"peer0.taobao.com",
	"peer0.wenjin.com",
}
```

### 2. 连接池配置
在配置文件中调整：
```yaml
peers:
  peer0.jd.com:
    grpcOptions:
      keep-alive-time: 10s
      keep-alive-timeout: 20s
```

### 3. 批量操作
尽量减少区块链调用次数，一次调用处理多条数据。

## 十、安全注意事项

1. **证书管理**：妥善保管 `crypto-config` 目录
2. **权限控制**：使用正确的用户身份（Admin/User）
3. **数据验证**：调用链码前验证数据格式
4. **错误处理**：不要在错误信息中暴露敏感信息
5. **日志脱敏**：生产环境不要记录完整的用户数据

## 十一、常见问题

**Q: SDK 初始化失败？**
A: 检查配置文件路径和证书路径是否正确

**Q: 连接节点超时？**
A: 检查节点是否启动，端口是否正确

**Q: 链码函数不存在？**
A: 确认链码已安装并实例化，函数名拼写正确

**Q: 权限不足？**
A: 检查用户身份和 MSP 配置

**Q: 数据解析失败？**
A: 检查返回数据格式，使用正确的结构体解析
