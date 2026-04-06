# GORM 实战练习：保险理赔统计系统

## 练习目标
通过实现一个完整的理赔统计分析模块，掌握 GORM 的核心功能。

## 练习环境准备

1. 确保数据库已启动
2. 查看现有的模型定义（参考 `application/server/model/` 目录）
3. 阅读 `docs/GORM学习指南.md` 了解基础知识

## 练习任务

### 第一阶段：基础查询（难度：⭐）

#### 任务 1.1：时间范围查询
实现 `GetClaimsByDateRange` 函数，查询指定时间段内的理赔记录。

**知识点**：
- `Where` 条件查询
- `BETWEEN` 操作符
- 时间类型处理

**参考代码**：
```go
var claims []Claim
err := cs.db.Where("created_at BETWEEN ? AND ?", startDate, endDate).Find(&claims).Error
```

**测试方法**：
```go
// 查询最近 30 天的理赔
startDate := time.Now().AddDate(0, 0, -30)
endDate := time.Now()
claims, err := cs.GetClaimsByDateRange(startDate, endDate)
```

---

#### 任务 1.2：金额阈值查询
实现 `GetHighValueClaims` 函数，查找高额理赔记录。

**知识点**：
- 比较运算符
- 排序 `Order`

**提示**：可以添加按金额降序排序

---

### 第二阶段：聚合统计（难度：⭐⭐）

#### 任务 2.1：分组统计
实现 `GetClaimAmountByType` 函数，按保险类型统计理赔总额。

**知识点**：
- `Group By` 分组
- `Select` 选择字段
- 聚合函数 `SUM`

**参考代码结构**：
```go
type TypeSummary struct {
    InsuranceType string
    TotalAmount   float64
    ClaimCount    int
}

var results []TypeSummary
err := cs.db.Model(&Claim{}).
    Select("insurance_type, SUM(amount) as total_amount, COUNT(*) as claim_count").
    Group("insurance_type").
    Find(&results).Error
```

---

#### 任务 2.2：计算平均值
实现 `GetAverageClaimAmount` 函数。

**知识点**：
- `AVG` 聚合函数
- `Scan` 扫描单个值

---

### 第三阶段：关联查询（难度：⭐⭐⭐）

#### 任务 3.1：预加载关联数据
实现 `GetUserWithPolicies` 和 `GetPolicyWithClaims` 函数。

**知识点**：
- `Preload` 预加载
- 一对多关系
- N+1 查询问题

**对比示例**：
```go
// ❌ 不好的做法（N+1 查询）
var user User
db.First(&user, userID)
var policies []Policy
db.Where("user_id = ?", userID).Find(&policies)

// ✅ 好的做法（使用 Preload）
var user User
db.Preload("Policies").First(&user, userID)
```

---

#### 任务 3.2：跨表统计
实现 `GetUserClaimSummary` 函数，统计每个用户的理赔情况。

**知识点**：
- `Joins` 连接查询
- 多表关联
- 分组统计

**参考思路**：
```go
// User -> Policy -> Claim 的关联统计
db.Model(&User{}).
    Select("users.id, users.name, SUM(claims.amount) as total_amount").
    Joins("LEFT JOIN policies ON policies.user_id = users.id").
    Joins("LEFT JOIN claims ON claims.policy_id = policies.id").
    Group("users.id").
    Find(&results)
```

---

### 第四阶段：复杂查询（难度：⭐⭐⭐）

#### 任务 4.1：动态多条件查询
实现 `SearchClaims` 函数，支持多个可选条件。

**知识点**：
- 链式调用
- 条件判断
- 动态构建查询

**参考模式**：
```go
query := cs.db.Model(&Claim{})

if status != "" {
    query = query.Where("status = ?", status)
}

if minAmount > 0 {
    query = query.Where("amount >= ?", minAmount)
}

// ... 更多条件

query.Find(&claims)
```

---

#### 任务 4.2：按月统计
实现 `GetMonthlyClaimStats` 函数。

**知识点**：
- 数据库日期函数
- 复杂的 `Select` 语句

**MySQL 示例**：
```go
Select("DATE_FORMAT(created_at, '%Y-%m') as month, COUNT(*) as count, SUM(amount) as total")
```

**PostgreSQL 示例**：
```go
Select("TO_CHAR(created_at, 'YYYY-MM') as month, COUNT(*) as count, SUM(amount) as total")
```

---

### 第五阶段：事务处理（难度：⭐⭐⭐⭐）

#### 任务 5.1：批量审批
实现 `BatchApproveClaims` 函数，使用事务确保数据一致性。

**知识点**：
- `Transaction` 事务
- 错误处理
- 回滚机制

**参考代码**：
```go
err := cs.db.Transaction(func(tx *gorm.DB) error {
    // 在事务中执行多个操作
    for _, claimID := range claimIDs {
        if err := tx.Model(&Claim{}).
            Where("id = ?", claimID).
            Updates(map[string]interface{}{
                "status": "approved",
                "approver_id": approverID,
                "approved_at": time.Now(),
            }).Error; err != nil {
            return err // 返回错误会自动回滚
        }
    }
    return nil // 返回 nil 会自动提交
})
```

---

## 进阶挑战

完成基础任务后，尝试以下挑战：

1. **性能优化**：使用 `EXPLAIN` 分析查询，添加合适的索引
2. **分页查询**：实现带分页的理赔列表查询
3. **软删除**：为理赔记录添加软删除功能
4. **钩子函数**：使用 `BeforeCreate`、`AfterUpdate` 等钩子
5. **原生 SQL**：对于复杂统计，尝试使用 `Raw` 执行原生 SQL

## 测试建议

1. 为每个函数编写单元测试
2. 使用测试数据库（避免污染开发数据）
3. 测试边界情况（空结果、大数据量等）
4. 验证事务的回滚机制

## 参考资源

- 项目内文档：`docs/GORM学习指南.md`
- 现有代码：`application/server/model/` 和 `application/server/api/`
- GORM 官方文档：https://gorm.io/zh_CN/docs/

## 完成标准

- [ ] 所有 TODO 函数都已实现
- [ ] 代码能够编译通过
- [ ] 至少完成 5 个测试用例
- [ ] 能够正确处理错误情况
- [ ] 代码符合 Go 语言规范

祝你练习顺利！遇到问题随时查阅文档或寻求帮助。
