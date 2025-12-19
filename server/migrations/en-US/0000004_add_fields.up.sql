-- Add need_a_hand field to t_topic table (1 byte numeric)
SET @col_exists = (SELECT COUNT(*) FROM information_schema.columns WHERE table_schema = DATABASE() AND table_name = 't_topic' AND column_name = 'need_a_hand');
SET @sql = IF(@col_exists = 0, 'ALTER TABLE `t_topic` ADD COLUMN `need_a_hand` TINYINT(1) NOT NULL DEFAULT 0 COMMENT ''Whether help is needed: 0-no, 1-yes''', 'SELECT ''Column already exists''');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add valuable field to t_comment table (16 byte character)
SET @col_exists = (SELECT COUNT(*) FROM information_schema.columns WHERE table_schema = DATABASE() AND table_name = 't_comment' AND column_name = 'valuable');
SET @sql = IF(@col_exists = 0, 'ALTER TABLE `t_comment` ADD COLUMN `valuable` VARCHAR(16) NOT NULL DEFAULT '' COMMENT ''Value identifier, empty string means no value''', 'SELECT ''Column already exists''');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add indexes
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