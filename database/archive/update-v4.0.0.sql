-- 删除用户名重复的用户，只保留一条
DELETE FROM chatgpt_users
WHERE username IN (
    SELECT username
    FROM (
             SELECT username
             FROM chatgpt_users
             GROUP BY username
             HAVING COUNT(*) > 1
         ) AS temp
) AND id NOT IN (
    SELECT MIN(id)
    FROM (
             SELECT id, username
             FROM chatgpt_users
             GROUP BY id, username
             HAVING COUNT(*) > 1
         ) AS temp
    GROUP BY username
);

-- 给 username 字段建立唯一索引
ALTER TABLE `chatgpt_users` ADD UNIQUE(`username`);

-- 当前用户剩余算力
ALTER TABLE `chatgpt_users` CHANGE `calls` `power` INT NOT NULL DEFAULT '0' COMMENT '剩余算力';
ALTER TABLE `chatgpt_users`
DROP `total_tokens`,
  DROP `tokens`,
  DROP `img_calls`;

ALTER TABLE `chatgpt_chat_models` CHANGE `weight` `power` TINYINT NOT NULL COMMENT '消耗算力点数';
ALTER TABLE `chatgpt_chat_models` ADD `temperature` FLOAT(3,1) NOT NULL DEFAULT '1' COMMENT '模型创意度' AFTER `power`, ADD `max_tokens` INT(11) NOT NULL DEFAULT '1024' COMMENT '最大响应长度' AFTER `temperature`, ADD `max_context` INT(11) NOT NULL DEFAULT '4096' COMMENT '最大上下文长度' AFTER `max_tokens`;

CREATE TABLE `chatgpt_plus`.`chatgpt_power_logs` ( `id` INT(11) NOT NULL AUTO_INCREMENT , `user_id` INT(11) NOT NULL COMMENT '用户ID' , `username` VARCHAR(30) NOT NULL COMMENT '用户名' , `type` TINYINT(1) NOT NULL COMMENT '类型（1：充值，2：消费，3：退费）' , `amount` SMALLINT(3) NOT NULL COMMENT '算力花费' , `balance` INT(11) NOT NULL COMMENT '余额' , `model` VARCHAR(30) NOT NULL COMMENT '模型' , `remark` VARCHAR(255) NOT NULL COMMENT '备注' , `created_at` DATETIME NOT NULL COMMENT '创建时间' , PRIMARY KEY (`id`)) ENGINE = InnoDB COMMENT = '用户算力消费日志';

ALTER TABLE `chatgpt_products` CHANGE `calls` `power` INT(11) NOT NULL DEFAULT '0' COMMENT '增加算力值';

ALTER TABLE `chatgpt_products` DROP `img_calls`;

ALTER TABLE `chatgpt_power_logs` CHANGE `amount` `amount` SMALLINT NOT NULL COMMENT '算力数值';
ALTER TABLE `chatgpt_power_logs` ADD `mark` TINYINT(1) NOT NULL COMMENT '资金类型（0：支出，1：收入）' AFTER `remark`;
ALTER TABLE `chatgpt_mj_jobs` ADD `power` SMALLINT(5) NOT NULL DEFAULT '0' COMMENT '消耗算力' AFTER `err_msg`;
ALTER TABLE `chatgpt_sd_jobs` ADD `power` SMALLINT(5) NOT NULL DEFAULT '0' COMMENT '消耗算力' AFTER `err_msg`;

ALTER TABLE `chatgpt_invite_logs` CHANGE `reward_json` `remark` VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '备注';

CREATE TABLE `chatgpt_admin_users`  (
                                        `id` int(0) NOT NULL AUTO_INCREMENT,
                                        `username` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '用户名',
                                        `password` char(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '密码',
                                        `salt` char(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '密码盐',
                                        `status` tinyint(1) NOT NULL COMMENT '当前状态',
                                        `last_login_at` int(0) NOT NULL COMMENT '最后登录时间',
                                        `last_login_ip` char(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '最后登录 IP',
                                        `created_at` datetime(0) NOT NULL COMMENT '创建时间',
                                        `updated_at` datetime(0) NOT NULL COMMENT '更新时间',
                                        PRIMARY KEY (`id`) USING BTREE,
                                        UNIQUE INDEX `username`(`username`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 108 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '系统用户' ROW_FORMAT = Dynamic;

INSERT INTO `chatgpt_admin_users` VALUES (1, 'admin', '6d17e80c87d209efb84ca4b2e0824f549d09fac8b2e1cc698de5bb5e1d75dfd0', 'mmrql75o', 1, 1710238055, '172.22.11.29', '2024-03-11 16:30:20', '2024-03-12 18:07:35');


ALTER TABLE `chatgpt_api_keys` CHANGE `use_proxy` `proxy_url` VARCHAR(100) NULL DEFAULT NULL COMMENT '代理地址';
-- 重置 proxy_url
UPDATE chatgpt_api_keys  set proxy_url='';

-- 重置系统配置，系统配置的数据结构变了，旧数据解析会失败。
UPDATE `chatgpt_configs` SET `config_json` = '{\"title\":\"ChatPlus AI 智能助手\",\"admin_title\":\"ChatPlus 控制台\",\"logo\":\"http://localhost:5678/static/upload/2024/3/1710732653645531.png\",\"init_power\":100,\"daily_power\":10,\"invite_power\":10,\"vip_month_power\":1000,\"register_ways\":[\"mobile\",\"username\",\"email\"],\"enabled_register\":true,\"reward_img\":\"http://localhost:5678/static/upload/2024/3/1710753716309668.jpg\",\"enabled_reward\":true,\"power_price\":0.1,\"order_pay_timeout\":1800,\"default_models\":[11,7,1,10,12,19,18,17],\"mj_power\":20,\"sd_power\":5,\"dall_power\":15,\"wechat_card_url\":\"/images/wx.png\",\"enable_context\":true,\"context_deep\":4}' WHERE `chatgpt_configs`.`id` = 1;

-- 重置用户默认模型
UPDATE `chatgpt_users` set chat_models_json = '[1]';

