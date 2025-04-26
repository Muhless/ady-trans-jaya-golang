CREATE TYPE volume_unit as ENUM ('kg', 'ton', 'm3', 'liter');

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
                    ID SERIAL PRIMARY KEY,
                    transaction_id INTEGER NOT NULL REFERENCES transactions (id) ON DELETE CASCADE,
                    driver_id INTEGER NOT NULL REFERENCES drivers (id) ON DELETE CASCADE,
                    vehicle_id INTEGER NOT NULL REFERENCES vehicles (id) ON DELETE CASCADE,
                    content VARCHAR(30) NOT NULL,
                    volume NUMERIC(10, 2) NOT NULL,
                    volume_unit volume_unit NOT NULL,
                    address_from TEXT NOT NULL,
                    address_to TEXT NOT NULL,
                    delivery_date TIMESTAMP NOT NULL,
                    delivery_deadline_date TIMESTAMP NOT NULL,
                    delivery_status delivery_status DEFAULT 'menunggu persetujuan',
                    total NUMERIC(15, 2) NOT NULL,
                    approved_at TIMESTAMP,
                    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
          )