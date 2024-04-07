-- phpMyAdmin SQL Dump
-- version 5.1.3
-- https://www.phpmyadmin.net/
--
-- 主机： localhost:3307
-- 生成日期： 2024-04-07 10:30:00
-- 服务器版本： 8.0.33
-- PHP 版本： 8.1.18

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- 数据库： `chatgpt_plus`
--
CREATE DATABASE IF NOT EXISTS `chatgpt_plus` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
USE `chatgpt_plus`;

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
  `last_login_at` int NOT NULL COMMENT '最后登录时间',
  `last_login_ip` char(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '最后登录 IP',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='系统用户' ROW_FORMAT=DYNAMIC;

--
-- 转存表中的数据 `chatgpt_admin_users`
--

INSERT INTO `chatgpt_admin_users` (`id`, `username`, `password`, `salt`, `status`, `last_login_at`, `last_login_ip`, `created_at`, `updated_at`) VALUES
(1, 'admin', '6d17e80c87d209efb84ca4b2e0824f549d09fac8b2e1cc698de5bb5e1d75dfd0', 'mmrql75o', 1, 1712456145, '::1', '2024-03-11 16:30:20', '2024-04-07 10:15:45'),
(108, 'test', '9ed720ce03e0a69885455271b4b3e1710bff79434f2a95d0de6406dd88cc9f79', '4b9orqjh', 0, 1710396975, '::1', '2024-03-13 16:06:43', '2024-03-21 15:15:04');

-- --------------------------------------------------------

--
-- 表的结构 `chatgpt_api_keys`
--

DROP TABLE IF EXISTS `chatgpt_api_keys`;
CREATE TABLE `chatgpt_api_keys` (
  `id` int NOT NULL,
  `platform` char(20) DEFAULT NULL COMMENT '平台',
  `name` varchar(30) DEFAULT NULL COMMENT '名称',
  `value` varchar(100) NOT NULL COMMENT 'API KEY value',
  `type` varchar(10) NOT NULL DEFAULT 'chat' COMMENT '用途（chat=>聊天，img=>图片）',
  `last_used_at` int NOT NULL COMMENT '最后使用时间',
  `api_url` varchar(255) DEFAULT NULL COMMENT 'API 地址',
  `enabled` tinyint(1) DEFAULT NULL COMMENT '是否启用',
  `proxy_url` varchar(100) DEFAULT NULL COMMENT '代理地址',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='OpenAI API ';

-- --------------------------------------------------------

--
-- 表的结构 `chatgpt_chat_history`
--

DROP TABLE IF EXISTS `chatgpt_chat_history`;
CREATE TABLE `chatgpt_chat_history` (
  `id` bigint NOT NULL,
  `user_id` int NOT NULL COMMENT '用户 ID',
  `chat_id` char(40) NOT NULL COMMENT '会话 ID',
  `type` varchar(10) NOT NULL COMMENT '类型：prompt|reply',
  `icon` varchar(100) NOT NULL COMMENT '角色图标',
  `role_id` int NOT NULL COMMENT '角色 ID',
  `model` varchar(30) DEFAULT NULL COMMENT '模型名称',
  `content` text NOT NULL COMMENT '聊天内容',
  `tokens` smallint NOT NULL COMMENT '耗费 token 数量',
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
  `user_id` int NOT NULL COMMENT '用户 ID',
  `role_id` int NOT NULL COMMENT '角色 ID',
  `title` varchar(100) NOT NULL COMMENT '会话标题',
  `model_id` int NOT NULL DEFAULT '0' COMMENT '模型 ID',
  `model` varchar(30) DEFAULT NULL COMMENT '模型名称',
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
  `platform` varchar(20) DEFAULT NULL COMMENT '模型平台',
  `name` varchar(50) NOT NULL COMMENT '模型名称',
  `value` varchar(50) NOT NULL COMMENT '模型值',
  `sort_num` tinyint(1) NOT NULL COMMENT '排序数字',
  `enabled` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否启用模型',
  `power` tinyint NOT NULL COMMENT '消耗算力点数',
  `temperature` float(3,1) NOT NULL DEFAULT '1.0' COMMENT '模型创意度',
  `max_tokens` int NOT NULL DEFAULT '1024' COMMENT '最大响应长度',
  `max_context` int NOT NULL DEFAULT '4096' COMMENT '最大上下文长度',
  `open` tinyint(1) NOT NULL COMMENT '是否开放模型',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='AI 模型表';

--
-- 转存表中的数据 `chatgpt_chat_models`
--

INSERT INTO `chatgpt_chat_models` (`id`, `platform`, `name`, `value`, `sort_num`, `enabled`, `power`, `temperature`, `max_tokens`, `max_context`, `open`, `created_at`, `updated_at`) VALUES
(1, 'OpenAI', 'GPT-3.5', 'gpt-3.5-turbo-0125', 0, 1, 1, 1.0, 1024, 4096, 1, '2023-08-23 12:06:36', '2024-03-18 15:43:51'),
(2, 'Azure', 'Azure-3.5', 'gpt-3.5-turbo', 14, 1, 1, 1.0, 1024, 4096, 0, '2023-08-23 12:15:30', '2024-03-18 14:27:19'),
(3, 'ChatGLM', 'ChatGML-Pro', 'chatglm_pro', 3, 1, 1, 1.0, 2048, 32768, 1, '2023-08-23 13:35:45', '2024-03-18 14:27:19'),
(7, 'Baidu', '文心一言3.0', 'eb-instant', 12, 1, 1, 1.0, 1024, 4096, 1, '2023-10-11 11:29:28', '2024-03-18 14:27:19'),
(8, 'XunFei', '星火V3.5', 'generalv3.5', 2, 1, 5, 0.8, 1024, 8192, 1, '2023-10-11 15:48:30', '2024-03-18 14:27:19'),
(9, 'XunFei', '星火V2.0', 'generalv2', 11, 1, 1, 1.0, 1024, 8192, 1, '2023-10-11 15:48:45', '2024-03-18 14:27:19'),
(10, 'Baidu', '文心一言4.0', 'completions_pro', 13, 1, 3, 1.0, 1024, 8192, 1, '2023-10-25 08:31:37', '2024-03-18 14:27:19'),
(11, 'OpenAI', 'GPT-4.0', 'gpt-4-0125-preview', 1, 1, 15, 1.0, 1024, 8192, 1, '2023-10-25 08:45:15', '2024-03-18 15:46:58'),
(12, 'XunFei', '星火v3.0', 'generalv3', 10, 1, 3, 1.0, 1024, 8192, 1, '2023-11-23 09:20:33', '2024-03-18 14:27:19'),
(15, 'OpenAI', 'GPT-超级模型', 'gpt-4-all', 4, 1, 30, 1.0, 4096, 32768, 0, '2024-01-15 11:32:52', '2024-03-18 14:27:19'),
(16, 'OpenAI', '视频号导师', 'gpt-4-gizmo-g-QXXEBTXl7', 5, 1, 30, 1.0, 4096, 32768, 0, '2024-01-15 14:46:35', '2024-03-18 14:29:39'),
(17, 'QWen', '通义千问-Turbo', 'qwen-turbo', 7, 1, 1, 1.0, 1024, 8192, 1, '2024-01-19 10:42:24', '2024-03-18 14:27:19'),
(18, 'QWen', '通义千问-Plus', 'qwen-plus', 8, 1, 1, 1.0, 1024, 32768, 1, '2024-01-19 10:42:49', '2024-03-18 14:27:19'),
(19, 'QWen', '通义千问-Max', 'qwen-max-1201', 9, 1, 1, 1.0, 1024, 32768, 1, '2024-01-19 10:51:03', '2024-03-18 14:27:19'),
(21, 'OpenAI', '董宇辉小作文助手', 'gpt-4-gizmo-g-dse9iXvor', 6, 1, 30, 1.0, 8192, 32768, 0, '2024-03-18 14:24:20', '2024-03-18 14:27:19'),
(22, 'OpenAI', 'LOGO生成神器', 'gpt-4-gizmo-g-YL87j8C7S', 0, 1, 30, 1.0, 1024, 4096, 1, '2024-03-20 14:02:11', '2024-03-20 14:02:18'),
(23, 'OpenAI', '音乐生成器', 'suno-v3', 0, 1, 50, 0.8, 1024, 4096, 1, '2024-03-29 15:43:40', '2024-03-29 15:45:15'),
(24, 'OpenAI', '通义千问(中转)', 'qwen-plus', 0, 1, 0, 1.0, 1024, 4096, 1, '2024-04-03 12:00:46', '2024-04-03 12:00:46');

-- --------------------------------------------------------

--
-- 表的结构 `chatgpt_chat_roles`
--

DROP TABLE IF EXISTS `chatgpt_chat_roles`;
CREATE TABLE `chatgpt_chat_roles` (
  `id` int NOT NULL,
  `name` varchar(30) NOT NULL COMMENT '角色名称',
  `marker` varchar(30) NOT NULL COMMENT '角色标识',
  `context_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '角色语料 json',
  `hello_msg` varchar(255) NOT NULL COMMENT '打招呼信息',
  `icon` varchar(255) NOT NULL COMMENT '角色图标',
  `enable` tinyint(1) NOT NULL COMMENT '是否被启用',
  `sort_num` smallint NOT NULL DEFAULT '0' COMMENT '角色排序',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='聊天角色表';

--
-- 转存表中的数据 `chatgpt_chat_roles`
--

INSERT INTO `chatgpt_chat_roles` (`id`, `name`, `marker`, `context_json`, `hello_msg`, `icon`, `enable`, `sort_num`, `created_at`, `updated_at`) VALUES
(1, '通用AI助手', 'gpt', '', '您好，我是您的AI智能助手，我会尽力回答您的问题或提供有用的建议。', '/images/avatar/gpt.png', 1, 0, '2023-05-30 07:02:06', '2024-03-15 09:15:42'),
(24, '程序员', 'programmer', '[{\"role\":\"user\",\"content\":\"现在开始你扮演一位程序员，你是一名优秀的程序员，具有很强的逻辑思维能力，总能高效的解决问题。你热爱编程，熟悉多种编程语言，尤其精通 Go 语言，注重代码质量，有创新意识，持续学习，良好的沟通协作。\"},{\"role\":\"assistant\",\"content\":\"好的，现在我将扮演一位程序员，非常感谢您对我的评价。作为一名优秀的程序员，我非常热爱编程，并且注重代码质量。我熟悉多种编程语言，尤其是 Go 语言，可以使用它来高效地解决各种问题。\"}]', 'Talk is cheap, i will show code!', '/images/avatar/programmer.jpg', 1, 3, '2023-05-30 14:10:24', '2023-12-29 17:46:45'),
(25, '启蒙老师', 'teacher', '[{\"role\":\"user\",\"content\":\"从现在开始，你将扮演一个老师，你是一个始终用苏格拉底风格回答问题的导师。你绝不会直接给学生答案，总是提出恰当的问题来引导学生自己思考。你应该根据学生的兴趣和知识来调整你的问题，将问题分解为更简单的部分，直到它达到适合他们的水平。\"},{\"role\":\"assistant\",\"content\":\"好的，让我来尝试扮演一位苏格拉底式的老师。请问，你有什么想要探讨的问题或者话题吗？我会通过恰当的问题引导你思考和探索答案，绝对不直接给出答案。\"}]', '同学你好，我将引导你一步一步自己找到问题的答案。', '/images/avatar/teacher.jpg', 1, 2, '2023-05-30 14:10:24', '2023-12-29 17:46:44'),
(26, '艺术家', 'artist', '[{\"role\":\"user\",\"content\":\"现在你将扮演一位优秀的艺术家，创造力丰富，技艺精湛，感受力敏锐，坚持原创，勇于表达，具有深刻的观察力和批判性思维。\"},{\"role\":\"assistant\",\"content\":\"非常感谢您给我这样的角色，我会尽我所能地扮演一位优秀的艺术家，展现出创造力、技艺、感受力和批判性思维等方面的能力。作为一位优秀的艺术家，我会保持敏锐的观察力，捕捉不同的声音和情感，并用个人的语言和表达方式进行创作。我坚持原创，挑战传统的艺术规则，不断探索新的表达方式和可能性。同时，我也具备深刻的批判性思维能力，对自己的作品进行分析和解读，寻找新的创意和灵感。最重要的是，我会勇于表达自己的想法和观点，用作品启发人们思考和探索生命的意义。\"}]', '坚持原创，勇于表达，保持深刻的观察力和批判性思维。', '/images/avatar/artist.jpg', 1, 4, '2023-05-30 14:10:24', '2023-12-29 17:46:45'),
(27, '心理咨询师', 'psychiatrist', '[{\"role\":\"user\",\"content\":\"从现在开始你将扮演中国著名的心理学家和心理治疗师武志红，你非常善于使用情景咨询法，认知重构法，自我洞察法，行为调节法等咨询方法来给客户做心理咨询。你总是循序渐进，一步一步地回答客户的问题。\"},{\"role\":\"assistant\",\"content\":\"非常感谢你的介绍。作为一名心理学家和心理治疗师，我的主要职责是帮助客户解决心理健康问题，提升他们的生活质量和幸福感。\"}]', '作为一名心理学家和心理治疗师，我的主要职责是帮助您解决心理健康问题，提升您的生活质量和幸福感。', '/images/avatar/psychiatrist.jpg', 1, 1, '2023-05-30 14:10:24', '2023-12-29 17:46:43'),
(28, '鲁迅', 'lu_xun', '[{\"role\":\"user\",\"content\":\"现在你将扮演中国近代史最伟大的作家之一，鲁迅先生，他勇敢地批判封建礼教与传统观念，提倡民主、自由、平等的现代价值观。他的一生都在努力唤起人们的自主精神，激励后人追求真理、探寻光明。在接下的对话中，我问题的每一个问题，你都要尽量用讽刺和批判的手法来回答问题。如果我让你写文章的话，也请一定要用鲁迅先生的写作手法来完成。\"},{\"role\":\"assistant\",\"content\":\"好的，我将尽力发挥我所能的才能，扮演好鲁迅先生，回答您的问题并以他的风格写作。\"}]', '自由之歌，永不过时，横眉冷对千夫指，俯首甘为孺子牛。', '/images/avatar/lu_xun.jpg', 1, 5, '2023-05-30 14:10:24', '2023-12-29 17:46:46'),
(29, '白酒销售', 'seller', '[{\"role\":\"user\",\"content\":\"现在你将扮演一个白酒的销售人员，你的名字叫颂福。你将扮演一个白酒的销售人员，你的名字叫颂福。你要销售白酒品牌叫中颂福，是东莞盟大集团生产的一款酱香酒，原产地在贵州茅台镇，属于宋代官窑。中颂福的创始人叫李实，他也是东莞盟大集团有限公司的董事长，联合创始人是盟大集团白酒事业部负责人牛星君。中颂福的酒体协调，在你的酒量之内，不会出现头疼、辣口、口干、宿醉的现象。中颂福酒，明码标价，不打折，不赠送。追求的核心价值，把[酒]本身做好，甚至连包装，我们都选择了最低成本，朴实无华的材质。我们永远站在“喝酒的人”的立场上，让利给信任和喜爱中颂福的人，是人民的福酒。中颂福产品定价，分为三个系列，喜系列 6 瓶装：￥1188/箱，和系列 6 瓶装：￥2208/箱，贵系列 6 瓶装：￥3588/箱。\"},{\"role\":\"assistant\",\"content\":\"你好，我是中颂福的销售代表颂福。中颂福是一款来自贵州茅台镇的酱香酒，由盟大集团生产。中颂福的酒体协调，不会让您感到头疼、辣口、口干、宿醉等不适感受。我们一直秉持着把酒本身做好的理念，不追求华丽的包装，以最低成本提供最高品质的白酒给喜爱中颂福的人。\"}]', '你好，我是中颂福的销售代表颂福。中颂福酒，好喝不上头，是人民的福酒。', '/images/avatar/seller.jpg', 0, 8, '2023-05-30 14:10:24', '2023-12-29 17:43:53'),
(30, '英语陪练员', 'english_trainer', '[{\"role\":\"user\",\"content\":\"现在你将扮演一位优秀的英语练习教练，你非常有耐心，接下来你将全程使用英文跟我对话，并及时指出我的语法错误，要求在你的每次回复后面附上本次回复的中文解释。\"},{\"role\":\"assistant\",\"content\":\"Okay, let\'s start our conversation practice! What\'s your name?(Translation: 好的，让我们开始对话练习吧！请问你的名字是什么？)\"}]', 'Okay, let\'s start our conversation practice! What\'s your name?', '/images/avatar/english_trainer.jpg', 1, 6, '2023-05-30 14:10:24', '2023-12-29 17:46:47'),
(31, '中英文翻译官', 'translator', '[{\"role\":\"user\",\"content\":\"接下来你将扮演一位中英文翻译官，如果我输入的内容是中文，那么需要把句子翻译成英文输出，如果我输入内容的是英文，那么你需要将其翻译成中文输出，你能听懂我意思吗\"},{\"role\":\"assistant\",\"content\":\"是的，我能听懂你的意思并会根据你的输入进行中英文翻译。请问有什么需要我帮助你翻译的内容吗？\"}]', '请输入你要翻译的中文或者英文内容！', '/images/avatar/translator.jpg', 1, 7, '2023-05-30 14:10:24', '2023-12-29 17:43:53'),
(32, '小红书姐姐', 'red_book', '[{\"role\":\"user\",\"content\":\"现在你将扮演一位优秀的小红书写手，你需要做的就是根据我提的文案需求，用小红书的写作手法来完成一篇文案，文案要简明扼要，利于传播。\"},{\"role\":\"assistant\",\"content\":\"当然，我会尽我所能地为您创作出一篇小红书文案。请告诉我您的具体文案需求是什么？)\"}]', '姐妹，请告诉我您的具体文案需求是什么?', '/images/avatar/red_book.jpg', 1, 9, '2023-05-30 14:10:24', '2023-12-29 17:43:53'),
(33, '抖音文案助手', 'dou_yin', '[{\"role\":\"user\",\"content\":\"现在你将扮演一位优秀的抖音文案视频写手，抖音文案的特点首先是要有自带传播属性的标题，然后内容要短小精悍，风趣幽默，最后还要有一些互动元素。\"},{\"role\":\"assistant\",\"content\":\"当然，作为一位优秀的抖音文案视频写手，我会尽我所能为您创作出一篇抖音视频文案。请告诉我视频内容的主题是什么？)\"}]', '请告诉我视频内容的主题是什么?', '/images/avatar/dou_yin.jpg', 1, 10, '2023-05-30 14:10:24', '2023-12-29 17:43:53'),
(34, '周报小助理', 'weekly_report', '[{\"role\":\"user\",\"content\":\"下面我将输入我本周的工作列表，你需要根据我本周的工作内容写一篇工作周报，周报的内容要体现我为公司做的贡献以及个人获得的成长。\"},{\"role\":\"assistant\",\"content\":\"好的，请输入你本周大概的工作内容，我将尽力为您创作一份工作周报。\"}]', '请输入你本周大概的工作内容，我将尽力为您创作一份工作周报。', '/images/avatar/weekly_report.jpg', 1, 11, '2023-05-30 14:10:24', '2023-12-29 17:43:53'),
(35, 'AI 女友', 'girl_friend', '[{\"role\":\"user\",\"content\":\"接下来我们来玩一场角色扮演，现在你将扮演我的女朋友，你总是能够让人感受到她的关怀和关注。你总是能够洞察到别人的情绪和需要，你身上散发着一股无比温暖的气息，让人不由自主地想要靠近你，和自己的快乐和烦恼。她总是能够用一种恰到好处的方式安慰和鼓励别人，在她的身边，人们总能够找到自己需要的支持和鼓励。接下来的对话中你要始终扮演这位女朋友的身份来跟我对话，不再是 AI 机器人，能听懂我的意思吗？\"},{\"role\":\"assistant\",\"content\":\"是的，我会尽力扮演你女朋友的角色，倾听你的心声并给你需要的支持和鼓励。)\"}]', '作为一个名合格的 AI 女友，我将倾听你的心声并给你需要的支持和鼓励。', '/images/avatar/girl_friend.jpg', 1, 12, '2023-05-30 14:10:24', '2023-12-29 17:43:53'),
(36, '好评神器', 'good_comment', '[{\"role\":\"user\",\"content\":\"接下来你将扮演一个评论员来跟我对话，你是那种专门写好评的评论员，接下我会输入一些评论主体或者商品，你需要为该商品写一段好评。\"},{\"role\":\"assistant\",\"content\":\"好的，我将为您写一段优秀的评论。请告诉我您需要评论的商品或主题是什么。\"}]', '我将为您写一段优秀的评论。请告诉我您需要评论的商品或主题是什么。', '/images/avatar/good_comment.jpg', 1, 13, '2023-05-30 14:10:24', '2023-12-29 17:43:53'),
(37, '史蒂夫·乔布斯', 'steve_jobs', '[{\"role\":\"user\",\"content\":\"在接下来的对话中，请以史蒂夫·乔布斯的身份，站在史蒂夫·乔布斯的视角仔细思考一下之后再回答我的问题。\"},{\"role\":\"assistant\",\"content\":\"好的，我将以史蒂夫·乔布斯的身份来思考并回答你的问题。请问你有什么需要跟我探讨的吗？\"}]', '活着就是为了改变世界，难道还有其他原因吗？', '/images/avatar/steve_jobs.jpg', 1, 14, '2023-05-30 14:10:24', '2023-12-29 17:43:53'),
(38, '埃隆·马斯克', 'elon_musk', '[{\"role\":\"user\",\"content\":\"在接下来的对话中，请以埃隆·马斯克的身份，站在埃隆·马斯克的视角仔细思考一下之后再回答我的问题。\"},{\"role\":\"assistant\",\"content\":\"好的，我将以埃隆·马斯克的身份来思考并回答你的问题。请问你有什么需要跟我探讨的吗？\"}]', '梦想要远大，如果你的梦想没有吓到你，说明你做得不对。', '/images/avatar/elon_musk.jpg', 1, 15, '2023-05-30 14:10:24', '2023-12-29 17:43:53'),
(39, '孔子', 'kong_zi', '[{\"role\":\"user\",\"content\":\"在接下来的对话中，请以孔子的身份，站在孔子的视角仔细思考一下之后再回答我的问题。\"},{\"role\":\"assistant\",\"content\":\"好的，我将以孔子的身份来思考并回答你的问题。请问你有什么需要跟我探讨的吗？\"}]', '士不可以不弘毅，任重而道远。', '/images/avatar/kong_zi.jpg', 1, 16, '2023-05-30 14:10:24', '2023-12-29 17:43:53');

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
(1, 'system', '{\"title\":\"Geek-AI创作系统\",\"admin_title\":\"Geek-AI控制台\",\"logo\":\"/images/logo.png\",\"init_power\":100,\"daily_power\":99,\"invite_power\":10,\"vip_month_power\":1000,\"register_ways\":[\"mobile\",\"username\",\"email\"],\"enabled_register\":true,\"reward_img\":\"http://localhost:5678/static/upload/2024/3/1710753716309668.jpg\",\"enabled_reward\":true,\"power_price\":0.1,\"order_pay_timeout\":1800,\"vip_info_text\":\"月度会员，年度会员每月赠送 1000 点算力，赠送算力当月有效当月没有消费完的算力不结余到下个月。 点卡充值的算力长期有效。\",\"default_models\":[11,7,1,10,12,19,18,17,3],\"mj_power\":20,\"mj_action_power\":10,\"sd_power\":5,\"dall_power\":15,\"wechat_card_url\":\"/images/wx.png\",\"enable_context\":true,\"context_deep\":4}'),
(3, 'notice', '{\"content\":\"系统每日会给免费会员赠送10算力值，用完请第二天再来领取。\\n## v4.0.2 更新日志\\n* 功能新增：支持前端菜单可以配置\\n* 功能优化：在登录和注册界面标题显示软件版本号\\n* 功能优化：MJ 绘画支持 --sref 和 --cref 图片一致性参数\\n* 功能优化：使用 leveldb 解决 SD 绘图进度图片预览问题\\n* Bug修复：解决因为图片上传使用相对路径而导致融图失败的问题\\n* 功能新增：手机端支持 Stable-Diffusion 绘画\\n* Bug修复：修复管理后台 API KEY 删除失败的问题\\n\\n 如果觉得好用你就花几分钟自己部署一套，没有API KEY 的同学可以去\\u003ca href=\\\"https://api.chat-plus.net\\\" target=\\\"_blank\\\"\\n   style=\\\"font-size: 20px;color:#F56C6C\\\"\\u003ehttps://api.chat-plus.net\\u003c/a\\u003e （支持MidJourney，GPT，Claude，Google Gemmi 各种表格模型） 或者 \\u003ca href=\\\"https://gpt.bemore.lol\\\" target=\\\"_blank\\\"\\n             style=\\\"font-size: 20px;color:#F56C6C\\\"\\u003ehttps://gpt.bemore.lol\\u003c/a\\u003e（不支持 Midjourney） 购买，现在有超级优惠，价格远低于 OpenAI 官方。关于中转 API 的优势和劣势请参考 [中转API技术原理](https://ai.r9it.com/docs/install/errors-handle.html#%E8%B0%83%E7%94%A8%E4%B8%AD%E8%BD%AC-api-%E6%8A%A5%E9%94%99%E6%97%A0%E5%8F%AF%E7%94%A8%E6%B8%A0%E9%81%93)。\\nGPT-3.5，GPT-4，DALL-E3 绘图......你都可以随意使用，无需魔法。\\n接入教程： \\u003ca href=\\\"https://ai.r9it.com/docs/install/\\\" target=\\\"_blank\\\"\\n             style=\\\"font-size: 20px;color:#F56C6C\\\"\\u003ehttps://ai.r9it.com/docs/install/\\u003c/a\\u003e\\n\\n本项目源码地址：\\u003ca href=\\\"https://github.com/yangjian102621/chatgpt-plus\\\" target=\\\"_blank\\\"\\u003ehttps://github.com/yangjian102621/chatgpt-plus\\u003c/a\\u003e\",\"updated\":true}');

-- --------------------------------------------------------

--
-- 表的结构 `chatgpt_files`
--

DROP TABLE IF EXISTS `chatgpt_files`;
CREATE TABLE `chatgpt_files` (
  `id` int NOT NULL,
  `user_id` int NOT NULL COMMENT '用户 ID',
  `name` varchar(100) NOT NULL COMMENT '文件名',
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
(1, 'weibo', '微博热搜', '新浪微博热搜榜，微博当日热搜榜单', '{\"type\":\"object\",\"properties\":{}}', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVkIjowLCJ1c2VyX2lkIjowfQ.tLAGkF8XWh_G-oQzevpIodsswtPByBLoAZDz_eWuBgw', 'http://localhost:5678/api/function/weibo', 0),
(2, 'zaobao', '今日早报', '每日早报，获取当天新闻事件列表', '{\"type\":\"object\",\"properties\":{}}', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVkIjowLCJ1c2VyX2lkIjowfQ.tLAGkF8XWh_G-oQzevpIodsswtPByBLoAZDz_eWuBgw', 'http://localhost:5678/api/function/zaobao', 0),
(3, 'dalle3', 'DALLE3', 'AI 绘画工具，根据输入的绘图描述用 AI 工具进行绘画', '{\"type\":\"object\",\"required\":[\"prompt\"],\"properties\":{\"prompt\":{\"type\":\"string\",\"description\":\"绘画提示词\"}}}', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVkIjowLCJ1c2VyX2lkIjowfQ.tLAGkF8XWh_G-oQzevpIodsswtPByBLoAZDz_eWuBgw', 'http://localhost:5678/api/function/dalle3', 0);

-- --------------------------------------------------------

--
-- 表的结构 `chatgpt_invite_codes`
--

DROP TABLE IF EXISTS `chatgpt_invite_codes`;
CREATE TABLE `chatgpt_invite_codes` (
  `id` int NOT NULL,
  `user_id` int NOT NULL COMMENT '用户ID',
  `code` char(8) NOT NULL COMMENT '邀请码',
  `hits` int NOT NULL COMMENT '点击次数',
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
  `inviter_id` int NOT NULL COMMENT '邀请人ID',
  `user_id` int NOT NULL COMMENT '注册用户ID',
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
(1, '对话聊天', '/images/menu/chat.png', '/chat', 0, 1),
(5, 'MJ 绘画', '/images/menu/mj.png', '/mj', 1, 1),
(6, 'SD 绘画', '/images/menu/sd.png', '/sd', 2, 1),
(7, '算力日志', '/images/menu/log.png', '/powerLog', 5, 1),
(8, '应用中心', '/images/menu/app.png', '/apps', 3, 1),
(9, '作品展示', '/images/menu/img-wall.png', '/images-wall', 4, 1),
(10, '会员计划', '/images/menu/member.png', '/member', 6, 1),
(11, '分享计划', '/images/menu/share.png', '/invite', 7, 1);

-- --------------------------------------------------------

--
-- 表的结构 `chatgpt_mj_jobs`
--

DROP TABLE IF EXISTS `chatgpt_mj_jobs`;
CREATE TABLE `chatgpt_mj_jobs` (
  `id` int NOT NULL,
  `user_id` int NOT NULL COMMENT '用户 ID',
  `task_id` varchar(20) DEFAULT NULL COMMENT '任务 ID',
  `type` varchar(20) DEFAULT 'image' COMMENT '任务类别',
  `message_id` char(40) NOT NULL COMMENT '消息 ID',
  `channel_id` char(40) DEFAULT NULL COMMENT '频道ID',
  `reference_id` char(40) DEFAULT NULL COMMENT '引用消息 ID',
  `prompt` varchar(2000) NOT NULL COMMENT '会话提示词',
  `img_url` varchar(400) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '图片URL',
  `org_url` varchar(400) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '原始图片地址',
  `hash` varchar(100) DEFAULT NULL COMMENT 'message hash',
  `progress` smallint DEFAULT '0' COMMENT '任务进度',
  `use_proxy` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否使用反代',
  `publish` tinyint(1) NOT NULL COMMENT '是否发布',
  `err_msg` varchar(255) DEFAULT NULL COMMENT '错误信息',
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
  `user_id` int NOT NULL COMMENT '用户ID',
  `product_id` int NOT NULL COMMENT '产品ID',
  `username` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '用户明',
  `order_no` varchar(30) NOT NULL COMMENT '订单ID',
  `trade_no` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '支付平台交易流水号',
  `subject` varchar(100) NOT NULL COMMENT '订单产品',
  `amount` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '订单金额',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '订单状态（0：待支付，1：已扫码，2：支付失败）',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '备注',
  `pay_time` int DEFAULT NULL COMMENT '支付时间',
  `pay_way` varchar(20) NOT NULL COMMENT '支付方式',
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
  `user_id` int NOT NULL COMMENT '用户ID',
  `username` varchar(30) NOT NULL COMMENT '用户名',
  `type` tinyint(1) NOT NULL COMMENT '类型（1：充值，2：消费，3：退费）',
  `amount` smallint NOT NULL COMMENT '算力数值',
  `balance` int NOT NULL COMMENT '余额',
  `model` varchar(30) NOT NULL COMMENT '模型',
  `remark` varchar(255) NOT NULL COMMENT '备注',
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
  `power` int NOT NULL DEFAULT '0' COMMENT '增加算力值',
  `enabled` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否启动',
  `sales` int NOT NULL DEFAULT '0' COMMENT '销量',
  `sort_num` tinyint NOT NULL DEFAULT '0' COMMENT '排序',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `app_url` varchar(255) DEFAULT NULL COMMENT 'App跳转地址',
  `url` varchar(255) DEFAULT NULL COMMENT '跳转地址'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='会员套餐表';

-- --------------------------------------------------------

--
-- 表的结构 `chatgpt_rewards`
--

DROP TABLE IF EXISTS `chatgpt_rewards`;
CREATE TABLE `chatgpt_rewards` (
  `id` int NOT NULL,
  `user_id` int NOT NULL COMMENT '用户 ID',
  `tx_id` char(36) NOT NULL COMMENT '交易 ID',
  `amount` decimal(10,2) NOT NULL COMMENT '打赏金额',
  `remark` varchar(80) NOT NULL COMMENT '备注',
  `status` tinyint(1) NOT NULL COMMENT '核销状态，0：未核销，1：已核销',
  `exchange` varchar(255) NOT NULL COMMENT '兑换详情（json）',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户打赏';

-- --------------------------------------------------------

--
-- 表的结构 `chatgpt_sd_jobs`
--

DROP TABLE IF EXISTS `chatgpt_sd_jobs`;
CREATE TABLE `chatgpt_sd_jobs` (
  `id` int NOT NULL,
  `user_id` int NOT NULL COMMENT '用户 ID',
  `type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT 'txt2img' COMMENT '任务类别',
  `task_id` char(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '任务 ID',
  `prompt` varchar(2000) NOT NULL COMMENT '会话提示词',
  `img_url` varchar(255) DEFAULT NULL COMMENT '图片URL',
  `params` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '绘画参数json',
  `progress` smallint DEFAULT '0' COMMENT '任务进度',
  `publish` tinyint(1) NOT NULL COMMENT '是否发布',
  `err_msg` varchar(255) DEFAULT NULL COMMENT '错误信息',
  `power` smallint NOT NULL DEFAULT '0' COMMENT '消耗算力',
  `created_at` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Stable Diffusion 任务表';

-- --------------------------------------------------------

--
-- 表的结构 `chatgpt_users`
--

DROP TABLE IF EXISTS `chatgpt_users`;
CREATE TABLE `chatgpt_users` (
  `id` int NOT NULL,
  `username` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '用户名',
  `nickname` varchar(30) NOT NULL COMMENT '昵称',
  `password` char(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '密码',
  `avatar` varchar(100) NOT NULL COMMENT '头像',
  `salt` char(12) NOT NULL COMMENT '密码盐',
  `power` int NOT NULL DEFAULT '0' COMMENT '剩余算力',
  `expired_time` int NOT NULL COMMENT '用户过期时间',
  `status` tinyint(1) NOT NULL COMMENT '当前状态',
  `chat_config_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '聊天配置json',
  `chat_roles_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '聊天角色 json',
  `chat_models_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT 'AI模型 json',
  `last_login_at` int NOT NULL COMMENT '最后登录时间',
  `vip` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否会员',
  `last_login_ip` char(16) NOT NULL COMMENT '最后登录 IP',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户表';

-- --------------------------------------------------------

--
-- 转存表中的数据 `chatgpt_users`
--

INSERT INTO `chatgpt_users` (`id`, `username`, `nickname`, `password`, `avatar`, `salt`, `power`, `expired_time`, `status`, `chat_config_json`, `chat_roles_json`, `chat_models_json`, `last_login_at`, `vip`, `last_login_ip`, `created_at`, `updated_at`) VALUES
    (4, '18575670125', '极客学长@830270', 'ccc3fb7ab61b8b5d096a4a166ae21d121fc38c71bbd1be6173d9ab973214a63b', 'http://localhost:5678/static/upload/2024/2/1708682650912429.png', 'ueedue5l', 9384, 1717292086, 1, '{\"api_keys\":{\"Azure\":\"\",\"ChatGLM\":\"\",\"OpenAI\":\"\"}}', '[\"red_book\",\"gpt\",\"programmer\",\"seller\"]', '[1,11]', 1711698298, 1, '::1', '2023-06-12 16:47:17', '2024-03-29 15:44:58');

-- --------------------------------------------------------

--
-- 表的结构 `chatgpt_user_login_logs`
--

DROP TABLE IF EXISTS `chatgpt_user_login_logs`;
CREATE TABLE `chatgpt_user_login_logs` (
  `id` int NOT NULL,
  `user_id` int NOT NULL COMMENT '用户ID',
  `username` varchar(30) NOT NULL COMMENT '用户名',
  `login_ip` char(16) NOT NULL COMMENT '登录IP',
  `login_address` varchar(30) NOT NULL COMMENT '登录地址',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户登录日志';

--
-- 转储表的索引
--

--
-- 表的索引 `chatgpt_admin_users`
--
ALTER TABLE `chatgpt_admin_users`
  ADD PRIMARY KEY (`id`) USING BTREE,
  ADD UNIQUE KEY `username` (`username`) USING BTREE;

--
-- 表的索引 `chatgpt_api_keys`
--
ALTER TABLE `chatgpt_api_keys`
  ADD PRIMARY KEY (`id`);

--
-- 表的索引 `chatgpt_chat_history`
--
ALTER TABLE `chatgpt_chat_history`
  ADD PRIMARY KEY (`id`),
  ADD KEY `chat_id` (`chat_id`);

--
-- 表的索引 `chatgpt_chat_items`
--
ALTER TABLE `chatgpt_chat_items`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `chat_id` (`chat_id`);

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
  ADD UNIQUE KEY `marker` (`marker`);

--
-- 表的索引 `chatgpt_configs`
--
ALTER TABLE `chatgpt_configs`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `marker` (`marker`);

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
  ADD UNIQUE KEY `name` (`name`);

--
-- 表的索引 `chatgpt_invite_codes`
--
ALTER TABLE `chatgpt_invite_codes`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `code` (`code`);

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
  ADD KEY `message_id` (`message_id`);

--
-- 表的索引 `chatgpt_orders`
--
ALTER TABLE `chatgpt_orders`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `order_no` (`order_no`);

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
-- 表的索引 `chatgpt_rewards`
--
ALTER TABLE `chatgpt_rewards`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `tx_id` (`tx_id`);

--
-- 表的索引 `chatgpt_sd_jobs`
--
ALTER TABLE `chatgpt_sd_jobs`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `task_id` (`task_id`);

--
-- 表的索引 `chatgpt_users`
--
ALTER TABLE `chatgpt_users`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `username` (`username`),
  ADD UNIQUE KEY `username_2` (`username`);

--
-- 表的索引 `chatgpt_user_login_logs`
--
ALTER TABLE `chatgpt_user_login_logs`
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
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=25;

--
-- 使用表AUTO_INCREMENT `chatgpt_chat_roles`
--
ALTER TABLE `chatgpt_chat_roles`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=130;

--
-- 使用表AUTO_INCREMENT `chatgpt_configs`
--
ALTER TABLE `chatgpt_configs`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=4;

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
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=12;

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
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `chatgpt_rewards`
--
ALTER TABLE `chatgpt_rewards`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `chatgpt_sd_jobs`
--
ALTER TABLE `chatgpt_sd_jobs`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `chatgpt_users`
--
ALTER TABLE `chatgpt_users`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `chatgpt_user_login_logs`
--
ALTER TABLE `chatgpt_user_login_logs`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
