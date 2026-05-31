-- ============================================
-- Fabric MIMS 数据库完整表结构恢复
-- 基于《数据库设计详细说明.md》
-- ============================================

USE fabric_mims;

-- 1. 用户表（users）
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT '用户ID',
    username VARCHAR(255) UNIQUE NOT NULL COMMENT '用户名（登录账号）',
    password VARCHAR(255) NOT NULL COMMENT '密码（SHA256哈希）',
    account_name VARCHAR(255) NOT NULL COMMENT '账号名称（显示名）',
    role VARCHAR(50) NOT NULL COMMENT '角色类型',
    organization VARCHAR(255) DEFAULT NULL COMMENT '组织MSPID',
    organization_name VARCHAR(255) DEFAULT NULL COMMENT '组织名称',
    department VARCHAR(255) DEFAULT NULL COMMENT '科室（医生专用）',
    doctor_title VARCHAR(100) DEFAULT NULL COMMENT '医生职称',
    age INT DEFAULT NULL COMMENT '年龄（患者专用）',
    gender VARCHAR(10) DEFAULT NULL COMMENT '性别（患者专用）',
    status INT DEFAULT 1 COMMENT '状态：1-正常，0-禁用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_role (role),
    INDEX idx_organization (organization),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 2. 病历缓存表（prescription_cache）
CREATE TABLE IF NOT EXISTS prescription_cache (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT '自增ID',
    prescription_id VARCHAR(255) UNIQUE NOT NULL COMMENT '病历ID（区块链）',
    prescription_no VARCHAR(255) UNIQUE NOT NULL COMMENT '病历编号',
    patient_id VARCHAR(255) NOT NULL COMMENT '患者ID',
    patient_name VARCHAR(255) NOT NULL COMMENT '患者姓名',
    doctor_id VARCHAR(255) NOT NULL COMMENT '医生ID',
    doctor_name VARCHAR(255) NOT NULL COMMENT '医生姓名',
    hospital_id VARCHAR(255) DEFAULT NULL COMMENT '医院ID',
    hospital_name VARCHAR(255) DEFAULT NULL COMMENT '医院名称',
    organization_id VARCHAR(255) DEFAULT NULL COMMENT '组织MSPID',
    chief_complaint TEXT DEFAULT NULL COMMENT '主诉',
    diagnosis TEXT DEFAULT NULL COMMENT '诊断结果',
    medical_advice TEXT DEFAULT NULL COMMENT '医嘱',
    tx_id VARCHAR(255) DEFAULT NULL COMMENT '区块链交易ID',
    block_number BIGINT DEFAULT NULL COMMENT '区块号',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    synced_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '同步时间',
    INDEX idx_patient_id (patient_id),
    INDEX idx_doctor_id (doctor_id),
    INDEX idx_hospital_id (hospital_id),
    INDEX idx_created_at (created_at),
    INDEX idx_composite (patient_id, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='病历缓存表';

-- 3. 补充记录缓存表（supplement_record_cache）
CREATE TABLE IF NOT EXISTS supplement_record_cache (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT '自增ID',
    supplement_id VARCHAR(255) UNIQUE NOT NULL COMMENT '补充记录ID',
    prescription_id VARCHAR(255) NOT NULL COMMENT '原始病历ID',
    doctor_id VARCHAR(255) NOT NULL COMMENT '医生ID',
    doctor_name VARCHAR(255) NOT NULL COMMENT '医生姓名',
    record_type VARCHAR(50) NOT NULL COMMENT '记录类型',
    content TEXT DEFAULT NULL COMMENT '补充内容',
    tx_id VARCHAR(255) DEFAULT NULL COMMENT '区块链交易ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    synced_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '同步时间',
    INDEX idx_prescription_id (prescription_id),
    INDEX idx_doctor_id (doctor_id),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='补充记录缓存表';

-- 4. 授权请求缓存表（access_request_cache）
CREATE TABLE IF NOT EXISTS access_request_cache (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT '自增ID',
    request_id VARCHAR(255) UNIQUE NOT NULL COMMENT '请求ID',
    prescription_id VARCHAR(255) NOT NULL COMMENT '病历ID',
    patient_id VARCHAR(255) NOT NULL COMMENT '患者ID',
    doctor_id VARCHAR(255) NOT NULL COMMENT '医生ID',
    doctor_org VARCHAR(255) DEFAULT NULL COMMENT '医生组织',
    reason TEXT DEFAULT NULL COMMENT '申请理由',
    status VARCHAR(50) NOT NULL DEFAULT 'pending' COMMENT '状态',
    reject_reason TEXT DEFAULT NULL COMMENT '拒绝理由',
    tx_id VARCHAR(255) DEFAULT NULL COMMENT '区块链交易ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '申请时间',
    response_at TIMESTAMP DEFAULT NULL COMMENT '响应时间',
    synced_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '同步时间',
    INDEX idx_prescription_id (prescription_id),
    INDEX idx_patient_id (patient_id),
    INDEX idx_doctor_id (doctor_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='授权请求缓存表';

-- 5. 药品订单缓存表（drug_order_cache）
CREATE TABLE IF NOT EXISTS drug_order_cache (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT '自增ID',
    order_id VARCHAR(255) UNIQUE NOT NULL COMMENT '订单ID',
    order_no VARCHAR(255) UNIQUE NOT NULL COMMENT '订单编号',
    patient_id VARCHAR(255) NOT NULL COMMENT '患者ID',
    patient_name VARCHAR(255) NOT NULL COMMENT '患者姓名',
    prescription_id VARCHAR(255) NOT NULL COMMENT '病历ID',
    drug_name VARCHAR(255) NOT NULL COMMENT '药品名称',
    quantity INT NOT NULL COMMENT '数量',
    price DECIMAL(10,2) NOT NULL COMMENT '单价',
    total_amount DECIMAL(10,2) NOT NULL COMMENT '总金额',
    pharmacy_id VARCHAR(255) NOT NULL COMMENT '药店ID',
    pharmacy_name VARCHAR(255) NOT NULL COMMENT '药店名称',
    status VARCHAR(50) NOT NULL DEFAULT 'pending' COMMENT '订单状态',
    tx_id VARCHAR(255) DEFAULT NULL COMMENT '区块链交易ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    synced_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '同步时间',
    INDEX idx_patient_id (patient_id),
    INDEX idx_prescription_id (prescription_id),
    INDEX idx_pharmacy_id (pharmacy_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='药品订单缓存表';

-- 6. 保险报销缓存表（insurance_cover_cache）
CREATE TABLE IF NOT EXISTS insurance_cover_cache (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT '自增ID',
    cover_id VARCHAR(255) UNIQUE NOT NULL COMMENT '报销ID',
    prescription_id VARCHAR(255) NOT NULL COMMENT '病历ID',
    patient_id VARCHAR(255) NOT NULL COMMENT '患者ID',
    insurance_id VARCHAR(255) DEFAULT NULL COMMENT '保险机构ID',
    amount DECIMAL(10,2) DEFAULT NULL COMMENT '报销金额',
    status VARCHAR(50) NOT NULL DEFAULT 'pending' COMMENT '状态',
    tx_id VARCHAR(255) DEFAULT NULL COMMENT '区块链交易ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    synced_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '同步时间',
    INDEX idx_prescription_id (prescription_id),
    INDEX idx_patient_id (patient_id),
    INDEX idx_insurance_id (insurance_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='保险报销缓存表';

-- 7. 会话表（sessions）
CREATE TABLE IF NOT EXISTS sessions (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT '自增ID',
    user_id INT NOT NULL COMMENT '用户ID',
    token VARCHAR(500) UNIQUE NOT NULL COMMENT 'JWT令牌',
    ip_address VARCHAR(50) DEFAULT NULL COMMENT '登录IP',
    user_agent VARCHAR(500) DEFAULT NULL COMMENT '用户代理',
    expires_at TIMESTAMP NOT NULL COMMENT '过期时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_user_id (user_id),
    INDEX idx_expires_at (expires_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='会话表';

-- 8. 操作日志表（operation_logs）
CREATE TABLE IF NOT EXISTS operation_logs (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '自增ID',
    user_id INT DEFAULT NULL COMMENT '操作用户ID',
    username VARCHAR(255) DEFAULT NULL COMMENT '用户名',
    operation VARCHAR(255) NOT NULL COMMENT '操作类型',
    resource_type VARCHAR(100) DEFAULT NULL COMMENT '资源类型',
    resource_id VARCHAR(255) DEFAULT NULL COMMENT '资源ID',
    method VARCHAR(10) DEFAULT NULL COMMENT 'HTTP方法',
    path VARCHAR(500) DEFAULT NULL COMMENT '请求路径',
    ip_address VARCHAR(50) DEFAULT NULL COMMENT 'IP地址',
    user_agent VARCHAR(500) DEFAULT NULL COMMENT '用户代理',
    status VARCHAR(50) DEFAULT NULL COMMENT '操作状态',
    error_message TEXT DEFAULT NULL COMMENT '错误信息',
    duration INT DEFAULT NULL COMMENT '执行时长(ms)',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '操作时间',
    INDEX idx_user_id (user_id),
    INDEX idx_operation (operation),
    INDEX idx_resource_type (resource_type),
    INDEX idx_created_at (created_at),
    INDEX idx_composite (user_id, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='操作日志表';

-- 9. 病历访问日志表（prescription_access_logs）
CREATE TABLE IF NOT EXISTS prescription_access_logs (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '自增ID',
    log_id VARCHAR(100) UNIQUE NOT NULL COMMENT '日志ID（区块链ID）',
    prescription_id VARCHAR(100) NOT NULL COMMENT '病历ID',
    prescription_no VARCHAR(100) NOT NULL COMMENT '病历编号',
    patient_id VARCHAR(100) NOT NULL COMMENT '患者ID',
    patient_name VARCHAR(100) NOT NULL COMMENT '患者姓名',
    accessor_id VARCHAR(100) NOT NULL COMMENT '访问者ID',
    accessor_name VARCHAR(100) NOT NULL COMMENT '访问者姓名',
    accessor_role VARCHAR(50) NOT NULL COMMENT '访问者角色',
    accessor_organization VARCHAR(100) DEFAULT NULL COMMENT '访问者组织MSPID',
    accessor_organization_name VARCHAR(200) DEFAULT NULL COMMENT '访问者组织名称',
    access_type VARCHAR(50) NOT NULL COMMENT '访问类型：view-查看, edit-编辑, download-下载',
    access_reason VARCHAR(500) DEFAULT NULL COMMENT '访问原因',
    ip_address VARCHAR(50) DEFAULT NULL COMMENT 'IP地址',
    user_agent TEXT DEFAULT NULL COMMENT '用户代理',
    tx_id VARCHAR(100) DEFAULT NULL COMMENT '区块链交易ID',
    accessed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '访问时间',
    INDEX idx_access_prescription (prescription_id),
    INDEX idx_access_patient (patient_id),
    INDEX idx_access_accessor (accessor_id),
    INDEX idx_access_time (accessed_at),
    INDEX idx_access_type (access_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='病历访问日志表';

-- 10. 系统配置表（system_config）
CREATE TABLE IF NOT EXISTS system_config (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT '自增ID',
    config_key VARCHAR(255) UNIQUE NOT NULL COMMENT '配置键',
    config_value TEXT DEFAULT NULL COMMENT '配置值',
    description VARCHAR(500) DEFAULT NULL COMMENT '配置说明',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统配置表';
