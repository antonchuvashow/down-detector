ALTER TABLE reports
    ADD COLUMN IF NOT EXISTS latitude Float64,
    ADD COLUMN IF NOT EXISTS longitude Float64,
    ADD COLUMN IF NOT EXISTS ip TEXT,
    ADD COLUMN IF NOT EXISTS platform LowCardinality(String) DEFAULT 'unknown';

ALTER TABLE reports
    ADD INDEX IF NOT EXISTS idx_platform platform TYPE set(0) GRANULARITY 4;

ALTER TABLE reports
    ADD INDEX IF NOT EXISTS idx_geo (latitude, longitude) TYPE minmax GRANULARITY 1;

ALTER TABLE reports
    ADD INDEX IF NOT EXISTS idx_ip ip TYPE bloom_filter(0.01) GRANULARITY 1;