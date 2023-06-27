create table bookmarks
(
    id         bigserial not null,
    url        varchar   not null,
    title      varchar   not null,
    created_at timestamp not null,
    updated_at timestamp,
    primary key (id)
);
