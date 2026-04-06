package v2

import (
	bc "application/blockchain"
	"application/db"
	"application/pkg/app"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RecordPrescriptionAccess 记录病历访问（在查看病历时调用）
func RecordPrescriptionAccess(c *gin.Context) {
	appG := app.Gin{C: c}

	type AccessRequest struct {
		PrescriptionID   string `json:"prescription_id"`
		PrescriptionNo   string `json:"prescription_no"`
		PatientID        string `json:"patient_id"`
		PatientName      string `json:"patient_name"`
		AccessorID       string `json:"accessor_id"`
		AccessorName     string `json:"accessor_name"`
		AccessorRole     string `json:"accessor_role"`
		AccessorOrg      string `json:"accessor_organization"`
		AccessorOrgName  string `json:"accessor_organization_name"`
		AccessType       string `json:"access_type"` // view, edit, download
		AccessReason     string `json:"access_reason"`
	}

	var req AccessRequest
	if err := c.ShouldBind(&req); err != nil {
		fmt.Printf("❌ 参数解析错误: %v\n", err)
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数解析错误: %s", err.Error()))
		return
	}

	fmt.Printf("📥 收到访问记录请求:\n")
	fmt.Printf("   病历ID: %s\n", req.PrescriptionID)
	fmt.Printf("   病历编号: %s\n", req.PrescriptionNo)
	fmt.Printf("   患者ID: %s, 患者姓名: %s\n", req.PatientID, req.PatientName)
	fmt.Printf("   访问者ID: %s, 访问者姓名: %s\n", req.AccessorID, req.AccessorName)
	fmt.Printf("   访问者角色: %s, 访问者组织: %s\n", req.AccessorRole, req.AccessorOrgName)

	// 如果访问者是患者本人，不记录（患者查看自己的病历不需要溯源）
	if req.AccessorID == req.PatientID {
		fmt.Printf("⏭️  患者本人访问，无需记录\n")
		appG.Response(http.StatusOK, "成功", map[string]interface{}{
			"message": "患者本人访问，无需记录",
		})
		return
	}

	// 生成日志ID
	logID := uuid.New().String()

	// 调用区块链记录访问日志
	args := [][]byte{
		[]byte(logID),
		[]byte(req.PrescriptionID),
		[]byte(req.PrescriptionNo),
		[]byte(req.PatientID),
		[]byte(req.PatientName),
		[]byte(req.AccessorID),
		[]byte(req.AccessorName),
		[]byte(req.AccessorRole),
		[]byte(req.AccessorOrg),
		[]byte(req.AccessorOrgName),
		[]byte(req.AccessType),
		[]byte(req.AccessReason),
	}

	response, err := bc.ChannelExecute("recordPrescriptionAccess", args)
	if err != nil {
		fmt.Printf("❌ 区块链记录失败: %v\n", err)
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("区块链记录失败: %v", err))
		return
	}

	txID := string(response.TransactionID)
	fmt.Printf("✅ 区块链记录成功 - TxID: %s\n", txID)

	// 同时保存到MySQL数据库（用于快速查询）
	accessLog := &db.PrescriptionAccessLog{
		LogID:                    logID,
		PrescriptionID:           req.PrescriptionID,
		PrescriptionNo:           req.PrescriptionNo,
		PatientID:                req.PatientID,
		PatientName:              req.PatientName,
		AccessorID:               req.AccessorID,
		AccessorName:             req.AccessorName,
		AccessorRole:             req.AccessorRole,
		AccessorOrganization:     req.AccessorOrg,
		AccessorOrganizationName: req.AccessorOrgName,
		AccessType:               req.AccessType,
		AccessReason:             req.AccessReason,
		IPAddress:                c.ClientIP(),
		UserAgent:                c.Request.UserAgent(),
		TxID:                     txID,
	}

	err = db.CreateAccessLog(accessLog)
	if err != nil {
		fmt.Printf("⚠️  MySQL保存失败: %v (区块链已记录)\n", err)
	} else {
		fmt.Printf("✅ MySQL保存成功\n")
	}

	appG.Response(http.StatusOK, "成功", map[string]interface{}{
		"log_id":  logID,
		"tx_id":   txID,
		"message": "访问日志已记录到区块链",
	})
}

// GetMyAccessLogs 获取我的病历访问日志（患者查看）- 从区块链查询
func GetMyAccessLogs(c *gin.Context) {
	appG := app.Gin{C: c}

	type QueryParams struct {
		PatientID string `json:"patient_id" form:"patient_id"`
		Page      int    `json:"page" form:"page"`
		PageSize  int    `json:"page_size" form:"page_size"`
	}

	params := QueryParams{
		Page:     1,
		PageSize: 20,
	}

	if err := c.ShouldBind(&params); err != nil {
		fmt.Printf("❌ 参数解析错误: %v\n", err)
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数解析错误: %s", err.Error()))
		return
	}

	if params.PatientID == "" {
		fmt.Printf("❌ 患者ID为空\n")
		appG.Response(http.StatusBadRequest, "失败", "患者ID不能为空")
		return
	}

	fmt.Printf("📥 查询访问日志: 患者ID=%s, 页码=%d, 每页=%d\n", params.PatientID, params.Page, params.PageSize)

	// 优先从MySQL查询（快速）
	logs, total, err := db.GetAccessLogsByPatient(params.PatientID, params.Page, params.PageSize)
	if err != nil {
		fmt.Printf("❌ MySQL查询失败: %v，尝试从区块链查询\n", err)
		
		// MySQL失败，从区块链查询
		args := [][]byte{[]byte(params.PatientID)}
		response, err := bc.ChannelQuery("queryAccessLogsByPatient", args)
		if err != nil {
			fmt.Printf("❌ 区块链查询也失败: %v\n", err)
			appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("查询失败: %v", err))
			return
		}

		// 解析区块链返回的数据
		var blockchainLogs []map[string]interface{}
		if err := json.Unmarshal(response.Payload, &blockchainLogs); err != nil {
			fmt.Printf("❌ 解析区块链数据失败: %v\n", err)
			appG.Response(http.StatusInternalServerError, "失败", "解析数据失败")
			return
		}

		fmt.Printf("✅ 从区块链查询成功: 找到 %d 条记录\n", len(blockchainLogs))
		
		appG.Response(http.StatusOK, "成功", map[string]interface{}{
			"logs":      blockchainLogs,
			"total":     len(blockchainLogs),
			"page":      params.Page,
			"page_size": params.PageSize,
			"source":    "blockchain",
		})
		return
	}

	fmt.Printf("✅ 从MySQL查询成功: 找到 %d 条记录 (总共 %d 条)\n", len(logs), total)

	// 转换为前端需要的格式
	var data []map[string]interface{}
	for _, log := range logs {
		logMap := map[string]interface{}{
			"id":                         log.ID,
			"log_id":                     log.LogID,
			"prescription_id":            log.PrescriptionID,
			"prescription_no":            log.PrescriptionNo,
			"accessor_name":              log.AccessorName,
			"accessor_role":              log.AccessorRole,
			"accessor_organization_name": log.AccessorOrganizationName,
			"access_type":                log.AccessType,
			"access_reason":              log.AccessReason,
			"accessed_at":                log.AccessedAt.Format("2006-01-02 15:04:05"),
			"ip_address":                 log.IPAddress,
			"tx_id":                      log.TxID,
		}
		data = append(data, logMap)
	}

	appG.Response(http.StatusOK, "成功", map[string]interface{}{
		"logs":      data,
		"total":     total,
		"page":      params.Page,
		"page_size": params.PageSize,
		"source":    "mysql",
	})
}

// GetAccessStatistics 获取访问统计
func GetAccessStatistics(c *gin.Context) {
	appG := app.Gin{C: c}

	patientID := c.Query("patient_id")
	if patientID == "" {
		appG.Response(http.StatusBadRequest, "失败", "患者ID不能为空")
		return
	}

	stats, err := db.GetAccessStatistics(patientID)
	if err != nil {
		fmt.Printf("查询访问统计失败: %v\n", err)
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("查询失败: %v", err))
		return
	}

	appG.Response(http.StatusOK, "成功", stats)
}
