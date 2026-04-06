-- 更新users表,添加缺失的字段

USE fabric_mims;

-- 添加 doctor_license 字段
ALTER TABLE users ADD COLUMN IF NOT EXISTS doctor_license VARCHAR(100) DEFAULT '' COMMENT '医师执业证号';

-- 添加 age 字段
ALTER TABLE users ADD COLUMN IF NOT EXISTS age INT DEFAULT NULL COMMENT '年龄';

-- 添加 status 字段
ALTER TABLE users ADD COLUMN IF NOT EXISTS status TINYINT DEFAULT 1 COMMENT '状态: 1-正常, 0-禁用';

-- 添加 updated_at 字段
ALTER TABLE users ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间';

-- 查看更新后的表结构
DESC users;

-- 查看现有数据
SELECT * FROM users;
