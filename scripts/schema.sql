CREATE TABLE `gallery` (
  `id` int NOT NULL AUTO_INCREMENT,
  `gallery_no` varchar(32) NOT NULL DEFAULT '' COMMENT 'gallery no',
  `user_no` varchar(64) NOT NULL DEFAULT '' COMMENT 'user''s no',
  `name` varchar(255) NOT NULL COMMENT 'gallery name',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'when the record is created',
  `create_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'who created this record',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'when the record is updated',
  `update_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'who updated this record',
  `is_del` tinyint NOT NULL DEFAULT '0' COMMENT '0-normal, 1-deleted',
  PRIMARY KEY (`id`),
  UNIQUE KEY `gallery_no_uniq` (`gallery_no`),
  UNIQUE KEY `name_uk` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Gallery';

CREATE TABLE `gallery_image` (
  `id` int NOT NULL AUTO_INCREMENT,
  `gallery_no` varchar(32) NOT NULL DEFAULT '' COMMENT 'gallery no',
  `image_no` varchar(32) NOT NULL DEFAULT '' COMMENT 'image no',
  `name` varchar(255) NOT NULL COMMENT 'name of the file',
  `file_key` varchar(255) NOT NULL COMMENT 'file key (file-service)',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'when the record is created',
  `create_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'who created this record',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'when the record is updated',
  `update_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'who updated this record',
  `is_del` tinyint NOT NULL DEFAULT '0' COMMENT '0-normal, 1-deleted',
  `status` varchar(20) NOT NULL DEFAULT 'NORMAL' COMMENT 'status',
  PRIMARY KEY (`id`),
  UNIQUE KEY `image_no_uniq` (`image_no`),
  UNIQUE KEY `gallery_no_file_key_uk` (`gallery_no`,`file_key`),
  KEY `gallery_no_idx` (`gallery_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT="Gallery''s Image";

CREATE TABLE `gallery_user_access` (
  `id` int NOT NULL AUTO_INCREMENT,
  `gallery_no` varchar(32) NOT NULL DEFAULT '' COMMENT 'gallery no',
  `user_no` varchar(64) NOT NULL DEFAULT '' COMMENT 'user''s no',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'when the record is created',
  `create_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'who created this record',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'when the record is updated',
  `update_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'who updated this record',
  `is_del` tinyint NOT NULL DEFAULT '0' COMMENT '0-normal, 1-deleted',
  PRIMARY KEY (`id`),
  UNIQUE KEY `gallery_user` (`gallery_no`,`user_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='User access to gallery';

