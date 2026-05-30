-- =============================================================
-- Schema: Inventaris Manajemen
-- =============================================================

CREATE EXTENSION IF NOT EXISTS "pgcrypto";


-- =============================================================
-- USERS
-- =============================================================

CREATE TABLE users (
    id            UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    name          VARCHAR(100) NOT NULL,
    username      VARCHAR(50)  NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,  -- bcrypt hash
    created_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);


-- =============================================================
-- CATEGORIES
-- =============================================================

CREATE TABLE categories (
    id   UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(80) NOT NULL UNIQUE
);


-- =============================================================
-- SUPPLIERS
-- =============================================================

CREATE TABLE suppliers (
    id      UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    name    VARCHAR(100) NOT NULL,
    contact VARCHAR(100),  -- nomor HP / email
    address TEXT
);


-- =============================================================
-- PRODUCTS
-- =============================================================

CREATE TABLE products (
    id          UUID           PRIMARY KEY DEFAULT gen_random_uuid(),
    category_id UUID           REFERENCES categories(id) ON DELETE SET NULL,
    name        VARCHAR(150)   NOT NULL,
    sku         VARCHAR(80)    NOT NULL UNIQUE,
    unit        VARCHAR(20)    NOT NULL DEFAULT 'pcs',  -- pcs, kg, liter, dll
    price_buy   NUMERIC(15, 2) NOT NULL DEFAULT 0,
    price_sell  NUMERIC(15, 2) NOT NULL DEFAULT 0,
    stock       INTEGER        NOT NULL DEFAULT 0,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Junction table: produk <-> supplier (many-to-many)
CREATE TABLE product_suppliers (
    product_id  UUID NOT NULL REFERENCES products(id)  ON DELETE CASCADE,
    supplier_id UUID NOT NULL REFERENCES suppliers(id) ON DELETE CASCADE,
    PRIMARY KEY (product_id, supplier_id)
);


-- =============================================================
-- STOCK MOVEMENTS
-- Tercatat otomatis setiap kali stok produk berubah.
-- User tidak perlu tau tabel ini ada — mereka cukup
-- interact lewat halaman produk.
-- =============================================================

CREATE TABLE stock_movements (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID        NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    changed_by UUID        NOT NULL REFERENCES users(id),
    type       VARCHAR(20) NOT NULL,  -- 'in', 'out', 'adjustment'
    qty        INTEGER     NOT NULL,  -- positif = stok naik, negatif = stok turun
    note       TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_movement_type CHECK (type IN ('in', 'out', 'adjustment'))
);


-- =============================================================
-- INDEXES
-- =============================================================

CREATE INDEX idx_products_category    ON products(category_id);
CREATE INDEX idx_products_sku         ON products(sku);
CREATE INDEX idx_stock_mov_product    ON stock_movements(product_id);
CREATE INDEX idx_stock_mov_created_at ON stock_movements(created_at);
