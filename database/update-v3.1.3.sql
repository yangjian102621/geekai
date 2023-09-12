ALTER TABLE `chatgpt_mj_jobs` DROP `image`;
ALTER TABLE `chatgpt_mj_jobs` DROP `hash`;
ALTER TABLE `chatgpt_mj_jobs` DROP `content`;
ALTER TABLE `chatgpt_mj_jobs` DROP `chat_id`;
ALTER TABLE `chatgpt_mj_jobs` ADD `progress` SMALLINT(5) NULL DEFAULT '0' COMMENT '任务进度' AFTER `prompt`;
ALTER TABLE `chatgpt_mj_jobs` ADD `hash` VARCHAR(100) NULL DEFAULT NULL COMMENT 'message hash' AFTER `prompt`;
ALTER TABLE `chatgpt_mj_jobs` ADD `img_url` VARCHAR(255) NULL DEFAULT NULL COMMENT '图片URL' AFTER `prompt`;