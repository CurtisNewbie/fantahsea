alter table gallery_image add constraint gallery_no_file_key_uk unique (gallery_no, file_key);