-- 增加字段保存用户开通的 AI 模型
ALTER TABLE `chatgpt_users` ADD `chat_models_json` TEXT NOT NULL COMMENT 'AI模型 json' AFTER `chat_roles_json`;
UPDATE `chatgpt_users` SET chat_models_json = '["completions_pro","eb-instant","general","generalv2","chatglm_pro","chatglm_lite","chatglm_std","gpt-3.5-turbo-16k"]';
-- 为每个模型设置对话权重
ALTER TABLE `chatgpt_chat_models` ADD `weight` TINYINT(3) NOT NULL COMMENT '对话权重，每次对话扣减多少次对话额度' AFTER `enabled`;
UPDATE `chatgpt_chat_models` SET weight = 1;

-- 更新系统配置，支持文心4.0模型
UPDATE `chatgpt_configs` SET config_json = '{"azure":{"api_url":"https://chat-bot-api.openai.azure.com/openai/deployments/{model}/chat/completions?api-version=2023-05-15","max_tokens":1024,"temperature":1},"baidu":{"api_url":"https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/{model}","max_tokens":1024,"temperature":0.95},"chat_gml":{"api_url":"https://open.bigmodel.cn/api/paas/v3/model-api/{model}/sse-invoke","max_tokens":1024,"temperature":0.95},"context_deep":4,"enable_context":true,"enable_history":true,"open_ai":{"api_url":"https://api.openai.com/v1/chat/completions","max_tokens":1024,"temperature":1},"xun_fei":{"api_url":"wss://spark-api.xf-yun.com/{version}/chat","max_tokens":1024,"temperature":0.5}}' WHERE marker ='chat';
