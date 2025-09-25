-- Modify "users" table
ALTER TABLE "public"."users" ADD COLUMN "verified_at" timestamptz NULL;
-- Create "tokens" table
CREATE TABLE "public"."tokens" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "token" text NOT NULL,
  "expired_at" timestamptz NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_users_tokens" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_tokens_token" to table: "tokens"
CREATE UNIQUE INDEX "idx_tokens_token" ON "public"."tokens" ("token");
