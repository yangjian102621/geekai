ALTER TABLE `chatgpt_mj_jobs` ADD `channel_id` CHAR(40) NULL DEFAULT NULL COMMENT '频道ID' AFTER `message_id`;
ALTER TABLE `chatgpt_mj_jobs` DROP INDEX `task_id`;
ALTER TABLE `chatgpt_products` ADD `img_calls` INT(11) NOT NULL DEFAULT '0' COMMENT '绘图次数' AFTER `calls`;

CREATE TABLE `chatgpt_functions` (
                                     `id` int NOT NULL,
                                     `name` varchar(30) NOT NULL COMMENT '函数名称',
                                     `description` varchar(255) DEFAULT NULL COMMENT '函数描述',
                                     `parameters` text COMMENT '函数参数（JSON）',
                                     `required` varchar(255) NOT NULL COMMENT '必填参数（JSON）',
                                     `action` varchar(255) DEFAULT NULL COMMENT '函数处理 API'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='函数插件表';

ALTER TABLE `chatgpt_functions` ADD PRIMARY KEY (`id`);
ALTER TABLE `chatgpt_functions` MODIFY `id` int NOT NULL AUTO_INCREMENT;

ALTER TABLE `chatgpt_functions` ADD UNIQUE(`name`);

ALTER TABLE `chatgpt_functions` ADD `enabled` TINYINT(1) NOT NULL DEFAULT '0' COMMENT '是否启用' AFTER `action`;

ALTER TABLE `chatgpt_functions` ADD `label` VARCHAR(30) NULL COMMENT '函数标签' AFTER `name`;

ALTER TABLE `chatgpt_mj_jobs` ADD `use_proxy` TINYINT(1) NOT NULL DEFAULT '0' COMMENT '是否使用反代' AFTER `progress`;
ALTER TABLE `chatgpt_mj_jobs` CHANGE `img_url` `img_url` VARCHAR(400) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '图片URL';

ALTER TABLE `chatgpt_functions` ADD `token` VARCHAR(255) NULL COMMENT 'API授权token' AFTER `action`;

ALTER TABLE `chatgpt_functions` DROP `required`;
