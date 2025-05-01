CREATE TABLE
          transactions (
                    id SERIAL PRIMARY KEY,
                    total_delivery INTEGER,
                    cost NUMERIC(15, 2),
                    payment_deadline TIMESTAMP,
                    down_payment NUMERIC(15, 2),
                    down_payment_status VARCHAR(20),
                    down_payment_time TIMESTAMP,
                    full_payment NUMERIC(15, 2),
                    full_payment_status VARCHAR(20),
                    full_payment_time TIMESTAMP,
                    transaction_status VARCHAR(20),
                    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                    customer_id INTEGER NOT NULL REFERENCES customers (id) ON DELETE CASCADE
          )