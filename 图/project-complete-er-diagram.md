# 基于区块链的医疗信息管理系统 - 完整项目ER图

## 项目概述

这是一个基于Hyperledger Fabric区块链的社区医疗管理信息系统，实现了跨组织的医疗数据共享、隐私保护和授权管理。

## 核心实体关系图 (Mermaid格式)

```mermaid
erDiagram
    %% 核心实体定义
    USER ||--o{ PRESCRIPTION : "创建/拥有"
    USER ||--o{ SUPPLEMENT_RECORD : "添加"
    USER ||--o{ ACCESS_REQUEST : "申请/审批"
    USER ||--o{ DRUG_ORDER : "下单/处理"
    USER ||--o{ INSURANCE_CLAIM : "申请/审批"
    USER ||--o{ AUDIT_LOG : "产生"
    
    PRESCRIPTION ||--o{ SUPPLEMENT_RECORD : "包含"
    PRESCRIPTION ||--o{ ACCESS_REQUEST : "关联"
    PRESCRIPTION ||--o{ DRUG_ORDER : "基于"
    PRESCRIPTION ||--o{ INSURANCE_CLAIM : "关联"
    
    ORGANIZATION ||--o{ USER : "所属"
    ORGANIZATION ||--o{ PRESCRIPTION : "创建"
    
    %% 用户实体
    USER {
        int id PK "用户ID"
        string username UK "用户名(登录账号)"
        string password "密码(SHA256)"
        string account_name "账户名称"
        string role "角色:医生/病人/管理员/药店/保险"
        string organization FK "所属组织MSPID"
        string organization_name "组织名称"
        string department "科室(医生)"
        string doctor_title "医生职称"
        string doctor_license "医师执业证号"
        int age "年龄(病人)"
        string gender "性别"
        string phone "联系电话"
        string email "电子邮箱"
        tinyint status "状态:1-正常,0-禁用"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
    }
    
    %% 病历实体
    PRESCRIPTION {
        int id PK "自增ID"
        string prescription_id UK "病历ID(区块链)"
        string prescription_no "病历编号"
        string patient_id FK "患者ID"
        string patient_name "患者姓名"
        string doctor_id FK "医生ID"
        string doctor_name "医生姓名"
        string doctor_title "医生职称"
        string hospital_id "医院ID"
        string hospital_name "医院名称"
        string organization_id FK "组织MSPID"
        string organization_name "组织名称"
        string department "科室"
        text chief_complaint "主诉"
        text present_illness "现病史"
        text physical_exam "体格检查"
        text diagnosis "诊断"
        text treatment_plan "治疗方案"
        text prescription_drugs "处方药品(JSON)"
        text medical_advice "医嘱"
        text comment "备注"
        timestamp created_at "创建时间"
        string tx_id "区块链交易ID"
        string creator_mspid "创建者MSPID"
    }
    
    %% 补充记录实体(弱实体)
    SUPPLEMENT_RECORD {
        int id PK "自增ID"
        string record_id UK "补充记录ID(区块链)"
        string original_prescription_id FK "原始病历ID"
        string record_type "记录类型:复诊/随访/急诊"
        string doctor_id FK "医生ID"
        string doctor_name "医生姓名"
        string doctor_title "医生职称"
        string department "科室"
        string hospital_name "医院名称"
        string organization_id FK "组织MSPID"
        string organization_name "组织名称"
        text symptoms "症状描述"
        text diagnosis "诊断"
        text treatment "治疗方案"
        text prescription_drugs "处方药品(JSON)"
        text medical_advice "医嘱"
        text comment "备注"
        timestamp created_at "创建时间"
        string tx_id "区块链交易ID"
        string creator_mspid "创建者MSPID"
        tinyint is_read_only "是否只读"
    }
    
    %% 授权请求实体
    ACCESS_REQUEST {
        int id PK "自增ID"
        string request_id UK "请求ID(区块链)"
        string prescription_id FK "病历ID"
        string patient_id FK "患者ID"
        string patient_name "患者姓名"
        string doctor_id FK "申请医生ID"
        string doctor_name "申请医生姓名"
        string doctor_org FK "医生所属组织MSPID"
        string doctor_org_name "医生所属组织名称"
        text reason "申请理由"
        string status "状态:pending/approved/rejected"
        text reject_reason "拒绝理由"
        timestamp created_at "申请时间"
        timestamp approved_at "审批时间"
    }
    
    %% 药品订单实体
    DRUG_ORDER {
        int id PK "自增ID"
        string order_id UK "订单ID(区块链)"
        string drug_name "药品名称"
        int quantity "数量"
        decimal price "价格"
        string patient_id FK "患者ID"
        string patient_name "患者姓名"
        string prescription_id FK "关联处方ID"
        string drug_store_id FK "药店ID"
        string drug_store_name "药店名称"
        string status "订单状态:pending/processing/completed/cancelled"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
    }
    
    %% 保险报销实体
    INSURANCE_CLAIM {
        int id PK "自增ID"
        string claim_id UK "报销订单ID(区块链)"
        string prescription_id FK "处方ID"
        string patient_id FK "患者ID"
        string patient_name "患者姓名"
        string insurance_id FK "保险机构ID"
        string insurance_name "保险机构名称"
        decimal claim_amount "报销金额"
        string status "状态:pending/approved/rejected"
        text reject_reason "拒绝理由"
        timestamp created_at "申请时间"
        timestamp approved_at "审批时间"
    }
    
    %% 组织机构实体
    ORGANIZATION {
        int id PK "自增ID"
        string org_id UK "组织ID(MSPID)"
        string org_name "组织名称"
        string org_type "组织类型:医院/药店/保险/监管"
        string domain "域名"
        string address "地址"
        string contact_person "联系人"
        string contact_phone "联系电话"
        string contact_email "联系邮箱"
        text description "组织描述"
        tinyint status "状态:1-正常,0-禁用"
        timestamp created_at "创建时间"
        timestamp updated_at "更新时间"
    }
    
    %% 审计日志实体
    AUDIT_LOG {
        bigint id PK "自增ID"
        string user_id FK "操作用户ID"
        string username "操作用户名"
        string user_role "用户角色"
        string action "操作动作"
        string resource_type "资源类型"
        string resource_id "资源ID"
        text description "操作描述"
        string ip_address "IP地址"
        text user_agent "用户代理"
        string tx_id "区块链交易ID"
        string status "操作状态:success/failed"
        text error_message "错误信息"
        timestamp created_at "操作时间"
    }
```

## 系统架构层次

### 1. 用户层 (User Layer)
- **医生 (Doctor)**: 创建病历、添加补充记录、申请跨组织授权
- **患者 (Patient)**: 拥有病历、审批授权请求、下药品订单、申请保险报销
- **药店 (Pharmacy)**: 处理药品订单
- **保险机构 (Insurance)**: 审批保险报销
- **管理员 (Admin)**: 系统管理、用户管理、审计监控

### 2. 业务层 (Business Layer)
- **病历管理**: 创建、查询、更新病历信息
- **补充记录**: 复诊、随访、急诊记录
- **授权管理**: 跨组织访问授权申请与审批
- **药品订单**: 基于处方的药品购买
- **保险报销**: 基于处方的保险理赔

### 3. 区块链层 (Blockchain Layer)
- **Hyperledger Fabric**: 联盟链底层
- **智能合约 (Chaincode)**: 业务逻辑执行
- **通道 (Channel)**: appchannel
- **组织 (Organizations)**: 
  - TaobaoMSP (协和医院)
  - JDMSP (301医院)
  - WenjinMSP (温江社区医疗中心)
  - RegCenterMSP (监管中心)

### 4. 数据层 (Data Layer)
- **MySQL数据库**: 缓存区块链数据，提供快速查询
- **区块链账本**: 不可篡改的数据存储

## 核心业务流程

### 流程1: 病历创建与管理
```
医生登录 → 创建病历 → 写入区块链 → 缓存到MySQL → 患者可查看
```

### 流程2: 跨组织授权
```
外院医生申请授权 → 患者收到通知 → 患者审批 → 授权写入区块链 → 医生可访问病历
```

### 流程3: 补充记录
```
医生查看病历 → 添加补充记录 → 写入区块链 → 关联到原病历 → 更新病历历史
```

### 流程4: 药品订单
```
患者查看处方 → 选择药店下单 → 订单写入区块链 → 药店处理订单 → 订单完成
```

### 流程5: 保险报销
```
患者提交报销申请 → 保险机构审核 → 审批结果写入区块链 → 患者收到通知
```

## 数据完整性约束

### 实体完整性
- 所有表的主键不能为空且唯一
- 区块链ID (prescription_id, record_id等) 作为唯一键

### 参照完整性
- PRESCRIPTION.doctor_id → USER.id
- SUPPLEMENT_RECORD.original_prescription_id → PRESCRIPTION.prescription_id
- ACCESS_REQUEST.prescription_id → PRESCRIPTION.prescription_id
- DRUG_ORDER.prescription_id → PRESCRIPTION.prescription_id
- INSURANCE_CLAIM.prescription_id → PRESCRIPTION.prescription_id
- USER.organization → ORGANIZATION.org_id

### 用户定义完整性
- 角色字段: 医生/病人/管理员/药店/保险机构
- 状态字段: pending/approved/rejected/completed/cancelled
- 数量、价格必须大于0
- 时间戳字段不能为空

## 技术栈

### 后端
- **语言**: Go (Golang)
- **框架**: Gin Web Framework
- **数据库**: MySQL 8.0
- **区块链**: Hyperledger Fabric 2.x
- **SDK**: Fabric SDK Go

### 前端
- **框架**: Vue.js 3
- **UI库**: Element Plus
- **路由**: Vue Router
- **状态管理**: Pinia

### 区块链
- **平台**: Hyperledger Fabric
- **智能合约**: Go Chaincode
- **共识算法**: Raft
- **通道**: appchannel

## 安全特性

1. **身份认证**: 基于用户名密码的认证机制
2. **权限控制**: 基于角色的访问控制(RBAC)
3. **数据加密**: 密码SHA256哈希存储
4. **区块链不可篡改**: 所有关键操作写入区块链
5. **审计日志**: 完整的操作日志记录
6. **跨组织授权**: 患者控制的数据访问权限

## 系统特点

1. **去中心化**: 多组织联盟链架构
2. **数据隐私**: 患者数据加密存储，授权访问
3. **可追溯性**: 所有操作记录在区块链上
4. **跨组织协作**: 支持不同医疗机构间的数据共享
5. **高性能**: MySQL缓存提供快速查询
6. **可扩展**: 支持新组织加入联盟链

---

**生成时间**: 2026-04-03  
**系统版本**: v1.0  
**文档作者**: Kiro AI Assistant
