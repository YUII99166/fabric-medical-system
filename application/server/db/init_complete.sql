-- ============================================
-- 基于区块链的医疗信息管理系统 - 完整数据库表结构
-- ============================================

-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS fabric;

-- 使用数据库
USE fabric;

-- ============================================
-- 1. 用户表 (users)
-- ============================================
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT '用户ID',
    username VARCHAR(255) UNIQUE NOT NULL COMMENT '用户名（登录账号）',
    password VARCHAR(255) NOT NULL COMMENT '密码（SHA256哈希）',
    account_name VARCHAR(255) NOT NULL COMMENT '账户名称（显示名称）',
    role VARCHAR(50) NOT NULL COMMENT '角色：医生/病人/管理员/药店/保险机构',
    organization VARCHAR(100) DEFAULT '' COMMENT '所属组织MSPID（如TaobaoMSP）',
    organization_name VARCHAR(200) DEFAULT '' COMMENT '组织名称（如协和医院）',
    department VARCHAR(100) DEFAULT '' COMMENT '科室（医生专用）',
    doctor_title VARCHAR(50) DEFAULT '' COMMENT '医生职称',
    doctor_license VARCHAR(100) DEFAULT '' COMMENT '医师执业证号',
    age INT DEFAULT NULL COMMENT '年龄（病人专用）',
    gender VARCHAR(10) DEFAULT NULL COMMENT '性别（病人专用）',
    phone VARCHAR(20) DEFAULT NULL COMMENT '联系电话',
    email VARCHAR(100) DEFAULT NULL COMMENT '电子邮箱',
    status TINYINT DEFAULT 1 COMMENT '状态：1-正常，0-禁用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_users_username (username),
    INDEX idx_users_role (role),
    INDEX idx_users_organization (organization),
    INDEX idx_users_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- ============================================
-- 2. 病历表 (prescriptions) - 缓存区块链数据
-- ============================================
CREATE TABLE IF NOT EXISTS prescriptions (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT '自增ID',
    prescription_id VARCHAR(100) UNIQUE NOT NULL COMMENT '病历ID（区块链ID）',
    prescription_no VARCHAR(100) NOT NULL COMMENT '病历编号',
    patient_id VARCHAR(100) NOT NULL COMMENT '患者ID',
    patient_name VARCHAR(100) NOT NULL COMMENT '患者姓名',
    doctor_id VARCHAR(100) NOT NULL COMMENT '医生ID',
    doctor_name VARCHAR(100) NOT NULL COMMENT '医生姓名',
    doctor_title VARCHAR(50) DEFAULT NULL COMMENT '医生职称',
    hospital_id VARCHAR(100) DEFAULT NULL COMMENT '医院ID',
    hospital_name VARCHAR(200) DEFAULT NULL COMMENT '医院名称',
    organization_id VARCHAR(100) NOT NULL COMMENT '组织MSPID',
    organization_name VARCHAR(200) NOT NULL COMMENT '组织名称',
    department VARCHAR(100) DEFAULT NULL COMMENT '科室',
    chief_complaint TEXT COMMENT '主诉',
    present_illness TEXT COMMENT '现病史',
    physical_exam TEXT COMMENT '体格检查',
    diagnosis TEXT COMMENT '诊断',
    treatment_plan TEXT COMMENT '治疗方案',
    prescription_drugs TEXT COMMENT '处方药品（JSON格式）',
    drug_order_info TEXT COMMENT '药品订单信息（JSON格式）',
    medical_advice TEXT COMMENT '医嘱',
    comment TEXT COMMENT '备注',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    tx_id VARCHAR(100) DEFAULT NULL COMMENT '区块链交易ID',
    creator_mspid VARCHAR(100) DEFAULT NULL COMMENT '创建者MSPID',
    INDEX idx_prescriptions_patient (patient_id),
    INDEX idx_prescriptions_doctor (doctor_id),
    INDEX idx_prescriptions_organization (organization_id),
    INDEX idx_prescriptions_created (created_at),
    INDEX idx_prescriptions_no (prescription_no)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='病历表';

-- ============================================
-- 3. 授权请求表 (access_requests)
-- ============================================
CREATE TABLE IF NOT EXISTS access_requests (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT '自增ID',
    request_id VARCHAR(100) UNIQUE NOT NULL COMMENT '请求ID（区块链ID）',
    prescription_id VARCHAR(100) NOT NULL COMMENT '病历ID',
    patient_id VARCHAR(100) NOT NULL COMMENT '患者ID',
    patient_name VARCHAR(100) NOT NULL COMMENT '患者姓名',
    doctor_id VARCHAR(100) NOT NULL COMMENT '申请医生ID',
    doctor_name VARCHAR(100) NOT NULL COMMENT '申请医生姓名',
    doctor_org VARCHAR(100) NOT NULL COMMENT '医生所属组织MSPID',
    doctor_org_name VARCHAR(200) NOT NULL COMMENT '医生所属组织名称',
    reason TEXT COMMENT '申请理由',
    status VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '状态：pending-待审批，approved-已批准，rejected-已拒绝',
    reject_reason TEXT COMMENT '拒绝理由',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '申请时间',
    approved_at TIMESTAMP NULL DEFAULT NULL COMMENT '审批时间',
    INDEX idx_access_patient (patient_id),
    INDEX idx_access_doctor (doctor_id),
    INDEX idx_access_prescription (prescription_id),
    INDEX idx_access_status (status),
    INDEX idx_access_created (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='授权请求表';

-- ============================================
-- 4. 补充诊疗记录表 (supplement_records)
-- ============================================
CREATE TABLE IF NOT EXISTS supplement_records (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT '自增ID',
    record_id VARCHAR(100) UNIQUE NOT NULL COMMENT '补充记录ID（区块链ID）',
    original_prescription_id VARCHAR(100) NOT NULL COMMENT '原始病历ID',
    record_type VARCHAR(50) NOT NULL COMMENT '记录类型：复诊/随访/急诊',
    doctor_id VARCHAR(100) NOT NULL COMMENT '医生ID',
    doctor_name VARCHAR(100) NOT NULL COMMENT '医生姓名',
    doctor_title VARCHAR(50) DEFAULT NULL COMMENT '医生职称',
    department VARCHAR(100) DEFAULT NULL COMMENT '科室',
    hospital_name VARCHAR(200) DEFAULT NULL COMMENT '医院名称',
    organization_id VARCHAR(100) NOT NULL COMMENT '组织MSPID',
    organization_name VARCHAR(200) NOT NULL COMMENT '组织名称',
    symptoms TEXT COMMENT '症状描述',
    diagnosis TEXT COMMENT '诊断',
    treatment TEXT COMMENT '治疗方案',
    prescription_drugs TEXT COMMENT '处方药品（JSON格式）',
    medical_advice TEXT COMMENT '医嘱',
    comment TEXT COMMENT '备注',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    tx_id VARCHAR(100) DEFAULT NULL COMMENT '区块链交易ID',
    creator_mspid VARCHAR(100) DEFAULT NULL COMMENT '创建者MSPID',
    is_read_only TINYINT DEFAULT 1 COMMENT '是否只读：1-是，0-否',
    INDEX idx_supplement_original (original_prescription_id),
    INDEX idx_supplement_doctor (doctor_id),
    INDEX idx_supplement_organization (organization_id),
    INDEX idx_supplement_created (created_at),
    INDEX idx_supplement_type (record_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='补充诊疗记录表';

-- ============================================
-- 5. 药品订单表 (drug_orders)
-- ============================================
CREATE TABLE IF NOT EXISTS drug_orders (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT '自增ID',
    order_id VARCHAR(100) UNIQUE NOT NULL COMMENT '订单ID（区块链ID）',
    drug_name VARCHAR(200) NOT NULL COMMENT '药品名称',
    quantity INT NOT NULL COMMENT '数量',
    price DECIMAL(10, 2) NOT NULL COMMENT '价格',
    patient_id VARCHAR(100) NOT NULL COMMENT '患者ID',
    patient_name VARCHAR(100) NOT NULL COMMENT '患者姓名',
    prescription_id VARCHAR(100) NOT NULL COMMENT '关联处方ID',
    drug_store_id VARCHAR(100) NOT NULL COMMENT '药店ID',
    drug_store_name VARCHAR(200) NOT NULL COMMENT '药店名称',
    status VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '订单状态：pending-待处理，processing-处理中，completed-已完成，cancelled-已取消',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_drug_orders_patient (patient_id),
    INDEX idx_drug_orders_prescription (prescription_id),
    INDEX idx_drug_orders_drugstore (drug_store_id),
    INDEX idx_drug_orders_status (status),
    INDEX idx_drug_orders_created (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='药品订单表';

-- ============================================
-- 6. 保险报销表 (insurance_claims)
-- ============================================
CREATE TABLE IF NOT EXISTS insurance_claims (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT '自增ID',
    claim_id VARCHAR(100) UNIQUE NOT NULL COMMENT '报销订单ID（区块链ID）',
    prescription_id VARCHAR(100) NOT NULL COMMENT '处方ID',
    patient_id VARCHAR(100) NOT NULL COMMENT '患者ID',
    patient_name VARCHAR(100) NOT NULL COMMENT '患者姓名',
    insurance_id VARCHAR(100) NOT NULL COMMENT '保险机构ID',
    insurance_name VARCHAR(200) NOT NULL COMMENT '保险机构名称',
    claim_amount DECIMAL(10, 2) NOT NULL COMMENT '报销金额',
    status VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '状态：pending-待审核，approved-已批准，rejected-已拒绝',
    reject_reason TEXT COMMENT '拒绝理由',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '申请时间',
    approved_at TIMESTAMP NULL DEFAULT NULL COMMENT '审批时间',
    INDEX idx_insurance_patient (patient_id),
    INDEX idx_insurance_prescription (prescription_id),
    INDEX idx_insurance_company (insurance_id),
    INDEX idx_insurance_status (status),
    INDEX idx_insurance_created (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='保险报销表';

-- ============================================
-- 7. 组织机构表 (organizations)
-- ============================================
CREATE TABLE IF NOT EXISTS organizations (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT '自增ID',
    org_id VARCHAR(100) UNIQUE NOT NULL COMMENT '组织ID（MSPID）',
    org_name VARCHAR(200) NOT NULL COMMENT '组织名称',
    org_type VARCHAR(50) NOT NULL COMMENT '组织类型：医院/药店/保险机构/监管中心',
    domain VARCHAR(100) DEFAULT NULL COMMENT '域名',
    address VARCHAR(500) DEFAULT NULL COMMENT '地址',
    contact_person VARCHAR(100) DEFAULT NULL COMMENT '联系人',
    contact_phone VARCHAR(20) DEFAULT NULL COMMENT '联系电话',
    contact_email VARCHAR(100) DEFAULT NULL COMMENT '联系邮箱',
    description TEXT COMMENT '组织描述',
    status TINYINT DEFAULT 1 COMMENT '状态：1-正常，0-禁用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_organizations_type (org_type),
    INDEX idx_organizations_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='组织机构表';

-- ============================================
-- 8. 审计日志表 (audit_logs)
-- ============================================
CREATE TABLE IF NOT EXISTS audit_logs (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '自增ID',
    user_id VARCHAR(100) NOT NULL COMMENT '操作用户ID',
    username VARCHAR(255) NOT NULL COMMENT '操作用户名',
    user_role VARCHAR(50) NOT NULL COMMENT '用户角色',
    action VARCHAR(100) NOT NULL COMMENT '操作动作',
    resource_type VARCHAR(50) NOT NULL COMMENT '资源类型：prescription/access_request/drug_order等',
    resource_id VARCHAR(100) DEFAULT NULL COMMENT '资源ID',
    description TEXT COMMENT '操作描述',
    ip_address VARCHAR(50) DEFAULT NULL COMMENT 'IP地址',
    user_agent TEXT COMMENT '用户代理',
    tx_id VARCHAR(100) DEFAULT NULL COMMENT '区块链交易ID（如果有）',
    status VARCHAR(20) DEFAULT 'success' COMMENT '操作状态：success-成功，failed-失败',
    error_message TEXT COMMENT '错误信息（如果失败）',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '操作时间',
    INDEX idx_audit_user (user_id),
    INDEX idx_audit_action (action),
    INDEX idx_audit_resource (resource_type, resource_id),
    INDEX idx_audit_created (created_at),
    INDEX idx_audit_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='审计日志表';

-- ============================================
-- 9. 系统配置表 (system_configs)
-- ============================================
CREATE TABLE IF NOT EXISTS system_configs (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT '自增ID',
    config_key VARCHAR(100) UNIQUE NOT NULL COMMENT '配置键',
    config_value TEXT NOT NULL COMMENT '配置值',
    config_type VARCHAR(50) DEFAULT 'string' COMMENT '配置类型：string/number/boolean/json',
    description VARCHAR(500) DEFAULT NULL COMMENT '配置描述',
    is_editable TINYINT DEFAULT 1 COMMENT '是否可编辑：1-是，0-否',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_configs_key (config_key)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统配置表';

-- ============================================
-- 初始化组织机构数据
-- ============================================
INSERT INTO organizations (org_id, org_name, org_type, domain, description) VALUES
('TaobaoMSP', '协和医院', '医院', 'taobao.com', '北京协和医院，三甲综合医院'),
('JDMSP', '301医院', '医院', 'jd.com', '中国人民解放军总医院（301医院）'),
('WenjinMSP', '温江社区医疗中心', '医院', 'wenjin.com', '成都市温江区社区医疗服务中心'),
('RegCenterMSP', '监管中心', '监管中心', 'regcenter.com', '医疗数据监管中心')
ON DUPLICATE KEY UPDATE org_name=VALUES(org_name);

-- ============================================
-- 初始化系统配置数据
-- ============================================
INSERT INTO system_configs (config_key, config_value, config_type, description) VALUES
('system_name', '基于区块链的社区医疗管理信息系统', 'string', '系统名称'),
('blockchain_network', 'appchannel', 'string', '区块链通道名称'),
('max_upload_size', '10485760', 'number', '最大上传文件大小（字节）'),
('session_timeout', '3600', 'number', '会话超时时间（秒）'),
('enable_audit_log', 'true', 'boolean', '是否启用审计日志')
ON DUPLICATE KEY UPDATE config_value=VALUES(config_value);

-- ============================================
-- 插入默认测试账户
-- ============================================
INSERT IGNORE INTO users (username, password, account_name, role, organization, organization_name) VALUES
('admin', '240be518fabd2724ddb6f04eeb1da5967448d7e831c08c8fa822809f74c720a9', '系统管理员', '管理员', 'TaobaoMSP', '协和医院'),
('doctor1', 'ef92b778bafe771e89245b89ecbc08a44a4e166c06659911881f383d4473e94f', '张医生', '医生', 'TaobaoMSP', '协和医院'),
('doctor2', 'ef92b778bafe771e89245b89ecbc08a44a4e166c06659911881f383d4473e94f', '李医生', '医生', 'JDMSP', '301医院'),
('patient1', 'a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3', '张三', '病人', '', ''),
('patient2', 'a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3', '李四', '病人', '', ''),
('drugstore1', 'a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3', '仁和药店', '药店', '', ''),
('insurance1', 'a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3', '平安保险', '保险机构', '', '');

-- ============================================
-- 创建视图 - 病历统计视图
-- ============================================
CREATE OR REPLACE VIEW v_prescription_statistics AS
SELECT 
    organization_id,
    organization_name,
    COUNT(*) as total_prescriptions,
    COUNT(DISTINCT patient_id) as total_patients,
    COUNT(DISTINCT doctor_id) as total_doctors,
    DATE(created_at) as date
FROM prescriptions
GROUP BY organization_id, organization_name, DATE(created_at);

-- ============================================
-- 创建视图 - 用户统计视图
-- ============================================
CREATE OR REPLACE VIEW v_user_statistics AS
SELECT 
    role,
    organization,
    organization_name,
    COUNT(*) as user_count,
    SUM(CASE WHEN status = 1 THEN 1 ELSE 0 END) as active_count,
    SUM(CASE WHEN status = 0 THEN 1 ELSE 0 END) as inactive_count
FROM users
GROUP BY role, organization, organization_name;

-- ============================================
-- 完成
-- ============================================
