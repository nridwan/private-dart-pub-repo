-- Create "pub_tokens" table
CREATE TABLE "pub_tokens" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "remarks" text NOT NULL,
  "write" boolean NOT NULL DEFAULT false,
  "expired_at" timestamptz NOT NULL,
  "user_id" uuid NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_pub_tokens_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE CASCADE ON DELETE SET NULL
);
-- Create index "idx_pub_tokens_deleted_at" to table: "pub_tokens"
CREATE INDEX "idx_pub_tokens_deleted_at" ON "pub_tokens" ("deleted_at");
-- Create "pub_packages" table
CREATE TABLE "pub_packages" (
  "name" text NOT NULL,
  "private" boolean NOT NULL DEFAULT true,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("name")
);
-- Create index "idx_pub_packages_deleted_at" to table: "pub_packages"
CREATE INDEX "idx_pub_packages_deleted_at" ON "pub_packages" ("deleted_at");
-- Create "pub_versions" table
CREATE TABLE "pub_versions" (
  "package_name" text NOT NULL,
  "version" text NOT NULL,
  "version_number_major" bigint NOT NULL,
  "version_number_minor" bigint NOT NULL,
  "version_number_patch" bigint NOT NULL,
  "prerelease" boolean NOT NULL DEFAULT false,
  "pubspec" jsonb NOT NULL DEFAULT '{}',
  "uploader_id" uuid NULL,
  "readme" text NULL,
  "changelog" text NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  CONSTRAINT "fk_pub_packages_versions" FOREIGN KEY ("package_name") REFERENCES "pub_packages" ("name") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_pub_versions_uploader" FOREIGN KEY ("uploader_id") REFERENCES "users" ("id") ON UPDATE CASCADE ON DELETE SET NULL
);
-- Create index "idx_pub_versions_deleted_at" to table: "pub_versions"
CREATE INDEX "idx_pub_versions_deleted_at" ON "pub_versions" ("deleted_at");
-- Create index "idx_pub_versions_pubversion" to table: "pub_versions"
CREATE UNIQUE INDEX "idx_pub_versions_pubversion" ON "pub_versions" ("package_name", "version");
