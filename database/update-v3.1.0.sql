ALTER TABLE `chatgpt_chat_items` CHANGE `model` `model_id` INT(11) NOT NULL DEFAULT '0' COMMENT '模型 ID';
ALTER TABLE `chatgpt_api_keys` ADD `platform` CHAR(20)  DEFAULT NULL COMMENT '平台' AFTER id;
ALTER TABLE `chatgpt_users` CHANGE `tokens` `total_tokens` BIGINT NOT NULL DEFAULT '0' COMMENT '累计消耗 tokens';
ALTER TABLE `chatgpt_chat_items` ADD `deleted_at` DATETIME NULL DEFAULT NULL AFTER `updated_at`;
ALTER TABLE `chatgpt_chat_history` ADD `deleted_at` DATETIME NULL DEFAULT NULL AFTER `updated_at`;

CREATE TABLE `chatgpt_chat_models` (
                                       `id` int NOT NULL,
                                       `platform` varchar(20) DEFAULT NULL COMMENT '模型平台',
                                       `name` varchar(50) NOT NULL COMMENT '模型名称',
                                       `value` varchar(50) NOT NULL COMMENT '模型值',
                                       `sort_num` tinyint(1) NOT NULL COMMENT '排序数字',
                                       `enabled` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否启用模型',
                                       `created_at` datetime DEFAULT NULL,
                                       `updated_at` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='AI 模型表';
ALTER TABLE `chatgpt_chat_models`
    ADD PRIMARY KEY (`id`);
ALTER TABLE `chatgpt_chat_models`
    MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=7;

INSERT INTO `chatgpt_chat_models` (`id`, `platform`, `name`, `value`, `sort_num`, `enabled`, `created_at`, `updated_at`) VALUES
 (1, 'OpenAI', 'Bot GPT-3.5', 'gpt-3.5-turbo', 0, 1, '2023-08-23 12:06:36', '2023-09-02 16:49:36'),
 (2, 'Azure', 'Bot Azure-3.5', 'gpt-3.5-turbo', 0, 1, '2023-08-23 12:15:30', '2023-09-02 16:49:46'),
 (3, 'ChatGML', 'ChatGML-Pro', 'chatglm_pro', 3, 1, '2023-08-23 13:35:45', '2023-08-29 11:41:29'),
 (5, 'ChatGML', 'ChatGLM-Std', 'chatglm_std', 2, 1, '2023-08-24 15:05:38', '2023-08-29 11:41:28'),
 (6, 'ChatGML', 'ChatGLM-Lite', 'chatglm_lite', 4, 1, '2023-08-24 15:06:15', '2023-08-29 11:41:29');

ALTER TABLE `chatgpt_users`
DROP `username`,
DROP `nickname`;