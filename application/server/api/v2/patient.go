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
	
	fmt.Printf("获取到 %d 条病历，时间范围: %s\n", len(prescriptions), timeRange)
	
	// 根据时间范围筛选
	prescriptions = filterByTimeRange(prescriptions, timeRange)
	
	fmt.Printf("筛选后 %d 条病历\n", len(prescriptions))
	
	// 计算健康档案统计
	profile := calculateHealthProfile(prescriptions, query.AccountID)
	
	fmt.Printf("统计完成: 病历=%d, 授权医生=%d, 订单=%d\n", 
		profile.TotalPrescriptions, profile.AuthorizedDoctors, profile.TotalOrders)
	
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

// calculateHealthProfile 计算健康档案统计
func calculateHealthProfile(prescriptions []Prescription, accountID string) HealthProfile {
	profile := HealthProfile{
		TotalPrescriptions: len(prescriptions),
		TotalVisits:        len(prescriptions),
		Timeline:           make([]TimelineItem, 0),
		DiseaseStats:       make([]DiseaseStatItem, 0),
		HospitalStats:      make([]HospitalStatItem, 0),
	}
	
	// 疾病统计 map
	diseaseMap := make(map[string]int)
	// 医院统计 map
	hospitalMap := make(map[string]int)
	
	// 按时间倒序排序
	sort.Slice(prescriptions, func(i, j int) bool {
		return prescriptions[i].Created > prescriptions[j].Created
	})
	
	// 生成时间线和统计数据
	for _, presc := range prescriptions {
		// 添加到时间线
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
		hospitalName := presc.OrganizationName
		if hospitalName != "" {
			hospitalMap[hospitalName]++
		}
	}
	
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
func getAuthorizedDoctorsCount(accountID string) int {
	// 从区块链查询授权记录 - 添加重试机制
	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, []byte(accountID))
	
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		resp, err := bc.ChannelQuery("queryAuthorizationsByPatient", bodyBytes)
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
			return 0
		}
		
		var authorizations []interface{}
		if err = json.Unmarshal(resp.Payload, &authorizations); err != nil {
			return 0
		}
		
		return len(authorizations)
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
