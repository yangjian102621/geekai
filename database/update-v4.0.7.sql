ALTER TABLE `chatgpt_mj_jobs` CHANGE `err_msg` `err_msg` VARCHAR(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '错误信息';
ALTER TABLE `chatgpt_sd_jobs` CHANGE `err_msg` `err_msg` VARCHAR(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '错误信息';
ALTER TABLE `chatgpt_dall_jobs` CHANGE `err_msg` `err_msg` VARCHAR(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '错误信息';
