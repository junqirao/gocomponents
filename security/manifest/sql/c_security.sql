create table if not exists test.c_security
(
    id      int auto_increment primary key,
    type    varchar(20) not null,
    name    varchar(20) null,
    content mediumtext  null
);

