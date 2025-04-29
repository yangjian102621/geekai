ALTER TABLE `chatgpt_suno_jobs` MODIFY `id` INT AUTO_INCREMENT;
ALTER TABLE `chatgpt_mj_jobs` CHANGE `channel_id` `channel_id` VARCHAR(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '频道ID';

truncate table chatgpt_rewards; -- 清空数据表

ALTER TABLE chatgpt_rewards RENAME TO chatgpt_redeems;
ALTER TABLE chatgpt_redeems COMMENT '兑换码';
ALTER TABLE `chatgpt_redeems` CHANGE `tx_id` `power` INT NOT NULL COMMENT '算力';
ALTER TABLE `chatgpt_redeems` DROP `remark`;
ALTER TABLE `chatgpt_redeems` DROP `exchange`;
ALTER TABLE `chatgpt_redeems` CHANGE `updated_at` `redeemed_at` INT NOT NULL COMMENT '兑换时间';
ALTER TABLE `chatgpt_redeems` CHANGE `amount` `code` VARCHAR(100) NOT NULL COMMENT '兑换码';
ALTER TABLE `chatgpt_redeems` DROP INDEX `tx_id`;
ALTER TABLE `chatgpt_redeems` ADD UNIQUE(`code`);
ALTER TABLE `chatgpt_redeems` ADD `name` VARCHAR(30) NOT NULL COMMENT '兑换码名称' AFTER `user_id`;
ALTER TABLE `chatgpt_redeems` CHANGE `status` `enabled` TINYINT(1) NOT NULL COMMENT '是否启用';