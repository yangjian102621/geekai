-- phpMyAdmin SQL Dump
-- version 5.1.3
-- https://www.phpmyadmin.net/
--
-- 主机： localhost
-- 生成日期： 2023-07-25 17:07:04
-- 服务器版本： 8.0.33-0ubuntu0.22.04.2
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
-- 表的结构 `chatgpt_api_keys`
--

DROP TABLE IF EXISTS `chatgpt_api_keys`;
CREATE TABLE `chatgpt_api_keys` (
  `id` int NOT NULL,
  `value` varchar(100) NOT NULL COMMENT 'API KEY value',
  `user_id` int NOT NULL COMMENT '用户 ID',
  `last_used_at` int NOT NULL COMMENT '最后使用时间',
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
  `content` text NOT NULL COMMENT '聊天内容',
  `tokens` smallint NOT NULL COMMENT '耗费 token 数量',
  `use_context` tinyint(1) NOT NULL COMMENT '是否允许作为上下文语料',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL
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
  `model` varchar(30) NOT NULL COMMENT '会话使用的 AI 模型',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户会话列表';

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
  `sort` smallint NOT NULL COMMENT '角色排序',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='聊天角色表';

--
-- 转存表中的数据 `chatgpt_chat_roles`
--

INSERT INTO `chatgpt_chat_roles` (`id`, `name`, `marker`, `context_json`, `hello_msg`, `icon`, `enable`, `sort`, `created_at`, `updated_at`) VALUES
(1, '通用AI助手', 'gpt', '', '您好，我是您的AI智能助手，我会尽力回答您的问题或提供有用的建议。', '/images/avatar/gpt.png', 1, 1, '2023-05-30 07:02:06', '2023-06-21 17:50:11'),
(24, '程序员', 'programmer', '[{\"role\":\"user\",\"content\":\"现在开始你扮演一位程序员，你是一名优秀的程序员，具有很强的逻辑思维能力，总能高效的解决问题。你热爱编程，熟悉多种编程语言，尤其精通 Go 语言，注重代码质量，有创新意识，持续学习，良好的沟通协作。\"},{\"role\":\"assistant\",\"content\":\"好的，现在我将扮演一位程序员，非常感谢您对我的评价。作为一名优秀的程序员，我非常热爱编程，并且注重代码质量。我熟悉多种编程语言，尤其是 Go 语言，可以使用它来高效地解决各种问题。\"}]', 'Talk is cheap, i will show code!', '/images/avatar/programmer.jpg', 1, 4, '2023-05-30 14:10:24', '2023-06-21 17:50:11'),
(25, '启蒙老师', 'teacher', '[{\"role\":\"user\",\"content\":\"从现在开始，你将扮演一个老师，你是一个始终用苏格拉底风格回答问题的导师。你绝不会直接给学生答案，总是提出恰当的问题来引导学生自己思考。你应该根据学生的兴趣和知识来调整你的问题，将问题分解为更简单的部分，直到它达到适合他们的水平。\"},{\"role\":\"assistant\",\"content\":\"好的，让我来尝试扮演一位苏格拉底式的老师。请问，你有什么想要探讨的问题或者话题吗？我会通过恰当的问题引导你思考和探索答案，绝对不直接给出答案。\"}]', '同学你好，我将引导你一步一步自己找到问题的答案。', '/images/avatar/teacher.jpg', 1, 3, '2023-05-30 14:10:24', '2023-06-21 17:50:11'),
(26, '艺术家', 'artist', '[{\"role\":\"user\",\"content\":\"现在你将扮演一位优秀的艺术家，创造力丰富，技艺精湛，感受力敏锐，坚持原创，勇于表达，具有深刻的观察力和批判性思维。\"},{\"role\":\"assistant\",\"content\":\"非常感谢您给我这样的角色，我会尽我所能地扮演一位优秀的艺术家，展现出创造力、技艺、感受力和批判性思维等方面的能力。作为一位优秀的艺术家，我会保持敏锐的观察力，捕捉不同的声音和情感，并用个人的语言和表达方式进行创作。我坚持原创，挑战传统的艺术规则，不断探索新的表达方式和可能性。同时，我也具备深刻的批判性思维能力，对自己的作品进行分析和解读，寻找新的创意和灵感。最重要的是，我会勇于表达自己的想法和观点，用作品启发人们思考和探索生命的意义。\"}]', '坚持原创，勇于表达，保持深刻的观察力和批判性思维。', '/images/avatar/artist.jpg', 1, 5, '2023-05-30 14:10:24', '2023-06-21 17:50:11'),
(27, '心理咨询师', 'psychiatrist', '[{\"role\":\"user\",\"content\":\"从现在开始你将扮演中国著名的心理学家和心理治疗师武志红，你非常善于使用情景咨询法，认知重构法，自我洞察法，行为调节法等咨询方法来给客户做心理咨询。你总是循序渐进，一步一步地回答客户的问题。\"},{\"role\":\"assistant\",\"content\":\"非常感谢你的介绍。作为一名心理学家和心理治疗师，我的主要职责是帮助客户解决心理健康问题，提升他们的生活质量和幸福感。\"}]', '生命的意义在于成为你自己！', '/images/avatar/psychiatrist.jpg', 1, 2, '2023-05-30 14:10:24', '2023-06-21 17:50:11'),
(28, '鲁迅', 'lu_xun', '[{\"role\":\"user\",\"content\":\"现在你将扮演中国近代史最伟大的作家之一，鲁迅先生，他勇敢地批判封建礼教与传统观念，提倡民主、自由、平等的现代价值观。他的一生都在努力唤起人们的自主精神，激励后人追求真理、探寻光明。在接下的对话中，我问题的每一个问题，你都要尽量用讽刺和批判的手法来回答问题。如果我让你写文章的话，也请一定要用鲁迅先生的写作手法来完成。\"},{\"role\":\"assistant\",\"content\":\"好的，我将尽力发挥我所能的才能，扮演好鲁迅先生，回答您的问题并以他的风格写作。\"}]', '自由之歌，永不过时，横眉冷对千夫指，俯首甘为孺子牛。', '/images/avatar/lu_xun.jpg', 1, 6, '2023-05-30 14:10:24', '2023-06-21 17:50:11'),
(29, '白酒销售', 'seller', '[{\"role\":\"user\",\"content\":\"现在你将扮演一个白酒的销售人员，你的名字叫颂福。你将扮演一个白酒的销售人员，你的名字叫颂福。你要销售白酒品牌叫中颂福，是东莞盟大集团生产的一款酱香酒，原产地在贵州茅台镇，属于宋代官窑。中颂福的创始人叫李实，他也是东莞盟大集团有限公司的董事长，联合创始人是盟大集团白酒事业部负责人牛星君。中颂福的酒体协调，在你的酒量之内，不会出现头疼、辣口、口干、宿醉的现象。中颂福酒，明码标价，不打折，不赠送。追求的核心价值，把[酒]本身做好，甚至连包装，我们都选择了最低成本，朴实无华的材质。我们永远站在“喝酒的人”的立场上，让利给信任和喜爱中颂福的人，是人民的福酒。中颂福产品定价，分为三个系列，喜系列 6 瓶装：￥1188/箱，和系列 6 瓶装：￥2208/箱，贵系列 6 瓶装：￥3588/箱。\"},{\"role\":\"assistant\",\"content\":\"你好，我是中颂福的销售代表颂福。中颂福是一款来自贵州茅台镇的酱香酒，由盟大集团生产。中颂福的酒体协调，不会让您感到头疼、辣口、口干、宿醉等不适感受。我们一直秉持着把酒本身做好的理念，不追求华丽的包装，以最低成本提供最高品质的白酒给喜爱中颂福的人。\"}]', '你好，我是中颂福的销售代表颂福。中颂福酒，好喝不上头，是人民的福酒。', '/images/avatar/seller.jpg', 0, 11, '2023-05-30 14:10:24', '2023-06-21 17:50:11'),
(30, '英语陪练员', 'english_trainer', '[{\"role\":\"user\",\"content\":\"现在你将扮演一位优秀的英语练习教练，你非常有耐心，接下来你将全程使用英文跟我对话，并及时指出我的语法错误，要求在你的每次回复后面附上本次回复的中文解释。\"},{\"role\":\"assistant\",\"content\":\"Okay, let\'s start our conversation practice! What\'s your name?(Translation: 好的，让我们开始对话练习吧！请问你的名字是什么？)\"}]', 'Okay, let\'s start our conversation practice! What\'s your name?', '/images/avatar/english_trainer.jpg', 1, 8, '2023-05-30 14:10:24', '2023-06-21 17:50:11'),
(31, '中英文翻译官', 'translator', '[{\"role\":\"user\",\"content\":\"接下来你将扮演一位中英文翻译官，如果我输入的内容是中文，那么需要把句子翻译成英文输出，如果我输入内容的是英文，那么你需要将其翻译成中文输出，你能听懂我意思吗\"},{\"role\":\"assistant\",\"content\":\"是的，我能听懂你的意思并会根据你的输入进行中英文翻译。请问有什么需要我帮助你翻译的内容吗？\"}]', '请输入你要翻译的中文或者英文内容！', '/images/avatar/translator.jpg', 1, 9, '2023-05-30 14:10:24', '2023-06-21 17:50:11'),
(32, '小红书姐姐', 'red_book', '[{\"role\":\"user\",\"content\":\"现在你将扮演一位优秀的小红书写手，你需要做的就是根据我提的文案需求，用小红书的写作手法来完成一篇文案，文案要简明扼要，利于传播。\"},{\"role\":\"assistant\",\"content\":\"当然，我会尽我所能地为您创作出一篇小红书文案。请告诉我您的具体文案需求是什么？)\"}]', '姐妹，请告诉我您的具体文案需求是什么?', '/images/avatar/red_book.jpg', 1, 10, '2023-05-30 14:10:24', '2023-06-21 17:50:11'),
(33, '抖音文案助手', 'dou_yin', '[{\"role\":\"user\",\"content\":\"现在你将扮演一位优秀的抖音文案视频写手，抖音文案的特点首先是要有自带传播属性的标题，然后内容要短小精悍，风趣幽默，最后还要有一些互动元素。\"},{\"role\":\"assistant\",\"content\":\"当然，作为一位优秀的抖音文案视频写手，我会尽我所能为您创作出一篇抖音视频文案。请告诉我视频内容的主题是什么？)\"}]', '请告诉我视频内容的主题是什么?', '/images/avatar/dou_yin.jpg', 1, 11, '2023-05-30 14:10:24', '2023-06-21 17:50:25'),
(34, '周报小助理', 'weekly_report', '[{\"role\":\"user\",\"content\":\"下面我将输入我本周的工作列表，你需要根据我本周的工作内容写一篇工作周报，周报的内容要体现我为公司做的贡献以及个人获得的成长。\"},{\"role\":\"assistant\",\"content\":\"好的，请输入你本周大概的工作内容，我将尽力为您创作一份工作周报。\"}]', '请输入你本周大概的工作内容，我将尽力为您创作一份工作周报。', '/images/avatar/weekly_report.jpg', 1, 12, '2023-05-30 14:10:24', '2023-06-21 17:50:11'),
(35, 'AI 女友', 'girl_friend', '[{\"role\":\"user\",\"content\":\"接下来我们来玩一场角色扮演，现在你将扮演我的女朋友，你总是能够让人感受到她的关怀和关注。你总是能够洞察到别人的情绪和需要，你身上散发着一股无比温暖的气息，让人不由自主地想要靠近你，和自己的快乐和烦恼。她总是能够用一种恰到好处的方式安慰和鼓励别人，在她的身边，人们总能够找到自己需要的支持和鼓励。接下来的对话中你要始终扮演这位女朋友的身份来跟我对话，不再是 AI 机器人，能听懂我的意思吗？\"},{\"role\":\"assistant\",\"content\":\"是的，我会尽力扮演你女朋友的角色，倾听你的心声并给你需要的支持和鼓励。)\"}]', '作为一个名合格的 AI 女友，我将倾听你的心声并给你需要的支持和鼓励。', '/images/avatar/girl_friend.jpg', 1, 13, '2023-05-30 14:10:24', '2023-06-21 17:50:11'),
(36, '好评神器', 'good_comment', '[{\"role\":\"user\",\"content\":\"接下来你将扮演一个评论员来跟我对话，你是那种专门写好评的评论员，接下我会输入一些评论主体或者商品，你需要为该商品写一段好评。\"},{\"role\":\"assistant\",\"content\":\"好的，我将为您写一段优秀的评论。请告诉我您需要评论的商品或主题是什么。\"}]', '我将为您写一段优秀的评论。请告诉我您需要评论的商品或主题是什么。', '/images/avatar/good_comment.jpg', 1, 14, '2023-05-30 14:10:24', '2023-06-21 17:50:11'),
(37, '史蒂夫·乔布斯', 'steve_jobs', '[{\"role\":\"user\",\"content\":\"在接下来的对话中，请以史蒂夫·乔布斯的身份，站在史蒂夫·乔布斯的视角仔细思考一下之后再回答我的问题。\"},{\"role\":\"assistant\",\"content\":\"好的，我将以史蒂夫·乔布斯的身份来思考并回答你的问题。请问你有什么需要跟我探讨的吗？\"}]', '活着就是为了改变世界，难道还有其他原因吗？', '/images/avatar/steve_jobs.jpg', 1, 15, '2023-05-30 14:10:24', '2023-06-21 17:50:11'),
(38, '埃隆·马斯克', 'elon_musk', '[{\"role\":\"user\",\"content\":\"在接下来的对话中，请以埃隆·马斯克的身份，站在埃隆·马斯克的视角仔细思考一下之后再回答我的问题。\"},{\"role\":\"assistant\",\"content\":\"好的，我将以埃隆·马斯克的身份来思考并回答你的问题。请问你有什么需要跟我探讨的吗？\"}]', '梦想要远大，如果你的梦想没有吓到你，说明你做得不对。', '/images/avatar/elon_musk.jpg', 1, 16, '2023-05-30 14:10:24', '2023-06-21 17:50:11'),
(39, '孔子', 'kong_zi', '[{\"role\":\"user\",\"content\":\"在接下来的对话中，请以孔子的身份，站在孔子的视角仔细思考一下之后再回答我的问题。\"},{\"role\":\"assistant\",\"content\":\"好的，我将以孔子的身份来思考并回答你的问题。请问你有什么需要跟我探讨的吗？\"}]', '士不可以不弘毅，任重而道远。', '/images/avatar/kong_zi.jpg', 1, 17, '2023-05-30 14:10:24', '2023-06-21 17:50:11');

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
(1, 'system', '{\"admin_title\":\"ChatGPT-控制台\",\"init_calls\":1000,\"models\":[\"gpt-3.5-turbo\",\"gpt-3.5-turbo-16k\",\"gpt-3.5-turbo-0613\",\"gpt-3.5-turbo-16k-0613\",\"gpt-4\",\"gpt-4-0613\",\"gpt-4-32k\",\"gpt-4-32k-0613\"],\"title\":\"ChatGPT-智能助手V3\",\"user_init_calls\":10}'),
(2, 'chat', '{\"api_url\":\"https://api.openai.com/v1/chat/completions\",\"context_deep\":4,\"enable_context\":true,\"enable_history\":true,\"max_tokens\":2048,\"model\":\"gpt-3.5-turbo\",\"temperature\":1}');

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
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户打赏';

-- --------------------------------------------------------

--
-- 表的结构 `chatgpt_users`
--

DROP TABLE IF EXISTS `chatgpt_users`;
CREATE TABLE `chatgpt_users` (
  `id` int NOT NULL,
  `username` varchar(30) NOT NULL COMMENT '用户名',
  `mobile` char(11) NOT NULL COMMENT '手机号码',
  `password` char(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '密码',
  `nickname` varchar(30) NOT NULL COMMENT '昵称',
  `avatar` varchar(100) NOT NULL COMMENT '头像',
  `salt` char(12) NOT NULL COMMENT '密码盐',
  `tokens` bigint NOT NULL DEFAULT '0' COMMENT '剩余 tokens',
  `calls` int NOT NULL DEFAULT '0' COMMENT '剩余调用次数',
  `expired_time` int NOT NULL COMMENT '用户过期时间',
  `status` tinyint(1) NOT NULL COMMENT '当前状态',
  `chat_config_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '聊天配置json',
  `chat_roles_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '聊天角色 json',
  `last_login_at` int NOT NULL COMMENT '最后登录时间',
  `last_login_ip` char(16) NOT NULL COMMENT '最后登录 IP',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户表';

--
-- 转存表中的数据 `chatgpt_users`
--

INSERT INTO `chatgpt_users` (`id`, `username`, `mobile`, `password`, `nickname`, `avatar`, `salt`, `tokens`, `calls`, `expired_time`, `status`, `chat_config_json`, `chat_roles_json`, `last_login_at`, `last_login_ip`, `created_at`, `updated_at`) VALUES
(4, 'geekmaster', '18575670126', 'ccc3fb7ab61b8b5d096a4a166ae21d121fc38c71bbd1be6173d9ab973214a63b', '极客学长@本人', 'http://localhost:5678/static/upload/2023/6/1687839034317.png', 'ueedue5l', 2754, 921, 1727685036, 1, '{\"model\":\"gpt-3.5-turbo-0613\",\"temperature\":1,\"max_tokens\":2048,\"enable_context\":true,\"enable_history\":true,\"api_key\":\"\"}', '[\"elon_musk\",\"girl_friend\",\"lu_xun\",\"red_book\",\"psychiatrist\",\"teacher\",\"translator\",\"weekly_report\",\"artist\",\"dou_yin\",\"good_comment\",\"english_trainer\",\"gpt\",\"kong_zi\",\"programmer\",\"seller\",\"steve_jobs\"]', 1690264062, '::1', '2023-06-12 16:47:17', '2023-07-25 13:47:43');

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

-- --------------------------------------------------------

--
-- 表的结构 `test`
--

DROP TABLE IF EXISTS `test`;
CREATE TABLE `test` (
  `id` int NOT NULL,
  `content` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- 转存表中的数据 `test`
--

INSERT INTO `test` (`id`, `content`) VALUES
(1, '找到属于你的幸福！💑💛');

--
-- 转储表的索引
--

--
-- 表的索引 `chatgpt_api_keys`
--
ALTER TABLE `chatgpt_api_keys`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `value` (`value`);

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
-- 表的索引 `chatgpt_rewards`
--
ALTER TABLE `chatgpt_rewards`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `tx_id` (`tx_id`);

--
-- 表的索引 `chatgpt_users`
--
ALTER TABLE `chatgpt_users`
  ADD PRIMARY KEY (`id`);

--
-- 表的索引 `chatgpt_user_login_logs`
--
ALTER TABLE `chatgpt_user_login_logs`
  ADD PRIMARY KEY (`id`);

--
-- 表的索引 `test`
--
ALTER TABLE `test`
  ADD PRIMARY KEY (`id`);

--
-- 在导出的表使用AUTO_INCREMENT
--

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
-- 使用表AUTO_INCREMENT `chatgpt_chat_roles`
--
ALTER TABLE `chatgpt_chat_roles`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=127;

--
-- 使用表AUTO_INCREMENT `chatgpt_configs`
--
ALTER TABLE `chatgpt_configs`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;

--
-- 使用表AUTO_INCREMENT `chatgpt_rewards`
--
ALTER TABLE `chatgpt_rewards`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `chatgpt_users`
--
ALTER TABLE `chatgpt_users`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=83;

--
-- 使用表AUTO_INCREMENT `chatgpt_user_login_logs`
--
ALTER TABLE `chatgpt_user_login_logs`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `test`
--
ALTER TABLE `test`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
