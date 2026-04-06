-- ============================================
-- 病历访问日志表 (prescription_access_logs)
-- 用于隐私保护和溯源
-- ============================================

USE fabric_mims;

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
    user_agent TEXT COMMENT '用户代理',
    tx_id VARCHAR(100) DEFAULT NULL COMMENT '区块链交易ID',
    accessed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '访问时间',
    INDEX idx_access_prescription (prescription_id),
    INDEX idx_access_patient (patient_id),
    INDEX idx_access_accessor (accessor_id),
    INDEX idx_access_time (accessed_at),
    INDEX idx_access_type (access_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='病历访问日志表';

-- 查看表结构
DESCRIBE prescription_access_logs;
