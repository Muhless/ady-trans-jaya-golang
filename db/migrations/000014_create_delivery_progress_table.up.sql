CREATE TABLE
          delivery_progresses (
                    id SERIAL PRIMARY KEY,
                    delivery_id INT NOT NULL UNIQUE,
                    delivery_start_time TIMESTAMP,
                    pickup_time TIMESTAMP,
                    pickup_photo_url TEXT,
                    arrival_time TIMESTAMP,
                    arrival_photo_url TEXT,
                    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                    CONSTRAINT fk_delivery FOREIGN KEY (delivery_id) REFERENCES deliveries (id) ON DELETE CASCADE
          );