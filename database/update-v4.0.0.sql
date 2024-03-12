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
ALTER TABLE `chatgpt_chat_models` ADD `temperature` FLOAT(3,2) NOT NULL DEFAULT '1' COMMENT '模型创意度' AFTER `power`, ADD `max_tokens` INT(11) NOT NULL DEFAULT '1024' COMMENT '最大响应长度' AFTER `temperature`, ADD `max_context` INT(11) NOT NULL DEFAULT '4096' COMMENT '最大上下文长度' AFTER `max_tokens`;

CREATE TABLE `chatgpt_plus`.`chatgpt_power_logs` ( `id` INT(11) NOT NULL AUTO_INCREMENT , `user_id` INT(11) NOT NULL COMMENT '用户ID' , `username` VARCHAR(30) NOT NULL COMMENT '用户名' , `type` TINYINT(1) NOT NULL COMMENT '类型（1：充值，2：消费，3：退费）' , `amount` SMALLINT(3) NOT NULL COMMENT '算力花费' , `balance` INT(11) NOT NULL COMMENT '余额' , `model` VARCHAR(30) NOT NULL COMMENT '模型' , `remark` VARCHAR(255) NOT NULL COMMENT '备注' , `created_at` DATETIME NOT NULL COMMENT '创建时间' , PRIMARY KEY (`id`)) ENGINE = InnoDB COMMENT = '用户算力消费日志';

ALTER TABLE `chatgpt_products` CHANGE `calls` `power` INT(11) NOT NULL DEFAULT '0' COMMENT '增加算力值';

ALTER TABLE `chatgpt_products` DROP `img_calls`;

ALTER TABLE `chatgpt_power_logs` CHANGE `amount` `amount` SMALLINT NOT NULL COMMENT '算力数值';
ALTER TABLE `chatgpt_power_logs` ADD `mark` TINYINT(1) NOT NULL COMMENT '资金类型（0：支出，1：收入）' AFTER `remark`;
ALTER TABLE `chatgpt_mj_jobs` ADD `power` SMALLINT(5) NOT NULL DEFAULT '0' COMMENT '消耗算力' AFTER `err_msg`;
ALTER TABLE `chatgpt_sd_jobs` ADD `power` SMALLINT(5) NOT NULL DEFAULT '0' COMMENT '消耗算力' AFTER `err_msg`;