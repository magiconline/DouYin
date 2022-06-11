CREATE DATABASE IF NOT EXISTS `douyin`;

USE `douyin`;

DROP TABLE IF EXISTS `video`;

CREATE TABLE `video`
(
    `video_id`       BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id`        BIGINT UNSIGNED NOT NULL,
    `play_url`       TEXT            NOT NULL,
    `cover_url`      TEXT            NOT NULL,
    `favorite_count` INT UNSIGNED    NOT NULL DEFAULT 0,
    `comment_count`  INT UNSIGNED    NOT NULL DEFAULT 0,
    `upload_time`    BIGINT UNSIGNED    NOT NULL,
    `title` TEXT,
    PRIMARY KEY (`video_id`)
);

DROP TABLE IF EXISTS `user`;

CREATE TABLE `user`
(
    `user_id`        BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_name`      VARCHAR(45)     NOT NULL,
    `password`       VARCHAR(45)     NOT NULL,
    `follow_count`   INT UNSIGNED    NOT NULL DEFAULT 0,
    `follower_count` INT UNSIGNED    NOT NULL DEFAULT 0,
    PRIMARY KEY (`user_id`)
);

create table star
(
    id           bigint unsigned auto_increment comment '主键ID' primary key,
    user_id      bigint unsigned  not null comment '用户ID',
    video_id     bigint unsigned  not null comment '视频ID'
 );

-- 插入点赞状态表
INSERT INTO douyin.star (user_id, video_id) VALUES (1, 1);

-- 插入视频
-- 需要插入数据库、拷贝文件、生成缩略图，sql无法模拟
-- INSERT INTO `video` (`video_id`, `user_id`, `play_url`, `cover_url`, `favorite_count`, `comment_count`, `upload_time`, `title`) VALUES (1, 1, "/static/2022/05/16/抖音-记录美好生活.mp4", "/static/2022/05/16/抖音-记录美好生活.jpg", 0, 0, REPLACE(unix_timestamp(current_timestamp(3)),'.',''), "测试title");

-- INSERT INTO `video` (`video_id`, `user_id`, `play_url`, `cover_url`, `favorite_count`, `comment_count`, `upload_time`, `title`) VALUES (2, 1, "/static/2022/05/16/抖音-记录美好生活(1).mp4", "/static/2022/05/16/抖音-记录美好生活.jpg", 0, 0, REPLACE(unix_timestamp(current_timestamp(3)),'.',''), "测试title");

-- INSERT INTO `video` (`video_id`, `user_id`, `play_url`, `cover_url`, `favorite_count`, `comment_count`, `upload_time`, `title`) VALUES (3, 1, "/static/2022/05/16/抖音-记录美好生活(2).mp4", "/static/2022/05/16/抖音-记录美好生活.jpg", 0, 0, REPLACE(unix_timestamp(current_timestamp(3)),'.',''), "测试title");

DROP TABLE IF EXISTS `remark`;
CREATE TABLE `remark`  (
                           `comment_id` integer(255) NOT NULL  auto_increment COMMENT '评论id',
                           `video_id` integer(255) NULL COMMENT '视频id',
                           `user_id` integer(255) NULL COMMENT '发出该评论的用户id',
                           `action_type` integer(255) NULL COMMENT '1-发布评论，2-删除评论',
                           `comment_text` varchar(255) NULL COMMENT '用户填写的评论内容',
                           `create_time` datetime NULL COMMENT '评论时间',
                           PRIMARY KEY (`comment_id`)
);
