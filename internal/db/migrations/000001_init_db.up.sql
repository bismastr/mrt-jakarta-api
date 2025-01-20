CREATE TABLE stations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(250)
);

CREATE TABLE lines (
    id SERIAL PRIMARY KEY,
    stations_id_start INTEGER REFERENCES stations (id),
    stations_id_end INTEGER REFERENCES stations (id)
);

CREATE TABLE schedules (
    id SERIAL PRIMARY KEY,
    line_id INTEGER REFERENCES lines (id),
    time TIME
);