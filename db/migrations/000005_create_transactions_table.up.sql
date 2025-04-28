CREATE TYPE payment_status AS ENUM ('pending', 'lunas', 'gagal');

CREATE TYPE transaction_status AS ENUM (
          'menunggu persetujuan',
          'sedang berlangsung',
          'selesai',
          'dibatalkan'
);

CREATE TABLE
          transactions (
                    ID SERIAL PRIMARY KEY,
                    customer_id INTEGER NOT NULL REFERENCES customers (id) ON DELETE CASCADE,
                    total_delivery INTEGER,
                    cost NUMERIC(15, 2),
                    payment_deadline TIMESTAMP,
                    down_payment NUMERIC(15, 2),
                    down_payment_status payment_status,
                    down_payment_time TIMESTAMP,
                    down_payment_proof TEXT,
                    full_payment NUMERIC(15, 2),
                    full_payment_status payment_status,
                    full_payment_time TIMESTAMP,
                    full_payment_proof TEXT,
                    transaction_status transaction_status,
                    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
          )