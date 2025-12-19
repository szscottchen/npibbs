-- Remove Owner role's batch import menu permission
DELETE FROM t_role_menu WHERE menu_id = (SELECT id FROM t_menu WHERE name = 'UserBatchImport' LIMIT 1);

-- Remove batch import menu
DELETE FROM t_menu WHERE name = 'UserBatchImport';
