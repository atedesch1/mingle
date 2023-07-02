CREATE TABLE "users" (
    "id" bigserial PRIMARY KEY,
    "name" varchar NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "messages" (
    "id" bigserial PRIMARY KEY,
    "user_id" bigserial NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE,
    "content" varchar NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);
