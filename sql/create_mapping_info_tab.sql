create table mapping_info_tab (
    id bigint unsigned not null auto_increment,
    file_name varchar(256) not null,
    hash_code int unsigned not null,
    check_sum varchar(40) not null, 
    create_time bigint unsigned not null,
    modify_time bigint unsigned not null,
    file_id bigint unsigned not null,
    primary key(id),
    unique key idx_code_ck(hash_code, check_sum)
) ENGINE=InnoDB CHARACTER SET utf8mb4;