ALTER TABLE `chatgpt_chat_roles` ADD `model_id` INT NOT NULL DEFAULT '0' COMMENT '绑定模型ID' AFTER `sort_num`;
ALTER TABLE `chatgpt_chat_models` ADD `key_id` INT(11) NOT NULL COMMENT '绑定API KEY ID' AFTER `open`;
INSERT INTO `chatgpt_plus`.`chatgpt_menus`(`id`, `name`, `icon`, `url`, `sort_num`, `enabled`) VALUES (12, '思维导图', '/images/menu/xmind.png', '/xmind', 3, 1);