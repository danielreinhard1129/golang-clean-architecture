-- Create "verification_codes" table
CREATE TABLE "public"."verification_codes" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "code" character varying(6) NOT NULL,
  "purpose" character varying(50) NOT NULL,
  "expired_at" timestamptz NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_users_verification_codes" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Drop "tokens" table
DROP TABLE "public"."tokens";
