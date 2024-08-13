ALTER TABLE `chatgpt_users` ADD `mobile` CHAR(11) NULL COMMENT '手机号' AFTER `username`;
ALTER TABLE `chatgpt_users` ADD `email` VARCHAR(50) NULL COMMENT '邮箱地址' AFTER `mobile`;