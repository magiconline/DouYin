CREATE DATABASE IF NOT EXISTS `douyin`;

USE `douyin`;

DROP TABLE IF EXISTS `video`;

CREATE TABLE `video`(
	`video_id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
	`user_id` BIGINT UNSIGNED NOT NULL,
	`play_url` TEXT NOT NULL,
	`cover_url` TEXT NOT NULL,
	`favorite_count` INT UNSIGNED NOT NULL,
	`comment_count` INT UNSIGNED NOT NULL,
	`upload_time` INT UNSIGNED NOT NULL,
	PRIMARY KEY (`video_id`)
);