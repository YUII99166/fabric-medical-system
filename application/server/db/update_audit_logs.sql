-- ============================================
-- 更新 audit_logs 表，添加组织信息字段
-- ============================================

USE fabric;

-- 添加组织相关字段
ALTER TABLE audit_logs 
ADD COLUMN organization VARCHAR(100) DEFAULT NULL COMMENT '所属组织MSPID' AFTER user_role,
ADD COLUMN organization_name VARCHAR(200) DEFAULT NULL COMMENT '组织名称' AFTER organization,
ADD COLUMN department VARCHAR(100) DEFAULT NULL COMMENT '科室' AFTER organization_name,
ADD COLUMN related_user_id VARCHAR(100) DEFAULT NULL COMMENT '关联用户ID（如患者ID）' AFTER resource_id,
ADD COLUMN related_user_name VARCHAR(255) DEFAULT NULL COMMENT '关联用户名（如患者姓名）' AFTER related_user_id;

-- 添加索引
CREATE INDEX idx_audit_organization ON audit_logs(organization);
CREATE INDEX idx_audit_org_name ON audit_logs(organization_name);

-- 查看更新后的表结构
DESCRIBE audit_logs;
