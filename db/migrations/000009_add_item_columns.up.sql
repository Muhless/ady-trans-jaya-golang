ALTER TABLE deliveries
ADD COLUMN delivery_code VARCHAR(25) UNIQUE,
ADD COLUMN total_weight INTEGER DEFAULT 0;