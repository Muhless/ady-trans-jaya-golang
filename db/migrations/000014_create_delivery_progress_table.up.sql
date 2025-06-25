CREATE TABLE
          delivery_progresses (
                    id SERIAL PRIMARY KEY,
                    delivery_id INT NOT NULL,
                    pickup_time TIMESTAMP,
                    destination_arrival_time TIMESTAMP,
                    receiver_name VARCHAR(50),
                    receiver_phone VARCHAR(20),
                    received_at TIMESTAMP,
                    receiver_signature_url TEXT,
                    delivery_photo_url TEXT,
                    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                    CONSTRAINT fk_delivery FOREIGN KEY (delivery_id) REFERENCES deliveries (id) ON DELETE CASCADE
          );