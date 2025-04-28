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
                    transaction_id INTEGER NOT NULL UNIQUE REFERENCES transactions (ID) ON DELETE CASCADE,
                    driver_id INTEGER NOT NULL UNIQUE REFERENCES drivers (id) ON DELETE CASCADE,
                    vehicle_id INTEGER NOT NULL UNIQUE REFERENCES vehicles (id) ON DELETE CASCADE load_type VARCHAR(30) NOT NULL,
                    load VARCHAR(30) NOT NULL,
                    quantity INTEGER NOT NULL,
                    weight INTEGER NOT NULL,
                    pickup_location TEXT NOT NULL,
                    destination TEXT NOT NULL,
                    delivery_date TIMESTAMP NOT NULL,
                    delivery_deadline_date TIMESTAMP NOT NULL,
                    delivery_status delivery_status DEFAULT 'menunggu persetujuan',
                    delivery_cost NUMERIC(15, 2) NOT NULL,
                    approved_at TIMESTAMP,
                    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
          )