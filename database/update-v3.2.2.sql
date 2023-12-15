ALTER TABLE `chatgpt_mj_jobs` ADD `org_url` VARCHAR(400) NULL DEFAULT NULL COMMENT '原始图片地址' AFTER `img_url`;
ALTER TABLE `chatgpt_mj_jobs` DROP `started`;
ALTER TABLE `chatgpt_mj_jobs` ADD `task_id` VARCHAR(20) NULL DEFAULT NULL COMMENT '任务 ID' AFTER `user_id`;
ALTER TABLE `chatgpt_mj_jobs` ADD UNIQUE(`task_id`);
ALTER TABLE `chatgpt_sd_jobs` DROP `started`;