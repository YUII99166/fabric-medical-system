# Fabric SDK 实战示例

## 示例1：完整的用户注册流程

### 前端代码（Vue.js）
```javascript
// application/web/src/api/accountV2.js
export function register(data) {
  return request({
    url: '/register',
    method: 'post',
    data
  })
}

// 组件中调用
async handleRegister() {
  const userData = {
    username: 'patient001',
    password: '123456',
    account_name: '病人-张三',
    role: 'patient',
    age: 35,
    gender: '男'
  }
  
  const response = await register(userData)
  if (response.code === 200) {
    this.$message.success('注册成功')
  }
}
```

### 后端代码（Go）
```go
// application/server/api/v2/account.go
func Register(c *gin.Context) {
	appG := app.Gin{C: c}
	
	// 1. 接收前端数据
	var req struct {
		Username    string `json:"username" binding:"required"`
		Password    string `json:"password" binding:"required"`
		AccountName string `json:"account_name" binding:"required"`
		Role        string `json:"role" binding:"required"`
		Age         int    `json:"age"`
		Gender      string `json:"gender"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		appG.Response(http.StatusBadRequest, "参数错误", err.Error())
		return
	}
	
	// 2. 生成区块链账户ID
	accountID := generateAccountID()
	
	// 3. 准备区块链数据
	blockchainData := map[string]interface{}{
		"account_id":   accountID,
		"account_name": req.AccountName,
		"username":     req.Username,
		"password":     hashPassword(req.Password),
		"role":         req.Role,
		"age":          req.Age,
		"gender":       req.Gender,
	}
	
	// 4. 序列化为JSON
	jsonData, _ := json.Marshal(blockchainData)
	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, jsonData)
	
	// 5. 调用区块链
	resp, err := bc.ChannelExecute("register", bodyBytes)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "区块链注册失败", err.Error())
		return
	}
	
	// 6. 保存到MySQL（用户认证信息）
	dbUser := model.User{
		Username:  req.Username,
		Password:  hashPassword(req.Password),
		AccountID: accountID,
		Role:      req.Role,
	}
	db.CreateUser(&dbUser)
	
	// 7. 返回结果
	appG.Response(http.StatusOK, "注册成功", map[string]interface{}{
		"account_id":     accountID,
		"transaction_id": resp.TransactionID,
	})
}
```

## 示例2：创建病历并查询

### 创建病历
```go
func CreatePrescriptionExample() {
	// 准备病历数据
	prescription := map[string]interface{}{
		"id":                generateID(),
		"prescription_no":   "P2026022201",
		"patient":           "patient_id_123",
		"patient_name":      "张三",
		"chief_complaint":   "头痛、发热",
		"present_illness":   "患者3天前出现头痛、发热症状",
		"physical_exam":     "体温38.5℃，血压正常",
		"diagnosis":         "上呼吸道感染",
		"drug": []map[string]string{
			{"Name": "阿莫西林", "amount": "2盒"},
			{"Name": "布洛芬", "amount": "1盒"},
		},
		"medical_advice":    "多休息，多喝水",
		"doctor":            "doctor_id_456",
		"doctor_name":       "李医生",
		"doctor_title":      "主治医师",
		"hospital_name":     "协和医院",
		"department":        "内科",
		"created":           time.Now().Format("2006-01-02 15:04:05"),
	}
	
	// 序列化
	jsonData, _ := json.Marshal(prescription)
	var args [][]byte
	args = append(args, jsonData)
	
	// 调用区块链
	resp, err := bc.ChannelExecute("createPrescription", args)
	if err != nil {
		fmt.Printf("创建失败: %v\n", err)
		return
	}
	
	fmt.Printf("病历创建成功！\n")
	fmt.Printf("交易ID: %s\n", resp.TransactionID)
}
```

### 查询病历
```go
func QueryPrescriptionExample(patientID string) {
	// 准备参数
	var args [][]byte
	args = append(args, []byte(patientID))
	
	// 查询区块链
	resp, err := bc.ChannelQuery("queryPrescription", args)
	if err != nil {
		fmt.Printf("查询失败: %v\n", err)
		return
	}
	
	// 解析结果
	var prescriptions []map[string]interface{}
	json.Unmarshal(resp.Payload, &prescriptions)
	
	// 打印结果
	fmt.Printf("查询到 %d 条病历\n", len(prescriptions))
	for i, p := range prescriptions {
		fmt.Printf("\n病历 %d:\n", i+1)
		fmt.Printf("  病历号: %s\n", p["prescription_no"])
		fmt.Printf("  诊断: %s\n", p["diagnosis"])
		fmt.Printf("  医生: %s\n", p["doctor_name"])
		fmt.Printf("  时间: %s\n", p["created"])
	}
}
```

## 示例3：授权访问流程

### 1. 病人授权医生访问病历
```go
func RequestAccessExample() {
	// 准备授权请求数据
	request := map[string]interface{}{
		"id":                generateID(),
		"prescription_id":   "prescription_123",
		"patient_id":        "patient_456",
		"patient_name":      "张三",
		"doctor_id":         "doctor_789",
		"doctor_name":       "王医生",
		"doctor_org":        "JDMSP",
		"doctor_org_name":   "301医院",
		"reason":            "跨院会诊",
		"status":            "pending",
		"request_time":      time.Now().Format("2006-01-02 15:04:05"),
	}
	
	// 序列化
	jsonData, _ := json.Marshal(request)
	var args [][]byte
	args = append(args, jsonData)
	
	// 调用区块链
	resp, err := bc.ChannelExecute("requestAccess", args)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return
	}
	
	fmt.Printf("授权请求已提交，交易ID: %s\n", resp.TransactionID)
}
```

### 2. 批准授权
```go
func ApproveAccessExample(requestID string) {
	// 准备批准数据
	approval := map[string]interface{}{
		"request_id": requestID,
		"status":     "approved",
	}
	
	jsonData, _ := json.Marshal(approval)
	var args [][]byte
	args = append(args, jsonData)
	
	// 调用区块链
	resp, err := bc.ChannelExecute("approveAccess", args)
	if err != nil {
		fmt.Printf("批准失败: %v\n", err)
		return
	}
	
	fmt.Printf("授权已批准\n")
}
```

## 示例4：统计查询

### 获取系统统计数据
```go
func GetStatisticsExample() {
	// 1. 查询用户总数
	userResp, _ := bc.ChannelQuery("queryAccountV2List", [][]byte{})
	var users []map[string]interface{}
	json.Unmarshal(userResp.Payload, &users)
	
	// 2. 查询病历总数
	prescResp, _ := bc.ChannelQuery("queryPrescription", [][]byte{})
	var prescriptions []map[string]interface{}
	json.Unmarshal(prescResp.Payload, &prescriptions)
	
	// 3. 查询药品订单总数
	drugResp, _ := bc.ChannelQuery("queryDrugOrder", [][]byte{})
	var orders []map[string]interface{}
	json.Unmarshal(drugResp.Payload, &orders)
	
	// 4. 统计今日新增
	today := time.Now().Format("2006-01-02")
	todayPrescriptions := 0
	for _, p := range prescriptions {
		if created, ok := p["created"].(string); ok {
			if strings.HasPrefix(created, today) {
				todayPrescriptions++
			}
		}
	}
	
	// 5. 返回统计结果
	stats := map[string]interface{}{
		"total_users":          len(users),
		"total_prescriptions":  len(prescriptions),
		"total_orders":         len(orders),
		"today_prescriptions":  todayPrescriptions,
	}
	
	fmt.Printf("系统统计: %+v\n", stats)
}
```

## 示例5：批量操作

### 批量创建用户
```go
func BatchCreateUsersExample() {
	users := []map[string]interface{}{
		{
			"account_id":   "user001",
			"account_name": "病人-张三",
			"role":         "patient",
		},
		{
			"account_id":   "user002",
			"account_name": "病人-李四",
			"role":         "patient",
		},
		{
			"account_id":   "user003",
			"account_name": "医生-王五",
			"role":         "doctor",
		},
	}
	
	successCount := 0
	for _, user := range users {
		jsonData, _ := json.Marshal(user)
		var args [][]byte
		args = append(args, jsonData)
		
		_, err := bc.ChannelExecute("createAccountV2", args)
		if err != nil {
			fmt.Printf("创建用户 %s 失败: %v\n", user["account_name"], err)
		} else {
			successCount++
		}
		
		// 避免过快调用
		time.Sleep(100 * time.Millisecond)
	}
	
	fmt.Printf("批量创建完成，成功 %d/%d\n", successCount, len(users))
}
```

## 示例6：错误处理和重试

### 带重试机制的查询
```go
func QueryWithRetry(fcn string, args [][]byte, maxRetries int) (channel.Response, error) {
	var resp channel.Response
	var err error
	
	for i := 0; i < maxRetries; i++ {
		resp, err = bc.ChannelQuery(fcn, args)
		if err == nil {
			return resp, nil
		}
		
		fmt.Printf("第 %d 次尝试失败: %v\n", i+1, err)
		
		// 等待后重试
		time.Sleep(time.Duration(i+1) * time.Second)
	}
	
	return resp, fmt.Errorf("重试 %d 次后仍然失败: %v", maxRetries, err)
}

// 使用示例
func Example() {
	resp, err := QueryWithRetry("queryAccountV2List", [][]byte{}, 3)
	if err != nil {
		fmt.Printf("查询失败: %v\n", err)
		return
	}
	
	fmt.Printf("查询成功\n")
}
```

## 示例7：事务处理

### 创建病历并创建药品订单（原子操作）
```go
func CreatePrescriptionWithOrder() error {
	// 1. 创建病历
	prescription := map[string]interface{}{
		"id":            "presc_001",
		"patient":       "patient_123",
		"diagnosis":     "感冒",
		// ... 其他字段
	}
	
	prescData, _ := json.Marshal(prescription)
	var prescArgs [][]byte
	prescArgs = append(prescArgs, prescData)
	
	prescResp, err := bc.ChannelExecute("createPrescription", prescArgs)
	if err != nil {
		return fmt.Errorf("创建病历失败: %v", err)
	}
	
	// 2. 创建药品订单
	order := map[string]interface{}{
		"id":           "order_001",
		"prescription": "presc_001",
		"patient":      "patient_123",
		"Name":         "阿莫西林",
		"amount":       "2盒",
	}
	
	orderData, _ := json.Marshal(order)
	var orderArgs [][]byte
	orderArgs = append(orderArgs, orderData)
	
	orderResp, err := bc.ChannelExecute("createDrugOrder", orderArgs)
	if err != nil {
		// 注意：病历已经创建，这里需要考虑补偿机制
		return fmt.Errorf("创建订单失败: %v", err)
	}
	
	fmt.Printf("病历和订单创建成功\n")
	fmt.Printf("病历交易ID: %s\n", prescResp.TransactionID)
	fmt.Printf("订单交易ID: %s\n", orderResp.TransactionID)
	
	return nil
}
```

## 示例8：测试代码

### 单元测试
```go
package blockchain_test

import (
	bc "application/blockchain"
	"testing"
)

func TestQueryAccountList(t *testing.T) {
	// 初始化 SDK
	bc.Init()
	
	// 查询用户列表
	resp, err := bc.ChannelQuery("queryAccountV2List", [][]byte{})
	if err != nil {
		t.Fatalf("查询失败: %v", err)
	}
	
	// 验证返回数据
	if len(resp.Payload) == 0 {
		t.Error("返回数据为空")
	}
	
	t.Logf("查询成功，数据长度: %d", len(resp.Payload))
}
```

## 总结

使用 Fabric SDK 的关键步骤：

1. **初始化 SDK**：加载配置文件
2. **准备参数**：序列化为 JSON 字节数组
3. **调用函数**：ChannelQuery（查询）或 ChannelExecute（写入）
4. **处理结果**：解析返回的 Payload
5. **错误处理**：捕获并处理各种错误

记住：
- 查询操作用 `ChannelQuery`
- 写入操作用 `ChannelExecute`
- 参数必须是 `[][]byte` 类型
- 返回数据在 `resp.Payload` 中
