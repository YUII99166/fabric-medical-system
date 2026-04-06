package v2

import (
	bc "application/blockchain"
	"application/model"
	"application/pkg/app"
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"github.com/gin-gonic/gin"
	"net/http"
)

// RequestAccess 申请授权
func RequestAccess(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(model.AccessRequestBody)

	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数解析错误: %s", err.Error()))
		return
	}

	if body.PrescriptionID == "" || body.DoctorID == "" || body.Reason == "" {
		appG.Response(http.StatusBadRequest, "失败", "参数存在空值")
		return
	}

	fmt.Printf("RequestAccess - 申请授权: PrescriptionID=%s, DoctorID=%s\n", body.PrescriptionID, body.DoctorID)

	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, []byte(body.PrescriptionID))
	bodyBytes = append(bodyBytes, []byte(body.DoctorID))
	bodyBytes = append(bodyBytes, []byte(body.Reason))

	// 添加重试机制处理各种区块链错误
	maxRetries := 5 // 增加重试次数
	for i := 0; i < maxRetries; i++ {
		fmt.Printf("RequestAccess - 尝试申请授权 (第%d次)\n", i+1)
		
		// 调用智能合约
		resp, blockchainErr := bc.ChannelExecute("requestAccess", bodyBytes)
		if blockchainErr != nil {
			fmt.Printf("RequestAccess - 第%d次尝试失败: %v\n", i+1, blockchainErr)
			
			// 检查是否是可重试的错误
			errorStr := blockchainErr.Error()
			isRetryableError := strings.Contains(errorStr, "txid") && strings.Contains(errorStr, "exists") ||
				strings.Contains(errorStr, "CONNECTION_FAILED") ||
				strings.Contains(errorStr, "TRANSIENT_FAILURE") ||
				strings.Contains(errorStr, "cannot retrieve package") ||
				strings.Contains(errorStr, "chaincode") ||
				strings.Contains(errorStr, "timeout") ||
				strings.Contains(errorStr, "unavailable")
			
			if isRetryableError && i < maxRetries-1 {
				fmt.Printf("RequestAccess - 检测到可重试错误，等待后重试\n")
				// 增加等待时间，避免频繁重试
				time.Sleep(time.Duration(i+1) * 200 * time.Millisecond)
				continue
			}
			
			// 如果不是可重试错误或者已经是最后一次尝试，直接返回错误
			appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("区块链调用失败: %v", blockchainErr))
			return
		}
		
		// 成功了，处理响应并返回
		fmt.Printf("RequestAccess - 第%d次尝试成功\n", i+1)
		
		var data map[string]interface{}
		if unmarshalErr := json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); unmarshalErr != nil {
			fmt.Printf("RequestAccess - 解析响应失败: %v\n", unmarshalErr)
			appG.Response(http.StatusInternalServerError, "失败", unmarshalErr.Error())
			return
		}

		fmt.Printf("RequestAccess - 申请授权成功\n")
		appG.Response(http.StatusOK, "成功", data)
		return
	}
	
	// 如果所有重试都失败了
	appG.Response(http.StatusInternalServerError, "失败", "申请授权失败，请稍后重试")
}

// ApproveAccess 审批授权
func ApproveAccess(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(model.ApproveAccessBody)

	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数解析错误: %s", err.Error()))
		return
	}

	if body.RequestID == "" || body.PatientID == "" || body.Approved == "" {
		appG.Response(http.StatusBadRequest, "失败", "参数存在空值")
		return
	}

	fmt.Printf("ApproveAccess - 审批授权: RequestID=%s, PatientID=%s, Approved=%s\n", body.RequestID, body.PatientID, body.Approved)

	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, []byte(body.RequestID))
	bodyBytes = append(bodyBytes, []byte(body.PatientID))
	bodyBytes = append(bodyBytes, []byte(body.Approved))
	bodyBytes = append(bodyBytes, []byte(body.RejectReason))

	// 添加重试机制处理各种区块链错误
	maxRetries := 5 // 增加重试次数
	for i := 0; i < maxRetries; i++ {
		fmt.Printf("ApproveAccess - 尝试审批 (第%d次)\n", i+1)
		
		// 调用智能合约
		resp, err := bc.ChannelExecute("approveAccess", bodyBytes)
		if err != nil {
			fmt.Printf("ApproveAccess - 第%d次尝试失败: %v\n", i+1, err)
			
			// 检查是否是可重试的错误
			errorStr := err.Error()
			isRetryableError := strings.Contains(errorStr, "txid") && strings.Contains(errorStr, "exists") ||
				strings.Contains(errorStr, "CONNECTION_FAILED") ||
				strings.Contains(errorStr, "TRANSIENT_FAILURE") ||
				strings.Contains(errorStr, "cannot retrieve package") ||
				strings.Contains(errorStr, "chaincode") ||
				strings.Contains(errorStr, "timeout") ||
				strings.Contains(errorStr, "unavailable")
			
			if isRetryableError && i < maxRetries-1 {
				fmt.Printf("ApproveAccess - 检测到可重试错误，等待后重试\n")
				// 增加等待时间，避免频繁重试
				time.Sleep(time.Duration(i+1) * 200 * time.Millisecond)
				continue
			}
			
			// 如果不是可重试错误或者已经是最后一次尝试，直接返回错误
			appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("区块链调用失败: %v", err))
			return
		}
		
		// 成功了，处理响应并返回
		fmt.Printf("ApproveAccess - 第%d次尝试成功\n", i+1)
		
		var data map[string]interface{}
		if err := json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
			fmt.Printf("ApproveAccess - 解析响应失败: %v\n", err)
			appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("解析响应失败: %v", err))
			return
		}

		fmt.Printf("ApproveAccess - 审批成功\n")
		appG.Response(http.StatusOK, "成功", data)
		return
	}
	
	// 如果所有重试都失败了
	appG.Response(http.StatusInternalServerError, "失败", "审批授权失败，请稍后重试")
}

// QueryAccessRequests 查询授权请求列表
func QueryAccessRequests(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(model.QueryAccessRequestsBody)

	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数解析错误: %s", err.Error()))
		return
	}

	if body.UserID == "" || body.Role == "" {
		appG.Response(http.StatusBadRequest, "失败", "参数存在空值")
		return
	}

	fmt.Printf("QueryAccessRequests - 查询授权请求: UserID=%s, Role=%s\n", body.UserID, body.Role)

	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, []byte(body.UserID))
	bodyBytes = append(bodyBytes, []byte(body.Role))

	// 添加重试机制处理各种区块链错误
	maxRetries := 5 // 增加重试次数
	for i := 0; i < maxRetries; i++ {
		fmt.Printf("QueryAccessRequests - 尝试查询 (第%d次)\n", i+1)
		
		// 调用智能合约
		resp, err := bc.ChannelQuery("queryAccessRequests", bodyBytes)
		if err != nil {
			fmt.Printf("QueryAccessRequests - 第%d次尝试失败: %v\n", i+1, err)
			
			// 检查是否是可重试的错误
			errorStr := err.Error()
			isRetryableError := strings.Contains(errorStr, "txid") && strings.Contains(errorStr, "exists") ||
				strings.Contains(errorStr, "CONNECTION_FAILED") ||
				strings.Contains(errorStr, "TRANSIENT_FAILURE") ||
				strings.Contains(errorStr, "cannot retrieve package") ||
				strings.Contains(errorStr, "chaincode") ||
				strings.Contains(errorStr, "timeout") ||
				strings.Contains(errorStr, "unavailable")
			
			if isRetryableError && i < maxRetries-1 {
				fmt.Printf("QueryAccessRequests - 检测到可重试错误，等待后重试\n")
				// 增加等待时间，避免频繁重试
				time.Sleep(time.Duration(i+1) * 200 * time.Millisecond)
				continue
			}
			
			// 如果不是可重试错误或者已经是最后一次尝试，直接返回错误
			appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("区块链调用失败: %v", err))
			return
		}
		
		// 成功了，处理响应并返回
		fmt.Printf("QueryAccessRequests - 第%d次尝试成功\n", i+1)
		
		var data []map[string]interface{}
		if err := json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
			fmt.Printf("QueryAccessRequests - 解析响应失败: %v\n", err)
			appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("解析响应失败: %v", err))
			return
		}

		fmt.Printf("QueryAccessRequests - 查询成功，返回 %d 条记录\n", len(data))
		appG.Response(http.StatusOK, "成功", data)
		return
	}
	
	// 如果所有重试都失败了
	appG.Response(http.StatusInternalServerError, "失败", "查询授权请求失败，请稍后重试")
}

// QueryPrescriptionsByPatient 根据患者姓名或ID查询病历
func QueryPrescriptionsByPatient(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(model.QueryPrescriptionsByPatientBody)

	if err := c.ShouldBind(body); err != nil {
		fmt.Printf("QueryPrescriptionsByPatient - 参数解析失败: %v\n", err)
		appG.Response(http.StatusBadRequest, fmt.Sprintf("参数解析错误: %s", err.Error()), nil)
		return
	}

	fmt.Printf("QueryPrescriptionsByPatient - 接收参数: SearchKey=%s, DoctorID=%s\n", body.SearchKey, body.DoctorID)

	if body.SearchKey == "" {
		fmt.Printf("QueryPrescriptionsByPatient - 搜索关键词为空\n")
		appG.Response(http.StatusBadRequest, "搜索关键词不能为空", nil)
		return
	}

	// 简化搜索逻辑：如果医生ID未同步，使用空字符串作为占位符
	doctorID := body.DoctorID
	if doctorID == "未同步到区块链" || doctorID == "" {
		doctorID = "ANY_DOCTOR" // 使用特殊标识符表示任意医生
		fmt.Printf("QueryPrescriptionsByPatient - 医生ID未同步，使用通用搜索模式\n")
	}

	// 准备多个搜索关键词进行尝试
	searchKeys := []string{}
	originalKey := body.SearchKey
	
	// 添加原始搜索词
	searchKeys = append(searchKeys, originalKey)
	
	// 如果不包含"病人-"前缀，添加带前缀的版本
	if !strings.HasPrefix(originalKey, "病人-") {
		searchKeys = append(searchKeys, "病人-"+originalKey)
	}
	
	// 如果包含"病人-"前缀，添加不带前缀的版本
	if strings.HasPrefix(originalKey, "病人-") {
		searchKeys = append(searchKeys, strings.TrimPrefix(originalKey, "病人-"))
	}

	fmt.Printf("QueryPrescriptionsByPatient - 将尝试搜索关键词: %v\n", searchKeys)

	// 依次尝试不同的搜索关键词
	var finalData []map[string]interface{}
	var lastError error
	
	for i, searchKey := range searchKeys {
		fmt.Printf("QueryPrescriptionsByPatient - 尝试搜索关键词 %d/%d: %s\n", i+1, len(searchKeys), searchKey)
		
		// 为每个搜索关键词添加重试机制
		maxRetries := 5 // 增加重试次数
		var searchSuccess bool
		
		for retry := 0; retry < maxRetries; retry++ {
			if retry > 0 {
				fmt.Printf("QueryPrescriptionsByPatient - 搜索关键词 '%s' 第%d次重试\n", searchKey, retry+1)
			}
			
			var bodyBytes [][]byte
			bodyBytes = append(bodyBytes, []byte(searchKey))
			bodyBytes = append(bodyBytes, []byte(doctorID))

			// 调用智能合约
			resp, err := bc.ChannelQuery("queryPrescriptionsByPatient", bodyBytes)
			if err != nil {
				fmt.Printf("QueryPrescriptionsByPatient - 搜索关键词 '%s' 第%d次尝试失败: %v\n", searchKey, retry+1, err)
				
				// 检查是否是可重试的错误
				errorStr := err.Error()
				isRetryableError := strings.Contains(errorStr, "txid") && strings.Contains(errorStr, "exists") ||
					strings.Contains(errorStr, "CONNECTION_FAILED") ||
					strings.Contains(errorStr, "TRANSIENT_FAILURE") ||
					strings.Contains(errorStr, "cannot retrieve package") ||
					strings.Contains(errorStr, "chaincode") ||
					strings.Contains(errorStr, "timeout") ||
					strings.Contains(errorStr, "unavailable")
				
				if isRetryableError && retry < maxRetries-1 {
					fmt.Printf("QueryPrescriptionsByPatient - 检测到可重试错误，等待后重试\n")
					// 增加等待时间，避免频繁重试
					time.Sleep(time.Duration(retry+1) * 200 * time.Millisecond)
					continue
				}
				
				lastError = err
				break // 跳出重试循环，尝试下一个搜索关键词
			}

			var data []map[string]interface{}
			if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
				fmt.Printf("QueryPrescriptionsByPatient - 搜索关键词 '%s' 解析响应失败: %v\n", searchKey, err)
				lastError = err
				break // 跳出重试循环，尝试下一个搜索关键词
			}

			fmt.Printf("QueryPrescriptionsByPatient - 搜索关键词 '%s' 返回 %d 条记录\n", searchKey, len(data))
			
			if len(data) > 0 {
				// 找到了结果，使用这个结果
				finalData = data
				searchSuccess = true
				fmt.Printf("QueryPrescriptionsByPatient - 使用搜索关键词 '%s' 的结果，共 %d 条记录\n", searchKey, len(finalData))
				break // 跳出重试循环
			}
			
			// 如果没有数据但没有错误，也跳出重试循环，尝试下一个搜索关键词
			break
		}
		
		// 如果找到了结果，跳出搜索关键词循环
		if searchSuccess {
			break
		}
	}

	// 如果所有搜索都没有结果
	if len(finalData) == 0 {
		if lastError != nil {
			fmt.Printf("QueryPrescriptionsByPatient - 所有搜索都失败，最后错误: %v\n", lastError)
			appG.Response(http.StatusInternalServerError, fmt.Sprintf("区块链查询失败: %v", lastError), nil)
		} else {
			fmt.Printf("QueryPrescriptionsByPatient - 所有搜索都没有找到结果\n")
			appG.Response(http.StatusOK, "成功", []map[string]interface{}{})
		}
		return
	}

	fmt.Printf("QueryPrescriptionsByPatient - 最终返回 %d 条记录\n", len(finalData))
	appG.Response(http.StatusOK, "成功", finalData)
}
