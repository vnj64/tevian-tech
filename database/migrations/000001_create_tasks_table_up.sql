CREATE TABLE IF NOT EXISTS tasks
(
    id uuid not null PRIMARY KEY,
    status varchar(255) not null default 'FORMING',
    image_address varchar(255),
    image_name varchar(255)
);