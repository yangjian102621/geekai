CREATE TABLE `chatgpt_suno_jobs` (
                                     `id` int NOT NULL,
                                     `user_id` int NOT NULL COMMENT '用户 ID',
                                     `title` varchar(100) DEFAULT NULL COMMENT '歌曲标题',
                                     `type` tinyint(1) DEFAULT '0' COMMENT '任务类型,1:灵感创作,2:自定义创作',
                                     `task_id` varchar(50) DEFAULT NULL COMMENT '任务 ID',
                                     `reference_id` char(50) DEFAULT NULL COMMENT '引用任务 ID',
                                     `tags` varchar(100) DEFAULT NULL COMMENT '歌曲风格',
                                     `instrumental` tinyint(1) DEFAULT '0' COMMENT '是否为纯音乐',
                                     `extend_secs` smallint DEFAULT '0' COMMENT '延长秒数',
                                     `song_id` varchar(50) DEFAULT NULL COMMENT '要续写的歌曲 ID',
                                     `prompt` varchar(2000) NOT NULL COMMENT '提示词',
                                     `thumb_img_url` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '缩略图地址',
                                     `cover_img_url` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '封面图地址',
                                     `audio_url` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '音频地址',
                                     `model_name` varchar(30) DEFAULT NULL COMMENT '模型地址',
                                     `progress` smallint DEFAULT '0' COMMENT '任务进度',
                                     `duration` smallint NOT NULL DEFAULT '0' COMMENT '歌曲时长',
                                     `publish` tinyint(1) NOT NULL COMMENT '是否发布',
                                     `err_msg` varchar(255) DEFAULT NULL COMMENT '错误信息',
                                     `raw_data` text COMMENT '原始数据',
                                     `power` smallint NOT NULL DEFAULT '0' COMMENT '消耗算力',
                                     `created_at` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='MidJourney 任务表';

--
-- 转储表的索引
--

--
-- 表的索引 `chatgpt_suno_jobs`
--
ALTER TABLE `chatgpt_suno_jobs`
    ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `task_id` (`task_id`);
