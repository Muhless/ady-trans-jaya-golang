ALTER TABLE delivery_items
ALTER COLUMN quantity TYPE VARCHAR(15)
USING quantity::VARCHAR(15);

ALTER TABLE delivery_items
DROP COLUMN unit;
