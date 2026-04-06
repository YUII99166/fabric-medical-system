package api

import (
	"chaincode/model"
	"chaincode/pkg/utils"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// RequestAccess 申请查看病历授权
func RequestAccess(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// args: [prescriptionID, doctorID, reason]
	if len(args) != 3 {
		return shim.Error("参数个数不满足，需要3个参数")
	}
	
	prescriptionID := args[0]
	doctorID := args[1]
	reason := args[2]
	
	if prescriptionID == "" || doctorID == "" || reason == "" {
		return shim.Error("参数存在空值")
	}
	
	// 获取医生信息
	resultsDoctor, err := utils.GetStateByPartialCompositeKeys(stub, model.AccountV2Key, []string{doctorID})
	if err != nil || len(resultsDoctor) != 1 {
		return shim.Error(fmt.Sprintf("医生信息验证失败: %s", err))
	}
	var doctor model.AccountV2
	if err = json.Unmarshal(resultsDoctor[0], &doctor); err != nil {
		return shim.Error(fmt.Sprintf("查询医生信息-反序列化出错: %s", err))
	}
	
	// 验证医生角色
	if doctor.Role != "医生" {
		return shim.Error("操作人权限不足，需要医生角色")
	}
	
	// 获取病历信息
	resultsPrescription, err := utils.GetStateByPartialCompositeKeys2(stub, model.PrescriptionKey, []string{})
	if err != nil {
		return shim.Error(fmt.Sprintf("查询病历失败: %s", err))
	}
	
	var prescription model.Prescription
	found := false
	for _, v := range resultsPrescription {
		var p model.Prescription
		if err = json.Unmarshal(v, &p); err != nil {
			continue
		}
		if p.ID == prescriptionID {
			prescription = p
			found = true
			break
		}
	}
	
	if !found {
		return shim.Error("病历不存在")
	}
	
	// 检查是否已经有权限
	for _, org := range prescription.AuthorizedOrgs {
		if org == doctor.Organization {
			return shim.Error("您已经有权限查看该病历")
		}
	}
	
	// 检查是否已经有待处理的申请
	existingRequests, err := utils.GetStateByPartialCompositeKeys2(stub, model.AccessRequestKey, []string{})
	if err == nil {
		for _, v := range existingRequests {
			var req model.AccessRequest
			if err = json.Unmarshal(v, &req); err != nil {
				continue
			}
			if req.PrescriptionID == prescriptionID && req.DoctorID == doctorID && req.Status == "pending" {
				return shim.Error("您已经提交过申请，请等待患者审批")
			}
		}
	}
	
	// 创建授权请求
	requestID := stub.GetTxID()[:16]
	accessRequest := &model.AccessRequest{
		ID:             requestID,
		PrescriptionID: prescriptionID,
		PatientID:      prescription.Patient,
		PatientName:    prescription.PatientName,
		DoctorID:       doctorID,
		DoctorName:     doctor.AccountName,
		DoctorOrg:      doctor.Organization,
		DoctorOrgName:  doctor.OrganizationName,
		Reason:         reason,
		Status:         "pending",
		RequestTime:    time.Now().Format("2006-01-02 15:04:05"),
		ResponseTime:   "",
		TxID:           stub.GetTxID(),
	}
	
	// 写入账本
	if err := utils.WriteLedger(accessRequest, stub, model.AccessRequestKey, []string{requestID}); err != nil {
		return shim.Error(fmt.Sprintf("写入授权请求失败: %s", err))
	}
	
	// 返回成功信息
	requestByte, err := json.Marshal(accessRequest)
	if err != nil {
		return shim.Error(fmt.Sprintf("序列化授权请求出错: %s", err))
	}
	
	return shim.Success(requestByte)
}

// ApproveAccess 患者审批授权
func ApproveAccess(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// args: [requestID, patientID, approved(true/false), rejectReason]
	if len(args) != 4 {
		return shim.Error("参数个数不满足，需要4个参数")
	}
	
	requestID := args[0]
	patientID := args[1]
	approved := args[2]
	rejectReason := args[3]
	
	if requestID == "" || patientID == "" || approved == "" {
		return shim.Error("参数存在空值")
	}
	
	// 获取患者信息
	resultsPatient, err := utils.GetStateByPartialCompositeKeys(stub, model.AccountV2Key, []string{patientID})
	if err != nil || len(resultsPatient) != 1 {
		return shim.Error(fmt.Sprintf("患者信息验证失败: %s", err))
	}
	var patient model.AccountV2
	if err = json.Unmarshal(resultsPatient[0], &patient); err != nil {
		return shim.Error(fmt.Sprintf("查询患者信息-反序列化出错: %s", err))
	}
	
	// 获取授权请求
	resultsRequest, err := utils.GetStateByPartialCompositeKeys(stub, model.AccessRequestKey, []string{requestID})
	if err != nil || len(resultsRequest) != 1 {
		return shim.Error(fmt.Sprintf("授权请求不存在: %s", err))
	}
	var accessRequest model.AccessRequest
	if err = json.Unmarshal(resultsRequest[0], &accessRequest); err != nil {
		return shim.Error(fmt.Sprintf("查询授权请求-反序列化出错: %s", err))
	}
	
	// 验证患者身份
	if accessRequest.PatientID != patientID {
		return shim.Error("您无权审批此授权请求")
	}
	
	// 检查请求状态
	if accessRequest.Status != "pending" {
		return shim.Error("该授权请求已经处理过了")
	}
	
	// 更新请求状态
	if approved == "true" {
		accessRequest.Status = "approved"
	} else {
		accessRequest.Status = "rejected"
		accessRequest.Reason = rejectReason // 拒绝理由
	}
	accessRequest.ResponseTime = time.Now().Format("2006-01-02 15:04:05")
	
	// 写入更新后的授权请求
	if err := utils.WriteLedger(&accessRequest, stub, model.AccessRequestKey, []string{requestID}); err != nil {
		return shim.Error(fmt.Sprintf("更新授权请求失败: %s", err))
	}
	
	// 如果同意，更新病历的授权列表
	if approved == "true" {
		// 获取病历
		resultsPrescription, err := utils.GetStateByPartialCompositeKeys2(stub, model.PrescriptionKey, []string{})
		if err != nil {
			return shim.Error(fmt.Sprintf("查询病历失败: %s", err))
		}
		
		var prescription model.Prescription
		found := false
		
		for _, v := range resultsPrescription {
			var p model.Prescription
			if err = json.Unmarshal(v, &p); err != nil {
				continue
			}
			if p.ID == accessRequest.PrescriptionID {
				prescription = p
				found = true
				break
			}
		}
		
		if !found {
			return shim.Error("病历不存在")
		}
		
		// 添加授权组织
		alreadyAuthorized := false
		for _, org := range prescription.AuthorizedOrgs {
			if org == accessRequest.DoctorOrg {
				alreadyAuthorized = true
				break
			}
		}
		
		if !alreadyAuthorized {
			prescription.AuthorizedOrgs = append(prescription.AuthorizedOrgs, accessRequest.DoctorOrg)
		}
		
		// 写入更新后的病历
		if err := utils.WriteLedger(&prescription, stub, model.PrescriptionKey, []string{prescription.Patient, prescription.ID}); err != nil {
			return shim.Error(fmt.Sprintf("更新病历授权列表失败: %s", err))
		}
	}
	
	// 返回成功信息
	requestByte, err := json.Marshal(accessRequest)
	if err != nil {
		return shim.Error(fmt.Sprintf("序列化授权请求出错: %s", err))
	}
	
	return shim.Success(requestByte)
}

// QueryAccessRequests 查询授权请求列表
func QueryAccessRequests(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// args: [userID, role] - role可以是"patient"或"doctor"
	if len(args) != 2 {
		return shim.Error("参数个数不满足，需要2个参数")
	}
	
	userID := args[0]
	role := args[1]
	
	if userID == "" || role == "" {
		return shim.Error("参数存在空值")
	}
	
	// 获取所有授权请求
	results, err := utils.GetStateByPartialCompositeKeys2(stub, model.AccessRequestKey, []string{})
	if err != nil {
		return shim.Error(fmt.Sprintf("查询授权请求失败: %s", err))
	}
	
	var requestList []model.AccessRequest
	for _, v := range results {
		if v != nil {
			var req model.AccessRequest
			err := json.Unmarshal(v, &req)
			if err != nil {
				continue
			}
			
			// 根据角色过滤
			if role == "patient" && req.PatientID == userID {
				requestList = append(requestList, req)
			} else if role == "doctor" && req.DoctorID == userID {
				requestList = append(requestList, req)
			}
		}
	}
	
	requestByte, err := json.Marshal(requestList)
	if err != nil {
		return shim.Error(fmt.Sprintf("序列化授权请求列表出错: %s", err))
	}
	
	return shim.Success(requestByte)
}

// QueryPrescriptionsByPatient 根据患者姓名或ID查询病历
func QueryPrescriptionsByPatient(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// args: [searchKey, doctorID] - searchKey可以是患者姓名或ID，doctorID可以为空或"ANY_DOCTOR"
	if len(args) != 2 {
		return shim.Error("参数个数不满足，需要2个参数")
	}
	
	searchKey := args[0]
	doctorID := args[1]
	
	if searchKey == "" {
		return shim.Error("搜索关键词不能为空")
	}
	
	// 如果doctorID为空或特殊值，则使用通用搜索模式
	var doctor model.AccountV2
	var hasValidDoctor bool = false
	
	if doctorID != "" && doctorID != "ANY_DOCTOR" {
		// 获取医生信息
		resultsDoctor, err := utils.GetStateByPartialCompositeKeys(stub, model.AccountV2Key, []string{doctorID})
		if err != nil || len(resultsDoctor) != 1 {
			return shim.Error(fmt.Sprintf("医生信息验证失败: %s", err))
		}
		if err = json.Unmarshal(resultsDoctor[0], &doctor); err != nil {
			return shim.Error(fmt.Sprintf("查询医生信息-反序列化出错: %s", err))
		}
		hasValidDoctor = true
	}
	
	// 获取所有病历
	results, err := utils.GetStateByPartialCompositeKeys2(stub, model.PrescriptionKey, []string{})
	if err != nil {
		return shim.Error(fmt.Sprintf("查询病历失败: %s", err))
	}
	
	type PrescriptionWithAccess struct {
		model.Prescription
		NeedAuthorization bool   `json:"need_authorization"` // 是否需要授权
		CanView           bool   `json:"can_view"`           // 是否可以查看
		IsSameOrg         bool   `json:"is_same_org"`        // 是否同院
	}
	
	var prescriptionList []PrescriptionWithAccess
	for _, v := range results {
		if v != nil {
			var p model.Prescription
			err := json.Unmarshal(v, &p)
			if err != nil {
				continue
			}
			
			// 匹配患者姓名或ID（支持多种匹配方式）
			// 1. 精确匹配患者ID
			// 2. 精确匹配患者姓名
			// 3. 患者姓名包含搜索关键词
			// 4. 搜索关键词包含患者姓名（去掉前缀后）
			var patientNameWithoutPrefix string
			if strings.HasPrefix(p.PatientName, "病人-") {
				patientNameWithoutPrefix = strings.TrimPrefix(p.PatientName, "病人-")
			} else {
				patientNameWithoutPrefix = p.PatientName
			}
			
			var searchKeyWithoutPrefix string
			if strings.HasPrefix(searchKey, "病人-") {
				searchKeyWithoutPrefix = strings.TrimPrefix(searchKey, "病人-")
			} else {
				searchKeyWithoutPrefix = searchKey
			}
			
			isMatch := p.Patient == searchKey || // 精确匹配患者ID
				p.PatientName == searchKey || // 精确匹配患者姓名
				strings.Contains(p.PatientName, searchKey) || // 患者姓名包含搜索关键词
				strings.Contains(patientNameWithoutPrefix, searchKeyWithoutPrefix) || // 去掉前缀后匹配
				patientNameWithoutPrefix == searchKeyWithoutPrefix // 去掉前缀后精确匹配
			
			if isMatch {
				// 检查是否有权限查看
				var isSameOrg, hasAuth, needAuth, canView bool
				
				if hasValidDoctor {
					isSameOrg = p.OrganizationID == doctor.Organization
					for _, org := range p.AuthorizedOrgs {
						if org == doctor.Organization {
							hasAuth = true
							break
						}
					}
					needAuth = !isSameOrg && !hasAuth
					canView = isSameOrg || hasAuth
				} else {
					// 如果没有有效的医生信息，默认可以查看但需要授权
					isSameOrg = false
					hasAuth = false
					needAuth = true
					canView = true // 允许查看，但标记需要授权
				}
				
				prescriptionList = append(prescriptionList, PrescriptionWithAccess{
					Prescription:      p,
					NeedAuthorization: needAuth,
					CanView:           canView,
					IsSameOrg:         isSameOrg,
				})
			}
		}
	}
	
	prescriptionByte, err := json.Marshal(prescriptionList)
	if err != nil {
		return shim.Error(fmt.Sprintf("序列化病历列表出错: %s", err))
	}
	
	return shim.Success(prescriptionByte)
}
