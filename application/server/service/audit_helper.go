package service

import (
	"application/db"
	"fmt"

	"github.com/gin-gonic/gin"
)

// LogPrescriptionCreation 记录病历创建日志
func LogPrescriptionCreation(c *gin.Context, userID, username, userRole, organization, organizationName, department, prescriptionID, patientID, patientName, txID string) error {
	log := &db.AuditLog{
		UserID:           userID,
		Username:         username,
		UserRole:         userRole,
		Organization:     organization,
		OrganizationName: organizationName,
		Department:       department,
		Action:           "create",
		ResourceType:     "prescription",
		ResourceID:       prescriptionID,
		RelatedUserID:    patientID,
		RelatedUserName:  patientName,
		Description:      fmt.Sprintf("医生 %s 为患者 %s 创建了病历", username, patientName),
		IPAddress:        c.ClientIP(),
		UserAgent:        c.Request.UserAgent(),
		TxID:             txID,
		Status:           "success",
	}

	return db.CreateAuditLog(log)
}

// LogSupplementCreation 记录补充记录创建日志
func LogSupplementCreation(c *gin.Context, userID, username, userRole, organization, organizationName, department, recordID, recordType, prescriptionID, txID string) error {
	log := &db.AuditLog{
		UserID:           userID,
		Username:         username,
		UserRole:         userRole,
		Organization:     organization,
		OrganizationName: organizationName,
		Department:       department,
		Action:           "create",
		ResourceType:     "supplement",
		ResourceID:       recordID,
		RelatedUserID:    prescriptionID,
		RelatedUserName:  "",
		Description:      fmt.Sprintf("医生 %s 添加了%s记录", username, recordType),
		IPAddress:        c.ClientIP(),
		UserAgent:        c.Request.UserAgent(),
		TxID:             txID,
		Status:           "success",
	}

	return db.CreateAuditLog(log)
}

// LogAccessRequest 记录授权请求日志
func LogAccessRequest(c *gin.Context, userID, username, userRole, organization, organizationName, requestID, patientID, patientName, prescriptionID, txID string) error {
	log := &db.AuditLog{
		UserID:           userID,
		Username:         username,
		UserRole:         userRole,
		Organization:     organization,
		OrganizationName: organizationName,
		Department:       "",
		Action:           "request",
		ResourceType:     "access_request",
		ResourceID:       requestID,
		RelatedUserID:    patientID,
		RelatedUserName:  patientName,
		Description:      fmt.Sprintf("医生 %s 请求访问患者 %s 的病历", username, patientName),
		IPAddress:        c.ClientIP(),
		UserAgent:        c.Request.UserAgent(),
		TxID:             txID,
		Status:           "success",
	}

	return db.CreateAuditLog(log)
}

// LogAccessApproval 记录授权批准日志
func LogAccessApproval(c *gin.Context, userID, username, userRole, requestID, doctorName string, approved bool, txID string) error {
	action := "approve"
	description := fmt.Sprintf("患者 %s 批准了医生 %s 的访问请求", username, doctorName)
	if !approved {
		action = "reject"
		description = fmt.Sprintf("患者 %s 拒绝了医生 %s 的访问请求", username, doctorName)
	}

	log := &db.AuditLog{
		UserID:           userID,
		Username:         username,
		UserRole:         userRole,
		Organization:     "",
		OrganizationName: "",
		Department:       "",
		Action:           action,
		ResourceType:     "access_request",
		ResourceID:       requestID,
		RelatedUserID:    "",
		RelatedUserName:  doctorName,
		Description:      description,
		IPAddress:        c.ClientIP(),
		UserAgent:        c.Request.UserAgent(),
		TxID:             txID,
		Status:           "success",
	}

	return db.CreateAuditLog(log)
}

// LogDrugOrderCreation 记录药品订单创建日志
func LogDrugOrderCreation(c *gin.Context, userID, username, userRole, orderID, drugName, patientID, patientName, drugStoreName, txID string) error {
	log := &db.AuditLog{
		UserID:           userID,
		Username:         username,
		UserRole:         userRole,
		Organization:     "",
		OrganizationName: drugStoreName,
		Department:       "",
		Action:           "create",
		ResourceType:     "drug_order",
		ResourceID:       orderID,
		RelatedUserID:    patientID,
		RelatedUserName:  patientName,
		Description:      fmt.Sprintf("患者 %s 在 %s 下单购买 %s", patientName, drugStoreName, drugName),
		IPAddress:        c.ClientIP(),
		UserAgent:        c.Request.UserAgent(),
		TxID:             txID,
		Status:           "success",
	}

	return db.CreateAuditLog(log)
}

// LogInsuranceClaimCreation 记录保险报销创建日志
func LogInsuranceClaimCreation(c *gin.Context, userID, username, userRole, claimID, patientID, patientName, insuranceName, txID string) error {
	log := &db.AuditLog{
		UserID:           userID,
		Username:         username,
		UserRole:         userRole,
		Organization:     "",
		OrganizationName: insuranceName,
		Department:       "",
		Action:           "create",
		ResourceType:     "insurance_claim",
		ResourceID:       claimID,
		RelatedUserID:    patientID,
		RelatedUserName:  patientName,
		Description:      fmt.Sprintf("患者 %s 向 %s 申请报销", patientName, insuranceName),
		IPAddress:        c.ClientIP(),
		UserAgent:        c.Request.UserAgent(),
		TxID:             txID,
		Status:           "success",
	}

	return db.CreateAuditLog(log)
}

// LogUserCreation 记录用户创建日志
func LogUserCreation(c *gin.Context, adminID, adminName, adminRole, newUserID, newUsername, newUserRole, organization, organizationName string) error {
	log := &db.AuditLog{
		UserID:           adminID,
		Username:         adminName,
		UserRole:         adminRole,
		Organization:     organization,
		OrganizationName: organizationName,
		Department:       "",
		Action:           "create",
		ResourceType:     "user",
		ResourceID:       newUserID,
		RelatedUserID:    newUserID,
		RelatedUserName:  newUsername,
		Description:      fmt.Sprintf("管理员 %s 创建了新用户 %s（角色：%s）", adminName, newUsername, newUserRole),
		IPAddress:        c.ClientIP(),
		UserAgent:        c.Request.UserAgent(),
		TxID:             "",
		Status:           "success",
	}

	return db.CreateAuditLog(log)
}

// LogUserDeletion 记录用户删除日志
func LogUserDeletion(c *gin.Context, adminID, adminName, adminRole, targetUserID, targetUsername, targetUserRole string) error {
	log := &db.AuditLog{
		UserID:           adminID,
		Username:         adminName,
		UserRole:         adminRole,
		Organization:     "",
		OrganizationName: "",
		Department:       "",
		Action:           "delete",
		ResourceType:     "user",
		ResourceID:       targetUserID,
		RelatedUserID:    targetUserID,
		RelatedUserName:  targetUsername,
		Description:      fmt.Sprintf("管理员 %s 停用了用户 %s（角色：%s）", adminName, targetUsername, targetUserRole),
		IPAddress:        c.ClientIP(),
		UserAgent:        c.Request.UserAgent(),
		TxID:             "",
		Status:           "success",
	}

	return db.CreateAuditLog(log)
}

// LogUserRestoration 记录用户恢复日志
func LogUserRestoration(c *gin.Context, adminID, adminName, adminRole, targetUserID, targetUsername, targetUserRole string) error {
	log := &db.AuditLog{
		UserID:           adminID,
		Username:         adminName,
		UserRole:         adminRole,
		Organization:     "",
		OrganizationName: "",
		Department:       "",
		Action:           "restore",
		ResourceType:     "user",
		ResourceID:       targetUserID,
		RelatedUserID:    targetUserID,
		RelatedUserName:  targetUsername,
		Description:      fmt.Sprintf("管理员 %s 恢复了用户 %s（角色：%s）", adminName, targetUsername, targetUserRole),
		IPAddress:        c.ClientIP(),
		UserAgent:        c.Request.UserAgent(),
		TxID:             "",
		Status:           "success",
	}

	return db.CreateAuditLog(log)
}

// LogError 记录错误日志
func LogError(c *gin.Context, userID, username, userRole, action, resourceType, resourceID, errorMessage string) error {
	log := &db.AuditLog{
		UserID:           userID,
		Username:         username,
		UserRole:         userRole,
		Organization:     "",
		OrganizationName: "",
		Department:       "",
		Action:           action,
		ResourceType:     resourceType,
		ResourceID:       resourceID,
		RelatedUserID:    "",
		RelatedUserName:  "",
		Description:      fmt.Sprintf("操作失败: %s", errorMessage),
		IPAddress:        c.ClientIP(),
		UserAgent:        c.Request.UserAgent(),
		TxID:             "",
		Status:           "failed",
		ErrorMessage:     errorMessage,
	}

	return db.CreateAuditLog(log)
}
