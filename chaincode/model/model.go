package model

// Account 账户，虚拟管理员和若干业主账号
type Account struct {
	AccountId string  `json:"accountId"` //账号ID
	UserName  string  `json:"userName"`  //账号名
	Balance   float64 `json:"balance"`   //余额
}

// objectType  对象类型，用于创建复合主键
const (
	AccountKey = "account-key"

	AccountV2Key         = "account-v2-key"
	PrescriptionKey      = "prescription-key"
	AccessRequestKey     = "access-request-key"
	PatientKey           = "patient-key"
	InsuranceKey         = "insurance-key"
	DrugKey              = "drug-key"
	SupplementRecordKey  = "supplement-record-key"
	AccessLogKey         = "access-log-key"
)

// --------------------------------------------------------------------

// AccountV2 账号
type AccountV2 struct {
	AccountId      string `json:"account_id"`      // 账号ID
	AccountName    string `json:"account_name"`    // 账号名
	Username       string `json:"username"`        // 用户名
	Password       string `json:"password"`        // 密码
	Role           string `json:"role"`            // 角色：医生/病人/管理员
	Organization   string `json:"organization"`    // 所属组织（MSPID，如TaobaoMSP）
	OrganizationName string `json:"organization_name"` // 组织名称（如协和医院）
	Department     string `json:"department"`      // 科室（医生专用）
	DoctorTitle    string `json:"doctor_title"`    // 医生职称
	DoctorLicense  string `json:"doctor_license"`  // 医师执业证号
	Age            int    `json:"age"`             // 年龄（病人专用）
	Gender         string `json:"gender"`          // 性别（病人专用）
}

// Hospital 医院
type Hospital struct {
	ID      string          `json:"id"`
	Name    string          `json:"name"`
	Admins  []HospitalAdmin `json:"admins"`
	Doctors []Doctor        `json:"doctors"`
}

// HospitalAdmin 医院管理员
type HospitalAdmin struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Doctor 医生
type Doctor struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Prescription 医疗处方
type Prescription struct {
	ID               string   `json:"id"`                // 医疗处方ID
	PrescriptionNo   string   `json:"prescription_no"`   // 病历编号
	Patient          string   `json:"patient"`           // 患者ID
	PatientName      string   `json:"patient_name"`      // 患者姓名
	ChiefComplaint   string   `json:"chief_complaint"`   // 主诉
	PresentIllness   string   `json:"present_illness"`   // 现病史
	PhysicalExam     string   `json:"physical_exam"`     // 体格检查
	Diagnosis        string   `json:"diagnosis"`         // 诊断结果
	Drug             []Drug   `json:"drug"`              // 药品列表及用量
	DrugOrderInfo    string   `json:"drug_order_info"`   // 药品订单信息（JSON格式）
	MedicalAdvice    string   `json:"medical_advice"`    // 医嘱
	Doctor           string   `json:"doctor"`            // 开方医师ID
	DoctorName       string   `json:"doctor_name"`       // 医生姓名
	DoctorTitle      string   `json:"doctor_title"`      // 医生职称
	Hospital         string   `json:"hospital"`          // 医院ID
	HospitalName     string   `json:"hospital_name"`     // 医院名称
	OrganizationID   string   `json:"organization_id"`   // 组织MSPID（体现联盟）
	OrganizationName string   `json:"organization_name"` // 组织名称
	Department       string   `json:"department"`        // 科室
	Created          string   `json:"created"`           // 创建时间
	Comment          string   `json:"comment"`           // 备注
	TxID             string   `json:"tx_id"`             // 交易ID（区块链特性）
	CreatorMSPID     string   `json:"creator_msp_id"`    // 创建者MSPID（证明来源）
	AuthorizedOrgs   []string `json:"authorized_orgs"`   // 授权的组织列表（跨院授权）
}

// AccessRequest 病历访问授权请求
type AccessRequest struct {
	ID               string `json:"id"`                // 请求ID
	PrescriptionID   string `json:"prescription_id"`   // 病历ID
	PatientID        string `json:"patient_id"`        // 患者ID
	PatientName      string `json:"patient_name"`      // 患者姓名
	DoctorID         string `json:"doctor_id"`         // 申请医生ID
	DoctorName       string `json:"doctor_name"`       // 申请医生姓名
	DoctorOrg        string `json:"doctor_org"`        // 医生所属组织
	DoctorOrgName    string `json:"doctor_org_name"`   // 医生所属组织名称
	Reason           string `json:"reason"`            // 申请理由
	Status           string `json:"status"`            // 状态：pending/approved/rejected
	RequestTime      string `json:"request_time"`      // 申请时间
	ResponseTime     string `json:"response_time"`     // 响应时间
	TxID             string `json:"tx_id"`             // 交易ID
}

// Patient 患者
type Patient struct {
	ID     string `json:"id"`     // 患者 AccountV2Id
	Name   string `json:"name"`   // 患者姓名
	Age    int    `json:"age"`    // 患者年龄
	Gender string `json:"gender"` // 患者性别
}

// Drug 药品
type Drug struct {
	//ID      string `json:"id"`
	Name   string `json:"Name"`   // 药品名
	Amount string `json:"amount"` // 药品数量
}

// DrugOrder 药品订单
type DrugOrder struct {
	ID           string `json:"id"`           // 订单ID
	Name         string `json:"Name"`         // 药品名
	Amount       string `json:"amount"`       // 药品数量
	Prescription string `json:"prescription"` // 处方ID
	Patient      string `json:"patient"`      // 患者ID
	DrugStore    string `json:"drug_store"`   // 药店id
	Created      string `json:"created"`      // 创建时间
}

// DrugStore 药店
type DrugStore struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Insurance 保险机构
type Insurance struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// InsuranceCover 保险报销订单
type InsuranceCover struct {
	ID           string `json:"id"`           // 订单ID
	Prescription string `json:"prescription"` // 处方ID
	Patient      string `json:"patient"`      // 患者ID
	Status       string `json:"status"`       // 订单状态
	Created      string `json:"created"`      // 创建时间
}

// InsuranceStatusConstant 保险状态
var InsuranceStatusConstant = func() map[string]string {
	return map[string]string{
		"processing": "处理中", // 患者发起保险报销申请，等待保险公司确认报销
		"cancelled":  "已取消", // 患者在保险公司确认报销之前取消保险报销申请
		"refused":    "已拒绝", // 保险公司拒绝确认报销
		"approved":   "已通过", // 保险公司确认报销，保险报销完成
	}
}

// DrugStatusConstant 药品状态
//var DrugStatusConstant = func() map[string]string {
//	return map[string]string{
//		"processing": "处理中", //
//		"done":       "完成",   //
//	}
//}

// SupplementRecord 补充诊疗记录
type SupplementRecord struct {
	ID                     string   `json:"id"`                       // 补充记录ID(独立的区块链ID)
	OriginalPrescriptionID string   `json:"original_prescription_id"` // 关联的原始病历ID
	RecordType             string   `json:"record_type"`              // 记录类型: consultation/followup/emergency
	PatientID              string   `json:"patient_id"`               // 患者ID
	PatientName            string   `json:"patient_name"`             // 患者姓名
	
	// 诊疗内容
	ChiefComplaint         string   `json:"chief_complaint"`          // 主诉
	PresentIllness         string   `json:"present_illness"`          // 现病史
	PhysicalExam           string   `json:"physical_exam"`            // 体格检查
	Diagnosis              string   `json:"diagnosis"`                // 诊断
	Treatment              string   `json:"treatment"`                // 治疗方案
	Drug                   []Drug   `json:"drug"`                     // 用药
	MedicalAdvice          string   `json:"medical_advice"`           // 医嘱
	
	// 医生信息
	DoctorID               string   `json:"doctor_id"`                // 医生ID
	DoctorName             string   `json:"doctor_name"`              // 医生姓名
	DoctorTitle            string   `json:"doctor_title"`             // 医生职称
	Department             string   `json:"department"`               // 科室
	
	// 医院信息
	HospitalName           string   `json:"hospital_name"`            // 医院名称
	OrganizationID         string   `json:"organization_id"`          // 组织MSPID
	OrganizationName       string   `json:"organization_name"`        // 组织名称
	
	// 区块链信息
	Created                string   `json:"created"`                  // 创建时间
	TxID                   string   `json:"tx_id"`                    // 交易ID
	CreatorMSPID           string   `json:"creator_msp_id"`           // 创建者MSPID
	
	// 关联信息
	IsReadOnly             bool     `json:"is_read_only"`             // 是否只读(创建后不可修改)
	Comment                string   `json:"comment"`                  // 备注
}

// MedicalHistory 完整病历历史(用于查询)
type MedicalHistory struct {
	OriginalPrescription Prescription        `json:"original_prescription"` // 原始病历
	SupplementRecords    []SupplementRecord  `json:"supplement_records"`    // 补充记录列表
	TotalRecords         int                 `json:"total_records"`         // 总记录数
}

// PrescriptionAccessLog 病历访问日志（隐私溯源）
type PrescriptionAccessLog struct {
	LogID                    string `json:"log_id"`                     // 日志ID
	PrescriptionID           string `json:"prescription_id"`            // 病历ID
	PrescriptionNo           string `json:"prescription_no"`            // 病历编号
	PatientID                string `json:"patient_id"`                 // 患者ID
	PatientName              string `json:"patient_name"`               // 患者姓名
	AccessorID               string `json:"accessor_id"`                // 访问者ID
	AccessorName             string `json:"accessor_name"`              // 访问者姓名
	AccessorRole             string `json:"accessor_role"`              // 访问者角色
	AccessorOrganization     string `json:"accessor_organization"`      // 访问者组织MSPID
	AccessorOrganizationName string `json:"accessor_organization_name"` // 访问者组织名称
	AccessType               string `json:"access_type"`                // 访问类型：view/edit/download
	AccessReason             string `json:"access_reason"`              // 访问原因
	AccessedAt               string `json:"accessed_at"`                // 访问时间
}
