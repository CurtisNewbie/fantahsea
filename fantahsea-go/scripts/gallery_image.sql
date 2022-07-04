CREATE TABLE IF NOT EXISTS gallery (
    id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    gallery_no VARCHAR(32) NOT NULL DEFAULT '' comment 'Gallery No', 
    name VARCHAR(255) NOT NULL COMMENT "Gallery Name",
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT 'when the record is created',
    create_by VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'who created this record',
    update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'when the record is updated',
    update_by VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'who updated this record',
    is_del TINYINT NOT NULL DEFAULT '0' COMMENT '0-normal, 1-deleted'
) ENGINE=InnoDB COMMENT 'Gallery';

CREATE TABLE IF NOT EXISTS gallery_image (
    id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    gallery_no VARCHAR(32) NOT NULL DEFAULT '' comment 'Gallery No', 
    name VARCHAR(255) NOT NULL COMMENT "name of the file",
    file_key VARCHAR(255) NOT NULL COMMENT "file key",
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT 'when the record is created',
    create_by VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'who created this record',
    update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'when the record is updated',
    update_by VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'who updated this record',
    is_del TINYINT NOT NULL DEFAULT '0' COMMENT '0-normal, 1-deleted'
) ENGINE=InnoDB COMMENT "Gallery's Image";
