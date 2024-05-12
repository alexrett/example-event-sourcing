create table if not exists users
(
    id       uuid
        constraint users_pk
            primary key,
    username varchar(100),
    password varchar(255),
    email    varchar(100)
);