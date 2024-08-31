CREATE TABLE IF NOT EXISTS tasks
(
    id uuid not null PRIMARY KEY,
    status varchar(255) not null default 'FORMING',
    all_faces_quantity int,
    male_quantity int,
    female_quantity int,
    average_male_age float,
    average_female_age float
);