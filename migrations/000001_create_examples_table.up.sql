CREATE TABLE IF NOT EXISTS Examples (
    ID        UInt64,
    Name      String,
    Timestamp DateTime,
    Value     Float64,
    Tags      Array(String)
) ENGINE = MergeTree()
ORDER BY (Timestamp, ID)
PARTITION BY toYYYYMM(Timestamp);
