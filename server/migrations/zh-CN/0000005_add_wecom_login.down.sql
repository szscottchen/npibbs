-- 删除企业微信配置项
DELETE FROM t_sys_config WHERE `key` LIKE 'wecom.%';
