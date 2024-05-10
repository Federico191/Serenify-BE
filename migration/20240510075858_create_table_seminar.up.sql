create table seminars (
    id varchar(100) primary key,
    title varchar(100) not null,
    time varchar(20) not null,
    place varchar(50) not null,
    price int not null,
    description varchar(255) not null,
    photo_link varchar(255),
    created_at timestamp default current_timestamp
);