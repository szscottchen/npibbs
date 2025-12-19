-- 删除t_topic表的need_a_hand字段
ALTER TABLE `t_topic` DROP COLUMN `need_a_hand`;

-- 删除t_comment表的valuable字段
ALTER TABLE `t_comment` DROP COLUMN `valuable`;