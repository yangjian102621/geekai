ALTER TABLE `chatgpt_mj_jobs` ADD `started` TINYINT(1) NOT NULL DEFAULT '0' COMMENT '任务是否开始' AFTER `progress`;
UPDATE `chatgpt_mj_jobs` SET started = 1

-- 创建 SD 绘图任务表
CREATE TABLE `chatgpt_sd_jobs` (
                                   `id` int NOT NULL,
                                   `user_id` int NOT NULL COMMENT '用户 ID',
                                   `type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT 'txt2img' COMMENT '任务类别',
                                   `task_id` char(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '任务 ID',
                                   `prompt` varchar(2000) NOT NULL COMMENT '会话提示词',
                                   `img_url` varchar(255) DEFAULT NULL COMMENT '图片URL',
                                   `params` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '绘画参数json',
                                   `progress` smallint DEFAULT '0' COMMENT '任务进度',
                                   `started` tinyint(1) NOT NULL DEFAULT '0' COMMENT '任务是否开始',
                                   `created_at` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='StableDiffusion 任务表';
--
-- 表的索引 `chatgpt_sd_jobs`
--
ALTER TABLE `chatgpt_sd_jobs`
    ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `task_id` (`task_id`);

ALTER TABLE `chatgpt_sd_jobs`
    MODIFY `id` int NOT NULL AUTO_INCREMENT;