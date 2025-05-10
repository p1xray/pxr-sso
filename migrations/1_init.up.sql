CREATE TABLE IF NOT EXISTS users
(
    id INTEGER PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    fio VARCHAR(255) NOT NULL,
    date_of_birth TIMESTAMP,
    gender INTEGER,
    avatar_file_key VARCHAR(255),
    deleted BOOL NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_users_username_password_hash ON users (username, password_hash);

CREATE TABLE IF NOT EXISTS clients
(
    id INTEGER PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(255) NOT NULL UNIQUE,
    secret_key VARCHAR(255) NOT NULL,
    deleted BOOL NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_clients_code ON clients (code);

CREATE TABLE IF NOT EXISTS user_clients
(
    id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL,
    client_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id)  REFERENCES users (id),
    FOREIGN KEY (client_id)  REFERENCES clients (id)
);
CREATE INDEX IF NOT EXISTS idx_user_clients_user_id ON user_clients (user_id);
CREATE INDEX IF NOT EXISTS idx_user_clients_client_id ON user_clients (client_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_clients_user_id_client_id ON user_clients(user_id, client_id);

CREATE TABLE IF NOT EXISTS sessions
(
    id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL,
    refresh_token VARCHAR(255) NOT NULL UNIQUE,
    user_agent VARCHAR(255) NOT NULL,
    fingerprint VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id)  REFERENCES users (id)
);
CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions (user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_refresh_token ON sessions (refresh_token);

CREATE TABLE IF NOT EXISTS roles
(
    id INTEGER PRIMARY KEY,
    code VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(1000),
    active BOOL NOT NULL,
    deleted BOOL NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_roles_active ON roles (active);

CREATE TABLE IF NOT EXISTS permissions
(
    id INTEGER PRIMARY KEY,
    code VARCHAR(255) NOT NULL UNIQUE,
    description VARCHAR(1000),
    active BOOL NOT NULL,
    deleted BOOL NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_permissions_active ON permissions (active);

CREATE TABLE IF NOT EXISTS role_permissions
(
    id INTEGER PRIMARY KEY,
    role_id INTEGER NOT NULL,
    permission_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    FOREIGN KEY (role_id)  REFERENCES roles (id),
    FOREIGN KEY (permission_id)  REFERENCES permissions (id)
);
CREATE INDEX IF NOT EXISTS idx_role_permissions_role_id ON role_permissions (role_id);
CREATE INDEX IF NOT EXISTS idx_role_permissions_permission_id ON role_permissions (permission_id);

CREATE TABLE IF NOT EXISTS user_roles
(
    id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL,
    role_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id)  REFERENCES users (id),
    FOREIGN KEY (role_id)  REFERENCES roles (id)
);
CREATE INDEX IF NOT EXISTS idx_user_roles_user_id ON user_roles (user_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_role_id ON user_roles (role_id);
