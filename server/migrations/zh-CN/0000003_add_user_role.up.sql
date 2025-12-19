-- 添加普通用户角色-- 添加普通用户角色
INSERT INTO `t_role` (`id`, `type`, `name`, `code`, `sort_no`, `remark`, `status`, `create_time`, `update_time`) VALUES
(3, 0, '用户', 'user', 2, '普通用户，可以发言发帖', 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000));