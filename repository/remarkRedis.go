package repository

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strconv"
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
	ID      	  uint64              `json:"id"`
	User          UserResponse        `json:"user"`
	Content  	  string     	     `json:"content"`
	CreateDate    time.Time          `json:"create_date"`
}

//将reamrk信息更新到redis
func UpdateReamrkToRedis(remark Remark)  {

}



//从redis缓存获取用户信息
//key:user_{user_id}   value:map{"user_name", "follow_count", "follower_count"}
func GetUserInfoFromRedis(userid uint64)(*UserResponse,error)  {
	fmt.Println("GetUserInfoFromRedis")
	user, _ := UserInfo(userid)
	return &UserResponse{ID: uint64(user.UserId), Name:user.UserName, FollowCount: 0,FollowerCount: 0,IsFollow: false},nil
}

//从redis缓存获取视频的当前评论总数
//key :video_{video_id} value:map{"user_id","play_url","cover_url","favorite_count","comment_count","upload_time","title"}
//返回值：count   未完成
func GetCurrentVideoCommentCountFromRedis(videoid uint64) (int64,error) {

	return 0,nil
}



//从redis 缓存获取评论信息
// key :com_video_{video_id}  value:"comment_id",
// key :comment_{comment_id}  valuemap{"user_id","action_type","comment_text","create_time"}
func GetAllCommentFromRedis(videoid uint64,token string)(*[]map[string]interface{},error)  {
	fmt.Println("GetAllCommentFromRedis")
	var result Remark
	var commentList []map[string]interface{}
	key := fmt.Sprintf("com_video_%v", videoid)
	rdbResult, err := RDB.HGetAll(CTX, key).Result()
	if err == nil && len(rdbResult) != 0 {
		// 缓存找到
		for _,v := range rdbResult{
			commentId, _ := strconv.Atoi(v)
			key1 := fmt.Sprintf("comment_%v",commentId)
			rdbResult1,err :=RDB.HGetAll(CTX,key1).Result()
			if err == nil && len(rdbResult1) != 0{
				//根据commentid查到了对应的值
				userid := rdbResult1["user_id"]
				userid1, _ :=strconv.Atoi(userid)
				actiontype := rdbResult1["action_type"]
				actiontype1, _ :=strconv.Atoi(actiontype)
				commenttext := rdbResult1["comment_text"]
				createtime := rdbResult1["create_time"]
				currenttime:=time.Now().String()
				loc, _ := time.LoadLocation("Local")  //CST
				createtime1, _ := time.ParseInLocation(currenttime,createtime,loc)

				result = Remark{
					Id : 		 int64(commentId),
					VideoId:     int64(videoid),
					UserId:      int64(userid1),
					ActionType:  int32(actiontype1),
					CommentText: commenttext,
					CreateTime:  createtime1,
				}
				//将结构体序列化成map
				data, _ := json.Marshal(&result)
				m := make(map[string]interface{})
				json.Unmarshal(data, &m)
				commentList = append(commentList,m)
			}
		}
	}
	fmt.Println("缓存拿取结果")
	return &commentList,nil
}


//从mysql获取用户信息,返回值*UserResponse
func GetUserInfoFromMysql(userid uint64)(*User,error)  {
	var result User
	err := DB.Table("user").Where(User{UserId: uint64(int64(userid))}).Select("user_name", "follow_count", "follower_count").Take(&result).Error
	return &result, err
}


//从mysql读取评论信息数据，返回值：*[]map[string]interface{}
func GetRemarkFromMysql(videoId uint64)(*[]map[string]interface{},error)  {
	fmt.Println("GetRemarkFromMysql")
	//执行sql查询语句,根据video查询
	var commentList []map[string]interface{}
	err := DB.Table("remark").Where("video_id = ?",videoId).Order("create_time desc").Find(&commentList).Error
	return &commentList,err

}



//将mysql中的评论查询记录更新到redies， 传入*[]map[string]interface{}
// key :com_video_{video_id}  value:"comment_id",
// key :comment_{comment_id}  valuemap{"user_id","action_type","comment_text","create_time"}
func CachRemark2Redis(commentlist *[]map[string]interface{})  {

	//清除原有缓冲：del remark
	for i := range *commentlist {
		videoid := (*commentlist)[i]["video_id"]
		commentid := (*commentlist)[i]["comment_id"]
		key3 := fmt.Sprintf("com_video_%v", videoid)
		key4 := fmt.Sprintf("comment_%v", commentid)
		RDB.Do(CTX,"del",key3)
		RDB.Do(CTX,"del",key4)
	}

	for i := range *commentlist{
		videoid :=(*commentlist)[i]["video_id"]
		commentid :=(*commentlist)[i]["comment_id"]
		userid :=(*commentlist)[i]["video_id"]
		actiontype := (*commentlist)[i]["action_type"]
		commenttext :=(*commentlist)[i]["comment_text"]
		createtime := (*commentlist)[i]["create_time"]
		//RDB.Do(CTX,"hkeys",videoid)
		key1 := fmt.Sprintf("com_video_%v", videoid)
		key2 := fmt.Sprintf("comment_%v", commentid)
		RDB.Do(CTX,"hmset",key1,"comment_id",commentid)
		RDB.Do(CTX,"hmset",key2,"user_id",userid,"action_type",actiontype,"comment_text",commenttext,"create_time",createtime)
	}
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
