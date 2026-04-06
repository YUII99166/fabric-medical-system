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

// CreatePrescription 创建处方(医生)
func CreatePrescription(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 验证参数
	if len(args) != 11 {
		return shim.Error("参数个数不满足")
	}
	doctorID := args[0]        // 医生id
	patientID := args[1]       // 患者id
	chiefComplaint := args[2]  // 主诉
	presentIllness := args[3]  // 现病史
	physicalExam := args[4]    // 体格检查
	diagnosis := args[5]       // 诊断结果
	drugName := args[6]        // 药品名
	drugAmount := args[7]      // 药品数量
	medicalAdvice := args[8]   // 医嘱
	hospitalID := args[9]      // 医院ID
	comment := args[10]        // 备注
	if doctorID == "" || patientID == "" || diagnosis == "" || drugName == "" || drugAmount == "" || hospitalID == "" {
		return shim.Error("参数存在空值")
	}

	// 参数数据格式转换
	var drugs []model.Drug
	drugNames := strings.Split(drugName, ",")
	drugAmounts := strings.Split(drugAmount, ",")
	for i, v := range drugNames {
		drug := model.Drug{
			Name:   v,
			Amount: drugAmounts[i],
		}
		drugs = append(drugs, drug)
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

	// 获取患者信息
	resultsPatient, err := utils.GetStateByPartialCompositeKeys(stub, model.AccountV2Key, []string{patientID})
	if err != nil || len(resultsPatient) != 1 {
		return shim.Error(fmt.Sprintf("患者信息验证失败: %s", err))
	}
	var patient model.AccountV2
	if err = json.Unmarshal(resultsPatient[0], &patient); err != nil {
		return shim.Error(fmt.Sprintf("查询患者信息-反序列化出错: %s", err))
	}

	// 生成病历编号
	prescriptionNo := fmt.Sprintf("%s-%s-%s", 
		doctor.Organization, 
		time.Now().Format("20060102"), 
		stub.GetTxID()[:6])

	prescription := &model.Prescription{
		ID:               stub.GetTxID()[:16],
		PrescriptionNo:   prescriptionNo,
		Patient:          patientID,
		PatientName:      patient.AccountName,
		ChiefComplaint:   chiefComplaint,
		PresentIllness:   presentIllness,
		PhysicalExam:     physicalExam,
		Diagnosis:        diagnosis,
		Drug:             drugs,
		MedicalAdvice:    medicalAdvice,
		Doctor:           doctorID,
		DoctorName:       doctor.AccountName,
		DoctorTitle:      doctor.DoctorTitle,
		Hospital:         hospitalID,
		HospitalName:     doctor.OrganizationName,
		OrganizationID:   doctor.Organization,
		OrganizationName: doctor.OrganizationName,
		Department:       doctor.Department,
		Created:          time.Now().Format("2006-01-02 15:04:05"),
		Comment:          comment,
		TxID:             stub.GetTxID(),
		CreatorMSPID:     doctor.Organization,
		AuthorizedOrgs:   []string{doctor.Organization}, // 默认只有本院可见
	}
	
	// 写入账本
	if err := utils.WriteLedger(prescription, stub, model.PrescriptionKey, []string{prescription.Patient, prescription.ID}); err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	
	// 将成功创建的信息返回
	prescriptionByte, err := json.Marshal(prescription)
	if err != nil {
		return shim.Error(fmt.Sprintf("序列化成功创建的信息出错: %s", err))
	}
	// 成功返回
	return shim.Success(prescriptionByte)
}

// QueryPrescription 查询处方(可查询所有，也可根据所有人查询名下处方)
func QueryPrescription(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var prescriptionList []model.Prescription
	results, err := utils.GetStateByPartialCompositeKeys2(stub, model.PrescriptionKey, args)
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	for _, v := range results {
		if v != nil {
			var p model.Prescription
			err := json.Unmarshal(v, &p)
			if err != nil {
				return shim.Error(fmt.Sprintf("QueryPrescription-反序列化出错: %s", err))
			}
			prescriptionList = append(prescriptionList, p)
		}
	}
	prescriptionByte, err := json.Marshal(prescriptionList)
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryPrescription-序列化出错: %s", err))
	}
	return shim.Success(prescriptionByte)
}

// QueryPatient 查询患者
func QueryPatient(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var patientList []model.Patient
	results, err := utils.GetStateByPartialCompositeKeys2(stub, model.PatientKey, args)
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	for _, v := range results {
		if v != nil {
			var p model.Patient
			err := json.Unmarshal(v, &p)
			if err != nil {
				return shim.Error(fmt.Sprintf("QueryPatient-反序列化出错: %s", err))
			}
			patientList = append(patientList, p)
		}
	}
	patientByte, err := json.Marshal(patientList)
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryPatient-序列化出错: %s", err))
	}
	return shim.Success(patientByte)
}
