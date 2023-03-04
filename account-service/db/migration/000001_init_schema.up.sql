CREATE TABLE "accounts" (
    "id" BIGSERIAL PRIMARY KEY,
    "account_id" varchar UNIQUE NOT NULL,
    "firstname" varchar UNIQUE NOT NULL,
    "lastname" varchar UNIQUE NOT NULL,
    "email" varchar UNIQUE NOT NULL,
    "hashed_password" varchar UNIQUE NOT NULL,
    "type" varchar DEFAULT 'user',
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now())
);
