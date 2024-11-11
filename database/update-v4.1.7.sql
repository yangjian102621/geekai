ALTER TABLE `chatgpt_sd_jobs` ADD `task_info` TEXT NOT NULL COMMENT '任务详情' AFTER `task_id`;
ALTER TABLE `chatgpt_mj_jobs` ADD `task_info` TEXT NOT NULL COMMENT '任务详情' AFTER `task_id`;
ALTER TABLE `chatgpt_dall_jobs` ADD `task_info` TEXT NOT NULL COMMENT '任务详情' AFTER `prompt`;
ALTER TABLE `chatgpt_suno_jobs` ADD `task_info` TEXT NOT NULL COMMENT '任务详情' AFTER `task_id`;
ALTER TABLE `chatgpt_video_jobs` CHANGE `params` `task_info` TEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '原始任务信息' AFTER `task_id`;