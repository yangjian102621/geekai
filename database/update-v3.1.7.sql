-- 增加字段保存用户开通的 AI 模型
ALTER TABLE `chatgpt_users` ADD `chat_models_json` TEXT NOT NULL COMMENT 'AI模型 json' AFTER `chat_roles_json`;

-- 为每个模型设置对话权重
ALTER TABLE `chatgpt_chat_models` ADD `weight` TINYINT(3) NOT NULL COMMENT '对话权重，每次对话扣减多少次对话额度' AFTER `enabled`;
UPDATE `chatgpt_chat_models` SET weight = 1;