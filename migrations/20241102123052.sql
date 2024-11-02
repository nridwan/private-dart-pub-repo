-- Create "user_otps" table
CREATE TABLE "user_otps" (
  "id" uuid NOT NULL,
  "purpose" text NOT NULL,
  "otp" text NOT NULL,
  "expired_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_users_user_otp" FOREIGN KEY ("id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
