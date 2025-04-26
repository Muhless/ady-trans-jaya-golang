CREATE TABLE
          customers (
                    ID SERIAL PRIMARY KEY,
                    name VARCHAR(30) NOT NULL,
                    company VARCHAR(30) NOT NULL,
                    email VARCHAR(30),
                    phone VARCHAR(15) UNIQUE NOT NULL,
                    address TEXT NOT NULL,
                    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
          )