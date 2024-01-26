ALTER TABLE `chatgpt_mj_jobs` ADD `err_msg` VARCHAR(255) DEFAULT NULL COMMENT '错误信息' AFTER `publish`;
ALTER TABLE `chatgpt_sd_jobs` ADD `err_msg` VARCHAR(255) DEFAULT NULL COMMENT '错误信息' AFTER `publish`;

ALTER TABLE `chatgpt_chat_items` ADD `model` VARCHAR(30) NULL COMMENT '模型名称' AFTER `model_id`;