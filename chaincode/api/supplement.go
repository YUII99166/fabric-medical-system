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

// AddSupplementRecord 添加补充诊疗记录
// args: [0]originalPrescriptionID [1]doctorID [2]recordType [3]chiefComplaint [4]presentIllness 
//       [5]physicalExam [6]diagnosis [7]treatment [8]drugName [9]drugAmount [10]medicalAdvice [11]comment
func AddSupplementRecord(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 验证参数
	if len(args) != 12 {
		return shim.Error(fmt.Sprintf("参数个数不满足，需要12个参数，实际收到%d个", len(args)))
	}

	originalPrescriptionID := args[0]
	doctorID := args[1]
	recordType := args[2]
	chiefComplaint := args[3]
	presentIllness := args[4]
	physicalExam := args[5]
	diagnosis := args[6]
	treatment := args[7]
	drugName := args[8]
	drugAmount := args[9]
	medicalAdvice := args[10]
	comment := args[11]

	// 验证必填字段
	if originalPrescriptionID == "" || doctorID == "" || diagnosis == "" {
		return shim.Error("参数存在空值（原始病历ID、医生ID、诊断不能为空）")
	}

	// 验证记录类型
	validTypes := map[string]bool{
		"consultation": true, // 会诊
		"followup":     true, // 复诊
		"emergency":    true, // 急诊
		"transfer":     true, // 转院
	}
	if !validTypes[recordType] {
		recordType = "followup" // 默认为复诊
	}

	// 获取原始病历信息
	prescriptionResults, err := utils.GetStateByPartialCompositeKeys2(stub, model.PrescriptionKey, []string{})
	if err != nil {
		return shim.Error(fmt.Sprintf("查询原始病历失败: %s", err))
	}

	var originalPrescription model.Prescription
	found := false
	for _, v := range prescriptionResults {
		var p model.Prescription
		if err := json.Unmarshal(v, &p); err != nil {
			continue
		}
		if p.ID == originalPrescriptionID {
			originalPrescription = p
			found = true
			break
		}
	}

	if !found {
		return shim.Error(fmt.Sprintf("未找到原始病历，ID: %s", originalPrescriptionID))
	}

	// 获取医生信息
	doctorResults, err := utils.GetStateByPartialCompositeKeys(stub, model.AccountV2Key, []string{doctorID})
	if err != nil || len(doctorResults) != 1 {
		return shim.Error(fmt.Sprintf("医生信息验证失败: %s", err))
	}
	var doctor model.AccountV2
	if err = json.Unmarshal(doctorResults[0], &doctor); err != nil {
		return shim.Error(fmt.Sprintf("查询医生信息-反序列化出错: %s", err))
	}

	// 验证医生角色
	if doctor.Role != "医生" {
		return shim.Error("操作人权限不足，需要医生角色")
	}

	// 验证医生是否有权限查看该病历（通过授权）
	// 检查是否是本院病历或已授权
	hasPermission := false
	if doctor.Organization == originalPrescription.OrganizationID {
		// 本院医生可以添加补充记录
		hasPermission = true
	} else {
		// 检查是否已授权
		for _, org := range originalPrescription.AuthorizedOrgs {
			if org == doctor.Organization {
				hasPermission = true
				break
			}
		}
	}

	if !hasPermission {
		return shim.Error("无权限：该病历未授权给您所在的医院")
	}

	// 处理药品信息
	var drugs []model.Drug
	if drugName != "" && drugAmount != "" {
		drugNames := strings.Split(drugName, ",")
		drugAmounts := strings.Split(drugAmount, ",")
		for i, v := range drugNames {
			if i < len(drugAmounts) {
				drug := model.Drug{
					Name:   v,
					Amount: drugAmounts[i],
				}
				drugs = append(drugs, drug)
			}
		}
	}

	// 创建补充记录
	supplementRecord := &model.SupplementRecord{
		ID:                     stub.GetTxID()[:16],
		OriginalPrescriptionID: originalPrescriptionID,
		RecordType:             recordType,
		PatientID:              originalPrescription.Patient,
		PatientName:            originalPrescription.PatientName,
		ChiefComplaint:         chiefComplaint,
		PresentIllness:         presentIllness,
		PhysicalExam:           physicalExam,
		Diagnosis:              diagnosis,
		Treatment:              treatment,
		Drug:                   drugs,
		MedicalAdvice:          medicalAdvice,
		DoctorID:               doctorID,
		DoctorName:             doctor.AccountName,
		DoctorTitle:            doctor.DoctorTitle,
		Department:             doctor.Department,
		HospitalName:           doctor.OrganizationName,
		OrganizationID:         doctor.Organization,
		OrganizationName:       doctor.OrganizationName,
		Created:                time.Now().Format("2006-01-02 15:04:05"),
		TxID:                   stub.GetTxID(),
		CreatorMSPID:           doctor.Organization,
		IsReadOnly:             true, // 创建后不可修改
		Comment:                comment,
	}

	// 写入账本
	if err := utils.WriteLedger(supplementRecord, stub, model.SupplementRecordKey, []string{originalPrescriptionID, supplementRecord.ID}); err != nil {
		return shim.Error(fmt.Sprintf("写入账本失败: %s", err))
	}

	// 返回成功创建的信息
	recordByte, err := json.Marshal(supplementRecord)
	if err != nil {
		return shim.Error(fmt.Sprintf("序列化补充记录出错: %s", err))
	}

	return shim.Success(recordByte)
}

// QuerySupplementRecords 查询补充记录
// args: [0]originalPrescriptionID (可选，如果为空则查询所有)
func QuerySupplementRecords(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var recordList []model.SupplementRecord
	
	results, err := utils.GetStateByPartialCompositeKeys2(stub, model.SupplementRecordKey, args)
	if err != nil {
		return shim.Error(fmt.Sprintf("查询补充记录失败: %s", err))
	}

	for _, v := range results {
		if v != nil {
			var record model.SupplementRecord
			err := json.Unmarshal(v, &record)
			if err != nil {
				return shim.Error(fmt.Sprintf("QuerySupplementRecords-反序列化出错: %s", err))
			}
			recordList = append(recordList, record)
		}
	}

	recordByte, err := json.Marshal(recordList)
	if err != nil {
		return shim.Error(fmt.Sprintf("QuerySupplementRecords-序列化出错: %s", err))
	}

	return shim.Success(recordByte)
}

// QueryFullMedicalHistory 查询完整病历历史（原始病历+所有补充记录）
// args: [0]prescriptionID
func QueryFullMedicalHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("参数个数不满足，需要1个参数（病历ID）")
	}

	prescriptionID := args[0]
	if prescriptionID == "" {
		return shim.Error("病历ID不能为空")
	}

	// 查询原始病历
	prescriptionResults, err := utils.GetStateByPartialCompositeKeys2(stub, model.PrescriptionKey, []string{})
	if err != nil {
		return shim.Error(fmt.Sprintf("查询原始病历失败: %s", err))
	}

	var originalPrescription model.Prescription
	found := false
	for _, v := range prescriptionResults {
		var p model.Prescription
		if err := json.Unmarshal(v, &p); err != nil {
			continue
		}
		if p.ID == prescriptionID {
			originalPrescription = p
			found = true
			break
		}
	}

	if !found {
		return shim.Error(fmt.Sprintf("未找到病历，ID: %s", prescriptionID))
	}

	// 查询该病历的所有补充记录
	supplementResults, err := utils.GetStateByPartialCompositeKeys2(stub, model.SupplementRecordKey, []string{prescriptionID})
	if err != nil {
		return shim.Error(fmt.Sprintf("查询补充记录失败: %s", err))
	}

	var supplementRecords []model.SupplementRecord
	for _, v := range supplementResults {
		if v != nil {
			var record model.SupplementRecord
			if err := json.Unmarshal(v, &record); err != nil {
				continue
			}
			supplementRecords = append(supplementRecords, record)
		}
	}

	// 构建完整病历历史
	history := model.MedicalHistory{
		OriginalPrescription: originalPrescription,
		SupplementRecords:    supplementRecords,
		TotalRecords:         len(supplementRecords),
	}

	historyByte, err := json.Marshal(history)
	if err != nil {
		return shim.Error(fmt.Sprintf("序列化病历历史出错: %s", err))
	}

	return shim.Success(historyByte)
}
