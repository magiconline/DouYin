# 数据库设计

## 用户信息
user_id 用户id key
username 用户名 唯一
password 密码
token
follow_count 关注总数
follower_count 粉丝总数



## 视频信息
video_id key 
author_id/user_id 外键
play_url 视频文件
cover_url 封面文件
favorite_count 点赞总数
comment_count 评论总数
upload_time 投稿时间

## 视频存储
static/日期/video_id/video.mp4
static/日期/video_id/cover.jpg

