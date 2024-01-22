ALTER TABLE `chatgpt_mj_jobs` ADD `publish` TINYINT(1) NOT NULL COMMENT '是否发布' AFTER `use_proxy`;
ALTER TABLE `chatgpt_sd_jobs` ADD `publish` TINYINT(1) NOT NULL COMMENT '是否发布' AFTER `progress`;

ALTER TABLE `chatgpt_orders` ADD `trade_no` VARCHAR(60) NOT NULL COMMENT '支付平台交易流水号' AFTER `order_no`;