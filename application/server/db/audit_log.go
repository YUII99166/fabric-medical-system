package db

import (
	"fmt"
	"time"
)

// AuditLog 审计日志结构
type AuditLog struct {
	ID              int64     `json:"id"`
	UserID          string    `json:"user_id"`
	Username        string    `json:"username"`
	UserRole        string    `json:"user_role"`
	Organization    string    `json:"organization"`
	OrganizationName string   `json:"organization_name"`
	Department      string    `json:"department"`
	Action          string    `json:"action"`
	ResourceType    string    `json:"resource_type"`
	ResourceID      string    `json:"resource_id"`
	RelatedUserID   string    `json:"related_user_id"`
	RelatedUserName string    `json:"related_user_name"`
	Description     string    `json:"description"`
	IPAddress       string    `json:"ip_address"`
	UserAgent       string    `json:"user_agent"`
	TxID            string    `json:"tx_id"`
	Status          string    `json:"status"`
	ErrorMessage    string    `json:"error_message"`
	CreatedAt       time.Time `json:"created_at"`
}

// CreateAuditLog 创建审计日志
func CreateAuditLog(log *AuditLog) error {
	query := `
		INSERT INTO audit_logs (
			user_id, username, user_role, organization, organization_name, department,
			action, resource_type, resource_id, related_user_id, related_user_name,
			description, ip_address, user_agent, tx_id, status, error_message
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := DB.Exec(query,
		log.UserID, log.Username, log.UserRole, log.Organization, log.OrganizationName, log.Department,
		log.Action, log.ResourceType, log.ResourceID, log.RelatedUserID, log.RelatedUserName,
		log.Description, log.IPAddress, log.UserAgent, log.TxID, log.Status, log.ErrorMessage,
	)

	if err != nil {
		return fmt.Errorf("创建审计日志失败: %v", err)
	}

	return nil
}

// GetAuditLogs 获取审计日志列表（支持分页和筛选）
func GetAuditLogs(page, pageSize int, organization, resourceType, startDate, endDate string) ([]AuditLog, int, error) {
	var logs []AuditLog
	var total int

	// 构建查询条件
	whereClause := "WHERE 1=1"
	var args []interface{}

	if organization != "" {
		whereClause += " AND organization = ?"
		args = append(args, organization)
	}

	if resourceType != "" {
		whereClause += " AND resource_type = ?"
		args = append(args, resourceType)
	}

	if startDate != "" && endDate != "" {
		whereClause += " AND created_at BETWEEN ? AND ?"
		args = append(args, startDate+" 00:00:00", endDate+" 23:59:59")
	}

	// 查询总数
	countQuery := "SELECT COUNT(*) FROM audit_logs " + whereClause
	err := DB.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("查询审计日志总数失败: %v", err)
	}

	// 查询数据
	offset := (page - 1) * pageSize
	query := `
		SELECT 
			id, user_id, username, user_role, organization, organization_name, department,
			action, resource_type, resource_id, related_user_id, related_user_name,
			description, ip_address, user_agent, tx_id, status, error_message, created_at
		FROM audit_logs
	` + whereClause + " ORDER BY created_at DESC LIMIT ? OFFSET ?"

	args = append(args, pageSize, offset)
	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("查询审计日志失败: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var log AuditLog
		err := rows.Scan(
			&log.ID, &log.UserID, &log.Username, &log.UserRole, &log.Organization, &log.OrganizationName, &log.Department,
			&log.Action, &log.ResourceType, &log.ResourceID, &log.RelatedUserID, &log.RelatedUserName,
			&log.Description, &log.IPAddress, &log.UserAgent, &log.TxID, &log.Status, &log.ErrorMessage, &log.CreatedAt,
		)
		if err != nil {
			continue
		}
		logs = append(logs, log)
	}

	return logs, total, nil
}

// GetAuditLogStatistics 获取审计日志统计
func GetAuditLogStatistics() (map[string]int, error) {
	stats := make(map[string]int)

	// 统计各类资源的操作次数
	query := `
		SELECT resource_type, COUNT(*) as count
		FROM audit_logs
		WHERE status = 'success'
		GROUP BY resource_type
	`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("查询审计日志统计失败: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var resourceType string
		var count int
		if err := rows.Scan(&resourceType, &count); err != nil {
			continue
		}
		stats[resourceType] = count
	}

	return stats, nil
}

