-- 添加批量导入菜单
INSERT INTO `t_menu` (
    `parent_id`,
    `type`,
    `name`,
    `title`,
    `icon`,
    `path`,
    `component`,
    `sort_no`,
    `status`,
    `create_time`,
    `update_time`
) VALUES (
    2,
    'menu',
    'UserBatchImport',
    '批量导入',
    'icon-upload',
    '/user/batch-import',
    'user/batch-import',
    3,
    0,
    (UNIX_TIMESTAMP(now()) * 1000),
    (UNIX_TIMESTAMP(now()) * 1000)
);

-- 为Owner角色授权
INSERT INTO `t_role_menu` (
    `role_id`,
    `menu_id`,
    `create_time`
) VALUES (
    1,
    (SELECT id FROM t_menu WHERE name = 'UserBatchImport' LIMIT 1),
    (UNIX_TIMESTAMP(now()) * 1000)
);
