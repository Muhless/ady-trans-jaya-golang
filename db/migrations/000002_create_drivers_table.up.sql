CREATE TYPE driver_status AS ENUM ('tersedia', 'tidak tersedia');

CREATE TABLE
  drivers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    phone VARCHAR(15) NOT NULL UNIQUE,
    address TEXT,
    photo TEXT,
    status driver_status DEFAULT 'tersedia',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
  );