create table if not exists test
(
    id    int generated always as identity primary key,
    ptype text not null,
    v0    text not null,
    v1    text not null,
    v2    text not null,
    v3    text not null,
    v4    text not null,
    v5    text not null
);
create unique index if not exists test_uk on test (ptype, v0, v1, v2, v3, v4, v5);

create table if not exists test
(
    id    int generated always as identity primary key,
    ptype text not null,
    v0    text not null,
    v1    text not null,
    v2    text not null,
    v3    text not null,
    v4    text not null,
    v5    text not null,
    v6    text not null,
    v7    text not null,
    v8    text not null,
    v9    text not null,
    v10   text not null,
    v11   text not null,
    v12   text not null,
    v13   text not null,
    v14   text not null,
    v15   text not null,
    v16   text not null,
    v17   text not null,
    v18   text not null,
    v19   text not null
);
create unique index if not exists test_uk on test (ptype, v0, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13,
                                                   v14, v15, v16, v17, v18, v19);
