DROP INDEX IF EXISTS idx_user_roles_user_id;
DROP INDEX IF EXISTS idx_user_roles_role_id;
DROP TABLE IF EXISTS user_roles;

DROP INDEX IF EXISTS idx_role_permissions_role_id;
DROP INDEX IF EXISTS idx_role_permissions_permission_id;
DROP TABLE IF EXISTS role_permissions;

DROP INDEX IF EXISTS idx_permissions_active;
DROP TABLE IF EXISTS permissions;

DROP INDEX IF EXISTS idx_roles_active;
DROP TABLE IF EXISTS roles;

DROP INDEX IF EXISTS idx_sessions_user_id;
DROP INDEX IF EXISTS idx_sessions_refresh_token;
DROP TABLE IF EXISTS sessions;

DROP INDEX IF EXISTS idx_user_clients_user_id;
DROP INDEX IF EXISTS idx_user_clients_client_id;
DROP TABLE IF EXISTS user_clients;

DROP INDEX IF EXISTS idx_clients_code;
DROP TABLE IF EXISTS clients;

DROP INDEX IF EXISTS idx_users_username_password_hash;
DROP TABLE IF EXISTS users;
