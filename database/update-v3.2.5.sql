ALTER TABLE `chatgpt_orders` CHANGE `mobile` `username` VARCHAR(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '用户明';
ALTER TABLE `chatgpt_invite_logs` CHANGE `username` `username` VARCHAR(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '用户名';

CREATE TABLE `chatgpt_plus`.`chatgpt_files` (
    `id` INT(11) NOT NULL AUTO_INCREMENT ,
    `user_id` INT(11) NOT NULL COMMENT '用户 ID' ,
    `url` VARCHAR(255) NOT NULL COMMENT '文件地址' ,
    `ext` VARCHAR(10) NOT NULL COMMENT '文件后缀' ,
    `created_at` DATETIME NOT NULL COMMENT '创建时间' ,
    PRIMARY KEY (`id`)) ENGINE = InnoDB COMMENT = '用户文件表';

ALTER TABLE `chatgpt_files` ADD `size` BIGINT(20) NOT NULL DEFAULT '0' COMMENT '文件大小' AFTER `ext`;
ALTER TABLE `chatgpt_files` ADD `name` VARCHAR(100) NOT NULL COMMENT '文件名' AFTER `user_id`;