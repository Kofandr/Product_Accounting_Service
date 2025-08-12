
-- +goose Up
ALTER TABLE products
    RENAME COLUMN categoryid TO category_id;

-- Обновляем FK констрейнт
ALTER TABLE products
DROP CONSTRAINT products_categoryid_fkey;

ALTER TABLE products
    ADD CONSTRAINT products_category_id_fkey
        FOREIGN KEY (category_id) REFERENCES categories(id);

-- +goose Down
ALTER TABLE products
    RENAME COLUMN category_id TO categoryid;

ALTER TABLE products
DROP CONSTRAINT products_category_id_fkey;

ALTER TABLE products
    ADD CONSTRAINT products_categoryid_fkey
        FOREIGN KEY (categoryid) REFERENCES categories(id);
