CREATE TABLE `chatgpt_plus`.`chatgpt_dall_jobs` ( `id` INT(11) NOT NULL AUTO_INCREMENT , `user_id` INT(11) NOT NULL COMMENT '用户ID' , `task_id` VARCHAR(20) NOT NULL COMMENT '任务ID' , `prompt` VARCHAR(2000) NOT NULL COMMENT '提示词' , `img_url` VARCHAR(255) NOT NULL COMMENT '图片地址' , `publish` TINYINT(1) NOT NULL COMMENT '是否发布' , `power` SMALLINT(3) NOT NULL COMMENT '消耗算力' , `progress` SMALLINT(3) NOT NULL COMMENT '任务进度' , `err_msg` VARCHAR(255) NOT NULL COMMENT '错误信息' , `created_at` DATETIME NOT NULL , PRIMARY KEY (`id`)) ENGINE = InnoDB COMMENT = 'DALLE 绘图任务表';

ALTER TABLE `chatgpt_dall_jobs` ADD `org_url` VARCHAR(400) NULL COMMENT '原图地址' AFTER `img_url`;
ALTER TABLE `chatgpt_dall_jobs` DROP `task_id`;


