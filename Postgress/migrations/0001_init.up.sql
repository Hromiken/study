CREATE TABLE IF NOT EXISTS public.example
(
    order_uid          UUID PRIMARY KEY,
    track_number       VARCHAR(255),
    entry              VARCHAR(50),
    delivery           JSONB,
    payment            JSONB,
    items              JSONB,
    locale             VARCHAR(10),
    internal_signature VARCHAR(255),
    customer_id        VARCHAR(255),
    delivery_service   VARCHAR(255),
    shardkey           VARCHAR(10),
    sm_id              INTEGER,
    date_created       TIMESTAMP WITH TIME ZONE,
    oof_shard          VARCHAR(10)
);