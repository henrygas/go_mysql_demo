CREATE DATABASE IF NOT EXISTS `ys` CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `ys`.`square_num` (
    num INT(11) NOT NULL DEFAULT 0 COMMENT 'number',
    square_num INT(11) NOT NULL DEFAULT 0 COMMENT 'square of num',
    PRIMARY KEY (`num`)
) ENGINE=InnoDB;