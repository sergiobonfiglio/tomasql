CREATE TABLE IF NOT EXISTS account
(
    id         bigint GENERATED ALWAYS AS IDENTITY NOT NULL,
    uuid       char(36)                            NOT NULL,
    created_ts int                                 NOT NULL,
    CONSTRAINT pk_account PRIMARY KEY ("id"),
    CONSTRAINT uq_account_uuid UNIQUE ("uuid")
);

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
