# 数据库设计

点赞表

| 字段名称     |       类型        | 允许空值 | 主键 | 说明                             |
| ------------ | :---------------: | -------- | ---- | -------------------------------- |
| id           |  BIGINT UNSIGNED  | ×        | √    | ID                               |
| user_id      |  BIGINT UNSIGNED  | ×        |      | 用户ID                           |
| star_type    | SMALLINT UNSIGNED | ×        |      | 点赞种类 0-视频点赞 1-评论区点赞 |
| star_type_id |  BIGINT UNSIGNED  | ×        |      | 评论或者视频ID                   |
| status       | SMALLINT UNSIGNED | ×        |      | 点赞状态 0-取消赞 1-有效点赞     |
