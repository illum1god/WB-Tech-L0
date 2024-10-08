CREATE TABLE IF NOT EXISTS deliveries
(
    id        SERIAL PRIMARY KEY,
    order_uid VARCHAR(64) UNIQUE,
    name      VARCHAR(64)  NOT NULL,
    phone     VARCHAR(16)  NOT NULL,
    zip       VARCHAR(255) NOT NULL,
    city      VARCHAR(255) NOT NULL,
    address   VARCHAR(255) NOT NULL,
    region    VARCHAR(255) NOT NULL,
    email     VARCHAR(255) NOT NULL
    );