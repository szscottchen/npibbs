-- Remove need_a_hand field from t_topic table
ALTER TABLE `t_topic` DROP COLUMN `need_a_hand`;

-- Remove valuable field from t_comment table
ALTER TABLE `t_comment` DROP COLUMN `valuable`;