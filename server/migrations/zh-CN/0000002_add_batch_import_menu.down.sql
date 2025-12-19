-- 删除Owner角色的批量导入菜单权限
DELETE FROM t_role_menu WHERE menu_id = (SELECT id FROM t_menu WHERE name = 'UserBatchImport' LIMIT 1);

-- 删除批量导入菜单
DELETE FROM t_menu WHERE name = 'UserBatchImport';
