-- Add user role
INSERT INTO `t_role` (`id`, `type`, `name`, `code`, `sort_no`, `remark`, `status`, `create_time`, `update_time`) VALUES
(3, 0, 'User', 'user', 2, 'Regular user who can post topics and comments', 0, (UNIX_TIMESTAMP(now()) * 1000), (UNIX_TIMESTAMP(now()) * 1000));