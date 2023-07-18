ALTER TABLE gallery_image MODIFY COLUMN file_key VARCHAR(64) NOT NULL COMMENT 'file key (vfm)';
ALTER TABLE gallery_image ADD KEY idx_file_key (file_key);
ALTER TABLE gallery ADD COLUMN dir_file_key VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'directory file_key (vfm)';
ALTER TABLE gallery ADD KEY idx_dir_file_key (dir_file_key);
