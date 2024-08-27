CREATE TABLE IF NOT EXISTS tasks
(
    id uuid not null PRIMARY KEY,
    status varchar(255) not null default 'FORMING'
);