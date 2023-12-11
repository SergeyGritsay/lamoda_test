
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

CREATE TABLE IF NOT EXISTS res_cen (
	id serial PRIMARY KEY,
	product_code int4 NULL,
	stock_id uuid NOT NULL,
	value int4 NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE PRODUCT;

DROP TABLE warehouse;

DROP TABLE res_cen;

-- +goose StatementEnd