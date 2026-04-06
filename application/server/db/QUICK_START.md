# 快速启动指南

## 1. 更新数据库表结构

在启动服务器之前，需要先更新`audit_logs`表结构。

### 方法一：使用MySQL命令行

```bash
# 进入MySQL容器
docker exec -it mysql bash

# 登录MySQL
mysql -u root -p123456

# 执行以下SQL
USE fabric;

ALTER TABLE audit_logs 
ADD COLUMN organization VARCHAR(100) DEFAULT NULL COMMENT '所属组织MSPID' AFTER user_role,
ADD COLUMN organization_name VARCHAR(200) DEFAULT NULL COMMENT '组织名称' AFTER organization,
ADD COLUMN department VARCHAR(100) DEFAULT NULL COMMENT '科室' AFTER organization_name,
ADD COLUMN related_user_id VARCHAR(100) DEFAULT NULL COMMENT '关联用户ID（如患者ID）' AFTER resource_id,
ADD COLUMN related_user_name VARCHAR(255) DEFAULT NULL COMMENT '关联用户名（如患者姓名）' AFTER related_user_id;

CREATE INDEX idx_audit_organization ON audit_logs(organization);
CREATE INDEX idx_audit_org_name ON audit_logs(organization_name);

# 验证表结构
DESCRIBE audit_logs;

# 退出
exit
exit
```

### 方法二：使用SQL文件

```bash
docker exec -i mysql mysql -u root -p123456 fabric < update_audit_logs.sql
```

## 2. 启动服务器

```bash
cd fabric-mims-main/application/server
go run main.go
```

## 3. 访问活动监控页面

1. 登录管理员账户
2. 访问：http://localhost:8080/#/account/activity-monitor
3. 或通过侧边栏：监管中心 > 活动监控

## 注意事项

- 目前活动监控页面会从`audit_logs`表读取数据
- 如果表中没有数据，页面会显示为空
- 需要在各个业务API中集成审计日志记录功能（参考`docs/审计日志集成指南.md`）
- 暂时可以先查看页面效果，后续再集成日志记录功能
