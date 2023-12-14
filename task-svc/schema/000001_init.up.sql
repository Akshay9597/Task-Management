CREATE TABLE tasks (
    "id" serial not null unique,
    "title" varchar(255) not null,
    "creation_time" timestamp not null default now(),
    "user_id" int not null
);