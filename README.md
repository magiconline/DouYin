# DouYin

# 接口文档
https://www.apifox.cn/apidoc/shared-8cc50618-0da6-4d5e-a398-76f3b8f766c5/api-18345145

# SQL安全
- 非string类型进行类型转换
- 转义（https://gorm.io/zh_CN/docs/security.html）
- 使用预编译（https://blog.csdn.net/qq_39384184/article/details/108144309）

# TODO
1. 测试
2. 缓存


# 测试
https://www.1024sou.com/article/915651.html  
https://gin-gonic.com/zh-cn/docs/testing/  
`go test -bench . -benchmem -cpuprofile prof.cpu -memprofile prof.mem`

统计代码覆盖率：  
`go test -coverpkg=./... -coverprofile=cover.out`  

转换为html  
`go tool cover -html=cover.out -o cover.html`

# 缓存
- user信息 哈希表
  - k: "user_{user_id}"
  - v: ["user_name", "follow_count", "follower_count"]

- 关注信息 
  - k: "relation_{user_id}_{user_id}"
  - v: 1 or 0