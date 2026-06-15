-- +goose Up
-- +goose StatementBegin

CREATE TABLE "sso"."sessions" (
    "id" bigserial PRIMARY KEY,
    "client_id" bigint NOT NULL,
    "user_id" bigint NOT NULL,
    "code" varchar(255) NOT NULL,
    "auth_time" timestamp NOT NULL DEFAULT (now()),
    "identity_provider" varchar(255) NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT (now()),
    "updated_at" timestamp NOT NULL DEFAULT (now())
);

CREATE INDEX "idx_sessions_client_id" ON "sso"."sessions" ("client_id");
CREATE INDEX "idx_sessions_user_id" ON "sso"."sessions" ("user_id");
CREATE UNIQUE INDEX "uidx_sessions_code" ON "sso"."sessions" ("code");

ALTER TABLE "sso"."sessions" ADD CONSTRAINT "fk_sessions_client_id" FOREIGN KEY ("client_id") REFERENCES "sso"."clients" ("id") DEFERRABLE INITIALLY IMMEDIATE;
ALTER TABLE "sso"."sessions" ADD CONSTRAINT "fk_sessions_user_id" FOREIGN KEY ("user_id") REFERENCES "sso"."users" ("id") DEFERRABLE INITIALLY IMMEDIATE;


CREATE TABLE "sso"."session_granted_scope_links" (
    "id" bigserial PRIMARY KEY,
    "session_id" bigint NOT NULL,
    "scope_id" bigint NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT (now()),
    "updated_at" timestamp NOT NULL DEFAULT (now())
);

CREATE INDEX "idx_session_granted_scope_links_session_id" ON "sso"."session_granted_scope_links" ("session_id");
CREATE INDEX "idx_session_granted_scope_links_scope_id" ON "sso"."session_granted_scope_links" ("scope_id");
CREATE UNIQUE INDEX "uidx_session_granted_scope_links_session_id_scope_id" ON "sso"."session_granted_scope_links" ("session_id", "scope_id");

ALTER TABLE "sso"."session_granted_scope_links" ADD CONSTRAINT "fk_session_granted_scope_links_session_id" FOREIGN KEY ("session_id") REFERENCES "sso"."sessions" ("id") DEFERRABLE INITIALLY IMMEDIATE;
ALTER TABLE "sso"."session_granted_scope_links" ADD CONSTRAINT "fk_session_granted_scope_links_scope_id" FOREIGN KEY ("scope_id") REFERENCES "sso"."scopes" ("id") DEFERRABLE INITIALLY IMMEDIATE;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE "sso"."session_granted_scope_links" DROP CONSTRAINT IF EXISTS "fk_session_granted_scope_links_scope_id";
ALTER TABLE "sso"."session_granted_scope_links" DROP CONSTRAINT IF EXISTS "fk_session_granted_scope_links_session_id";

DROP INDEX IF EXISTS "uidx_session_granted_scope_links_session_id_scope_id";
DROP INDEX IF EXISTS "idx_session_granted_scope_links_scope_id";
DROP INDEX IF EXISTS "idx_session_granted_scope_links_session_id";

DROP TABLE IF EXISTS "sso"."session_granted_scope_links";


ALTER TABLE "sso"."sessions" DROP CONSTRAINT IF EXISTS "fk_sessions_user_id";
ALTER TABLE "sso"."sessions" DROP CONSTRAINT IF EXISTS "fk_sessions_client_id";

DROP INDEX IF EXISTS "uidx_sessions_code";
DROP INDEX IF EXISTS "idx_sessions_user_id";
DROP INDEX IF EXISTS "idx_sessions_client_id";

DROP TABLE IF EXISTS "sso"."sessions";

-- +goose StatementEnd
