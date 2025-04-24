ALTER TABLE `chatgpt_chat_models` ADD `category` VARCHAR(1024) NOT NULL DEFAULT '' COMMENT '模型类别' AFTER `id`;
ALTER TABLE `chatgpt_chat_models` ADD `description` VARCHAR(1024) NOT NULL DEFAULT '' COMMENT '模型类型描述' AFTER `id`;
