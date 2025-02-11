CREATE TABLE "vallet" (
  "vallet_id" bigserial PRIMARY KEY,
  "balance" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);
