package v2

import (
	bc "application/blockchain"
	"application/db"
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

func CreateDrugOrder(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(model.DrugOrderRequestBody)

	//解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	if body.DrugName == "" || body.DrugAmount == "" || body.Prescription == "" || body.Patient == "" || body.DrugStore == "" {
		appG.Response(http.StatusBadRequest, "失败", "参数存在空值")
		return
	}
	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, []byte(body.DrugName))
	bodyBytes = append(bodyBytes, []byte(body.DrugAmount))
	bodyBytes = append(bodyBytes, []byte(body.Prescription))
	bodyBytes = append(bodyBytes, []byte(body.Patient))
	bodyBytes = append(bodyBytes, []byte(body.DrugStore))
	//bodyBytes = append(bodyBytes, []byte(strconv.FormatFloat(body.TotalArea, 'E', -1, 64)))

	// 调用智能合约 - 添加重试机制
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		resp, err := bc.ChannelExecute("createDrugOrder", bodyBytes)
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

func QueryDrugOrderList(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(model.DrugOrderQueryRequestBody)
	//解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	var bodyBytes [][]byte
	if body.Patient != "" {
		bodyBytes = append(bodyBytes, []byte(body.Patient))
	}
	//调用智能合约 - 添加重试机制
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		resp, err := bc.ChannelQuery("queryDrugOrder", bodyBytes)
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
		
		// 补充药店名称信息
		for i := range data {
			if drugStoreID, ok := data[i]["drug_store"].(string); ok && drugStoreID != "" {
				// 尝试从数据库查询药店名称
				var drugStoreName string
				err := db.DB.QueryRow("SELECT account_name FROM users WHERE account_id = ? OR username = ?", drugStoreID, drugStoreID).Scan(&drugStoreName)
				if err == nil && drugStoreName != "" {
					data[i]["drug_store_name"] = drugStoreName
					data[i]["DrugStoreName"] = drugStoreName
				} else {
					// 如果查询失败，使用ID作为名称
					data[i]["drug_store_name"] = drugStoreID
					data[i]["DrugStoreName"] = drugStoreID
				}
			}
		}
		
		appG.Response(http.StatusOK, "成功", data)
		return
	}
}