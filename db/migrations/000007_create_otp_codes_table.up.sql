CREATE TABLE
          otp_codes (
                    id SERIAL PRIMARY KEY,
                    phone VARCHAR(20) NOT NULL,
                    otp_code VARCHAR(6) NOT NULL,
                    otp_expired_at TIMESTAMP NOT NULL,
                    is_used BOOLEAN DEFAULT FALSE,
                    created_at TIMESTAMP DEFAULT NOW ()
          );