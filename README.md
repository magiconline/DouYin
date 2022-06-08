# DouYin

# 接口文档
https://www.apifox.cn/apidoc/shared-8cc50618-0da6-4d5e-a398-76f3b8f766c5/api-18345145

# bug
- 使用 users表 还是 user 表

# 注册避免重名
https://zhuanlan.zhihu.com/p/111824045
1. 字段加唯一约束
2. 分布式锁
    - 一把锁 k, v=random
    - 多个锁 k = name, v=random


# SQL安全
- 非string类型进行类型转换
- 转义（https://gorm.io/zh_CN/docs/security.html）
- 使用预编译（https://blog.csdn.net/qq_39384184/article/details/108144309）

# TODO
1. 测试关注功能
   1. ~~对应user的粉丝count +- 1~~
   2. 检查user_id是否存在
   3. 自己关注自己
   4. 关注列表| 粉丝列表token验证失败
2. 使用联结查询
3. ~~SQL安全~~
4. 性能测试
5. 缓存

# 数据库对不存在的记录加锁查询
https://blog.csdn.net/u013269532/article/details/96841410  
https://www.cnblogs.com/jian0110/p/12721924.html  
行锁是基于索引的。
如果有索引，对不存在的数据会升级为范围锁；
如果没有索引，会升级为表锁。

# 测试
https://www.1024sou.com/article/915651.html  
https://gin-gonic.com/zh-cn/docs/testing/