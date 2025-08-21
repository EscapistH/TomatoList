// database/database.go
package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"TomatoList/models"
)

// DB 全局数据库实例
var DB *gorm.DB

// InitDatabase 初始化数据库连接
// 与Python对比：
// # # 使用SQLAlchemy初始化数据库
// # from sqlalchemy import create_engine
// # from sqlalchemy.orm import sessionmaker
// #
// # engine = create_engine(DATABASE_URL)
// # SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)
// #
// # # 创建表
// # Base.metadata.create_all(bind=engine)
func InitDatabase() {
	var err error
	var dialect gorm.Dialector

	// PostgreSQL连接
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", "127.0.0.1", "postgres", "postgresql-pwd", "tomato-list", "5432")
	dialect = postgres.Open(dsn)

	// 配置GORM日志
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // 慢SQL阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略记录不存在的错误
			Colorful:                  true,        // 彩色打印
		},
	)

	// 连接数据库
	DB, err = gorm.Open(dialect, &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 自动迁移模型
	err = DB.AutoMigrate(
		&models.User{},
		&models.Task{},
		&models.Pomodoro{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database connected successfully")
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}
