CREATE TABLE events (
    id          String,
    operation   String,
    detail      JSON,
    options     Nested (
        name    String,
        type    String,
        action  String
    ),
    records     Array(JSON),
    timestamp   DateTime64(9)
) ENGINE = MergeTree()
ORDER BY (id);
