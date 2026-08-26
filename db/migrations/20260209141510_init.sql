-- +goose Up
-- +goose StatementBegin

CREATE SCHEMA "sso";

CREATE TABLE "sso"."users" (
    "id" bigserial PRIMARY KEY,
    "username" varchar(255) UNIQUE NOT NULL,
    "password_hash" varchar(255) NOT NULL,
    "full_name" varchar(255) NOT NULL,
    "date_of_birth" date,
    "gender" smallint,
    "avatar_file_key" varchar(255),
    "deleted" bool NOT NULL DEFAULT false,
    "created_at" timestamp NOT NULL DEFAULT (now()),
    "updated_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "sso"."clients" (
    "id" bigserial PRIMARY KEY,
    "name" varchar(255) NOT NULL,
    "code" varchar(255) UNIQUE NOT NULL,
    "secret_key" varchar(255) NOT NULL,
    "deleted" bool NOT NULL DEFAULT false,
    "created_at" timestamp NOT NULL DEFAULT (now()),
    "updated_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "sso"."client_audiences" (
    "id" bigserial PRIMARY KEY,
    "client_id" bigint NOT NULL,
    "uri" varchar(255) NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT (now()),
    "updated_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "sso"."client_redirect_uris" (
    "id" bigserial PRIMARY KEY,
    "client_id" bigint NOT NULL,
    "uri" varchar(255) NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT (now()),
    "updated_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "sso"."user_client_links" (
    "id" bigserial PRIMARY KEY,
    "user_id" bigint NOT NULL,
    "client_id" bigint NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT (now()),
    "updated_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "sso"."scopes" (
    "id" bigserial PRIMARY KEY,
    "code" varchar(255) UNIQUE NOT NULL,
    "name" varchar(255) NOT NULL,
    "description" varchar(255) NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT (now()),
    "updated_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "sso"."client_scope_links" (
    "id" bigserial PRIMARY KEY,
    "client_id" bigint NOT NULL,
    "scope_id" bigint NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT (now()),
    "updated_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "sso"."roles" (
    "id" bigserial PRIMARY KEY,
    "code" varchar(255) NOT NULL,
    "name" varchar(255) NOT NULL,
    "description" varchar(255),
    "active" bool NOT NULL DEFAULT true,
    "deleted" bool NOT NULL DEFAULT false,
    "created_at" timestamp NOT NULL DEFAULT (now()),
    "updated_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "sso"."user_role_links" (
    "id" bigserial PRIMARY KEY,
    "user_id" bigint NOT NULL,
    "role_id" bigint NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT (now()),
    "updated_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "sso"."permissions" (
    "id" bigserial PRIMARY KEY,
    "code" varchar(255) NOT NULL,
    "description" varchar(255),
    "active" bool NOT NULL DEFAULT true,
    "deleted" bool NOT NULL DEFAULT false,
    "created_at" timestamp NOT NULL DEFAULT (now()),
    "updated_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "sso"."role_permission_links" (
    "id" bigserial PRIMARY KEY,
    "role_id" bigint NOT NULL,
    "permission_id" bigint NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT (now()),
    "updated_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "sso"."client_default_role_links" (
    "id" bigserial PRIMARY KEY,
    "client_id" bigint NOT NULL,
    "role_id" bigint NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT (now()),
    "updated_at" timestamp NOT NULL DEFAULT (now())
);


CREATE INDEX "idx_users_username_deleted" ON "sso"."users" ("username", "deleted");

CREATE INDEX "idx_users_code_deleted" ON "sso"."clients" ("code", "deleted");

CREATE INDEX "idx_client_audiences_client_id" ON "sso"."client_audiences" ("client_id");

CREATE UNIQUE INDEX "uidx_client_audiences_client_id_uri" ON "sso"."client_audiences" ("client_id", "uri");

CREATE INDEX "idx_client_redirect_uris_client_id" ON "sso"."client_redirect_uris" ("client_id");

CREATE UNIQUE INDEX "uidx_client_redirect_uris_client_id_uri" ON "sso"."client_redirect_uris" ("client_id", "uri");

CREATE INDEX "idx_user_client_links_user_id" ON "sso"."user_client_links" ("user_id");

CREATE INDEX "idx_user_client_links_client_id" ON "sso"."user_client_links" ("client_id");

CREATE UNIQUE INDEX "uidx_user_client_links_user_id_client_id" ON "sso"."user_client_links" ("user_id", "client_id");

CREATE INDEX "idx_scopes_code" ON "sso"."scopes" ("code");

CREATE INDEX "idx_client_scope_links_client_id" ON "sso"."client_scope_links" ("client_id");

CREATE INDEX "idx_client_scope_links_scope_id" ON "sso"."client_scope_links" ("scope_id");

CREATE UNIQUE INDEX "uidx_client_scope_links_client_id_scope_id" ON "sso"."client_scope_links" ("client_id", "scope_id");

CREATE INDEX "idx_roles_code_active" ON "sso"."roles" ("code", "active");

CREATE INDEX "idx_user_role_links_user_id" ON "sso"."user_role_links" ("user_id");

CREATE INDEX "idx_user_role_links_role_id" ON "sso"."user_role_links" ("role_id");

CREATE UNIQUE INDEX "uidx_user_role_links_user_id_role_id" ON "sso"."user_role_links" ("user_id", "role_id");

CREATE INDEX "idx_permissions_code_active" ON "sso"."permissions" ("code", "active");

CREATE INDEX "idx_role_permission_links_role_id" ON "sso"."role_permission_links" ("role_id");

CREATE INDEX "idx_role_permission_links_permission_id" ON "sso"."role_permission_links" ("permission_id");

CREATE UNIQUE INDEX "uidx_role_permission_links_role_id_permission_id" ON "sso"."role_permission_links" ("role_id", "permission_id");

CREATE INDEX "idx_client_default_role_links_client_id" ON "sso"."client_default_role_links" ("client_id");

CREATE INDEX "idx_client_default_role_links_role_id" ON "sso"."client_default_role_links" ("role_id");

CREATE UNIQUE INDEX "uidx_client_default_role_links_client_id_role_id" ON "sso"."client_default_role_links" ("client_id", "role_id");


ALTER TABLE "sso"."client_audiences" ADD CONSTRAINT "fk_client_audiences_client_id" FOREIGN KEY ("client_id") REFERENCES "sso"."clients" ("id");

ALTER TABLE "sso"."client_redirect_uris" ADD CONSTRAINT "fk_client_redirect_uris_client_id" FOREIGN KEY ("client_id") REFERENCES "sso"."clients" ("id");

ALTER TABLE "sso"."user_client_links" ADD CONSTRAINT "fk_user_client_links_user_id" FOREIGN KEY ("user_id") REFERENCES "sso"."users" ("id");

ALTER TABLE "sso"."user_client_links" ADD CONSTRAINT "fk_user_client_links_client_id" FOREIGN KEY ("client_id") REFERENCES "sso"."clients" ("id");

ALTER TABLE "sso"."client_scope_links" ADD CONSTRAINT "fk_client_scope_links_client_id" FOREIGN KEY ("client_id") REFERENCES "sso"."clients" ("id");

ALTER TABLE "sso"."client_scope_links" ADD CONSTRAINT "fk_client_scope_links_scope_id" FOREIGN KEY ("scope_id") REFERENCES "sso"."scopes" ("id");

ALTER TABLE "sso"."user_role_links" ADD CONSTRAINT "fk_user_role_links_user_id" FOREIGN KEY ("user_id") REFERENCES "sso"."users" ("id");

ALTER TABLE "sso"."user_role_links" ADD CONSTRAINT "fk_user_role_links_role_id" FOREIGN KEY ("role_id") REFERENCES "sso"."roles" ("id");

ALTER TABLE "sso"."role_permission_links" ADD CONSTRAINT "fk_role_permission_links_role_id" FOREIGN KEY ("role_id") REFERENCES "sso"."roles" ("id");

ALTER TABLE "sso"."role_permission_links" ADD CONSTRAINT "fk_role_permission_links_permission_id" FOREIGN KEY ("permission_id") REFERENCES "sso"."permissions" ("id");

ALTER TABLE "sso"."client_default_role_links" ADD CONSTRAINT "fk_client_default_role_links_client_id" FOREIGN KEY ("client_id") REFERENCES "sso"."clients" ("id");

ALTER TABLE "sso"."client_default_role_links" ADD CONSTRAINT "fk_client_default_role_links_role_id" FOREIGN KEY ("role_id") REFERENCES "sso"."roles" ("id");

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE "sso"."client_audiences" DROP CONSTRAINT IF EXISTS "fk_client_audiences_client_id";

ALTER TABLE "sso"."client_redirect_uris" DROP CONSTRAINT IF EXISTS "fk_client_redirect_uris_client_id";

ALTER TABLE "sso"."user_client_links" DROP CONSTRAINT IF EXISTS "fk_user_client_links_user_id";

ALTER TABLE "sso"."user_client_links" DROP CONSTRAINT IF EXISTS "fk_user_client_links_client_id";

ALTER TABLE "sso"."client_scope_links" DROP CONSTRAINT IF EXISTS "fk_client_scope_links_client_id";

ALTER TABLE "sso"."client_scope_links" DROP CONSTRAINT IF EXISTS "fk_client_scope_links_scope_id";

ALTER TABLE "sso"."user_role_links" DROP CONSTRAINT IF EXISTS "fk_user_role_links_user_id";

ALTER TABLE "sso"."user_role_links" DROP CONSTRAINT IF EXISTS "fk_user_role_links_role_id";

ALTER TABLE "sso"."role_permission_links" DROP CONSTRAINT IF EXISTS "fk_role_permission_links_role_id";

ALTER TABLE "sso"."role_permission_links" DROP CONSTRAINT IF EXISTS "fk_role_permission_links_permission_id";

ALTER TABLE "sso"."client_default_role_links" DROP CONSTRAINT IF EXISTS "fk_client_default_role_links_client_id";

ALTER TABLE "sso"."client_default_role_links" DROP CONSTRAINT IF EXISTS "fk_client_default_role_links_role_id";


DROP INDEX IF EXISTS "idx_users_username_deleted";

DROP INDEX IF EXISTS "idx_users_code_deleted";

DROP INDEX IF EXISTS "idx_client_audiences_client_id";

DROP INDEX IF EXISTS "uidx_client_audiences_client_id_uri";

DROP INDEX IF EXISTS "idx_client_redirect_uris_client_id";

DROP INDEX IF EXISTS "uidx_client_redirect_uris_client_id_uri";

DROP INDEX IF EXISTS "idx_user_client_links_user_id";

DROP INDEX IF EXISTS "idx_user_client_links_client_id";

DROP INDEX IF EXISTS "uidx_user_client_links_user_id_client_id";

DROP INDEX IF EXISTS "idx_scopes_code";

DROP INDEX IF EXISTS "idx_client_scope_links_client_id";

DROP INDEX IF EXISTS "idx_client_scope_links_scope_id";

DROP INDEX IF EXISTS "uidx_client_scope_links_client_id_scope_id";

DROP INDEX IF EXISTS "idx_roles_code_active";

DROP INDEX IF EXISTS "idx_user_role_links_user_id";

DROP INDEX IF EXISTS "idx_user_role_links_role_id";

DROP INDEX IF EXISTS "uidx_user_role_links_user_id_role_id";

DROP INDEX IF EXISTS "idx_permissions_code_active";

DROP INDEX IF EXISTS "idx_role_permission_links_role_id";

DROP INDEX IF EXISTS "idx_role_permission_links_permission_id";

DROP INDEX IF EXISTS "uidx_role_permission_links_role_id_permission_id";

DROP INDEX IF EXISTS "idx_client_default_role_links_client_id";

DROP INDEX IF EXISTS "idx_client_default_role_links_role_id";

DROP INDEX IF EXISTS "uidx_client_default_role_links_client_id_role_id";


DROP TABLE IF EXISTS "sso"."client_default_role_links";

DROP TABLE IF EXISTS "sso"."role_permission_links";

DROP TABLE IF EXISTS "sso"."permissions";

DROP TABLE IF EXISTS "sso"."user_role_links";

DROP TABLE IF EXISTS "sso"."roles";

DROP TABLE IF EXISTS "sso"."client_scope_links";

DROP TABLE IF EXISTS "sso"."scopes";

DROP TABLE IF EXISTS "sso"."user_client_links";

DROP TABLE IF EXISTS "sso"."client_redirect_uris";

DROP TABLE IF EXISTS "sso"."client_audiences";

DROP TABLE IF EXISTS "sso"."clients";

DROP TABLE IF EXISTS "sso"."users";


DROP SCHEMA IF EXISTS "sso";

-- +goose StatementEnd
