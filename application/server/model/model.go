package model

// ----------------------         Account 用户   ----------------------------------

type AccountIdBody struct {
	AccountId string `json:"account_id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

type AccountRequestBody struct {
	Args []AccountIdBody `json:"args"`
}

type CreateAccountBody struct {
	AccountName string `json:"account_name"`
	Operator    string `json:"operator"`
}

type RegisterBody struct {
	AccountName      string `json:"account_name"`
	Username         string `json:"username"`
	Password         string `json:"password"`
	Role             string `json:"role"`
	Organization     string `json:"organization"`      // 组织MSPID
	OrganizationName string `json:"organization_name"` // 组织名称
	Department       string `json:"department"`        // 科室
	DoctorTitle      string `json:"doctor_title"`      // 医生职称
	Age              int    `json:"age"`               // 年龄（病人）
	Gender           string `json:"gender"`            // 性别（病人）
	Operator         string `json:"operator"`
}

type LoginBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateUserBody struct {
	ID               int    `json:"id"`
	AccountName      string `json:"account_name"`
	Role             string `json:"role"`
	Organization     string `json:"organization"`
	OrganizationName string `json:"organization_name"`
	Department       string `json:"department"`
	DoctorTitle      string `json:"doctor_title"`
	Password         string `json:"password"` // 可选，仅在修改密码时提供
}

type DeleteUserBody struct {
	ID int `json:"id"`
}

type BatchDeleteUsersBody struct {
	IDs []int `json:"ids"`
}

type SyncAccountBody struct {
	Username string `json:"username"`
}

// ----------------------         Prescription 病历   ----------------------------------

type PrescriptionRequestBody struct {
	Doctor         string `json:"doctor"`          // 医生ID
	Patient        string `json:"patient"`         // 患者Id
	ChiefComplaint string `json:"chief_complaint"` // 主诉
	PresentIllness string `json:"present_illness"` // 现病史
	PhysicalExam   string `json:"physical_exam"`   // 体格检查
	Diagnosis      string `json:"diagnosis"`       // 诊断结果
	DrugName       string `json:"drug_name"`       // 药品名
	DrugAmount     string `json:"drug_amount"`     // 药品用量
	DrugOrderInfo  string `json:"drug_order_info"` // 药品订单信息（JSON格式）
	MedicalAdvice  string `json:"medical_advice"`  // 医嘱
	Hospital       string `json:"hospital"`        // 医院 ID
	Comment        string `json:"comment"`         // 备注
}

type PrescriptionQueryRequestBody struct {
	Patient string `json:"patient"` // 患者AccountId
}

// ----------------------         DrugOrder 药品订单   ----------------------------------

type DrugOrderRequestBody struct {
	//Drug      []Drug `json:"drug"`      // 药品列表及用量
	DrugName     string `json:"drug_name"`    // 药品名
	DrugAmount   string `json:"drug_amount"`  // 药品用量
	Prescription string `json:"prescription"` // 处方ID
	Patient      string `json:"patient"`      // 患者Id
	DrugStore    string `json:"drug_store"`   // 药店Id
}

type DrugOrderQueryRequestBody struct {
	Patient   string `json:"patient"` // 患者AccountId
	DrugStore string `json:"drug_store"`
}

// ----------------------         InsuranceCover 保险报销   ----------------------------------

type InsuranceCoverRequestBody struct {
	Prescription string `json:"prescription"` // 处方ID
	Patient      string `json:"patient"`      // 患者Id
	Status       string `json:"status"`       // 订单状态
}

type InsuranceCoverQueryRequestBody struct {
	Patient        string `json:"patient"`         // 患者Id
	InsuranceCover string `json:"insurance_cover"` // 报销订单ID
}

type UpdateInsuranceCoverRequestBody struct {
	InsuranceCover string `json:"insurance_cover"` // 报销订单ID
	Patient        string `json:"patient"`         // 病人ID
	InsuranceID    string `json:"insurance_id"`    // 保险机构ID
	Status         string `json:"status"`          // 订单状态
}

// ----------------------         Access 授权管理   ----------------------------------

type AccessRequestBody struct {
	PrescriptionID string `json:"prescription_id"` // 病历ID
	DoctorID       string `json:"doctor_id"`       // 医生ID
	Reason         string `json:"reason"`          // 申请理由
}

type ApproveAccessBody struct {
	RequestID    string `json:"request_id"`    // 请求ID
	PatientID    string `json:"patient_id"`    // 患者ID
	Approved     string `json:"approved"`      // 是否同意（true/false）
	RejectReason string `json:"reject_reason"` // 拒绝理由
}

type QueryAccessRequestsBody struct {
	UserID string `json:"user_id"` // 用户ID
	Role   string `json:"role"`    // 角色（patient/doctor）
}

type QueryPrescriptionsByPatientBody struct {
	SearchKey string `json:"search_key"` // 患者姓名或ID
	DoctorID  string `json:"doctor_id"`  // 医生ID
}

// ----------------------         SupplementRecord 补充诊疗记录   ----------------------------------

type SupplementRecordRequestBody struct {
	OriginalPrescriptionID string `json:"original_prescription_id"` // 原始病历ID
	DoctorID               string `json:"doctor_id"`                // 医生ID
	RecordType             string `json:"record_type"`              // 记录类型
	ChiefComplaint         string `json:"chief_complaint"`          // 主诉
	PresentIllness         string `json:"present_illness"`          // 现病史
	PhysicalExam           string `json:"physical_exam"`            // 体格检查
	Diagnosis              string `json:"diagnosis"`                // 诊断
	Treatment              string `json:"treatment"`                // 治疗方案
	DrugName               string `json:"drug_name"`                // 药品名
	DrugAmount             string `json:"drug_amount"`              // 药品用量
	MedicalAdvice          string `json:"medical_advice"`           // 医嘱
	Comment                string `json:"comment"`                  // 备注
}

type QuerySupplementRecordsBody struct {
	OriginalPrescriptionID string `json:"original_prescription_id"` // 原始病历ID（可选）
}

type QueryFullMedicalHistoryBody struct {
	PrescriptionID string `json:"prescription_id"` // 病历ID
}

