ALTER TABLE deliveries 
  DROP COLUMN destination_address,
  DROP COLUMN destination_address_lat,
  DROP COLUMN destination_address_lang,
  DROP COLUMN delivery_deadline_date,
  DROP COLUMN load_type;

ALTER TABLE deliveries
  ADD COLUMN delivery_start_time TIMESTAMP,
  ADD COLUMN pickup_time TIMESTAMP,
  ADD COLUMN pickup_photo_url TEXT;
