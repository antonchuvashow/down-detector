CREATE TABLE IF NOT EXISTS routes
(
    id  String,
    url String
)
ENGINE = MergeTree()
ORDER BY id;