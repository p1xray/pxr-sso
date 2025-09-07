CREATE TABLE IF NOT EXISTS client_default_roles
(
    id INTEGER PRIMARY KEY,
    client_id INTEGER NOT NULL,
    role_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    FOREIGN KEY (client_id)  REFERENCES clients (id),
    FOREIGN KEY (role_id)  REFERENCES roles (id)
);
CREATE INDEX IF NOT EXISTS idx_client_default_roles_client_id ON client_default_roles (client_id);
CREATE INDEX IF NOT EXISTS idx_client_default_roles_role_id ON client_default_roles (role_id);