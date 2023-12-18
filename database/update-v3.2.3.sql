ALTER TABLE `chatgpt_mj_jobs` ADD `channel_id` CHAR(40) NULL DEFAULT NULL COMMENT '频道ID' AFTER `message_id`;
ALTER TABLE `chatgpt_mj_jobs` DROP INDEX `task_id`;