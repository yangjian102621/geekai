CREATE TABLE `chatgpt_mj_jobs` (
                                   `id` int NOT NULL,
                                   `user_id` int NOT NULL COMMENT '用户 ID',
                                   `chat_id` char(40) NOT NULL COMMENT '聊天会话 ID',
                                   `message_id` char(40) NOT NULL COMMENT '消息 ID',
                                   `hash` char(40) NOT NULL COMMENT '图片哈希',
                                   `content` varchar(2000) NOT NULL COMMENT '消息内容',
                                   `prompt` varchar(2000) NOT NULL COMMENT '会话提示词',
                                   `image` text NOT NULL COMMENT '图片信息 json',
                                   `created_at` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='MidJourney 任务表';

--
-- 转储表的索引
--

--
-- 表的索引 `chatgpt_mj_jobs`
--
ALTER TABLE `chatgpt_mj_jobs`
    ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `message_id` (`message_id`),
  ADD UNIQUE KEY `hash` (`hash`);

--
-- 在导出的表使用AUTO_INCREMENT
--

--
-- 使用表AUTO_INCREMENT `chatgpt_mj_jobs`
--
ALTER TABLE `chatgpt_mj_jobs`
    MODIFY `id` int NOT NULL AUTO_INCREMENT;
COMMIT;
