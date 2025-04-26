CREATE TYPE vehicle_type AS ENUM ('pick up', 'cde', 'cdd', 'fuso', 'wingbox');

CREATE TYPE vehicle_status AS ENUM ('tersedia', 'tidak tersedia');

CREATE TABLE
          vehicles (
                    ID SERIAL PRIMARY KEY,
                    name VARCHAR(30) NOT NULL,
                    type vehicle_type NOT NULL,
                    license_plate VARCHAR(8) NOT NULL,
                    capacity VARCHAR(15) NOT NULL,
                    price DECIMAL(10, 2) NOT NULL,
                    status vehicle_status DEFAULT 'tersedia',
                    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
          )