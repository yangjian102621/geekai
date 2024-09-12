CREATE TABLE `chatgpt_app_types` (
                                      `id` int NOT NULL,
                                      `name` varchar(50) NOT NULL COMMENT '名称',
                                      `icon` varchar(255) NOT NULL COMMENT '图标URL',
                                      `sort_num` tinyint(3) NOT NULL COMMENT '排序',
                                      `enabled` tinyint(1) NOT NULL COMMENT '是否启用',
                                      `created_at` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='应用分类表';

ALTER TABLE `chatgpt_app_types`ADD PRIMARY KEY (`id`);
ALTER TABLE `chatgpt_app_types` MODIFY `id` int NOT NULL AUTO_INCREMENT;
ALTER TABLE `chatgpt_chat_roles` ADD `tid` INT NOT NULL COMMENT '分类ID' AFTER `name`;