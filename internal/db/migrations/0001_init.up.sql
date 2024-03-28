CREATE TABLE IF NOT EXISTS categories (
    category_id SERIAL PRIMARY KEY,
    category TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS items (
    item_id SERIAL PRIMARY KEY,
    item TEXT NOT NULL,
    price INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS categories_items (
    category_id INTEGER REFERENCES categories (category_id) ON UPDATE CASCADE ON DELETE CASCADE,
    item_id INTEGER REFERENCES items (item_id) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT categories_items_pkey PRIMARY KEY (category_id, item_id)
);

CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    password VARCHAR(250) NOT NULL,
    username TEXT NOT NULL UNIQUE
);