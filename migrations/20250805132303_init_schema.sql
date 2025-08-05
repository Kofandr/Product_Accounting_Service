-- +goose Up
CREATE TABLE categories (
                            id SERIAL PRIMARY KEY,
                            name TEXT NOT NULL,
                            description TEXT,
                            CONSTRAINT categories_name_key UNIQUE (name)
);

CREATE TABLE products (
                          id SERIAL PRIMARY KEY,
                          name VARCHAR(255) NOT NULL,
                          amount INTEGER NOT NULL CHECK (amount >= 0),
                          categoryid INTEGER NOT NULL,
                          CONSTRAINT products_name_key UNIQUE (name),
                          CONSTRAINT products_categoryid_fkey FOREIGN KEY (categoryid) REFERENCES categories(id) ON DELETE RESTRICT
);

-- +goose Down
DROP TABLE products;
DROP TABLE categories;
