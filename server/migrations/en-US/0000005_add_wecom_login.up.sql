-- Insert WeChat Work configuration items into system config table
INSERT IGNORE INTO t_sys_config (`key`, `value`, `name`, `description`, `create_time`, `update_time`) VALUES
('wecom.enabled', 'false', 'Enable WeChat Work login', 'WeChat Work login switch', UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000),
('wecom.corp_id', '', 'WeChat Work Corp ID', 'WeChat Work CorpID', UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000),
('wecom.secret', '', 'WeChat Work App Secret', 'WeChat Work Secret', UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000),
('wecom.agent_id', '', 'WeChat Work App AgentId', 'WeChat Work AgentId', UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000),
('wecom.redirect_uri', '', 'WeChat Work callback URL', 'WeChat Work callback URL', UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000);
