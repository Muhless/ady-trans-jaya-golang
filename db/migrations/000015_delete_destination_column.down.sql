ALTER TABLE deliveries
ADD COLUMN destination_address TEXT,
ADD COLUMN destination_address_lat FLOAT,
ADD COLUMN destination_address_lang FLOAT,
ADD COLUMN load_type VARCHAR(30),
ADD COLUMN delivery_deadline_date TIMESTAMP;

ALTER TABLE deliveries
DROP COLUMN delivery_start_time,
DROP COLUMN pickup_time,
DROP COLUMN pickup_photo_url;