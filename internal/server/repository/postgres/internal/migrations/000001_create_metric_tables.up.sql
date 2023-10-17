CREATE TABLE gauges
(
    name  TEXT PRIMARY KEY,
    value FLOAT8 NOT NULL
);
CREATE TABLE counters
(
    name  TEXT PRIMARY KEY,
    value BIGINT NOT NULL
);
