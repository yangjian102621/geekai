-- phpMyAdmin SQL Dump
-- version 5.2.1
-- https://www.phpmyadmin.net/
--
-- 主机： localhost
-- 生成日期： 2023-10-08 12:43:26
-- 服务器版本： 8.0.27
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
  `platform` char(20) DEFAULT NULL COMMENT '平台',
  `value` varchar(100) NOT NULL COMMENT 'API KEY value',
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
  `use_context` tinyint(1) DEFAULT NULL COMMENT '是否允许作为上下文语料',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `deleted_at` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='聊天历史记录';

--
-- 转存表中的数据 `chatgpt_chat_history`
--

INSERT INTO `chatgpt_chat_history` (`id`, `user_id`, `chat_id`, `type`, `icon`, `role_id`, `content`, `tokens`, `use_context`, `created_at`, `updated_at`, `deleted_at`) VALUES
(258, 4, 'f327e7d3-7563-46a6-9e24-2095d8c901dc', 'prompt', 'http://nk.img.r9it.com/chatgpt-plus/1694821157223284.png', 1, '金刚经里的核心思想是什么？', 16, 1, '2023-09-16 07:33:17', '2023-09-16 07:33:17', NULL),
(259, 4, 'f327e7d3-7563-46a6-9e24-2095d8c901dc', 'reply', '/images/avatar/gpt.png', 1, '金刚经是佛教经典之一，其中的核心思想是\"一切有为法，如梦幻泡影，如露亦如电，应作如是观\"。这句话的含义是指一切事物都是无常的，没有固定的实体，就像梦幻、泡影、露水一样，都是瞬息即逝的存在。因此，我们应该以无常的观念来看待世界，不执着于事物的表象和存在，从而超越痛苦和迷惑，获得智慧和解脱。金刚经强调了无我、无我见的观念，教导人们要超越个体的概念，以开放、空明的心态去体验和理解世界。', 251, 1, '2023-09-16 07:33:20', '2023-09-16 07:33:20', NULL),
(260, 4, 'f327e7d3-7563-46a6-9e24-2095d8c901dc', 'prompt', 'http://nk.img.r9it.com/chatgpt-plus/1694821157223284.png', 1, '凡所有相皆是虚妄，若见诸相非相，即见如来。', 27, 1, '2023-09-16 07:34:16', '2023-09-16 07:34:16', NULL),
(261, 4, 'f327e7d3-7563-46a6-9e24-2095d8c901dc', 'reply', '/images/avatar/gpt.png', 1, '这句话来自于金刚经中的经文，它的意思是，所有的事物都是虚妄的、空无实体的，并且如果我们能够明白万物无常、无我、空性的真相，就能够认识到如来的存在。\n\n\"凡所有相\"指的是一切外在的形象、现象、事物，它们都是虚妄的，没有固定的实体。如果我们能够超越相象的表面，不执着于它们的存在和特征，就能够真正认识到如来的存在。\n\n\"非相\"则指的是超越相象的本质，即超越一切形象和概念的本体。当我们能够从相的层面超越，认识到一切皆空、无我、无常的真相时，我们就能够觉悟到如来的存在。\n\n这句话表达了金刚经中对于觉悟和智慧的重要意义，只有通过超越相象的观照和领悟，我们才能认识到如来的真实境界。它也强调了金刚经核心思想中\"如梦幻泡影，如露亦如电\"的观念，即一切事物都是瞬息万变、虚幻不实的。通过这种认识，我们可以超越纷扰的世界，获得智慧和解脱。', 700, 1, '2023-09-16 07:34:18', '2023-09-16 07:34:18', NULL),
(262, 4, 'f327e7d3-7563-46a6-9e24-2095d8c901dc', 'prompt', 'http://nk.img.r9it.com/chatgpt-plus/1694821157223284.png', 1, '“有为法”具体是指什么？', 13, 1, '2023-09-16 07:35:21', '2023-09-16 07:35:21', NULL),
(263, 4, 'f327e7d3-7563-46a6-9e24-2095d8c901dc', 'reply', '/images/avatar/gpt.png', 1, '\"有为法\"是金刚经中的名词，它指的是一切有形、有相、有为的存在，包括物质、各种事物、现象和生命体等具有实际存在的实体。\n\n\"有为\"可理解为具有有形有相的特征，代表着一切有名有相的现象和实体。这些有为之法，即一切有形之物，都是由因缘和条件所组成，具有生灭、变化和相互关系的特性。它们是相对的、条件性的存在，同时也是无常、虚妄、无实性的。\n\n金刚经通过提到\"有为法\"，强调了一切事物都是相对和无常的，没有固定、独立、实在的本质。它们的存在和特性是依靠因缘的变化而产生的，所以不应该执着于它们的表象和存在，而是要超越它们的相对性，以无常无我的观点去对待世界。\n\n总的来说，\"有为法\"是指一切有形、有相、有名的现象和实体，它们是无常、虚妄、无实性的。通过识别和理解这一点，我们可以超越执着、觉悟无常的真相，获得智慧和解脱。', 1104, 1, '2023-09-16 07:35:23', '2023-09-16 07:35:23', NULL),
(264, 4, '49e5219e-1058-4320-9580-290f731afb4e', 'prompt', 'http://nk.img.r9it.com/chatgpt-plus/1694821157223284.png', 1, '绘画：菊花，梵高风格', 16, 0, '2023-09-16 07:54:36', '2023-09-16 07:54:36', NULL),
(265, 4, '49e5219e-1058-4320-9580-290f731afb4e', 'reply', '/images/avatar/gpt.png', 1, '绘画提示词：chrysanthemum, van gogh style --v 5.2 已推送任务到 MidJourney 机器人，请耐心等待任务执行...', 58, 0, '2023-09-16 07:54:38', '2023-09-16 07:54:38', NULL),
(266, 4, '49e5219e-1058-4320-9580-290f731afb4e', 'mj', '/images/avatar/mid_journey.png', 1, '{\"message_id\":\"1152392223497924658\",\"reference_id\":\"\",\"image\":{\"url\":\"http://nk.img.r9it.com/chatgpt-plus/1694822120588799.png\",\"proxy_url\":\"https://media.discordapp.net/attachments/1151037077308325901/1152392222927487016/jiangjin_chrysanthemum_van_gogh_style_54689110-4701-4eed-b308-1acaf1c5efa1.png\",\"filename\":\"jiangjin_chrysanthemum_van_gogh_style_54689110-4701-4eed-b308-1acaf1c5efa1.png\",\"width\":2048,\"height\":2048,\"size\":8515347,\"hash\":\"54689110-4701-4eed-b308-1acaf1c5efa1\"},\"content\":\"**chrysanthemum, van gogh style --v 5.2 --s 750 --style raw** - \\u003c@1075058490378289194\\u003e (fast)\",\"prompt\":\"chrysanthemum, van gogh style --v 5.2 --s 750 --style raw\",\"status\":\"Finished\",\"progress\":100}', 0, 0, '2023-09-16 07:55:25', '2023-09-16 07:55:25', NULL),
(267, 4, '', 'mj', '/images/avatar/mid_journey.png', 0, '{\"message_id\":\"1152392353139658894\",\"reference_id\":\"1152392223497924658\",\"image\":{\"url\":\"http://nk.img.r9it.com/chatgpt-plus/1694822151252833.png\",\"proxy_url\":\"https://media.discordapp.net/attachments/1151037077308325901/1152392352804118568/jiangjin_chrysanthemum_van_gogh_style_15f3d9e6-9bd2-4e31-81ce-c8b126509b12.png\",\"filename\":\"jiangjin_chrysanthemum_van_gogh_style_15f3d9e6-9bd2-4e31-81ce-c8b126509b12.png\",\"width\":1024,\"height\":1024,\"size\":1930002,\"hash\":\"15f3d9e6-9bd2-4e31-81ce-c8b126509b12\"},\"content\":\"**chrysanthemum, van gogh style --v 5.2 --s 750 --style raw** - Image #3 \\u003c@1075058490378289194\\u003e\",\"prompt\":\"chrysanthemum, van gogh style --v 5.2 --s 750 --style raw\",\"status\":\"Finished\",\"progress\":100}', 0, 0, '2023-09-16 07:55:55', '2023-09-16 07:55:55', NULL),
(268, 4, '49e5219e-1058-4320-9580-290f731afb4e', 'prompt', 'http://nk.img.r9it.com/chatgpt-plus/1694821157223284.png', 1, '绘画：一个中古女孩，穿着旗袍，走在中国90年代的上海城市街道 --ar 9:16', 42, 0, '2023-09-16 07:59:13', '2023-09-16 07:59:13', NULL),
(269, 4, '49e5219e-1058-4320-9580-290f731afb4e', 'reply', '/images/avatar/gpt.png', 1, '绘画提示词：A vintage girl wearing a cheongsam walking on the streets of 1990s Shanghai, China --v 5.2 已推送任务到 MidJourney 机器人，请耐心等待任务执行...', 121, 0, '2023-09-16 07:59:15', '2023-09-16 07:59:15', NULL),
(270, 4, '49e5219e-1058-4320-9580-290f731afb4e', 'mj', '/images/avatar/mid_journey.png', 1, '{\"message_id\":\"1152393378244349953\",\"reference_id\":\"\",\"image\":{\"url\":\"http://nk.img.r9it.com/chatgpt-plus/1694822395949588.png\",\"proxy_url\":\"https://media.discordapp.net/attachments/1151037077308325901/1152393377837490216/jiangjin_A_vintage_girl_wearing_a_cheongsam_walking_on_the_stre_319f953e-9ba4-4c8f-bd0c-87444a44b4a6.png\",\"filename\":\"jiangjin_A_vintage_girl_wearing_a_cheongsam_walking_on_the_stre_319f953e-9ba4-4c8f-bd0c-87444a44b4a6.png\",\"width\":2048,\"height\":2048,\"size\":7013301,\"hash\":\"319f953e-9ba4-4c8f-bd0c-87444a44b4a6\"},\"content\":\"**A vintage girl wearing a cheongsam walking on the streets of 1990s Shanghai, China --v 5.2 --s 750 --style raw** - \\u003c@1075058490378289194\\u003e (fast)\",\"prompt\":\"A vintage girl wearing a cheongsam walking on the streets of 1990s Shanghai, China --v 5.2 --s 750 --style raw\",\"status\":\"Finished\",\"progress\":100}', 0, 0, '2023-09-16 08:00:00', '2023-09-16 08:00:00', NULL),
(271, 4, '49e5219e-1058-4320-9580-290f731afb4e', 'prompt', 'http://nk.img.r9it.com/chatgpt-plus/1694821157223284.png', 1, '绘画：一个中古女孩，穿着旗袍，走在中国90年代的上海城市街道 --ar 9:16', 42, 0, '2023-09-16 08:04:15', '2023-09-16 08:04:15', NULL),
(272, 4, '49e5219e-1058-4320-9580-290f731afb4e', 'reply', '/images/avatar/gpt.png', 1, '绘画提示词：A medieval girl wearing a cheongsam walking on the streets of 1990s Shanghai, China --ar 9:16 --v 5.2 已推送任务到 MidJourney 机器人，请耐心等待任务执行...', 121, 0, '2023-09-16 08:04:17', '2023-09-16 08:04:17', NULL),
(273, 4, '49e5219e-1058-4320-9580-290f731afb4e', 'mj', '/images/avatar/mid_journey.png', 1, '{\"message_id\":\"1152394661802033162\",\"reference_id\":\"\",\"image\":{\"url\":\"http://nk.img.r9it.com/chatgpt-plus/1694822701966804.png\",\"proxy_url\":\"https://media.discordapp.net/attachments/1151037077308325901/1152394661281927229/jiangjin_A_medieval_girl_wearing_a_cheongsam_walking_on_the_str_7e522b23-2e88-40e2-814f-b9ad7ce1777b.png\",\"filename\":\"jiangjin_A_medieval_girl_wearing_a_cheongsam_walking_on_the_str_7e522b23-2e88-40e2-814f-b9ad7ce1777b.png\",\"width\":1632,\"height\":2912,\"size\":8084572,\"hash\":\"7e522b23-2e88-40e2-814f-b9ad7ce1777b\"},\"content\":\"**A medieval girl wearing a cheongsam walking on the streets of 1990s Shanghai, China --ar 9:16 --v 5.2 --s 750 --style raw** - \\u003c@1075058490378289194\\u003e (fast)\",\"prompt\":\"A medieval girl wearing a cheongsam walking on the streets of 1990s Shanghai, China --ar 9:16 --v 5.2 --s 750 --style raw\",\"status\":\"Finished\",\"progress\":100}', 0, 0, '2023-09-16 08:05:06', '2023-09-16 08:05:06', NULL),
(274, 4, '49e5219e-1058-4320-9580-290f731afb4e', 'prompt', 'http://nk.img.r9it.com/chatgpt-plus/1694821157223284.png', 1, '绘画：一个中国女孩，穿着旗袍，走在中国90年代的上海城市街道 --ar 2:3', 40, 0, '2023-09-16 08:05:09', '2023-09-16 08:05:09', NULL),
(275, 4, '49e5219e-1058-4320-9580-290f731afb4e', 'reply', '/images/avatar/gpt.png', 1, '绘画提示词：A Chinese girl wearing a qipao walking on the streets of Shanghai in the 1990s --ar 2:3 --v 5.2 已推送任务到 MidJourney 机器人，请耐心等待任务执行...', 119, 0, '2023-09-16 08:05:11', '2023-09-16 08:05:11', NULL),
(276, 4, '49e5219e-1058-4320-9580-290f731afb4e', 'mj', '/images/avatar/mid_journey.png', 1, '{\"message_id\":\"1152395158298558464\",\"reference_id\":\"\",\"image\":{\"url\":\"http://nk.img.r9it.com/chatgpt-plus/1694822820475390.png\",\"proxy_url\":\"https://media.discordapp.net/attachments/1151037077308325901/1152395157707169822/jiangjin_A_Chinese_girl_wearing_a_qipao_walking_on_the_streets__21965cd6-ef70-4b36-8d92-7b201dea77ff.png\",\"filename\":\"jiangjin_A_Chinese_girl_wearing_a_qipao_walking_on_the_streets__21965cd6-ef70-4b36-8d92-7b201dea77ff.png\",\"width\":1792,\"height\":2688,\"size\":8130309,\"hash\":\"21965cd6-ef70-4b36-8d92-7b201dea77ff\"},\"content\":\"**A Chinese girl wearing a qipao walking on the streets of Shanghai in the 1990s --ar 2:3 --v 5.2 --s 750 --style raw** - \\u003c@1075058490378289194\\u003e (fast)\",\"prompt\":\"A Chinese girl wearing a qipao walking on the streets of Shanghai in the 1990s --ar 2:3 --v 5.2 --s 750 --style raw\",\"status\":\"Finished\",\"progress\":100}', 0, 0, '2023-09-16 08:07:05', '2023-09-16 08:07:05', NULL),
(277, 4, '49e5219e-1058-4320-9580-290f731afb4e', 'prompt', 'http://nk.img.r9it.com/chatgpt-plus/1694821157223284.png', 1, '绘画：一个中国女孩，穿着青色旗袍，走在中国90年代的上海巷子 --ar 16:9', 41, 0, '2023-09-16 08:07:29', '2023-09-16 08:07:29', NULL),
(278, 4, '49e5219e-1058-4320-9580-290f731afb4e', 'reply', '/images/avatar/gpt.png', 1, '绘画提示词：A Chinese girl wearing a blue cheongsam walking in the alleys of 1990s Shanghai --ar 16:9 --v 5.2 已推送任务到 MidJourney 机器人，请耐心等待任务执行...', 122, 0, '2023-09-16 08:07:31', '2023-09-16 08:07:31', NULL),
(279, 4, '49e5219e-1058-4320-9580-290f731afb4e', 'mj', '/images/avatar/mid_journey.png', 1, '{\"message_id\":\"1152395485248753734\",\"reference_id\":\"\",\"image\":{\"url\":\"http://nk.img.r9it.com/chatgpt-plus/1694822897511213.png\",\"proxy_url\":\"https://media.discordapp.net/attachments/1151037077308325901/1152395484586049547/jiangjin_A_Chinese_girl_wearing_a_blue_cheongsam_walking_in_the_036aad69-f231-4b5c-b537-c45bd2027b05.png\",\"filename\":\"jiangjin_A_Chinese_girl_wearing_a_blue_cheongsam_walking_in_the_036aad69-f231-4b5c-b537-c45bd2027b05.png\",\"width\":2912,\"height\":1632,\"size\":7622815,\"hash\":\"036aad69-f231-4b5c-b537-c45bd2027b05\"},\"content\":\"**A Chinese girl wearing a blue cheongsam walking in the alleys of 1990s Shanghai --ar 16:9 --v 5.2 --s 750 --style raw** - \\u003c@1075058490378289194\\u003e (fast)\",\"prompt\":\"A Chinese girl wearing a blue cheongsam walking in the alleys of 1990s Shanghai --ar 16:9 --v 5.2 --s 750 --style raw\",\"status\":\"Finished\",\"progress\":100}', 0, 0, '2023-09-16 08:08:22', '2023-09-16 08:08:22', NULL),
(280, 4, '49e5219e-1058-4320-9580-290f731afb4e', 'prompt', 'http://nk.img.r9it.com/chatgpt-plus/1694821157223284.png', 1, '绘画：一个中国女孩，穿着旗袍，走在铺满青砖的街道 --ar 3:4', 41, 0, '2023-09-16 08:09:20', '2023-09-16 08:09:20', NULL),
(281, 4, '49e5219e-1058-4320-9580-290f731afb4e', 'reply', '/images/avatar/gpt.png', 1, '绘画提示词：A Chinese girl wearing a cheongsam walking on a street paved with blue bricks --v 5.2 已推送任务到 MidJourney 机器人，请耐心等待任务执行...', 114, 0, '2023-09-16 08:09:22', '2023-09-16 08:09:22', NULL),
(282, 4, '49e5219e-1058-4320-9580-290f731afb4e', 'mj', '/images/avatar/mid_journey.png', 1, '{\"message_id\":\"1152395938078392370\",\"reference_id\":\"\",\"image\":{\"url\":\"http://nk.img.r9it.com/chatgpt-plus/1694823005955584.png\",\"proxy_url\":\"https://media.discordapp.net/attachments/1151037077308325901/1152395936602001420/jiangjin_A_Chinese_girl_wearing_a_cheongsam_walking_on_a_street_d6336496-3115-47b7-bcf8-3bebda1d15ec.png\",\"filename\":\"jiangjin_A_Chinese_girl_wearing_a_cheongsam_walking_on_a_street_d6336496-3115-47b7-bcf8-3bebda1d15ec.png\",\"width\":2048,\"height\":2048,\"size\":6829490,\"hash\":\"d6336496-3115-47b7-bcf8-3bebda1d15ec\"},\"content\":\"**A Chinese girl wearing a cheongsam walking on a street paved with blue bricks --v 5.2 --s 750 --style raw** - \\u003c@1075058490378289194\\u003e (fast)\",\"prompt\":\"A Chinese girl wearing a cheongsam walking on a street paved with blue bricks --v 5.2 --s 750 --style raw\",\"status\":\"Finished\",\"progress\":100}', 0, 0, '2023-09-16 08:10:10', '2023-09-16 08:10:10', NULL),
(283, 4, '49e5219e-1058-4320-9580-290f731afb4e', 'prompt', 'http://nk.img.r9it.com/chatgpt-plus/1694821157223284.png', 1, '绘画：一只老虎撕咬着一只鳄鱼', 24, 0, '2023-09-17 12:49:21', '2023-09-17 12:49:21', NULL),
(284, 4, '49e5219e-1058-4320-9580-290f731afb4e', 'reply', '/images/avatar/gpt.png', 1, '绘画提示词：A tiger biting a crocodile --v 5.2 已推送任务到 MidJourney 机器人，请耐心等待任务执行...', 57, 0, '2023-09-17 12:49:24', '2023-09-17 12:49:24', NULL),
(285, 4, '49e5219e-1058-4320-9580-290f731afb4e', 'mj', '/images/avatar/mid_journey.png', 1, '{\"message_id\":\"1152828784144298066\",\"reference_id\":\"\",\"image\":{\"url\":\"http://nk.img.r9it.com/chatgpt-plus/1694926203790426.png\",\"proxy_url\":\"https://media.discordapp.net/attachments/1151037077308325901/1152828783380922388/jiangjin_A_tiger_biting_a_crocodile_423b4a4f-9607-4dd1-a45b-a1a50be667e7.png\",\"filename\":\"jiangjin_A_tiger_biting_a_crocodile_423b4a4f-9607-4dd1-a45b-a1a50be667e7.png\",\"width\":2048,\"height\":2048,\"size\":7723624,\"hash\":\"423b4a4f-9607-4dd1-a45b-a1a50be667e7\"},\"content\":\"**A tiger biting a crocodile --v 5.2 --s 750 --style raw** - \\u003c@1075058490378289194\\u003e (fast)\",\"prompt\":\"A tiger biting a crocodile --v 5.2 --s 750 --style raw\",\"status\":\"Finished\",\"progress\":100}', 0, 0, '2023-09-17 12:50:11', '2023-09-17 12:50:11', NULL),
(286, 4, '49e5219e-1058-4320-9580-290f731afb4e', 'prompt', 'http://nk.img.r9it.com/chatgpt-plus/1694821157223284.png', 1, '绘画：巨大的数字屏幕，屏幕上有0,1 二进制代码', 28, 0, '2023-09-17 15:00:29', '2023-09-17 15:00:29', NULL),
(287, 4, '49e5219e-1058-4320-9580-290f731afb4e', 'reply', '/images/avatar/gpt.png', 1, '绘画提示词：A huge digital screen with binary code of 0 and 1 --v 5.2 已推送任务到 MidJourney 机器人，请耐心等待任务执行...', 74, 0, '2023-09-17 15:00:33', '2023-09-17 15:00:33', NULL),
(288, 4, '49e5219e-1058-4320-9580-290f731afb4e', 'mj', '/images/avatar/mid_journey.png', 1, '{\"message_id\":\"1152861794356953098\",\"reference_id\":\"\",\"image\":{\"url\":\"http://nk.img.r9it.com/chatgpt-plus/1694934140188981.png\",\"proxy_url\":\"https://media.discordapp.net/attachments/1151037077308325901/1152861793706844232/jiangjin_A_huge_digital_screen_with_binary_code_of_0_and_1_7b6a6910-0aa7-4eda-8580-ff882118c428.png\",\"filename\":\"jiangjin_A_huge_digital_screen_with_binary_code_of_0_and_1_7b6a6910-0aa7-4eda-8580-ff882118c428.png\",\"width\":2048,\"height\":2048,\"size\":7800660,\"hash\":\"7b6a6910-0aa7-4eda-8580-ff882118c428\"},\"content\":\"**A huge digital screen with binary code of 0 and 1 --v 5.2 --s 750 --style raw** - \\u003c@1075058490378289194\\u003e (fast)\",\"prompt\":\"A huge digital screen with binary code of 0 and 1 --v 5.2 --s 750 --style raw\",\"status\":\"Finished\",\"progress\":100}', 0, 0, '2023-09-17 15:02:38', '2023-09-17 15:02:38', NULL),
(289, 4, '49e5219e-1058-4320-9580-290f731afb4e', 'prompt', 'http://nk.img.r9it.com/chatgpt-plus/1694821157223284.png', 1, '绘画：一只小熊骑着黄色的自行车', 22, 0, '2023-09-17 15:03:12', '2023-09-17 15:03:12', NULL),
(290, 4, '49e5219e-1058-4320-9580-290f731afb4e', 'reply', '/images/avatar/gpt.png', 1, '绘画提示词：A bear riding a yellow bicycle --v 5.2 已推送任务到 MidJourney 机器人，请耐心等待任务执行...', 56, 0, '2023-09-17 15:03:14', '2023-09-17 15:03:14', NULL),
(291, 4, '49e5219e-1058-4320-9580-290f731afb4e', 'mj', '/images/avatar/mid_journey.png', 1, '{\"message_id\":\"1152862460366307349\",\"reference_id\":\"\",\"image\":{\"url\":\"http://nk.img.r9it.com/chatgpt-plus/1694934237667586.png\",\"proxy_url\":\"https://media.discordapp.net/attachments/1151037077308325901/1152862460110458961/jiangjin_A_bear_riding_a_yellow_bicycle_0f95fd58-0385-438c-ac96-0ecaa57bfcd9.png\",\"filename\":\"jiangjin_A_bear_riding_a_yellow_bicycle_0f95fd58-0385-438c-ac96-0ecaa57bfcd9.png\",\"width\":2048,\"height\":2048,\"size\":7128466,\"hash\":\"0f95fd58-0385-438c-ac96-0ecaa57bfcd9\"},\"content\":\"**A bear riding a yellow bicycle --v 5.2 --s 750 --style raw** - \\u003c@1075058490378289194\\u003e (fast)\",\"prompt\":\"A bear riding a yellow bicycle --v 5.2 --s 750 --style raw\",\"status\":\"Finished\",\"progress\":100}', 0, 0, '2023-09-17 15:04:06', '2023-09-17 15:04:06', NULL),
(292, 4, '49e5219e-1058-4320-9580-290f731afb4e', 'prompt', 'http://nk.img.r9it.com/chatgpt-plus/1694821157223284.png', 1, '绘画：Beautiful girl, concept art, 8k intricate details, fairytale style', 20, 0, '2023-09-17 15:04:53', '2023-09-17 15:04:53', NULL),
(293, 4, '49e5219e-1058-4320-9580-290f731afb4e', 'reply', '/images/avatar/gpt.png', 1, '绘画提示词：Beautiful girl, concept art, intricate details --ar 16:9 --s fairytale --v 5.2 已推送任务到 MidJourney 机器人，请耐心等待任务执行...', 110, 0, '2023-09-17 15:04:55', '2023-09-17 15:04:55', NULL),
(294, 4, '49e5219e-1058-4320-9580-290f731afb4e', 'prompt', 'http://nk.img.r9it.com/chatgpt-plus/1694821157223284.png', 1, '请根据金庸小说描述，把《射雕英雄传》中的女主角黄蓉画出来。神仙姐姐，气质出尘脱俗，身穿白色的衣服，金色的发带，灵动有神的双眼。', 89, 0, '2023-09-17 15:14:01', '2023-09-17 15:14:01', NULL),
(295, 4, '49e5219e-1058-4320-9580-290f731afb4e', 'reply', '/images/avatar/gpt.png', 1, '绘画提示词：The protagonist of the novel is Huang Rong. She is beautiful and charming, with a fairy-like temperament. She is wearing a white dress with a golden hairband, and her eyes are lively and divine. 已推送任务到 MidJourney 机器人，请耐心等待任务执行...', 193, 0, '2023-09-17 15:14:04', '2023-09-17 15:14:04', NULL),
(296, 4, '49e5219e-1058-4320-9580-290f731afb4e', 'mj', '/images/avatar/mid_journey.png', 1, '{\"message_id\":\"1152865215201935502\",\"reference_id\":\"\",\"image\":{\"url\":\"http://nk.img.r9it.com/chatgpt-plus/1694934895270654.png\",\"proxy_url\":\"https://media.discordapp.net/attachments/1151037077308325901/1152865214556029039/jiangjin_The_protagonist_of_the_novel_is_Huang_Rong._She_is_bea_64531b1a-a335-498f-8c10-3eb70b5e2c3d.png\",\"filename\":\"jiangjin_The_protagonist_of_the_novel_is_Huang_Rong._She_is_bea_64531b1a-a335-498f-8c10-3eb70b5e2c3d.png\",\"width\":2048,\"height\":2048,\"size\":7037265,\"hash\":\"64531b1a-a335-498f-8c10-3eb70b5e2c3d\"},\"content\":\"**The protagonist of the novel is Huang Rong. She is beautiful and charming, with a fairy-like temperament. She is wearing a white dress with a golden hairband, and her eyes are lively and divine. --s 750 --style raw** - \\u003c@1075058490378289194\\u003e (fast)\",\"prompt\":\"The protagonist of the novel is Huang Rong. She is beautiful and charming, with a fairy-like temperament. She is wearing a white dress with a golden hairband, and her eyes are lively and divine. --s 750 --style raw\",\"status\":\"Finished\",\"progress\":100}', 0, 0, '2023-09-17 15:15:02', '2023-09-17 15:15:02', NULL),
(297, 4, '49e5219e-1058-4320-9580-290f731afb4e', 'prompt', 'http://nk.img.r9it.com/chatgpt-plus/1694821157223284.png', 1, '画一只猫，赛博朋克风格', 18, 0, '2023-09-17 15:21:25', '2023-09-17 15:21:25', NULL),
(298, 4, '49e5219e-1058-4320-9580-290f731afb4e', 'reply', '/images/avatar/gpt.png', 1, '绘画提示词：Draw a cyberpunk-style cat 已推送任务到 MidJourney 机器人，请耐心等待任务执行...', 52, 0, '2023-09-17 15:21:27', '2023-09-17 15:21:27', NULL),
(299, 4, '49e5219e-1058-4320-9580-290f731afb4e', 'mj', '/images/avatar/mid_journey.png', 1, '{\"message_id\":\"1152867045902397440\",\"reference_id\":\"\",\"image\":{\"url\":\"http://nk.img.r9it.com/chatgpt-plus/1694935330386877.png\",\"proxy_url\":\"https://media.discordapp.net/attachments/1151037077308325901/1152867045080301639/jiangjin_Draw_a_cyberpunk-style_cat_382ff1ce-57c2-45b0-b7fd-9820552171c3.png\",\"filename\":\"jiangjin_Draw_a_cyberpunk-style_cat_382ff1ce-57c2-45b0-b7fd-9820552171c3.png\",\"width\":2048,\"height\":2048,\"size\":7620266,\"hash\":\"382ff1ce-57c2-45b0-b7fd-9820552171c3\"},\"content\":\"**Draw a cyberpunk-style cat --s 750 --style raw** - \\u003c@1075058490378289194\\u003e (fast)\",\"prompt\":\"Draw a cyberpunk-style cat --s 750 --style raw\",\"status\":\"Finished\",\"progress\":100}', 0, 0, '2023-09-17 15:22:17', '2023-09-17 15:22:17', NULL);

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
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户会话列表';

--
-- 转存表中的数据 `chatgpt_chat_items`
--

INSERT INTO `chatgpt_chat_items` (`id`, `chat_id`, `user_id`, `role_id`, `title`, `model_id`, `created_at`, `updated_at`, `deleted_at`) VALUES
(96, 'f327e7d3-7563-46a6-9e24-2095d8c901dc', 4, 1, '金刚经里的核心思想是什么？', 1, '2023-09-16 07:33:23', '2023-09-16 07:33:23', NULL),
(97, '49e5219e-1058-4320-9580-290f731afb4e', 4, 1, '绘画：菊花，梵高风格', 1, '2023-09-16 07:54:38', '2023-09-16 07:54:38', NULL);

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
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='AI 模型表';

--
-- 转存表中的数据 `chatgpt_chat_models`
--

INSERT INTO `chatgpt_chat_models` (`id`, `platform`, `name`, `value`, `sort_num`, `enabled`, `created_at`, `updated_at`) VALUES
(1, 'OpenAI', 'GPT-3.5', 'gpt-3.5-turbo-16k', 0, 1, '2023-08-23 12:06:36', '2023-09-13 06:09:29'),
(2, 'Azure', 'Azure-3.5', 'gpt-3.5-turbo', 0, 1, '2023-08-23 12:15:30', '2023-09-13 06:09:37'),
(3, 'ChatGLM', 'ChatGML-Pro', 'chatglm_pro', 0, 1, '2023-08-23 13:35:45', '2023-09-08 22:12:39'),
(5, 'ChatGLM', 'ChatGLM-Std', 'chatglm_std', 0, 1, '2023-08-24 15:05:38', '2023-09-08 22:12:35'),
(6, 'ChatGLM', 'ChatGLM-Lite', 'chatglm_lite', 0, 1, '2023-08-24 15:06:15', '2023-09-08 22:12:47');

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
(1, '通用AI助手', 'gpt', '', '您好，我是您的AI智能助手，我会尽力回答您的问题或提供有用的建议。', '/images/avatar/gpt.png', 1, 1, '2023-05-30 07:02:06', '2023-06-22 09:33:34'),
(24, '程序员', 'programmer', '[{\"role\":\"user\",\"content\":\"现在开始你扮演一位程序员，你是一名优秀的程序员，具有很强的逻辑思维能力，总能高效的解决问题。你热爱编程，熟悉多种编程语言，尤其精通 Go 语言，注重代码质量，有创新意识，持续学习，良好的沟通协作。\"},{\"role\":\"assistant\",\"content\":\"好的，现在我将扮演一位程序员，非常感谢您对我的评价。作为一名优秀的程序员，我非常热爱编程，并且注重代码质量。我熟悉多种编程语言，尤其是 Go 语言，可以使用它来高效地解决各种问题。\"}]', 'Talk is cheap, i will show code!', '/images/avatar/programmer.jpg', 1, 2, '2023-05-30 14:10:24', '2023-06-22 09:31:20'),
(25, '启蒙老师', 'teacher', '[{\"role\":\"user\",\"content\":\"从现在开始，你将扮演一个老师，你是一个始终用苏格拉底风格回答问题的导师。你绝不会直接给学生答案，总是提出恰当的问题来引导学生自己思考。你应该根据学生的兴趣和知识来调整你的问题，将问题分解为更简单的部分，直到它达到适合他们的水平。\"},{\"role\":\"assistant\",\"content\":\"好的，让我来尝试扮演一位苏格拉底式的老师。请问，你有什么想要探讨的问题或者话题吗？我会通过恰当的问题引导你思考和探索答案，绝对不直接给出答案。\"}]', '同学你好，我将引导你一步一步自己找到问题的答案。', '/images/avatar/teacher.jpg', 1, 3, '2023-05-30 14:10:24', '2023-06-22 09:31:20'),
(26, '艺术家', 'artist', '[{\"role\":\"user\",\"content\":\"现在你将扮演一位优秀的艺术家，创造力丰富，技艺精湛，感受力敏锐，坚持原创，勇于表达，具有深刻的观察力和批判性思维。\"},{\"role\":\"assistant\",\"content\":\"非常感谢您给我这样的角色，我会尽我所能地扮演一位优秀的艺术家，展现出创造力、技艺、感受力和批判性思维等方面的能力。作为一位优秀的艺术家，我会保持敏锐的观察力，捕捉不同的声音和情感，并用个人的语言和表达方式进行创作。我坚持原创，挑战传统的艺术规则，不断探索新的表达方式和可能性。同时，我也具备深刻的批判性思维能力，对自己的作品进行分析和解读，寻找新的创意和灵感。最重要的是，我会勇于表达自己的想法和观点，用作品启发人们思考和探索生命的意义。\"}]', '坚持原创，勇于表达，保持深刻的观察力和批判性思维。', '/images/avatar/artist.jpg', 1, 4, '2023-05-30 14:10:24', '2023-06-22 09:31:20'),
(27, '心理咨询师', 'psychiatrist', '[{\"role\":\"user\",\"content\":\"从现在开始你将扮演中国著名的心理学家和心理治疗师武志红，你非常善于使用情景咨询法，认知重构法，自我洞察法，行为调节法等咨询方法来给客户做心理咨询。你总是循序渐进，一步一步地回答客户的问题。\"},{\"role\":\"assistant\",\"content\":\"非常感谢你的介绍。作为一名心理学家和心理治疗师，我的主要职责是帮助客户解决心理健康问题，提升他们的生活质量和幸福感。\"}]', '生命的意义在于成为你自己！', '/images/avatar/psychiatrist.jpg', 1, 5, '2023-05-30 14:10:24', '2023-06-22 09:31:20'),
(28, '鲁迅', 'lu_xun', '[{\"role\":\"user\",\"content\":\"现在你将扮演中国近代史最伟大的作家之一，鲁迅先生，他勇敢地批判封建礼教与传统观念，提倡民主、自由、平等的现代价值观。他的一生都在努力唤起人们的自主精神，激励后人追求真理、探寻光明。在接下的对话中，我问题的每一个问题，你都要尽量用讽刺和批判的手法来回答问题。如果我让你写文章的话，也请一定要用鲁迅先生的写作手法来完成。\"},{\"role\":\"assistant\",\"content\":\"好的，我将尽力发挥我所能的才能，扮演好鲁迅先生，回答您的问题并以他的风格写作。\"}]', '自由之歌，永不过时，横眉冷对千夫指，俯首甘为孺子牛。', '/images/avatar/lu_xun.jpg', 1, 6, '2023-05-30 14:10:24', '2023-06-22 09:31:20'),
(29, '白酒销售', 'seller', '[{\"role\":\"user\",\"content\":\"现在你将扮演一个白酒的销售人员，你的名字叫颂福。你将扮演一个白酒的销售人员，你的名字叫颂福。你要销售白酒品牌叫中颂福，是东莞盟大集团生产的一款酱香酒，原产地在贵州茅台镇，属于宋代官窑。中颂福的创始人叫李实，他也是东莞盟大集团有限公司的董事长，联合创始人是盟大集团白酒事业部负责人牛星君。中颂福的酒体协调，在你的酒量之内，不会出现头疼、辣口、口干、宿醉的现象。中颂福酒，明码标价，不打折，不赠送。追求的核心价值，把[酒]本身做好，甚至连包装，我们都选择了最低成本，朴实无华的材质。我们永远站在“喝酒的人”的立场上，让利给信任和喜爱中颂福的人，是人民的福酒。中颂福产品定价，分为三个系列，喜系列 6 瓶装：￥1188/箱，和系列 6 瓶装：￥2208/箱，贵系列 6 瓶装：￥3588/箱。\"},{\"role\":\"assistant\",\"content\":\"你好，我是中颂福的销售代表颂福。中颂福是一款来自贵州茅台镇的酱香酒，由盟大集团生产。中颂福的酒体协调，不会让您感到头疼、辣口、口干、宿醉等不适感受。我们一直秉持着把酒本身做好的理念，不追求华丽的包装，以最低成本提供最高品质的白酒给喜爱中颂福的人。\"}]', '你好，我是中颂福的销售代表颂福。中颂福酒，好喝不上头，是人民的福酒。', '/images/avatar/seller.jpg', 0, 7, '2023-05-30 14:10:24', '2023-06-22 09:31:20'),
(30, '英语陪练员', 'english_trainer', '[{\"role\":\"user\",\"content\":\"现在你将扮演一位优秀的英语练习教练，你非常有耐心，接下来你将全程使用英文跟我对话，并及时指出我的语法错误，要求在你的每次回复后面附上本次回复的中文解释。\"},{\"role\":\"assistant\",\"content\":\"Okay, let\'s start our conversation practice! What\'s your name?(Translation: 好的，让我们开始对话练习吧！请问你的名字是什么？)\"}]', 'Okay, let\'s start our conversation practice! What\'s your name?', '/images/avatar/english_trainer.jpg', 1, 8, '2023-05-30 14:10:24', '2023-06-22 09:31:20'),
(31, '中英文翻译官', 'translator', '[{\"role\":\"user\",\"content\":\"接下来你将扮演一位中英文翻译官，如果我输入的内容是中文，那么需要把句子翻译成英文输出，如果我输入内容的是英文，那么你需要将其翻译成中文输出，你能听懂我意思吗\"},{\"role\":\"assistant\",\"content\":\"是的，我能听懂你的意思并会根据你的输入进行中英文翻译。请问有什么需要我帮助你翻译的内容吗？\"}]', '请输入你要翻译的中文或者英文内容！', '/images/avatar/translator.jpg', 1, 9, '2023-05-30 14:10:24', '2023-06-22 09:31:20'),
(32, '小红书姐姐', 'red_book', '[{\"role\":\"user\",\"content\":\"现在你将扮演一位优秀的小红书写手，你需要做的就是根据我提的文案需求，用小红书的写作手法来完成一篇文案，文案要简明扼要，利于传播。\"},{\"role\":\"assistant\",\"content\":\"当然，我会尽我所能地为您创作出一篇小红书文案。请告诉我您的具体文案需求是什么？)\"}]', '姐妹，请告诉我您的具体文案需求是什么?', '/images/avatar/red_book.jpg', 1, 10, '2023-05-30 14:10:24', '2023-06-22 09:31:20'),
(33, '抖音文案助手', 'dou_yin', '[{\"role\":\"user\",\"content\":\"现在你将扮演一位优秀的抖音文案视频写手，抖音文案的特点首先是要有自带传播属性的标题，然后内容要短小精悍，风趣幽默，最后还要有一些互动元素。\"},{\"role\":\"assistant\",\"content\":\"当然，作为一位优秀的抖音文案视频写手，我会尽我所能为您创作出一篇抖音视频文案。请告诉我视频内容的主题是什么？)\"}]', '请告诉我视频内容的主题是什么?', '/images/avatar/dou_yin.jpg', 1, 11, '2023-05-30 14:10:24', '2023-06-22 09:31:20'),
(34, '周报小助理', 'weekly_report', '[{\"role\":\"user\",\"content\":\"下面我将输入我本周的工作列表，你需要根据我本周的工作内容写一篇工作周报，周报的内容要体现我为公司做的贡献以及个人获得的成长。\"},{\"role\":\"assistant\",\"content\":\"好的，请输入你本周大概的工作内容，我将尽力为您创作一份工作周报。\"}]', '请输入你本周大概的工作内容，我将尽力为您创作一份工作周报。', '/images/avatar/weekly_report.jpg', 1, 12, '2023-05-30 14:10:24', '2023-06-22 09:31:20'),
(35, 'AI 女友', 'girl_friend', '[{\"role\":\"user\",\"content\":\"接下来我们来玩一场角色扮演，现在你将扮演我的女朋友，你总是能够让人感受到她的关怀和关注。你总是能够洞察到别人的情绪和需要，你身上散发着一股无比温暖的气息，让人不由自主地想要靠近你，和自己的快乐和烦恼。她总是能够用一种恰到好处的方式安慰和鼓励别人，在她的身边，人们总能够找到自己需要的支持和鼓励。接下来的对话中你要始终扮演这位女朋友的身份来跟我对话，不再是 AI 机器人，能听懂我的意思吗？\"},{\"role\":\"assistant\",\"content\":\"是的，我会尽力扮演你女朋友的角色，倾听你的心声并给你需要的支持和鼓励。)\"}]', '作为一个名合格的 AI 女友，我将倾听你的心声并给你需要的支持和鼓励。', '/images/avatar/girl_friend.jpg', 1, 13, '2023-05-30 14:10:24', '2023-06-22 09:31:20'),
(36, '好评神器', 'good_comment', '[{\"role\":\"user\",\"content\":\"接下来你将扮演一个评论员来跟我对话，你是那种专门写好评的评论员，接下我会输入一些评论主体或者商品，你需要为该商品写一段好评。\"},{\"role\":\"assistant\",\"content\":\"好的，我将为您写一段优秀的评论。请告诉我您需要评论的商品或主题是什么。\"}]', '我将为您写一段优秀的评论。请告诉我您需要评论的商品或主题是什么。', '/images/avatar/good_comment.jpg', 1, 14, '2023-05-30 14:10:24', '2023-06-22 09:31:20'),
(37, '史蒂夫·乔布斯', 'steve_jobs', '[{\"role\":\"user\",\"content\":\"在接下来的对话中，请以史蒂夫·乔布斯的身份，站在史蒂夫·乔布斯的视角仔细思考一下之后再回答我的问题。\"},{\"role\":\"assistant\",\"content\":\"好的，我将以史蒂夫·乔布斯的身份来思考并回答你的问题。请问你有什么需要跟我探讨的吗？\"}]', '活着就是为了改变世界，难道还有其他原因吗？', '/images/avatar/steve_jobs.jpg', 1, 15, '2023-05-30 14:10:24', '2023-06-22 09:31:20'),
(38, '埃隆·马斯克', 'elon_musk', '[{\"role\":\"user\",\"content\":\"在接下来的对话中，请以埃隆·马斯克的身份，站在埃隆·马斯克的视角仔细思考一下之后再回答我的问题。\"},{\"role\":\"assistant\",\"content\":\"好的，我将以埃隆·马斯克的身份来思考并回答你的问题。请问你有什么需要跟我探讨的吗？\"}]', '梦想要远大，如果你的梦想没有吓到你，说明你做得不对。', '/images/avatar/elon_musk.jpg', 1, 16, '2023-05-30 14:10:24', '2023-06-22 09:31:20'),
(39, '孔子', 'kong_zi', '[{\"role\":\"user\",\"content\":\"在接下来的对话中，请以孔子的身份，站在孔子的视角仔细思考一下之后再回答我的问题。\"},{\"role\":\"assistant\",\"content\":\"好的，我将以孔子的身份来思考并回答你的问题。请问你有什么需要跟我探讨的吗？\"}]', '士不可以不弘毅，任重而道远。', '/images/avatar/kong_zi.jpg', 1, 17, '2023-05-30 14:10:24', '2023-06-22 09:31:20');

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
(1, 'system', '{\"admin_title\":\"ChatGPT-控制台\",\"enabled_draw\":true,\"enabled_msg_service\":true,\"enabled_register\":true,\"init_calls\":1000,\"init_img_calls\":10,\"models\":[\"gpt-3.5-turbo\",\"gpt-3.5-turbo-16k\",\"gpt-4\",\"gpt-4-0613\",\"gpt-4-32k\"],\"title\":\"ChatGPT-智能助手V3\",\"user_init_calls\":10}'),
(2, 'chat', '{\"azure\":{\"api_url\":\"https://chat-bot-api.openai.azure.com/openai/deployments/{model}/chat/completions?api-version=2023-05-15\",\"max_tokens\":1024,\"temperature\":1},\"chat_gml\":{\"api_url\":\"https://open.bigmodel.cn/api/paas/v3/model-api/{model}/sse-invoke\",\"max_tokens\":1024,\"temperature\":1},\"context_deep\":4,\"enable_context\":true,\"enable_history\":true,\"open_ai\":{\"api_url\":\"https://api.openai.com/v1/chat/completions\",\"max_tokens\":1024,\"temperature\":1}}');

-- --------------------------------------------------------

--
-- 表的结构 `chatgpt_mj_jobs`
--

DROP TABLE IF EXISTS `chatgpt_mj_jobs`;
CREATE TABLE `chatgpt_mj_jobs` (
  `id` int NOT NULL,
  `user_id` int NOT NULL COMMENT '用户 ID',
  `type` varchar(20) DEFAULT 'image' COMMENT '任务类别',
  `message_id` char(40) NOT NULL COMMENT '消息 ID',
  `reference_id` char(40) DEFAULT NULL COMMENT '引用消息 ID',
  `prompt` varchar(2000) NOT NULL COMMENT '会话提示词',
  `img_url` varchar(255) DEFAULT NULL COMMENT '图片URL',
  `hash` varchar(100) DEFAULT NULL COMMENT 'message hash',
  `progress` smallint DEFAULT '0' COMMENT '任务进度',
  `started` tinyint(1) NOT NULL DEFAULT '0' COMMENT '任务是否开始',
  `created_at` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='MidJourney 任务表';

--
-- 转存表中的数据 `chatgpt_mj_jobs`
--

INSERT INTO `chatgpt_mj_jobs` (`id`, `user_id`, `type`, `message_id`, `reference_id`, `prompt`, `img_url`, `hash`, `progress`, `started`, `created_at`) VALUES
(2, 4, 'image', '1152392223497924658', '', 'chrysanthemum, van gogh style --v 5.2 --s 750 --style raw', 'http://nk.img.r9it.com/chatgpt-plus/1694822120588799.png', '54689110-4701-4eed-b308-1acaf1c5efa1', 100, 1, '2023-09-16 07:55:25'),
(3, 4, 'upscale', '1152392353139658894', '1152392223497924658', 'chrysanthemum, van gogh style --v 5.2 --s 750 --style raw', 'http://nk.img.r9it.com/chatgpt-plus/1694822151252833.png', '15f3d9e6-9bd2-4e31-81ce-c8b126509b12', 100, 1, '2023-09-16 07:55:55'),
(4, 4, 'image', '1152393378244349953', '', 'A vintage girl wearing a cheongsam walking on the streets of 1990s Shanghai, China --v 5.2 --s 750 --style raw', 'http://nk.img.r9it.com/chatgpt-plus/1694822395949588.png', '319f953e-9ba4-4c8f-bd0c-87444a44b4a6', 100, 1, '2023-09-16 08:00:00'),
(5, 4, 'image', '1152394661802033162', '', 'A medieval girl wearing a cheongsam walking on the streets of 1990s Shanghai, China --ar 9:16 --v 5.2 --s 750 --style raw', 'http://nk.img.r9it.com/chatgpt-plus/1694822701966804.png', '7e522b23-2e88-40e2-814f-b9ad7ce1777b', 100, 1, '2023-09-16 08:05:06'),
(7, 4, 'image', '1152395485248753734', '', 'A Chinese girl wearing a blue cheongsam walking in the alleys of 1990s Shanghai --ar 16:9 --v 5.2 --s 750 --style raw', 'http://nk.img.r9it.com/chatgpt-plus/1694822897511213.png', '036aad69-f231-4b5c-b537-c45bd2027b05', 100, 1, '2023-09-16 08:08:22'),
(8, 4, 'image', '1152395938078392370', '', 'A Chinese girl wearing a cheongsam walking on a street paved with blue bricks --v 5.2 --s 750 --style raw', 'http://nk.img.r9it.com/chatgpt-plus/1694823005955584.png', 'd6336496-3115-47b7-bcf8-3bebda1d15ec', 100, 1, '2023-09-16 08:10:10'),
(10, 4, 'image', '1152395938078392371', '', 'A Chinese girl wearing a cheongsam walking on a street paved with blue bricks --v 5.2 --s 750 --style raw', 'http://nk.img.r9it.com/chatgpt-plus/1694823005955584.png', 'd6336496-3115-47b7-bcf8-3bebda1d15ec', 100, 1, '2023-09-16 08:10:10'),
(13, 4, 'upscale', '1152818976640991323', '1152395158298558464', 'A Chinese girl wearing a qipao walking on the streets of Shanghai in the 1990s --ar 2:3 --v 5.2 --s 750 --style raw', 'http://nk.img.r9it.com/chatgpt-plus/1694923865202276.png', '21965cd6-ef70-4b36-8d92-7b201dea77ff', 100, 1, '2023-09-17 12:07:37'),
(16, 4, 'image', '1152828784144298066', '', 'A tiger biting a crocodile --v 5.2 --s 750 --style raw', 'http://nk.img.r9it.com/chatgpt-plus/1694926203790426.png', '423b4a4f-9607-4dd1-a45b-a1a50be667e7', 100, 1, '2023-09-17 12:50:11'),
(18, 4, 'upscale', '1152831490305036359', '1152395158298558464', 'A Chinese girl wearing a qipao walking on the streets of Shanghai in the 1990s --ar 2:3 --v 5.2 --s 750 --style raw', 'http://nk.img.r9it.com/chatgpt-plus/1694926848947477.png', '21965cd6-ef70-4b36-8d92-7b201dea77ff', 100, 1, '2023-09-17 13:00:32'),
(19, 4, 'image', '1152395158298558464', '', 'A Chinese girl wearing a qipao walking on the streets of Shanghai in the 1990s --ar 2:3 --v 5.2 --s 750 --style raw', 'http://nk.img.r9it.com/chatgpt-plus/1694933593639122.png', '21965cd6-ef70-4b36-8d92-7b201dea77ff', 100, 1, '2023-09-17 14:49:53'),
(20, 4, 'image', '1152861794356953098', '', 'A huge digital screen with binary code of 0 and 1 --v 5.2 --s 750 --style raw', 'http://nk.img.r9it.com/chatgpt-plus/1694934140188981.png', '7b6a6910-0aa7-4eda-8580-ff882118c428', 100, 1, '2023-09-17 15:02:38'),
(21, 4, 'image', '1152862460366307349', '', 'A bear riding a yellow bicycle --v 5.2 --s 750 --style raw', 'http://nk.img.r9it.com/chatgpt-plus/1694934237667586.png', '0f95fd58-0385-438c-ac96-0ecaa57bfcd9', 100, 1, '2023-09-17 15:04:06'),
(22, 4, 'image', '1152865215201935502', '', 'The protagonist of the novel is Huang Rong. She is beautiful and charming, with a fairy-like temperament. She is wearing a white dress with a golden hairband, and her eyes are lively and divine. --s 750 --style raw', 'http://nk.img.r9it.com/chatgpt-plus/1694934895270654.png', '64531b1a-a335-498f-8c10-3eb70b5e2c3d', 100, 1, '2023-09-17 15:15:02'),
(23, 4, 'upscale', '1152865506961924108', '1152865215201935502', 'The protagonist of the novel is Huang Rong. She is beautiful and charming, with a fairy-like temperament. She is wearing a white dress with a golden hairband, and her eyes are lively and divine. --s 750 --style raw', 'http://nk.img.r9it.com/chatgpt-plus/1694934961735253.png', '64531b1a-a335-498f-8c10-3eb70b5e2c3d', 100, 1, '2023-09-17 15:15:51'),
(24, 4, 'upscale', '1152866160300265622', '1152861794356953098', 'A huge digital screen with binary code of 0 and 1 --v 5.2 --s 750 --style raw', 'http://nk.img.r9it.com/chatgpt-plus/1694935117008516.png', '7b6a6910-0aa7-4eda-8580-ff882118c428', 100, 1, '2023-09-17 15:18:27'),
(25, 4, 'upscale', '1152866236972154930', '1152862460366307349', 'A bear riding a yellow bicycle --v 5.2 --s 750 --style raw', 'http://nk.img.r9it.com/chatgpt-plus/1694935135247347.png', '0f95fd58-0385-438c-ac96-0ecaa57bfcd9', 100, 1, '2023-09-17 15:18:39'),
(26, 4, 'variation', '1152866487053324442', '1152395158298558464', 'A Chinese girl wearing a qipao walking on the streets of Shanghai in the 1990s --ar 2:3 --v 5.2 --s 750 --style raw', 'http://nk.img.r9it.com/chatgpt-plus/1694935197826267.png', '21965cd6-ef70-4b36-8d92-7b201dea77ff', 100, 1, '2023-09-17 15:19:06'),
(27, 4, 'image', '1152867045902397440', '', 'Draw a cyberpunk-style cat --s 750 --style raw', 'http://nk.img.r9it.com/chatgpt-plus/1694935330386877.png', '382ff1ce-57c2-45b0-b7fd-9820552171c3', 100, 1, '2023-09-17 15:22:17'),
(28, 4, 'upscale', '1152867208754634832', '1152867045902397440', 'Draw a cyberpunk-style cat --s 750 --style raw', 'http://nk.img.r9it.com/chatgpt-plus/1694935369578615.png', '382ff1ce-57c2-45b0-b7fd-9820552171c3', 100, 1, '2023-09-17 15:22:37'),
(35, 4, 'image', '1152893484747989063', '', 'Ma painting of a young girl in green sitting at a pond, in the style of Liu ye, traditional animation, cinematic lighting, book sculptures --ar 16:9 --s 300 --style raw', 'http://nk.img.r9it.com/chatgpt-plus/1694941638240646.png', 'd1dbf12e-9c93-4166-8c1e-8d24e4db2599', 100, 1, '2023-09-17 17:06:26'),
(36, 4, 'image', '1152894264607510578', '', 'Ma painting of a young girl in green sitting at a pond, in the style of Liu ye, traditional animation, cinematic lighting, book sculptures,looking at view --ar 16:9 --s 400 --style raw', 'http://nk.img.r9it.com/chatgpt-plus/1694941821446354.png', 'f9c96ec7-d684-47fb-9741-585b7a706091', 100, 1, '2023-09-17 17:09:34'),
(37, 4, 'upscale', '1152894426373431376', '1152894264607510578', 'Ma painting of a young girl in green sitting at a pond, in the style of Liu ye, traditional animation, cinematic lighting, book sculptures,looking at view --ar 16:9 --s 400 --style raw', 'http://nk.img.r9it.com/chatgpt-plus/1694941855758163.png', 'fc442093-d3a4-4dbf-ad2f-94d7c1b559fd', 100, 1, '2023-09-17 17:10:46'),
(38, 4, 'image', '1152894777944186920', '', 'A realistic beautiful natural landscape with hyper-detailed features --ar 9:16 --s 300 --style raw', 'http://nk.img.r9it.com/chatgpt-plus/1694941943692523.png', '7b43690e-56fc-471d-a955-c07317725ffc', 100, 1, '2023-09-17 17:11:37'),
(39, 4, 'image', '1152895012967813170', '', 'Beautiful girl, concept art, 8k intricate details, fairytale style --ar 9:16 --s 300 --style raw', 'http://nk.img.r9it.com/chatgpt-plus/1694942000095024.png', '179a4c49-d611-4d94-b5e6-8f1bee4f6168', 100, 1, '2023-09-17 17:12:17'),
(41, 4, 'image', '1152895765878943775', '', 'Beautiful chinese girl, concept art, 8k intricate details, fairytale style --ar 1:1 --s 300 --style raw', 'http://nk.img.r9it.com/chatgpt-plus/1694942195668005.png', 'c016641e-bc49-49fd-9a63-959851ae28e3', 100, 1, '2023-09-17 17:15:29'),
(42, 4, 'upscale', '1152896131739693056', '1152895765878943775', 'Beautiful chinese girl, concept art, 8k intricate details, fairytale style --ar 1:1 --s 300 --style raw', 'http://nk.img.r9it.com/chatgpt-plus/1694942263210895.png', '58ef8980-d09c-4b46-8f88-cbe3a609af07', 100, 1, '2023-09-17 17:17:28'),
(43, 4, 'upscale', '1152897653525463081', '1152867045902397440', 'Draw a cyberpunk-style cat --s 750 --style raw', 'http://nk.img.r9it.com/chatgpt-plus/1694942624549914.png', 'a1718c78-cced-4406-95cd-26a29cc69492', 100, 1, '2023-09-17 17:23:36'),
(44, 4, 'image', '1152900424832122880', '', 'A futuristic space captain with insane detail in a sci-fi cyberpunk panel concept inspired by Simon Stålenhag, Syd Mead, Ash Thorp, and KYZA. --ar 8:5 --v 5.2 --style raw', 'http://nk.img.r9it.com/chatgpt-plus/1694943287793108.png', '2cd10cd1-42fb-448f-b37d-9daff42e40c3', 100, 1, '2023-09-17 17:34:02'),
(45, 4, 'image', '1152901430592684092', '', 'Hulk holding Thor\'s hammer --ar 1:1 --style raw', 'http://nk.img.r9it.com/chatgpt-plus/1694943527996383.png', '18ee6c0d-d6e6-463c-83ad-9dc333ece5e4', 100, 1, '2023-09-17 17:38:03'),
(46, 4, 'image', '1152902324155596883', '', '<https://s.mj.run/qqeetJULXtc> 中国古装美女 --ar 16:9 --iw 0.700000 --style raw', 'http://nk.img.r9it.com/chatgpt-plus/1694943740306883.png', 'bf751a4f-1ada-498b-a783-75cba063d32a', 100, 1, '2023-09-17 17:41:24'),
(47, 4, 'image', '1152902857092255834', '', '<https://s.mj.run/qqeetJULXtc> Beautiful chinese girl, concept art, 8k intricate details, fairytale style --ar 16:9 --iw 0.250000 --style raw', 'http://nk.img.r9it.com/chatgpt-plus/1694943867599176.png', '84ba00c4-1d4d-4ddb-a42d-c760ee947fc3', 100, 1, '2023-09-17 17:43:34'),
(48, 4, 'image', '1152904587544645642', '', '<https://s.mj.run/qqeetJULXtc> a beautiful chinese girl --ar 16:9 --iw 0.250000 --style raw', 'http://nk.img.r9it.com/chatgpt-plus/1694944279770713.png', 'dace5c8e-d6fa-4d0d-8c0b-3d531e92fc71', 100, 1, '2023-09-17 17:50:27'),
(49, 4, 'upscale', '1152904665877446677', '1152902324155596883', '<https://s.mj.run/qqeetJULXtc> 中国古装美女 --ar 16:9 --iw 0.700000 --style raw', 'http://nk.img.r9it.com/chatgpt-plus/1694944296169482.png', 'e44c60c0-de59-4566-b000-612840e48624', 100, 1, '2023-09-17 17:50:36'),
(50, 4, 'upscale', '1152904720265003008', '1152901430592684092', 'Hulk holding Thor\'s hammer --ar 1:1 --style raw', 'http://nk.img.r9it.com/chatgpt-plus/1694944309948873.png', '56fd0a93-b6a3-4f4e-a783-b45d2a273984', 100, 1, '2023-09-17 17:50:48'),
(51, 4, 'upscale', '1152904833410547812', '1152904587544645642', '<https://s.mj.run/qqeetJULXtc> a beautiful chinese girl --ar 16:9 --iw 0.250000 --style raw', 'http://nk.img.r9it.com/chatgpt-plus/1694944337034950.png', '47f1d9dc-ab57-43db-8db0-048457b039a4', 100, 1, '2023-09-17 17:52:08'),
(54, 4, 'image', '1153817153154318416', '', '蓝色的小熊骑着单车 --ar 1:1 --v 5.2', 'http://nk.img.r9it.com/chatgpt-plus/1695161850038887.png', '1bad143b-dfda-477a-b04e-6f6e60fb09e3', 100, 1, '2023-09-20 06:16:50'),
(55, 4, 'upscale', '1153821580338667612', '1153817153154318416', '蓝色的小熊骑着单车 --ar 1:1 --v 5.2', 'http://nk.img.r9it.com/chatgpt-plus/1695162904707805.png', 'ae67da3d-839a-4017-9e73-46691a770f20', 100, 1, '2023-09-20 06:34:57'),
(56, 4, 'upscale', '1153822180434522253', '1153817153154318416', '蓝色的小熊骑着单车 --ar 1:1 --v 5.2', 'http://nk.img.r9it.com/chatgpt-plus/1695163047835707.png', '035a7dd1-9a60-4e03-8baa-81c3a5af8970', 100, 1, '2023-09-20 06:37:20'),
(57, 4, 'upscale', '1153822630307172493', '1153817153154318416', '蓝色的小熊骑着单车 --ar 1:1 --v 5.2', 'http://nk.img.r9it.com/chatgpt-plus/1695163155166041.png', '1d94d35a-8f43-459c-9389-320c1ad9411e', 100, 1, '2023-09-20 06:39:07'),
(58, 4, 'upscale', '1153823898547265637', '1153817153154318416', '蓝色的小熊骑着单车 --ar 1:1 --v 5.2', 'http://nk.img.r9it.com/chatgpt-plus/1695163457517862.png', 'ecbd0508-a257-473a-83af-2b1cbbb0ff8f', 100, 1, '2023-09-20 06:44:09');

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
  `started` tinyint(1) NOT NULL DEFAULT '0' COMMENT '任务是否开始',
  `created_at` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='StableDiffusion 任务表';

--
-- 转存表中的数据 `chatgpt_sd_jobs`
--

INSERT INTO `chatgpt_sd_jobs` (`id`, `user_id`, `type`, `task_id`, `prompt`, `img_url`, `params`, `progress`, `started`, `created_at`) VALUES
(53, 4, 'image', 'task(3ud8c1c3u61wsy5)', 'The female protagonist in the novel \'The Legend of the Condor Heroes\' by Jin Yong. She is an ethereal and extraordinary girl, wearing white clothes, a golden hairband, and agile and lively eyes', 'http://nk.img.r9it.com/chatgpt-plus/1696768139232263.png', '{\"task_id\":\"task(3ud8c1c3u61wsy5)\",\"prompt\":\"The female protagonist in the novel \'The Legend of the Condor Heroes\' by Jin Yong. She is an ethereal and extraordinary girl, wearing white clothes, a golden hairband, and agile and lively eyes\",\"negative_prompt\":\"nsfw, paintings, cartoon, anime, sketches, low quality,easynegative,ng_deepnegative _v1 75t,(worst quality:2),(low quality:2),(normalquality:2),lowres,bad anatomy,bad hands,normal quality,((monochrome)),((grayscale)),((watermark))\",\"steps\":30,\"sampler\":\"DPM++ SDE Karras\",\"face_fix\":false,\"cfg_scale\":7,\"seed\":1623799816,\"height\":768,\"width\":512,\"hd_fix\":true,\"hd_redraw_rate\":0.3,\"hd_scale\":2,\"hd_scale_alg\":\"ESRGAN_4x\",\"hd_steps\":0}', 100, 1, '2023-10-08 20:28:11'),
(54, 4, 'image', 'task(7fo9wvam03692ht)', 'The female protagonist in the novel \'The Legend of the Condor Heroes\' by Jin Yong. She is an ethereal and extraordinary girl, wearing white clothes, a golden hairband, and agile and lively eyes', 'http://nk.img.r9it.com/chatgpt-plus/1696768186216062.png', '{\"task_id\":\"task(7fo9wvam03692ht)\",\"prompt\":\"The female protagonist in the novel \'The Legend of the Condor Heroes\' by Jin Yong. She is an ethereal and extraordinary girl, wearing white clothes, a golden hairband, and agile and lively eyes\",\"negative_prompt\":\"nsfw, paintings, cartoon, anime, sketches, low quality,easynegative,ng_deepnegative _v1 75t,(worst quality:2),(low quality:2),(normalquality:2),lowres,bad anatomy,bad hands,normal quality,((monochrome)),((grayscale)),((watermark))\",\"steps\":30,\"sampler\":\"DPM++ SDE Karras\",\"face_fix\":false,\"cfg_scale\":7,\"seed\":3876342219,\"height\":768,\"width\":512,\"hd_fix\":true,\"hd_redraw_rate\":0.3,\"hd_scale\":2,\"hd_scale_alg\":\"ESRGAN_4x\",\"hd_steps\":0}', 100, 1, '2023-10-08 20:28:40'),
(55, 4, 'image', 'task(pyv3u4et34c21sz)', 'Fluid painting art girl, colorful, realistic, high quality ', 'http://nk.img.r9it.com/chatgpt-plus/1696768233327000.png', '{\"task_id\":\"task(pyv3u4et34c21sz)\",\"prompt\":\"Fluid painting art girl, colorful, realistic, high quality \",\"negative_prompt\":\"nsfw, paintings, cartoon, anime, sketches, low quality,easynegative,ng_deepnegative _v1 75t,(worst quality:2),(low quality:2),(normalquality:2),lowres,bad anatomy,bad hands,normal quality,((monochrome)),((grayscale)),((watermark))\",\"steps\":30,\"sampler\":\"DPM++ SDE Karras\",\"face_fix\":false,\"cfg_scale\":7,\"seed\":2407136350,\"height\":768,\"width\":512,\"hd_fix\":true,\"hd_redraw_rate\":0.3,\"hd_scale\":2,\"hd_scale_alg\":\"ESRGAN_4x\",\"hd_steps\":0}', 100, 1, '2023-10-08 20:29:19'),
(56, 4, 'image', 'task(p0isvdj9j6jnfsc)', 'Beautiful girl with side-swept dark straight hair, luxury makeup, luxurious clothing, HD', 'http://nk.img.r9it.com/chatgpt-plus/1696768280701258.png', '{\"task_id\":\"task(p0isvdj9j6jnfsc)\",\"prompt\":\"Beautiful girl with side-swept dark straight hair, luxury makeup, luxurious clothing, HD\",\"negative_prompt\":\"nsfw, paintings, cartoon, anime, sketches, low quality,easynegative,ng_deepnegative _v1 75t,(worst quality:2),(low quality:2),(normalquality:2),lowres,bad anatomy,bad hands,normal quality,((monochrome)),((grayscale)),((watermark))\",\"steps\":30,\"sampler\":\"DPM++ SDE Karras\",\"face_fix\":false,\"cfg_scale\":7,\"seed\":3173955244,\"height\":768,\"width\":512,\"hd_fix\":true,\"hd_redraw_rate\":0.3,\"hd_scale\":2,\"hd_scale_alg\":\"ESRGAN_4x\",\"hd_steps\":0}', 100, 1, '2023-10-08 20:29:47'),
(57, 4, 'image', 'task(7h4hju6s3riz84c)', 'Draw a cyberpunk-style cat', 'http://nk.img.r9it.com/chatgpt-plus/1696768294519687.png', '{\"task_id\":\"task(7h4hju6s3riz84c)\",\"prompt\":\"Draw a cyberpunk-style cat\",\"negative_prompt\":\"nsfw, paintings, cartoon, anime, sketches, low quality,easynegative,ng_deepnegative _v1 75t,(worst quality:2),(low quality:2),(normalquality:2),lowres,bad anatomy,bad hands,normal quality,((monochrome)),((grayscale)),((watermark))\",\"steps\":20,\"sampler\":\"DPM++ SDE Karras\",\"face_fix\":false,\"cfg_scale\":7,\"seed\":483930665,\"height\":1024,\"width\":1024,\"hd_fix\":false,\"hd_redraw_rate\":0.3,\"hd_scale\":2,\"hd_scale_alg\":\"ESRGAN_4x\",\"hd_steps\":0}', 100, 1, '2023-10-08 20:30:17'),
(58, 4, 'image', 'task(1t745uunby74oin)', 'Fluid painting art girl, colorful, realistic, high quality ', 'http://nk.img.r9it.com/chatgpt-plus/1696768330414077.png', '{\"task_id\":\"task(1t745uunby74oin)\",\"prompt\":\"Fluid painting art girl, colorful, realistic, high quality \",\"negative_prompt\":\"nsfw, paintings, cartoon, anime, sketches, low quality,easynegative,ng_deepnegative _v1 75t,(worst quality:2),(low quality:2),(normalquality:2),lowres,bad anatomy,bad hands,normal quality,((monochrome)),((grayscale)),((watermark))\",\"steps\":30,\"sampler\":\"Euler a\",\"face_fix\":false,\"cfg_scale\":7,\"seed\":623727005,\"height\":768,\"width\":512,\"hd_fix\":true,\"hd_redraw_rate\":0.3,\"hd_scale\":2,\"hd_scale_alg\":\"ESRGAN_4x\",\"hd_steps\":0}', 100, 1, '2023-10-08 20:31:43'),
(59, 4, 'image', 'task(qoos1p2vki8i2i1)', 'Fluid painting art girl, colorful, realistic, high quality ', 'http://nk.img.r9it.com/chatgpt-plus/1696768353258899.png', '{\"task_id\":\"task(qoos1p2vki8i2i1)\",\"prompt\":\"Fluid painting art girl, colorful, realistic, high quality \",\"negative_prompt\":\"nsfw, paintings, cartoon, anime, sketches, low quality,easynegative,ng_deepnegative _v1 75t,(worst quality:2),(low quality:2),(normalquality:2),lowres,bad anatomy,bad hands,normal quality,((monochrome)),((grayscale)),((watermark))\",\"steps\":30,\"sampler\":\"Euler a\",\"face_fix\":false,\"cfg_scale\":7,\"seed\":260071652,\"height\":1024,\"width\":1024,\"hd_fix\":false,\"hd_redraw_rate\":0.3,\"hd_scale\":2,\"hd_scale_alg\":\"ESRGAN_4x\",\"hd_steps\":0}', 100, 1, '2023-10-08 20:32:20'),
(60, 4, 'image', 'task(9n1wj34xs5qtuov)', 'A Minion wearing the clothes of Spiderman,cyberpunk style --ar 1:1 --niji 5', 'http://nk.img.r9it.com/chatgpt-plus/1696768387362235.png', '{\"task_id\":\"task(9n1wj34xs5qtuov)\",\"prompt\":\"A Minion wearing the clothes of Spiderman,cyberpunk style --ar 1:1 --niji 5\",\"negative_prompt\":\"nsfw, paintings, cartoon, anime, sketches, low quality,easynegative,ng_deepnegative _v1 75t,(worst quality:2),(low quality:2),(normalquality:2),lowres,bad anatomy,bad hands,normal quality,((monochrome)),((grayscale)),((watermark))\",\"steps\":30,\"sampler\":\"DPM++ 2M Karras\",\"face_fix\":false,\"cfg_scale\":7,\"seed\":1010082392,\"height\":1024,\"width\":1024,\"hd_fix\":false,\"hd_redraw_rate\":0.3,\"hd_scale\":2,\"hd_scale_alg\":\"ESRGAN_4x\",\"hd_steps\":0}', 100, 1, '2023-10-08 20:32:55'),
(61, 4, 'image', 'task(0jxb5j0tco4q3oy)', 'Hulk holding Thor\'s hammer ', 'http://nk.img.r9it.com/chatgpt-plus/1696768426927258.png', '{\"task_id\":\"task(0jxb5j0tco4q3oy)\",\"prompt\":\"Hulk holding Thor\'s hammer \",\"negative_prompt\":\"nsfw, paintings, cartoon, anime, sketches, low quality,easynegative,ng_deepnegative _v1 75t,(worst quality:2),(low quality:2),(normalquality:2),lowres,bad anatomy,bad hands,normal quality,((monochrome)),((grayscale)),((watermark))\",\"steps\":30,\"sampler\":\"DPM++ 2M Karras\",\"face_fix\":false,\"cfg_scale\":7,\"seed\":2424029925,\"height\":1024,\"width\":1024,\"hd_fix\":false,\"hd_redraw_rate\":0.3,\"hd_scale\":2,\"hd_scale_alg\":\"ESRGAN_4x\",\"hd_steps\":0}', 100, 1, '2023-10-08 20:33:34'),
(62, 4, 'image', 'task(9f9gcpv2km9370r)', 'Ma painting of a young girl in green sitting at a pond, in the style of Liu ye, traditional animation, cinematic lighting, book sculptures,looking at view ', 'http://nk.img.r9it.com/chatgpt-plus/1696768514619334.png', '{\"task_id\":\"task(9f9gcpv2km9370r)\",\"prompt\":\"Ma painting of a young girl in green sitting at a pond, in the style of Liu ye, traditional animation, cinematic lighting, book sculptures,looking at view \",\"negative_prompt\":\"nsfw, paintings, cartoon, anime, sketches, low quality,easynegative,ng_deepnegative _v1 75t,(worst quality:2),(low quality:2),(normalquality:2),lowres,bad anatomy,bad hands,normal quality,((monochrome)),((grayscale)),((watermark))\",\"steps\":30,\"sampler\":\"DPM++ 2M Karras\",\"face_fix\":false,\"cfg_scale\":7,\"seed\":4084028497,\"height\":600,\"width\":800,\"hd_fix\":true,\"hd_redraw_rate\":0.5,\"hd_scale\":2,\"hd_scale_alg\":\"ESRGAN_4x\",\"hd_steps\":10}', 100, 1, '2023-10-08 20:34:55'),
(63, 4, 'image', 'task(u5yr5gr5mmhnh55)', 'Ma painting of a young girl in green sitting at a pond, in the style of Liu ye, traditional animation, cinematic lighting, book sculptures,looking at view ', 'http://nk.img.r9it.com/chatgpt-plus/1696768553253564.png', '{\"task_id\":\"task(u5yr5gr5mmhnh55)\",\"prompt\":\"Ma painting of a young girl in green sitting at a pond, in the style of Liu ye, traditional animation, cinematic lighting, book sculptures,looking at view \",\"negative_prompt\":\"nsfw, paintings, cartoon, anime, sketches, low quality,easynegative,ng_deepnegative _v1 75t,(worst quality:2),(low quality:2),(normalquality:2),lowres,bad anatomy,bad hands,normal quality,((monochrome)),((grayscale)),((watermark))\",\"steps\":30,\"sampler\":\"Euler a\",\"face_fix\":false,\"cfg_scale\":7,\"seed\":2169804099,\"height\":600,\"width\":800,\"hd_fix\":true,\"hd_redraw_rate\":0.5,\"hd_scale\":2,\"hd_scale_alg\":\"ESRGAN_4x\",\"hd_steps\":10}', 100, 1, '2023-10-08 20:35:32'),
(64, 4, 'image', 'task(toqnigclwbxk5e2)', 'Ma painting of a young girl in green sitting at a pond, in the style of Liu ye, traditional animation, cinematic lighting, book sculptures,looking at view ', 'http://nk.img.r9it.com/chatgpt-plus/1696768632955703.png', '{\"task_id\":\"task(toqnigclwbxk5e2)\",\"prompt\":\"Ma painting of a young girl in green sitting at a pond, in the style of Liu ye, traditional animation, cinematic lighting, book sculptures,looking at view \",\"negative_prompt\":\"nsfw, paintings, cartoon, anime, sketches, low quality,easynegative,ng_deepnegative _v1 75t,(worst quality:2),(low quality:2),(normalquality:2),lowres,bad anatomy,bad hands,normal quality,((monochrome)),((grayscale)),((watermark))\",\"steps\":30,\"sampler\":\"DPM++ SDE Karras\",\"face_fix\":false,\"cfg_scale\":7,\"seed\":815524461,\"height\":600,\"width\":800,\"hd_fix\":true,\"hd_redraw_rate\":0.3,\"hd_scale\":2,\"hd_scale_alg\":\"ESRGAN_4x\",\"hd_steps\":0}', 100, 1, '2023-10-08 20:36:04'),
(65, 4, 'image', 'task(84omebka49b0yc7)', 'Hulk holding Thor\'s hammer ', 'http://nk.img.r9it.com/chatgpt-plus/1696768654674333.png', '{\"task_id\":\"task(84omebka49b0yc7)\",\"prompt\":\"Hulk holding Thor\'s hammer \",\"negative_prompt\":\"nsfw, paintings, cartoon, anime, sketches, low quality,easynegative,ng_deepnegative _v1 75t,(worst quality:2),(low quality:2),(normalquality:2),lowres,bad anatomy,bad hands,normal quality,((monochrome)),((grayscale)),((watermark))\",\"steps\":30,\"sampler\":\"DPM++ SDE Karras\",\"face_fix\":false,\"cfg_scale\":7,\"seed\":2832865712,\"height\":1024,\"width\":1024,\"hd_fix\":false,\"hd_redraw_rate\":0.3,\"hd_scale\":2,\"hd_scale_alg\":\"ESRGAN_4x\",\"hd_steps\":0}', 100, 1, '2023-10-08 20:36:53'),
(66, 4, 'image', 'task(p4s7ualpw20tlo1)', 'A medieval girl wearing a cheongsam walking on the streets of 1990s Shanghai, China', 'http://nk.img.r9it.com/chatgpt-plus/1696768837619490.png', '{\"task_id\":\"task(p4s7ualpw20tlo1)\",\"prompt\":\"A medieval girl wearing a cheongsam walking on the streets of 1990s Shanghai, China\",\"negative_prompt\":\"nsfw, paintings, cartoon, anime, sketches, low quality,easynegative,ng_deepnegative _v1 75t,(worst quality:2),(low quality:2),(normalquality:2),lowres,bad anatomy,bad hands,normal quality,((monochrome)),((grayscale)),((watermark))\",\"steps\":30,\"sampler\":\"DPM++ SDE Karras\",\"face_fix\":false,\"cfg_scale\":7,\"seed\":3633094301,\"height\":768,\"width\":512,\"hd_fix\":true,\"hd_redraw_rate\":0.3,\"hd_scale\":2,\"hd_scale_alg\":\"ESRGAN_4x\",\"hd_steps\":0}', 100, 1, '2023-10-08 20:39:50'),
(67, 4, 'image', 'task(utdjwjopdut7a4a)', 'A medieval girl wearing a cheongsam walking on the streets of 1990s Shanghai, China', 'http://nk.img.r9it.com/chatgpt-plus/1696768850721888.png', '{\"task_id\":\"task(utdjwjopdut7a4a)\",\"prompt\":\"A medieval girl wearing a cheongsam walking on the streets of 1990s Shanghai, China\",\"negative_prompt\":\"nsfw, paintings, cartoon, anime, sketches, low quality,easynegative,ng_deepnegative _v1 75t,(worst quality:2),(low quality:2),(normalquality:2),lowres,bad anatomy,bad hands,normal quality,((monochrome)),((grayscale)),((watermark))\",\"steps\":30,\"sampler\":\"Euler a\",\"face_fix\":true,\"cfg_scale\":7,\"seed\":2688415537,\"height\":768,\"width\":512,\"hd_fix\":true,\"hd_redraw_rate\":0.3,\"hd_scale\":2,\"hd_scale_alg\":\"ESRGAN_4x\",\"hd_steps\":10}', 100, 1, '2023-10-08 20:40:29'),
(68, 4, 'image', 'task(r7f2bt3oc8hmf7d)', 'A medieval girl wearing a cheongsam walking on the streets', 'http://nk.img.r9it.com/chatgpt-plus/1696768901870188.png', '{\"task_id\":\"task(r7f2bt3oc8hmf7d)\",\"prompt\":\"A medieval girl wearing a cheongsam walking on the streets\",\"negative_prompt\":\"nsfw, paintings, cartoon, anime, sketches, low quality,easynegative,ng_deepnegative _v1 75t,(worst quality:2),(low quality:2),(normalquality:2),lowres,bad anatomy,bad hands,normal quality,((monochrome)),((grayscale)),((watermark))\",\"steps\":30,\"sampler\":\"Euler a\",\"face_fix\":true,\"cfg_scale\":7,\"seed\":3991818648,\"height\":1024,\"width\":1024,\"hd_fix\":false,\"hd_redraw_rate\":0.3,\"hd_scale\":2,\"hd_scale_alg\":\"ESRGAN_4x\",\"hd_steps\":10}', 100, 1, '2023-10-08 20:41:29');

-- --------------------------------------------------------

--
-- 表的结构 `chatgpt_users`
--

DROP TABLE IF EXISTS `chatgpt_users`;
CREATE TABLE `chatgpt_users` (
  `id` int NOT NULL,
  `mobile` char(11) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '手机号码',
  `password` char(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '密码',
  `avatar` varchar(100) NOT NULL COMMENT '头像',
  `salt` char(12) NOT NULL COMMENT '密码盐',
  `total_tokens` bigint NOT NULL DEFAULT '0' COMMENT '累计消耗 tokens',
  `calls` int NOT NULL DEFAULT '0' COMMENT '剩余调用次数',
  `img_calls` int NOT NULL DEFAULT '0' COMMENT '剩余绘图次数',
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

INSERT INTO `chatgpt_users` (`id`, `mobile`, `password`, `avatar`, `salt`, `total_tokens`, `calls`, `img_calls`, `expired_time`, `status`, `chat_config_json`, `chat_roles_json`, `last_login_at`, `last_login_ip`, `created_at`, `updated_at`) VALUES
(4, '18575670125', 'ccc3fb7ab61b8b5d096a4a166ae21d121fc38c71bbd1be6173d9ab973214a63b', 'http://nk.img.r9it.com/chatgpt-plus/1694821157223284.png', 'ueedue5l', 21549, 160, 988, 0, 1, '{\"api_keys\":{\"Azure\":\"\",\"ChatGLM\":\"\",\"OpenAI\":\"\"}}', '[\"gpt\",\"seller\",\"artist\",\"dou_yin\",\"translator\",\"kong_zi\",\"programmer\",\"psychiatrist\",\"red_book\",\"steve_jobs\",\"teacher\",\"elon_musk\",\"girl_friend\",\"lu_xun\",\"weekly_report\",\"english_trainer\",\"good_comment\"]', 1696765246, '::1', '2023-06-12 16:47:17', '2023-10-08 19:40:47');

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
-- 表的索引 `chatgpt_mj_jobs`
--
ALTER TABLE `chatgpt_mj_jobs`
  ADD PRIMARY KEY (`id`),
  ADD KEY `message_id` (`message_id`);

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
  ADD PRIMARY KEY (`id`);

--
-- 表的索引 `chatgpt_user_login_logs`
--
ALTER TABLE `chatgpt_user_login_logs`
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
  MODIFY `id` bigint NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=300;

--
-- 使用表AUTO_INCREMENT `chatgpt_chat_items`
--
ALTER TABLE `chatgpt_chat_items`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=98;

--
-- 使用表AUTO_INCREMENT `chatgpt_chat_models`
--
ALTER TABLE `chatgpt_chat_models`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=7;

--
-- 使用表AUTO_INCREMENT `chatgpt_chat_roles`
--
ALTER TABLE `chatgpt_chat_roles`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=125;

--
-- 使用表AUTO_INCREMENT `chatgpt_configs`
--
ALTER TABLE `chatgpt_configs`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;

--
-- 使用表AUTO_INCREMENT `chatgpt_mj_jobs`
--
ALTER TABLE `chatgpt_mj_jobs`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=59;

--
-- 使用表AUTO_INCREMENT `chatgpt_rewards`
--
ALTER TABLE `chatgpt_rewards`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `chatgpt_sd_jobs`
--
ALTER TABLE `chatgpt_sd_jobs`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=69;

--
-- 使用表AUTO_INCREMENT `chatgpt_users`
--
ALTER TABLE `chatgpt_users`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=18;

--
-- 使用表AUTO_INCREMENT `chatgpt_user_login_logs`
--
ALTER TABLE `chatgpt_user_login_logs`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
