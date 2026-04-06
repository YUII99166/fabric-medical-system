package service

import (
	"application/db"
	"fmt"
	"time"
)

// GetAllianceActivities 获取联盟活动列表（从业务表读取）
func GetAllianceActivities(page, pageSize int, organization, activityType, startDate, endDate string) ([]map[string]interface{}, int, error) {
	var activities []map[string]interface{}

	// 构建日期过滤条件
	dateFilter := ""
	if startDate != "" && endDate != "" {
		dateFilter = fmt.Sprintf(" AND created_at BETWEEN '%s 00:00:00' AND '%s 23:59:59'", startDate, endDate)
	}

	// 构建组织过滤条件
	orgFilter := ""
	if organization != "" {
		orgFilter = fmt.Sprintf(" AND organization_id = '%s'", organization)
	}

	// 根据活动类型查询不同的表
	if activityType == "" || activityType == "prescription" {
		// 查询病历创建活动
		query := fmt.Sprintf(`
			SELECT 
				'prescription' as type,
				CONCAT('医生 ', doctor_name, ' 为患者 ', patient_name, ' 创建了病历') as description,
				prescription_id as resource_id,
				doctor_name,
				patient_name,
				organization_id,
				organization_name,
				created_at
			FROM prescriptions
			WHERE 1=1 %s %s
			ORDER BY created_at DESC
		`, orgFilter, dateFilter)

		rows, err := db.DB.Query(query)
		if err != nil {
			return nil, 0, fmt.Errorf("查询病历活动失败: %v", err)
		}
		defer rows.Close()

		for rows.Next() {
			var activity map[string]interface{} = make(map[string]interface{})
			var actType, description, resourceID, doctorName, patientName, orgID, orgName string
			var createdAt time.Time

			err := rows.Scan(&actType, &description, &resourceID, &doctorName, &patientName, &orgID, &orgName, &createdAt)
			if err != nil {
				continue
			}

			activity["type"] = actType
			activity["description"] = description
			activity["resource_id"] = resourceID
			activity["doctor_name"] = doctorName
			activity["patient_name"] = patientName
			activity["organization_id"] = orgID
			activity["organization_name"] = orgName
			activity["department"] = ""
			activity["created_at"] = createdAt.Format("2006-01-02 15:04:05")

			activities = append(activities, activity)
		}
	}

	if activityType == "" || activityType == "supplement" {
		// 查询补充记录活动
		query := fmt.Sprintf(`
			SELECT 
				'supplement' as type,
				CONCAT('医生 ', doctor_name, ' 添加了', record_type, '记录') as description,
				record_id as resource_id,
				doctor_name,
				'' as patient_name,
				organization_id,
				organization_name,
				created_at
			FROM supplement_records
			WHERE 1=1 %s %s
			ORDER BY created_at DESC
		`, orgFilter, dateFilter)

		rows, err := db.DB.Query(query)
		if err != nil {
			return nil, 0, fmt.Errorf("查询补充记录活动失败: %v", err)
		}
		defer rows.Close()

		for rows.Next() {
			var activity map[string]interface{} = make(map[string]interface{})
			var actType, description, resourceID, doctorName, patientName, orgID, orgName string
			var createdAt time.Time

			err := rows.Scan(&actType, &description, &resourceID, &doctorName, &patientName, &orgID, &orgName, &createdAt)
			if err != nil {
				continue
			}

			activity["type"] = actType
			activity["description"] = description
			activity["resource_id"] = resourceID
			activity["doctor_name"] = doctorName
			activity["patient_name"] = patientName
			activity["organization_id"] = orgID
			activity["organization_name"] = orgName
			activity["department"] = ""
			activity["created_at"] = createdAt.Format("2006-01-02 15:04:05")

			activities = append(activities, activity)
		}
	}

	if activityType == "" || activityType == "access_request" {
		// 查询授权请求活动
		query := fmt.Sprintf(`
			SELECT 
				'access_request' as type,
				CONCAT('医生 ', doctor_name, ' 请求访问患者 ', patient_name, ' 的病历') as description,
				request_id as resource_id,
				doctor_name,
				patient_name,
				doctor_org as organization_id,
				doctor_org_name as organization_name,
				'' as department,
				created_at
			FROM access_requests
			WHERE 1=1 %s %s
			ORDER BY created_at DESC
		`, orgFilter, dateFilter)

		rows, err := db.DB.Query(query)
		if err != nil {
			return nil, 0, fmt.Errorf("查询授权请求活动失败: %v", err)
		}
		defer rows.Close()

		for rows.Next() {
			var activity map[string]interface{} = make(map[string]interface{})
			var actType, description, resourceID, doctorName, patientName, orgID, orgName, dept string
			var createdAt time.Time

			err := rows.Scan(&actType, &description, &resourceID, &doctorName, &patientName, &orgID, &orgName, &dept, &createdAt)
			if err != nil {
				continue
			}

			activity["type"] = actType
			activity["description"] = description
			activity["resource_id"] = resourceID
			activity["doctor_name"] = doctorName
			activity["patient_name"] = patientName
			activity["organization_id"] = orgID
			activity["organization_name"] = orgName
			activity["department"] = dept
			activity["created_at"] = createdAt.Format("2006-01-02 15:04:05")

			activities = append(activities, activity)
		}
	}

	if activityType == "" || activityType == "drug_order" {
		// 查询药品订单活动
		query := fmt.Sprintf(`
			SELECT 
				'drug_order' as type,
				CONCAT('患者 ', patient_name, ' 在 ', drug_store_name, ' 下单购买 ', drug_name) as description,
				order_id as resource_id,
				'' as doctor_name,
				patient_name,
				'' as organization_id,
				drug_store_name as organization_name,
				'' as department,
				created_at
			FROM drug_orders
			WHERE 1=1 %s
			ORDER BY created_at DESC
		`, dateFilter)

		rows, err := db.DB.Query(query)
		if err != nil {
			return nil, 0, fmt.Errorf("查询药品订单活动失败: %v", err)
		}
		defer rows.Close()

		for rows.Next() {
			var activity map[string]interface{} = make(map[string]interface{})
			var actType, description, resourceID, doctorName, patientName, orgID, orgName, dept string
			var createdAt time.Time

			err := rows.Scan(&actType, &description, &resourceID, &doctorName, &patientName, &orgID, &orgName, &dept, &createdAt)
			if err != nil {
				continue
			}

			activity["type"] = actType
			activity["description"] = description
			activity["resource_id"] = resourceID
			activity["doctor_name"] = doctorName
			activity["patient_name"] = patientName
			activity["organization_id"] = orgID
			activity["organization_name"] = orgName
			activity["department"] = dept
			activity["created_at"] = createdAt.Format("2006-01-02 15:04:05")

			activities = append(activities, activity)
		}
	}

	if activityType == "" || activityType == "insurance_claim" {
		// 查询保险报销活动
		query := fmt.Sprintf(`
			SELECT 
				'insurance_claim' as type,
				CONCAT('患者 ', patient_name, ' 向 ', insurance_name, ' 申请报销') as description,
				claim_id as resource_id,
				'' as doctor_name,
				patient_name,
				'' as organization_id,
				insurance_name as organization_name,
				'' as department,
				created_at
			FROM insurance_claims
			WHERE 1=1 %s
			ORDER BY created_at DESC
		`, dateFilter)

		rows, err := db.DB.Query(query)
		if err != nil {
			return nil, 0, fmt.Errorf("查询保险报销活动失败: %v", err)
		}
		defer rows.Close()

		for rows.Next() {
			var activity map[string]interface{} = make(map[string]interface{})
			var actType, description, resourceID, doctorName, patientName, orgID, orgName, dept string
			var createdAt time.Time

			err := rows.Scan(&actType, &description, &resourceID, &doctorName, &patientName, &orgID, &orgName, &dept, &createdAt)
			if err != nil {
				continue
			}

			activity["type"] = actType
			activity["description"] = description
			activity["resource_id"] = resourceID
			activity["doctor_name"] = doctorName
			activity["patient_name"] = patientName
			activity["organization_id"] = orgID
			activity["organization_name"] = orgName
			activity["department"] = dept
			activity["created_at"] = createdAt.Format("2006-01-02 15:04:05")

			activities = append(activities, activity)
		}
	}

	// 按时间排序（已经在SQL中排序）
	total := len(activities)

	// 分页
	start := (page - 1) * pageSize
	end := start + pageSize
	if start > total {
		return []map[string]interface{}{}, total, nil
	}
	if end > total {
		end = total
	}

	return activities[start:end], total, nil
}

// GetActivityStatistics 获取活动统计数据（从业务表统计）
func GetActivityStatistics() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 统计病历总数
	var prescriptionCount int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM prescriptions").Scan(&prescriptionCount)
	if err != nil {
		prescriptionCount = 0
	}
	stats["totalPrescriptions"] = prescriptionCount

	// 统计补充记录总数
	var supplementCount int
	err = db.DB.QueryRow("SELECT COUNT(*) FROM supplement_records").Scan(&supplementCount)
	if err != nil {
		supplementCount = 0
	}
	stats["totalSupplements"] = supplementCount

	// 统计药品订单总数
	var orderCount int
	err = db.DB.QueryRow("SELECT COUNT(*) FROM drug_orders").Scan(&orderCount)
	if err != nil {
		orderCount = 0
	}
	stats["totalOrders"] = orderCount

	// 统计保险报销总数
	var claimCount int
	err = db.DB.QueryRow("SELECT COUNT(*) FROM insurance_claims").Scan(&claimCount)
	if err != nil {
		claimCount = 0
	}
	stats["totalClaims"] = claimCount

	return stats, nil
}


