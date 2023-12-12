
-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS warehouse (
	id serial PRIMARY KEY,
	name text NULL,
	available bool NOT NULL DEFAULT true
);

CREATE TABLE IF NOT EXISTS product (
    code serial PRIMARY KEY,
    name text,
    size numeric,
    value integer,
    stock_id int
);

CREATE TABLE IF NOT EXISTS reserver (
	id serial PRIMARY KEY,
	product_code int NULL,
	stock_id int NOT NULL,
	value int NULL
);

-- +goose StatementEnd

-- -- +goose down
-- -- +goose StatementBegin
-- DROP TABLE PRODUCT;

-- DROP TABLE warehouse;

-- DROP TABLE resever;

-- -- +goose StatementEnd