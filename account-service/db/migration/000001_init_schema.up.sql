CREATE TABLE "accounts" (
    "id" BIGSERIAL PRIMARY KEY,
    "account_id" varchar UNIQUE NOT NULL,
    "user_id" varchar NOT NULL,
    "balance" int NOT NULL DEFAULT 0,
    "currency" varchar NOT NULL DEFAULT 'EUR',
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "transactions" (
    "id" BIGSERIAL PRIMARY KEY,
    "transaction_id" VARCHAR NOT NULL,
    "from_account_id" varchar NOT NULL,
    "to_account_id" varchar NOT NULL,
    "transaction_amount" int NOT NULL,
    "description" varchar,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now())
);