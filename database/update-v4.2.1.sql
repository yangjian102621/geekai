ALTER TABLE `chatgpt_video_jobs` CHANGE `prompt` `prompt` TEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '提示词';
ALTER TABLE `chatgpt_video_jobs` CHANGE `prompt_ext` `prompt_ext` TEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '优化后提示词';

ALTER TABLE `chatgpt_mj_jobs` CHANGE `prompt` `prompt` TEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '会话提示词';
ALTER TABLE `chatgpt_sd_jobs` CHANGE `prompt` `prompt` TEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '会话提示词';
ALTER TABLE `chatgpt_dall_jobs` CHANGE `prompt` `prompt` TEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '提示词';

ALTER TABLE `chatgpt_files` CHANGE `name` `name` VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '文件名';
ALTER TABLE `chatgpt_chat_models` CHANGE `name` `name` VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '模型名称';