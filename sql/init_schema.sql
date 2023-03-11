CREATE SCHEMA sales;
SET search_path TO sales;

CREATE TYPE products_status AS ENUM ( 'out_of_stock', 'in_stock', 'running_low' );

CREATE TYPE order_status AS ENUM ( 'pending', 'processing', 'shipped', 'delivered', 'canceled' );


CREATE TABLE merchants
(
    id            SERIAL PRIMARY KEY,
    merchant_name VARCHAR(255),
    phone         INT UNIQUE          NOT NULL,
    password      VARCHAR(255)        NOT NULL,
    email         VARCHAR(255) UNIQUE NOT NULL,
    created_at    TIMESTAMP WITHOUT TIME ZONE DEFAULT now()
);

CREATE TABLE buyers
(
    id         SERIAL PRIMARY KEY,
    full_name  VARCHAR(255),
    phone      INT UNIQUE          NOT NULL,
    email      VARCHAR(255) UNIQUE NOT NULL,
    password   VARCHAR(255)        NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT now()
);

CREATE TABLE products
(
    id          SERIAL PRIMARY KEY,
    merchant_id SERIAL NOT NULL references merchants (id) ON DELETE CASCADE,
    name        VARCHAR(255),
    description TEXT,
    price       BIGINT CHECK (price > 0),
    status      products_status,
    created_at  TIMESTAMP WITHOUT TIME ZONE DEFAULT now()
);

CREATE INDEX IF NOT EXISTS product_id ON products USING BTREE (id);
CREATE INDEX IF NOT EXISTS product_status ON products (merchant_id, status);
CREATE INDEX IF NOT EXISTS product_name ON products USING GIN (to_tsvector('simple', name));


CREATE TABLE orders
(
    id         SERIAL PRIMARY KEY,
    buyer_id   INT NOT NULL REFERENCES buyers (id),
    status     order_status,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT now(),
    deleted_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL
);

CREATE TABLE order_items
(
    order_id   SERIAL REFERENCES orders (id) ON DELETE CASCADE,
    product_id SERIAL REFERENCES products (id) ON DELETE CASCADE,
    quantity   INT DEFAULT 1
);
