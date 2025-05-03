CREATE TABLE provinces (
    id VARCHAR PRIMARY KEY,
    name VARCHAR NOT NULL
);

CREATE TABLE cities (
    id VARCHAR PRIMARY KEY,
    name VARCHAR NOT NULL,
    province_id VARCHAR NOT NULL REFERENCES provinces(id)
);

CREATE TABLE districts (
    id VARCHAR PRIMARY KEY,
    name VARCHAR NOT NULL,
    city_id VARCHAR NOT NULL REFERENCES cities(id)
);

CREATE TABLE subdistricts (
    id VARCHAR PRIMARY KEY,
    name VARCHAR NOT NULL,
    district_id VARCHAR NOT NULL REFERENCES districts(id)
);
