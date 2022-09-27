alter table gallery add unique user_name_uk (user_no, name);
alter table gallery_image add column status varchar(20) NOT NULL default 'NORMAL' COMMENT 'status';