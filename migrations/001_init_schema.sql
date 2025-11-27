CREATE TYPE product_type AS ENUM (
    'clothing_headwear', 'clothing_body', 'clothing_pants', 'clothing_shoes',
    'food', 'furniture', 'electronics', 'adult', 'home_goods'
);

CREATE TABLE products (
                          id BIGSERIAL PRIMARY KEY,
                          name VARCHAR(255) NOT NULL,
                          weight DECIMAL(10,3) NOT NULL,
                          unit VARCHAR(10) NOT NULL CHECK (unit IN ('g', 'kg', 'l', 'piece')),
                          color VARCHAR(50),
                          type product_type NOT NULL,
                          price DECIMAL(10,2) NOT NULL CHECK (price >= 0),
                          attributes JSONB NOT NULL DEFAULT '{}',
                          created_at TIMESTAMPTZ DEFAULT NOW(),
                          updated_at TIMESTAMPTZ DEFAULT NOW(),

                          CONSTRAINT valid_weight CHECK (
                              (unit = 'piece' AND weight >= 1) OR
                              (unit != 'piece' AND weight > 0)
                              )
);

CREATE INDEX idx_products_type ON products(type);
CREATE INDEX idx_products_price ON products(price);
CREATE INDEX idx_products_color ON products(color);
CREATE INDEX idx_products_created_at ON products(created_at);
CREATE INDEX idx_products_attributes ON products USING GIN(attributes);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_products_updated_at
    BEFORE UPDATE ON products
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();