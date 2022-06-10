package repository

import (
	"bytes"
	"encoding/gob"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
)

//redis key-value
//key : remark:videoid value: {}
//key : registr:userid value
type UserInfomation struct {
	UserId        int64  `gorm:"column:user_id"`
	UserName      string `gorm:"column:user_name"`
	FollowCount   int    `gorm:"column:follow_count"`
	FollowerCount int    `gorm:"column:follower_count"`
}

//序列化
func serialize(value interface{}) ([]byte, error) {
	if bytes, ok := value.([]byte); ok {
		return bytes, nil
	}

	switch v := reflect.ValueOf(value); v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return []byte(strconv.FormatInt(v.Int(), 10)), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return []byte(strconv.FormatUint(v.Uint(), 10)), nil
	}

	var b bytes.Buffer
	encoder := gob.NewEncoder(&b)
	if err := encoder.Encode(value); err != nil { //编码
		return nil, err
	}
	return b.Bytes(), nil
}

//反序列化
func deserialize(byt []byte, ptr interface{}) (err error) {
	if bytes, ok := ptr.(*[]byte); ok {
		*bytes = byt

		return nil
	}

	if v := reflect.ValueOf(ptr); v.Kind() == reflect.Ptr { // 通过反射得到ptr类型，判断ptr是指针类型
		switch p := v.Elem(); p.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64: //符号整型
			var i int64
			i, err = strconv.ParseInt(string(byt), 10, 64)
			if err != nil {
				return err
			} else {
				p.SetInt(i)
			}
			return nil

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64: //无符号整型
			var i uint64
			i, err = strconv.ParseUint(string(byt), 10, 64)
			if err != nil {
				return err
			} else {
				p.SetUint(i)
			}
			return nil
		}
	}

	b := bytes.NewBuffer(byt)
	decoder := gob.NewDecoder(b)
	if err = decoder.Decode(ptr); err != nil { //解码
		return err
	}
	return nil
}

func FindUserFromRedis(c *gin.Context) (string, error) {
	userID := c.Query("user_id")
	return RDB.Get(c, userID).Result()

}
