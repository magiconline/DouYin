package repository

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	mysqlUrl = "root:123456@tcp(127.0.0.1:3306)/douyin"
	db       *gorm.DB
)

func Init() error {
	var err error
	db, err = gorm.Open(mysql.Open(mysqlUrl), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
