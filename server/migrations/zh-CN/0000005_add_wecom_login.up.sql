-- 插入企业微信配置项到系统配置表
INSERT IGNORE INTO t_sys_config (`key`, `value`, `name`, `description`, `create_time`, `update_time`) VALUES
('wecom.enabled', 'false', '是否启用企业微信登录', '企业微信登录开关', UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000),
('wecom.corp_id', '', '企业微信企业ID', '企业微信CorpID', UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000),
('wecom.secret', '', '企业微信应用Secret', '企业微信Secret', UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000),
('wecom.agent_id', '', '企业微信应用AgentId', '企业微信AgentId', UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000),
('wecom.redirect_uri', '', '企业微信回调地址', '企业微信回调URL', UNIX_TIMESTAMP(NOW()) * 1000, UNIX_TIMESTAMP(NOW()) * 1000);
