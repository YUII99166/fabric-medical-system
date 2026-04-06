package api

import (
	"chaincode/model"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// RecordPrescriptionAccess 记录病历访问日志
func RecordPrescriptionAccess(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// args: [logID, prescriptionID, prescriptionNo, patientID, patientName, accessorID, accessorName, accessorRole, accessorOrg, accessorOrgName, accessType, accessReason]
	if len(args) != 12 {
		return shim.Error(fmt.Sprintf("参数个数不满足，需要12个参数，实际收到%d个", len(args)))
	}

	logID := args[0]
	prescriptionID := args[1]
	prescriptionNo := args[2]
	patientID := args[3]
	patientName := args[4]
	accessorID := args[5]
	accessorName := args[6]
	accessorRole := args[7]
	accessorOrg := args[8]
	accessorOrgName := args[9]
	accessType := args[10]
	accessReason := args[11]

	if logID == "" || prescriptionID == "" || patientID == "" || accessorID == "" {
		return shim.Error("关键参数存在空值")
	}

	// 创建访问日志对象
	accessLog := &model.PrescriptionAccessLog{
		LogID:                    logID,
		PrescriptionID:           prescriptionID,
		PrescriptionNo:           prescriptionNo,
		PatientID:                patientID,
		PatientName:              patientName,
		AccessorID:               accessorID,
		AccessorName:             accessorName,
		AccessorRole:             accessorRole,
		AccessorOrganization:     accessorOrg,
		AccessorOrganizationName: accessorOrgName,
		AccessType:               accessType,
		AccessReason:             accessReason,
		AccessedAt:               time.Now().Format("2006-01-02 15:04:05"),
	}

	// 序列化并写入账本
	accessLogBytes, err := json.Marshal(accessLog)
	if err != nil {
		return shim.Error(fmt.Sprintf("序列化失败: %s", err))
	}

	// 创建复合主键
	key, err := stub.CreateCompositeKey(model.AccessLogKey, []string{logID})
	if err != nil {
		return shim.Error(fmt.Sprintf("创建主键失败: %s", err))
	}

	// 写入账本
	if err := stub.PutState(key, accessLogBytes); err != nil {
		return shim.Error(fmt.Sprintf("写入账本失败: %s", err))
	}

	// 获取交易ID
	txID := stub.GetTxID()

	// 返回成功响应，包含交易ID
	return shim.Success([]byte(fmt.Sprintf(`{"message":"访问日志已记录到区块链","logId":"%s","txId":"%s"}`, logID, txID)))
}

// QueryAccessLogsByPatient 查询患者的访问日志
func QueryAccessLogsByPatient(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// args: [patientID]
	if len(args) != 1 {
		return shim.Error(fmt.Sprintf("参数个数不满足，需要1个参数，实际收到%d个", len(args)))
	}

	patientID := args[0]
	if patientID == "" {
		return shim.Error("患者ID不能为空")
	}

	// 构建查询条件
	queryString := fmt.Sprintf(`{"selector":{"docType":"AccessLog","patientID":"%s"},"sort":[{"accessedAt":"desc"}]}`, patientID)

	// 执行富查询
	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return shim.Error(fmt.Sprintf("查询失败: %s", err))
	}
	defer resultsIterator.Close()

	// 收集结果
	var logs []model.PrescriptionAccessLog
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(fmt.Sprintf("迭代失败: %s", err))
		}

		var log model.PrescriptionAccessLog
		if err := json.Unmarshal(queryResponse.Value, &log); err != nil {
			continue
		}
		logs = append(logs, log)
	}

	// 序列化结果
	logsBytes, err := json.Marshal(logs)
	if err != nil {
		return shim.Error(fmt.Sprintf("序列化失败: %s", err))
	}

	return shim.Success(logsBytes)
}
