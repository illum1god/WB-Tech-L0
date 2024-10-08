CREATE TABLE IF NOT EXISTS items
(
    id           SERIAL PRIMARY KEY,
    order_uid    VARCHAR(64) NOT NULL,
    chrt_id      INT         NOT NULL,
    track_number VARCHAR(64) NOT NULL,
    price        INT         NOT NULL,
    rid          VARCHAR(64) NOT NULL,
    name         VARCHAR(64) NOT NULL,
    sale         INT         NOT NULL,
    size         VARCHAR(64) NOT NULL,
    total_price  INT         NOT NULL,
    nm_id        INT         NOT NULL,
    brand        VARCHAR(64) NOT NULL,
    status       INT         NOT NULL
    );