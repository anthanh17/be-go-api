CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "user_name" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "mail" varchar NOT NULL,
  "role" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);
