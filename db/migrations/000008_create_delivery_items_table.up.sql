CREATE TABLE
          delivery_items (
                    id SERIAL PRIMARY KEY,
                    delivery_id INTEGER NOT NULL REFERENCES deliveries (id) ON DELETE CASCADE,
                    item_name VARCHAR(100) NOT NULL,
                    quantity VARCHAR(15) NOT NULL,
                    weight VARCHAR(10) NOT NULL
          );