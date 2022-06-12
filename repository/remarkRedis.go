package repository

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"reflect"
	"time"
)


type UserResponse struct {
	ID            uint64 `json:"id"`
	Name          string `json:"name"`
	FollowCount   uint64 `json:"follow_count"`
	FollowerCount uint64 `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

type CommentResponse struct {
	ID      	  uint64         `json:"id"`
	User        UserResponse   `json:"user"`
	Content  string     	 `json:"content"`
	CreateDate    time.Time          `json:"create_date"`
}

//从redis缓存获取用户信息，返回值，user 的json格式
func GetUserInfoFromRedis(userid uint64)(*UserResponse,error)  {
	fmt.Println("GetUserInfoFromRedis")


	return nil,nil
}




//从redis 缓存获取评论信息，返回值：评论信息字符串，用[]string返回
func GetAllCommentFromRedis()(strs []string)  {
	fmt.Println("GetAllCommentFromRedis")
	reply := RDB.Do(CTX, "lrange", "remark", "0", "-1")
	//strs, _ = redis.Strings(reply)
	fmt.Println("缓存拿取结果",reply)

	return strs
}



//从mysql获取用户信息,返回值*UserResponse
func GetUserInfoFromMysql(userid uint64)(*User,error)  {
	var result User
	err := DB.Table("user").Where(User{UserId: int64(userid)}).Select("user_name", "follow_count", "follower_count").Take(&result).Error
	return &result, err
}


//从mysql读取数据，返回值：remark []remark
func GetRemarkFromMysql(videoId uint64)(*[]map[string]interface{},error)  {
	fmt.Println("GetRemarkFromMysql")
	//执行sql查询语句,根据video查询
	var commentList []map[string]interface{}
	err := DB.Table("remark").Where("video_id = ?",videoId).Order("create_time desc").Find(&commentList).Error
	return &commentList,err

}

//结构体转为map
func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

//缓存查询结果到redies
//传入remark的map,将其转换为key-value
//key : videoId  value :remark []map[string]interface{}
func CachRemark2Redis(commentlist *[]map[string]interface{})  {

	//清除原有缓冲：del remark
	RDB.Do(CTX,"del","videoId")
	//var remark []map[string]interface{}
	//遍历remark中的每一个remark对象,存入key value
	//

	fmt.Println(commentlist)
	for i := range *commentlist{

		videoid :=(*commentlist)[i]["video_id"]
		//RDB.Do(CTX,"hdel",videoid)
		commentid :=(*commentlist)[i]["comment_id"]
		userid :=(*commentlist)[i]["video_id"]
		actiontype := (*commentlist)[i]["action_type"]
		commenttext :=(*commentlist)[i]["comment_text"]
		createtime := (*commentlist)[i]["create_time"]
		RDB.Do(CTX,"hmset",videoid,"commentId",commentid,"userId",userid,"actionType",actiontype,"commentText",commenttext,"createTime",createtime)
		//RDB.Do(CTX,"hkeys",videoid)
		fmt.Println(RDB.Do(CTX,"hvals",videoid))
		//RDB.Do(CTX,"hmset","")
		//tmpdata := Struct2Map(v)
		//str, err := json.Marshal(tmpdata)
		//HandleError(err,"json to str error")
		//remarkStr := string(str)
		//
		////执行redis命令：rpush remark xxx，把每一个personStr存入列表remark
		////_, err = conn.Do("rpush", "remark", remarkStr)
		//HandleError(err, "@ rpush people "+ remarkStr)
	}


	//执行redis命令：expire remark 20，使remark在20秒后过期
	//_,err := conn.Do("expire", "remark",20)
	//HandleError(err, "@ expire remark 60")
	fmt.Println("缓存remark成功")
}


/*
错误处理函数
参数：传入错误、出错的场景
只要有错误，就打印错误并暴力退出程序
*/
func HandleError(err error, when string) {
	if err != nil {
		fmt.Println(err, when)

		//暴力结束程序
		os.Exit(1)
	}
}
