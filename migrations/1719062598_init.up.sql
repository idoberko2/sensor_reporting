CREATE TABLE sensors_data (
    t               TIMESTAMPTZ NOT NULL PRIMARY KEY,
    sensor          VARCHAR NOT NULL,
    value           FLOAT NOT NULL
);
SELECT create_hypertable('sensors_data', 't');
