CREATE TABLE IF NOT EXISTS client_audiences
(
    id INTEGER PRIMARY KEY,
    client_id INTEGER NOT NULL,
    url VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    FOREIGN KEY (client_id)  REFERENCES clients (id)
);
CREATE INDEX IF NOT EXISTS idx_client_audiences_client_id ON client_audiences (client_id);