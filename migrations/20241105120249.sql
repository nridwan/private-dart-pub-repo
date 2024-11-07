-- Modify "users" table
ALTER TABLE "users" ADD COLUMN "can_write" boolean NOT NULL DEFAULT false;
