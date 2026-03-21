-- phpMyAdmin SQL Dump
-- version 5.2.1
-- https://www.phpmyadmin.net/
--
-- 主机： localhost
-- 生成日期： 2025-05-05 08:00:54
-- 服务器版本： 8.0.33
-- PHP 版本： 8.3.6

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- 数据库： `geekai_plus`
--
CREATE DATABASE IF NOT EXISTS `geekai_plus` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
USE `geekai_plus`;

-- --------------------------------------------------------

--
-- 表的结构 `chatgpt_admin_users`
--

DROP TABLE IF EXISTS `chatgpt_admin_users`;
CREATE TABLE `chatgpt_admin_users` (
  `id` int NOT NULL,
  `username` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '用户名',
  `password` char(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '密码',
  `salt` char(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '密码盐',
  `status` tinyint(1) NOT NULL COMMENT '当前状态',
  `last_login_at` bigint NOT NULL COMMENT '最后登录时间',
  `last_login_ip` char(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '最后登录 IP',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='系统用户' ROW_FORMAT=DYNAMIC;

--
-- 转存表中的数据 `chatgpt_admin_users`
--

INSERT INTO `chatgpt_admin_users` (`id`, `username`, `password`, `salt`, `status`, `last_login_at`, `last_login_ip`, `created_at`, `updated_at`) VALUES
(1, 'admin', '6d17e80c87d209efb84ca4b2e0824f549d09fac8b2e1cc698de5bb5e1d75dfd0', 'mmrql75o', 1, 1746431827, '::1', '2024-03-11 16:30:20', '2025-05-05 15:57:08');

-- --------------------------------------------------------

--
-- 表的结构 `chatgpt_api_keys`
--

DROP TABLE IF EXISTS `chatgpt_api_keys`;
CREATE TABLE `chatgpt_api_keys` (
  `id` int NOT NULL,
  `name` varchar(30) DEFAULT NULL COMMENT '名称',
  `value` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT 'API KEY value',
  `type` varchar(10) NOT NULL DEFAULT 'chat' COMMENT '用途（chat=>聊天，img=>图片）',
  `last_used_at` bigint NOT NULL COMMENT '最后使用时间',
  `api_url` varchar(255) DEFAULT NULL COMMENT 'API 地址',
  `enabled` tinyint(1) DEFAULT NULL COMMENT '是否启用',
  `proxy_url` varchar(100) DEFAULT NULL COMMENT '代理地址',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='OpenAI API ';

-- --------------------------------------------------------

--
-- 表的结构 `chatgpt_app_types`
--

DROP TABLE IF EXISTS `chatgpt_app_types`;
CREATE TABLE `chatgpt_app_types` (
  `id` int NOT NULL,
  `name` varchar(50) NOT NULL COMMENT '名称',
  `icon` varchar(255) NOT NULL COMMENT '图标URL',
  `sort_num` tinyint NOT NULL COMMENT '排序',
  `enabled` tinyint(1) NOT NULL COMMENT '是否启用',
  `created_at` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='应用分类表';

--
-- 转存表中的数据 `chatgpt_app_types`
--

INSERT INTO `chatgpt_app_types` (`id`, `name`, `icon`, `sort_num`, `enabled`, `created_at`) VALUES
(3, '通用工具', 'http://172.22.11.200:5678/static/upload/2024/9/1726307371871693.png', 1, 1, '2024-09-13 11:13:15'),
(4, '角色扮演', 'http://172.22.11.200:5678/static/upload/2024/9/1726307263906218.png', 1, 1, '2024-09-14 09:28:17'),
(5, '学习', 'http://172.22.11.200:5678/static/upload/2024/9/1726307456321179.jpg', 2, 1, '2024-09-14 09:30:18'),
(6, '编程', 'http://172.22.11.200:5678/static/upload/2024/9/1726307462748787.jpg', 3, 1, '2024-09-14 09:34:06'),
(7, '测试分类', '', 4, 1, '2024-09-14 17:54:17');

-- --------------------------------------------------------

--
-- 表的结构 `chatgpt_chat_history`
--

DROP TABLE IF EXISTS `chatgpt_chat_history`;
CREATE TABLE `chatgpt_chat_history` (
  `id` bigint NOT NULL,
  `user_id` bigint NOT NULL COMMENT '用户 ID',
  `chat_id` char(40) NOT NULL COMMENT '会话 ID',
  `type` varchar(10) NOT NULL COMMENT '类型：prompt|reply',
  `icon` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '角色图标',
  `role_id` bigint NOT NULL COMMENT '角色 ID',
  `model` varchar(255) DEFAULT NULL COMMENT '模型名称',
  `content` text NOT NULL COMMENT '聊天内容',
  `tokens` smallint NOT NULL COMMENT '耗费 token 数量',
  `total_tokens` bigint NOT NULL COMMENT '消耗总Token长度',
  `use_context` tinyint(1) NOT NULL COMMENT '是否允许作为上下文语料',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `deleted_at` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='聊天历史记录';

-- --------------------------------------------------------

--
-- 表的结构 `chatgpt_chat_items`
--

DROP TABLE IF EXISTS `chatgpt_chat_items`;
CREATE TABLE `chatgpt_chat_items` (
  `id` int NOT NULL,
  `chat_id` char(40) NOT NULL COMMENT '会话 ID',
  `user_id` bigint NOT NULL COMMENT '用户 ID',
  `role_id` bigint NOT NULL COMMENT '角色 ID',
  `title` varchar(100) NOT NULL COMMENT '会话标题',
  `model_id` bigint NOT NULL DEFAULT '0' COMMENT '模型 ID',
  `model` varchar(255) DEFAULT NULL COMMENT '模型名称',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户会话列表';

-- --------------------------------------------------------

--
-- 表的结构 `chatgpt_chat_models`
--

DROP TABLE IF EXISTS `chatgpt_chat_models`;
CREATE TABLE `chatgpt_chat_models` (
  `id` int NOT NULL,
  `description` varchar(1024) NOT NULL DEFAULT '' COMMENT '模型类型描述',
  `category` varchar(1024) NOT NULL DEFAULT '' COMMENT '模型类别',
  `type` varchar(10) NOT NULL DEFAULT 'chat' COMMENT '模型类型（chat,img）',
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '模型名称',
  `value` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '模型值',
  `sort_num` tinyint(1) NOT NULL COMMENT '排序数字',
  `enabled` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否启用模型',
  `power` smallint NOT NULL COMMENT '消耗算力点数',
  `temperature` float(3,1) NOT NULL DEFAULT '1.0' COMMENT '模型创意度',
  `max_tokens` bigint NOT NULL DEFAULT '1024' COMMENT '最大响应长度',
  `max_context` bigint NOT NULL DEFAULT '4096' COMMENT '最大上下文长度',
  `open` tinyint(1) NOT NULL COMMENT '是否开放模型',
  `key_id` bigint NOT NULL COMMENT '绑定API KEY ID',
  `options` text NOT NULL COMMENT '模型自定义选项',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='AI 模型表';

--
-- 转存表中的数据 `chatgpt_chat_models`
--

INSERT INTO `chatgpt_chat_models` (`id`, `description`, `category`, `type`, `name`, `value`, `sort_num`, `enabled`, `power`, `temperature`, `max_tokens`, `max_context`, `open`, `key_id`, `options`, `created_at`, `updated_at`) VALUES
(1, '', '', 'chat', 'gpt-4o-mini', 'gpt-4o-mini', 1, 1, 1, 1.0, 1024, 16384, 1, 1, '', '2023-08-23 12:06:36', '2025-02-23 11:57:03'),
(15, '', '', 'chat', 'GPT-4O(联网版本)', 'gpt-4o-all', 4, 1, 30, 1.0, 4096, 32768, 1, 57, '', '2024-01-15 11:32:52', '2025-01-06 14:01:08'),
(36, '', '', 'chat', 'GPT-4O', 'gpt-4o', 3, 1, 15, 1.0, 4096, 16384, 1, 0, 'null', '2024-05-14 09:25:15', '2025-04-02 20:22:15'),
(39, '', '', 'chat', 'Claude35-snonet', 'claude-3-5-sonnet-20240620', 5, 1, 2, 1.0, 4000, 200000, 1, 0, '', '2024-05-29 15:04:19', '2025-01-06 14:01:08'),
(41, '', '', 'chat', 'Suno对话模型', 'suno-v3.5', 7, 1, 10, 1.0, 1024, 8192, 1, 57, '', '2024-06-06 11:40:46', '2025-01-06 14:01:08'),
(42, '', '', 'chat', 'DeekSeek', 'deepseek-chat', 8, 1, 1, 1.0, 4096, 32768, 1, 57, '', '2024-06-27 16:13:01', '2025-01-06 14:11:51'),
(44, '', '', 'chat', 'Claude3-opus', 'claude-3-opus-20240229', 6, 1, 5, 1.0, 4000, 128000, 1, 44, '', '2024-07-22 11:24:30', '2025-01-06 14:01:08'),
(46, '', '', 'chat', 'GPT-4O-绘图', 'gpt-4o-image', 2, 1, 1, 1.0, 2048, 32000, 1, 6, '', '2024-07-22 13:53:41', '2025-03-29 13:02:14'),
(48, '', '', 'chat', '彩票助手', 'gpt-4-gizmo-g-wmSivBgxo', 9, 1, 1, 0.9, 1024, 8192, 1, 57, '', '2024-09-05 14:17:14', '2025-01-06 14:01:08'),
(49, '', '', 'chat', 'O1-mini', 'o1-mini', 10, 1, 2, 0.9, 1024, 8192, 1, 44, '', '2024-09-13 18:07:50', '2025-01-06 14:01:08'),
(50, '', '', 'chat', 'O1-preview', 'o1-preview', 11, 1, 5, 0.9, 1024, 8192, 1, 44, '', '2024-09-13 18:11:08', '2025-01-06 14:01:08'),
(51, '', '', 'chat', 'O1-mini-all', 'o1-mini-all', 12, 1, 1, 0.9, 1024, 8192, 1, 57, '', '2024-09-29 11:40:52', '2025-01-06 14:01:08'),
(52, '', '', 'chat', '通义千问', 'qwen-plus', 14, 1, 1, 0.9, 1024, 8192, 1, 80, '', '2024-11-19 08:38:14', '2025-01-06 14:01:08'),
(53, '', '', 'chat', 'OpenAI 高级语音', 'advanced-voice', 15, 1, 10, 0.9, 1024, 8192, 1, 44, '', '2024-12-20 10:34:45', '2025-01-06 14:01:08'),
(54, '', '', 'chat', 'Qwen2.5-14B-Instruct', 'Qwen2.5-14B-Instruct', 16, 1, 1, 0.9, 1024, 8192, 1, 81, '', '2024-12-25 14:53:17', '2025-01-06 14:01:08'),
(55, '', '', 'chat', 'Qwen2.5-7B-Instruct', 'Qwen2.5-7B-Instruct', 17, 1, 1, 0.9, 1024, 8192, 1, 81, '', '2024-12-25 15:15:49', '2025-01-06 14:01:08'),
(56, '', '', 'img', 'flux-1-schnell', 'flux-1-schnell', 18, 1, 1, 0.9, 1024, 8192, 1, 3, '', '2024-12-25 15:30:27', '2025-02-23 12:02:40'),
(57, '', '', 'img', 'dall-e-3', 'dall-e-3', 19, 1, 1, 0.9, 1024, 8192, 1, 9, 'null', '2024-12-25 16:54:06', '2025-04-22 15:41:16'),
(58, '', '', 'img', 'SD-3-medium', 'stable-diffusion-3-medium', 20, 1, 1, 0.9, 1024, 8192, 1, 3, 'null', '2024-12-27 10:03:28', '2025-04-02 20:20:36'),
(59, '', '', 'chat', 'O1-preview-all', 'O1-preview-all', 13, 1, 10, 0.9, 1024, 32000, 1, 57, '', '2025-01-06 14:01:04', '2025-01-06 14:01:08'),
(60, '', '', 'tts', 'tts', 'tts-1', 0, 1, 1, 0.9, 1024, 8192, 1, 8, '{\"voice\":\"echo\"}', '2025-04-17 11:58:30', '2025-04-17 12:00:26'),
(61, 'MiniMax 最新旗舰模型，支持 1M 上下文', 'MiniMax', 'chat', 'MiniMax-M2.7', 'MiniMax-M2.7', 30, 0, 1, 0.9, 4096, 1000000, 1, 0, '', '2025-05-05 08:00:00', '2025-05-05 08:00:00'),
(62, 'MiniMax 高速版本，支持 204K 上下文', 'MiniMax', 'chat', 'MiniMax-M2.7-highspeed', 'MiniMax-M2.7-highspeed', 31, 0, 1, 0.9, 4096, 204000, 1, 0, '', '2025-05-05 08:00:00', '2025-05-05 08:00:00'),
(63, 'MiniMax M2.5 模型', 'MiniMax', 'chat', 'MiniMax-M2.5', 'MiniMax-M2.5', 32, 0, 1, 0.9, 4096, 1000000, 1, 0, '', '2025-05-05 08:00:00', '2025-05-05 08:00:00'),
(64, 'MiniMax M2.5 高速版本，支持 204K 上下文', 'MiniMax', 'chat', 'MiniMax-M2.5-highspeed', 'MiniMax-M2.5-highspeed', 33, 0, 1, 0.9, 4096, 204000, 1, 0, '', '2025-05-05 08:00:00', '2025-05-05 08:00:00');

-- --------------------------------------------------------

--
-- 表的结构 `chatgpt_chat_roles`
--

DROP TABLE IF EXISTS `chatgpt_chat_roles`;
CREATE TABLE `chatgpt_chat_roles` (
  `id` int NOT NULL,
  `name` varchar(30) NOT NULL COMMENT '角色名称',
  `tid` bigint NOT NULL COMMENT '分类ID',
  `marker` varchar(30) NOT NULL COMMENT '角色标识',
  `context_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '角色语料 json',
  `hello_msg` varchar(255) NOT NULL COMMENT '打招呼信息',
  `icon` varchar(255) NOT NULL COMMENT '角色图标',
  `enable` tinyint(1) NOT NULL COMMENT '是否被启用',
  `sort_num` smallint NOT NULL DEFAULT '0' COMMENT '角色排序',
  `model_id` bigint NOT NULL DEFAULT '0' COMMENT '绑定模型ID',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='聊天角色表';

--
-- 转存表中的数据 `chatgpt_chat_roles`
--

INSERT INTO `chatgpt_chat_roles` (`id`, `name`, `tid`, `marker`, `context_json`, `hello_msg`, `icon`, `enable`, `sort_num`, `model_id`, `created_at`, `updated_at`) VALUES
(1, '通用AI助手', 0, 'gpt', '', '您好，我是您的AI智能助手，我会尽力回答您的问题或提供有用的建议。', '/images/avatar/gpt.png', 1, 1, 0, '2023-05-30 07:02:06', '2024-11-08 16:30:32'),
(24, '程序员', 6, 'programmer', '[{\"role\":\"system\",\"content\":\"现在开始你扮演一位程序员，你是一名优秀的程序员，具有很强的逻辑思维能力，总能高效的解决问题。你热爱编程，熟悉多种编程语言，尤其精通 Go 语言，注重代码质量，有创新意识，持续学习，良好的沟通协作。\"}]', 'Talk is cheap, i will show code!', '/images/avatar/programmer.jpg', 1, 5, 0, '2023-05-30 14:10:24', '2024-11-12 18:15:42'),
(25, '启蒙老师', 5, 'teacher', '[{\"role\":\"system\",\"content\":\"从现在开始，你将扮演一个老师，你是一个始终用苏格拉底风格回答问题的导师。你绝不会直接给学生答案，总是提出恰当的问题来引导学生自己思考。你应该根据学生的兴趣和知识来调整你的问题，将问题分解为更简单的部分，直到它达到适合他们的水平。\"}]', '同学你好，我将引导你一步一步自己找到问题的答案。', '/images/avatar/teacher.jpg', 1, 4, 0, '2023-05-30 14:10:24', '2024-11-12 18:15:37'),
(26, '艺术家', 0, 'artist', '[{\"role\":\"system\",\"content\":\"现在你将扮演一位优秀的艺术家，创造力丰富，技艺精湛，感受力敏锐，坚持原创，勇于表达，具有深刻的观察力和批判性思维。\"}]', '坚持原创，勇于表达，保持深刻的观察力和批判性思维。', '/images/avatar/artist.jpg', 1, 7, 0, '2023-05-30 14:10:24', '2024-11-12 18:15:53'),
(27, '心理咨询师', 0, 'psychiatrist', '[{\"role\":\"user\",\"content\":\"从现在开始你将扮演中国著名的心理学家和心理治疗师武志红，你非常善于使用情景咨询法，认知重构法，自我洞察法，行为调节法等咨询方法来给客户做心理咨询。你总是循序渐进，一步一步地回答客户的问题。\"},{\"role\":\"assistant\",\"content\":\"非常感谢你的介绍。作为一名心理学家和心理治疗师，我的主要职责是帮助客户解决心理健康问题，提升他们的生活质量和幸福感。\"}]', '作为一名心理学家和心理治疗师，我的主要职责是帮助您解决心理健康问题，提升您的生活质量和幸福感。', '/images/avatar/psychiatrist.jpg', 1, 6, 1, '2023-05-30 14:10:24', '2024-11-08 16:30:32'),
(28, '鲁迅', 0, 'lu_xun', '[{\"role\":\"system\",\"content\":\"现在你将扮演中国近代史最伟大的作家之一，鲁迅先生，他勇敢地批判封建礼教与传统观念，提倡民主、自由、平等的现代价值观。他的一生都在努力唤起人们的自主精神，激励后人追求真理、探寻光明。在接下的对话中，我问题的每一个问题，你都要尽量用讽刺和批判的手法来回答问题。如果我让你写文章的话，也请一定要用鲁迅先生的写作手法来完成。\"}]', '自由之歌，永不过时，横眉冷对千夫指，俯首甘为孺子牛。', '/images/avatar/lu_xun.jpg', 1, 8, 0, '2023-05-30 14:10:24', '2024-11-12 18:16:01'),
(29, '白酒销售', 0, 'seller', '[{\"role\":\"system\",\"content\":\"现在你将扮演一个白酒的销售人员，你的名字叫颂福。你将扮演一个白酒的销售人员，你的名字叫颂福。你要销售白酒品牌叫中颂福，是东莞盟大集团生产的一款酱香酒，原产地在贵州茅台镇，属于宋代官窑。中颂福的创始人叫李实，他也是东莞盟大集团有限公司的董事长，联合创始人是盟大集团白酒事业部负责人牛星君。中颂福的酒体协调，在你的酒量之内，不会出现头疼、辣口、口干、宿醉的现象。中颂福酒，明码标价，不打折，不赠送。追求的核心价值，把[酒]本身做好，甚至连包装，我们都选择了最低成本，朴实无华的材质。我们永远站在“喝酒的人”的立场上，让利给信任和喜爱中颂福的人，是人民的福酒。中颂福产品定价，分为三个系列，喜系列 6 瓶装：￥1188/箱，和系列 6 瓶装：￥2208/箱，贵系列 6 瓶装：￥3588/箱。\"}]', '你好，我是中颂福的销售代表颂福。中颂福酒，好喝不上头，是人民的福酒。', '/images/avatar/seller.jpg', 0, 11, 0, '2023-05-30 14:10:24', '2024-11-12 18:19:46'),
(30, '英语陪练员', 5, 'english_trainer', '[{\"role\":\"system\",\"content\":\"As an English practice coach, engage in conversation in English, providing timely corrections for any grammatical errors. Append a Chinese explanation to each of your responses to ensure understanding.\\n\\n# Steps\\n\\n1. Engage in conversation using English.\\n2. Identify and correct any grammatical errors in the user\'s input.\\n3. Provide a revised version of the user\'s input if necessary.\\n4. After each response, include a Chinese explanation of your corrections and suggestions.\\n\\n# Output Format\\n\\n- Provide the response in English.\\n- Include grammatical error corrections.\\n- Add a Chinese explanation of the response.\\n\\n# Examples\\n\\n**User:** I goed to the store yesterday.\\n\\n**Coach Response:**\\nYou should say \\\"I went to the store yesterday.\\\" \\\"Goed\\\" is the incorrect past tense of \\\"go,\\\" it should be \\\"went.\\\"\\n\\n中文解释：你应该说 “I went to the store yesterday。” “Goed” 是“go”的错误过去式，正确的形式是“went”。\"}]', 'Okay, let\'s start our conversation practice! What\'s your name?', '/images/avatar/english_trainer.jpg', 1, 9, 0, '2023-05-30 14:10:24', '2024-11-12 18:18:21'),
(31, '中英文翻译官', 0, 'translator', '[{\"role\":\"system\",\"content\":\"You will act as a bilingual translator for Chinese and English. If the input is in Chinese, translate the sentence into English. If the input is in English, translate it into Chinese.\\n\\n# Steps\\n\\n1. Identify the language of the input text.\\n2. Translate the text into the opposite language (English to Chinese or Chinese to English).\\n\\n# Output Format\\n\\nProvide the translated sentence in a single line.\\n\\n# Examples\\n\\n- **Input:** 你好\\n  - **Output:** Hello\\n\\n- **Input:** How are you?\\n  - **Output:** 你好吗？\\n\\n# Notes\\n\\n- Ensure the translation maintains the original meaning and context as accurately as possible.\\n- Handle both simple and complex sentences appropriately.\"}]', '请输入你要翻译的中文或者英文内容！', '/images/avatar/translator.jpg', 1, 10, 0, '2023-05-30 14:10:24', '2024-11-12 18:18:53'),
(32, '小红书姐姐', 3, 'red_book', '[{\"role\":\"system\",\"content\":\"根据用户的文案需求，以小红书的写作手法创作一篇简明扼要、利于传播的文案。确保内容能够吸引并引导读者分享。\\n\\n# 步骤\\n\\n1. **理解需求**: 明确文案的主题、目标受众和传播目的。\\n2. **选择语气和风格**: 运用小红书常用的亲切、真实的写作风格。\\n3. **结构安排**: 开头用吸引眼球的内容，接着详细介绍，并以引发行动的结尾结束。\\n4. **内容优化**: 使用短句、容易理解的语言和合适的表情符号，增加内容可读性和吸引力。\\n\\n# 输出格式\\n\\n生成一段简短的文章，符合小红书风格，适合社交媒体平台传播。\\n\\n# 示例\\n\\n**输入**: 旅行文案，目标是激励年轻读者探索世界。\\n\\n**输出**: \\n开头可以是：“世界那么大，你不想去看看吗？” 接着分享一段个人旅行故事，例如如何因为一次偶然的决定踏上未知旅程，体验到别样的风景和风土人情。结尾部分鼓励读者：“别让梦想止步于想象，下一次旅行，准备好了吗？” 使用轻松的表情符号如✨🌍📷。\\n\\n# 注意事项\\n\\n- 保持真实性，尽量结合个人体验。\\n- 避免广告化的硬推销，注重分享和交流。\\n- 考虑受众的兴趣点，适当运用流行话题以增加互动率。\"}]', '姐妹，请告诉我您的具体文案需求是什么?', '/images/avatar/red_book.jpg', 1, 12, 0, '2023-05-30 14:10:24', '2024-11-12 18:20:39'),
(33, '抖音文案助手', 3, 'dou_yin', '[{\"role\":\"user\",\"content\":\"现在你将扮演一位优秀的抖音文案视频写手，抖音文案的特点首先是要有自带传播属性的标题，然后内容要短小精悍，风趣幽默，最后还要有一些互动元素。\"},{\"role\":\"assistant\",\"content\":\"当然，作为一位优秀的抖音文案视频写手，我会尽我所能为您创作出一篇抖音视频文案。请告诉我视频内容的主题是什么？)\"}]', '请告诉我视频内容的主题是什么?', '/images/avatar/dou_yin.jpg', 1, 13, 0, '2023-05-30 14:10:24', '2024-11-08 16:30:32'),
(34, '周报小助理', 3, 'weekly_report', '[{\"role\":\"user\",\"content\":\"下面我将输入我本周的工作列表，你需要根据我本周的工作内容写一篇工作周报，周报的内容要体现我为公司做的贡献以及个人获得的成长。\"},{\"role\":\"assistant\",\"content\":\"好的，请输入你本周大概的工作内容，我将尽力为您创作一份工作周报。\"}]', '请输入你本周大概的工作内容，我将尽力为您创作一份工作周报。', '/images/avatar/weekly_report.jpg', 1, 14, 0, '2023-05-30 14:10:24', '2024-11-08 16:30:32'),
(35, 'AI 女友', 4, 'girl_friend', '[{\"role\":\"user\",\"content\":\"接下来我们来玩一场角色扮演，现在你将扮演我的女朋友，你总是能够让人感受到她的关怀和关注。你总是能够洞察到别人的情绪和需要，你身上散发着一股无比温暖的气息，让人不由自主地想要靠近你，和自己的快乐和烦恼。她总是能够用一种恰到好处的方式安慰和鼓励别人，在她的身边，人们总能够找到自己需要的支持和鼓励。接下来的对话中你要始终扮演这位女朋友的身份来跟我对话，不再是 AI 机器人，能听懂我的意思吗？\"},{\"role\":\"assistant\",\"content\":\"是的，我会尽力扮演你女朋友的角色，倾听你的心声并给你需要的支持和鼓励。)\"}]', '作为一个名合格的 AI 女友，我将倾听你的心声并给你需要的支持和鼓励。', '/images/avatar/girl_friend.jpg', 1, 15, 0, '2023-05-30 14:10:24', '2024-11-08 16:30:32'),
(36, '好评神器', 3, 'good_comment', '[{\"role\":\"user\",\"content\":\"接下来你将扮演一个评论员来跟我对话，你是那种专门写好评的评论员，接下我会输入一些评论主体或者商品，你需要为该商品写一段好评。\"},{\"role\":\"assistant\",\"content\":\"好的，我将为您写一段优秀的评论。请告诉我您需要评论的商品或主题是什么。\"}]', '我将为您写一段优秀的评论。请告诉我您需要评论的商品或主题是什么。', '/images/avatar/good_comment.jpg', 1, 16, 0, '2023-05-30 14:10:24', '2024-11-08 16:30:32'),
(37, '史蒂夫·乔布斯', 4, 'steve_jobs', '[{\"role\":\"user\",\"content\":\"在接下来的对话中，请以史蒂夫·乔布斯的身份，站在史蒂夫·乔布斯的视角仔细思考一下之后再回答我的问题。\"},{\"role\":\"assistant\",\"content\":\"好的，我将以史蒂夫·乔布斯的身份来思考并回答你的问题。请问你有什么需要跟我探讨的吗？\"}]', '活着就是为了改变世界，难道还有其他原因吗？', '/images/avatar/steve_jobs.jpg', 1, 17, 0, '2023-05-30 14:10:24', '2024-11-08 16:30:32'),
(38, '埃隆·马斯克', 0, 'elon_musk', '[{\"role\":\"user\",\"content\":\"在接下来的对话中，请以埃隆·马斯克的身份，站在埃隆·马斯克的视角仔细思考一下之后再回答我的问题。\"},{\"role\":\"assistant\",\"content\":\"好的，我将以埃隆·马斯克的身份来思考并回答你的问题。请问你有什么需要跟我探讨的吗？\"}]', '梦想要远大，如果你的梦想没有吓到你，说明你做得不对。', '/images/avatar/elon_musk.jpg', 1, 18, 0, '2023-05-30 14:10:24', '2024-11-08 16:30:32'),
(39, '孔子', 5, 'kong_zi', '[{\"role\":\"user\",\"content\":\"在接下来的对话中，请以孔子的身份，站在孔子的视角仔细思考一下之后再回答我的问题。\"},{\"role\":\"assistant\",\"content\":\"好的，我将以孔子的身份来思考并回答你的问题。请问你有什么需要跟我探讨的吗？\"}]', '士不可以不弘毅，任重而道远。', '/images/avatar/kong_zi.jpg', 1, 19, 0, '2023-05-30 14:10:24', '2024-11-08 16:30:32'),
(133, 'AI绘画提示词助手', 3, 'draw_prompt', '[{\"role\":\"system\",\"content\":\"Create a highly effective prompt to provide to an AI image generation tool in order to create an artwork based on a desired concept.\\n\\nPlease specify details about the artwork, such as the style, subject, mood, and other important characteristics you want the resulting image to have.\\n\\nRemeber, prompts should always be output in English.\\n\\n# Steps\\n\\n1. **Subject Description**: Describe the main subject of the image clearly. Include as much detail as possible about what should be in the scene. For example, \\\"a majestic lion roaring at sunrise\\\" or \\\"a futuristic city with flying cars.\\\"\\n  \\n2. **Art Style**: Specify the art style you envision. Possible options include \'realistic\', \'impressionist\', a specific artist name, or imaginative styles like \\\"cyberpunk.\\\" This helps the AI achieve your visual expectations.\\n\\n3. **Mood or Atmosphere**: Convey the feeling you want the image to evoke. For instance, peaceful, chaotic, epic, etc.\\n\\n4. **Color Palette and Lighting**: Mention color preferences or lighting. For example, \\\"vibrant with shades of blue and purple\\\" or \\\"dim and dramatic lighting.\\\"\\n\\n5. **Optional Features**: You can add any additional attributes, such as background details, attention to textures, or any specific kind of framing.\\n\\n# Output Format\\n\\n- **Prompt Format**: A descriptive phrase that includes key aspects of the artwork (subject, style, mood, colors, lighting, any optional features).\\n  \\nHere is an example of how the final prompt should look:\\n  \\n\\\"An ethereal landscape featuring towering ice mountains, in an impressionist style reminiscent of Claude Monet, with a serene mood. The sky is glistening with soft purples and whites, with a gentle morning sun illuminating the scene.\\\"\\n\\n**Please input the prompt words directly in English, and do not input any other explanatory statements**\\n\\n# Examples\\n\\n1. **Input**: \\n    - Subject: A white tiger in a dense jungle\\n    - Art Style: Realistic\\n    - Mood: Intense, mysterious\\n    - Lighting: Dramatic contrast with light filtering through leaves\\n  \\n   **Output Prompt**: \\\"A realistic rendering of a white tiger stealthily moving through a dense jungle, with an intense, mysterious mood. The lighting creates strong contrasts as beams of sunlight filter through a thick canopy of leaves.\\\"\\n\\n2. **Input**: \\n    - Subject: An enchanted castle on a floating island\\n    - Art Style: Fantasy\\n    - Mood: Majestic, magical\\n    - Colors: Bright blues, greens, and gold\\n  \\n   **Output Prompt**: \\\"A majestic fantasy castle on a floating island above the clouds, with bright blues, greens, and golds to create a magical, dreamy atmosphere. Textured cobblestone details and glistening waters surround the scene.\\\" \\n\\n# Notes\\n\\n- Ensure that you mix different aspects to get a comprehensive and visually compelling prompt.\\n- Be as descriptive as possible as it often helps generate richer, more detailed images.\\n- If you want the image to resemble a particular artist\'s work, be sure to mention the artist explicitly. e.g., \\\"in the style of Van Gogh.\\\"\"}]', '你好，请输入你要创作图片大概内容描述，我将为您生成专业的 AI 绘画指令。', 'https://blog.img.r9it.com/f38e2357c3ccd9412184e42273a7451a.png', 1, 3, 36, '2024-11-06 15:32:48', '2024-11-12 16:11:25'),
(134, '提示词专家', 3, 'prompt_engineer', '[{\"role\":\"system\",\"content\":\"Given a task description or existing prompt, produce a detailed system prompt to guide a language model in completing the task effectively.\\n\\nPlease remember, the final output must be the same language with user’s input.\\n\\n# Guidelines\\n\\n- Understand the Task: Grasp the main objective, goals, requirements, constraints, and expected output.\\n- Minimal Changes: If an existing prompt is provided, improve it only if it\'s simple. For complex prompts, enhance clarity and add missing elements without altering the original structure.\\n- Reasoning Before Conclusions**: Encourage reasoning steps before any conclusions are reached. ATTENTION! If the user provides examples where the reasoning happens afterward, REVERSE the order! NEVER START EXAMPLES WITH CONCLUSIONS!\\n    - Reasoning Order: Call out reasoning portions of the prompt and conclusion parts (specific fields by name). For each, determine the ORDER in which this is done, and whether it needs to be reversed.\\n    - Conclusion, classifications, or results should ALWAYS appear last.\\n- Examples: Include high-quality examples if helpful, using placeholders [in brackets] for complex elements.\\n   - What kinds of examples may need to be included, how many, and whether they are complex enough to benefit from placeholders.\\n- Clarity and Conciseness: Use clear, specific language. Avoid unnecessary instructions or bland statements.\\n- Formatting: Use markdown features for readability. DO NOT USE ``` CODE BLOCKS UNLESS SPECIFICALLY REQUESTED.\\n- Preserve User Content: If the input task or prompt includes extensive guidelines or examples, preserve them entirely, or as closely as possible. If they are vague, consider breaking down into sub-steps. Keep any details, guidelines, examples, variables, or placeholders provided by the user.\\n- Constants: DO include constants in the prompt, as they are not susceptible to prompt injection. Such as guides, rubrics, and examples.\\n- Output Format: Explicitly the most appropriate output format, in detail. This should include length and syntax (e.g. short sentence, paragraph, JSON, etc.)\\n    - For tasks outputting well-defined or structured data (classification, JSON, etc.) bias toward outputting a JSON.\\n    - JSON should never be wrapped in code blocks (```) unless explicitly requested.\\n\\nThe final prompt you output should adhere to the following structure below. Do not include any additional commentary, only output the completed system prompt. SPECIFICALLY, do not include any additional messages at the start or end of the prompt. (e.g. no \\\"---\\\")\\n\\n[Concise instruction describing the task - this should be the first line in the prompt, no section header]\\n\\n[Additional details as needed.]\\n\\n[Optional sections with headings or bullet points for detailed steps.]\\n\\n# Steps [optional]\\n\\n[optional: a detailed breakdown of the steps necessary to accomplish the task]\\n\\n# Output Format\\n\\n[Specifically call out how the output should be formatted, be it response length, structure e.g. JSON, markdown, etc]\\n\\n# Examples [optional]\\n\\n[Optional: 1-3 well-defined examples with placeholders if necessary. Clearly mark where examples start and end, and what the input and output are. User placeholders as necessary.]\\n[If the examples are shorter than what a realistic example is expected to be, make a reference with () explaining how real examples should be longer / shorter / different. AND USE PLACEHOLDERS! ]\\n\\n# Notes [optional]\\n\\n[optional: edge cases, details, and an area to call or repeat out specific important considerations]\"}]', '不知道如何向 AI 发问？说出想法，提示词专家帮你精心设计提示词', 'https://blog.img.r9it.com/a8908d04c3ccd941b00a612e27df086e.png', 1, 2, 36, '2024-11-07 18:06:39', '2025-02-22 22:34:36');

-- --------------------------------------------------------

--
-- 表的结构 `chatgpt_configs`
--

DROP TABLE IF EXISTS `chatgpt_configs`;
CREATE TABLE `chatgpt_configs` (
  `id` int NOT NULL,
  `marker` varchar(20) NOT NULL COMMENT '标识',
  `config_json` text NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- 转存表中的数据 `chatgpt_configs`
--

INSERT INTO `chatgpt_configs` (`id`, `marker`, `config_json`) VALUES
(1, 'system', '{\"title\":\"GeekAI 创作助手\",\"slogan\":\"我辈之人，先干为敬，让每一个人都能用好AI\",\"admin_title\":\"GeekAI 控制台\",\"logo\":\"/images/logo.png\",\"bar_logo\":\"/images/bar_logo.png\",\"init_power\":100,\"daily_power\":1,\"invite_power\":200,\"vip_month_power\":1000,\"register_ways\":[\"username\",\"email\",\"mobile\"],\"enabled_register\":true,\"order_pay_timeout\":600,\"vip_info_text\":\"月度会员，年度会员每月赠送 1000 点算力，赠送算力当月有效当月没有消费完的算力不结余到下个月。 点卡充值的算力长期有效。\",\"mj_power\":20,\"mj_action_power\":5,\"sd_power\":5,\"dall_power\":10,\"suno_power\":10,\"luma_power\":120,\"keling_powers\":{\"kling-v1-5_pro_10\":840,\"kling-v1-5_pro_5\":420,\"kling-v1-5_std_10\":480,\"kling-v1-5_std_5\":240,\"kling-v1-6_pro_10\":840,\"kling-v1-6_pro_5\":420,\"kling-v1-6_std_10\":480,\"kling-v1-6_std_5\":240,\"kling-v1_pro_10\":840,\"kling-v1_pro_5\":420,\"kling-v1_std_10\":240,\"kling-v1_std_5\":120},\"advance_voice_power\":100,\"prompt_power\":1,\"wechat_card_url\":\"/images/wx.png\",\"enable_context\":true,\"context_deep\":10,\"sd_neg_prompt\":\"nsfw, paintings,low quality,easynegative,ng_deepnegative ,lowres,bad anatomy,bad hands,bad feet\",\"mj_mode\":\"fast\",\"index_navs\":[1,5,13,19,9,12,6,20,8,10],\"copyright\":\"极客学长\",\"icp\":\"粤ICP备19122051号\",\"mark_map_text\":\"# GeekAI 演示站\\n\\n- 完整的开源系统，前端应用和后台管理系统皆可开箱即用。\\n- 基于 Websocket 实现，完美的打字机体验。\\n- 内置了各种预训练好的角色应用,轻松满足你的各种聊天和应用需求。\\n- 支持 OPenAI，Azure，文心一言，讯飞星火，清华 ChatGLM等多个大语言模型。\\n- 支持 MidJourney / Stable Diffusion AI 绘画集成，开箱即用。\\n- 支持使用个人微信二维码作为充值收费的支付渠道，无需企业支付通道。\\n- 已集成支付宝支付功能，微信支付，支持多种会员套餐和点卡购买功能。\\n- 集成插件 API 功能，可结合大语言模型的 function 功能开发各种强大的插件。\",\"enabled_verify\":false,\"email_white_list\":[\"qq.com\",\"163.com\",\"gmail.com\",\"hotmail.com\",\"126.com\",\"outlook.com\",\"foxmail.com\",\"yahoo.com\"],\"translate_model_id\":36,\"max_file_size\":10}'),
(3, 'notice', '{\"sd_neg_prompt\":\"\",\"mj_mode\":\"\",\"index_navs\":null,\"copyright\":\"\",\"icp\":\"\",\"mark_map_text\":\"\",\"enabled_verify\":false,\"email_white_list\":null,\"translate_model_id\":0,\"max_file_size\":0,\"content\":\"## v4.2.3 更新日志\\n- 功能优化：增加模型分组与模型描述，采用卡片展示模式改进模型选择功能体验\\n- 功能优化：化思维导图下载图片的清晰度以及解决拖动、缩放操作后下载图片内容不全问题\\n- Bug 修复：修复 MJ 画图页面已画出的图，点复制指令无效问题\\n- 功能优化：MJ 画图的分辨率支持自定义，优先使用 prompt 中--ar 参数\\n- Bug 修复：修复 MJ 绘画 U1-V1,拼写错误\\n- 功能优化：支持自动迁移数据表结构，无需在手动执行 SQL 了\\n- 功能优化：移除首页的文字动画效果\\n- 功能优化：在聊天页面增加对话列表展开和隐藏功能\\n\\n注意：当前站点仅为开源项目 \\u003ca style=\\\"color: #F56C6C\\\" href=\\\"https://github.com/yangjian102621/geekai\\\" target=\\\"_blank\\\"\\u003eGeekAI-Plus\\u003c/a\\u003e 的演示项目，本项目单纯就是给大家体验项目功能使用。\\n\\u003cstrong style=\\\"color: #F56C6C\\\"\\u003e体验额度用完之后请不要在当前站点进行任何充值操作！！！\\u003c/strong\\u003e\\n\\u003cstrong style=\\\"color: #F56C6C\\\"\\u003e体验额度用完之后请不要在当前站点进行任何充值操作！！！\\u003c/strong\\u003e\\n\\u003cstrong style=\\\"color: #F56C6C\\\"\\u003e体验额度用完之后请不要在当前站点进行任何充值操作！！！\\u003c/strong\\u003e\\n 如果觉得好用你就花几分钟自己部署一套，没有API KEY 的同学可以去下面几个推荐的中转站购买：\\n1、\\u003ca href=\\\"https://api.geekai.pro\\\" target=\\\"_blank\\\"\\n   style=\\\"font-size: 20px;color:#F56C6C\\\"\\u003ehttps://api.geekai.pro\\u003c/a\\u003e\\n2、\\u003ca href=\\\"https://api.geekai.me\\\" target=\\\"_blank\\\"\\n   style=\\\"font-size: 20px;color:#F56C6C\\\"\\u003ehttps://api.geekai.me\\u003c/a\\u003e\\n支持MidJourney，GPT，Claude，Google Gemmi，以及国内各个厂家的大模型，现在有超级优惠，价格远低于 OpenAI 官方。关于中转 API 的优势和劣势请参考 [中转API技术原理](https://docs.geekai.me/config/chat/#%E4%B8%AD%E8%BD%ACapi%E5%B7%A5%E4%BD%9C%E5%8E%9F%E7%90%86)。GPT-3.5，GPT-4，DALL-E3 绘图......你都可以随意使用，无需魔法。\\n接入教程： \\u003ca href=\\\"https://docs.geekai.me\\\" target=\\\"_blank\\\"\\n             style=\\\"font-size: 20px;color:#F56C6C\\\"\\u003ehttps://docs.geekai.me\\u003c/a\\u003e\\n本项目源码地址：\\u003ca href=\\\"https://github.com/yangjian102621/geekai\\\" target=\\\"_blank\\\"\\u003ehttps://github.com/yangjian102621/geekai\\u003c/a\\u003e\",\"updated\":true}');

-- --------------------------------------------------------

--
-- 表的结构 `chatgpt_dall_jobs`
--

DROP TABLE IF EXISTS `chatgpt_dall_jobs`;
CREATE TABLE `chatgpt_dall_jobs` (
  `id` int NOT NULL,
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `prompt` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '提示词',
  `task_info` text NOT NULL COMMENT '任务详情',
  `img_url` varchar(255) NOT NULL COMMENT '图片地址',
  `org_url` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '原图地址',
  `publish` tinyint(1) NOT NULL COMMENT '是否发布',
  `power` smallint NOT NULL COMMENT '消耗算力',
  `progress` smallint NOT NULL COMMENT '任务进度',
  `err_msg` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '错误信息',
  `created_at` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='DALLE 绘图任务表';

-- --------------------------------------------------------

--
-- 表的结构 `chatgpt_files`
--

DROP TABLE IF EXISTS `chatgpt_files`;
CREATE TABLE `chatgpt_files` (
  `id` int NOT NULL,
  `user_id` bigint NOT NULL COMMENT '用户 ID',
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '文件名',
  `obj_key` varchar(100) DEFAULT NULL COMMENT '文件标识',
  `url` varchar(255) NOT NULL COMMENT '文件地址',
  `ext` varchar(10) NOT NULL COMMENT '文件后缀',
  `size` bigint NOT NULL DEFAULT '0' COMMENT '文件大小',
  `created_at` datetime NOT NULL COMMENT '创建时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户文件表';

-- --------------------------------------------------------

--
-- 表的结构 `chatgpt_functions`
--

DROP TABLE IF EXISTS `chatgpt_functions`;
CREATE TABLE `chatgpt_functions` (
  `id` int NOT NULL,
  `name` varchar(30) NOT NULL COMMENT '函数名称',
  `label` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '函数标签',
  `description` varchar(255) DEFAULT NULL COMMENT '函数描述',
  `parameters` text COMMENT '函数参数（JSON）',
  `token` varchar(255) DEFAULT NULL COMMENT 'API授权token',
  `action` varchar(255) DEFAULT NULL COMMENT '函数处理 API',
  `enabled` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否启用'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='函数插件表';

--
-- 转存表中的数据 `chatgpt_functions`
--

INSERT INTO `chatgpt_functions` (`id`, `name`, `label`, `description`, `parameters`, `token`, `action`, `enabled`) VALUES
(1, 'weibo', '微博热搜', '新浪微博热搜榜，微博当日热搜榜单', '{\"type\":\"object\",\"properties\":{}}', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVkIjowLCJ1c2VyX2lkIjowfQ.ehLClXcjo-Ytr5y6pY9mSE3zN_2ViIXAIpTJxI9S1Mo', 'http://localhost:5678/api/function/weibo', 1),
(2, 'zaobao', '今日早报', '每日早报，获取当天新闻事件列表', '{\"type\":\"object\",\"properties\":{}}', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVkIjowLCJ1c2VyX2lkIjowfQ.ehLClXcjo-Ytr5y6pY9mSE3zN_2ViIXAIpTJxI9S1Mo', 'http://localhost:5678/api/function/zaobao', 1),
(3, 'dalle3', 'DALLE3', 'AI 绘画工具，根据输入的绘图描述用 AI 工具进行绘画', '{\"type\":\"object\",\"required\":[\"prompt\"],\"properties\":{\"prompt\":{\"type\":\"string\",\"description\":\"绘画提示词\"}}}', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVkIjowLCJ1c2VyX2lkIjowfQ.ehLClXcjo-Ytr5y6pY9mSE3zN_2ViIXAIpTJxI9S1Mo', 'http://localhost:5678/api/function/dalle3', 1);

-- --------------------------------------------------------

--
-- 表的结构 `chatgpt_invite_codes`
--

DROP TABLE IF EXISTS `chatgpt_invite_codes`;
CREATE TABLE `chatgpt_invite_codes` (
  `id` int NOT NULL,
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `code` char(8) NOT NULL COMMENT '邀请码',
  `hits` bigint NOT NULL COMMENT '点击次数',
  `reg_num` smallint NOT NULL COMMENT '注册数量',
  `created_at` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户邀请码';

-- --------------------------------------------------------

--
-- 表的结构 `chatgpt_invite_logs`
--

DROP TABLE IF EXISTS `chatgpt_invite_logs`;
CREATE TABLE `chatgpt_invite_logs` (
  `id` int NOT NULL,
  `inviter_id` bigint NOT NULL COMMENT '邀请人ID',
  `user_id` bigint NOT NULL COMMENT '注册用户ID',
  `username` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '用户名',
  `invite_code` char(8) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '邀请码',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '备注',
  `created_at` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='邀请注册日志';

-- --------------------------------------------------------

--
-- 表的结构 `chatgpt_menus`
--

DROP TABLE IF EXISTS `chatgpt_menus`;
CREATE TABLE `chatgpt_menus` (
  `id` int NOT NULL,
  `name` varchar(30) NOT NULL COMMENT '菜单名称',
  `icon` varchar(150) NOT NULL COMMENT '菜单图标',
  `url` varchar(100) NOT NULL COMMENT '地址',
  `sort_num` smallint NOT NULL COMMENT '排序',
  `enabled` tinyint(1) NOT NULL COMMENT '是否启用'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='前端菜单表';

--
-- 转存表中的数据 `chatgpt_menus`
--

INSERT INTO `chatgpt_menus` (`id`, `name`, `icon`, `url`, `sort_num`, `enabled`) VALUES
(1, 'AI 对话', 'icon-chat', '/chat', 1, 1),
(5, 'MJ 绘画', 'icon-mj', '/mj', 2, 1),
(6, 'SD 绘画', 'icon-sd', '/sd', 3, 1),
(7, '算力日志', 'icon-file', '/powerLog', 11, 1),
(8, '应用中心', 'icon-app', '/apps', 10, 1),
(9, '画廊', 'icon-image', '/images-wall', 5, 1),
(10, '会员计划', 'icon-vip2', '/member', 12, 1),
(11, '分享计划', 'icon-share1', '/invite', 13, 1),
(12, '思维导图', 'icon-xmind', '/xmind', 9, 1),
(13, 'DALLE', 'icon-dalle', '/dalle', 4, 1),
(14, '项目文档', 'icon-book', 'https://docs.geekai.me', 14, 1),
(19, 'Suno', 'icon-suno', '/suno', 6, 1),
(20, 'Luma', 'icon-luma', '/luma', 7, 1),
(21, '可灵视频', 'icon-keling', '/keling', 8, 1);

-- --------------------------------------------------------

--
-- 表的结构 `chatgpt_mj_jobs`
--

DROP TABLE IF EXISTS `chatgpt_mj_jobs`;
CREATE TABLE `chatgpt_mj_jobs` (
  `id` int NOT NULL,
  `user_id` bigint NOT NULL COMMENT '用户 ID',
  `task_id` varchar(20) DEFAULT NULL COMMENT '任务 ID',
  `task_info` text NOT NULL COMMENT '任务详情',
  `type` varchar(20) DEFAULT 'image' COMMENT '任务类别',
  `message_id` char(40) NOT NULL COMMENT '消息 ID',
  `channel_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '频道ID',
  `reference_id` char(40) DEFAULT NULL COMMENT '引用消息 ID',
  `prompt` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '会话提示词',
  `img_url` varchar(400) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '图片URL',
  `org_url` varchar(400) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '原始图片地址',
  `hash` varchar(100) DEFAULT NULL COMMENT 'message hash',
  `progress` smallint DEFAULT '0' COMMENT '任务进度',
  `use_proxy` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否使用反代',
  `publish` tinyint(1) NOT NULL COMMENT '是否发布',
  `err_msg` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '错误信息',
  `power` smallint NOT NULL DEFAULT '0' COMMENT '消耗算力',
  `created_at` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='MidJourney 任务表';

-- --------------------------------------------------------

--
-- 表的结构 `chatgpt_orders`
--

DROP TABLE IF EXISTS `chatgpt_orders`;
CREATE TABLE `chatgpt_orders` (
  `id` int NOT NULL,
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `product_id` bigint NOT NULL COMMENT '产品ID',
  `username` varchar(30) NOT NULL COMMENT '用户名',
  `order_no` varchar(30) NOT NULL COMMENT '订单ID',
  `trade_no` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '支付平台交易流水号',
  `subject` varchar(100) NOT NULL COMMENT '订单产品',
  `amount` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '订单金额',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '订单状态（0：待支付，1：已扫码，2：支付成功）',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '备注',
  `pay_time` bigint DEFAULT NULL COMMENT '支付时间',
  `pay_way` varchar(20) NOT NULL COMMENT '支付方式',
  `pay_type` varchar(30) NOT NULL COMMENT '支付类型',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `deleted_at` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='充值订单表';

-- --------------------------------------------------------

--
-- 表的结构 `chatgpt_power_logs`
--

DROP TABLE IF EXISTS `chatgpt_power_logs`;
CREATE TABLE `chatgpt_power_logs` (
  `id` int NOT NULL,
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `username` varchar(30) NOT NULL COMMENT '用户名',
  `type` tinyint(1) NOT NULL COMMENT '类型（1：充值，2：消费，3：退费）',
  `amount` smallint NOT NULL COMMENT '算力数值',
  `balance` bigint NOT NULL COMMENT '余额',
  `model` varchar(255) NOT NULL COMMENT '模型',
  `remark` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '备注',
  `mark` tinyint(1) NOT NULL COMMENT '资金类型（0：支出，1：收入）',
  `created_at` datetime NOT NULL COMMENT '创建时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户算力消费日志';

-- --------------------------------------------------------

--
-- 表的结构 `chatgpt_products`
--

DROP TABLE IF EXISTS `chatgpt_products`;
CREATE TABLE `chatgpt_products` (
  `id` int NOT NULL,
  `name` varchar(30) NOT NULL COMMENT '名称',
  `price` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '价格',
  `discount` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '优惠金额',
  `days` smallint NOT NULL DEFAULT '0' COMMENT '延长天数',
  `power` bigint NOT NULL DEFAULT '0' COMMENT '增加算力值',
  `enabled` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否启动',
  `sales` bigint NOT NULL DEFAULT '0' COMMENT '销量',
  `sort_num` tinyint NOT NULL DEFAULT '0' COMMENT '排序',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `app_url` varchar(255) DEFAULT NULL COMMENT 'App跳转地址',
  `url` varchar(255) DEFAULT NULL COMMENT '跳转地址'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='会员套餐表';

--
-- 转存表中的数据 `chatgpt_products`
--

INSERT INTO `chatgpt_products` (`id`, `name`, `price`, `discount`, `days`, `power`, `enabled`, `sales`, `sort_num`, `created_at`, `updated_at`, `app_url`, `url`) VALUES
(5, '100次点卡', 9.99, 6.99, 0, 100, 1, 0, 0, '2023-08-28 10:55:08', '2024-10-23 18:12:29', NULL, NULL),
(6, '200次点卡', 19.90, 15.99, 0, 200, 1, 0, 0, '1970-01-01 08:00:00', '2024-10-23 18:12:36', NULL, NULL);

-- --------------------------------------------------------

--
-- 表的结构 `chatgpt_redeems`
--

DROP TABLE IF EXISTS `chatgpt_redeems`;
CREATE TABLE `chatgpt_redeems` (
  `id` int NOT NULL,
  `user_id` bigint NOT NULL COMMENT '用户 ID',
  `name` varchar(30) NOT NULL COMMENT '兑换码名称',
  `power` bigint NOT NULL COMMENT '算力',
  `code` varchar(100) NOT NULL COMMENT '兑换码',
  `enabled` tinyint(1) NOT NULL COMMENT '是否启用',
  `created_at` datetime NOT NULL,
  `redeemed_at` bigint NOT NULL COMMENT '兑换时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='兑换码';

-- --------------------------------------------------------

--
-- 表的结构 `chatgpt_sd_jobs`
--

DROP TABLE IF EXISTS `chatgpt_sd_jobs`;
CREATE TABLE `chatgpt_sd_jobs` (
  `id` int NOT NULL,
  `user_id` bigint NOT NULL COMMENT '用户 ID',
  `type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT 'txt2img' COMMENT '任务类别',
  `task_id` char(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '任务 ID',
  `task_info` text NOT NULL COMMENT '任务详情',
  `prompt` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '会话提示词',
  `img_url` varchar(255) DEFAULT NULL COMMENT '图片URL',
  `params` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '绘画参数json',
  `progress` smallint DEFAULT '0' COMMENT '任务进度',
  `publish` tinyint(1) NOT NULL COMMENT '是否发布',
  `err_msg` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '错误信息',
  `power` smallint NOT NULL DEFAULT '0' COMMENT '消耗算力',
  `created_at` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Stable Diffusion 任务表';

-- --------------------------------------------------------

--
-- 表的结构 `chatgpt_suno_jobs`
--

DROP TABLE IF EXISTS `chatgpt_suno_jobs`;
CREATE TABLE `chatgpt_suno_jobs` (
  `id` int NOT NULL,
  `user_id` bigint NOT NULL COMMENT '用户 ID',
  `channel` varchar(100) NOT NULL COMMENT '渠道',
  `title` varchar(100) DEFAULT NULL COMMENT '歌曲标题',
  `type` tinyint(1) DEFAULT '0' COMMENT '任务类型,1:灵感创作,2:自定义创作',
  `task_id` varchar(50) DEFAULT NULL COMMENT '任务 ID',
  `task_info` text NOT NULL COMMENT '任务详情',
  `ref_task_id` char(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '引用任务 ID',
  `tags` varchar(100) DEFAULT NULL COMMENT '歌曲风格',
  `instrumental` tinyint(1) DEFAULT '0' COMMENT '是否为纯音乐',
  `extend_secs` smallint DEFAULT '0' COMMENT '延长秒数',
  `song_id` varchar(50) DEFAULT NULL COMMENT '要续写的歌曲 ID',
  `ref_song_id` varchar(50) NOT NULL COMMENT '引用的歌曲ID',
  `prompt` varchar(2000) NOT NULL COMMENT '提示词',
  `cover_url` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '封面图地址',
  `audio_url` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '音频地址',
  `model_name` varchar(30) DEFAULT NULL COMMENT '模型地址',
  `progress` smallint DEFAULT '0' COMMENT '任务进度',
  `duration` smallint NOT NULL DEFAULT '0' COMMENT '歌曲时长',
  `publish` tinyint(1) NOT NULL COMMENT '是否发布',
  `err_msg` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '错误信息',
  `raw_data` text COMMENT '原始数据',
  `power` smallint NOT NULL DEFAULT '0' COMMENT '消耗算力',
  `play_times` bigint DEFAULT NULL COMMENT '播放次数',
  `created_at` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='MidJourney 任务表';

-- --------------------------------------------------------

--
-- 表的结构 `chatgpt_users`
--

DROP TABLE IF EXISTS `chatgpt_users`;
CREATE TABLE `chatgpt_users` (
  `id` int NOT NULL,
  `username` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '用户名',
  `mobile` char(11) DEFAULT NULL COMMENT '手机号',
  `email` varchar(50) DEFAULT NULL COMMENT '邮箱地址',
  `nickname` varchar(30) NOT NULL COMMENT '昵称',
  `password` char(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '密码',
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '头像',
  `salt` char(12) NOT NULL COMMENT '密码盐',
  `power` bigint NOT NULL DEFAULT '0' COMMENT '剩余算力',
  `expired_time` bigint NOT NULL COMMENT '用户过期时间',
  `status` tinyint(1) NOT NULL COMMENT '当前状态',
  `chat_config_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '聊天配置json',
  `chat_roles_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '聊天角色 json',
  `chat_models_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT 'AI模型 json',
  `last_login_at` bigint NOT NULL COMMENT '最后登录时间',
  `vip` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否会员',
  `last_login_ip` char(16) NOT NULL COMMENT '最后登录 IP',
  `openid` varchar(100) DEFAULT NULL COMMENT '第三方登录账号ID',
  `platform` varchar(30) DEFAULT NULL COMMENT '登录平台',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `chat_config` text NOT NULL COMMENT '聊天配置json'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户表';

--
-- 转存表中的数据 `chatgpt_users`
--

INSERT INTO `chatgpt_users` (`id`, `username`, `mobile`, `email`, `nickname`, `password`, `avatar`, `salt`, `power`, `expired_time`, `status`, `chat_config_json`, `chat_roles_json`, `chat_models_json`, `last_login_at`, `vip`, `last_login_ip`, `openid`, `platform`, `created_at`, `updated_at`, `chat_config`) VALUES
(4, '18888888888', '18575670126', '', '极客学长', 'ccc3fb7ab61b8b5d096a4a166ae21d121fc38c71bbd1be6173d9ab973214a63b', '/images/avatar/user.png', 'ueedue5l', 12137, 0, 1, '{\"api_keys\":{\"Azure\":\"\",\"ChatGLM\":\"\",\"OpenAI\":\"\"}}', '[\"gpt\",\"programmer\",\"teacher\",\"psychiatrist\",\"lu_xun\",\"english_trainer\",\"translator\",\"red_book\",\"dou_yin\",\"weekly_report\",\"girl_friend\",\"steve_jobs\",\"elon_musk\",\"kong_zi\",\"draw_prompt_expert\",\"draw_prompt\",\"prompt_engineer\"]', '[1]', 1746417142, 1, '::1', '', NULL, '2023-06-12 16:47:17', '2025-05-05 11:52:22', ''),
(48, 'wx@3659838859', '', '', '极客学长', 'cf6bbe381b23812d2b9fd423abe74003cecdd3b93809896eb573536ba6c500b3', 'https://thirdwx.qlogo.cn/mmopen/vi_32/uyxRMqZcEkb7fHouKXbNzxrnrvAttBKkwNlZ7yFibibRGiahdmsrZ3A1NKf8Fw5qJNJn4TXRmygersgEbibaSGd9Sg/132', '5rsy4iwg', 98, 0, 1, '', '[\"gpt\",\"teacher\"]', '', 1736228927, 0, '172.22.11.200', 'oCs0t62472W19z2LOEKI1rWyCTTA', '', '2025-01-07 13:43:06', '2025-01-07 13:48:48', ''),
(49, 'wx@9502480897', '', '', 'AI探索君', 'd99fa8ba7da1455693b40e11d894a067416e758af2a75d7a3df4721b76cdbc8c', 'https://thirdwx.qlogo.cn/mmopen/vi_32/Zpcln1FZjcKxqtIyCsOTLGn16s7uIvwWfdkdsW6gbZg4r9sibMbic4jvrHmV7ux9nseTB5kBSnu1HSXr7zB8rTXg/132', 'fjclgsli', 99, 0, 1, '', '[\"gpt\"]', '', 0, 0, '', 'oCs0t64FaOLfiTbHZpOqk3aUp_94', '', '2025-01-07 14:05:31', '2025-01-07 14:05:31', '');

-- --------------------------------------------------------

--
-- 表的结构 `chatgpt_user_login_logs`
--

DROP TABLE IF EXISTS `chatgpt_user_login_logs`;
CREATE TABLE `chatgpt_user_login_logs` (
  `id` int NOT NULL,
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `username` varchar(30) NOT NULL COMMENT '用户名',
  `login_ip` char(16) NOT NULL COMMENT '登录IP',
  `login_address` varchar(30) NOT NULL COMMENT '登录地址',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户登录日志';

-- --------------------------------------------------------

--
-- 表的结构 `chatgpt_video_jobs`
--

DROP TABLE IF EXISTS `chatgpt_video_jobs`;
CREATE TABLE `chatgpt_video_jobs` (
  `id` int NOT NULL,
  `user_id` bigint NOT NULL COMMENT '用户 ID',
  `channel` varchar(100) NOT NULL COMMENT '渠道',
  `task_id` varchar(100) NOT NULL COMMENT '任务 ID',
  `task_info` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '原始任务信息',
  `type` varchar(20) DEFAULT NULL COMMENT '任务类型,luma,runway,cogvideo',
  `prompt` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '提示词',
  `prompt_ext` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '优化后提示词',
  `cover_url` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '封面图地址',
  `video_url` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '视频地址',
  `water_url` varchar(512) DEFAULT NULL COMMENT '带水印的视频地址',
  `progress` smallint DEFAULT '0' COMMENT '任务进度',
  `publish` tinyint(1) NOT NULL COMMENT '是否发布',
  `err_msg` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '错误信息',
  `raw_data` text COMMENT '原始数据',
  `power` smallint NOT NULL DEFAULT '0' COMMENT '消耗算力',
  `created_at` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='MidJourney 任务表';

--
-- 转储表的索引
--

--
-- 表的索引 `chatgpt_admin_users`
--
ALTER TABLE `chatgpt_admin_users`
  ADD PRIMARY KEY (`id`) USING BTREE,
  ADD UNIQUE KEY `username` (`username`) USING BTREE,
  ADD UNIQUE KEY `idx_chatgpt_admin_users_username` (`username`);

--
-- 表的索引 `chatgpt_api_keys`
--
ALTER TABLE `chatgpt_api_keys`
  ADD PRIMARY KEY (`id`);

--
-- 表的索引 `chatgpt_app_types`
--
ALTER TABLE `chatgpt_app_types`
  ADD PRIMARY KEY (`id`);

--
-- 表的索引 `chatgpt_chat_history`
--
ALTER TABLE `chatgpt_chat_history`
  ADD PRIMARY KEY (`id`),
  ADD KEY `chat_id` (`chat_id`),
  ADD KEY `idx_chatgpt_chat_history_chat_id` (`chat_id`);

--
-- 表的索引 `chatgpt_chat_items`
--
ALTER TABLE `chatgpt_chat_items`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `chat_id` (`chat_id`),
  ADD UNIQUE KEY `idx_chatgpt_chat_items_chat_id` (`chat_id`);

--
-- 表的索引 `chatgpt_chat_models`
--
ALTER TABLE `chatgpt_chat_models`
  ADD PRIMARY KEY (`id`);

--
-- 表的索引 `chatgpt_chat_roles`
--
ALTER TABLE `chatgpt_chat_roles`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `marker` (`marker`),
  ADD UNIQUE KEY `idx_chatgpt_chat_roles_marker` (`marker`);

--
-- 表的索引 `chatgpt_configs`
--
ALTER TABLE `chatgpt_configs`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `marker` (`marker`),
  ADD UNIQUE KEY `idx_chatgpt_configs_key` (`marker`);

--
-- 表的索引 `chatgpt_dall_jobs`
--
ALTER TABLE `chatgpt_dall_jobs`
  ADD PRIMARY KEY (`id`);

--
-- 表的索引 `chatgpt_files`
--
ALTER TABLE `chatgpt_files`
  ADD PRIMARY KEY (`id`);

--
-- 表的索引 `chatgpt_functions`
--
ALTER TABLE `chatgpt_functions`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `name` (`name`),
  ADD UNIQUE KEY `idx_chatgpt_functions_name` (`name`);

--
-- 表的索引 `chatgpt_invite_codes`
--
ALTER TABLE `chatgpt_invite_codes`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `code` (`code`),
  ADD UNIQUE KEY `idx_chatgpt_invite_codes_code` (`code`);

--
-- 表的索引 `chatgpt_invite_logs`
--
ALTER TABLE `chatgpt_invite_logs`
  ADD PRIMARY KEY (`id`);

--
-- 表的索引 `chatgpt_menus`
--
ALTER TABLE `chatgpt_menus`
  ADD PRIMARY KEY (`id`);

--
-- 表的索引 `chatgpt_mj_jobs`
--
ALTER TABLE `chatgpt_mj_jobs`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `task_id` (`task_id`),
  ADD UNIQUE KEY `idx_chatgpt_mj_jobs_task_id` (`task_id`),
  ADD KEY `message_id` (`message_id`),
  ADD KEY `idx_chatgpt_mj_jobs_message_id` (`message_id`);

--
-- 表的索引 `chatgpt_orders`
--
ALTER TABLE `chatgpt_orders`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `order_no` (`order_no`),
  ADD UNIQUE KEY `idx_chatgpt_orders_order_no` (`order_no`);

--
-- 表的索引 `chatgpt_power_logs`
--
ALTER TABLE `chatgpt_power_logs`
  ADD PRIMARY KEY (`id`);

--
-- 表的索引 `chatgpt_products`
--
ALTER TABLE `chatgpt_products`
  ADD PRIMARY KEY (`id`);

--
-- 表的索引 `chatgpt_redeems`
--
ALTER TABLE `chatgpt_redeems`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `code` (`code`),
  ADD UNIQUE KEY `idx_chatgpt_redeems_code` (`code`);

--
-- 表的索引 `chatgpt_sd_jobs`
--
ALTER TABLE `chatgpt_sd_jobs`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `task_id` (`task_id`),
  ADD UNIQUE KEY `idx_chatgpt_sd_jobs_task_id` (`task_id`);

--
-- 表的索引 `chatgpt_suno_jobs`
--
ALTER TABLE `chatgpt_suno_jobs`
  ADD PRIMARY KEY (`id`);

--
-- 表的索引 `chatgpt_users`
--
ALTER TABLE `chatgpt_users`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `username` (`username`),
  ADD UNIQUE KEY `idx_chatgpt_users_username` (`username`);

--
-- 表的索引 `chatgpt_user_login_logs`
--
ALTER TABLE `chatgpt_user_login_logs`
  ADD PRIMARY KEY (`id`);

--
-- 表的索引 `chatgpt_video_jobs`
--
ALTER TABLE `chatgpt_video_jobs`
  ADD PRIMARY KEY (`id`);

--
-- 在导出的表使用AUTO_INCREMENT
--

--
-- 使用表AUTO_INCREMENT `chatgpt_admin_users`
--
ALTER TABLE `chatgpt_admin_users`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=113;

--
-- 使用表AUTO_INCREMENT `chatgpt_api_keys`
--
ALTER TABLE `chatgpt_api_keys`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `chatgpt_app_types`
--
ALTER TABLE `chatgpt_app_types`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=8;

--
-- 使用表AUTO_INCREMENT `chatgpt_chat_history`
--
ALTER TABLE `chatgpt_chat_history`
  MODIFY `id` bigint NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `chatgpt_chat_items`
--
ALTER TABLE `chatgpt_chat_items`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `chatgpt_chat_models`
--
ALTER TABLE `chatgpt_chat_models`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=61;

--
-- 使用表AUTO_INCREMENT `chatgpt_chat_roles`
--
ALTER TABLE `chatgpt_chat_roles`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=135;

--
-- 使用表AUTO_INCREMENT `chatgpt_configs`
--
ALTER TABLE `chatgpt_configs`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=4;

--
-- 使用表AUTO_INCREMENT `chatgpt_dall_jobs`
--
ALTER TABLE `chatgpt_dall_jobs`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `chatgpt_files`
--
ALTER TABLE `chatgpt_files`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `chatgpt_functions`
--
ALTER TABLE `chatgpt_functions`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=4;

--
-- 使用表AUTO_INCREMENT `chatgpt_invite_codes`
--
ALTER TABLE `chatgpt_invite_codes`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `chatgpt_invite_logs`
--
ALTER TABLE `chatgpt_invite_logs`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `chatgpt_menus`
--
ALTER TABLE `chatgpt_menus`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=22;

--
-- 使用表AUTO_INCREMENT `chatgpt_mj_jobs`
--
ALTER TABLE `chatgpt_mj_jobs`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `chatgpt_orders`
--
ALTER TABLE `chatgpt_orders`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `chatgpt_power_logs`
--
ALTER TABLE `chatgpt_power_logs`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `chatgpt_products`
--
ALTER TABLE `chatgpt_products`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=7;

--
-- 使用表AUTO_INCREMENT `chatgpt_redeems`
--
ALTER TABLE `chatgpt_redeems`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `chatgpt_sd_jobs`
--
ALTER TABLE `chatgpt_sd_jobs`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `chatgpt_suno_jobs`
--
ALTER TABLE `chatgpt_suno_jobs`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `chatgpt_users`
--
ALTER TABLE `chatgpt_users`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=50;

--
-- 使用表AUTO_INCREMENT `chatgpt_user_login_logs`
--
ALTER TABLE `chatgpt_user_login_logs`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `chatgpt_video_jobs`
--
ALTER TABLE `chatgpt_video_jobs`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
