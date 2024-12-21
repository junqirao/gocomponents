create table if not exists c_user
(
    id            varchar(50)            not null primary key,
    username      varchar(20)            not null,
    password      varchar(200)           not null,
    created_at    datetime               null,
    updated_at    datetime               null,
    administrator tinyint(1)  default 0  null,
    source        varchar(20) default '' null,
    status        tinyint     default 0  null,
    extra         json                   null,
    unique index c_user_uk_username (username)
);