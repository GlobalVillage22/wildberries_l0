CREATE TABLE orders (
    "order_uid" text NOT NULL PRIMARY KEY,
    "track_number" text NOT NULL,
    "entry" text NOT NULL,
    "delivery" jsonb NOT NULL,
    "payment" jsonb NOT NULL,
    "items" jsonb NOT NULL,
    "locale" text NOT NULL,
    "internal_signature" text NOT NULL,
    "customer_id" text NOT NULL,
    "delivery_service" text NOT NULL,
    "shardkey" text NOT NULL,
    "sm_id" integer NOT NULL,
    "date_created" text NOT NULL,
    "oof_shard" text NOT NULL
);