create table articles (
    id varchar(100) primary key,
    title varchar(100) not null,
    content text not null,
    photo_link varchar(255),
    created_at timestamptz default current_timestamp
);
