-- 在t_topic表添加need_a_hand字段（1字节数值）
SET @col_exists = (SELECT COUNT(*) FROM information_schema.columns WHERE table_schema = DATABASE() AND table_name = 't_topic' AND column_name = 'need_a_hand');
SET @sql = IF(@col_exists = 0, 'ALTER TABLE `t_topic` ADD COLUMN `need_a_hand` TINYINT(1) NOT NULL DEFAULT 0 COMMENT ''是否需要帮助：0-不需要，1-需要''', 'SELECT ''Column already exists''');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 在t_comment表添加valuable字段（16字节字符）
SET @col_exists = (SELECT COUNT(*) FROM information_schema.columns WHERE table_schema = DATABASE() AND table_name = 't_comment' AND column_name = 'valuable');
SET @sql = IF(@col_exists = 0, 'ALTER TABLE `t_comment` ADD COLUMN `valuable` VARCHAR(16) NOT NULL DEFAULT '' COMMENT ''价值标识，空字符串表示无价值''', 'SELECT ''Column already exists''');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 添加索引
SET @idx_exists = (SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 't_topic' AND index_name = 'idx_topic_need_a_hand');
SET @sql = IF(@idx_exists = 0, 'ALTER TABLE `t_topic` ADD INDEX `idx_topic_need_a_hand` (`need_a_hand`)', 'SELECT ''Index already exists''');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @idx_exists = (SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 't_comment' AND index_name = 'idx_comment_valuable');
SET @sql = IF(@idx_exists = 0, 'ALTER TABLE `t_comment` ADD INDEX `idx_comment_valuable` (`valuable`)', 'SELECT ''Index already exists''');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;