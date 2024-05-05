create table users (
	id varchar(100) primary key, 
	full_name varchar(100) not null,
	birth_date date not null,
	email varchar(30) unique not null,
	password varchar(100) not null,
	photo_link varchar(100),
	verification_code varchar(40),
	is_verified boolean default false,
	created_at timestamptz default current_timestamp,
	updated_at timestamptz default current_timestamp
);