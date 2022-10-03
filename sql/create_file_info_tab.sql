CREATE DATABASE IF NOT EXISTS filedb;
use filedb;
create table if not exists file_info_tab (
    id bigint unsigned not null auto_increment,
    file_name varchar(256) not null,
    hash varchar(40) not null,
    file_size bigint unsigned not null, 
    create_time bigint unsigned not null,
    down_key bigint unsigned not null,
    file_key varchar(128) not null,
    extra blob not null,
    st_type tinyint not null,
    primary key(id),
    unique key idx_downkey(down_key)
) ENGINE=InnoDB CHARACTER SET utf8mb4;