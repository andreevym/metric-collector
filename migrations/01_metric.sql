CREATE TABLE IF NOT EXISTS metric
(
    id   VARCHAR(50) PRIMARY KEY,
    type VARCHAR(50),
    delta integer,
    value decimal
);