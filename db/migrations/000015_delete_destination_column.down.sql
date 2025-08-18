ALTER TABLE deliveries
ADD COLUMN load_type VARCHAR(30),
ADD COLUMN delivery_deadline_date TIMESTAMP;