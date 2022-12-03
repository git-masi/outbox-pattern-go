-- Init tables for this project once the outbox DB has been created
CREATE TYPE order_status AS ENUM ('created', 'shipped', 'backorder', 'canceled');
CREATE TYPE item_fee_description AS ENUM ('heavy item', 'large item', 'fragile item');

CREATE TABLE IF NOT EXISTS clients (
	id serial PRIMARY KEY,
	created text NOT NULL,
	last_updated text NOT NULL,
	"name" text NOT NULL
);

CREATE TABLE IF NOT EXISTS orders (
	id serial PRIMARY KEY,
	created text NOT NULL,
	last_updated text NOT NULL,
	"status" order_status NOT NULL,
	client_id int REFERENCES clients
);

CREATE TABLE IF NOT EXISTS order_fulfillment_messages (
	id serial PRIMARY KEY,
	created text NOT NULL,
	message_body jsonb NOT NULL
);

CREATE TABLE IF NOT EXISTS items (
	id serial PRIMARY KEY,
	created text NOT NULL,
	last_updated text NOT NULL,
	"name" text NOT NULL,
	"description" text NOT NULL,
	price int NOT NULL,
	client_id int REFERENCES clients
);

CREATE TABLE IF NOT EXISTS order_items (
	order_id int REFERENCES orders,
	item_id int REFERENCES items
);

CREATE TABLE IF NOT EXISTS fulfillment_fees (
	id serial PRIMARY KEY,
	created text NOT NULL,
	last_updated text NOT NULL,
	rate int NOT NULL,
	client_id int REFERENCES clients
);

CREATE TABLE IF NOT EXISTS item_fees (
	id serial PRIMARY KEY,
	created text NOT NULL,
	last_updated text NOT NULL,
	rate int NOT NULL,
	fee_description item_fee_description NOT NULL,
	item_id int REFERENCES items,
	client_id int REFERENCES clients
);

CREATE TABLE IF NOT EXISTS charges (
	id serial PRIMARY KEY,
	created text NOT NULL,
	amount int NOT NULL,
	order_id int REFERENCES orders,
	client_id int REFERENCES clients
);