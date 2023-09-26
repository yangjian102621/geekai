ALTER TABLE `chatgpt_mj_jobs` ADD `started` TINYINT(1) NOT NULL DEFAULT '0' COMMENT '任务是否开始' AFTER `progress`;
UPDATE `chatgpt_mj_jobs` SET started = 1