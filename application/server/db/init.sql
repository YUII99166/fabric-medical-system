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

-- 注意：生产环境中不应该插入默认测试账户
-- 建议在系统部署后通过注册接口创建管理员账号
-- 如需初始化管理员账号，请取消下面的注释并修改密码

-- 示例：创建初始管理员账号（请修改密码）
-- INSERT IGNORE INTO users (username, password, account_name, role) VALUES
-- ('admin', 'YOUR_HASHED_PASSWORD_HERE', '系统管理员', '管理员');

-- 密码哈希说明：
-- 系统使用 bcrypt 算法加密密码
-- 请使用后端提供的密码加密工具生成密码哈希值
-- 不要在生产环境中使用明文密码或简单密码
