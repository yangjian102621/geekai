ALTER TABLE `chatgpt_mj_jobs` ADD `err_msg` VARCHAR(255) DEFAULT NULL COMMENT '错误信息' AFTER `publish`;
ALTER TABLE `chatgpt_sd_jobs` ADD `err_msg` VARCHAR(255) DEFAULT NULL COMMENT '错误信息' AFTER `publish`;

ALTER TABLE `chatgpt_chat_items` ADD `model` VARCHAR(30) NULL COMMENT '模型名称' AFTER `model_id`;
ALTER TABLE `chatgpt_chat_history` ADD `model` VARCHAR(30) NULL COMMENT '模型名称' AFTER `role_id`;

-- 初始化对话数据
UPDATE chatgpt_chat_items s SET model=(SELECT value FROM chatgpt_chat_models WHERE id = s.model_id);
-- 初始化聊天记录数据
UPDATE chatgpt_chat_history s SET model=(SELECT model FROM chatgpt_chat_items WHERE chat_id = s.chat_id);

-- 清理对话已删除的聊天记录（可选）
-- DELETE FROM `chatgpt_chat_history` WHERE model is NULL;

ALTER TABLE `chatgpt_files` ADD `obj_key` VARCHAR(100) NULL COMMENT '文件标识' AFTER `name`;