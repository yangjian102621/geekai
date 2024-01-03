ALTER TABLE `chatgpt_users` ADD `nickname` VARCHAR(30) NOT NULL COMMENT '昵称' AFTER `mobile`;
ALTER TABLE `chatgpt_rewards` ADD `exchange` VARCHAR(255) NOT NULL COMMENT '兑换详情（json）' AFTER `status`;
