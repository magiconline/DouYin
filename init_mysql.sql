CREATE DATABASE IF NOT EXISTS `douyin`;

USE `douyin`;

DROP TABLE IF EXISTS `video`;

CREATE TABLE `video`
(
    `video_id`       BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id`        BIGINT UNSIGNED NOT NULL,
    `play_url`       TEXT            NOT NULL,
    `cover_url`      TEXT            NOT NULL,
    `favorite_count` INT UNSIGNED    NOT NULL,
    `comment_count`  INT UNSIGNED    NOT NULL,
    `upload_time`    BIGINT UNSIGNED    NOT NULL,
    PRIMARY KEY (`video_id`)
);

CREATE TABLE `user`
(
    `user_id`        BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_name`      VARCHAR(45)     NOT NULL,
    `password`       VARCHAR(45)     NOT NULL,
    `follow_count`   INT UNSIGNED    NOT NULL,
    `follower_count` INT UNSIGNED    NOT NULL,
    `token`          VARCHAR(45)     NOT NULL,
    PRIMARY KEY (`user_id`)
);

create table star
(
    id           bigint unsigned auto_increment comment '主键ID' primary key,
    user_id      bigint unsigned  not null comment '用户ID',
    video_id     bigint unsigned  not null comment '视频ID'
 );


-- 插入用户
INSERT INTO `users` (`id`, `user_name`, `password`, `follow_count`, `follower_count`, `token`) VALUES (1, "test_user", "123", 0, 0, "abc");
-- 插入点赞状态表
INSERT INTO douyin.star (user_id, video_id) VALUES (1, 1);
-- 插入视频
INSERT INTO `video` (`video_id`, `user_id`, `play_url`, `cover_url`, `favorite_count`, `comment_count`, `upload_time`) VALUES (1, 1, "/static/2022/05/16/抖音-记录美好生活.mp4", "/static/2022/05/16/抖音-记录美好生活.jpg", 0, 0, UNIX_TIMESTAMP());
INSERT INTO `video` (`video_id`, `user_id`, `play_url`, `cover_url`, `favorite_count`, `comment_count`, `upload_time`) VALUES (2, 1, "/static/2022/05/16/抖音-记录美好生活(1).mp4", "/static/2022/05/16/抖音-记录美好生活.jpg", 0, 0, UNIX_TIMESTAMP());
INSERT INTO `video` (`video_id`, `user_id`, `play_url`, `cover_url`, `favorite_count`, `comment_count`, `upload_time`) VALUES (3, 1, "/static/2022/05/16/抖音-记录美好生活(2).mp4", "/static/2022/05/16/抖音-记录美好生活.jpg", 0, 0, UNIX_TIMESTAMP());