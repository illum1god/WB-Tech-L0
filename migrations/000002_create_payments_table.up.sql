CREATE TABLE IF NOT EXISTS payments
(
    id            SERIAL PRIMARY KEY,
    order_uid     VARCHAR(64) UNIQUE,
    transaction   VARCHAR(64) NOT NULL,
    request_id    VARCHAR(64) NOT NULL,
    currency      VARCHAR(6)  NOT NULL,
    provider      VARCHAR(64) NOT NULL,
    amount        INT         NOT NULL,
    payment_dt    INT         NOT NULL,
    bank          VARCHAR(64) NOT NULL,
    delivery_cost INT         NOT NULL,
    goods_total   INT         NOT NULL,
    custom_fee    INT         NOT NULL
    );