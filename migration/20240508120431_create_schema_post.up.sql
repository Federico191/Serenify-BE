create table posts (
    ID varchar(100) primary key,
    user_id varchar(100) not null,
    content varchar(255) not null,
    photo_link varchar(255),
    created_at timestamptz default current_timestamp,
    updated_at timestamptz default current_timestamp
);


create table post_likes (
    user_id varchar(100) not null,
    post_id varchar(100) not null,
    primary key (user_id, post_id)
);


create table comments (
    ID SERIAL primary key,
    user_id varchar(100) not null,
    post_id varchar(100) not null,
    comment varchar(255) not null,
    created_at timestamptz default current_timestamp
);

create table comment_likes (
    user_id varchar(100) not null,
    comment_id int not null,
    primary key (user_id, comment_id)
);

alter table Posts add foreign key (user_id) references users(id);
alter table post_likes add foreign key (user_id) references users(id);
alter table post_likes add foreign key (post_id) references posts(ID);
alter table Comments add foreign key (user_id) references users(id);
alter table Comments add foreign key (post_id) references posts(ID);
alter table comment_likes add foreign key (user_id) references users(id);
alter table comment_likes add foreign key (comment_id) references Comments(ID);
