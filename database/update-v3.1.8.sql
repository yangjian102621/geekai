ALTER TABLE `chatgpt_users` ADD `vip` TINYINT(1) NOT NULL DEFAULT '0' COMMENT '是否会员' AFTER `last_login_at`;
ALTER TABLE `chatgpt_users` ADD `tokens` BIGINT NOT NULL DEFAULT '0' COMMENT '当月消耗 tokens' AFTER `total_tokens`;

CREATE TABLE `chatgpt_orders` (
  `id` int NOT NULL,
  `user_id` int NOT NULL COMMENT '用户ID',
  `product_id` int NOT NULL COMMENT '产品ID',
  `mobile` char(11) NOT NULL COMMENT '用户手机号',
  `order_no` varchar(30) NOT NULL COMMENT '订单ID',
  `subject` varchar(100) NOT NULL COMMENT '订单产品',
  `amount` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '订单金额',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '订单状态（0：待支付，1：已扫码，2：支付失败）',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '备注',
  `pay_time` int DEFAULT NULL COMMENT '支付时间',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `deleted_at` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='充值订单表';

-- 创建索引
ALTER TABLE `chatgpt_orders` ADD PRIMARY KEY (`id`), ADD UNIQUE KEY `order_no` (`order_no`);
ALTER TABLE `chatgpt_orders` MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=13;

CREATE TABLE `chatgpt_products` (
    `id` int NOT NULL,
    `name` varchar(30) NOT NULL COMMENT '名称',
    `price` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '价格',
    `discount` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '优惠金额',
    `days` smallint NOT NULL DEFAULT '0' COMMENT '延长天数',
    `calls` int NOT NULL DEFAULT '0' COMMENT '调用次数',
    `enabled` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否启动',
    `sales` int NOT NULL DEFAULT '0' COMMENT '销量',
    `sort_num` tinyint NOT NULL DEFAULT '0' COMMENT '排序',
    `created_at` datetime NOT NULL,
    `updated_at` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='会员套餐表';

INSERT INTO `chatgpt_products` (`id`, `name`, `price`, `discount`, `days`, `calls`, `enabled`, `sales`, `sort_num`, `created_at`, `updated_at`) VALUES
(1, '会员1个月', '1.01', '1.00', 30, 0, 1, 0, 0, '2023-08-28 10:48:57', '2023-08-31 16:24:26'),
(2, '会员3个月', '140.00', '30.00', 90, 0, 1, 0, 0, '2023-08-28 10:52:22', '2023-08-31 16:24:31'),
(3, '会员6个月', '290.00', '100.00', 180, 0, 1, 0, 0, '2023-08-28 10:53:39', '2023-08-31 16:24:36'),
(4, '会员12个月', '580.00', '200.00', 365, 0, 1, 0, 0, '2023-08-28 10:54:15', '2023-08-31 16:24:42'),
(5, '100次点卡', '10.03', '10.00', 0, 100, 1, 0, 0, '2023-08-28 10:55:08', '2023-08-31 17:34:43');

ALTER TABLE `chatgpt_products` ADD PRIMARY KEY (`id`);
ALTER TABLE `chatgpt_products` MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=6;

ALTER TABLE `chatgpt_orders` ADD `pay_way` VARCHAR(20) DEFAULT '0'  NOT NULL COMMENT '支付方式' AFTER `pay_time`;