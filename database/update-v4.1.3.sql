ALTER TABLE `chatgpt_users` ADD `mobile` CHAR(11) NULL COMMENT '手机号' AFTER `username`;
ALTER TABLE `chatgpt_users` ADD `email` VARCHAR(50) NULL COMMENT '邮箱地址' AFTER `mobile`;

CREATE TABLE `chatgpt_video_jobs` (
                                      `id` int NOT NULL,
                                      `user_id` int NOT NULL COMMENT '用户 ID',
                                      `channel` varchar(100) NOT NULL COMMENT '渠道',
                                      `task_id` varchar(100) NOT NULL COMMENT '任务 ID',
                                      `type` varchar(20) DEFAULT NULL COMMENT '任务类型,luma,runway,cogvideo',
                                      `prompt` varchar(2000) NOT NULL COMMENT '提示词',
                                      `prompt_ext` varchar(2000) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '优化后提示词',
                                      `cover_url` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '封面图地址',
                                      `video_url` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '视频地址',
                                      `water_url` varchar(512) DEFAULT NULL COMMENT '带水印的视频地址',
                                      `progress` smallint DEFAULT '0' COMMENT '任务进度',
                                      `publish` tinyint(1) NOT NULL COMMENT '是否发布',
                                      `err_msg` varchar(255) DEFAULT NULL COMMENT '错误信息',
                                      `raw_data` text COMMENT '原始数据',
                                      `power` smallint NOT NULL DEFAULT '0' COMMENT '消耗算力',
                                      `created_at` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='MidJourney 任务表';

ALTER TABLE `chatgpt_video_jobs`ADD PRIMARY KEY (`id`);

ALTER TABLE `chatgpt_video_jobs` MODIFY `id` int NOT NULL AUTO_INCREMENT;

ALTER TABLE `chatgpt_video_jobs` ADD `params` VARCHAR(512) NULL COMMENT '参数JSON' AFTER `raw_data`;