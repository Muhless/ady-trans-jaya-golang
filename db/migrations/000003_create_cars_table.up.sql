CREATE TYPE car_type AS ENUM ('pick up', 'cde', 'cdd', 'fuso', 'wingbox');

CREATE TYPE car_status AS ENUM ('tersedia', 'tidak tersedia');

CREATE TABLE
          cars (
                    ID SERIAL PRIMARY KEY,
                    name VARCHAR(30) NOT NULL,
                    type car_type NOT NULL,
                    license_plat VARCHAR(8) NOT NULL,
                    price INTEGER NOT NULL,
                    status car_status NOT NULL
          )