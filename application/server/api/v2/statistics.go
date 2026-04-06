package v2

import (
	bc "application/blockchain"
	"application/db"
	"application/pkg/app"
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"net/http"
)

// GetStatistics 获取系统统计数据
func GetStatistics(c *gin.Context) {
	appG := app.Gin{C: c}

	// 1. 获取用户统计 - 添加重试机制
	var userRes channel.Response
	var err error
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		userRes, err = bc.ChannelQuery("queryAccountV2List", [][]byte{})
		if err != nil {
			errorStr := err.Error()
			isRetryableError := strings.Contains(errorStr, "txid") && strings.Contains(errorStr, "exists") ||
				strings.Contains(errorStr, "CONNECTION_FAILED") ||
				strings.Contains(errorStr, "TRANSIENT_FAILURE") ||
				strings.Contains(errorStr, "cannot retrieve package") ||
				strings.Contains(errorStr, "chaincode") ||
				strings.Contains(errorStr, "timeout") ||
				strings.Contains(errorStr, "unavailable")
			
			if isRetryableError && i < maxRetries-1 {
				time.Sleep(time.Duration(i+1) * 200 * time.Millisecond)
				continue
			}
			appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("查询用户失败: %v", err))
			return
		}
		break
	}

	var userList []map[string]interface{}
	if err = json.Unmarshal(bytes.NewBuffer(userRes.Payload).Bytes(), &userList); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("解析用户数据失败: %v", err))
		return
	}

	// 2. 获取病历统计 - 添加重试机制
	var prescriptionRes channel.Response
	for i := 0; i < maxRetries; i++ {
		prescriptionRes, err = bc.ChannelQuery("queryPrescription", [][]byte{})
		if err != nil {
			errorStr := err.Error()
			isRetryableError := strings.Contains(errorStr, "txid") && strings.Contains(errorStr, "exists") ||
				strings.Contains(errorStr, "CONNECTION_FAILED") ||
				strings.Contains(errorStr, "TRANSIENT_FAILURE") ||
				strings.Contains(errorStr, "cannot retrieve package") ||
				strings.Contains(errorStr, "chaincode") ||
				strings.Contains(errorStr, "timeout") ||
				strings.Contains(errorStr, "unavailable")
			
			if isRetryableError && i < maxRetries-1 {
				time.Sleep(time.Duration(i+1) * 200 * time.Millisecond)
				continue
			}
			appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("查询病历失败: %v", err))
			return
		}
		break
	}

	var prescriptionList []map[string]interface{}
	if err = json.Unmarshal(bytes.NewBuffer(prescriptionRes.Payload).Bytes(), &prescriptionList); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("解析病历数据失败: %v", err))
		return
	}

	// 3. 获取药品订单统计 - 添加重试机制
	var drugOrderRes channel.Response
	var drugOrderList []map[string]interface{}
	for i := 0; i < maxRetries; i++ {
		drugOrderRes, err = bc.ChannelQuery("queryDrugOrder", [][]byte{})
		if err != nil {
			errorStr := err.Error()
			isRetryableError := strings.Contains(errorStr, "txid") && strings.Contains(errorStr, "exists") ||
				strings.Contains(errorStr, "CONNECTION_FAILED") ||
				strings.Contains(errorStr, "TRANSIENT_FAILURE") ||
				strings.Contains(errorStr, "cannot retrieve package") ||
				strings.Contains(errorStr, "chaincode") ||
				strings.Contains(errorStr, "timeout") ||
				strings.Contains(errorStr, "unavailable")
			
			if isRetryableError && i < maxRetries-1 {
				time.Sleep(time.Duration(i+1) * 200 * time.Millisecond)
				continue
			}
			// 如果查询失败，设为空数组
			drugOrderList = []map[string]interface{}{}
			break
		}
		// 成功查询，解析数据
		if err = json.Unmarshal(bytes.NewBuffer(drugOrderRes.Payload).Bytes(), &drugOrderList); err != nil {
			// 如果解析失败，设为空数组
			drugOrderList = []map[string]interface{}{}
		}
		break
	}

	// 4. 计算今日新增数据
	today := time.Now().Format("2006-01-02")
	todayNewUsers := 0
	todayNewPrescriptions := 0
	todayNewOrders := 0

	// 统计今日新增用户（从数据库）
	allUsers, err := db.GetAllUsers()
	if err == nil {
		for _, user := range allUsers {
			if user.CreatedAt != "" && user.CreatedAt[:10] == today {
				todayNewUsers++
			}
		}
	}

	// 统计今日新增病历
	for _, prescription := range prescriptionList {
		if created, ok := prescription["created"].(string); ok {
			if len(created) >= 10 && created[:10] == today {
				todayNewPrescriptions++
			}
		}
	}

	// 统计今日新增药品订单
	for _, order := range drugOrderList {
		if created, ok := order["created"].(string); ok {
			if len(created) >= 10 && created[:10] == today {
				todayNewOrders++
			}
		}
	}

	// 5. 构建响应数据
	stats := map[string]interface{}{
		"userCount":              len(userList),
		"prescriptionCount":      len(prescriptionList),
		"drugOrderCount":         len(drugOrderList),
		"organizationCount":      4, // 固定4个组织
		"todayNewUsers":          todayNewUsers,
		"todayNewPrescriptions":  todayNewPrescriptions,
		"todayNewOrders":         todayNewOrders,
	}

	appG.Response(http.StatusOK, "成功", stats)
}

// GetRecentActivities 获取最近活动
func GetRecentActivities(c *gin.Context) {
	appG := app.Gin{C: c}

	activities := []map[string]interface{}{}

	// 1. 获取最近注册的用户
	allUsers, err := db.GetAllUsers()
	if err == nil && len(allUsers) > 0 {
		// 数据库已经按created_at DESC排序，直接取前5个
		count := 0
		for i := 0; i < len(allUsers) && count < 5; i++ {
			user := allUsers[i]
			activities = append(activities, map[string]interface{}{
				"timestamp": user.CreatedAt,
				"type":      "用户注册",
				"content":   fmt.Sprintf("新用户 \"%s\" 注册成功", user.Username),
				"icon":      "el-icon-user",
				"color":     "#67C23A",
			})
			count++
		}
		fmt.Printf("✅ 获取到 %d 个用户注册活动\n", count)
	} else {
		fmt.Printf("⚠️  获取用户失败或无用户: %v\n", err)
	}

	// 2. 获取最近创建的病历 - 添加重试机制
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		prescriptionRes, err := bc.ChannelQuery("queryPrescription", [][]byte{})
		if err != nil {
			errorStr := err.Error()
			isRetryableError := strings.Contains(errorStr, "txid") && strings.Contains(errorStr, "exists") ||
				strings.Contains(errorStr, "CONNECTION_FAILED") ||
				strings.Contains(errorStr, "TRANSIENT_FAILURE") ||
				strings.Contains(errorStr, "cannot retrieve package") ||
				strings.Contains(errorStr, "chaincode") ||
				strings.Contains(errorStr, "timeout") ||
				strings.Contains(errorStr, "unavailable")
			
			if isRetryableError && i < maxRetries-1 {
				time.Sleep(time.Duration(i+1) * 200 * time.Millisecond)
				continue
			}
			break // 如果失败，跳过病历活动
		}
		
		var prescriptionList []map[string]interface{}
		if err = json.Unmarshal(bytes.NewBuffer(prescriptionRes.Payload).Bytes(), &prescriptionList); err == nil {
			// 取最近5条病历
			count := 0
			for i := len(prescriptionList) - 1; i >= 0 && count < 5; i-- {
				p := prescriptionList[i]
				doctorName := "未知医生"
				patientName := "未知患者"
				if name, ok := p["doctor_name"].(string); ok {
					doctorName = name
				}
				if name, ok := p["patient_name"].(string); ok {
					patientName = name
				}
				created := ""
				if ct, ok := p["created"].(string); ok {
					created = ct
				}
				activities = append(activities, map[string]interface{}{
					"timestamp": created,
					"type":      "病历创建",
					"content":   fmt.Sprintf("医生 \"%s\" 为患者 \"%s\" 创建了新病历", doctorName, patientName),
					"icon":      "el-icon-document",
					"color":     "#409EFF",
				})
				count++
			}
		}
		break
	}

	// 3. 暂时跳过药品订单（chaincode 中没有实现该功能）

	// 按时间倒序排序（最新的在前面）
	// 这里简化处理，实际应该按 timestamp 排序

	// 限制返回最多15条活动
	if len(activities) > 15 {
		activities = activities[:15]
	}

	appG.Response(http.StatusOK, "成功", activities)
}
