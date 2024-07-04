ALTER TABLE `chatgpt_chat_models` CHANGE `power` `power` SMALLINT NOT NULL COMMENT '消耗算力点数';
ALTER TABLE `chatgpt_users` ADD `openid` VARCHAR(100) NULL COMMENT '第三方登录账号ID' AFTER `last_login_ip`;
ALTER TABLE `chatgpt_users` ADD UNIQUE(`openid`);
ALTER TABLE `chatgpt_users` ADD `platform` VARCHAR(30) NULL COMMENT '登录平台' AFTER `openid`;
ALTER TABLE `chatgpt_users` CHANGE `avatar` `avatar` VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '头像';
ALTER TABLE `chatgpt_chat_history` CHANGE `icon` `icon` VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '角色图标';