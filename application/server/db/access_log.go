package db

import (
	"fmt"
	"time"
)

// PrescriptionAccessLog 病历访问日志结构
type PrescriptionAccessLog struct {
	ID                       int64     `json:"id"`
	LogID                    string    `json:"log_id"`
	PrescriptionID           string    `json:"prescription_id"`
	PrescriptionNo           string    `json:"prescription_no"`
	PatientID                string    `json:"patient_id"`
	PatientName              string    `json:"patient_name"`
	AccessorID               string    `json:"accessor_id"`
	AccessorName             string    `json:"accessor_name"`
	AccessorRole             string    `json:"accessor_role"`
	AccessorOrganization     string    `json:"accessor_organization"`
	AccessorOrganizationName string    `json:"accessor_organization_name"`
	AccessType               string    `json:"access_type"`
	AccessReason             string    `json:"access_reason"`
	IPAddress                string    `json:"ip_address"`
	UserAgent                string    `json:"user_agent"`
	TxID                     string    `json:"tx_id"`
	AccessedAt               time.Time `json:"accessed_at"`
}

// CreateAccessLog 创建访问日志
func CreateAccessLog(log *PrescriptionAccessLog) error {
	query := `
		INSERT INTO prescription_access_logs (
			log_id, prescription_id, prescription_no, patient_id, patient_name,
			accessor_id, accessor_name, accessor_role, accessor_organization, accessor_organization_name,
			access_type, access_reason, ip_address, user_agent, tx_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := DB.Exec(query,
		log.LogID, log.PrescriptionID, log.PrescriptionNo, log.PatientID, log.PatientName,
		log.AccessorID, log.AccessorName, log.AccessorRole, log.AccessorOrganization, log.AccessorOrganizationName,
		log.AccessType, log.AccessReason, log.IPAddress, log.UserAgent, log.TxID,
	)

	if err != nil {
		return fmt.Errorf("创建访问日志失败: %v", err)
	}

	return nil
}

// GetAccessLogsByPatient 获取患者的病历访问日志
func GetAccessLogsByPatient(patientID string, page, pageSize int) ([]PrescriptionAccessLog, int, error) {
	var logs []PrescriptionAccessLog
	var total int

	// 查询总数
	countQuery := "SELECT COUNT(*) FROM prescription_access_logs WHERE patient_id = ?"
	err := DB.QueryRow(countQuery, patientID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("查询访问日志总数失败: %v", err)
	}

	// 查询数据
	offset := (page - 1) * pageSize
	query := `
		SELECT 
			id, log_id, prescription_id, prescription_no, patient_id, patient_name,
			accessor_id, accessor_name, accessor_role, accessor_organization, accessor_organization_name,
			access_type, access_reason, ip_address, user_agent, tx_id, accessed_at
		FROM prescription_access_logs
		WHERE patient_id = ?
		ORDER BY accessed_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := DB.Query(query, patientID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("查询访问日志失败: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var log PrescriptionAccessLog
		err := rows.Scan(
			&log.ID, &log.LogID, &log.PrescriptionID, &log.PrescriptionNo, &log.PatientID, &log.PatientName,
			&log.AccessorID, &log.AccessorName, &log.AccessorRole, &log.AccessorOrganization, &log.AccessorOrganizationName,
			&log.AccessType, &log.AccessReason, &log.IPAddress, &log.UserAgent, &log.TxID, &log.AccessedAt,
		)
		if err != nil {
			continue
		}
		logs = append(logs, log)
	}

	return logs, total, nil
}

// GetAccessLogsByPrescription 获取特定病历的访问日志
func GetAccessLogsByPrescription(prescriptionID string) ([]PrescriptionAccessLog, error) {
	var logs []PrescriptionAccessLog

	query := `
		SELECT 
			id, log_id, prescription_id, prescription_no, patient_id, patient_name,
			accessor_id, accessor_name, accessor_role, accessor_organization, accessor_organization_name,
			access_type, access_reason, ip_address, user_agent, tx_id, accessed_at
		FROM prescription_access_logs
		WHERE prescription_id = ?
		ORDER BY accessed_at DESC
	`

	rows, err := DB.Query(query, prescriptionID)
	if err != nil {
		return nil, fmt.Errorf("查询访问日志失败: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var log PrescriptionAccessLog
		err := rows.Scan(
			&log.ID, &log.LogID, &log.PrescriptionID, &log.PrescriptionNo, &log.PatientID, &log.PatientName,
			&log.AccessorID, &log.AccessorName, &log.AccessorRole, &log.AccessorOrganization, &log.AccessorOrganizationName,
			&log.AccessType, &log.AccessReason, &log.IPAddress, &log.UserAgent, &log.TxID, &log.AccessedAt,
		)
		if err != nil {
			continue
		}
		logs = append(logs, log)
	}

	return logs, nil
}

// GetAccessStatistics 获取访问统计
func GetAccessStatistics(patientID string) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 统计总访问次数
	var totalAccess int
	err := DB.QueryRow("SELECT COUNT(*) FROM prescription_access_logs WHERE patient_id = ?", patientID).Scan(&totalAccess)
	if err != nil {
		totalAccess = 0
	}
	stats["total_access"] = totalAccess

	// 统计不同访问者数量
	var uniqueAccessors int
	err = DB.QueryRow("SELECT COUNT(DISTINCT accessor_id) FROM prescription_access_logs WHERE patient_id = ? AND accessor_id != patient_id", patientID).Scan(&uniqueAccessors)
	if err != nil {
		uniqueAccessors = 0
	}
	stats["unique_accessors"] = uniqueAccessors

	// 统计最近7天的访问次数
	var recentAccess int
	err = DB.QueryRow("SELECT COUNT(*) FROM prescription_access_logs WHERE patient_id = ? AND accessed_at >= DATE_SUB(NOW(), INTERVAL 7 DAY)", patientID).Scan(&recentAccess)
	if err != nil {
		recentAccess = 0
	}
	stats["recent_access"] = recentAccess

	return stats, nil
}
