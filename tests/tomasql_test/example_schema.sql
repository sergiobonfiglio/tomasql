-- Simulate ENUM with CHECK constraint
CREATE TABLE account (
    id BIGINT GENERATED ALWAYS AS IDENTITY NOT NULL,
    uuid CHAR(36) NOT NULL,
    created_ts INTEGER NOT NULL,
    CONSTRAINT pk_account PRIMARY KEY (id),
    CONSTRAINT uq_account_uuid UNIQUE (uuid)
);
CREATE TABLE config (
    id BIGINT GENERATED ALWAYS AS IDENTITY NOT NULL,
    uuid CHAR(36) NOT NULL,
    account_id BIGINT NOT NULL,
    created_ts INTEGER NOT NULL,
    archived_ts INTEGER,
    CONSTRAINT pk_config PRIMARY KEY (id),
    CONSTRAINT uq_config_uuid UNIQUE (uuid),
    CONSTRAINT fk_config_account_id FOREIGN KEY (account_id) REFERENCES account (id)
);
CREATE TABLE shopping_cart (
    id BIGINT GENERATED ALWAYS AS IDENTITY NOT NULL,
    uuid CHAR(36) NOT NULL,
    owner_id BIGINT NOT NULL,
    created_ts INTEGER NOT NULL,
    archived_ts INTEGER,
    CONSTRAINT pk_shopping_cart PRIMARY KEY (id),
    CONSTRAINT uq_shopping_cart_uuid UNIQUE (uuid),
    CONSTRAINT fk_shopping_cart_owner_id FOREIGN KEY (owner_id) REFERENCES account (id)
);