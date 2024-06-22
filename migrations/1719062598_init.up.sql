CREATE TABLE sensors_data (
    t               TIMESTAMPTZ NOT NULL,
    sensor          VARCHAR(16) NOT NULL,
    value           FLOAT NOT NULL,
    CONSTRAINT sensors_data_pkey PRIMARY KEY (t, sensor)
);
SELECT create_hypertable('sensors_data', 't');
