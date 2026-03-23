package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var isMysqlInit = false
var DbClient *gorm.DB

// SetupDB 初始化MySQL连接
func SetupDB(cfg *DataBase) (err error) {
	if isMysqlInit != false {
		return nil
	}
	//fmt.Println("数据库 init")
	Dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=%t&loc=%s",
		cfg.UserName,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		true,
		"Local")
	//Dsn := "root:eyFMV6kub9TVe7Ld@tcp(127.0.0.1:3316)/plastics?charset=utf8&parseTime=True&loc=Local"
	//fmt.Println(Dsn)
	//fmt.Println(cfg.Prefix)

	gormConfig := &gorm.Config{
		// 命名策略
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   cfg.Prefix, // 表名前缀
			SingularTable: true,       // 使用单数表名
		},
		// 日志配置
		Logger: GetLogger(cfg.LogLevel),
		// 禁用外键约束（可选）
		DisableForeignKeyConstraintWhenMigrating: true,
		// 时间函数
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	}

	// 建立连接

	db, gormErr := gorm.Open(mysql.Open(Dsn), gormConfig)
	if gormErr != nil {
		return fmt.Errorf("连接数据库失败: %w", gormErr)
	}

	// 获取底层sql.DB
	sqlDB, gormErr := db.DB()
	if gormErr != nil {
		return fmt.Errorf("获取底层数据库连接失败: %w", gormErr)
	}

	// 连接池配置
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.MaxLifeTime) * time.Second)

	// 测试连接
	if gormErr := sqlDB.Ping(); gormErr != nil {
		return fmt.Errorf("数据库连接测试失败: %w", gormErr)
	}

	DbClient = db
	isMysqlInit = true
	//log.Println("数据库连接成功")

	return gormErr
}

// CloseDB 关闭MySQL连接
func CloseDB(db *gorm.DB) {
	sqlDB, _ := db.DB()
	_ = sqlDB.Close()
}

func GetMysql() *gorm.DB {
	return DbClient
}

func GetLogger(logLevel int) logger.Interface {
	var level logger.LogLevel

	switch logLevel {
	case 1:
		level = logger.Silent
	case 2:
		level = logger.Error
	case 3:
		level = logger.Warn
	case 4:
		level = logger.Info
	default:
		level = logger.Info
	}

	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second, // 慢查询阈值
			LogLevel:                  level,       // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略记录未找到错误
			Colorful:                  true,        // 彩色输出
		},
	)
}
