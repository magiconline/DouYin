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

CREATE TABLE `user` (
  `user_id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_name` VARCHAR(45) NOT NULL,
  `password` VARCHAR(45) NOT NULL,
  `follow_count` INT UNSIGNED NOT NULL,
  `follower_count` INT UNSIGNED NOT NULL,
  `token` VARCHAR(45) NOT NULL,
  PRIMARY KEY (`user_id`));