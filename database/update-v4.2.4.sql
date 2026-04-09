-- MiniMax AI 模型支持
-- 添加 MiniMax 大语言模型预设
-- MiniMax API 兼容 OpenAI 接口，使用前请在 API KEY 管理中添加 MiniMax 的 API KEY
-- API URL: https://api.minimax.io

INSERT INTO `chatgpt_chat_models` (`description`, `category`, `type`, `name`, `value`, `sort_num`, `enabled`, `power`, `temperature`, `max_tokens`, `max_context`, `open`, `key_id`, `options`, `created_at`, `updated_at`) VALUES
('MiniMax 最新旗舰模型，支持 1M 上下文', 'MiniMax', 'chat', 'MiniMax-M2.7', 'MiniMax-M2.7', 30, 0, 1, 0.9, 4096, 1000000, 1, 0, '', NOW(), NOW()),
('MiniMax 高速版本，支持 204K 上下文', 'MiniMax', 'chat', 'MiniMax-M2.7-highspeed', 'MiniMax-M2.7-highspeed', 31, 0, 1, 0.9, 4096, 204000, 1, 0, '', NOW(), NOW()),
('MiniMax M2.5 模型', 'MiniMax', 'chat', 'MiniMax-M2.5', 'MiniMax-M2.5', 32, 0, 1, 0.9, 4096, 1000000, 1, 0, '', NOW(), NOW()),
('MiniMax M2.5 高速版本，支持 204K 上下文', 'MiniMax', 'chat', 'MiniMax-M2.5-highspeed', 'MiniMax-M2.5-highspeed', 33, 0, 1, 0.9, 4096, 204000, 1, 0, '', NOW(), NOW());
