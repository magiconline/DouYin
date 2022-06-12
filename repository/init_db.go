package repository

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	mysqlUrl = "root:123456@tcp(127.0.0.1:3306)/douyin?charset=utf8mb4&parseTime=true&loc=Local"
	DB       *gorm.DB
	RDB      *redis.Client
	CTX      = context.Background()
)

// Init 初始化mysql
// 初始化redis
func Init() error {
	// 打开或创建日志文件
	gormLogFile, err := os.OpenFile("./log/gorm.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	gormLogger := logger.New(
		log.New(gormLogFile, "", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Silent,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	DB, err = gorm.Open(mysql.Open(mysqlUrl), &gorm.Config{
		PrepareStmt: true,
		Logger:      gormLogger,
	})
	if err != nil {
		return err
	}

	DB.AutoMigrate(&Star{})
	DB.AutoMigrate(&Relation{})

	RDB = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "123456",
		DB:       0,
	})

	if _, err := RDB.Ping(CTX).Result(); err != nil {
		return err
	}

	return nil
}
