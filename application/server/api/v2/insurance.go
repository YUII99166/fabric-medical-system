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

func CreateInsuranceCover(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(model.InsuranceCoverRequestBody)

	//解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	if body.Prescription == "" || body.Patient == "" || body.Status == "" {
		appG.Response(http.StatusBadRequest, "失败", "参数存在空值")
		return
	}
	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, []byte(body.Prescription))
	bodyBytes = append(bodyBytes, []byte(body.Patient))
	bodyBytes = append(bodyBytes, []byte(body.Status))
	//bodyBytes = append(bodyBytes, []byte(strconv.FormatFloat(body.TotalArea, 'E', -1, 64)))

	// 调用智能合约 - 添加重试机制
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		resp, err := bc.ChannelExecute("createInsuranceCover", bodyBytes)
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
			appG.Response(http.StatusInternalServerError, "失败", err.Error())
			return
		}
		
		var data map[string]interface{}
		if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
			appG.Response(http.StatusInternalServerError, "失败", err.Error())
			return
		}
		appG.Response(http.StatusOK, "成功", data)
		return
	}
}

func QueryInsuranceCoverList(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(model.InsuranceCoverQueryRequestBody)
	//解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	var bodyBytes [][]byte
	if body.Patient != "" {
		bodyBytes = append(bodyBytes, []byte(body.Patient))
	}
	if body.InsuranceCover != "" {
		bodyBytes = append(bodyBytes, []byte(body.InsuranceCover))
	}
	//调用智能合约 - 添加重试机制
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		resp, err := bc.ChannelQuery("queryInsuranceCover", bodyBytes)
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
			appG.Response(http.StatusInternalServerError, "失败", err.Error())
			return
		}
		
		// 反序列化json
		var data []map[string]interface{}
		if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
			appG.Response(http.StatusInternalServerError, "失败", err.Error())
			return
		}
		appG.Response(http.StatusOK, "成功", data)
		return
	}
}

func UpdateInsuranceCover(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(model.UpdateInsuranceCoverRequestBody)
	//解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, []byte(body.InsuranceCover))
	bodyBytes = append(bodyBytes, []byte(body.InsuranceID))
	bodyBytes = append(bodyBytes, []byte(body.Status))
	bodyBytes = append(bodyBytes, []byte(body.Patient))

	//调用智能合约 - 添加重试机制
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		resp, err := bc.ChannelExecute("updateInsuranceCover", bodyBytes)
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
			appG.Response(http.StatusInternalServerError, "失败", err.Error())
			return
		}
		
		// 反序列化json
		var data map[string]interface{}
		if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
			appG.Response(http.StatusInternalServerError, "失败", err.Error())
			return
		}
		appG.Response(http.StatusOK, "成功", data)
		return
	}
}

func DeleteInsuranceCover(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(model.UpdateInsuranceCoverRequestBody)
	//解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, []byte(body.InsuranceCover))
	bodyBytes = append(bodyBytes, []byte(body.InsuranceID))
	bodyBytes = append(bodyBytes, []byte(body.Status))
	bodyBytes = append(bodyBytes, []byte(body.Patient))

	//调用智能合约 - 添加重试机制
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		resp, err := bc.ChannelExecute("deleteInsuranceCover", bodyBytes)
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
			appG.Response(http.StatusInternalServerError, "失败", err.Error())
			return
		}
		
		// 反序列化json
		var data map[string]interface{}
		if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
			appG.Response(http.StatusInternalServerError, "失败", err.Error())
			return
		}
		appG.Response(http.StatusOK, "成功", data)
		return
	}
}
