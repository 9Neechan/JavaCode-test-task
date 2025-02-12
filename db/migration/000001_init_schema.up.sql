CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
SELECT uuid_generate_v4();

CREATE TABLE "wallet" (
  "wallet_uuid" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "balance" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);
