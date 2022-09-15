create table file_info_tab (
    id bigint unsigned not null auto_increment,
    file_name varchar(256) not null,
    hash char(32) not null,
    file_size bigint unsigned not null, 
    create_time bigint unsigned not null,
    down_key varchar(64) not null,
    extra blob not null,
    primary key(id),
    unique key idx_downkey(down_key)
) ENGINE=InnoDB CHARACTER SET utf8mb4;