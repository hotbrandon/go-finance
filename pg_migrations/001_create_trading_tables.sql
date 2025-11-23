-- purchases: one row per buy (one "position" created by each automated $10 buy)
CREATE TABLE crypto_purchases (
    id BIGSERIAL PRIMARY KEY,
    symbol TEXT NOT NULL,                  -- e.g. 'ETHUSDT' or 'ETH'
    exchange TEXT,                         -- e.g. 'binance'
    order_id TEXT,                         -- exchange order id for traceability
    quantity NUMERIC(20,8) NOT NULL,      -- crypto units purchased
    remaining_qty NUMERIC(20,8) NOT NULL, -- remaining quantity (after sells)
    fiat_invested NUMERIC(20,8) NOT NULL, -- USD amount invested (e.g. 10.00)
    buy_price NUMERIC(20,8) NOT NULL,     -- price per unit in fiat at buy time
    fee NUMERIC(20,8) DEFAULT 0,
    fee_currency TEXT,
    target_gain NUMERIC(6,5) NOT NULL DEFAULT 0.03, -- 3% = 0.03
    target_price NUMERIC(20,8) GENERATED ALWAYS AS (buy_price * (1 + target_gain)) STORED,
    status TEXT NOT NULL DEFAULT 'open' CHECK (status IN ('open','partial','closed')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    closed_at TIMESTAMPTZ,
    CONSTRAINT uq_exchange_order UNIQUE (exchange, order_id)
);

-- sells: record of sell trades that reference a purchase (allows partial sells)
CREATE TABLE crypto_sells (
    id BIGSERIAL PRIMARY KEY,
    purchase_id BIGINT NOT NULL REFERENCES crypto_purchases(id) ON DELETE CASCADE,
    exchange TEXT,
    order_id TEXT,
    quantity NUMERIC(20,8) NOT NULL,
    sell_price NUMERIC(20,8) NOT NULL,   -- price per unit in fiat at sell time
    fiat_amount NUMERIC(20,8) NOT NULL,  -- quantity * sell_price - fee
    fee NUMERIC(20,8) DEFAULT 0,
    fee_currency TEXT,
    executed_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Useful indexes
CREATE INDEX idx_purchases_symbol_status ON crypto_purchases (symbol, status);
CREATE INDEX idx_purchases_target_price ON crypto_purchases (target_price);
CREATE INDEX idx_sells_purchase_id ON crypto_sells (purchase_id);

-- Optional: partial index for open purchases (fast candidate lookup)
CREATE INDEX idx_purchases_open_symbol_target ON crypto_purchases (symbol, target_price)
    WHERE status = 'open' AND remaining_qty > 0;