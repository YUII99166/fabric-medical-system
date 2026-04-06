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

// AddSupplementRecord 添加补充诊疗记录
func AddSupplementRecord(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(model.SupplementRecordRequestBody)

	// 解析Body参数
	if err := c.ShouldBind(body); err != nil {
		fmt.Printf("添加补充记录 - 参数解析失败: %v\n", err)
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数解析错误: %s", err.Error()))
		return
	}

	fmt.Printf("添加补充记录 - 接收参数: OriginalPrescriptionID=%s, DoctorID=%s, RecordType=%s, Diagnosis=%s\n",
		body.OriginalPrescriptionID, body.DoctorID, body.RecordType, body.Diagnosis)

	// 验证必填字段
	if body.OriginalPrescriptionID == "" || body.DoctorID == "" || body.Diagnosis == "" {
		fmt.Printf("添加补充记录 - 参数存在空值\n")
		appG.Response(http.StatusBadRequest, "失败", "参数存在空值（原始病历ID、医生ID、诊断不能为空）")
		return
	}

	// 构建参数数组
	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, []byte(body.OriginalPrescriptionID))
	bodyBytes = append(bodyBytes, []byte(body.DoctorID))
	bodyBytes = append(bodyBytes, []byte(body.RecordType))
	bodyBytes = append(bodyBytes, []byte(body.ChiefComplaint))
	bodyBytes = append(bodyBytes, []byte(body.PresentIllness))
	bodyBytes = append(bodyBytes, []byte(body.PhysicalExam))
	bodyBytes = append(bodyBytes, []byte(body.Diagnosis))
	bodyBytes = append(bodyBytes, []byte(body.Treatment))
	bodyBytes = append(bodyBytes, []byte(body.DrugName))
	bodyBytes = append(bodyBytes, []byte(body.DrugAmount))
	bodyBytes = append(bodyBytes, []byte(body.MedicalAdvice))
	bodyBytes = append(bodyBytes, []byte(body.Comment))

	fmt.Printf("添加补充记录 - 调用区块链智能合约\n")
	
	// 添加重试机制处理各种区块链错误
	maxRetries := 5 // 增加重试次数
	for i := 0; i < maxRetries; i++ {
		fmt.Printf("AddSupplementRecord - 尝试添加 (第%d次)\n", i+1)
		
		// 调用智能合约
		resp, err := bc.ChannelExecute("addSupplementRecord", bodyBytes)
		if err != nil {
			fmt.Printf("AddSupplementRecord - 第%d次尝试失败: %v\n", i+1, err)
			
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
				fmt.Printf("AddSupplementRecord - 检测到可重试错误，等待后重试\n")
				// 增加等待时间，避免频繁重试
				time.Sleep(time.Duration(i+1) * 200 * time.Millisecond)
				continue
			}
			
			// 如果不是可重试错误或者已经是最后一次尝试，直接返回错误
			appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("区块链调用失败: %v", err))
			return
		}
		
		// 成功了，处理响应并返回
		fmt.Printf("AddSupplementRecord - 第%d次尝试成功\n", i+1)
		
		var data map[string]interface{}
		if err := json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
			fmt.Printf("AddSupplementRecord - 解析响应失败: %v\n", err)
			appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("解析响应失败: %v", err))
			return
		}

		fmt.Printf("AddSupplementRecord - 添加补充记录成功\n")
		appG.Response(http.StatusOK, "成功", data)
		return
	}
	
	// 如果所有重试都失败了
	appG.Response(http.StatusInternalServerError, "失败", "添加补充记录失败，请稍后重试")
}

// QuerySupplementRecords 查询补充记录
func QuerySupplementRecords(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(model.QuerySupplementRecordsBody)

	// 解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错: %s", err.Error()))
		return
	}

	var bodyBytes [][]byte
	if body.OriginalPrescriptionID != "" {
		bodyBytes = append(bodyBytes, []byte(body.OriginalPrescriptionID))
	}

	// 添加重试机制处理各种区块链错误
	maxRetries := 5 // 增加重试次数
	for i := 0; i < maxRetries; i++ {
		fmt.Printf("QuerySupplementRecords - 尝试查询 (第%d次)\n", i+1)
		
		// 调用智能合约
		resp, err := bc.ChannelQuery("querySupplementRecords", bodyBytes)
		if err != nil {
			fmt.Printf("QuerySupplementRecords - 第%d次尝试失败: %v\n", i+1, err)
			
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
				fmt.Printf("QuerySupplementRecords - 检测到可重试错误，等待后重试\n")
				// 增加等待时间，避免频繁重试
				time.Sleep(time.Duration(i+1) * 200 * time.Millisecond)
				continue
			}
			
			// 如果不是可重试错误或者已经是最后一次尝试，直接返回错误
			appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("区块链调用失败: %v", err))
			return
		}
		
		// 成功了，处理响应并返回
		fmt.Printf("QuerySupplementRecords - 第%d次尝试成功\n", i+1)
		
		// 反序列化json
		var data []map[string]interface{}
		if err := json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
			fmt.Printf("QuerySupplementRecords - 解析响应失败: %v\n", err)
			appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("解析响应失败: %v", err))
			return
		}

		fmt.Printf("QuerySupplementRecords - 查询成功，返回 %d 条记录\n", len(data))
		appG.Response(http.StatusOK, "成功", data)
		return
	}
	
	// 如果所有重试都失败了
	appG.Response(http.StatusInternalServerError, "失败", "查询补充记录失败，请稍后重试")
}

// QueryFullMedicalHistory 查询完整病历历史
func QueryFullMedicalHistory(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(model.QueryFullMedicalHistoryBody)

	// 解析Body参数
	if err := c.ShouldBind(body); err != nil {
		fmt.Printf("QueryFullMedicalHistory - 参数解析失败: %v\n", err)
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错: %s", err.Error()))
		return
	}

	fmt.Printf("QueryFullMedicalHistory - 接收参数: PrescriptionID=%s\n", body.PrescriptionID)

	if body.PrescriptionID == "" {
		fmt.Printf("QueryFullMedicalHistory - 病历ID为空\n")
		appG.Response(http.StatusBadRequest, "失败", "病历ID不能为空")
		return
	}

	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, []byte(body.PrescriptionID))

	// 添加重试机制处理各种区块链错误
	maxRetries := 5 // 增加重试次数
	for i := 0; i < maxRetries; i++ {
		fmt.Printf("QueryFullMedicalHistory - 尝试查询 (第%d次)\n", i+1)
		
		// 调用智能合约
		resp, err := bc.ChannelQuery("queryFullMedicalHistory", bodyBytes)
		if err != nil {
			fmt.Printf("QueryFullMedicalHistory - 第%d次尝试失败: %v\n", i+1, err)
			
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
				fmt.Printf("QueryFullMedicalHistory - 检测到可重试错误，等待后重试\n")
				// 增加等待时间，避免频繁重试
				time.Sleep(time.Duration(i+1) * 200 * time.Millisecond)
				continue
			}
			
			// 如果不是可重试错误或者已经是最后一次尝试，直接返回错误
			appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("区块链调用失败: %v", err))
			return
		}
		
		// 成功了，处理响应并返回
		fmt.Printf("QueryFullMedicalHistory - 第%d次尝试成功\n", i+1)
		fmt.Printf("QueryFullMedicalHistory - 区块链返回数据长度: %d\n", len(resp.Payload))
		fmt.Printf("QueryFullMedicalHistory - 区块链返回原始数据: %s\n", string(resp.Payload))

		// 反序列化json
		var data map[string]interface{}
		if err := json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
			fmt.Printf("QueryFullMedicalHistory - 解析响应失败: %v\n", err)
			appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("解析响应失败: %v", err))
			return
		}

		fmt.Printf("QueryFullMedicalHistory - 解析后的数据: %+v\n", data)
		appG.Response(http.StatusOK, "成功", data)
		return
	}
	
	// 如果所有重试都失败了
	appG.Response(http.StatusInternalServerError, "失败", "查询病历失败，请稍后重试")
}
