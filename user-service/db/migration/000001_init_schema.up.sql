CREATE TABLE "users" (
    "id" BIGSERIAL PRIMARY KEY,
    "user_id" VARCHAR NOT NULL,
    "firstname" varchar NOT NULL,
    "lastname" varchar NOT NULL,
    "password" varchar NOT NULL,
    "email" varchar UNIQUE NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now())
);