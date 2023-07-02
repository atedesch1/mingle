CREATE TABLE users (
    id bigserial PRIMARY KEY,
    name varchar NOT NULL,
    created_at timestamp DEFAULT now()
);
