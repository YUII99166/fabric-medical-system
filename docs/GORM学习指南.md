# GORM 学习指南

## 1. GORM 简介

GORM 是 Go 语言最流行的 ORM（对象关系映射）框架，它可以让你用面向对象的方式操作数据库，而不需要写原生 SQL。

### 为什么使用 GORM？

- **简化开发**：不需要手写 SQL 语句
- **类型安全**：编译时检查，减少运行时错误
- **自动迁移**：自动创建和更新表结构
- **关联处理**：轻松处理表之间的关系
- **钩子函数**：在操作前后执行自定义逻辑

## 2. 安装 GORM

```bash
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql
```

## 3. 基础使用

### 3.1 定义模型（Model）

```go
package model

import "gorm.io/gorm"

// User 用户模型
type User struct {
	gorm.Model              // 包含 ID, CreatedAt, UpdatedAt, DeletedAt
	Username       string   `gorm:"uniqueIndex;not null"` // 唯一索引
	Password       string   `gorm:"not null"`
	AccountName    string   `gorm:"not null"`
	Role           string   `gorm:"not null"`
	Organization   string   `gorm:"default:''"`
	OrganizationName string `gorm:"default:''"`
	Department     string   `gorm:"default:''"`
	DoctorTitle    string   `gorm:"default:''"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}
```

**gorm.Model 包含的字段**：
```go
type Model struct {
	ID        uint           `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
```

### 3.2 连接数据库

```go
package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func Init() error {
	// MySQL 连接字符串
	dsn := "user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	
	// 连接数据库
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("连接数据库失败: %v", err)
	}
	
	log.Println("数据库连接成功")
	
	// 自动迁移（创建表）
	err = DB.AutoMigrate(&User{})
	if err != nil {
		return fmt.Errorf("自动迁移失败: %v", err)
	}
	
	return nil
}
```

## 4. CRUD 操作

### 4.1 创建（Create）

```go
// 创建单条记录
func CreateUser(username, password, accountName, role string) error {
	user := User{
		Username:    username,
		Password:    password,
		AccountName: accountName,
		Role:        role,
	}
	
	result := DB.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	
	fmt.Printf("创建成功，ID: %d\n", user.ID)
	return nil
}

// 批量创建
func CreateUsers(users []User) error {
	result := DB.Create(&users)
	return result.Error
}
```

### 4.2 查询（Read）

```go
// 查询单条记录
func GetUserByID(id uint) (*User, error) {
	var user User
	result := DB.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// 根据条件查询
func GetUserByUsername(username string) (*User, error) {
	var user User
	result := DB.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// 查询所有记录
func GetAllUsers() ([]User, error) {
	var users []User
	result := DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

// 条件查询
func GetUsersByRole(role string) ([]User, error) {
	var users []User
	result := DB.Where("role = ?", role).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

// 复杂查询
func SearchUsers(keyword string, role string) ([]User, error) {
	var users []User
	query := DB.Model(&User{})
	
	if keyword != "" {
		query = query.Where("username LIKE ? OR account_name LIKE ?", 
			"%"+keyword+"%", "%"+keyword+"%")
	}
	
	if role != "" {
		query = query.Where("role = ?", role)
	}
	
	result := query.Find(&users)
	return users, result.Error
}
```

### 4.3 更新（Update）

```go
// 更新单个字段
func UpdateUserPassword(id uint, newPassword string) error {
	result := DB.Model(&User{}).Where("id = ?", id).Update("password", newPassword)
	return result.Error
}

// 更新多个字段（使用 map）
func UpdateUser(id uint, updates map[string]interface{}) error {
	result := DB.Model(&User{}).Where("id = ?", id).Updates(updates)
	return result.Error
}

// 更新多个字段（使用结构体）
func UpdateUserInfo(user *User) error {
	// 注意：零值字段不会被更新
	result := DB.Model(user).Updates(user)
	return result.Error
}

// 使用 Select 指定要更新的字段
func UpdateUserWithSelect(user *User) error {
	result := DB.Model(user).Select("account_name", "role").Updates(user)
	return result.Error
}
```

### 4.4 删除（Delete）

```go
// 软删除（推荐）- 只是标记 DeletedAt 字段
func DeleteUser(id uint) error {
	result := DB.Delete(&User{}, id)
	return result.Error
}

// 永久删除
func HardDeleteUser(id uint) error {
	result := DB.Unscoped().Delete(&User{}, id)
	return result.Error
}

// 批量删除
func DeleteUsersByRole(role string) error {
	result := DB.Where("role = ?", role).Delete(&User{})
	return result.Error
}
```

## 5. 高级查询

### 5.1 链式查询

```go
func AdvancedQuery() {
	var users []User
	
	DB.Where("role = ?", "医生").
		Where("organization = ?", "TaobaoMSP").
		Order("created_at DESC").
		Limit(10).
		Offset(0).
		Find(&users)
}
```

### 5.2 原生 SQL

```go
// 执行原生 SQL
func RawQuery() {
	var users []User
	DB.Raw("SELECT * FROM users WHERE role = ?", "医生").Scan(&users)
}

// 执行原生 SQL（无返回值）
func RawExec() {
	DB.Exec("UPDATE users SET role = ? WHERE id = ?", "管理员", 1)
}
```

### 5.3 事务

```go
func TransferExample() error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// 在事务中执行操作
		if err := tx.Create(&User{Username: "user1"}).Error; err != nil {
			return err // 返回错误会回滚
		}
		
		if err := tx.Create(&User{Username: "user2"}).Error; err != nil {
			return err
		}
		
		// 返回 nil 会提交事务
		return nil
	})
}
```

## 6. 将你的项目改造为 GORM

### 6.1 修改 db.go

```go
package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

type User struct {
	gorm.Model
	Username         string `gorm:"uniqueIndex;not null"`
	Password         string `gorm:"not null"`
	AccountName      string `gorm:"not null"`
	Role             string `gorm:"not null"`
	Organization     string `gorm:"default:''"`
	OrganizationName string `gorm:"default:''"`
	Department       string `gorm:"default:''"`
	DoctorTitle      string `gorm:"default:''"`
}

func Init() error {
	dsn := "fabric:123456@tcp(101.43.97.84:3306)/fabric_mims?charset=utf8mb4&parseTime=True&loc=Local"
	
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("连接数据库失败: %v", err)
	}
	
	log.Println("数据库连接成功")
	
	// 自动迁移
	err = DB.AutoMigrate(&User{})
	if err != nil {
		return fmt.Errorf("自动迁移失败: %v", err)
	}
	
	return nil
}

// RegisterUser 注册用户（GORM 版本）
func RegisterUser(username, password, accountName, role, org, orgName, dept, title string) error {
	user := User{
		Username:         username,
		Password:         HashPassword(password),
		AccountName:      accountName,
		Role:             role,
		Organization:     org,
		OrganizationName: orgName,
		Department:       dept,
		DoctorTitle:      title,
	}
	
	result := DB.Create(&user)
	return result.Error
}

// LoginUser 登录验证（GORM 版本）
func LoginUser(username, password string) (*User, error) {
	var user User
	hashedPassword := HashPassword(password)
	
	result := DB.Where("username = ? AND password = ?", username, hashedPassword).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("用户名或密码错误")
		}
		return nil, result.Error
	}
	
	return &user, nil
}

// GetUserByID 根据ID获取用户
func GetUserByID(id int) (*User, error) {
	var user User
	result := DB.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// GetUserByUsername 根据用户名获取用户
func GetUserByUsername(username string) (*User, error) {
	var user User
	result := DB.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// UpdateUserFull 更新用户完整信息
func UpdateUserFull(id int, accountName, role, org, orgName, dept, title string) error {
	updates := map[string]interface{}{
		"account_name":      accountName,
		"role":              role,
		"organization":      org,
		"organization_name": orgName,
		"department":        dept,
		"doctor_title":      title,
	}
	
	result := DB.Model(&User{}).Where("id = ?", id).Updates(updates)
	return result.Error
}

// DeleteUser 删除用户（软删除）
func DeleteUser(id int) error {
	result := DB.Delete(&User{}, id)
	return result.Error
}

// CountAdminsExcluding 统计除指定ID外的管理员数量
func CountAdminsExcluding(excludeID int) (int64, error) {
	var count int64
	result := DB.Model(&User{}).
		Where("role = ? AND id != ?", "管理员", excludeID).
		Count(&count)
	return count, result.Error
}
```

## 7. 常用标签（Tags）

```go
type User struct {
	ID        uint   `gorm:"primaryKey"`                    // 主键
	Username  string `gorm:"uniqueIndex;not null"`          // 唯一索引，非空
	Email     string `gorm:"type:varchar(100);unique"`      // 指定类型和唯一约束
	Age       int    `gorm:"default:18"`                    // 默认值
	Role      string `gorm:"size:50;not null"`              // 字符串长度
	CreatedAt time.Time `gorm:"autoCreateTime"`             // 自动设置创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime"`             // 自动更新时间
	DeletedAt gorm.DeletedAt `gorm:"index"`                 // 软删除
	Ignored   string `gorm:"-"`                             // 忽略该字段
}
```

## 8. 实战技巧

### 8.1 预加载（Preload）

```go
// 假设有关联关系
type User struct {
	gorm.Model
	Username string
	Posts    []Post // 一对多关系
}

type Post struct {
	gorm.Model
	Title  string
	UserID uint
}

// 预加载关联数据
var user User
DB.Preload("Posts").First(&user, 1)
```

### 8.2 钩子函数

```go
// BeforeCreate 创建前执行
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// 自动加密密码
	u.Password = HashPassword(u.Password)
	return nil
}

// AfterFind 查询后执行
func (u *User) AfterFind(tx *gorm.DB) error {
	// 可以做一些数据处理
	return nil
}
```

### 8.3 作用域（Scopes）

```go
// 定义作用域
func ActiveUsers(db *gorm.DB) *gorm.DB {
	return db.Where("deleted_at IS NULL")
}

func DoctorRole(db *gorm.DB) *gorm.DB {
	return db.Where("role = ?", "医生")
}

// 使用作用域
var users []User
DB.Scopes(ActiveUsers, DoctorRole).Find(&users)
```

## 9. 性能优化

### 9.1 批量插入

```go
users := []User{
	{Username: "user1"},
	{Username: "user2"},
	{Username: "user3"},
}

// 批量插入（一次性插入）
DB.Create(&users)

// 分批插入（每批100条）
DB.CreateInBatches(users, 100)
```

### 9.2 选择字段

```go
// 只查询需要的字段
var users []User
DB.Select("id", "username", "role").Find(&users)

// 排除某些字段
DB.Omit("password").Find(&users)
```

### 9.3 索引

```go
type User struct {
	Username string `gorm:"index"`           // 普通索引
	Email    string `gorm:"uniqueIndex"`     // 唯一索引
	Age      int    `gorm:"index:idx_age"`   // 命名索引
}
```

## 10. 调试技巧

```go
// 打印 SQL 语句
DB.Debug().Where("username = ?", "test").First(&user)

// 获取执行的 SQL
sql := DB.ToSQL(func(tx *gorm.DB) *gorm.DB {
	return tx.Where("username = ?", "test").First(&user)
})
fmt.Println(sql)
```

## 11. 总结

GORM 的优势：
- ✅ 代码更简洁，不需要手写 SQL
- ✅ 类型安全，编译时检查
- ✅ 自动迁移，维护表结构
- ✅ 丰富的查询方法
- ✅ 支持事务、钩子、关联等高级特性

适用场景：
- 中小型项目
- 快速开发
- 标准的 CRUD 操作

不适用场景：
- 复杂的 SQL 查询（可以用原生 SQL）
- 对性能要求极高的场景
- 需要精细控制 SQL 的场景

## 12. 参考资源

- 官方文档：https://gorm.io/zh_CN/docs/
- GitHub：https://github.com/go-gorm/gorm
