package db

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
)

type User struct {
	ID               int
	Username         string
	Password         string
	AccountName      string
	Role             string
	Organization     string // 所属组织MSPID
	OrganizationName string // 组织名称
	Department       string // 科室
	DoctorTitle      string // 医生职称
	DoctorLicense    string // 医师执业证号
	Status           int    // 状态：1-正常，0-禁用
	CreatedAt        string
	UpdatedAt        string
}

// HashPassword 对密码进行哈希处理
func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// RegisterUser 注册新用户
func RegisterUser(username, password, accountName, role, organization, organizationName, department, doctorTitle string) error {
	// 检查用户名是否已存在
	var exists bool
	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", username).Scan(&exists)
	if err != nil {
		return fmt.Errorf("检查用户名失败: %v", err)
	}
	if exists {
		return errors.New("用户名已存在")
	}

	// 对密码进行哈希处理
	hashedPassword := HashPassword(password)

	// 插入新用户
	_, err = DB.Exec(
		"INSERT INTO users (username, password, account_name, role, organization, organization_name, department, doctor_title) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		username, hashedPassword, accountName, role, organization, organizationName, department, doctorTitle,
	)
	if err != nil {
		return fmt.Errorf("注册用户失败: %v", err)
	}

	return nil
}

// LoginUser 用户登录验证（只允许status=1的用户登录）
func LoginUser(username, password string) (*User, error) {
	user := &User{}
	
	// 先检查用户是否存在
	var existingUsername string
	err := DB.QueryRow("SELECT username FROM users WHERE username = ?", username).Scan(&existingUsername)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("用户不存在，请先注册")
		}
		return nil, fmt.Errorf("查询用户失败: %v", err)
	}
	
	// 用户存在，验证密码
	hashedPassword := HashPassword(password)
	var status int

	err = DB.QueryRow(
		"SELECT id, username, account_name, role, organization, organization_name, department, doctor_title, status, created_at, updated_at FROM users WHERE username = ? AND password = ?",
		username, hashedPassword,
	).Scan(&user.ID, &user.Username, &user.AccountName, &user.Role, &user.Organization, &user.OrganizationName, &user.Department, &user.DoctorTitle, &status, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("密码错误，请重新输入")
		}
		return nil, fmt.Errorf("查询用户失败: %v", err)
	}

	// 检查用户状态
	if status == 0 {
		return nil, errors.New("该账户已被停用，请联系管理员")
	}

	return user, nil
}

// GetUserByUsername 根据用户名获取用户信息
func GetUserByUsername(username string) (*User, error) {
	user := &User{}
	err := DB.QueryRow(
		"SELECT id, username, account_name, role, organization, organization_name, department, doctor_title, status, created_at, updated_at FROM users WHERE username = ?",
		username,
	).Scan(&user.ID, &user.Username, &user.AccountName, &user.Role, &user.Organization, &user.OrganizationName, &user.Department, &user.DoctorTitle, &user.Status, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %v", err)
	}

	return user, nil
}

// GetUserByID 根据ID获取用户信息（包括已删除的用户）
func GetUserByID(id int) (*User, error) {
	user := &User{}
	var status int
	err := DB.QueryRow(
		"SELECT id, username, account_name, role, status, created_at, updated_at FROM users WHERE id = ?",
		id,
	).Scan(&user.ID, &user.Username, &user.AccountName, &user.Role, &status, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %v", err)
	}

	return user, nil
}

// GetUserDetail 根据ID获取用户详细信息
func GetUserDetail(id int) (*User, error) {
	user := &User{}
	err := DB.QueryRow(
		"SELECT id, username, account_name, role, organization, organization_name, department, doctor_title, status, created_at, updated_at FROM users WHERE id = ?",
		id,
	).Scan(&user.ID, &user.Username, &user.AccountName, &user.Role, &user.Organization, &user.OrganizationName, &user.Department, &user.DoctorTitle, &user.Status, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %v", err)
	}

	return user, nil
}

// UpdateUser 更新用户信息
func UpdateUser(id int, accountName, role string) error {
	_, err := DB.Exec(
		"UPDATE users SET account_name = ?, role = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		accountName, role, id,
	)
	if err != nil {
		return fmt.Errorf("更新用户失败: %v", err)
	}
	return nil
}

// UpdateUserFull 更新用户完整信息
func UpdateUserFull(id int, accountName, role, organization, organizationName, department, doctorTitle string) error {
	_, err := DB.Exec(
		"UPDATE users SET account_name = ?, role = ?, organization = ?, organization_name = ?, department = ?, doctor_title = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		accountName, role, organization, organizationName, department, doctorTitle, id,
	)
	if err != nil {
		return fmt.Errorf("更新用户失败: %v", err)
	}
	return nil
}

// UpdateUserPassword 更新用户密码
func UpdateUserPassword(id int, newPassword string) error {
	hashedPassword := HashPassword(newPassword)
	_, err := DB.Exec(
		"UPDATE users SET password = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		hashedPassword, id,
	)
	if err != nil {
		return fmt.Errorf("更新密码失败: %v", err)
	}
	return nil
}

// DeleteUser 软删除用户（将status设为0）
func DeleteUser(id int) error {
	_, err := DB.Exec("UPDATE users SET status = 0, updated_at = CURRENT_TIMESTAMP WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("删除用户失败: %v", err)
	}
	return nil
}

// RestoreUser 恢复已删除的用户（将status设为1）
func RestoreUser(id int) error {
	_, err := DB.Exec("UPDATE users SET status = 1, updated_at = CURRENT_TIMESTAMP WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("恢复用户失败: %v", err)
	}
	return nil
}

// HardDeleteUser 物理删除用户（真正从数据库删除，慎用）
func HardDeleteUser(id int) error {
	_, err := DB.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("物理删除用户失败: %v", err)
	}
	return nil
}

// CountAdmins 统计管理员数量
func CountAdmins() (int, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users WHERE role = '管理员'").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("统计管理员数量失败: %v", err)
	}
	return count, nil
}

// CountAdminsExcluding 统计除指定用户外的管理员数量（只统计正常状态的管理员）
func CountAdminsExcluding(excludeID int) (int, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users WHERE role = '管理员' AND status = 1 AND id != ?", excludeID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("统计管理员数量失败: %v", err)
	}
	return count, nil
}

// ListAllUsers 获取所有用户（默认只显示正常用户，可选显示已删除用户）
func ListAllUsers() ([]User, error) {
	return ListUsersWithStatus(1) // 默认只显示正常用户
}

// ListUsersWithStatus 根据状态获取用户列表
// status: -1=全部, 0=已删除, 1=正常
func ListUsersWithStatus(status int) ([]User, error) {
	var query string
	var args []interface{}
	
	if status == -1 {
		// 查询所有用户
		query = "SELECT id, username, account_name, role, organization, organization_name, department, doctor_title, status, created_at, updated_at FROM users ORDER BY status DESC, created_at DESC"
	} else {
		// 查询指定状态的用户
		query = "SELECT id, username, account_name, role, organization, organization_name, department, doctor_title, status, created_at, updated_at FROM users WHERE status = ? ORDER BY created_at DESC"
		args = append(args, status)
	}
	
	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("查询用户列表失败: %v", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		user := User{}
		err := rows.Scan(&user.ID, &user.Username, &user.AccountName, &user.Role, 
			&user.Organization, &user.OrganizationName, &user.Department, &user.DoctorTitle,
			&user.Status, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("扫描用户数据失败: %v", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历用户数据失败: %v", err)
	}

	return users, nil
}

// GetAllUsers 获取所有用户（别名）
func GetAllUsers() ([]User, error) {
	return ListAllUsers()
}
