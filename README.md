# 项目介绍
字节跳动青训营极简版抖音项目后端实现。  
演示视频 https://github.com/magiconline/DouYin/blob/main/demo.mp4

# 项目功能
https://www.apifox.cn/apidoc/shared-8cc50618-0da6-4d5e-a398-76f3b8f766c5/api-18345145
- 基础接口
- 扩展接口1
- 扩展接口2

# 项目依赖
- GO 1.18
  - GORM
  - Gin
  - JWT
  - uuid
  - redis
  - ffmpeg-go
- MySQL
- Redis
- ffmpeg


# 项目架构
```
|-- DouYin
    |-- .gitignore
    |-- README.md
    |-- cover.html              // 覆盖率报告
    |-- demo.mp4                // 演示视频
    |-- docker-compose.yml      // mysql与redis容器
    |-- go.mod
    |-- go.sum
    |-- init_mysql.sql          // mysql建库建表
    |-- redis.conf              // redis配置文件
    |-- server.go               // 项目启动文件
    |-- server_test.go          // 项目测试文件
    |-- test.sh                 // 覆盖率测试脚本
    |-- controller              // 控制层
    |   |-- relation.go
    |   |-- remark.go
    |   |-- star.go
    |   |-- user.go
    |   |-- video.go
    |-- db                      // docker容器存储
    |   |-- mysql
    |   |-- redis
    |-- doc                     // 项目记录
    |   |-- db_wzw.md
    |   |-- db_zc.md
    |   |-- note.md
    |   |-- 评论表remark.md
    |-- log                     // 日志存储
    |   |-- gin.log
    |   |-- gorm.log
    |   |-- log.log
    |-- logger                  // 配置自定义log
    |   |-- logger.go
    |-- repository              // 初始化、读写mysql与redis
    |   |-- init_db.go
    |   |-- redis.go
    |   |-- relation.go
    |   |-- remark.go
    |   |-- remarkRedis.go
    |   |-- star.go
    |   |-- user.go
    |   |-- video.go
    |-- service                 // 通用逻辑层
    |   |-- jwt.go
    |   |-- relation.go
    |   |-- remark.go
    |   |-- star.go
    |   |-- user.go
    |   |-- video.go
    |-- static                  // 存储视频、封面文件
        |-- 2022
            |-- 05
            |   |-- 30
            |-- 06
                |-- 08
                |-- 12
                |-- 13
            ......
```

# 数据库表结构
```
video 表
+----------------+-----------+-----------------+------------+----------------+
| COLUMN_NAME    | DATA_TYPE | COLUMN_TYPE     | COLUMN_KEY | EXTRA          |
+----------------+-----------+-----------------+------------+----------------+
| video_id       | bigint    | bigint unsigned | PRI        | auto_increment |
| user_id        | bigint    | bigint unsigned | MUL        |                |
| play_url       | text      | char(59)        |            |                |
| cover_url      | text      | char(59)        |            |                |
| favorite_count | int       | int unsigned    |            |                |
| comment_count  | int       | int unsigned    |            |                |
| upload_time    | bigint    | bigint unsigned |            |                |
| title          | text      | text            |            |                |
+----------------+-----------+-----------------+------------+----------------+

user 表
+----------------+-----------+-----------------+------------+----------------+
| COLUMN_NAME    | DATA_TYPE | COLUMN_TYPE     | COLUMN_KEY | EXTRA          |
+----------------+-----------+-----------------+------------+----------------+
| user_id        | bigint    | bigint unsigned | PRI        | auto_increment |
| user_name      | varchar   | varchar(45)     |            |                |
| password       | varchar   | varchar(45)     |            |                |
| follow_count   | int       | int unsigned    |            |                |
| follower_count | int       | int unsigned    |            |                |
+----------------+-----------+-----------------+------------+----------------+

star 表
+-------------+-----------+-----------------+------------+----------------+
| COLUMN_NAME | DATA_TYPE | COLUMN_TYPE     | COLUMN_KEY | EXTRA          |
+-------------+-----------+-----------------+------------+----------------+
| id          | bigint    | bigint unsigned | PRI        | auto_increment |
| user_id     | bigint    | bigint unsigned | MUL        |                |
| video_id    | bigint    | bigint unsigned |            |                |
+-------------+-----------+-----------------+------------+----------------+

remark 表
+--------------+-----------+--------------+------------+----------------+
| COLUMN_NAME  | DATA_TYPE | COLUMN_TYPE  | COLUMN_KEY | EXTRA          |
+--------------+-----------+--------------+------------+----------------+
| comment_id   | int       | int          | PRI        | auto_increment |
| video_id     | int       | int          | MUL        |                |
| user_id      | int       | int          |            |                |
| action_type  | int       | int          |            |                |
| comment_text | varchar   | varchar(255) |            |                |
| create_time  | datetime  | datetime     |            |                |
+--------------+-----------+--------------+------------+----------------+

relations 表
+-------------+-----------+-----------------+------------+----------------+
| COLUMN_NAME | DATA_TYPE | COLUMN_TYPE     | COLUMN_KEY | EXTRA          |
+-------------+-----------+-----------------+------------+----------------+
| id          | bigint    | bigint unsigned | PRI        | auto_increment |
| user_id     | bigint    | bigint unsigned | MUL        |                |
| to_user_id  | bigint    | bigint unsigned | MUL        |                |
| relation    | tinyint   | tinyint(1)      |            |                |
+-------------+-----------+-----------------+------------+----------------+
```

# 项目运行
1. 安装go, docker, ffmpeg
2. 修改 docker-compose.yml，repository/init_db.go 中的数据库密码，连接地址
3. 运行mysql与redis  
   `docker-compose up`
4. 安装go依赖  
   `go mod tidy`
5. 运行项目  
    `go run server.go`

# 项目测试
- 运行测试，生成覆盖率报告
  - `bash test.sh`
- 目前测试覆盖率50%左右

# 难点
1. 用户鉴权  
   通过 user_id 生成 JWT token，验证时检查过期时间和 user_id 是否存在。
2. Feed流  
   数据库以时间戳形式存储视频上传的时间并建立索引，通过url参数latest_time查找视频。
   当访问完所有视频后，latest_time 为最后一个视频的上传时间。  
   这时再次访问家口，首先查询获得空列表，然后设置latest_time为当前时间，尝试再次查询并返回信息，从而访问新上传的视频。
3. 视频存储  
   项目没有使用对象存储服务，仅将文件存储在 static 对应年月日目录下，数据库中存储文件以 `/static/...` 开始的相对目录，查询时动态添加 `http://ip:port` 前缀。
4. 视频封面  
   视频封面用ffmpeg抽取视频第1帧生成，与视频存储在相同目录下。
5. redis缓存  
    如用户信息数据等高频访问的数据使用redis缓存。  
    其中 key 为 `user_{user_id}` 的形式, value 为数据库中对应的字段组成的哈希表。  
    初始时不存在缓存信息，查询时第一步查询redis，没有缓存则查询mysql，查到后写入redis。  
    更新信息时先修改mysql，然后删除redis缓存，从而保证一致性。
6. 并发安全  
    涉及多请求并发修改数据库的功能，项目使用 redis 作为并发访问的锁，保证并发安全。
    通过类似 `lock_{table}_{id}` 的key表示相关表对应id的锁，使用 uuid 作为value，标识创建锁的协程，设置过期时间10s。

# 项目贡献
- 张弛：视频接口、关注接口、项目搭建
- 叶静：用户登录注册接口
- 郑海森：评论接口
- 王智威：点赞接口
