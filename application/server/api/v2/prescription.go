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

func CreatePrescription(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(model.PrescriptionRequestBody)

	//解析Body参数
	if err := c.ShouldBind(body); err != nil {
		fmt.Printf("创建病历 - 参数解析失败: %v\n", err)
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数解析错误: %s", err.Error()))
		return
	}
	
	fmt.Printf("创建病历 - 接收参数: Doctor=%s, Patient=%s, ChiefComplaint=%s, PresentIllness=%s, PhysicalExam=%s, Diagnosis=%s, DrugName=%s, DrugAmount=%s, MedicalAdvice=%s, Hospital=%s, Comment=%s\n", 
		body.Doctor, body.Patient, body.ChiefComplaint, body.PresentIllness, body.PhysicalExam, body.Diagnosis, body.DrugName, body.DrugAmount, body.MedicalAdvice, body.Hospital, body.Comment)
	
	if body.Doctor == "" || body.Patient == "" || body.Diagnosis == "" || body.Hospital == "" {
		fmt.Printf("创建病历 - 参数存在空值\n")
		appG.Response(http.StatusBadRequest, "失败", "参数存在空值（医生、病人、诊断、医院不能为空）")
		return
	}
	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, []byte(body.Doctor))
	bodyBytes = append(bodyBytes, []byte(body.Patient))
	bodyBytes = append(bodyBytes, []byte(body.ChiefComplaint))
	bodyBytes = append(bodyBytes, []byte(body.PresentIllness))
	bodyBytes = append(bodyBytes, []byte(body.PhysicalExam))
	bodyBytes = append(bodyBytes, []byte(body.Diagnosis))
	bodyBytes = append(bodyBytes, []byte(body.DrugName))
	bodyBytes = append(bodyBytes, []byte(body.DrugAmount))
	bodyBytes = append(bodyBytes, []byte(body.MedicalAdvice))
	bodyBytes = append(bodyBytes, []byte(body.Hospital))
	bodyBytes = append(bodyBytes, []byte(body.Comment))

	fmt.Printf("创建病历 - 调用区块链智能合约\n")
	
	// 添加重试机制处理各种区块链错误
	maxRetries := 5 // 增加重试次数
	for i := 0; i < maxRetries; i++ {
		fmt.Printf("CreatePrescription - 尝试创建 (第%d次)\n", i+1)
		
		// 调用智能合约
		resp, err := bc.ChannelExecute("createPrescription", bodyBytes)
		if err != nil {
			fmt.Printf("CreatePrescription - 第%d次尝试失败: %v\n", i+1, err)
			
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
				fmt.Printf("CreatePrescription - 检测到可重试错误，等待后重试\n")
				// 增加等待时间，避免频繁重试
				time.Sleep(time.Duration(i+1) * 200 * time.Millisecond)
				continue
			}
			
			// 如果不是可重试错误或者已经是最后一次尝试，直接返回错误
			appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("区块链调用失败: %v", err))
			return
		}
		
		// 成功了，处理响应并返回
		fmt.Printf("CreatePrescription - 第%d次尝试成功\n", i+1)
		
		var data map[string]interface{}
		if err := json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
			fmt.Printf("CreatePrescription - 解析响应失败: %v\n", err)
			appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("解析响应失败: %v", err))
			return
		}

		fmt.Printf("CreatePrescription - 创建病历成功\n")
		appG.Response(http.StatusOK, "成功", data)
		return
	}
	
	// 如果所有重试都失败了
	appG.Response(http.StatusInternalServerError, "失败", "创建病历失败，请稍后重试")
}

func QueryPrescriptionList(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(model.PrescriptionQueryRequestBody)
	//解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	var bodyBytes [][]byte
	if body.Patient != "" {
		bodyBytes = append(bodyBytes, []byte(body.Patient))
	}
	// 添加重试机制处理各种区块链错误
	maxRetries := 5 // 增加重试次数
	for i := 0; i < maxRetries; i++ {
		fmt.Printf("QueryPrescriptionList - 尝试查询 (第%d次)\n", i+1)
		
		resp, err := bc.ChannelQuery("queryPrescription", bodyBytes)
		if err != nil {
			fmt.Printf("QueryPrescriptionList - 第%d次尝试失败: %v\n", i+1, err)
			
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
				fmt.Printf("QueryPrescriptionList - 检测到可重试错误，等待后重试\n")
				// 增加等待时间，避免频繁重试
				time.Sleep(time.Duration(i+1) * 200 * time.Millisecond)
				continue
			}
			
			// 如果不是可重试错误或者已经是最后一次尝试，直接返回错误
			appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("区块链调用失败: %v", err))
			return
		}
		
		// 成功了，处理响应并返回
		fmt.Printf("QueryPrescriptionList - 第%d次尝试成功\n", i+1)
		
		// 反序列化json
		var data []map[string]interface{}
		if err := json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
			fmt.Printf("QueryPrescriptionList - 解析响应失败: %v\n", err)
			appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("解析响应失败: %v", err))
			return
		}

		fmt.Printf("QueryPrescriptionList - 查询成功，返回 %d 条记录\n", len(data))
		appG.Response(http.StatusOK, "成功", data)
		return
	}
	
	// 如果所有重试都失败了
	appG.Response(http.StatusInternalServerError, "失败", "查询病历列表失败，请稍后重试")
}
