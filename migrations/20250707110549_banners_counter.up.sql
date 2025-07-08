--
CREATE TABLE IF NOT EXISTS banners_counter(
    banner_id bigint NOT NULL,
    ts timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    v bigint NOT NULL DEFAULT 0,
    PRIMARY KEY (banner_id, ts)
)
PARTITION BY RANGE (ts);

--
CREATE INDEX IF NOT EXISTS idx_banners_counter_banner_id ON banners_counter(banner_id);

CREATE INDEX IF NOT EXISTS idx_banners_counter_ts ON banners_counter(ts);

CREATE INDEX IF NOT EXISTS idx_banners_counter_v ON banners_counter(v)
WHERE
    v > 0;

--
CREATE TABLE IF NOT EXISTS banners_counter_2025_07 PARTITION OF banners_counter
FOR VALUES FROM ('2025-07-01') TO ('2025-08-01');

--
CREATE TABLE IF NOT EXISTS banners_counter_default PARTITION OF banners_counter DEFAULT;

