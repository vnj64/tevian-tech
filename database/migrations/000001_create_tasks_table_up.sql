CREATE TABLE IF NOT EXISTS tasks
(
    id uuid not null PRIMARY KEY,
    status varchar(255) not null default 'FORMING',
    image_address varchar(255),
    image_name varchar(255),
    all_faces_quantity int,
    male_quantity int,
    female_quantity int,
    average_male_age float,
    average_female_age float
);

CREATE TABLE images (
    id uuid not null PRIMARY KEY,
    task_id uuid REFERENCES tasks(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE faces (
    id uuid not null PRIMARY KEY,
    image_id uuid REFERENCES images(id) ON DELETE CASCADE,
    bbox VARCHAR(255),
    gender VARCHAR(50),
    age INT
);
