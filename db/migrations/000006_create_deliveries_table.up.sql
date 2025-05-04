CREATE TYPE delivery_status AS ENUM (
          'menunggu persetujuan',
          'disetujui',
          'ditolak',
          'sedang berlangsung',
          'selesai',
          'dibatalkan'
);

CREATE TABLE
          deliveries (
                    id SERIAL PRIMARY KEY,
                    load_type VARCHAR(30),
                    load VARCHAR(30) NOT NULL,
                    quantity VARCHAR(15) NOT NULL,
                    weight VARCHAR(10) NOT NULL,
                    pickup_address TEXT NOT NULL,
                    pickup_address_lat FLOAT NOT NULL,
                    pickup_address_lang FLOAT NOT NULL,
                    destination_address TEXT NOT NULL,
                    destination_address_lat FLOAT NOT NULL,
                    destination_address_lang FLOAT NOT NULL,
                    delivery_date TIMESTAMP NOT NULL,
                    delivery_deadline_date TIMESTAMP NOT NULL,
                    delivery_status delivery_status DEFAULT 'menunggu persetujuan',
                    delivery_cost NUMERIC(15, 2) NOT NULL,
                    approved_at TIMESTAMP,
                    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                    transaction_id INTEGER NOT NULL REFERENCES transactions (ID) ON DELETE CASCADE,
                    driver_id INTEGER NOT NULL REFERENCES drivers (id) ON DELETE CASCADE,
                    vehicle_id INTEGER NOT NULL REFERENCES vehicles (id) ON DELETE CASCADE
          )