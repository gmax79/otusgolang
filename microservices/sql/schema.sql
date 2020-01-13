CREATE TABLE IF NOT EXISTS events (
    timer       TIMESTAMP NOT NULL,
    information VARCHAR(255) NOT NULL,
    CONSTRAINT ti UNIQUE (timer, information)
);
