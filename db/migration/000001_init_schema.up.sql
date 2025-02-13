CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
SELECT uuid_generate_v4();

CREATE TABLE "wallet" (
  "wallet_uuid" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "balance" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

INSERT INTO "wallet" ("wallet_uuid", "balance", "created_at") 
VALUES 
  ('7ff05ab9-80d5-40d0-8037-7133da806e49', 1000, NOW()), 
  (uuid_generate_v4(), 5000, NOW()), 
  (uuid_generate_v4(), 250, NOW());