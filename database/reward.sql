CREATE TABLE `chatgpt_rewards` (
                                  `id` int NOT NULL,
                                  `tx_id` char(36) NOT NULL COMMENT '交易 ID',
                                  `amount` decimal(10,2) NOT NULL COMMENT '打赏金额',
                                  `remark` varchar(80) NOT NULL COMMENT '备注',
                                  `status` tinyint(1) NOT NULL COMMENT '核销状态，0：未核销，1：已核销',
                                  `created_at` datetime NOT NULL,
                                  `updated_at` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户打赏';

ALTER TABLE `chatgpt_rewards`
    ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `tx_id` (`tx_id`);


ALTER TABLE `chatgpt_rewards`
    MODIFY `id` int NOT NULL AUTO_INCREMENT;

update chatgpt_users set calls=0

ALTER TABLE `chatgpt_rewards` ADD `user_id` INT(11) NOT NULL COMMENT '用户 ID' AFTER `id`;