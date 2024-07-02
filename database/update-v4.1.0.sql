ALTER TABLE `chatgpt_chat_models` CHANGE `power` `power` SMALLINT NOT NULL COMMENT '消耗算力点数';
ALTER TABLE `chatgpt_users` ADD `openid` VARCHAR(100) NULL COMMENT '第三方登录账号ID' AFTER `last_login_ip`;
ALTER TABLE `chatgpt_users` ADD `platform` VARCHAR(30) NULL COMMENT '登录平台' AFTER `openid`;