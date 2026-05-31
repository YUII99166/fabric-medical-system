package db

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
	"os"
	"time"
)

var (
	DB   *sql.DB
	once sync.Once
)

type DBConfig struct {
	DB struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
	} `yaml:"db"`
}

// Init 初始化数据库连接
func Init() error {
	var err error
	once.Do(func() {
		// 读取配置文件
		config, err := loadConfig()
		if err != nil {
			log.Printf("加载数据库配置失败: %v", err)
			return
		}

		// 构建 MySQL 连接字符串
		// charset/collation 在 DSN 中设置，驱动会对每个新连接自动应用
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=True&loc=Local",
			config.DB.User,
			config.DB.Password,
			config.DB.Host,
			config.DB.Port,
			config.DB.DBName,
		)

		log.Printf("正在连接数据库: %s:%d/%s", config.DB.Host, config.DB.Port, config.DB.DBName)

		// 连接数据库
		DB, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Printf("打开数据库连接失败: %v", err)
			return
		}

		// 设置连接池参数，确保连接定期刷新
		DB.SetMaxOpenConns(25)
		DB.SetMaxIdleConns(5)
		DB.SetConnMaxLifetime(5 * time.Minute)

		// 测试连接
		err = DB.Ping()
		if err != nil {
			log.Printf("数据库连接失败: %v", err)
			return
		}

		log.Println("数据库连接成功")

		// 初始化表
		err = initTables()
		if err != nil {
			log.Printf("初始化表失败: %v", err)
			return
		}
	})
	return err
}

// loadConfig 加载数据库配置
func loadConfig() (*DBConfig, error) {
	config := &DBConfig{}
	data, err := os.ReadFile("db/config.yaml")
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// initTables 初始化数据库表
func initTables() error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		account_name VARCHAR(255) NOT NULL,
		role VARCHAR(50) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`
	_, err := DB.Exec(createTableSQL)
	if err != nil {
		return err
	}
	log.Println("用户表初始化成功")

	// 创建病历访问日志表
	createAccessLogsTableSQL := `
	CREATE TABLE IF NOT EXISTS prescription_access_logs (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		log_id VARCHAR(100) UNIQUE NOT NULL,
		prescription_id VARCHAR(100) NOT NULL,
		prescription_no VARCHAR(100) NOT NULL,
		patient_id VARCHAR(100) NOT NULL,
		patient_name VARCHAR(100) NOT NULL,
		accessor_id VARCHAR(100) NOT NULL,
		accessor_name VARCHAR(100) NOT NULL,
		accessor_role VARCHAR(50) NOT NULL,
		accessor_organization VARCHAR(100) DEFAULT NULL,
		accessor_organization_name VARCHAR(200) DEFAULT NULL,
		access_type VARCHAR(50) NOT NULL,
		access_reason VARCHAR(500) DEFAULT NULL,
		ip_address VARCHAR(50) DEFAULT NULL,
		user_agent TEXT,
		tx_id VARCHAR(100) DEFAULT NULL,
		accessed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		INDEX idx_access_prescription (prescription_id),
		INDEX idx_access_patient (patient_id),
		INDEX idx_access_accessor (accessor_id),
		INDEX idx_access_time (accessed_at),
		INDEX idx_access_type (access_type)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`
	_, err = DB.Exec(createAccessLogsTableSQL)
	if err != nil {
		log.Printf("创建病历访问日志表失败: %v", err)
		return err
	}
	log.Println("病历访问日志表初始化成功")

	return nil
}

// Close 关闭数据库连接
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
