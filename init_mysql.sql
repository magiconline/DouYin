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
    PRIMARY KEY (`video_id`),
    INDEX(`user_id`)
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

DROP TABLE iF EXISTS `star`;
create table star
(
    id           bigint unsigned auto_increment comment '主键ID' primary key,
    user_id      bigint unsigned  not null comment '用户ID',
    video_id     bigint unsigned  not null comment '视频ID',
    INDEX(`user_id`, `video_id`)
 );

DROP TABLE IF EXISTS `remark`;
CREATE TABLE `remark`  (
                           `comment_id` integer(255) NOT NULL  auto_increment COMMENT '评论id',
                           `video_id` integer(255) NULL COMMENT '视频id',
                           `user_id` integer(255) NULL COMMENT '发出该评论的用户id',
                           `action_type` integer(255) NULL COMMENT '1-发布评论，2-删除评论',
                           `comment_text` varchar(255) NULL COMMENT '用户填写的评论内容',
                           `create_time` datetime NULL COMMENT '评论时间',
                           PRIMARY KEY (`comment_id`),
                           INDEX(`video_id`)
);
