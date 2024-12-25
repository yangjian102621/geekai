ALTER TABLE `chatgpt_suno_jobs` CHANGE `err_msg` `err_msg` VARCHAR(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '错误信息';

ALTER TABLE `chatgpt_sd_jobs` CHANGE `err_msg` `err_msg` VARCHAR(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '错误信息';

ALTER TABLE `chatgpt_mj_jobs` CHANGE `err_msg` `err_msg` VARCHAR(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '错误信息';

ALTER TABLE `chatgpt_dall_jobs` CHANGE `err_msg` `err_msg` VARCHAR(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '错误信息';
ALTER TABLE `chatgpt_video_jobs` CHANGE `err_msg` `err_msg` VARCHAR(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '错误信息';
ALTER TABLE `chatgpt_chat_models` ADD `type` VARCHAR(10) NOT NULL DEFAULT 'chat' COMMENT '模型类型（chat,img）' AFTER `id`;