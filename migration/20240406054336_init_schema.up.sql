create table users (
	id varchar(100) primary key, 
	full_name varchar(100) not null,
	birth_date date not null,
	email varchar(30) unique not null,
	password varchar(100) not null,
	photo_link varchar(100),
	verification_code varchar(40),
	is_verified boolean default false,
	score_test int default 0,
	created_at timestamptz default current_timestamp,
	updated_at timestamptz default current_timestamp
);

create table reset_password_tokens (
	token varchar(100) not null,
	user_id varchar(100) not null,
	created_at timestamptz default current_timestamp,
	primary key (token, user_id),
	foreign key (user_id) references users(id)
);