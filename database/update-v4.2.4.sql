-- MiniMax AI 模型支持
-- 添加 MiniMax 大语言模型预设
-- MiniMax API 兼容 OpenAI 接口，使用前请在 API KEY 管理中添加 MiniMax 的 API KEY
-- API URL: https://api.minimax.io

INSERT INTO `chatgpt_chat_models` (`description`, `category`, `type`, `name`, `value`, `sort_num`, `enabled`, `power`, `temperature`, `max_tokens`, `max_context`, `open`, `key_id`, `options`, `created_at`, `updated_at`) VALUES
('MiniMax 最新旗舰模型，512K 上下文，最大输出 128K，支持图片输入', 'MiniMax', 'chat', 'MiniMax-M3', 'MiniMax-M3', 30, 0, 1, 0.9, 128000, 512000, 1, 0, '', NOW(), NOW()),
('MiniMax 上一代旗舰模型', 'MiniMax', 'chat', 'MiniMax-M2.7', 'MiniMax-M2.7', 31, 0, 1, 0.9, 4096, 192000, 1, 0, '', NOW(), NOW()),
('MiniMax 高速版本，支持 204K 上下文', 'MiniMax', 'chat', 'MiniMax-M2.7-highspeed', 'MiniMax-M2.7-highspeed', 32, 0, 1, 0.9, 4096, 204000, 1, 0, '', NOW(), NOW());
