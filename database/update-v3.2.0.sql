
CREATE TABLE `chatgpt_invite_logs` (
                                      `id` int NOT NULL,
                                      `inviter_id` int NOT NULL COMMENT '邀请人ID',
                                      `user_id` int NOT NULL COMMENT '注册用户ID',
                                      `username` varchar(11) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '用户名',
                                      `invite_code` char(8) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '邀请码',
                                      `calls` smallint NOT NULL COMMENT '奖励对话次数',
                                      `created_at` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='邀请注册日志';
ALTER TABLE `chatgpt_invite_logs` ADD PRIMARY KEY (`id`);
ALTER TABLE `chatgpt_invite_logs` MODIFY `id` int NOT NULL AUTO_INCREMENT;
ALTER TABLE `chatgpt_invite_logs` CHANGE `calls` `reward_json` TEXT NOT NULL COMMENT '邀请奖励';


CREATE TABLE `chatgpt_invite_codes` (
                                        `id` int NOT NULL,
                                        `user_id` int NOT NULL COMMENT '用户ID',
                                        `code` char(8) NOT NULL COMMENT '邀请码',
                                        `hits` int NOT NULL COMMENT '点击次数',
                                        `reg_num` smallint NOT NULL COMMENT '注册数量',
                                        `created_at` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户邀请码';
ALTER TABLE `chatgpt_invite_codes` ADD PRIMARY KEY (`id`);
ALTER TABLE `chatgpt_invite_codes` MODIFY `id` int NOT NULL AUTO_INCREMENT;
ALTER TABLE `chatgpt_invite_codes` ADD UNIQUE(`code`);

ALTER TABLE `chatgpt_api_keys` ADD `type` VARCHAR(10) NOT NULL DEFAULT 'chat' COMMENT '用途（chat=>聊天，img=>图片）' AFTER `value`;