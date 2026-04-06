package v2

import (
	bc "application/blockchain"
	"application/db"
	"application/model"
	"application/pkg/app"
	"application/pkg/cache"
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"net/http"
)

func QueryAccountV2List(c *gin.Context) {
	appG := app.Gin{C: c}
	
	// 尝试从缓存获取
	cacheKey := "account_list"
	memCache := cache.GetGlobalCache()
	
	var data []map[string]interface{}
	if memCache.GetJSON(cacheKey, &data) {
		fmt.Println("✅ 从缓存返回账户列表:", len(data), "个账户")
		appG.Response(http.StatusOK, "成功", data)
		return
	}
	
	// 从区块链查询账户列表
	var bodyBytes [][]byte
	
	// 使用重试机制处理txid冲突
	var resp channel.Response
	var err error
	maxRetries := 3
	
	for i := 0; i < maxRetries; i++ {
		resp, err = bc.ChannelQuery("queryAccountV2List", bodyBytes)
		if err == nil {
			break
		}
		
		// 如果是txid冲突错误,等待一小段时间后重试
		if i < maxRetries-1 {
			fmt.Printf("查询区块链失败(第%d次尝试): %v, 等待后重试...\n", i+1, err)
			time.Sleep(time.Millisecond * 100)
		}
	}
	
	if err != nil {
		fmt.Printf("查询区块链账户列表失败(已重试%d次): %v\n", maxRetries, err)
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("查询区块链失败: %v", err))
		return
	}
	
	// 检查返回的数据
	if len(resp.Payload) == 0 {
		fmt.Println("区块链返回空数据,返回空数组")
		appG.Response(http.StatusOK, "成功", []map[string]interface{}{})
		return
	}
	
	if err = json.Unmarshal(resp.Payload, &data); err != nil {
		fmt.Printf("解析区块链响应失败: %v, 原始数据: %s\n", err, string(resp.Payload))
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("解析响应失败: %v", err))
		return
	}
	
	fmt.Printf("成功从区块链查询到 %d 条账户记录\n", len(data))
	
	// 添加调试信息，显示每个账户的关键字段
	for i, account := range data {
		if i < 5 { // 只显示前5个账户的详细信息
			fmt.Printf("账户 %d: username=%v, account_name=%v, account_id=%v\n", 
				i+1, account["username"], account["account_name"], account["account_id"])
		}
	}
	
	// 存入缓存，有效期5分钟
	memCache.SetJSON(cacheKey, data, 5*time.Minute)
	fmt.Println("✅ 账户列表已缓存，有效期5分钟")
	
	appG.Response(http.StatusOK, "成功", data)
}

func QueryAccountListFromDB(c *gin.Context) {
	appG := app.Gin{C: c}
	
	// 从请求中获取状态参数（可选）
	type QueryParams struct {
		Status int `json:"status" form:"status"` // -1=全部, 0=已删除, 1=正常
	}
	params := QueryParams{Status: 1} // 默认只查询正常用户
	
	// 尝试绑定参数，如果失败使用默认值
	_ = c.ShouldBind(&params)
	
	fmt.Printf("从数据库查询用户列表，状态筛选: %d\n", params.Status)
	
	// 从数据库查询用户列表
	users, err := db.ListUsersWithStatus(params.Status)
	if err != nil {
		fmt.Printf("查询数据库用户列表失败: %v\n", err)
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("查询失败: %v", err))
		return
	}
	
	// 转换为前端需要的格式
	var data []map[string]interface{}
	for _, user := range users {
		userMap := map[string]interface{}{
			"id":                user.ID,
			"username":          user.Username,
			"account_name":      user.AccountName,
			"account_id":        fmt.Sprintf("%d", user.ID), // 使用数据库ID作为account_id
			"role":              user.Role,
			"organization":      user.Organization,
			"organization_name": user.OrganizationName,
			"department":        user.Department,
			"doctor_title":      user.DoctorTitle,
			"status":            user.Status,
			"created_at":        user.CreatedAt,
			"updated_at":        user.UpdatedAt,
		}
		data = append(data, userMap)
	}
	
	fmt.Printf("成功从数据库查询到 %d 条用户记录\n", len(data))
	appG.Response(http.StatusOK, "成功", data)
}

func CreateAccountV2(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(model.CreateAccountBody)

	//解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	if body.AccountName == "" || body.Operator == "" {
		appG.Response(http.StatusBadRequest, "失败", "参数存在空值")
		return
	}
	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, []byte(body.AccountName))
	bodyBytes = append(bodyBytes, []byte(body.Operator))

	// 调用智能合约
	resp, err := bc.ChannelExecute("createAccountV2", bodyBytes)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	var data map[string]interface{}
	if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	appG.Response(http.StatusOK, "成功", data)
}

func Register(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(model.RegisterBody)

	//解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, fmt.Sprintf("参数解析错误: %s", err.Error()), nil)
		return
	}
	if body.AccountName == "" || body.Username == "" || body.Password == "" || body.Role == "" {
		appG.Response(http.StatusBadRequest, "参数存在空值，请填写完整信息", nil)
		return
	}

	// 病人和药店不需要组织信息，医生和管理员需要
	if (body.Role != "病人" && body.Role != "药店" && body.Organization == "") {
		appG.Response(http.StatusBadRequest, "医生和管理员必须选择所属组织", nil)
		return
	}

	// 先检查数据库中用户名是否已存在
	existingUser, _ := db.GetUserByUsername(body.Username)
	if existingUser != nil {
		appG.Response(http.StatusBadRequest, "用户名已存在，请更换用户名", nil)
		return
	}

	// 对密码进行哈希处理
	hashedPassword := db.HashPassword(body.Password)

	// 调用智能合约记录到区块链
	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, []byte(body.AccountName))
	bodyBytes = append(bodyBytes, []byte(body.Username))
	bodyBytes = append(bodyBytes, []byte(hashedPassword))
	bodyBytes = append(bodyBytes, []byte(body.Role))
	
	// 根据角色传递不同参数
	if body.Role == "病人" || body.Role == "药店" {
		// 病人/药店：6个参数
		bodyBytes = append(bodyBytes, []byte(fmt.Sprintf("%d", body.Age)))
		bodyBytes = append(bodyBytes, []byte(body.Gender))
		fmt.Printf("尝试注册到区块链(病人/药店): username=%s, account_name=%s, role=%s, age=%d, gender=%s\n", 
			body.Username, body.AccountName, body.Role, body.Age, body.Gender)
	} else {
		// 医生/管理员：10个参数
		bodyBytes = append(bodyBytes, []byte(body.Organization))
		bodyBytes = append(bodyBytes, []byte(body.OrganizationName))
		bodyBytes = append(bodyBytes, []byte(body.Department))
		bodyBytes = append(bodyBytes, []byte(body.DoctorTitle))
		bodyBytes = append(bodyBytes, []byte(fmt.Sprintf("%d", body.Age)))
		bodyBytes = append(bodyBytes, []byte(body.Gender))
		fmt.Printf("尝试注册到区块链(医生/管理员): username=%s, account_name=%s, role=%s, org=%s\n", 
			body.Username, body.AccountName, body.Role, body.OrganizationName)
	}
	
	// 调用智能合约记录到区块链
	// 注意：不使用重试机制，因为每次重试会生成新的txid，导致更多冲突
	resp, err := bc.ChannelExecute("register", bodyBytes)
	
	if err != nil {
		fmt.Printf("区块链注册失败: %v\n", err)
		// 检查是否是连接错误
		errMsg := err.Error()
		if strings.Contains(errMsg, "CONNECTION_FAILED") || strings.Contains(errMsg, "connection") {
			appG.Response(http.StatusInternalServerError, "区块链网络连接失败，请确认Docker已启动且区块链网络正在运行。请联系管理员启动区块链网络。", nil)
		} else if strings.Contains(errMsg, "exists") || strings.Contains(errMsg, "txid") {
			appG.Response(http.StatusInternalServerError, "区块链交易冲突，请稍后重试。如果问题持续存在，请联系管理员。", nil)
		} else {
			appG.Response(http.StatusInternalServerError, fmt.Sprintf("区块链注册失败: %v", err), nil)
		}
		return
	}

	var data map[string]interface{}
	if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
		fmt.Printf("解析区块链响应失败: %v\n", err)
		appG.Response(http.StatusInternalServerError, fmt.Sprintf("解析区块链响应失败: %v", err), nil)
		return
	}

	fmt.Printf("区块链注册成功，开始同步到数据库\n")

	// 区块链注册成功后，同步到数据库（用于登录验证）
	err = db.RegisterUser(body.Username, body.Password, body.AccountName, body.Role, 
		body.Organization, body.OrganizationName, body.Department, body.DoctorTitle)
	if err != nil {
		// 数据库同步失败只记录日志，不影响注册结果（因为区块链已经成功）
		fmt.Printf("警告: 数据库同步失败: %v\n", err)
		appG.Response(http.StatusOK, "成功", map[string]interface{}{
			"message":           "注册成功（区块链已记录，但数据库同步失败，登录可能受影响）",
			"account_name":      body.AccountName,
			"username":          body.Username,
			"role":              body.Role,
			"organization_name": body.OrganizationName,
		})
		return
	}

	fmt.Printf("数据库同步成功，注册完成\n")
	
	// 清除账户列表缓存
	memCache := cache.GetGlobalCache()
	memCache.Delete("account_list")
	fmt.Println("✅ 已清除账户列表缓存")
	
	appG.Response(http.StatusOK, "成功", data)
}

func LoginWithPassword(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(model.LoginBody)

	//解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, fmt.Sprintf("参数解析错误: %s", err.Error()), nil)
		return
	}
	if body.Username == "" || body.Password == "" {
		appG.Response(http.StatusBadRequest, "用户名或密码不能为空", nil)
		return
	}

	// 先从数据库验证用户
	user, err := db.LoginUser(body.Username, body.Password)
	if err != nil {
		// 返回更友好的错误信息
		errMsg := err.Error()
		if errMsg == "用户不存在，请先注册" {
			appG.Response(http.StatusUnauthorized, "用户不存在，请先注册", nil)
		} else if errMsg == "密码错误，请重新输入" {
			appG.Response(http.StatusUnauthorized, "密码错误，请重新输入", nil)
		} else if errMsg == "该账户已被停用，请联系管理员" {
			appG.Response(http.StatusForbidden, "该账户已被停用，请联系管理员", nil)
		} else {
			appG.Response(http.StatusInternalServerError, fmt.Sprintf("登录失败: %v", err), nil)
		}
		return
	}

	// 构建返回数据
	data := map[string]interface{}{
		"id":           user.ID,
		"username":     user.Username,
		"account_name": user.AccountName,
		"role":         user.Role,
		"created_at":   user.CreatedAt,
		"updated_at":   user.UpdatedAt,
	}

	appG.Response(http.StatusOK, "成功", data)
}

func GetUserInfo(c *gin.Context) {
	appG := app.Gin{C: c}
	
	// 从请求中获取用户 ID
	type UserIDBody struct {
		ID int `json:"id"`
	}
	body := new(UserIDBody)
	
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	
	if body.ID == 0 {
		appG.Response(http.StatusBadRequest, "失败", "用户ID不能为空")
		return
	}
	
	// 从数据库获取用户信息
	user, err := db.GetUserByID(body.ID)
	if err != nil {
		appG.Response(http.StatusNotFound, "失败", err.Error())
		return
	}
	
	// 构建返回数据
	data := map[string]interface{}{
		"id":           user.ID,
		"username":     user.Username,
		"account_name": user.AccountName,
		"role":         user.Role,
		"created_at":   user.CreatedAt,
		"updated_at":   user.UpdatedAt,
	}
	
	appG.Response(http.StatusOK, "成功", data)
}

func GetUserDetail(c *gin.Context) {
	appG := app.Gin{C: c}
	
	// 从请求中获取用户 ID 或用户名
	type UserQueryBody struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
	}
	body := new(UserQueryBody)
	
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错: %s", err.Error()))
		return
	}
	
	var user *db.User
	var err error
	
	// 优先使用ID查询，如果没有ID则使用用户名查询
	if body.ID != 0 {
		user, err = db.GetUserDetail(body.ID)
	} else if body.Username != "" {
		user, err = db.GetUserByUsername(body.Username)
	} else {
		appG.Response(http.StatusBadRequest, "失败", "用户ID或用户名不能为空")
		return
	}
	
	if err != nil {
		appG.Response(http.StatusNotFound, "失败", err.Error())
		return
	}
	
	// 构建返回数据
	data := map[string]interface{}{
		"id":                user.ID,
		"username":          user.Username,
		"account_name":      user.AccountName,
		"role":              user.Role,
		"organization":      user.Organization,
		"organization_name": user.OrganizationName,
		"department":        user.Department,
		"doctor_title":      user.DoctorTitle,
		"status":            user.Status,
		"created_at":        user.CreatedAt,
		"updated_at":        user.UpdatedAt,
	}
	
	appG.Response(http.StatusOK, "成功", data)
}

func UpdateUser(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(model.UpdateUserBody)

	// 解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数解析错误: %s", err.Error()))
		return
	}
	
	if body.ID == 0 {
		appG.Response(http.StatusBadRequest, "失败", "用户ID不能为空")
		return
	}
	
	if body.AccountName == "" || body.Role == "" {
		appG.Response(http.StatusBadRequest, "失败", "姓名和角色不能为空")
		return
	}
	
	// 检查用户是否存在
	_, err := db.GetUserByID(body.ID)
	if err != nil {
		appG.Response(http.StatusNotFound, "失败", "用户不存在")
		return
	}
	
	// 更新用户信息
	err = db.UpdateUserFull(body.ID, body.AccountName, body.Role, body.Organization, body.OrganizationName, body.Department, body.DoctorTitle)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("更新用户失败: %v", err))
		return
	}
	
	// 如果提供了新密码，则更新密码
	if body.Password != "" {
		err = db.UpdateUserPassword(body.ID, body.Password)
		if err != nil {
			appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("更新密码失败: %v", err))
			return
		}
	}
	
	appG.Response(http.StatusOK, "成功", map[string]interface{}{
		"message": "用户信息更新成功",
	})
}

func DeleteUser(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(model.DeleteUserBody)

	// 解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数解析错误: %s", err.Error()))
		return
	}
	
	if body.ID == 0 {
		appG.Response(http.StatusBadRequest, "失败", "用户ID不能为空")
		return
	}
	
	// 检查用户是否存在
	user, err := db.GetUserByID(body.ID)
	if err != nil {
		appG.Response(http.StatusNotFound, "失败", "用户不存在")
		return
	}
	
	// 检查是否为管理员，如果是，确保不是最后一个管理员
	if user.Role == "管理员" {
		adminCount, err := db.CountAdminsExcluding(body.ID)
		if err != nil {
			appG.Response(http.StatusInternalServerError, "失败", "检查管理员数量失败")
			return
		}
		if adminCount == 0 {
			appG.Response(http.StatusBadRequest, "失败", "不能删除最后一个管理员账户")
			return
		}
	}
	
	// 软删除用户（将status设为0）
	err = db.DeleteUser(body.ID)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("删除用户失败: %v", err))
		return
	}
	
	appG.Response(http.StatusOK, "成功", map[string]interface{}{
		"message": "用户已停用（软删除），可以通过恢复功能重新启用",
	})
}

func RestoreUser(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(model.DeleteUserBody) // 复用DeleteUserBody，只需要ID

	// 解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数解析错误: %s", err.Error()))
		return
	}
	
	if body.ID == 0 {
		appG.Response(http.StatusBadRequest, "失败", "用户ID不能为空")
		return
	}
	
	// 检查用户是否存在
	_, err := db.GetUserByID(body.ID)
	if err != nil {
		appG.Response(http.StatusNotFound, "失败", "用户不存在")
		return
	}
	
	// 恢复用户（将status设为1）
	err = db.RestoreUser(body.ID)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("恢复用户失败: %v", err))
		return
	}
	
	appG.Response(http.StatusOK, "成功", map[string]interface{}{
		"message": "用户已恢复，可以正常登录使用",
	})
}

func BatchDeleteUsers(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(model.BatchDeleteUsersBody)

	// 解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数解析错误: %s", err.Error()))
		return
	}
	
	if len(body.IDs) == 0 {
		appG.Response(http.StatusBadRequest, "失败", "用户ID列表不能为空")
		return
	}
	
	successCount := 0
	failedCount := 0
	var errors []string
	
	for _, id := range body.IDs {
		// 检查用户是否存在
		user, err := db.GetUserByID(id)
		if err != nil {
			failedCount++
			errors = append(errors, fmt.Sprintf("用户ID %d: %s", id, err.Error()))
			continue
		}
		
		// 检查是否为管理员
		if user.Role == "管理员" {
			adminCount, err := db.CountAdminsExcluding(id)
			if err != nil {
				failedCount++
				errors = append(errors, fmt.Sprintf("用户ID %d: 检查管理员数量失败", id))
				continue
			}
			if adminCount == 0 {
				failedCount++
				errors = append(errors, fmt.Sprintf("用户ID %d: 不能删除最后一个管理员", id))
				continue
			}
		}
		
		// 删除用户
		err = db.DeleteUser(id)
		if err != nil {
			failedCount++
			errors = append(errors, fmt.Sprintf("用户ID %d: %s", id, err.Error()))
			continue
		}
		
		successCount++
	}
	
	result := map[string]interface{}{
		"success_count": successCount,
		"failed_count":  failedCount,
		"message":       fmt.Sprintf("批量删除完成，成功 %d 个，失败 %d 个", successCount, failedCount),
	}
	
	if len(errors) > 0 {
		result["errors"] = errors
	}
	
	appG.Response(http.StatusOK, "成功", result)
}

// SyncAccountFromBlockchain 从区块链重新同步用户账户ID
func SyncAccountFromBlockchain(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(model.SyncAccountBody)

	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数解析错误: %s", err.Error()))
		return
	}

	if body.Username == "" {
		appG.Response(http.StatusBadRequest, "失败", "用户名不能为空")
		return
	}

	fmt.Printf("开始为用户 %s 同步区块链账户ID\n", body.Username)

	// 查询区块链账户列表，使用重试机制
	var bodyBytes [][]byte
	var resp channel.Response
	var err error
	maxRetries := 5
	
	for i := 0; i < maxRetries; i++ {
		resp, err = bc.ChannelQuery("queryAccountV2List", bodyBytes)
		if err == nil {
			break
		}
		
		// 检查是否是可重试的错误
		errMsg := err.Error()
		isRetryable := strings.Contains(errMsg, "txid") || 
			strings.Contains(errMsg, "CONNECTION_FAILED") || 
			strings.Contains(errMsg, "TRANSIENT_FAILURE") ||
			strings.Contains(errMsg, "timeout") ||
			strings.Contains(errMsg, "unavailable")
		
		if isRetryable && i < maxRetries-1 {
			fmt.Printf("查询区块链失败(第%d次尝试): %v, 等待后重试...\n", i+1, err)
			time.Sleep(time.Millisecond * 200 * time.Duration(i+1)) // 递增延迟
		} else {
			break
		}
	}
	
	if err != nil {
		fmt.Printf("查询区块链账户列表失败(已重试%d次): %v\n", maxRetries, err)
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("查询区块链失败: %v", err))
		return
	}

	var accountList []map[string]interface{}
	if err = json.Unmarshal(resp.Payload, &accountList); err != nil {
		fmt.Printf("解析区块链响应失败: %v\n", err)
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("解析响应失败: %v", err))
		return
	}

	// 查找匹配的用户
	var foundAccount map[string]interface{}
	for _, account := range accountList {
		username, _ := account["username"].(string)
		if username == body.Username {
			foundAccount = account
			break
		}
	}

	if foundAccount == nil {
		fmt.Printf("在区块链中未找到用户: %s\n", body.Username)
		appG.Response(http.StatusNotFound, "失败", "在区块链中未找到该用户，请确认用户已注册到区块链")
		return
	}

	accountID, _ := foundAccount["account_id"].(string)
	if accountID == "" {
		appG.Response(http.StatusInternalServerError, "失败", "区块链中的用户账户ID为空")
		return
	}

	fmt.Printf("找到用户 %s 的区块链账户ID: %s\n", body.Username, accountID)

	// 返回找到的账户信息
	result := map[string]interface{}{
		"username":          foundAccount["username"],
		"account_name":      foundAccount["account_name"],
		"account_id":        foundAccount["account_id"],
		"organization":      foundAccount["organization"],
		"organization_name": foundAccount["organization_name"],
		"role":              foundAccount["role"],
	}

	appG.Response(http.StatusOK, "成功", result)
}
