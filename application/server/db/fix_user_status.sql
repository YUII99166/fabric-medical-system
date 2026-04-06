-- 修复用户状态：将所有现有用户的status设置为1（正常）
UPDATE users SET status = 1 WHERE status IS NULL OR status = 0;

-- 验证修改
SELECT username, account_name, role, status FROM users;
