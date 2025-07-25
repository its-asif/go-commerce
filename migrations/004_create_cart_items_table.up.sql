CREATE TABLE IF NOT EXISTS cart_items (
    user_id INT REFERENCES users(id),
    product_id INT REFERENCES products(id),
    quantity INT NOT NULL,
    price NUMERIC(10,2) NOT NULL,
    added_at TIMESTAMP DEFAULT now(),
    PRIMARY KEY (user_id, product_id)
);
