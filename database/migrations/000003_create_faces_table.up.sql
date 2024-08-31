CREATE TABLE IF NOT EXISTS faces  (
                       id uuid not null PRIMARY KEY,
                       image_id uuid REFERENCES images(id) ON DELETE CASCADE,
                       bbox VARCHAR(255),
                       gender VARCHAR(50),
                       age INT
);