CREATE TYPE order_status AS ENUM ('created', 'shipped', 'backorder', 'canceled');

CREATE TABLE IF NOT EXISTS clients (
	id bigserial PRIMARY KEY,
	created text NOT NULL,
	last_updated text NOT NULL,
	`name` text NOT NULL
)

CREATE TABLE IF NOT EXISTS orders (
	id bigserial PRIMARY KEY,
	created text NOT NULL,
	last_updated text NOT NULL,
	`status` order_status NOT NULL,
	client_id bigserial REFERENCES clients
);

CREATE TABLE IF NOT EXISTS items (
	id bigserial PRIMARY KEY,
	created text NOT NULL,
	last_updated text NOT NULL,
	`name` text NOT NULL,
	`description` text NOT NULL,
	price int NOT NULL,
	client_id bigserial REFERENCES clients
)

CREATE TABLE IF NOT EXISTS order_items (
	order_id bigserial NOT NULL,
	item_id bigserial NOT NULL,
	FOREIGN KEY (order_id) REFERENCES orders,
	FOREIGN KEY (item_id) REFERENCES items
)

CREATE TABLE IF NOT EXISTS fulfillment_fees (
	id bigserial PRIMARY KEY,
	rate int NOT NULL,
	order_threshold int NOT NULL,
	`description` text NOT NULL,
	client_id bigserial REFERENCES clients
)

CREATE TABLE IF NOT EXISTS item_fees (
	id bigserial PRIMARY KEY,
	rate int NOT NULL,
	`description` text NOT NULL,
	FOREIGN KEY (item_id) REFERENCES items
)
