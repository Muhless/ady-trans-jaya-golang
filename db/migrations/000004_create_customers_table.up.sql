CREATE TABLE
          customers (
                    ID SERIAL PRIMARY KEY,
                    name VARCHAR(30) NOT NULL,
                    company VARCHAR(30) NOT NULL,
                    email VARCHAR(30) NOT NULL,
                    phone VARCHAR(15) UNIQUE NOT NULL,
                    address TEXT NOT NULL
          )