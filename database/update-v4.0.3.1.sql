-- 用户表添加OpenId和unionid
ALTER TABLE `chatgpt_users` ADD `official_openid` varchar(64) DEFAULT NULL COMMENT '公众号OpenId' AFTER `last_login_ip`;
ALTER TABLE `chatgpt_users` ADD `unionid` varchar(64) DEFAULT NULL COMMENT '公众号unionid' AFTER `official_openid`;
