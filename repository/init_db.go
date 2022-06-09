package repository

import (
	"context"

	"github.com/go-redis/redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	mysqlUrl = "root:1234@tcp(127.0.0.1:3306)/douyin?charset=utf8mb4&parseTime=true&loc=Local"
	DB       *gorm.DB
	RDB      *redis.Client
	CTX      = context.Background()
)

// 初始化mysql
// 初始化redis
func Init() error {
	var _ error
	DB, _ = gorm.Open(mysql.Open(mysqlUrl), &gorm.Config{
		PrepareStmt: true,
	})
	//if err != nil {
	//	return err
	//}
	//DB.AutoMigrate(&Star{})
	//// DB.AutoMigrate(&User{})
	//DB.AutoMigrate(&Relation{})
	//
	//RDB = redis.NewClient(&redis.Options{
	//	Addr:     "127.0.0.1:6379",
	//	Password: "1234",
	//	DB:       0,
	//})
	//
	//if _, err := RDB.Ping(CTX).Result(); err != nil {
	//	return err
	//}

	return nil
}
