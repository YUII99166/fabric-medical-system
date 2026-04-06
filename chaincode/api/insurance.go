package api

import (
	"chaincode/model"
	"chaincode/pkg/utils"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"time"
)

// CreateInsuranceCover 创建保险报销订单
func CreateInsuranceCover(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 验证参数
	if len(args) != 2 {
		return shim.Error("参数个数不满足")
	}
	prescriptionID := args[0] // 处方id
	patientID := args[1]      // 患者id

	if prescriptionID == "" || patientID == "" {
		return shim.Error("参数存在空值")
	}

	// 判断患者是否存在
	resultsAccount, err := utils.GetStateByPartialCompositeKeys(stub, model.AccountV2Key, []string{patientID})
	if err != nil || len(resultsAccount) != 1 {
		return shim.Error(fmt.Sprintf("患者权限验证失败%s", err))
	}

	record := &model.InsuranceCover{
		ID:           stub.GetTxID()[:16],
		Prescription: prescriptionID,
		Patient:      patientID,
		Status:       "processing", // 默认状态为处理中
		Created:      time.Now().Format("2006-01-02 15:04:05"),
	}

	// 写入账本
	if err := utils.WriteLedger(record, stub, model.InsuranceKey, []string{record.Patient, record.ID}); err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	//将成功创建的信息返回
	recordByte, err := json.Marshal(record)
	if err != nil {
		return shim.Error(fmt.Sprintf("序列化成功创建的信息出错: %s", err))
	}
	// 成功返回
	return shim.Success(recordByte)
}

// QueryInsuranceCover 查询保险报销订单(可查询所有，也可根据患者查询)
func QueryInsuranceCover(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var insuranceList []model.InsuranceCover
	results, err := utils.GetStateByPartialCompositeKeys2(stub, model.InsuranceKey, args)
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	for _, v := range results {
		if v != nil {
			var p model.InsuranceCover
			err := json.Unmarshal(v, &p)
			if err != nil {
				return shim.Error(fmt.Sprintf("QueryInsuranceCover-反序列化出错: %s", err))
			}
			insuranceList = append(insuranceList, p)
		}
	}
	insuranceByte, err := json.Marshal(insuranceList)
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryInsuranceCover-序列化出错: %s", err))
	}
	return shim.Success(insuranceByte)
}

// UpdateInsuranceCover 更新保险报销订单状态
func UpdateInsuranceCover(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 验证参数
	if len(args) != 3 {
		return shim.Error("参数个数不满足")
	}
	patientID := args[0]   // 患者id
	insuranceID := args[1] // 保险订单id
	status := args[2]      // 新状态

	if patientID == "" || insuranceID == "" || status == "" {
		return shim.Error("参数存在空值")
	}

	// 验证状态值是否合法
	validStatuses := map[string]bool{
		"processing": true,
		"cancelled":  true,
		"refused":    true,
		"approved":   true,
	}
	if !validStatuses[status] {
		return shim.Error("状态值不合法")
	}

	// 查询保险订单
	results, err := utils.GetStateByPartialCompositeKeys(stub, model.InsuranceKey, []string{patientID, insuranceID})
	if err != nil || len(results) != 1 {
		return shim.Error(fmt.Sprintf("保险订单查询失败: %s", err))
	}

	var insurance model.InsuranceCover
	if err = json.Unmarshal(results[0], &insurance); err != nil {
		return shim.Error(fmt.Sprintf("反序列化保险订单失败: %s", err))
	}

	// 更新状态
	insurance.Status = status

	// 写入账本
	if err := utils.WriteLedger(&insurance, stub, model.InsuranceKey, []string{insurance.Patient, insurance.ID}); err != nil {
		return shim.Error(fmt.Sprintf("更新保险订单失败: %s", err))
	}

	// 返回更新后的信息
	insuranceByte, err := json.Marshal(insurance)
	if err != nil {
		return shim.Error(fmt.Sprintf("序列化保险订单出错: %s", err))
	}
	return shim.Success(insuranceByte)
}

// DeleteInsuranceCover 删除保险报销订单
func DeleteInsuranceCover(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 验证参数
	if len(args) != 2 {
		return shim.Error("参数个数不满足")
	}
	patientID := args[0]   // 患者id
	insuranceID := args[1] // 保险订单id

	if patientID == "" || insuranceID == "" {
		return shim.Error("参数存在空值")
	}

	// 查询保险订单是否存在
	results, err := utils.GetStateByPartialCompositeKeys(stub, model.InsuranceKey, []string{patientID, insuranceID})
	if err != nil || len(results) != 1 {
		return shim.Error(fmt.Sprintf("保险订单查询失败: %s", err))
	}

	// 删除保险订单
	insuranceKey, err := stub.CreateCompositeKey(model.InsuranceKey, []string{patientID, insuranceID})
	if err != nil {
		return shim.Error(fmt.Sprintf("创建复合键失败: %s", err))
	}

	if err := stub.DelState(insuranceKey); err != nil {
		return shim.Error(fmt.Sprintf("删除保险订单失败: %s", err))
	}

	return shim.Success([]byte("删除成功"))
}
