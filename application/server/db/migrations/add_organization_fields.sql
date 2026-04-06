-- 添加组织相关字段到users表
-- 体现联盟链特性

ALTER TABLE users 
ADD COLUMN organization VARCHAR(100) DEFAULT '' COMMENT '所属组织MSPID（如TaobaoMSP）',
ADD COLUMN organization_name VARCHAR(200) DEFAULT '' COMMENT '组织名称（如协和医院）',
ADD COLUMN department VARCHAR(100) DEFAULT '' COMMENT '科室',
ADD COLUMN doctor_title VARCHAR(50) DEFAULT '' COMMENT '医生职称',
ADD COLUMN doctor_license VARCHAR(100) DEFAULT '' COMMENT '医师执业证号';

-- 为现有数据设置默认值
UPDATE users SET organization = 'TaobaoMSP', organization_name = '协和医院' WHERE role = '医生' OR role = '管理员';
UPDATE users SET organization = '', organization_name = '' WHERE role = '病人';

-- 添加索引
CREATE INDEX idx_organization ON users(organization);
CREATE INDEX idx_role_organization ON users(role, organization);
