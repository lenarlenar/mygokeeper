CREATE TABLE IF NOT EXISTS records (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    type TEXT NOT NULL,
    data BYTEA NOT NULL,
    meta TEXT
);
