package v2

import (
	bc "application/blockchain"
	"application/pkg/app"
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
)

// Prescription 病历结构（从区块链返回）
type Prescription struct {
	ID               string `json:"id"`
	Patient          string `json:"patient"`
	PatientName      string `json:"patient_name"`
	Diagnosis        string `json:"diagnosis"`
	Doctor           string `json:"doctor"`
	DoctorName       string `json:"doctor_name"`
	OrganizationID   string `json:"organization_id"`
	OrganizationName string `json:"organization_name"`
	Created          string `json:"created"`
}

// SupplementRecord 补充记录结构（从区块链返回）
type SupplementRecord struct {
	ID                     string `json:"id"`
	OriginalPrescriptionID string `json:"original_prescription_id"`
	RecordType             string `json:"record_type"`
	PatientID              string `json:"patient_id"`
	PatientName            string `json:"patient_name"`
	Diagnosis              string `json:"diagnosis"`
	DoctorID               string `json:"doctor_id"`
	DoctorName             string `json:"doctor_name"`
	OrganizationID         string `json:"organization_id"`
	OrganizationName       string `json:"organization_name"`
	Department             string `json:"department"`
	Created                string `json:"created"`
}

// Authorization 授权记录结构（从区块链返回）
type Authorization struct {
	ID             string `json:"id"`
	PrescriptionID string `json:"prescription_id"`
	PatientID      string `json:"patient_id"`
	PatientName    string `json:"patient_name"`
	DoctorID       string `json:"doctor_id"`
	DoctorName     string `json:"doctor_name"`
	DoctorOrg      string `json:"doctor_org"`           // 医生所属组织ID
	DoctorOrgName  string `json:"doctor_org_name"`      // 医生所属组织名称
	Reason         string `json:"reason"`
	Status         string `json:"status"`
	RequestTime    string `json:"request_time"`         // 申请时间
	ResponseTime   string `json:"response_time"`        // 响应时间
	TxID           string `json:"tx_id"`
}

// HealthProfile 健康档案统计数据
type HealthProfile struct {
	TotalPrescriptions int                  `json:"totalPrescriptions"` // 病历总数
	TotalVisits        int                  `json:"totalVisits"`        // 就诊次数
	AuthorizedDoctors  int                  `json:"authorizedDoctors"`  // 授权医生数量
	TotalOrders        int                  `json:"totalOrders"`        // 药品订单数量
	Timeline           []TimelineItem       `json:"timeline"`           // 就诊时间线
	DiseaseStats       []DiseaseStatItem    `json:"diseaseStats"`       // 疾病统计
	HospitalStats      []HospitalStatItem   `json:"hospitalStats"`      // 医院统计
}

// TimelineItem 时间线项目
type TimelineItem struct {
	Date      string `json:"date"`      // 就诊日期
	Hospital  string `json:"hospital"`  // 医院名称
	Doctor    string `json:"doctor"`    // 医生姓名
	Diagnosis string `json:"diagnosis"` // 诊断
	PrescID   string `json:"prescId"`   // 病历ID
}

// DiseaseStatItem 疾病统计项目
type DiseaseStatItem struct {
	Disease string `json:"disease"` // 疾病名称
	Count   int    `json:"count"`   // 出现次数
}

// HospitalStatItem 医院统计项目
type HospitalStatItem struct {
	Hospital string `json:"hospital"` // 医院名称
	Count    int    `json:"count"`    // 就诊次数
}

// GetHealthProfile 获取健康档案统计
func GetHealthProfile(c *gin.Context) {
	appG := app.Gin{C: c}
	
	// 从请求参数获取区块链账户ID
	type AccountIDQuery struct {
		AccountID string `form:"accountId" binding:"required"`
	}
	query := new(AccountIDQuery)
	
	if err := c.ShouldBindQuery(query); err != nil {
		appG.Response(http.StatusBadRequest, "失败", "区块链账户ID不能为空")
		return
	}
	
	fmt.Printf("获取健康档案，区块链账户ID: %s\n", query.AccountID)
	
	// 获取时间范围参数
	timeRange := c.DefaultQuery("timeRange", "all")
	
	// 从区块链获取病历数据 - 使用区块链账户ID
	prescriptions, err := getPrescriptionsFromBlockchain(query.AccountID)
	if err != nil {
		fmt.Printf("获取病历数据失败: %v\n", err)
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("获取病历数据失败: %v", err))
		return
	}
	
	fmt.Printf("获取到 %d 条原始病历，时间范围: %s\n", len(prescriptions), timeRange)
	
	// 获取补充记录数据
	supplementRecords, err := getSupplementRecordsFromBlockchain(query.AccountID)
	if err != nil {
		fmt.Printf("获取补充记录失败: %v\n", err)
		// 补充记录获取失败不影响主流程，继续执行
		supplementRecords = []SupplementRecord{}
	}
	
	fmt.Printf("获取到 %d 条补充记录\n", len(supplementRecords))
	
	// 获取授权记录（作为转诊记录）
	authorizations, err := getAuthorizationsFromBlockchain(query.AccountID)
	if err != nil {
		fmt.Printf("获取授权记录失败: %v\n", err)
		// 授权记录获取失败不影响主流程，继续执行
		authorizations = []Authorization{}
	}
	
	fmt.Printf("获取到 %d 条授权记录（转诊）\n", len(authorizations))
	
	// 根据时间范围筛选
	prescriptions = filterByTimeRange(prescriptions, timeRange)
	supplementRecords = filterSupplementRecordsByTimeRange(supplementRecords, timeRange)
	authorizations = filterAuthorizationsByTimeRange(authorizations, timeRange)
	
	fmt.Printf("筛选后: %d 条原始病历, %d 条补充记录, %d 条授权记录\n", len(prescriptions), len(supplementRecords), len(authorizations))
	
	// 计算健康档案统计（包含补充记录和授权记录）
	profile := calculateHealthProfile(prescriptions, supplementRecords, authorizations, query.AccountID)
	
	fmt.Printf("统计完成: 病历=%d, 就诊次数=%d, 授权医生=%d, 订单=%d\n", 
		profile.TotalPrescriptions, profile.TotalVisits, profile.AuthorizedDoctors, profile.TotalOrders)
	
	appG.Response(http.StatusOK, "成功", profile)
}

// getPrescriptionsFromBlockchain 从区块链获取病历数据
func getPrescriptionsFromBlockchain(patientID string) ([]Prescription, error) {
	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, []byte(patientID))
	
	// 添加重试机制
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		// 调用智能合约查询病历 - 使用 queryPrescription 函数
		// 注意：这里传入的应该是区块链上的患者ID，不是数据库的account_name
		resp, err := bc.ChannelQuery("queryPrescription", bodyBytes)
		if err != nil {
			fmt.Printf("区块链查询失败 (第%d次): %v\n", i+1, err)
			
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
				fmt.Printf("检测到可重试错误，等待后重试\n")
				time.Sleep(time.Duration(i+1) * 200 * time.Millisecond)
				continue
			}
			
			return nil, err
		}
		
		// 打印原始响应
		responseStr := string(resp.Payload)
		fmt.Printf("区块链响应: %s\n", responseStr)
		
		// 如果响应是 "null" 或空，返回空数组
		if responseStr == "null" || responseStr == "" || responseStr == "[]" {
			fmt.Printf("没有找到病历数据\n")
			return []Prescription{}, nil
		}
		
		var prescriptions []Prescription
		if err = json.Unmarshal(resp.Payload, &prescriptions); err != nil {
			fmt.Printf("JSON解析失败: %v\n", err)
			// 如果解析失败，返回空数组而不是错误
			return []Prescription{}, nil
		}
		
		fmt.Printf("成功获取 %d 条病历\n", len(prescriptions))
		return prescriptions, nil
	}
	
	return nil, fmt.Errorf("所有重试都失败了")
}

// getSupplementRecordsFromBlockchain 从区块链获取补充记录数据
func getSupplementRecordsFromBlockchain(patientID string) ([]SupplementRecord, error) {
	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, []byte(patientID))
	
	// 添加重试机制
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		// 调用智能合约查询补充记录 - 使用 querySupplementRecords 函数
		resp, err := bc.ChannelQuery("querySupplementRecords", bodyBytes)
		if err != nil {
			fmt.Printf("查询补充记录失败 (第%d次): %v\n", i+1, err)
			
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
				fmt.Printf("检测到可重试错误，等待后重试\n")
				time.Sleep(time.Duration(i+1) * 200 * time.Millisecond)
				continue
			}
			
			return nil, err
		}
		
		// 打印原始响应
		responseStr := string(resp.Payload)
		fmt.Printf("补充记录响应: %s\n", responseStr)
		
		// 如果响应是 "null" 或空，返回空数组
		if responseStr == "null" || responseStr == "" || responseStr == "[]" {
			fmt.Printf("没有找到补充记录\n")
			return []SupplementRecord{}, nil
		}
		
		var records []SupplementRecord
		if err = json.Unmarshal(resp.Payload, &records); err != nil {
			fmt.Printf("补充记录JSON解析失败: %v\n", err)
			// 如果解析失败，返回空数组而不是错误
			return []SupplementRecord{}, nil
		}
		
		fmt.Printf("成功获取 %d 条补充记录\n", len(records))
		return records, nil
	}
	
	return nil, fmt.Errorf("所有重试都失败了")
}

// getAuthorizationsFromBlockchain 从区块链获取授权记录数据
func getAuthorizationsFromBlockchain(patientID string) ([]Authorization, error) {
	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, []byte(patientID))
	bodyBytes = append(bodyBytes, []byte("patient")) // 角色参数
	
	// 添加重试机制
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		// 调用智能合约查询授权记录 - 使用 queryAccessRequests 函数
		resp, err := bc.ChannelQuery("queryAccessRequests", bodyBytes)
		if err != nil {
			fmt.Printf("查询授权记录失败 (第%d次): %v\n", i+1, err)
			
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
				fmt.Printf("检测到可重试错误，等待后重试\n")
				time.Sleep(time.Duration(i+1) * 200 * time.Millisecond)
				continue
			}
			
			return nil, err
		}
		
		// 打印原始响应
		responseStr := string(resp.Payload)
		fmt.Printf("授权记录响应: %s\n", responseStr)
		
		// 如果响应是 "null" 或空，返回空数组
		if responseStr == "null" || responseStr == "" || responseStr == "[]" {
			fmt.Printf("没有找到授权记录\n")
			return []Authorization{}, nil
		}
		
		var authorizations []Authorization
		if err = json.Unmarshal(resp.Payload, &authorizations); err != nil {
			fmt.Printf("授权记录JSON解析失败: %v\n", err)
			// 如果解析失败，返回空数组而不是错误
			return []Authorization{}, nil
		}
		
		// 只保留已批准的授权记录
		var approvedAuthorizations []Authorization
		for _, auth := range authorizations {
			if auth.Status == "approved" {
				approvedAuthorizations = append(approvedAuthorizations, auth)
			}
		}
		
		fmt.Printf("成功获取 %d 条授权记录（已批准: %d 条）\n", len(authorizations), len(approvedAuthorizations))
		return approvedAuthorizations, nil
	}
	
	return nil, fmt.Errorf("所有重试都失败了")
}

// filterByTimeRange 根据时间范围筛选病历
func filterByTimeRange(prescriptions []Prescription, timeRange string) []Prescription {
	if timeRange == "all" {
		return prescriptions
	}
	
	now := time.Now()
	var cutoffDate time.Time
	
	switch timeRange {
	case "1m":
		cutoffDate = now.AddDate(0, -1, 0) // 一个月前
	case "3m":
		cutoffDate = now.AddDate(0, -3, 0) // 三个月前
	case "1y":
		cutoffDate = now.AddDate(-1, 0, 0) // 一年前
	default:
		return prescriptions
	}
	
	var filtered []Prescription
	for _, presc := range prescriptions {
		// 解析病历创建时间，支持多种格式
		var prescTime time.Time
		var err error
		
		// 尝试解析 "2006-01-02 15:04:05" 格式
		prescTime, err = time.Parse("2006-01-02 15:04:05", presc.Created)
		if err != nil {
			// 尝试解析 "2006-01-02" 格式
			prescTime, err = time.Parse("2006-01-02", presc.Created)
			if err != nil {
				fmt.Printf("无法解析时间: %s, 错误: %v\n", presc.Created, err)
				continue
			}
		}
		
		// 使用 After 或 Equal（即 >=）
		if prescTime.After(cutoffDate) || prescTime.Equal(cutoffDate) {
			filtered = append(filtered, presc)
		}
	}
	
	fmt.Printf("时间筛选: 原始=%d, 筛选后=%d, 范围=%s, 截止日期=%s\n", 
		len(prescriptions), len(filtered), timeRange, cutoffDate.Format("2006-01-02"))
	
	return filtered
}

// filterSupplementRecordsByTimeRange 根据时间范围筛选补充记录
func filterSupplementRecordsByTimeRange(records []SupplementRecord, timeRange string) []SupplementRecord {
	if timeRange == "all" {
		return records
	}
	
	now := time.Now()
	var cutoffDate time.Time
	
	switch timeRange {
	case "1m":
		cutoffDate = now.AddDate(0, -1, 0)
	case "3m":
		cutoffDate = now.AddDate(0, -3, 0)
	case "1y":
		cutoffDate = now.AddDate(-1, 0, 0)
	default:
		return records
	}
	
	var filtered []SupplementRecord
	for _, record := range records {
		var recordTime time.Time
		var err error
		
		recordTime, err = time.Parse("2006-01-02 15:04:05", record.Created)
		if err != nil {
			recordTime, err = time.Parse("2006-01-02", record.Created)
			if err != nil {
				fmt.Printf("无法解析补充记录时间: %s, 错误: %v\n", record.Created, err)
				continue
			}
		}
		
		if recordTime.After(cutoffDate) || recordTime.Equal(cutoffDate) {
			filtered = append(filtered, record)
		}
	}
	
	fmt.Printf("补充记录时间筛选: 原始=%d, 筛选后=%d\n", len(records), len(filtered))
	
	return filtered
}

// filterAuthorizationsByTimeRange 根据时间范围筛选授权记录
func filterAuthorizationsByTimeRange(authorizations []Authorization, timeRange string) []Authorization {
	if timeRange == "all" {
		return authorizations
	}
	
	now := time.Now()
	var cutoffDate time.Time
	
	switch timeRange {
	case "1m":
		cutoffDate = now.AddDate(0, -1, 0)
	case "3m":
		cutoffDate = now.AddDate(0, -3, 0)
	case "1y":
		cutoffDate = now.AddDate(-1, 0, 0)
	default:
		return authorizations
	}
	
	var filtered []Authorization
	for _, auth := range authorizations {
		var authTime time.Time
		var err error
		
		// 使用 RequestTime 字段
		authTime, err = time.Parse("2006-01-02 15:04:05", auth.RequestTime)
		if err != nil {
			authTime, err = time.Parse("2006-01-02", auth.RequestTime)
			if err != nil {
				fmt.Printf("无法解析授权记录时间: %s, 错误: %v\n", auth.RequestTime, err)
				continue
			}
		}
		
		if authTime.After(cutoffDate) || authTime.Equal(cutoffDate) {
			filtered = append(filtered, auth)
		}
	}
	
	fmt.Printf("授权记录时间筛选: 原始=%d, 筛选后=%d\n", len(authorizations), len(filtered))
	
	return filtered
}

// calculateHealthProfile 计算健康档案统计
func calculateHealthProfile(prescriptions []Prescription, supplementRecords []SupplementRecord, authorizations []Authorization, accountID string) HealthProfile {
	profile := HealthProfile{
		TotalPrescriptions: len(prescriptions), // 只统计原始病历数量
		TotalVisits:        len(prescriptions) + len(supplementRecords) + len(authorizations), // 就诊次数 = 原始病历 + 补充记录 + 授权记录
		Timeline:           make([]TimelineItem, 0),
		DiseaseStats:       make([]DiseaseStatItem, 0),
		HospitalStats:      make([]HospitalStatItem, 0),
	}
	
	// 疾病统计 map
	diseaseMap := make(map[string]int)
	// 医院统计 map
	hospitalMap := make(map[string]int)
	
	// 将原始病历添加到时间线
	for _, presc := range prescriptions {
		timeline := TimelineItem{
			Date:      presc.Created,
			Hospital:  presc.OrganizationName,
			Doctor:    presc.DoctorName,
			Diagnosis: presc.Diagnosis,
			PrescID:   presc.ID,
		}
		profile.Timeline = append(profile.Timeline, timeline)
		
		// 统计疾病
		if presc.Diagnosis != "" {
			diseaseMap[presc.Diagnosis]++
		}
		
		// 统计医院
		if presc.OrganizationName != "" {
			hospitalMap[presc.OrganizationName]++
		}
	}
	
	// 将补充记录添加到时间线
	for _, record := range supplementRecords {
		// 获取记录类型的中文描述
		recordTypeDesc := getRecordTypeDescription(record.RecordType)
		
		timeline := TimelineItem{
			Date:      record.Created,
			Hospital:  record.OrganizationName,
			Doctor:    record.DoctorName,
			Diagnosis: fmt.Sprintf("[%s] %s", recordTypeDesc, record.Diagnosis),
			PrescID:   record.ID,
		}
		profile.Timeline = append(profile.Timeline, timeline)
		
		// 统计疾病（补充记录也计入）
		if record.Diagnosis != "" {
			diseaseMap[record.Diagnosis]++
		}
		
		// 统计医院（补充记录也计入）
		if record.OrganizationName != "" {
			hospitalMap[record.OrganizationName]++
		}
	}
	
	// 将授权记录添加到时间线（作为转诊记录）
	for _, auth := range authorizations {
		timeline := TimelineItem{
			Date:      auth.RequestTime,      // 使用 RequestTime
			Hospital:  auth.DoctorOrgName,    // 使用 DoctorOrgName
			Doctor:    auth.DoctorName,
			Diagnosis: fmt.Sprintf("[授权转诊] %s", auth.Reason),
			PrescID:   auth.ID,
		}
		profile.Timeline = append(profile.Timeline, timeline)
		
		// 统计医院（授权记录也计入）
		if auth.DoctorOrgName != "" {
			hospitalMap[auth.DoctorOrgName]++
		}
	}
	
	// 按时间倒序排序时间线
	sort.Slice(profile.Timeline, func(i, j int) bool {
		return profile.Timeline[i].Date > profile.Timeline[j].Date
	})
	
	// 转换疾病统计为数组
	for disease, count := range diseaseMap {
		profile.DiseaseStats = append(profile.DiseaseStats, DiseaseStatItem{
			Disease: disease,
			Count:   count,
		})
	}
	
	// 按次数排序疾病统计
	sort.Slice(profile.DiseaseStats, func(i, j int) bool {
		return profile.DiseaseStats[i].Count > profile.DiseaseStats[j].Count
	})
	
	// 转换医院统计为数组
	for hospital, count := range hospitalMap {
		profile.HospitalStats = append(profile.HospitalStats, HospitalStatItem{
			Hospital: hospital,
			Count:    count,
		})
	}
	
	// 按次数排序医院统计
	sort.Slice(profile.HospitalStats, func(i, j int) bool {
		return profile.HospitalStats[i].Count > profile.HospitalStats[j].Count
	})
	
	// 获取授权医生数量
	profile.AuthorizedDoctors = getAuthorizedDoctorsCount(accountID)
	
	// 获取药品订单数量
	profile.TotalOrders = getDrugOrdersCount(accountID)
	
	return profile
}

// getRecordTypeDescription 获取记录类型的中文描述
func getRecordTypeDescription(recordType string) string {
	typeMap := map[string]string{
		"consultation": "转诊",
		"followup":     "复诊",
		"emergency":    "急诊",
	}
	
	if desc, ok := typeMap[recordType]; ok {
		return desc
	}
	return recordType
}

// getOrganizationName 根据组织 MSP 获取组织名称
func getOrganizationName(msp string) string {
	orgMap := map[string]string{
		"TaobaoMSP":    "协和医院",
		"JDMSP":        "301医院",
		"WenjinMSP":    "成都温江社区医疗中心",
		"RegCenterMSP": "监管中心",
	}
	
	if name, ok := orgMap[msp]; ok {
		return name
	}
	return msp
}

// getAuthorizedDoctorsCount 获取授权医生数量
// getAuthorizedDoctorsCount 获取授权医生数量
func getAuthorizedDoctorsCount(accountID string) int {
	// 从区块链查询授权记录 - 使用 queryAccessRequests 函数
	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, []byte(accountID))
	bodyBytes = append(bodyBytes, []byte("patient")) // 角色参数
	
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		resp, err := bc.ChannelQuery("queryAccessRequests", bodyBytes)
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
			fmt.Printf("查询授权医生数量失败: %v\n", err)
			return 0
		}
		
		// 解析授权记录
		var authorizations []Authorization
		if err = json.Unmarshal(resp.Payload, &authorizations); err != nil {
			fmt.Printf("解析授权记录失败: %v\n", err)
			return 0
		}
		
		// 统计已批准的授权中不同医生的数量
		doctorSet := make(map[string]bool)
		for _, auth := range authorizations {
			if auth.Status == "approved" && auth.DoctorID != "" {
				doctorSet[auth.DoctorID] = true
			}
		}
		
		count := len(doctorSet)
		fmt.Printf("授权医生数量: %d (总授权记录: %d, 已批准: %d)\n", count, len(authorizations), len(doctorSet))
		return count
	}
	
	return 0
}

// getDrugOrdersCount 获取药品订单数量
func getDrugOrdersCount(patientID string) int {
	// 从区块链查询药品订单 - 使用 queryDrugOrder 函数，传入患者ID - 添加重试机制
	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, []byte(patientID))
	
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		resp, err := bc.ChannelQuery("queryDrugOrder", bodyBytes)
		if err != nil {
			fmt.Printf("查询药品订单失败 (第%d次): %v\n", i+1, err)
			
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
			return 0
		}
		
		// 打印响应
		responseStr := string(resp.Payload)
		fmt.Printf("药品订单响应: %s\n", responseStr)
		
		// 如果响应是 "null" 或空，返回 0
		if responseStr == "null" || responseStr == "" || responseStr == "[]" {
			return 0
		}
		
		var orders []interface{}
		if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &orders); err != nil {
			fmt.Printf("解析药品订单失败: %v\n", err)
			return 0
		}
		
		fmt.Printf("患者 %s 有 %d 个药品订单\n", patientID, len(orders))
		return len(orders)
	}
	
	return 0
}
