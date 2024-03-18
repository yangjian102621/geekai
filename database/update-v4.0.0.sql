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
ALTER TABLE `chatgpt_users` ADD UNIQUE(`username`)

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

CREATE TABLE `chatgpt_admin_user_roles` (
                                            `id` int NOT NULL AUTO_INCREMENT COMMENT 'ID',
                                            `admin_id` int NOT NULL COMMENT '用户ID',
                                            `role_id` int NOT NULL COMMENT '角色ID',
                                            `created_at` datetime NOT NULL,
                                            `updated_at` datetime NOT NULL,
                                            PRIMARY KEY (`id`) USING BTREE,
                                            KEY `idx_admin_id` (`admin_id`),
                                            KEY `idx_role_id` (`role_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户角色表';

CREATE TABLE `chatgpt_admin_roles` (
                                       `id` int NOT NULL AUTO_INCREMENT COMMENT 'ID',
                                       `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '角色名称',
                                       `description` varchar(255) DEFAULT '' COMMENT '描述',
                                       `created_at` datetime NOT NULL,
                                       `updated_at` datetime NOT NULL,
                                       PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='角色表';

CREATE TABLE `chatgpt_admin_role_permissions` (
                                                  `id` int NOT NULL AUTO_INCREMENT COMMENT 'ID',
                                                  `role_id` int NOT NULL COMMENT '角色ID',
                                                  `permission_id` int NOT NULL COMMENT '权限ID',
                                                  `created_at` datetime NOT NULL,
                                                  `updated_at` datetime NOT NULL,
                                                  PRIMARY KEY (`id`) USING BTREE,
                                                  KEY `idx_role_id` (`role_id`),
                                                  KEY `idx_permission_id` (`permission_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户角色表';

CREATE TABLE `chatgpt_admin_permissions` (
                                             `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
                                             `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '权限名称（模块@菜单名称@方法）',
                                             `slug` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '前端权限模块',
                                             `sort` int NOT NULL DEFAULT '1' COMMENT '排序(越大排越靠前)',
                                             `pid` int NOT NULL DEFAULT '0' COMMENT '父节点',
                                             `created_at` datetime DEFAULT NULL,
                                             `updated_at` datetime DEFAULT NULL,
                                             PRIMARY KEY (`id`),
                                             KEY `idx_slug` (`slug`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=35 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='权限表';
INSERT INTO `chatgpt_admin_permissions` VALUES (6, '用户管理', '', 2, 0, '2024-03-14 10:47:41', '2024-03-14 10:48:41');
INSERT INTO `chatgpt_admin_permissions` VALUES (7, '列表', 'api_admin_user_list', 1, 6, '2024-03-14 10:53:36', '2024-03-14 10:53:36');
INSERT INTO `chatgpt_admin_permissions` VALUES (10, '角色模型', '', 3, 0, '2024-03-14 11:02:45', '2024-03-14 11:02:45');
INSERT INTO `chatgpt_admin_permissions` VALUES (11, '列表', 'api_admin_role_list', 1, 10, '2024-03-14 11:03:43', '2024-03-14 11:03:45');
INSERT INTO `chatgpt_admin_permissions` VALUES (12, '语言模型', '', 4, 0, '2024-03-14 11:05:00', '2024-03-14 11:05:00');
INSERT INTO `chatgpt_admin_permissions` VALUES (13, '列表', 'api_admin_model_list', 1, 12, '2024-03-14 11:05:36', '2024-03-14 11:05:36');
INSERT INTO `chatgpt_admin_permissions` VALUES (14, '充值产品', '', 5, 0, '2024-03-14 11:16:37', '2024-03-14 11:16:37');
INSERT INTO `chatgpt_admin_permissions` VALUES (15, '列表', 'api_admin_product_list', 1, 14, '2024-03-14 11:17:05', '2024-03-14 11:17:05');
INSERT INTO `chatgpt_admin_permissions` VALUES (16, 'APIKEY', '', 6, 0, '2024-03-14 11:17:16', '2024-03-14 11:17:16');
INSERT INTO `chatgpt_admin_permissions` VALUES (17, '列表', 'api_admin_apikey_list', 1, 16, '2024-03-14 11:17:53', '2024-03-14 11:18:02');
INSERT INTO `chatgpt_admin_permissions` VALUES (18, '充值订单', '', 7, 0, '2024-03-14 11:18:20', '2024-03-14 11:18:20');
INSERT INTO `chatgpt_admin_permissions` VALUES (19, '列表', 'api_admin_order_list', 1, 18, '2024-03-14 11:18:44', '2024-03-14 11:18:44');
INSERT INTO `chatgpt_admin_permissions` VALUES (20, '众筹管理', '', 8, 0, '2024-03-14 11:18:59', '2024-03-14 11:18:59');
INSERT INTO `chatgpt_admin_permissions` VALUES (21, '列表', 'api_admin_reward_list', 1, 20, '2024-03-14 11:19:28', '2024-03-14 11:19:28');
INSERT INTO `chatgpt_admin_permissions` VALUES (22, '函数管理', '', 9, 0, '2024-03-14 11:19:42', '2024-03-14 11:19:42');
INSERT INTO `chatgpt_admin_permissions` VALUES (23, '列表', 'api_admin_function_list', 1, 22, '2024-03-14 11:20:07', '2024-03-14 11:20:07');
INSERT INTO `chatgpt_admin_permissions` VALUES (24, '对话管理', '', 10, 0, '2024-03-14 11:20:30', '2024-03-14 11:20:30');
INSERT INTO `chatgpt_admin_permissions` VALUES (25, '列表', 'api_admin_chat_list', 1, 24, '2024-03-14 11:20:51', '2024-03-14 11:20:51');
INSERT INTO `chatgpt_admin_permissions` VALUES (26, '网站设置', '', 11, 0, '2024-03-14 11:21:05', '2024-03-14 11:21:05');
INSERT INTO `chatgpt_admin_permissions` VALUES (27, '配置列表', 'api_admin_config_get', 1, 26, '2024-03-14 11:21:53', '2024-03-14 11:21:53');
INSERT INTO `chatgpt_admin_permissions` VALUES (28, '系统配置', '', 12, 0, '2024-03-14 11:22:07', '2024-03-14 15:28:22');
INSERT INTO `chatgpt_admin_permissions` VALUES (29, '列表', 'api_admin_sysUser_list', 1, 30, '2024-03-14 11:22:36', '2024-03-14 15:28:52');
INSERT INTO `chatgpt_admin_permissions` VALUES (30, '系统管理员', '', 1, 28, '2024-03-14 15:28:43', '2024-03-14 15:28:43');
INSERT INTO `chatgpt_admin_permissions` VALUES (31, '权限配置', '', 2, 28, '2024-03-14 15:29:05', '2024-03-14 15:29:05');
INSERT INTO `chatgpt_admin_permissions` VALUES (32, '角色配置', '', 3, 28, '2024-03-14 15:29:15', '2024-03-14 15:29:15');
INSERT INTO `chatgpt_admin_permissions` VALUES (33, '列表', 'api_admin_sysPermission_list', 1, 31, '2024-03-14 15:29:52', '2024-03-14 15:29:52');
INSERT INTO `chatgpt_admin_permissions` VALUES (34, '列表', 'api_admin_sysRole_list', 1, 32, '2024-03-14 15:30:21', '2024-03-14 15:30:21');


ALTER TABLE `chatgpt_api_keys` CHANGE `use_proxy` `proxy_url` VARCHAR(100) NULL DEFAULT NULL COMMENT '代理地址';
-- 重置 proxy_url
UPDATE chatgpt_api_keys  set proxy_url=''