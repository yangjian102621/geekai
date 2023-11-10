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

--
-- 转存表中的数据 `chatgpt_orders`
--
INSERT INTO `chatgpt_orders` (`id`, `user_id`, `product_id`, `mobile`, `order_no`, `subject`, `amount`, `status`, `remark`, `pay_time`, `created_at`, `updated_at`, `deleted_at`) VALUES
  (4, 4, 1, '18575670125', '202308317102915300416290816', '会员1个月', '0.01', 2, '{\"days\":30,\"calls\":500,\"name\":\"会员1个月\",\"discount\":10.99}', 1693466990, '2023-08-31 15:29:33', '2023-08-31 15:29:51', NULL),
  (5, 4, 5, '18575670125', '202308317102946758199607296', '100次点卡', '0.30', 2, '{\"days\":0,\"calls\":100,\"name\":\"100次点卡\"}', 1693466990, '2023-08-31 17:34:34', '2023-08-31 17:34:34', NULL),
  (6, 4, 5, '18575670125', '202308317102946843595636736', '100次点卡', '0.03', 2, '{\"days\":0,\"calls\":100,\"name\":\"100次点卡\"}', 1693474722, '2023-08-31 17:34:54', '2023-08-31 17:38:43', NULL),
  (7, 4, 1, '18575670125', '202309017103252664456052736', '会员1个月', '0.01', 2, '{\"days\":30,\"calls\":0,\"name\":\"会员1个月\"}', 1693466990, '2023-09-01 13:50:07', '2023-09-01 13:50:07', NULL),
  (8, 4, 1, '18575670125', '202309017103252894391992320', '会员1个月', '0.01', 2, '{\"days\":30,\"calls\":0,\"name\":\"会员1个月\"}', 1693466990, '2023-09-01 13:51:02', '2023-09-01 13:51:02', NULL),
  (9, 4, 5, '18575670125', '202309017103254657538981888', '100次点卡', '0.03', 2, '{\"days\":0,\"calls\":100,\"name\":\"100次点卡\"}', 1693474722, '2023-09-01 13:58:02', '2023-09-01 13:58:02', NULL),
  (10, 4, 1, '18575670125', '202309017103259375405367296', '会员1个月', '0.01', 2, '{\"days\":30,\"calls\":0,\"name\":\"会员1个月\"}', 1693474722, '2023-09-01 14:16:47', '2023-09-01 14:16:47', NULL),
  (11, 4, 3, '18575670125', '202309017103290730432430080', '会员6个月', '190.00', 2, '{\"days\":180,\"calls\":0,\"name\":\"会员6个月\",\"price\":290,\"discount\":100}', 1693474722, '2023-09-01 16:21:23', '2023-09-01 16:21:23', NULL),
  (12, 4, 4, '18575670125', '202309017103291707520712704', '会员12个月', '380.00', 2, '{\"days\":365,\"calls\":0,\"name\":\"会员12个月\",\"price\":580,\"discount\":200}', 1693466990, '2023-09-01 16:25:16', '2023-09-01 16:25:16', NULL);

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