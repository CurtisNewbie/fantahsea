CREATE TABLE IF NOT EXISTS gallery (
    id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    gallery_no VARCHAR(32) NOT NULL DEFAULT '' COMMENT 'gallery no', 
    user_no VARCHAR(64) NOT NULL DEFAULT '' COMMENT "user's no",
    name VARCHAR(255) NOT NULL COMMENT "gallery name",
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT 'when the record is created',
    create_by VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'who created this record',
    update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'when the record is updated',
    update_by VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'who updated this record',
    is_del TINYINT NOT NULL DEFAULT '0' COMMENT '0-normal, 1-deleted',
    UNIQUE gallery_no_uniq(gallery_no)
) ENGINE=InnoDB COMMENT 'Gallery';

CREATE TABLE IF NOT EXISTS gallery_image (
    id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    gallery_no VARCHAR(32) NOT NULL DEFAULT '' COMMENT 'gallery no', 
    image_no VARCHAR(32) NOT NULL DEFAULT '' COMMENT 'image no',
    name VARCHAR(255) NOT NULL COMMENT "name of the file",
    file_key VARCHAR(255) NOT NULL COMMENT "file key (file-service)",
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT 'when the record is created',
    create_by VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'who created this record',
    update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'when the record is updated',
    update_by VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'who updated this record',
    is_del TINYINT NOT NULL DEFAULT '0' COMMENT '0-normal, 1-deleted',
    UNIQUE image_no_uniq(image_no)
) ENGINE=InnoDB COMMENT "Gallery's Image";

CREATE TABLE IF NOT EXISTS gallery_user_access (
    id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    gallery_no VARCHAR(32) NOT NULL DEFAULT '' COMMENT 'gallery no', 
    user_no VARCHAR(64) NOT NULL DEFAULT '' COMMENT "user's no",
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT 'when the record is created',
    create_by VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'who created this record',
    update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'when the record is updated',
    update_by VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'who updated this record',
    is_del TINYINT NOT NULL DEFAULT '0' COMMENT '0-normal, 1-deleted',
    UNIQUE gallery_user (gallery_no, user_no) 
) ENGINE=InnoDB COMMENT 'User access to gallery';

    
