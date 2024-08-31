CREATE TABLE IF NOT EXISTS images (
                        id uuid not null PRIMARY KEY,
                        task_id uuid REFERENCES tasks(id) ON DELETE CASCADE,
                        image_name varchar(255),
                        image_address varchar(255)
);