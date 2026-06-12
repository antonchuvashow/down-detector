CREATE TABLE IF NOT EXISTS reports
(
    time        DateTime64(3)  NOT NULL,
    route_id    String         NOT NULL,
    success     Bool           NOT NULL,
    error_types Array(String)  NOT NULL,
    source      String         NOT NULL,
    latency_ms  Int64          NOT NULL
)
ENGINE = MergeTree()
ORDER BY (time, source)
PARTITION BY toYYYYMM(time);