ALTER TABLE `chatgpt_mj_jobs` DROP `image`;
ALTER TABLE `chatgpt_mj_jobs` DROP `hash`;
ALTER TABLE `chatgpt_mj_jobs` DROP `content`;
ALTER TABLE `chatgpt_mj_jobs` DROP `chat_id`;
ALTER TABLE `chatgpt_mj_jobs` ADD `progress` SMALLINT(5) NULL DEFAULT '0' COMMENT '任务进度' AFTER `prompt`;
ALTER TABLE `chatgpt_mj_jobs` ADD `hash` VARCHAR(100) NULL DEFAULT NULL COMMENT 'message hash' AFTER `prompt`;
ALTER TABLE `chatgpt_mj_jobs` ADD `img_url` VARCHAR(255) NULL DEFAULT NULL COMMENT '图片URL' AFTER `prompt`;

-- 2023-09-15
ALTER TABLE `chatgpt_mj_jobs` ADD `type` VARCHAR(20) NULL DEFAULT 'image' COMMENT '任务类别' AFTER `user_id`;
ALTER TABLE `chatgpt_mj_jobs` DROP INDEX `message_id`;
ALTER TABLE `chatgpt_mj_jobs` ADD INDEX(`message_id`);