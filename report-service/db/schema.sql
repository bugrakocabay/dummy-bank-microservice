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