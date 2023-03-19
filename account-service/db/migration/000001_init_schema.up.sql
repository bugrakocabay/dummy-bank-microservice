CREATE TABLE "accounts" (
    "id" BIGSERIAL PRIMARY KEY,
    "account_id" varchar UNIQUE NOT NULL,
    "user_id" varchar NOT NULL,
    "balance" float NOT NULL DEFAULT 0,
    "currency" varchar NOT NULL DEFAULT 'EUR',
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "transactions" (
    "id" BIGSERIAL PRIMARY KEY,
    "transaction_id" VARCHAR UNIQUE NOT NULL,
    "from_account_id" varchar NOT NULL,
    "to_account_id" varchar NOT NULL,
    "transaction_amount" float NOT NULL,
    "commission" float NOT NULL,
    "description" varchar,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "daily_transaction_report" (
    "id" BIGSERIAL PRIMARY KEY,
    "num_transactions" int NOT NULL,
    "avg_transaction_amount" float NOT NULL,
    "total_transaction_amount" int NOT NULL,
    "total_commission" float NOT NULL,
    "day" varchar NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now())
);