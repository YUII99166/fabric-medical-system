-- 测试查询用户数据
USE fabric_mims;

-- 查看users表结构
DESCRIBE users;

-- 查看所有用户
SELECT id, username, account_name, role, status, created_at 
FROM users 
ORDER BY created_at DESC 
LIMIT 10;

-- 统计用户数量
SELECT 
    COUNT(*) as total_users,
    SUM(CASE WHEN status = 1 THEN 1 ELSE 0 END) as active_users,
    SUM(CASE WHEN status = 0 THEN 1 ELSE 0 END) as deleted_users
FROM users;
