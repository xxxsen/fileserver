CREATE DATABASE IF NOT EXISTS filedb;
use filedb;
create table if not exists fs_system_tab (
    id bigint unsigned not null,
    parent_id bigint unsigned not null,
    name_code int unsigned not null,
    file_name varchar(256) not null,
    file_type tinyint unsigned not null,
    file_size bigint unsigned not null,
    ctime bigint unsigned not null,
    mtime bigint unsigned not null,
    down_key bigint unsigned not null,
    primary key(id),
    key idx_parentid_namecode(parent_id, name_code, file_type)
) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8mb4