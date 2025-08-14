CREATE TABLE otp (
    id serial PRIMARY KEY,
    password VARCHAR(255),
    expires_at INTEGER,
)