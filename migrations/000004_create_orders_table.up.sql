CREATE TABLE IF NOT EXISTS orders
(
    uid                VARCHAR(64) PRIMARY KEY,
    track_number       VARCHAR(64) NOT NULL,
    entry              VARCHAR(64) NOT NULL,
    locale             VARCHAR(6)  NOT NULL,
    internal_signature VARCHAR(64) NOT NULL,
    customer_id        VARCHAR(64) NOT NULL,
    delivery_service   VARCHAR(64) NOT NULL,
    shard_key           VARCHAR(64) NOT NULL,
    sm_id              INT         NOT NULL,
    date_created       TIMESTAMP   NOT NULL,
    oof_shard          VARCHAR(64) NOT NULL
    );