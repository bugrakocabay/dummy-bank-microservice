CREATE TABLE "accounts" (
    "id" BIGSERIAL PRIMARY KEY,
    "account_id" varchar UNIQUE NOT NULL,
    "firstname" varchar UNIQUE NOT NULL,
    "lastname" varchar UNIQUE NOT NULL,
    "balance" int NOT NULL DEFAULT 0,
    "email" varchar UNIQUE NOT NULL,
    "password" varchar UNIQUE NOT NULL,
    "type" varchar NOT NULL DEFAULT 'user',
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now())
);