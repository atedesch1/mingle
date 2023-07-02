CREATE TABLE messages (
    id bigserial PRIMARY KEY,
    user_id bigint REFERENCES users (id) ON DELETE CASCADE,
    content varchar NOT NULL,
    created_at timestamp
);
