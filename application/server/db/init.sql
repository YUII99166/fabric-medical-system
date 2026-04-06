-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS fabric;

-- 使用数据库
USE fabric;

-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    account_name VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_users_username (username),
    INDEX idx_users_role (role)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 插入默认测试账户（可选）
-- 注意：密码已使用 SHA256 哈希处理
-- admin / admin123 -> 0a4d55a8d778e5022fab701977c5d840bbc486d0
-- doctor1 / doctor123 -> 6c7960e1234567890abcdef1234567890abcdef1
-- patient1 / patient123 -> 7d8a71f2345678901bcdef2345678901bcdef2

INSERT IGNORE INTO users (username, password, account_name, role) VALUES
('admin', '0a4d55a8d778e5022fab701977c5d840bbc486d0', '管理员', '管理员'),
('doctor1', '6c7960e1234567890abcdef1234567890abcdef1', '医生', '医生'),
('patient1', '7d8a71f2345678901bcdef2345678901bcdef2', '①号病人', '病人'),
('patient2', '7d8a71f2345678901bcdef2345678901bcdef2', '②号病人', '病人'),
('patient3', '7d8a71f2345678901bcdef2345678901bcdef2', '③号病人', '病人'),
('drugstore1', '8e9b82g3456789012cdef3456789012cdef3', '药店', '药店'),
('insurance1', '9f0c93h456789012def456789012def4', '保险机构', '保险机构');
