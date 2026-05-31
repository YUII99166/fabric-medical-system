package main

import (
	"chaincode/api"
	"chaincode/model"
	"chaincode/pkg/utils"
	"fmt"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type BlockChainMedicalInfoManageSystem struct {
}

// Init 链码部署到链上并进行初始化时会执行该方法
func (t *BlockChainMedicalInfoManageSystem) Init(stub shim.ChaincodeStubInterface) pb.Response { // stub 是智能合约中的一个对象，用于与区块链网络进行交互
	fmt.Println("链码初始化")
	//初始化默认数据
	// 注意：生产环境中不应该在链码中硬编码测试账号
	// 建议通过注册接口动态创建用户账号

	// 测试账号已移除，请通过系统注册功能创建用户
	// 如需初始化管理员账号，请在部署后通过API接口创建

	return shim.Success(nil)
}

// Invoke 实现Invoke接口调用智能合约
func (t *BlockChainMedicalInfoManageSystem) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	// 获取调用该智能合约函数时传入的参数和函数名
	// 这些信息可以用于智能合约内部的逻辑处理，例如根据不同的函数名和参数执行不同的操作或者检查参数的有效性等
	funcName, args := stub.GetFunctionAndParameters() // 返回两个值：funcName 表示函数名，args 表示函数参数列表，以字符串数组的形式返回
	switch funcName {
	case "hello":
		return api.Hello(stub, args)
	// api v2
	case "createAccountV2":
		return api.CreateAccountV2(stub, args)
	case "queryAccountV2List":
		return api.QueryAccountV2List(stub, args)
	case "register":
		return api.Register(stub, args)
	case "loginWithPassword":
		return api.LoginWithPassword(stub, args)
	case "createPrescription":
		return api.CreatePrescription(stub, args)
	case "queryPrescription":
		return api.QueryPrescription(stub, args)
	// 保险功能暂时注释
	//case "createInsuranceCover":
	//	return api.CreateInsuranceCover(stub, args)
	//case "queryInsuranceCover":
	//	return api.QueryInsuranceCover(stub, args)
	//case "updateInsuranceCover":
	//	return api.UpdateInsuranceCover(stub, args)
	//case "deleteInsuranceCover":
	//	return api.DeleteInsuranceCover(stub, args)
	case "createDrugOrder":
		return api.CreateDrugOrder(stub, args)
	case "queryDrugOrder":
		return api.QueryDrugOrder(stub, args)
	case "requestAccess":
		return api.RequestAccess(stub, args)
	case "approveAccess":
		return api.ApproveAccess(stub, args)
	case "queryAccessRequests":
		return api.QueryAccessRequests(stub, args)
	case "queryPrescriptionsByPatient":
		return api.QueryPrescriptionsByPatient(stub, args)
	case "addSupplementRecord":
		return api.AddSupplementRecord(stub, args)
	case "querySupplementRecords":
		return api.QuerySupplementRecords(stub, args)
	case "queryFullMedicalHistory":
		return api.QueryFullMedicalHistory(stub, args)
	case "recordPrescriptionAccess":
		return api.RecordPrescriptionAccess(stub, args)
	case "queryAccessLogsByPatient":
		return api.QueryAccessLogsByPatient(stub, args)
	default:
		return shim.Error(fmt.Sprintf("没有该功能: %s", funcName))
	}
}

func main() {
	timeLocal, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
	time.Local = timeLocal
	err = shim.Start(new(BlockChainMedicalInfoManageSystem)) // 启动智能合约
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
