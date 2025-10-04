-- Sample schema for testing TomaSQL library
CREATE TYPE account_type AS ENUM (
    'regular',
    'premium',
    'vip'
    );

-- ACCOUNT
CREATE TABLE IF NOT EXISTS account
(
    id         bigint GENERATED ALWAYS AS IDENTITY NOT NULL,
    uuid       char(36)                            NOT NULL,
    type       account_type                        NOT NULL,
    created_ts int                                 NOT NULL,
    CONSTRAINT pk_account PRIMARY KEY ("id"),
    CONSTRAINT uq_account_uuid UNIQUE ("uuid")
);

-- CONFIG
CREATE TABLE IF NOT EXISTS config
(
    id                                      bigint GENERATED ALWAYS AS IDENTITY NOT NULL,
    uuid                                    char(36)                            NOT NULL,
    account_id                              bigint                              NOT NULL,
    created_ts                              int                                 NOT NULL,
    archived_ts                             int,
    CONSTRAINT pk_config PRIMARY KEY ("id"),
    CONSTRAINT uq_config_uuid UNIQUE ("uuid"),
    CONSTRAINT fk_config_account_id FOREIGN KEY ("account_id") REFERENCES account ("id")
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_config_archived ON config (LOWER("uuid"), "archived_ts") NULLS NOT DISTINCT;

CREATE TABLE IF NOT EXISTS shopping_cart
(
    id          bigint GENERATED ALWAYS AS IDENTITY NOT NULL,
    uuid        char(36)                            NOT NULL,
    owner_id    bigint                              NOT NULL,
    created_ts  int                                 NOT NULL,
    archived_ts int,
    CONSTRAINT pk_shopping_cart PRIMARY KEY ("id"),
    CONSTRAINT uq_shopping_cart_uuid UNIQUE ("uuid"),
    CONSTRAINT fk_shopping_cart_owner_id FOREIGN KEY ("owner_id") REFERENCES account ("id")
);